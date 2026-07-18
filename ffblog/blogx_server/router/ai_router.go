package router

import (
	"blogx_server/api"
	"blogx_server/api/ai_api"
	"blogx_server/middleware"
	"github.com/gin-gonic/gin"
)

func AiRouter(r *gin.RouterGroup) {
	app := api.App.AiApi
	r.POST("ai/analysis", middleware.AuthMiddleware, middleware.BindJsonMiddleware[ai_api.ArticleAnalysisRequest], app.ArticleAnalysisView)
	r.GET("ai/article", middleware.AuthMiddleware, middleware.BindQueryMiddleware[ai_api.ArticleAiRequest], app.ArticleAiView)
}
