package model

import "time"

// Tag 标签模型
type Tag struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"size:30;not null;uniqueIndex;comment:标签名称"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName 指定标签表名
func (Tag) TableName() string {
	return "tags"
}

// --- 请求 DTO ---

// CreateTagRequest 创建标签请求
type CreateTagRequest struct {
	Name string `json:"name" binding:"required,max=30"`
}

// UpdateTagRequest 更新标签请求
type UpdateTagRequest struct {
	Name string `json:"name" binding:"required,max=30"`
}
