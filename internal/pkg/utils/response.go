package utils

import (
	"net/http"
	"zhku-oj/internal/pkg/errors"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	TraceID string      `json:"trace_id,omitempty"`
}

// PageResponse 分页响应结构
type PageResponse struct {
	Response
	Pagination Pagination `json:"pagination"`
}

// Pagination 分页信息
type Pagination struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// SendSuccess 成功响应 - 新版本（根据用户需求）
func SendSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    errors.SUCCESS,
		Message: "成功",
		Data:    data,
	})
}

// SendError 错误响应 - 新版本（根据用户需求）
func SendError(c *gin.Context, errCode int) {
	c.JSON(http.StatusOK, Response{
		Code:    errCode,
		Message: errors.GetErrorMessage(errCode),
	})
}

// SendErrorWithDetail 带详情的错误响应
func SendErrorWithDetail(c *gin.Context, errCode int, detail string) {
	response := Response{
		Code:    errCode,
		Message: errors.GetErrorMessage(errCode),
	}

	if detail != "" {
		response.Data = gin.H{"detail": detail}
	}

	c.JSON(http.StatusOK, response)
}

// SendBusinessError 业务错误响应
func SendBusinessError(c *gin.Context, err *errors.BusinessError) {
	response := Response{
		Code:    err.GetCode(),
		Message: err.GetMessage(),
	}

	if detail := err.GetDetail(); detail != "" {
		response.Data = gin.H{"detail": detail}
	}

	c.JSON(http.StatusOK, response)
}

// SuccessResponse 成功响应 - 兼容老版本
func SuccessResponse(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    errors.SUCCESS,
		Message: message,
		Data:    data,
	})
}

// SuccessResponseWithPagination 分页成功响应
func SuccessResponseWithPagination(c *gin.Context, data interface{}, pagination Pagination) {
	c.JSON(http.StatusOK, PageResponse{
		Response: Response{
			Code:    errors.SUCCESS,
			Message: "success",
			Data:    data,
		},
		Pagination: pagination,
	})
}

// ErrorResponse 错误响应 - 兼容老版本
func ErrorResponse(c *gin.Context, statusCode int, message string, detail string) {
	response := Response{
		Code:    statusCode,
		Message: message,
	}

	if detail != "" {
		response.Data = gin.H{"detail": detail}
	}

	c.JSON(statusCode, response)
}

// ValidationErrorResponse 参数验证错误响应
func ValidationErrorResponse(c *gin.Context, validationErrors map[string]string) {
	c.JSON(http.StatusOK, Response{
		Code:    errors.INVALID_PARAMS,
		Message: errors.GetErrorMessage(errors.INVALID_PARAMS),
		Data: gin.H{
			"errors": validationErrors,
		},
	})
}

// HandleError 统一错误处理函数
func HandleError(c *gin.Context, err error) {
	if bizErr, ok := errors.GetBusinessError(err); ok {
		SendBusinessError(c, bizErr)
		return
	}

	// 其他类型错误作为系统错误处理
	SendErrorWithDetail(c, errors.SYSTEM_ERROR, err.Error())
}

// SendSuccessWithPagination 带分页的成功响应（简化版）
func SendSuccessWithPagination(c *gin.Context, data interface{}, page, pageSize int, total int64) {
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	c.JSON(http.StatusOK, PageResponse{
		Response: Response{
			Code:    errors.SUCCESS,
			Message: "成功",
			Data:    data,
		},
		Pagination: Pagination{
			Page:       page,
			PageSize:   pageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	})
}
