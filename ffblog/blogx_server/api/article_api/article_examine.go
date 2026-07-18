package article_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middleware"
	"blogx_server/models"
	"blogx_server/models/enum"
	"blogx_server/service/message_service"
	"fmt"
	"github.com/gin-gonic/gin"
)

type ArticleExamineRequest struct {
	ArticleID uint               `json:"articleID" binding:"required"`
	Status    enum.ArticleStatus `json:"status" binding:"required,oneof=3 4"`
	Msg       string             `json:"msg"` // 为4的时候，传递进来
}

func (ArticleApi) ArticleExamineView(c *gin.Context) {
	cr := middleware.GetBind[ArticleExamineRequest](c)

	var article models.ArticleModel
	err := global.DB.Take(&article, cr.ArticleID).Error
	if err != nil {
		res.FailWithMsg("文章不存在", c)
		return
	}

	global.DB.Model(&article).Update("status", cr.Status)

	switch cr.Status {
	case 3: // 审核成功
		message_service.InsertSystemMessage(article.UserID, "管理员审核了你的文章", "审核成功", article.Title, fmt.Sprintf("/article/%d", article.ID))
	case 4: // 审核失败
		message_service.InsertSystemMessage(article.UserID, "管理员审核了你的文章", fmt.Sprintf("审核失败 失败原因：%s", cr.Msg), "", "")
	}

	res.OkWithMsg("审核成功", c)
}
