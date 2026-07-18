package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/goagent/mojian/internal/config"
	"github.com/goagent/mojian/internal/database"
	"github.com/goagent/mojian/internal/model"
	"github.com/goagent/mojian/internal/repository"
	"github.com/goagent/mojian/pkg/utils"
)

// createadmin 命令行工具：通过终端创建管理员账号
// 用法: go run cmd/createadmin/main.go -username admin -password admin123 -email admin@example.com

func main() {
	// 解析命令行参数
	username := flag.String("username", "", "管理员用户名（必填）")
	password := flag.String("password", "", "管理员密码（必填，至少6位）")
	email := flag.String("email", "", "管理员邮箱（必填）")
	configPath := flag.String("config", "config/config.yaml", "配置文件路径")
	flag.Parse()

	// 校验必填参数
	if *username == "" || *password == "" || *email == "" {
		log.Fatal("参数不完整，用法: go run cmd/createadmin/main.go -username <用户名> -password <密码> -email <邮箱>")
	}
	if len(*password) < 6 {
		log.Fatal("密码长度不能少于6位")
	}

	// 1. 加载配置
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 2. 初始化数据库
	db, err := database.InitDB(&cfg.Database)
	if err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}

	// 3. 初始化 Repository
	userRepo := repository.NewUserRepository(db)

	// 4. 检查用户名是否已存在
	existingUser, err := userRepo.FindByUsername(*username)
	if err == nil && existingUser != nil {
		// 用户已存在，询问是否升级为管理员
		if existingUser.Role == 1 {
			log.Fatalf("用户 %s 已经是管理员", *username)
		}
		existingUser.Role = 1
		if err := userRepo.Update(existingUser); err != nil {
			log.Fatalf("升级用户为管理员失败: %v", err)
		}
		fmt.Printf("用户 %s 已成功升级为管理员！\n", *username)
		return
	}

	// 5. 检查邮箱是否已存在
	if _, err := userRepo.FindByEmail(*email); err == nil {
		log.Fatalf("邮箱 %s 已被使用", *email)
	}

	// 6. 密码哈希加密
	hash, err := utils.HashPassword(*password)
	if err != nil {
		log.Fatalf("密码加密失败: %v", err)
	}

	// 7. 创建管理员用户
	admin := &model.User{
		Username:     *username,
		PasswordHash: hash,
		Email:        *email,
		Role:         1, // 管理员角色
	}

	if err := userRepo.Create(admin); err != nil {
		log.Fatalf("创建管理员失败: %v", err)
	}

	fmt.Printf("管理员 %s 创建成功！\n", *username)
}
