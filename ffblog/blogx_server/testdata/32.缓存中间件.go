package main

import (
	"blogx_server/common/res"
	"blogx_server/core"
	"blogx_server/flags"
	"blogx_server/global"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func xxx(c *gin.Context) {
	fmt.Println(c.Request.URL.String(), time.Now().Format("15:04:05"))
	res.OkWithData("123", c)
}

type ResponseWriter struct {
	gin.ResponseWriter
	Body []byte
}

func (w *ResponseWriter) Write(data []byte) (int, error) {
	w.Body = append(w.Body, data...)
	return w.ResponseWriter.Write(data)
}
func CacheMiddleware(cacheTime time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.Request.URL.String()
		// 请求部分
		val, err := global.Redis.Get(key).Result()
		if err == nil {
			c.Abort()
			fmt.Println("走缓存了")
			c.Header("Content-Type", "application/json; charset=utf-8")
			c.Writer.Write([]byte(val))
			return
		}
		w := &ResponseWriter{
			ResponseWriter: c.Writer,
		}
		c.Writer = w
		c.Next()
		// 响应
		body := string(w.Body)
		// 加入到缓存里面
		global.Redis.Set(key, body, cacheTime)
	}
}

func main() {
	flags.Parse()
	global.Config = core.ReadConf()
	core.InitLogrus()
	global.Redis = core.InitRedis()

	r := gin.Default()
	r.GET("/cache/xxx", CacheMiddleware(10*time.Second), xxx)
	r.GET("/xxx", xxx)
	r.Run(":8081")
}
