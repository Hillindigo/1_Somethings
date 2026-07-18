package errcode

// 通用错误码 (1xxxx)
const (
	ErrInvalidParams = 10001 // 参数校验失败
	ErrTokenGenerate = 10002 // Token 生成失败
	ErrTokenInvalid  = 10003 // Token 无效或已过期
)

// 用户模块错误码 (2xxxx)
const (
	ErrUserNotFound      = 20001 // 用户不存在
	ErrUserAlreadyExists = 20002 // 用户名已存在
	ErrEmailAlreadyUsed  = 20003 // 邮箱已被使用
	ErrPasswordIncorrect = 20004 // 密码错误
	ErrUserCreateFailed  = 20005 // 创建用户失败
)

// 文章模块错误码 (3xxxx)
const (
	ErrArticleNotFound     = 30001 // 文章不存在
	ErrArticleCreateFailed = 30002 // 创建文章失败
	ErrArticleUpdateFailed = 30003 // 更新文章失败
	ErrArticleDeleteFailed = 30004 // 删除文章失败
)

// 分类模块错误码 (4xxxx)
const (
	ErrCategoryNotFound      = 40001 // 分类不存在
	ErrCategoryAlreadyExists = 40002 // 分类名已存在
	ErrCategoryCreateFailed  = 40003 // 创建分类失败
	ErrCategoryHasArticles   = 40004 // 分类下存在文章，无法删除
)

// 标签模块错误码 (5xxxx)
const (
	ErrTagNotFound      = 50001 // 标签不存在
	ErrTagAlreadyExists = 50002 // 标签名已存在
	ErrTagCreateFailed  = 50003 // 创建标签失败
)

// 评论模块错误码 (6xxxx)
const (
	ErrCommentNotFound     = 60001 // 评论不存在
	ErrCommentCreateFailed = 60002 // 创建评论失败
)

// 验证码模块错误码 (7xxxx)
const (
	ErrCaptchaInvalid  = 70001 // 验证码错误或已过期
	ErrCaptchaGenerate = 70002 // 验证码生成失败
	ErrCaptchaRequired = 70003 // 验证码参数缺失
)

// codeMessageMap 错误码与提示信息的映射
var codeMessageMap = map[int]string{
	ErrInvalidParams: "参数校验失败",
	ErrTokenGenerate: "Token 生成失败",
	ErrTokenInvalid:  "Token 无效或已过期",

	ErrUserNotFound:      "用户不存在",
	ErrUserAlreadyExists: "用户名已存在",
	ErrEmailAlreadyUsed:  "邮箱已被使用",
	ErrPasswordIncorrect: "密码错误",
	ErrUserCreateFailed:  "创建用户失败",

	ErrArticleNotFound:     "文章不存在",
	ErrArticleCreateFailed: "创建文章失败",
	ErrArticleUpdateFailed: "更新文章失败",
	ErrArticleDeleteFailed: "删除文章失败",

	ErrCategoryNotFound:      "分类不存在",
	ErrCategoryAlreadyExists: "分类名已存在",
	ErrCategoryCreateFailed:  "创建分类失败",
	ErrCategoryHasArticles:   "分类下存在文章，无法删除",

	ErrTagNotFound:      "标签不存在",
	ErrTagAlreadyExists: "标签名已存在",
	ErrTagCreateFailed:  "创建标签失败",

	ErrCommentNotFound:     "评论不存在",
	ErrCommentCreateFailed: "创建评论失败",

	ErrCaptchaInvalid:  "验证码错误或已过期",
	ErrCaptchaGenerate: "验证码生成失败",
	ErrCaptchaRequired: "验证码参数缺失",
}

// GetMessage 根据错误码获取对应的提示信息
func GetMessage(code int) string {
	if msg, ok := codeMessageMap[code]; ok {
		return msg
	}
	return "未知错误"
}
