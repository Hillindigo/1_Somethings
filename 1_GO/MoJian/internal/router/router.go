package router

import (
	"github.com/gin-gonic/gin"
	"github.com/goagent/mojian/internal/handler"
	"github.com/goagent/mojian/internal/middleware"
)

// Router 路由管理器，负责注册所有 API 路由
type Router struct {
	userHandler     *handler.UserHandler
	articleHandler  *handler.ArticleHandler
	categoryHandler *handler.CategoryHandler
	tagHandler      *handler.TagHandler
	commentHandler  *handler.CommentHandler
	captchaHandler  *handler.CaptchaHandler
}

// NewRouter 创建路由管理器实例，通过构造函数注入所有 Handler 依赖
func NewRouter(
	userHandler *handler.UserHandler,
	articleHandler *handler.ArticleHandler,
	categoryHandler *handler.CategoryHandler,
	tagHandler *handler.TagHandler,
	commentHandler *handler.CommentHandler,
	captchaHandler *handler.CaptchaHandler,
) *Router {
	return &Router{
		userHandler:     userHandler,
		articleHandler:  articleHandler,
		categoryHandler: categoryHandler,
		tagHandler:      tagHandler,
		commentHandler:  commentHandler,
		captchaHandler:  captchaHandler,
	}
}

// Setup 注册所有路由规则，返回配置好的 Gin Engine
func (r *Router) Setup() *gin.Engine {
	engine := gin.New()

	// 全局中间件
	engine.Use(gin.Recovery())             // panic 恢复
	engine.Use(middleware.RequestLogger()) // 请求日志
	engine.Use(middleware.CORS())          // 跨域

	// 健康检查
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	v1 := engine.Group("/api/v1")
	{
		// --- 公开接口（无需认证） ---
		auth := v1.Group("/auth")
		{
			auth.POST("/register", r.userHandler.Register)
			auth.POST("/login", r.userHandler.Login)
			auth.GET("/captcha", r.captchaHandler.GetCaptcha)
		}

		// 文章公开访问
		v1.GET("/articles", r.articleHandler.List)
		v1.GET("/articles/:id", r.articleHandler.GetByID)

		// 分类公开访问
		v1.GET("/categories", r.categoryHandler.List)
		v1.GET("/categories/:id", r.categoryHandler.GetByID)

		// 标签公开访问
		v1.GET("/tags", r.tagHandler.List)
		v1.GET("/tags/:id", r.tagHandler.GetByID)

		// 评论公开访问（使用 OptionalJWTAuth 识别游客身份，游客最多5条）
		v1.GET("/articles/:id/comments", middleware.OptionalJWTAuth(), r.commentHandler.ListByArticleID)

		// --- 需要认证的接口（管理员 + 普通用户） ---
		authenticated := v1.Group("")
		authenticated.Use(middleware.JWTAuth())
		{
			// 用户信息
			authenticated.GET("/user/profile", r.userHandler.GetUser)
			authenticated.PUT("/user/profile", r.userHandler.UpdateUser)

			// 文章管理（需要认证：登录用户可发布/编辑/删除自己的文章）
			authenticated.POST("/articles", r.articleHandler.Create)
			authenticated.PUT("/articles/:id", r.articleHandler.Update)
			authenticated.DELETE("/articles/:id", r.articleHandler.Delete)

			// 分类管理（需要认证：管理员和普通用户均可管理分类）
			authenticated.POST("/categories", r.categoryHandler.Create)
			authenticated.PUT("/categories/:id", r.categoryHandler.Update)
			authenticated.DELETE("/categories/:id", r.categoryHandler.Delete)

			// 评论（需要认证）
			authenticated.POST("/comments", r.commentHandler.Create)
			authenticated.DELETE("/comments/:id", r.commentHandler.Delete)
		}

		// --- 管理员接口 ---
		admin := v1.Group("/admin")
		admin.Use(middleware.JWTAuth(), middleware.AdminRequired())
		{
			// 用户管理
			admin.GET("/users", r.userHandler.ListUsers)
			admin.PUT("/users/:id/role", r.userHandler.UpdateUserRole)
			admin.DELETE("/users/:id", r.userHandler.DeleteUser)

			// 标签管理
			admin.POST("/tags", r.tagHandler.Create)
			admin.PUT("/tags/:id", r.tagHandler.Update)
			admin.DELETE("/tags/:id", r.tagHandler.Delete)

			// 评论审核
			admin.PUT("/comments/:id/status", r.commentHandler.UpdateStatus)
		}
	}

	return engine
}
