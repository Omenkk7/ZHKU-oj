package interfaces

import (
	"context"
	"zhku-oj/pkg/io/request"
	"zhku-oj/pkg/io/response"
)

/*
@Author: omenkk7
@Date: 2025/9/14 16:13
@Description:
*/

type SubmitService interface {
	//提交接口
	AddSubmit(ctx context.Context, req *request.AddSubmitReq) (*response.AddSubmitResp, error)
}
