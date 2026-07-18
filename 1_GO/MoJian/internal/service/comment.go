package service

import (
	"errors"

	"github.com/goagent/mojian/internal/model"
	"github.com/goagent/mojian/internal/repository"
	"github.com/goagent/mojian/pkg/errcode"
	"gorm.io/gorm"
)

// CommentService 评论业务逻辑层
type CommentService struct {
	repo *repository.CommentRepository
}

// NewCommentService 创建评论 Service 实例
func NewCommentService(repo *repository.CommentRepository) *CommentService {
	return &CommentService{repo: repo}
}

// Create 创建评论
func (s *CommentService) Create(userID uint, req *model.CreateCommentRequest) (*model.Comment, error) {
	comment := &model.Comment{
		Content:   req.Content,
		ArticleID: req.ArticleID,
		UserID:    userID,
		ParentID:  req.ParentID,
		Status:    1, // 默认通过审核，如需审核可改为 0
	}

	if err := s.repo.Create(comment); err != nil {
		return nil, errors.New(errcode.GetMessage(errcode.ErrCommentCreateFailed))
	}
	return comment, nil
}

// GetByID 根据 ID 获取评论
func (s *CommentService) GetByID(id uint) (*model.Comment, error) {
	comment, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(errcode.GetMessage(errcode.ErrCommentNotFound))
		}
		return nil, err
	}
	return comment, nil
}

// ListByArticleID 查询指定文章的评论列表
func (s *CommentService) ListByArticleID(articleID uint) ([]model.Comment, error) {
	return s.repo.ListByArticleID(articleID)
}

// Delete 删除评论
func (s *CommentService) Delete(id uint) error {
	return s.repo.Delete(id)
}

// UpdateStatus 更新评论审核状态（管理员操作）
func (s *CommentService) UpdateStatus(id uint, status int8) error {
	return s.repo.UpdateStatus(id, status)
}
