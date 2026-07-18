package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/goagent/mojian/internal/model"
	"github.com/goagent/mojian/internal/service"
	"github.com/goagent/mojian/pkg/errcode"
	"github.com/goagent/mojian/pkg/response"
)

// CategoryHandler 分类处理器
type CategoryHandler struct {
	svc *service.CategoryService
}

// NewCategoryHandler 创建分类 Handler 实例
func NewCategoryHandler(svc *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{svc: svc}
}

// Create 创建分类
// POST /api/v1/categories
func (h *CategoryHandler) Create(c *gin.Context) {
	var req model.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithBadRequest(c, err.Error())
		return
	}

	category, err := h.svc.Create(&req)
	if err != nil {
		response.FailWithServerError(c, err.Error())
		return
	}

	response.Success(c, category)
}

// GetByID 获取分类详情
// GET /api/v1/categories/:id
func (h *CategoryHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.FailWithBadRequest(c, "无效的分类 ID")
		return
	}

	category, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.FailWithNotFound(c, err.Error())
		return
	}

	response.Success(c, category)
}

// List 获取所有分类
// GET /api/v1/categories
func (h *CategoryHandler) List(c *gin.Context) {
	categories, err := h.svc.List()
	if err != nil {
		response.FailWithServerError(c, err.Error())
		return
	}

	response.Success(c, categories)
}

// Update 更新分类
// PUT /api/v1/categories/:id
func (h *CategoryHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.FailWithBadRequest(c, "无效的分类 ID")
		return
	}

	var req model.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithBadRequest(c, err.Error())
		return
	}

	category, err := h.svc.Update(uint(id), &req)
	if err != nil {
		response.FailWithServerError(c, err.Error())
		return
	}

	response.Success(c, category)
}

// Delete 删除分类，如果分类下存在关联文章则返回错误
// DELETE /api/v1/categories/:id
func (h *CategoryHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.FailWithBadRequest(c, "无效的分类 ID")
		return
	}

	if err := h.svc.Delete(uint(id)); err != nil {
		// 区分"分类不存在"和"分类下有文章"两种错误
		if err.Error() == errcode.GetMessage(errcode.ErrCategoryHasArticles) {
			response.FailWithBadRequest(c, err.Error())
			return
		}
		response.FailWithNotFound(c, err.Error())
		return
	}

	response.Success(c, nil)
}
