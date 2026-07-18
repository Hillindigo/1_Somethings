// api/enter.go
package api

import (
	"blogx_server/api/ai_api"
	"blogx_server/api/article_api"
	"blogx_server/api/banner_api"
	"blogx_server/api/captcha_api"
	"blogx_server/api/chat_api"
	"blogx_server/api/comment_api"
	"blogx_server/api/data_api"
	"blogx_server/api/focus_api"
	"blogx_server/api/global_notification_api"
	"blogx_server/api/image_api"
	"blogx_server/api/log_api"
	"blogx_server/api/search_api"
	"blogx_server/api/site_api"
	"blogx_server/api/site_msg_api"
	"blogx_server/api/user_api"
)

type Api struct {
	SiteApi               site_api.SiteApi
	LogApi                log_api.LogApi
	ImageApi              image_api.ImageApi
	BannerApi             banner_api.BannerApi
	CaptchaApi            captcha_api.CaptchaApi
	UserApi               user_api.UserApi
	ArticleApi            article_api.ArticleApi
	CommentApi            comment_api.CommentApi
	SiteMsgApi            site_msg_api.SiteMsgApi
	GlobalNotificationApi global_notification_api.GlobalNotificationApi
	FocusApi              focus_api.FocusApi
	ChatApi               chat_api.ChatApi
	SearchApi             search_api.SearchApi
	AiApi                 ai_api.AiApi
	DataApi               data_api.DataApi
}

var App = Api{}
