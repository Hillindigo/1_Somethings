package model

import "time"

// Article 文章模型
type Article struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	Title       string     `json:"title" gorm:"size:200;not null;comment:文章标题"`
	Content     string     `json:"content" gorm:"type:text;not null;comment:文章内容(Markdown)"`
	Summary     string     `json:"summary" gorm:"size:500;comment:文章摘要"`
	CoverImage  string     `json:"cover_image" gorm:"size:255;comment:封面图URL"`
	Status      int8       `json:"status" gorm:"default:0;comment:状态 0-草稿 1-已发布"`
	ViewCount   int        `json:"view_count" gorm:"default:0;comment:浏览次数"`
	UserID      uint       `json:"user_id" gorm:"index;not null;comment:作者ID"`
	CategoryID  *uint      `json:"category_id" gorm:"index;comment:分类ID"`
	PublishedAt *time.Time `json:"published_at" gorm:"comment:发布时间"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	// 关联关系
	User     User     `json:"user" gorm:"foreignKey:UserID"`
	Category Category `json:"category" gorm:"foreignKey:CategoryID"`
	Tags     []Tag    `json:"tags" gorm:"many2many:article_tags;"`
}

// TableName 指定文章表名
func (Article) TableName() string {
	return "articles"
}

// IsPublished 判断文章是否已发布
func (a *Article) IsPublished() bool {
	return a.Status == 1
}

// --- 请求/响应 DTO ---

// CreateArticleRequest 创建文章请求
type CreateArticleRequest struct {
	Title      string `json:"title" binding:"required,max=200"`
	Content    string `json:"content" binding:"required"`
	Summary    string `json:"summary" binding:"max=500"`
	CoverImage string `json:"cover_image" binding:"omitempty,url"`
	Status     int8   `json:"status" binding:"oneof=0 1"`
	CategoryID uint   `json:"category_id"`
	TagIDs     []uint `json:"tag_ids"`
}

// UpdateArticleRequest 更新文章请求
type UpdateArticleRequest struct {
	Title      string `json:"title" binding:"omitempty,max=200"`
	Content    string `json:"content"`
	Summary    string `json:"summary" binding:"max=500"`
	CoverImage string `json:"cover_image" binding:"omitempty,url"`
	Status     int8   `json:"status" binding:"omitempty,oneof=0 1"`
	CategoryID uint   `json:"category_id"`
	TagIDs     []uint `json:"tag_ids"`
}

// ArticleListRequest 文章列表查询请求
type ArticleListRequest struct {
	Page       int  `form:"page" binding:"min=1"`
	PageSize   int  `form:"page_size" binding:"min=1,max=100"`
	CategoryID uint `form:"category_id"`
	TagID      uint `form:"tag_id"`
	Status     int8 `form:"status"`
	UserID     uint `form:"user_id"`
}
