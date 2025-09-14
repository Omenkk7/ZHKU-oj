package response

import (
	"net/http"

	"zhku-oj/pkg/io/constanct"

	"github.com/gin-gonic/gin"
)

/**
 * @Author: omenkk7
 * @Date: 2025/9/13
 * @Desc:响应体基类
 */

// Response 统一响应结构
type Response[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

type RetType string

const (
	SUCCESS   RetType = "success"
	FAIL      RetType = "fail"
	ERROR     RetType = "error"
	NOT_FOUND RetType = "not_found"
)

// 成功统一用 StatusOK
func ResponseOK[T any](c *gin.Context, data T) {
	c.JSON(http.StatusOK, Response[T]{
		Code:    int(constanct.SuccessCode),
		Message: constanct.SuccessCode.Msg(),
		Data:    data,
	})
}

// 错误返回，data 永远是空
func ResponseError(c *gin.Context, code constanct.ResCode) {
	c.JSON(code.HttpCode(), Response[any]{
		Code:    int(code),
		Message: code.Msg(),
		Data:    nil, // 或 constanct.Empty{}
	})
}
