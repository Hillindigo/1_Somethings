package service

import (
	"errors"

	"github.com/goagent/mojian/internal/model"
	"github.com/goagent/mojian/internal/repository"
	"github.com/goagent/mojian/pkg/errcode"
	"gorm.io/gorm"
)

// CategoryService 分类业务逻辑层
type CategoryService struct {
	repo        *repository.CategoryRepository
	articleRepo *repository.ArticleRepository
}

// NewCategoryService 创建分类 Service 实例，注入分类和文章 Repository
func NewCategoryService(repo *repository.CategoryRepository, articleRepo *repository.ArticleRepository) *CategoryService {
	return &CategoryService{repo: repo, articleRepo: articleRepo}
}

// Create 创建分类，校验名称唯一性
func (s *CategoryService) Create(req *model.CreateCategoryRequest) (*model.Category, error) {
	// 检查分类名是否已存在
	if _, err := s.repo.FindByName(req.Name); err == nil {
		return nil, errors.New(errcode.GetMessage(errcode.ErrCategoryAlreadyExists))
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	category := &model.Category{
		Name:      req.Name,
		ParentID:  req.ParentID,
		SortOrder: req.SortOrder,
	}

	if err := s.repo.Create(category); err != nil {
		return nil, errors.New(errcode.GetMessage(errcode.ErrCategoryCreateFailed))
	}

	return category, nil
}

// GetByID 根据 ID 获取分类
func (s *CategoryService) GetByID(id uint) (*model.Category, error) {
	category, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(errcode.GetMessage(errcode.ErrCategoryNotFound))
		}
		return nil, err
	}
	return category, nil
}

// List 获取所有分类列表
func (s *CategoryService) List() ([]model.Category, error) {
	return s.repo.List()
}

// Update 更新分类
func (s *CategoryService) Update(id uint, req *model.UpdateCategoryRequest) (*model.Category, error) {
	category, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(errcode.GetMessage(errcode.ErrCategoryNotFound))
		}
		return nil, err
	}

	if req.Name != "" {
		category.Name = req.Name
	}
	if req.ParentID != nil {
		category.ParentID = req.ParentID
	}
	category.SortOrder = req.SortOrder

	if err := s.repo.Update(category); err != nil {
		return nil, err
	}
	return category, nil
}

// Delete 删除分类，如果分类下存在关联文章则拒绝删除，防止级联影响
func (s *CategoryService) Delete(id uint) error {
	// 检查分类是否存在
	if _, err := s.repo.FindByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New(errcode.GetMessage(errcode.ErrCategoryNotFound))
		}
		return err
	}

	// 级联保护：检查分类下是否有关联文章
	count, err := s.articleRepo.CountByCategoryID(id)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New(errcode.GetMessage(errcode.ErrCategoryHasArticles))
	}

	return s.repo.Delete(id)
}
