// router/site_router.go
package router

import (
	"blogx_server/api"
	"blogx_server/api/image_api"
	"blogx_server/middleware"
	"github.com/gin-gonic/gin"
)

func ImageRouter(r *gin.RouterGroup) {
	app := api.App.ImageApi
	r.POST("images", middleware.AuthMiddleware, app.ImageUploadView)
	r.POST("images/qiniu", middleware.AuthMiddleware, app.QiNiuGenToken)
	r.POST("images/transfer_deposit", middleware.AuthMiddleware, middleware.BindJsonMiddleware[image_api.TransferDepositRequest], app.TransferDepositView)
	r.GET("images", middleware.AdminMiddleware, app.ImageListView)
	r.DELETE("images", middleware.AdminMiddleware, app.ImageRemoveView)
}
