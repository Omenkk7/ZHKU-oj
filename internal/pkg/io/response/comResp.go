package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/**
 * @Author: omenkk7
 * @Date: 2025/9/13
 * @Desc:通用
 */

// Response 统一响应结构
type Response[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

// 通用Response构造函数
func NewResponse[T any](code int, msg string, data T) Response[T] {
	return Response[T]{
		Code:    code,
		Message: msg,
		Data:    data,
	}
}


type RetType string

const (
	SUCCESS   RetType = "success"
	FAIL      RetType = "fail"
	ERROR     RetType = "error"
	NOT_FOUND RetType = "not_found"
)

func ResponseOK(c *gin.Context, resp interface{}) {
	c.JSON(http.StatusOK, resp)
}

func ResponseError(c *gin.Context, code constant.ResCode) {
	c.JSON(code.HttpCode(), Response{
		Code:    int(code),
		Message: code.Msg(),
	})
}

/*
*
创建响应结构函数
*/
func CreateResponse(code constant.ResCode, data interface{}) Response {
	return Response{
		Code:    int(code),
		Message: code.Msg(),
		Data:    data,
	}
}
