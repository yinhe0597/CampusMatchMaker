package response

import (
	"net/http"
	"time"

	"campus_collab/pkg/errors"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	Timestamp int64       `json:"timestamp"`
}

// PageData 分页数据结构
type PageData struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

// OK 成功响应
func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:      0,
		Message:   "success",
		Data:      data,
		Timestamp: time.Now().Unix(),
	})
}

// OKWithMsg 成功响应（带自定义消息）
func OKWithMsg(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, Response{
		Code:      0,
		Message:   msg,
		Data:      nil,
		Timestamp: time.Now().Unix(),
	})
}

// Created 创建成功响应
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Code:      0,
		Message:   "success",
		Data:      data,
		Timestamp: time.Now().Unix(),
	})
}

// Page 分页响应
func Page(c *gin.Context, list interface{}, total int64, page, pageSize int) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data: PageData{
			List:     list,
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		},
		Timestamp: time.Now().Unix(),
	})
}

// BadRequest 参数错误
func BadRequest(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, Response{
		Code:      1001,
		Message:   msg,
		Data:      nil,
		Timestamp: time.Now().Unix(),
	})
}

// Unauthorized 未认证
func Unauthorized(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, Response{
		Code:      1100,
		Message:   "未登录或登录已过期",
		Data:      nil,
		Timestamp: time.Now().Unix(),
	})
}

// Forbidden 无权限
func Forbidden(c *gin.Context) {
	c.JSON(http.StatusForbidden, Response{
		Code:      1106,
		Message:   "无操作权限",
		Data:      nil,
		Timestamp: time.Now().Unix(),
	})
}

// NotFound 资源不存在
func NotFound(c *gin.Context, msg string) {
	c.JSON(http.StatusNotFound, Response{
		Code:      1002,
		Message:   msg,
		Data:      nil,
		Timestamp: time.Now().Unix(),
	})
}

// ServerError 服务器内部错误
func ServerError(c *gin.Context, msg string) {
	c.JSON(http.StatusInternalServerError, Response{
		Code:      1000,
		Message:   msg,
		Data:      nil,
		Timestamp: time.Now().Unix(),
	})
}

// HandleError 根据 AppError 类型自动返回对应状态码
func HandleError(c *gin.Context, err error) {
	if appErr, ok := err.(*errors.AppError); ok {
		status := codeToHTTPStatus(appErr.Code)
		c.JSON(status, Response{
			Code:      appErr.Code,
			Message:   appErr.Message,
			Data:      nil,
			Timestamp: time.Now().Unix(),
		})
		return
	}
	// 未知错误
	c.JSON(http.StatusInternalServerError, Response{
		Code:      1000,
		Message:   "服务器内部错误",
		Data:      nil,
		Timestamp: time.Now().Unix(),
	})
}

// codeToHTTPStatus 根据业务错误码映射 HTTP 状态码
func codeToHTTPStatus(code int) int {
	switch {
	case code == 1106:
		return http.StatusForbidden
	case code >= 1100 && code <= 1199:
		return http.StatusUnauthorized
	case code >= 1200 && code <= 1299:
		return http.StatusBadRequest
	case code == 1002:
		return http.StatusNotFound
	default:
		return http.StatusBadRequest
	}
}
