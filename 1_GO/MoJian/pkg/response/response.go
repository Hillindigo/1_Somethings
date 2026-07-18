package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构体
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// PageData 分页数据结构
type PageData struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// SuccessWithPage 成功响应（分页）
func SuccessWithPage(c *gin.Context, list interface{}, total int64, page, pageSize int) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data: PageData{
			List:     list,
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		},
	})
}

// Fail 失败响应
func Fail(c *gin.Context, httpCode int, code int, message string) {
	c.JSON(httpCode, Response{
		Code:    code,
		Message: message,
	})
}

// FailWithBadRequest 参数错误响应（400）
func FailWithBadRequest(c *gin.Context, message string) {
	Fail(c, http.StatusBadRequest, 400, message)
}

// FailWithUnauthorized 未认证响应（401）
func FailWithUnauthorized(c *gin.Context, message string) {
	Fail(c, http.StatusUnauthorized, 401, message)
}

// FailWithForbidden 无权限响应（403）
func FailWithForbidden(c *gin.Context, message string) {
	Fail(c, http.StatusForbidden, 403, message)
}

// FailWithNotFound 资源不存在响应（404）
func FailWithNotFound(c *gin.Context, message string) {
	Fail(c, http.StatusNotFound, 404, message)
}

// FailWithServerError 服务器内部错误响应（500）
func FailWithServerError(c *gin.Context, message string) {
	Fail(c, http.StatusInternalServerError, 500, message)
}
