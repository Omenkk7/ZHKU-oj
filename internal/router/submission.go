package router

import (
	"zhku-oj/internal/middleware"

	"github.com/gin-gonic/gin"
)

// setupSubmissionRoutes 设置代码提交相关路由
// 代码提交、判题结果查询等功能
func (rm *RouterManager) setupSubmissionRoutes(v1 *gin.RouterGroup) {
	submissionGroup := v1.Group("/submissions")
	submissionGroup.Use(middleware.AuthRequired()) // 所有提交接口都需要认证
	{
		// ========== 代码提交接口 ==========

		// 提交代码进行判题
		// POST /api/v1/submissions
		// 请求体: {"problem_id": "xxx", "language": "java", "code": "..."}
		// 响应码: 0-成功, 10002-参数错误, 30001-题目不存在, 40004-代码过长, 40007-重复提交
		submissionGroup.POST("", rm.submissionHandler.Submit)

		// ========== 提交记录查询 ==========

		// 获取提交详情
		// GET /api/v1/submissions/{id}
		// 响应码: 0-成功, 10002-参数错误, 40001-提交记录不存在, 40008-提交访问被拒绝
		submissionGroup.GET("/:id", rm.submissionHandler.GetSubmission)

		// 获取提交列表（当前用户）
		// GET /api/v1/submissions?page=1&page_size=20&problem_id=xxx&status=ACCEPTED&language=java
		// 响应码: 0-成功, 10002-参数错误
		submissionGroup.GET("", rm.submissionHandler.ListSubmissions)

		// 获取提交代码（需要是提交者本人或管理员）
		// GET /api/v1/submissions/{id}/code
		// 响应码: 0-成功, 10002-参数错误, 40001-提交记录不存在, 40008-提交访问被拒绝
		// submissionGroup.GET("/:id/code", rm.submissionHandler.GetSubmissionCode)

		// 重新判题（管理员权限）
		// POST /api/v1/submissions/{id}/rejudge
		// 权限: admin
		// 响应码: 0-成功, 10002-参数错误, 10004-权限不足, 40001-提交记录不存在
		// submissionGroup.POST("/:id/rejudge",
		// 	middleware.RoleRequired("admin"),
		// 	rm.submissionHandler.RejudgeSubmission)

		// ========== 实时判题状态 ==========

		// WebSocket连接获取实时判题状态
		// GET /api/v1/submissions/{id}/status (升级为WebSocket)
		// 用于实时推送判题进度和结果
		// submissionGroup.GET("/:id/status", rm.submissionHandler.GetSubmissionStatus)

		// 获取判题队列状态
		// GET /api/v1/submissions/queue/status
		// 权限: teacher, admin
		// 响应码: 0-成功, 10004-权限不足
		// submissionGroup.GET("/queue/status",
		// 	middleware.RoleRequired("teacher", "admin"),
		// 	rm.submissionHandler.GetJudgeQueueStatus)
	}
}
