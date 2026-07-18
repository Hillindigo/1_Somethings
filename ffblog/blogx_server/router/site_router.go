// router/site_router.go
package router

import (
	"blogx_server/api"
	"blogx_server/middleware"
	"github.com/gin-gonic/gin"
)

func SiteRouter(r *gin.RouterGroup) {
	app := api.App.SiteApi
	r.GET("site/qq_url", app.SiteInfoQQView)
	r.GET("site/ai_info", app.SiteInfoAiView)
	r.GET("site/:name", app.SiteInfoView)

	r.PUT("site/:name", middleware.AdminMiddleware, app.SiteUpdateView)
}
