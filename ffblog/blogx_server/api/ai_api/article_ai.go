package ai_api

import (
	"blogx_server/common"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middleware"
	"blogx_server/models"
	"blogx_server/service/ai_service"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
	"strings"
)

type ArticleAiRequest struct {
	Content string `form:"content" binding:"required"`
}

func (AiApi) ArticleAiView(c *gin.Context) {
	cr := middleware.GetBind[ArticleAiRequest](c)
	if !global.Config.Ai.Enable {
		res.SSEFail("站点未启用ai功能", c)
		return
	}
	var content string
	if global.ESClient == nil {
		// 服务降级
		list, _, _ := common.ListQuery(models.ArticleModel{}, common.Options{
			Likes: []string{"title", "abstract"},
			PageInfo: common.PageInfo{
				Page:  1,
				Limit: 10,
			},
		})
		byteData, _ := json.Marshal(list)
		content = string(byteData)
	} else {
		// 查这个内容关联的文章列表
		query := elastic.NewBoolQuery()
		query.Should(
			elastic.NewMatchQuery("title", cr.Content),
			elastic.NewMatchQuery("abstract", cr.Content),
			elastic.NewMatchQuery("content", cr.Content),
		)
		// 只能查发布的文章
		query.Must(elastic.NewTermQuery("status", 3))
		result, err := global.ESClient.
			Search(models.ArticleModel{}.Index()).
			Query(query).
			From(1).
			Size(1).
			Do(context.Background())
		if err != nil {
			source, _ := query.Source()
			byteData, _ := json.Marshal(source)
			logrus.Errorf("查询失败 %s \n %s", err, string(byteData))
			res.SSEFail("查询失败", c)
			return
		}
		var list []string
		for _, hit := range result.Hits.Hits {
			list = append(list, string(hit.Source))
		}
		content = "[" + strings.Join(list, ",") + "]"
	}
	msgChan, err := ai_service.ChatStream(cr.Content, content)
	if err != nil {
		res.SSEFail("ai分析失败", c)
		return
	}
	for s := range msgChan {
		res.SSEOk(s, c)
	}
}
