package comment_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middleware"
	"blogx_server/models"
	"blogx_server/models/enum"
	"blogx_server/service/comment_service"
	"blogx_server/service/message_service"
	"blogx_server/service/redis_service/redis_comment"
	"blogx_server/utils/jwts"
	"fmt"
	"github.com/gin-gonic/gin"
)

func (CommentApi) CommentRemoveApi(c *gin.Context) {
	var cr = middleware.GetBind[models.IDRequest](c)

	claims := jwts.GetClaims(c)

	var model models.CommentModel
	err := global.DB.Preload("ArticleModel").Take(&model, cr.ID).Error
	if err != nil {
		res.FailWithMsg("评论不存在", c)
		return
	}

	if claims.Role != enum.AdminRole {
		// 普通用户只能删自己发的评论，或者自己发的文章的评论
		if !(model.UserID == claims.UserID || model.ArticleModel.UserID == claims.UserID) {
			res.FailWithMsg("权限错误", c)
			return
		}
	}

	message_service.InsertSystemMessage(model.UserID, "管理员删除了你的评论", fmt.Sprintf("%s 内容不符合社区规范", model.Content), "", "")

	// 找所有的子评论，还要找所有的父评论
	subList := comment_service.GetCommentOneDimensional(model.ID)

	if model.ParentID != nil {
		// 有父评论
		parentList := comment_service.GetParents(*model.ParentID)
		for _, commentModel := range parentList {
			redis_comment.SetCacheApply(commentModel.ID, -len(subList))
		}
	}
	// 删评论
	global.DB.Delete(&subList)

	msg := fmt.Sprintf("删除成功，共删除评论%d条", len(subList))
	res.OkWithMsg(msg, c)

}
