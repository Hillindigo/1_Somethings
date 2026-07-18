package article_api

import (
	"blogx_server/common"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middleware"
	"blogx_server/models"
	"github.com/gin-gonic/gin"
)

type ArticleRecommendResponse struct {
	ID        uint   `json:"id" gorm:"column:id"`
	Title     string `json:"title" gorm:"column:title"`
	LookCount int    `json:"lookCount" gorm:"column:lookCount"`
}

func (ArticleApi) ArticleRecommendView(c *gin.Context) {
	cr := middleware.GetBind[common.PageInfo](c)

	var list = make([]ArticleRecommendResponse, 0)
	global.DB.Model(models.ArticleModel{}).
		Order("look_count desc").
		Where("date(created_at) = date(now())").
		Limit(cr.Limit).Select("id", "title", "look_count").Scan(&list)

	res.OkWithList(list, len(list), c)
}
