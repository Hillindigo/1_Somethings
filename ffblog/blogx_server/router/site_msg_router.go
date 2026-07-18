// router/site_router.go
package router

import (
	"blogx_server/api"
	"blogx_server/api/site_msg_api"
	"blogx_server/middleware"
	"github.com/gin-gonic/gin"
)

func SiteMsgRouter(r *gin.RouterGroup) {
	app := api.App.SiteMsgApi
	r.GET("site_msg", middleware.AuthMiddleware, middleware.BindQueryMiddleware[site_msg_api.SiteMsgListRequest], app.SiteMsgListView)

	r.GET("site_msg/conf", middleware.AuthMiddleware, app.UserSiteMessageConfView)
	r.PUT("site_msg/conf", middleware.AuthMiddleware, middleware.BindJsonMiddleware[site_msg_api.UserMessageConfUpdateRequest], app.UserSiteMessageConfUpdateView)
	r.POST("site_msg", middleware.AuthMiddleware, middleware.BindJsonMiddleware[site_msg_api.SiteMsgReadRequest], app.SiteMsgReadView)
	r.DELETE("site_msg", middleware.AuthMiddleware, middleware.BindJsonMiddleware[site_msg_api.SiteMsgRemoveRequest], app.SiteMsgRemoveView)

	r.GET("site_msg/user", middleware.AuthMiddleware, app.UserMsgView)

}
