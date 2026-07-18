package model

import "time"

// Category 分类模型，支持树形结构
type Category struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	Name      string     `json:"name" gorm:"size:50;not null;uniqueIndex;comment:分类名称"`
	ParentID  *uint      `json:"parent_id" gorm:"index;comment:父分类ID，顶级分类为NULL"`
	SortOrder int        `json:"sort_order" gorm:"default:0;comment:排序权重"`
	CreatedAt time.Time  `json:"created_at"`
	Children  []Category `json:"children,omitempty" gorm:"foreignKey:ParentID"`
}

// TableName 指定分类表名
func (Category) TableName() string {
	return "categories"
}

// --- 请求 DTO ---

// CreateCategoryRequest 创建分类请求
type CreateCategoryRequest struct {
	Name      string `json:"name" binding:"required,max=50"`
	ParentID  *uint  `json:"parent_id"`
	SortOrder int    `json:"sort_order"`
}

// UpdateCategoryRequest 更新分类请求
type UpdateCategoryRequest struct {
	Name      string `json:"name" binding:"omitempty,max=50"`
	ParentID  *uint  `json:"parent_id"`
	SortOrder int    `json:"sort_order"`
}
