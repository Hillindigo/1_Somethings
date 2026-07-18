package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/goagent/mojian/internal/model"
	"github.com/goagent/mojian/internal/service"
	"github.com/goagent/mojian/pkg/captcha"
	"github.com/goagent/mojian/pkg/errcode"
	"github.com/goagent/mojian/pkg/response"
)

// UserHandler 用户处理器，处理用户相关的 HTTP 请求
type UserHandler struct {
	svc *service.UserService
}

// NewUserHandler 创建用户 Handler 实例
func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

// Register 用户注册接口
// POST /api/v1/auth/register
func (h *UserHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithBadRequest(c, err.Error())
		return
	}

	// 校验验证码
	if !captcha.Verify(req.CaptchaID, req.CaptchaCode, true) {
		response.Fail(c, http.StatusBadRequest, errcode.ErrCaptchaInvalid, errcode.GetMessage(errcode.ErrCaptchaInvalid))
		return
	}

	user, err := h.svc.Register(&req)
	if err != nil {
		response.FailWithServerError(c, err.Error())
		return
	}

	response.Success(c, user)
}

// Login 用户登录接口
// POST /api/v1/auth/login
func (h *UserHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithBadRequest(c, err.Error())
		return
	}

	// 校验验证码
	if !captcha.Verify(req.CaptchaID, req.CaptchaCode, true) {
		response.Fail(c, http.StatusBadRequest, errcode.ErrCaptchaInvalid, errcode.GetMessage(errcode.ErrCaptchaInvalid))
		return
	}

	result, err := h.svc.Login(&req)
	if err != nil {
		response.Fail(c, http.StatusUnauthorized, 401, err.Error())
		return
	}

	response.Success(c, result)
}

// GetUser 获取当前登录用户信息
// GET /api/v1/user/profile
func (h *UserHandler) GetUser(c *gin.Context) {
	userID := c.GetUint("user_id")
	user, err := h.svc.GetUser(userID)
	if err != nil {
		response.FailWithNotFound(c, err.Error())
		return
	}

	response.Success(c, user)
}

// UpdateUser 更新当前登录用户信息
// PUT /api/v1/user/profile
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req model.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithBadRequest(c, err.Error())
		return
	}

	user, err := h.svc.UpdateUser(userID, &req)
	if err != nil {
		response.FailWithServerError(c, err.Error())
		return
	}

	response.Success(c, user)
}

// ListUsers 获取用户列表（管理员）
// GET /api/v1/admin/users
func (h *UserHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	users, total, err := h.svc.ListUsers(page, pageSize)
	if err != nil {
		response.FailWithServerError(c, err.Error())
		return
	}

	response.SuccessWithPage(c, users, total, page, pageSize)
}

// UpdateUserRole 更新用户角色（管理员）
// PUT /api/v1/admin/users/:id/role
// 禁止管理员降级自己的角色，防止误操作导致失去管理权限
func (h *UserHandler) UpdateUserRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.FailWithBadRequest(c, "无效的用户 ID")
		return
	}

	// 自我保护：禁止管理员降级自己
	currentUserID := c.GetUint("user_id")
	if uint(id) == currentUserID {
		response.FailWithForbidden(c, "不能修改自己的角色")
		return
	}

	var req struct {
		Role int8 `json:"role" binding:"oneof=0 1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithBadRequest(c, err.Error())
		return
	}

	user, err := h.svc.UpdateUserRole(uint(id), req.Role)
	if err != nil {
		response.FailWithServerError(c, err.Error())
		return
	}

	response.Success(c, user)
}

// DeleteUser 删除用户（管理员）
// DELETE /api/v1/admin/users/:id
// 禁止管理员删除自己，防止误操作导致无法管理
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.FailWithBadRequest(c, "无效的用户 ID")
		return
	}

	// 自我保护：禁止管理员删除自己
	currentUserID := c.GetUint("user_id")
	if uint(id) == currentUserID {
		response.FailWithForbidden(c, "不能删除自己")
		return
	}

	if err := h.svc.DeleteUser(uint(id)); err != nil {
		response.FailWithServerError(c, err.Error())
		return
	}

	response.Success(c, nil)
}
