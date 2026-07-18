package model

import "time"

// User 用户模型
type User struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Username     string    `json:"username" gorm:"uniqueIndex;size:50;not null;comment:用户名"`
	PasswordHash string    `json:"-" gorm:"size:255;not null;comment:密码哈希"`
	Email        string    `json:"email" gorm:"uniqueIndex;size:100;not null;comment:邮箱"`
	Avatar       string    `json:"avatar" gorm:"size:255;comment:头像URL"`
	Role         int8      `json:"role" gorm:"default:0;comment:角色 0-普通用户 1-管理员"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TableName 指定用户表名
func (User) TableName() string {
	return "users"
}

// IsAdmin 判断用户是否为管理员
func (u *User) IsAdmin() bool {
	return u.Role == 1
}

// --- 请求/响应 DTO ---

// RegisterRequest 用户注册请求
type RegisterRequest struct {
	Username    string `json:"username" binding:"required,min=2,max=50"`
	Password    string `json:"password" binding:"required,min=6,max=50"`
	Email       string `json:"email" binding:"required,email"`
	CaptchaID   string `json:"captcha_id" binding:"required"`   // 验证码唯一标识
	CaptchaCode string `json:"captcha_code" binding:"required"` // 用户输入的验证码
}

// LoginRequest 用户登录请求
type LoginRequest struct {
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required"`
	CaptchaID   string `json:"captcha_id" binding:"required"`   // 验证码唯一标识
	CaptchaCode string `json:"captcha_code" binding:"required"` // 用户输入的验证码
}

// LoginResponse 登录成功响应
type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// UpdateUserRequest 更新用户信息请求
type UpdateUserRequest struct {
	Email  string `json:"email" binding:"omitempty,email"`
	Avatar string `json:"avatar" binding:"omitempty,url"`
}
