package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/goagent/mojian/internal/config"
	"github.com/goagent/mojian/pkg/response"
	"github.com/goagent/mojian/pkg/utils"
)

// JWTAuth JWT 认证中间件，从请求头 Authorization 中提取并验证 Token
// 验证通过后将 user_id 和 role 写入 gin.Context
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.FailWithUnauthorized(c, "请求头中缺少 Authorization")
			c.Abort()
			return
		}

		// 格式: Bearer <token>
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.FailWithUnauthorized(c, "Authorization 格式错误，应为 Bearer <token>")
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(parts[1], config.GlobalConfig.JWT.Secret)
		if err != nil {
			response.FailWithUnauthorized(c, "Token 无效或已过期")
			c.Abort()
			return
		}

		// 将用户信息存入上下文，后续 Handler 可通过 c.Get 获取
		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}

// OptionalJWTAuth 可选 JWT 认证中间件，尝试解析 Token 但不强制要求登录
// 如果 Token 存在且有效，将 user_id 和 role 写入上下文；否则不设置
func OptionalJWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && parts[0] == "Bearer" {
				claims, err := utils.ParseToken(parts[1], config.GlobalConfig.JWT.Secret)
				if err == nil {
					c.Set("user_id", claims.UserID)
					c.Set("role", claims.Role)
				}
			}
		}
		c.Next()
	}
}

// AdminRequired 管理员权限校验中间件，必须在 JWTAuth 之后使用
func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role.(int8) != 1 {
			response.FailWithForbidden(c, "需要管理员权限")
			c.Abort()
			return
		}
		c.Next()
	}
}
