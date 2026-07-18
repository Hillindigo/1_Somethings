package chat_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middleware"
	"blogx_server/models"
	"blogx_server/models/ctype/chat_msg"
	"blogx_server/models/enum/chat_msg_type"
	"blogx_server/utils/jwts"
	"github.com/gin-gonic/gin"
)

func (ChatApi) ChatReadView(c *gin.Context) {
	cr := middleware.GetBind[models.IDRequest](c)

	var chat models.ChatModel
	err := global.DB.Take(&chat, cr.ID).Error
	if err != nil {
		res.FailWithMsg("消息不存在", c)
		return
	}
	item := ChatResponse{
		ChatListResponse: ChatListResponse{
			ChatModel: models.ChatModel{
				MsgType: chat_msg_type.MsgReadType,
				Msg: chat_msg.ChatMsg{
					MsgReadMsg: &chat_msg.MsgReadMsg{
						ReadChatID: chat.ID,
					},
				},
			},
		},
	}
	claims := jwts.GetClaims(c)
	var chatAc models.UserChatActionModel
	err = global.DB.Take(&chatAc, "user_id = ? and chat_id = ?", claims.UserID, cr.ID).Error
	if err != nil {
		global.DB.Create(&models.UserChatActionModel{
			UserID: claims.UserID,
			ChatID: cr.ID,
			IsRead: true,
		})
		res.SendWsMsg(OnlineMap, chat.SendUserID, item)
		res.OkWithMsg("消息读取成功", c)
		return
	}

	if chatAc.IsDelete {
		res.FailWithMsg("消息被删除", c)
		return
	}

	res.SendWsMsg(OnlineMap, chat.SendUserID, item)
	global.DB.Model(&chatAc).Update("is_read", true)
	res.OkWithMsg("消息读取成功", c)
}
