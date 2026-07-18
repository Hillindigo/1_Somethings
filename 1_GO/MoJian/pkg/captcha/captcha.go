package captcha

import (
	"time"

	"github.com/mojocn/base64Captcha"
)

// 默认全局验证码存储实例，最多缓存 10240 条，过期时间 5 分钟
var defaultStore = base64Captcha.NewMemoryStore(10240, 5*time.Minute)

// CaptchaResult 验证码生成结果
type CaptchaResult struct {
	CaptchaID string `json:"captcha_id"` // 验证码唯一标识
	Image     string `json:"image"`      // base64 编码的验证码图片（含 data:image/png;base64, 前缀）
}

// Generate 生成图片验证码，返回验证码ID和base64图片
// 使用 DriverString 生成字母数字混合验证码，4位字符
func Generate() (*CaptchaResult, error) {
	driver := base64Captcha.NewDriverString(
		80,  // height
		240, // width
		80,  // noiseCount 干扰点数
		base64Captcha.OptionShowHollowLine|base64Captcha.OptionShowSlimeLine, // showLineOptions
		4,                                // length 验证码字符数
		base64Captcha.TxtSimpleCharaters, // source 字符集（排除易混淆字符）
		nil,                              // bgColor
		nil,                              // fontsStorage
		nil,                              // fonts
	)

	captchaInst := base64Captcha.NewCaptcha(driver, defaultStore)
	id, b64s, _, err := captchaInst.Generate()
	if err != nil {
		return nil, err
	}

	return &CaptchaResult{
		CaptchaID: id,
		Image:     b64s,
	}, nil
}

// Verify 校验验证码
// id: 验证码ID  answer: 用户输入的验证码  clear: 是否清除已验证的验证码（一次性使用）
func Verify(id, answer string, clear bool) bool {
	return defaultStore.Verify(id, answer, clear)
}
