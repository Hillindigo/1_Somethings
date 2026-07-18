package service

import (
	"errors"
	"time"

	"github.com/goagent/mojian/internal/model"
	"github.com/goagent/mojian/internal/repository"
	"github.com/goagent/mojian/pkg/errcode"
	"gorm.io/gorm"
)

// ArticleService 文章业务逻辑层
type ArticleService struct {
	articleRepo  *repository.ArticleRepository
	categoryRepo *repository.CategoryRepository
	tagRepo      *repository.TagRepository
}

// NewArticleService 创建文章 Service 实例
func NewArticleService(articleRepo *repository.ArticleRepository, categoryRepo *repository.CategoryRepository, tagRepo *repository.TagRepository) *ArticleService {
	return &ArticleService{
		articleRepo:  articleRepo,
		categoryRepo: categoryRepo,
		tagRepo:      tagRepo,
	}
}

// Create 创建文章，处理标签关联和发布时间
func (s *ArticleService) Create(userID uint, req *model.CreateArticleRequest) (*model.Article, error) {
	article := &model.Article{
		Title:      req.Title,
		Content:    req.Content,
		Summary:    req.Summary,
		CoverImage: req.CoverImage,
		Status:     req.Status,
		UserID:     userID,
	}

	// 处理分类ID：0 表示无分类，需设为 nil 以避免外键约束冲突
	if req.CategoryID > 0 {
		article.CategoryID = &req.CategoryID
	}

	// 如果是发布状态，设置发布时间
	if req.Status == 1 {
		now := time.Now()
		article.PublishedAt = &now
	}

	// 查询并关联标签
	if len(req.TagIDs) > 0 {
		tags, err := s.tagRepo.FindByIDs(req.TagIDs)
		if err != nil {
			return nil, err
		}
		article.Tags = tags
	}

	if err := s.articleRepo.Create(article); err != nil {
		return nil, errors.New(errcode.GetMessage(errcode.ErrArticleCreateFailed))
	}

	return article, nil
}

// GetByID 根据 ID 获取文章详情，同时增加浏览次数
func (s *ArticleService) GetByID(id uint) (*model.Article, error) {
	article, err := s.articleRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(errcode.GetMessage(errcode.ErrArticleNotFound))
		}
		return nil, err
	}

	// 增加浏览次数（忽略错误，不影响主流程）
	_ = s.articleRepo.IncrementViewCount(id)

	return article, nil
}

// List 分页查询文章列表
func (s *ArticleService) List(req *model.ArticleListRequest) ([]model.Article, int64, error) {
	return s.articleRepo.List(req)
}

// Update 更新文章，处理标签关联变更和发布状态变更
func (s *ArticleService) Update(id uint, userID uint, req *model.UpdateArticleRequest) (*model.Article, error) {
	article, err := s.articleRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(errcode.GetMessage(errcode.ErrArticleNotFound))
		}
		return nil, err
	}

	// 权限校验：只有作者本人可以修改
	if article.UserID != userID {
		return nil, errors.New(errcode.GetMessage(errcode.ErrArticleUpdateFailed))
	}

	// 更新字段
	if req.Title != "" {
		article.Title = req.Title
	}
	if req.Content != "" {
		article.Content = req.Content
	}
	if req.Summary != "" {
		article.Summary = req.Summary
	}
	if req.CoverImage != "" {
		article.CoverImage = req.CoverImage
	}
	// 处理分类ID：0 表示无分类，需设为 nil；>0 表示指定分类
	if req.CategoryID > 0 {
		article.CategoryID = &req.CategoryID
	} else {
		article.CategoryID = nil
	}

	// 处理发布状态变更
	if req.Status == 1 && article.Status == 0 {
		now := time.Now()
		article.PublishedAt = &now
	}
	article.Status = req.Status

	if err := s.articleRepo.Update(article); err != nil {
		return nil, errors.New(errcode.GetMessage(errcode.ErrArticleUpdateFailed))
	}

	// 更新标签关联
	if req.TagIDs != nil {
		tags, err := s.tagRepo.FindByIDs(req.TagIDs)
		if err != nil {
			return nil, err
		}
		if err := s.articleRepo.ReplaceTags(id, tags); err != nil {
			return nil, err
		}
	}

	return article, nil
}

// Delete 删除文章
func (s *ArticleService) Delete(id uint, userID uint) error {
	article, err := s.articleRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New(errcode.GetMessage(errcode.ErrArticleNotFound))
		}
		return err
	}

	// 权限校验：只有作者本人可以删除
	if article.UserID != userID {
		return errors.New(errcode.GetMessage(errcode.ErrArticleDeleteFailed))
	}

	return s.articleRepo.Delete(id)
}
