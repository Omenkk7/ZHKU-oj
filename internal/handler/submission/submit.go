package submission

import (
	"zhku-oj/internal/service/interfaces"
	"zhku-oj/pkg/io/constanct"
	"zhku-oj/pkg/io/request"
	"zhku-oj/pkg/io/response"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// AddCommit 提交代码接口
// 用户提交解题代码，系统将进行在线判题
// 请求方法: POST
// 路径: /submitq
// 请求体:
// 响应:

type SubmitHandler struct {
	SubmitService interfaces.SubmitService
}

// AddCommit 提交申请接口
func (submitHandler *SubmitHandler) AddCommit(ctx *gin.Context) {
	//logger := utils.GetLogInstance()
	req := new(request.AddSubmitReq)

	//参数校验和绑定
	if err := ctx.ShouldBindWith(req, binding.JSON); err != nil {
		//logger.Errorf("call ShouldBindWith failed, err = %s", err.Error()
		response.ResponseError(ctx, constanct.InvalidParamCode)
		return
	}
	//业务逻辑
	resp, err := submitHandler.SubmitService.AddSubmit(ctx, req)
	if err != nil {
		//logger.Errorf("call AddSubmit failed, req=%+v, err=%s", utils.Sdump(req), err)
		response.ResponseError(ctx, constanct.ServerErrorCode) //服务器错误
		return
	}
	response.ResponseOK(ctx, resp) //返回结果

	return
}
