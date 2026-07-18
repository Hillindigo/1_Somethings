package repository

import (
	"github.com/goagent/mojian/internal/model"
	"gorm.io/gorm"
)

// CommentRepository 评论数据访问层
type CommentRepository struct {
	db *gorm.DB
}

// NewCommentRepository 创建评论 Repository 实例
func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

// Create 创建评论
func (r *CommentRepository) Create(comment *model.Comment) error {
	return r.db.Create(comment).Error
}

// FindByID 根据 ID 查询评论
func (r *CommentRepository) FindByID(id uint) (*model.Comment, error) {
	var comment model.Comment
	if err := r.db.Preload("User").First(&comment, id).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

// ListByArticleID 查询指定文章下的评论列表，预加载评论者信息
func (r *CommentRepository) ListByArticleID(articleID uint) ([]model.Comment, error) {
	var comments []model.Comment
	if err := r.db.Where("article_id = ? AND status = 1", articleID).
		Preload("User").
		Order("created_at ASC").
		Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

// Delete 根据 ID 删除评论
func (r *CommentRepository) Delete(id uint) error {
	return r.db.Delete(&model.Comment{}, id).Error
}

// UpdateStatus 更新评论审核状态
func (r *CommentRepository) UpdateStatus(id uint, status int8) error {
	return r.db.Model(&model.Comment{}).Where("id = ?", id).Update("status", status).Error
}
