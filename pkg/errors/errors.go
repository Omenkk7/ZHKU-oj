package errors

import (
	"fmt"
)

// BusinessError 业务错误结构体
type BusinessError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

// Error 实现error接口
func (e *BusinessError) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("Code: %d, Message: %s, Detail: %s", e.Code, e.Message, e.Detail)
	}
	return fmt.Sprintf("Code: %d, Message: %s", e.Code, e.Message)
}

// GetCode 获取错误码
func (e *BusinessError) GetCode() int {
	return e.Code
}

// GetMessage 获取错误消息
func (e *BusinessError) GetMessage() string {
	return e.Message
}

// GetDetail 获取错误详情
func (e *BusinessError) GetDetail() string {
	return e.Detail
}

// New 创建业务错误
func New(code int, detail ...string) *BusinessError {
	err := &BusinessError{
		Code:    code,
		Message: GetErrorMessage(code),
	}

	if len(detail) > 0 {
		err.Detail = detail[0]
	}

	return err
}

// Newf 创建带格式化详情的业务错误
func Newf(code int, format string, args ...interface{}) *BusinessError {
	return &BusinessError{
		Code:    code,
		Message: GetErrorMessage(code),
		Detail:  fmt.Sprintf(format, args...),
	}
}

// Wrap 包装已有错误
func Wrap(code int, err error) *BusinessError {
	return &BusinessError{
		Code:    code,
		Message: GetErrorMessage(code),
		Detail:  err.Error(),
	}
}

// 预定义常用错误实例

// 通用错误
func NewSystemError(detail ...string) *BusinessError {
	return New(SYSTEM_ERROR, detail...)
}

func NewInvalidParams(detail ...string) *BusinessError {
	return New(INVALID_PARAMS, detail...)
}

func NewUnauthorized(detail ...string) *BusinessError {
	return New(UNAUTHORIZED, detail...)
}

func NewForbidden(detail ...string) *BusinessError {
	return New(FORBIDDEN, detail...)
}

func NewNotFound(detail ...string) *BusinessError {
	return New(NOT_FOUND, detail...)
}

// 用户模块错误
func NewUserNotFound(detail ...string) *BusinessError {
	return New(USER_NOT_FOUND, detail...)
}

func NewUserAlreadyExists(detail ...string) *BusinessError {
	return New(USER_ALREADY_EXISTS, detail...)
}

func NewUsernameAlreadyExists(detail ...string) *BusinessError {
	return New(USERNAME_ALREADY_EXISTS, detail...)
}

func NewEmailAlreadyExists(detail ...string) *BusinessError {
	return New(EMAIL_ALREADY_EXISTS, detail...)
}

func NewInvalidPassword(detail ...string) *BusinessError {
	return New(INVALID_PASSWORD, detail...)
}

func NewUserDisabled(detail ...string) *BusinessError {
	return New(USER_DISABLED, detail...)
}

func NewLoginFailed(detail ...string) *BusinessError {
	return New(LOGIN_FAILED, detail...)
}

// 题目模块错误
func NewProblemNotFound(detail ...string) *BusinessError {
	return New(PROBLEM_NOT_FOUND, detail...)
}

func NewProblemAccessDenied(detail ...string) *BusinessError {
	return New(PROBLEM_ACCESS_DENIED, detail...)
}

// 提交模块错误
func NewSubmissionNotFound(detail ...string) *BusinessError {
	return New(SUBMISSION_NOT_FOUND, detail...)
}

func NewCodeTooLong(detail ...string) *BusinessError {
	return New(CODE_TOO_LONG, detail...)
}

func NewCodeEmpty(detail ...string) *BusinessError {
	return New(CODE_EMPTY, detail...)
}

func NewDuplicateSubmission(detail ...string) *BusinessError {
	return New(DUPLICATE_SUBMISSION, detail...)
}

func NewSubmissionTooFrequent(detail ...string) *BusinessError {
	return New(SUBMISSION_TOO_FREQUENT, detail...)
}

// 判题模块错误
func NewJudgeSystemError(detail ...string) *BusinessError {
	return New(JUDGE_SYSTEM_ERROR, detail...)
}

func NewJudgeTimeout(detail ...string) *BusinessError {
	return New(JUDGE_TIMEOUT, detail...)
}

func NewSandboxError(detail ...string) *BusinessError {
	return New(SANDBOX_ERROR, detail...)
}

// 管理模块错误
func NewAdminPermissionDenied(detail ...string) *BusinessError {
	return New(ADMIN_PERMISSION_DENIED, detail...)
}

func NewDatabaseError(detail ...string) *BusinessError {
	return New(DATABASE_ERROR, detail...)
}

// IsBusinessError 判断是否为业务错误
func IsBusinessError(err error) bool {
	_, ok := err.(*BusinessError)
	return ok
}

// GetBusinessError 获取业务错误
func GetBusinessError(err error) (*BusinessError, bool) {
	bizErr, ok := err.(*BusinessError)
	return bizErr, ok
}
