package main

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger
var sugarLogger *zap.SugaredLogger

// GinLogger 接收gin框架默认的日志
func GinLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		cost := time.Since(start)
		logger.Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
		)
	}
}

// GinRecovery recover掉项目可能出现的panic
func GinRecovery(logger *zap.Logger, stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logger.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}

func InitLogger() {
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	logger = zap.New(core) // 修复局部变量遮蔽问题
	sugarLogger = logger.Sugar()
}

func getEncoder() zapcore.Encoder {
	config := zap.NewDevelopmentEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return zapcore.NewConsoleEncoder(config)
}

func getLogWriter() zapcore.WriteSyncer {
	file, err := os.OpenFile("./test.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(fmt.Sprintf("创建日志文件失败: %v", err))
	}
	return zapcore.AddSync(file)
}

func simpleHttpGet(url string) {
	// 添加空行作为请求开始的分隔符
	sugarLogger.Info("")
	sugarLogger.Info("========== 开始请求 ==========")
	sugarLogger.Infow("请求URL", "url", url)

	start := time.Now()
	resp, err := http.Get(url)
	duration := time.Since(start)

	if err != nil {
		sugarLogger.Errorw(
			"HTTP请求失败",
			"url", url,
			"duration", duration,
			"error", err.Error(),
			"stack", zap.Stack("stack"),
		)
		// 添加空行作为错误结束的分隔符
		sugarLogger.Info("")
		return
	}

	defer resp.Body.Close()

	sugarLogger.Infow(
		"HTTP请求成功",
		"url", url,
		"status", resp.Status,
		"status_code", resp.StatusCode,
		"duration", duration,
		"content_length", resp.ContentLength,
	)

	// 添加空行作为成功请求的分隔符
	sugarLogger.Info("========== 请求结束 ==========")
	sugarLogger.Info("")
}

func main() {
	InitLogger()
	r := gin.New()
	r.Use(GinLogger(logger), GinRecovery(logger, true))
	r.GET("a", func(c *gin.Context) {
		c.String(http.StatusOK, "你好")
	})
	r.Run(":8888")
	getLogWriter()
}
