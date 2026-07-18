package repository

import (
	"github.com/goagent/mojian/internal/model"
	"gorm.io/gorm"
)

// CategoryRepository 分类数据访问层
type CategoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository 创建分类 Repository 实例
func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// Create 创建分类
func (r *CategoryRepository) Create(category *model.Category) error {
	return r.db.Create(category).Error
}

// FindByID 根据 ID 查询分类
func (r *CategoryRepository) FindByID(id uint) (*model.Category, error) {
	var category model.Category
	if err := r.db.First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

// FindByName 根据名称查询分类
func (r *CategoryRepository) FindByName(name string) (*model.Category, error) {
	var category model.Category
	if err := r.db.Where("name = ?", name).First(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

// List 查询所有分类，按排序权重升序排列
func (r *CategoryRepository) List() ([]model.Category, error) {
	var categories []model.Category
	if err := r.db.Order("sort_order ASC").Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// Update 更新分类
func (r *CategoryRepository) Update(category *model.Category) error {
	return r.db.Save(category).Error
}

// Delete 根据 ID 删除分类
func (r *CategoryRepository) Delete(id uint) error {
	return r.db.Delete(&model.Category{}, id).Error
}
