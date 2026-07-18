package site_msg_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middleware"
	"blogx_server/models"
	"blogx_server/models/enum/message_type_enum"
	"blogx_server/utils/jwts"
	"fmt"
	"github.com/gin-gonic/gin"
)

type SiteMsgReadRequest struct {
	ID uint `json:"id"`
	T  int8 `json:"t"` // 一键已读的类型
}

func (SiteMsgApi) SiteMsgReadView(c *gin.Context) {
	cr := middleware.GetBind[SiteMsgReadRequest](c)

	claims := jwts.GetClaims(c)
	if cr.ID != 0 {
		// 找这个消息是不是当前用户的
		var msg models.MessageModel
		err := global.DB.Take(&msg, "id = ? and rev_user_id = ?", cr.ID, claims.UserID).Error
		if err != nil {
			res.FailWithMsg("消息不存在", c)
			return
		}

		if msg.IsRead {
			res.FailWithMsg("消息已读取", c)
			return
		}

		global.DB.Model(&msg).Update("is_read", true)
		res.OkWithMsg("消息读取成功", c)
		return
	}

	var typeList []message_type_enum.Type
	switch cr.T {
	case 1:
		typeList = append(typeList, message_type_enum.CommentType, message_type_enum.ApplyType)
	case 2:
		typeList = append(typeList, message_type_enum.DiggArticleType, message_type_enum.DiggCommentType, message_type_enum.CollectArticleType)
	case 3:
		typeList = append(typeList, message_type_enum.SystemType)
	}

	var msgList []models.MessageModel
	global.DB.Find(&msgList, "rev_user_id = ? and type in ? and is_read = ?", claims.UserID, typeList, false)

	if len(msgList) > 0 {
		global.DB.Model(&msgList).Update("is_read", true)
	}

	res.OkWithMsg(fmt.Sprintf("批量读取%d条消息成功", len(msgList)), c)

}
