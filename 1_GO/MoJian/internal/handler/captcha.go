package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/goagent/mojian/pkg/captcha"
	"github.com/goagent/mojian/pkg/response"
)

// CaptchaHandler 验证码处理器，处理验证码相关的 HTTP 请求
type CaptchaHandler struct{}

// NewCaptchaHandler 创建验证码 Handler 实例
func NewCaptchaHandler() *CaptchaHandler {
	return &CaptchaHandler{}
}

// GetCaptcha 获取图片验证码接口
// GET /api/v1/auth/captcha
// 返回 captcha_id 和 base64 编码的验证码图片
func (h *CaptchaHandler) GetCaptcha(c *gin.Context) {
	result, err := captcha.Generate()
	if err != nil {
		response.FailWithServerError(c, "验证码生成失败")
		return
	}

	response.Success(c, result)
}
