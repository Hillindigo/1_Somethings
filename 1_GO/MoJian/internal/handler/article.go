package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/goagent/mojian/internal/model"
	"github.com/goagent/mojian/internal/service"
	"github.com/goagent/mojian/pkg/response"
)

// ArticleHandler 文章处理器
type ArticleHandler struct {
	svc *service.ArticleService
}

// NewArticleHandler 创建文章 Handler 实例
func NewArticleHandler(svc *service.ArticleService) *ArticleHandler {
	return &ArticleHandler{svc: svc}
}

// Create 创建文章
// POST /api/v1/articles
func (h *ArticleHandler) Create(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req model.CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithBadRequest(c, err.Error())
		return
	}

	article, err := h.svc.Create(userID, &req)
	if err != nil {
		response.FailWithServerError(c, err.Error())
		return
	}

	response.Success(c, article)
}

// GetByID 获取文章详情
// GET /api/v1/articles/:id
func (h *ArticleHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.FailWithBadRequest(c, "无效的文章 ID")
		return
	}

	article, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.FailWithNotFound(c, err.Error())
		return
	}

	response.Success(c, article)
}

// List 获取文章列表（分页）
// GET /api/v1/articles
func (h *ArticleHandler) List(c *gin.Context) {
	var req model.ArticleListRequest
	// 设置默认分页参数
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithBadRequest(c, err.Error())
		return
	}
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}

	articles, total, err := h.svc.List(&req)
	if err != nil {
		response.FailWithServerError(c, err.Error())
		return
	}

	response.SuccessWithPage(c, articles, total, req.Page, req.PageSize)
}

// Update 更新文章
// PUT /api/v1/articles/:id
func (h *ArticleHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.FailWithBadRequest(c, "无效的文章 ID")
		return
	}

	userID := c.GetUint("user_id")

	var req model.UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithBadRequest(c, err.Error())
		return
	}

	article, err := h.svc.Update(uint(id), userID, &req)
	if err != nil {
		response.FailWithServerError(c, err.Error())
		return
	}

	response.Success(c, article)
}

// Delete 删除文章
// DELETE /api/v1/articles/:id
func (h *ArticleHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.FailWithBadRequest(c, "无效的文章 ID")
		return
	}

	userID := c.GetUint("user_id")

	if err := h.svc.Delete(uint(id), userID); err != nil {
		response.FailWithServerError(c, err.Error())
		return
	}

	response.Success(c, nil)
}
