package main

import (
	"blogx_server/common/res"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func main() {
	r := gin.Default()
	r.GET("/sse", func(c *gin.Context) {
		for i := 0; i < 5; i++ {
			res.SSEOk(fmt.Sprintf("第%d条数据", i+1), c)
			time.Sleep(1 * time.Second)
		}
	})
	r.Run(":8081")
}
