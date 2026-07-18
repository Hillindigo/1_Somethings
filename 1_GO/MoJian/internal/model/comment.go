package model

import "time"

// Comment 评论模型，支持树形回复
type Comment struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Content   string    `json:"content" gorm:"type:text;not null;comment:评论内容"`
	ArticleID uint      `json:"article_id" gorm:"index;not null;comment:文章ID"`
	UserID    uint      `json:"user_id" gorm:"index;not null;comment:评论者ID"`
	ParentID  *uint     `json:"parent_id" gorm:"index;comment:父评论ID，顶级评论为NULL"`
	Status    int8      `json:"status" gorm:"default:1;comment:状态 0-待审核 1-已通过 2-已拒绝"`
	CreatedAt time.Time `json:"created_at"`

	// 关联关系
	User     User      `json:"user" gorm:"foreignKey:UserID"`
	Children []Comment `json:"children,omitempty" gorm:"foreignKey:ParentID"`
}

// TableName 指定评论表名
func (Comment) TableName() string {
	return "comments"
}

// --- 请求 DTO ---

// CreateCommentRequest 创建评论请求
type CreateCommentRequest struct {
	Content   string `json:"content" binding:"required,max=1000"`
	ArticleID uint   `json:"article_id" binding:"required"`
	ParentID  *uint  `json:"parent_id"`
}
