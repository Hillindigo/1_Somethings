// router/site_router.go
package router

import (
	"blogx_server/api"
	"blogx_server/api/comment_api"
	"blogx_server/middleware"
	"blogx_server/models"
	"github.com/gin-gonic/gin"
)

func CommentRouter(r *gin.RouterGroup) {
	app := api.App.CommentApi
	r.POST("comment", middleware.AuthMiddleware, middleware.BindJsonMiddleware[comment_api.CommentCreateRequest], app.CommentCreateView)
	r.GET("comment/tree/:id", middleware.BindUriMiddleware[models.IDRequest], app.CommentTreeView)
	r.GET("comment", middleware.AuthMiddleware, middleware.BindQueryMiddleware[comment_api.CommentListRequest], app.CommentListView)
	r.DELETE("comment/:id", middleware.AuthMiddleware, middleware.BindUriMiddleware[models.IDRequest], app.CommentRemoveApi)
	r.GET("comment/digg/:id", middleware.AuthMiddleware, middleware.BindUriMiddleware[models.IDRequest], app.CommentDiggView)
}
