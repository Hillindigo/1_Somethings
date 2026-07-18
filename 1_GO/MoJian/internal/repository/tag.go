package repository

import (
	"github.com/goagent/mojian/internal/model"
	"gorm.io/gorm"
)

// TagRepository 标签数据访问层
type TagRepository struct {
	db *gorm.DB
}

// NewTagRepository 创建标签 Repository 实例
func NewTagRepository(db *gorm.DB) *TagRepository {
	return &TagRepository{db: db}
}

// Create 创建标签
func (r *TagRepository) Create(tag *model.Tag) error {
	return r.db.Create(tag).Error
}

// FindByID 根据 ID 查询标签
func (r *TagRepository) FindByID(id uint) (*model.Tag, error) {
	var tag model.Tag
	if err := r.db.First(&tag, id).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

// FindByName 根据名称查询标签
func (r *TagRepository) FindByName(name string) (*model.Tag, error) {
	var tag model.Tag
	if err := r.db.Where("name = ?", name).First(&tag).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

// FindByIDs 根据 ID 列表批量查询标签
func (r *TagRepository) FindByIDs(ids []uint) ([]model.Tag, error) {
	var tags []model.Tag
	if err := r.db.Where("id IN ?", ids).Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

// List 查询所有标签
func (r *TagRepository) List() ([]model.Tag, error) {
	var tags []model.Tag
	if err := r.db.Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

// Update 更新标签
func (r *TagRepository) Update(tag *model.Tag) error {
	return r.db.Save(tag).Error
}

// Delete 根据 ID 删除标签
func (r *TagRepository) Delete(id uint) error {
	return r.db.Delete(&model.Tag{}, id).Error
}
