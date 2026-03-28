package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/schema"
	"github.com/joho/godotenv"
)

func ChatGenerate() {
	// 加载环境变量
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx := context.Background()
	timeout := 30 * time.Second

	// 初始化模型
	model, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey:  os.Getenv("API_KEY"),
		Model:   os.Getenv("MODEL"),
		Timeout: &timeout,
	})
	if err != nil {
		log.Fatalf("模型初始化失败: %v", err)
	}

	// 构造消息
	messages := []*schema.Message{
		schema.SystemMessage("你是一个智能助手"),
		schema.UserMessage("你好！，今天的天气怎么样？"),
	}

	// 调用生成
	response, err := model.Generate(ctx, messages)
	if err != nil {
		log.Fatalf("调用失败: %v", err)
	}

	fmt.Println("回复：", response.Content)
}

// 程序入口
func main() {
	ChatGenerate()
}
