package article_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middleware"
	"blogx_server/models"
	"blogx_server/models/enum"
	"blogx_server/service/message_service"
	"blogx_server/service/redis_service/redis_article"
	"blogx_server/utils/jwts"
	"github.com/gin-gonic/gin"
)

func (ArticleApi) ArticleDiggView(c *gin.Context) {
	cr := middleware.GetBind[models.IDRequest](c)

	var article models.ArticleModel
	err := global.DB.Take(&article, "status = ? and id = ?", enum.ArticleStatusPublished, cr.ID).Error
	if err != nil {
		res.FailWithMsg("文章不存在", c)
		return
	}

	claims := jwts.GetClaims(c)

	// 查一下之前有没有点过
	var userDiggArticle models.ArticleDiggModel
	err = global.DB.Take(&userDiggArticle, "user_id = ? and article_id = ?", claims.UserID, article.ID).Error
	if err != nil {
		// 点赞
		model := models.ArticleDiggModel{
			UserID:    claims.UserID,
			ArticleID: cr.ID,
		}
		err = global.DB.Create(&model).Error
		if err != nil {
			res.FailWithMsg("点赞失败", c)
			return
		}
		redis_article.SetCacheDigg(cr.ID, true)
		// 给文章拥有者发消息
		message_service.InsertDiggArticleMessage(model)
		res.OkWithMsg("点赞成功", c)
		return
	}
	// 取消点赞
	global.DB.Delete(&userDiggArticle)
	res.OkWithMsg("取消点赞成功", c)
	redis_article.SetCacheDigg(cr.ID, false)
	return
}
