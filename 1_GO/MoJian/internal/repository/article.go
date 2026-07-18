package repository

import (
	"github.com/goagent/mojian/internal/model"
	"gorm.io/gorm"
)

// ArticleRepository 文章数据访问层，封装文章相关的数据库操作
type ArticleRepository struct {
	db *gorm.DB
}

// NewArticleRepository 创建文章 Repository 实例
func NewArticleRepository(db *gorm.DB) *ArticleRepository {
	return &ArticleRepository{db: db}
}

// Create 创建文章，同时创建文章与标签的关联关系
func (r *ArticleRepository) Create(article *model.Article) error {
	return r.db.Create(article).Error
}

// FindByID 根据 ID 查询文章，预加载关联的用户、分类和标签
func (r *ArticleRepository) FindByID(id uint) (*model.Article, error) {
	var article model.Article
	if err := r.db.Preload("User").Preload("Category").Preload("Tags").First(&article, id).Error; err != nil {
		return nil, err
	}
	return &article, nil
}

// Update 更新文章信息
func (r *ArticleRepository) Update(article *model.Article) error {
	return r.db.Save(article).Error
}

// Delete 根据 ID 删除文章
func (r *ArticleRepository) Delete(id uint) error {
	return r.db.Delete(&model.Article{}, id).Error
}

// List 分页查询文章列表，支持按分类、标签、状态、作者筛选
// 返回文章列表和总记录数
func (r *ArticleRepository) List(req *model.ArticleListRequest) ([]model.Article, int64, error) {
	var articles []model.Article
	var total int64

	query := r.db.Model(&model.Article{})

	// 条件筛选
	if req.CategoryID > 0 {
		query = query.Where("category_id = ?", req.CategoryID)
	}
	if req.Status >= 0 {
		query = query.Where("status = ?", req.Status)
	}
	if req.UserID > 0 {
		query = query.Where("user_id = ?", req.UserID)
	}
	if req.TagID > 0 {
		// 通过标签筛选需要关联 article_tags 表
		query = query.Joins("JOIN article_tags ON article_tags.article_id = articles.id").
			Where("article_tags.tag_id = ?", req.TagID)
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询，预加载关联数据
	offset := (req.Page - 1) * req.PageSize
	if err := query.Preload("User").Preload("Category").Preload("Tags").
		Order("created_at DESC").
		Offset(offset).Limit(req.PageSize).
		Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

// ReplaceTags 替换文章的标签关联（先删除旧关联，再创建新关联）
func (r *ArticleRepository) ReplaceTags(articleID uint, tags []model.Tag) error {
	article := model.Article{ID: articleID}
	return r.db.Model(&article).Association("Tags").Replace(tags)
}

// IncrementViewCount 增加文章浏览次数
func (r *ArticleRepository) IncrementViewCount(id uint) error {
	return r.db.Model(&model.Article{}).Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}

// CountByCategoryID 统计指定分类下的文章数量，用于删除分类前的级联保护检查
func (r *ArticleRepository) CountByCategoryID(categoryID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&model.Article{}).Where("category_id = ?", categoryID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
