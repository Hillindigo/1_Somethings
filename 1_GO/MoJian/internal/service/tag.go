package service

import (
	"errors"

	"github.com/goagent/mojian/internal/model"
	"github.com/goagent/mojian/internal/repository"
	"github.com/goagent/mojian/pkg/errcode"
	"gorm.io/gorm"
)

// TagService 标签业务逻辑层
type TagService struct {
	repo *repository.TagRepository
}

// NewTagService 创建标签 Service 实例
func NewTagService(repo *repository.TagRepository) *TagService {
	return &TagService{repo: repo}
}

// Create 创建标签，校验名称唯一性
func (s *TagService) Create(req *model.CreateTagRequest) (*model.Tag, error) {
	if _, err := s.repo.FindByName(req.Name); err == nil {
		return nil, errors.New(errcode.GetMessage(errcode.ErrTagAlreadyExists))
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	tag := &model.Tag{Name: req.Name}
	if err := s.repo.Create(tag); err != nil {
		return nil, errors.New(errcode.GetMessage(errcode.ErrTagCreateFailed))
	}
	return tag, nil
}

// GetByID 根据 ID 获取标签
func (s *TagService) GetByID(id uint) (*model.Tag, error) {
	tag, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(errcode.GetMessage(errcode.ErrTagNotFound))
		}
		return nil, err
	}
	return tag, nil
}

// List 获取所有标签列表
func (s *TagService) List() ([]model.Tag, error) {
	return s.repo.List()
}

// Update 更新标签
func (s *TagService) Update(id uint, req *model.UpdateTagRequest) (*model.Tag, error) {
	tag, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(errcode.GetMessage(errcode.ErrTagNotFound))
		}
		return nil, err
	}

	tag.Name = req.Name
	if err := s.repo.Update(tag); err != nil {
		return nil, err
	}
	return tag, nil
}

// Delete 删除标签
func (s *TagService) Delete(id uint) error {
	return s.repo.Delete(id)
}
