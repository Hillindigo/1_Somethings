package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/goagent/mojian/internal/config"
	"github.com/goagent/mojian/internal/database"
	"github.com/goagent/mojian/internal/handler"
	"github.com/goagent/mojian/internal/repository"
	"github.com/goagent/mojian/internal/router"
	"github.com/goagent/mojian/internal/service"
)

func main() {
	// 1. 加载配置
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 2. 初始化数据库
	db, err := database.InitDB(&cfg.Database)
	if err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}

	// 3. 设置 Gin 运行模式
	gin.SetMode(cfg.Server.Mode)

	// 4. 依赖注入：Repository → Service → Handler
	// Repository 层
	userRepo := repository.NewUserRepository(db)
	articleRepo := repository.NewArticleRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	tagRepo := repository.NewTagRepository(db)
	commentRepo := repository.NewCommentRepository(db)

	// Service 层
	userSvc := service.NewUserService(userRepo)
	articleSvc := service.NewArticleService(articleRepo, categoryRepo, tagRepo)
	categorySvc := service.NewCategoryService(categoryRepo, articleRepo)
	tagSvc := service.NewTagService(tagRepo)
	commentSvc := service.NewCommentService(commentRepo)

	// Handler 层
	userHandler := handler.NewUserHandler(userSvc)
	articleHandler := handler.NewArticleHandler(articleSvc)
	categoryHandler := handler.NewCategoryHandler(categorySvc)
	tagHandler := handler.NewTagHandler(tagSvc)
	commentHandler := handler.NewCommentHandler(commentSvc)
	captchaHandler := handler.NewCaptchaHandler()

	// 5. 注册路由
	r := router.NewRouter(userHandler, articleHandler, categoryHandler, tagHandler, commentHandler, captchaHandler)
	engine := r.Setup()

	// 6. 启动服务
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("服务启动成功，监听端口 %s", addr)
	if err := engine.Run(addr); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
