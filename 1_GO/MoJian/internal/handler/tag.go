package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/goagent/mojian/internal/model"
	"github.com/goagent/mojian/internal/service"
	"github.com/goagent/mojian/pkg/response"
)

// TagHandler 标签处理器
type TagHandler struct {
	svc *service.TagService
}

// NewTagHandler 创建标签 Handler 实例
func NewTagHandler(svc *service.TagService) *TagHandler {
	return &TagHandler{svc: svc}
}

// Create 创建标签
// POST /api/v1/tags
func (h *TagHandler) Create(c *gin.Context) {
	var req model.CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithBadRequest(c, err.Error())
		return
	}

	tag, err := h.svc.Create(&req)
	if err != nil {
		response.FailWithServerError(c, err.Error())
		return
	}

	response.Success(c, tag)
}

// GetByID 获取标签详情
// GET /api/v1/tags/:id
func (h *TagHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.FailWithBadRequest(c, "无效的标签 ID")
		return
	}

	tag, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.FailWithNotFound(c, err.Error())
		return
	}

	response.Success(c, tag)
}

// List 获取所有标签
// GET /api/v1/tags
func (h *TagHandler) List(c *gin.Context) {
	tags, err := h.svc.List()
	if err != nil {
		response.FailWithServerError(c, err.Error())
		return
	}

	response.Success(c, tags)
}

// Update 更新标签
// PUT /api/v1/tags/:id
func (h *TagHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.FailWithBadRequest(c, "无效的标签 ID")
		return
	}

	var req model.UpdateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithBadRequest(c, err.Error())
		return
	}

	tag, err := h.svc.Update(uint(id), &req)
	if err != nil {
		response.FailWithServerError(c, err.Error())
		return
	}

	response.Success(c, tag)
}

// Delete 删除标签
// DELETE /api/v1/tags/:id
func (h *TagHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.FailWithBadRequest(c, "无效的标签 ID")
		return
	}

	if err := h.svc.Delete(uint(id)); err != nil {
		response.FailWithServerError(c, err.Error())
		return
	}

	response.Success(c, nil)
}
