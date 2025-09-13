package examples

import (
	"zhku-oj/internal/pkg/errors"
	"zhku-oj/internal/pkg/utils"

	"github.com/gin-gonic/gin"
)

// ExampleHandler 展示新响应体系统的使用示例
type ExampleHandler struct{}

// NewExampleHandler 创建示例控制器
func NewExampleHandler() *ExampleHandler {
	return &ExampleHandler{}
}

// SuccessExample 成功响应示例
// GET /api/v1/examples/success
func (h *ExampleHandler) SuccessExample(c *gin.Context) {
	data := gin.H{
		"message": "这是一个成功的响应示例",
		"user": gin.H{
			"id":       "507f1f77bcf86cd799439011",
			"username": "zhangsan",
			"email":    "zhangsan@example.com",
		},
	}

	// 使用新版本的 SendSuccess
	utils.SendSuccess(c, data)

	/* 响应格式：
	{
		"code": 0,
		"message": "成功",
		"data": {
			"message": "这是一个成功的响应示例",
			"user": {
				"id": "507f1f77bcf86cd799439011",
				"username": "zhangsan",
				"email": "zhangsan@example.com"
			}
		}
	}
	*/
}

// ErrorExample 错误响应示例
// GET /api/v1/examples/error
func (h *ExampleHandler) ErrorExample(c *gin.Context) {
	// 方式1：直接使用错误码
	utils.SendError(c, errors.USER_NOT_FOUND)

	/* 响应格式：
	{
		"code": 20001,
		"message": "用户不存在"
	}
	*/
}

// ErrorWithDetailExample 带详情的错误响应示例
// GET /api/v1/examples/error-detail
func (h *ExampleHandler) ErrorWithDetailExample(c *gin.Context) {
	// 方式2：带详情的错误响应
	utils.SendErrorWithDetail(c, errors.INVALID_PARAMS, "用户ID格式不正确")

	/* 响应格式：
	{
		"code": 10002,
		"message": "参数错误",
		"data": {
			"detail": "用户ID格式不正确"
		}
	}
	*/
}

// BusinessErrorExample 业务错误响应示例
// GET /api/v1/examples/business-error
func (h *ExampleHandler) BusinessErrorExample(c *gin.Context) {
	// 方式3：使用业务错误对象
	err := errors.NewDuplicateSubmission("相同代码在10分钟内已提交")
	utils.SendBusinessError(c, err)

	/* 响应格式：
	{
		"code": 40007,
		"message": "请勿重复提交相同代码",
		"data": {
			"detail": "相同代码在10分钟内已提交"
		}
	}
	*/
}

// HandleErrorExample 统一错误处理示例
// GET /api/v1/examples/handle-error
func (h *ExampleHandler) HandleErrorExample(c *gin.Context) {
	// 模拟一个业务错误
	err := errors.NewUserNotFound("用户ID: 507f1f77bcf86cd799439011")

	// 使用统一错误处理函数
	utils.HandleError(c, err)

	/* 响应格式：
	{
		"code": 20001,
		"message": "用户不存在",
		"data": {
			"detail": "用户ID: 507f1f77bcf86cd799439011"
		}
	}
	*/
}

// PaginationExample 分页响应示例
// GET /api/v1/examples/pagination
func (h *ExampleHandler) PaginationExample(c *gin.Context) {
	users := []gin.H{
		{
			"id":       "507f1f77bcf86cd799439011",
			"username": "zhangsan",
			"email":    "zhangsan@example.com",
		},
		{
			"id":       "507f1f77bcf86cd799439012",
			"username": "lisi",
			"email":    "lisi@example.com",
		},
	}

	// 使用分页响应
	utils.SendSuccessWithPagination(c, users, 1, 20, 100)

	/* 响应格式：
	{
		"code": 0,
		"message": "成功",
		"data": [
			{
				"id": "507f1f77bcf86cd799439011",
				"username": "zhangsan",
				"email": "zhangsan@example.com"
			},
			{
				"id": "507f1f77bcf86cd799439012",
				"username": "lisi",
				"email": "lisi@example.com"
			}
		],
		"pagination": {
			"page": 1,
			"page_size": 20,
			"total": 100,
			"total_pages": 5
		}
	}
	*/
}

// ValidationErrorExample 参数验证错误示例
// POST /api/v1/examples/validation
func (h *ExampleHandler) ValidationErrorExample(c *gin.Context) {
	// 模拟参数验证失败
	validationErrors := map[string]string{
		"username": "用户名不能为空",
		"email":    "邮箱格式不正确",
		"password": "密码长度至少6位",
	}

	utils.ValidationErrorResponse(c, validationErrors)

	/* 响应格式：
	{
		"code": 10002,
		"message": "参数错误",
		"data": {
			"errors": {
				"username": "用户名不能为空",
				"email": "邮箱格式不正确",
				"password": "密码长度至少6位"
			}
		}
	}
	*/
}

// ServiceLayerExample 服务层使用示例
func (h *ExampleHandler) ServiceLayerExample() {
	// 在服务层中创建业务错误的示例

	// 1. 简单错误
	_ = errors.NewUserNotFound()

	// 2. 带详情的错误
	_ = errors.NewInvalidPassword("密码长度不足6位")

	// 3. 自定义错误
	_ = errors.New(errors.JUDGE_TIMEOUT, "判题服务响应超时，请稍后重试")

	// 4. 格式化错误
	_ = errors.Newf(errors.CODE_TOO_LONG, "代码长度%d超过限制%d", 10000, 5000)

	// 5. 包装已有错误
	// originalErr := someOperation()
	// err := errors.Wrap(errors.DATABASE_ERROR, originalErr)
}

// 前端JavaScript使用示例
/*
// 前端处理响应的统一函数
function handleApiResponse(response) {
    if (response.code === 0) {
        // 成功处理
        console.log('Success:', response.data);
        return response.data;
    } else {
        // 错误处理
        const errorMessage = response.message;
        const errorDetail = response.data?.detail;

        // 根据错误码进行不同处理
        switch (response.code) {
            case 10003: // UNAUTHORIZED
                // 跳转登录页
                window.location.href = '/login';
                break;
            case 20001: // USER_NOT_FOUND
                showError('用户不存在');
                break;
            case 40007: // DUPLICATE_SUBMISSION
                showWarning('请勿重复提交相同代码');
                break;
            case 50006: // TIME_LIMIT_EXCEEDED
                showInfo('代码执行超时');
                break;
            default:
                // 显示通用错误提示
                showError(errorMessage + (errorDetail ? ': ' + errorDetail : ''));
        }

        throw new Error(errorMessage);
    }
}

// 使用示例
fetch('/api/v1/users/profile')
    .then(response => response.json())
    .then(handleApiResponse)
    .then(userData => {
        console.log('User data:', userData);
    })
    .catch(error => {
        console.error('Error:', error.message);
    });
*/
