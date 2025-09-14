package request

import (
	"zhku-oj/pkg/io/constanct"
)

/*
@Author: omenkk7
@Date: 2025/9/14 16:01
@Description: 提交解题请求模型
*/

type AddSubmitReq struct {
	ProblemID   string `json:"ProblemID"`   //题目ID
	UserID      string `json:"UserID"`      //用户ID
	Code        string `json:"Code"`        //源代码
	Lang        string `json:"language"`    //语言
	SubmittedAt int64  `json:"SubmittedAt"` //提交时间
}
