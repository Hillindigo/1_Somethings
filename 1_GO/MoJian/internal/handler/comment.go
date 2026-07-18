package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/goagent/mojian/internal/config"
	"github.com/goagent/mojian/internal/model"
	"github.com/goagent/mojian/internal/service"
	"github.com/goagent/mojian/pkg/response"
)

// CommentHandler 评论处理器
type CommentHandler struct {
	svc *service.CommentService
}

// NewCommentHandler 创建评论 Handler 实例
func NewCommentHandler(svc *service.CommentService) *CommentHandler {
	return &CommentHandler{svc: svc}
}

// Create 创建评论
// POST /api/v1/comments
func (h *CommentHandler) Create(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req model.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithBadRequest(c, err.Error())
		return
	}

	comment, err := h.svc.Create(userID, &req)
	if err != nil {
		response.FailWithServerError(c, err.Error())
		return
	}

	response.Success(c, comment)
}

// GetByID 获取评论详情
// GET /api/v1/comments/:id
func (h *CommentHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.FailWithBadRequest(c, "无效的评论 ID")
		return
	}

	comment, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.FailWithNotFound(c, err.Error())
		return
	}

	response.Success(c, comment)
}

// ListByArticleID 获取指定文章的评论列表
// GET /api/v1/articles/:article_id/comments
// 游客最多返回 guest_comment_limit 条评论（配置文件可调），登录用户无限制
func (h *CommentHandler) ListByArticleID(c *gin.Context) {
	articleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.FailWithBadRequest(c, "无效的文章 ID")
		return
	}

	comments, err := h.svc.ListByArticleID(uint(articleID))
	if err != nil {
		response.FailWithServerError(c, err.Error())
		return
	}

	// 游客（未登录用户）最多只能查看配置限制条数的评论
	_, isLoggedIn := c.Get("user_id")
	limit := config.GlobalConfig.Blog.GuestCommentLimit
	if limit <= 0 {
		limit = 5 // 兜底默认值
	}
	isLimited := false
	if !isLoggedIn && len(comments) > limit {
		comments = comments[:limit]
		isLimited = true
	}

	// 返回结构化对象，携带游客限制信息供前端动态显示
	response.Success(c, gin.H{
		"comments":            comments,
		"guest_comment_limit": limit,
		"is_limited":          isLimited,
	})
}

// Delete 删除评论
// DELETE /api/v1/comments/:id
func (h *CommentHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.FailWithBadRequest(c, "无效的评论 ID")
		return
	}

	if err := h.svc.Delete(uint(id)); err != nil {
		response.FailWithServerError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// UpdateStatus 更新评论审核状态（管理员）
// PUT /api/v1/comments/:id/status
func (h *CommentHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.FailWithBadRequest(c, "无效的评论 ID")
		return
	}

	var req struct {
		Status int8 `json:"status" binding:"oneof=0 1 2"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithBadRequest(c, err.Error())
		return
	}

	if err := h.svc.UpdateStatus(uint(id), req.Status); err != nil {
		response.FailWithServerError(c, err.Error())
		return
	}

	response.Success(c, nil)
}
