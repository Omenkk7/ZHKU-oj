package router

import (
	"zhku-oj/internal/middleware"

	"github.com/gin-gonic/gin"
)

// setupAdminRoutes 设置管理员相关路由
// 系统管理、用户管理、数据统计等功能
func (rm *RouterManager) setupAdminRoutes(v1 *gin.RouterGroup) {
	adminGroup := v1.Group("/admin")
	adminGroup.Use(middleware.AuthRequired(), middleware.RoleRequired("admin")) // 管理员权限
	{
		// ========== 系统管理 ==========

		// 管理员仪表板
		// GET /api/v1/admin/dashboard
		// 响应码: 0-成功, 10004-权限不足
		adminGroup.GET("/dashboard", rm.adminHandler.Dashboard)

		// 系统状态监控
		// GET /api/v1/admin/system/status
		// 响应码: 0-成功, 10004-权限不足
		adminGroup.GET("/system/status", rm.adminHandler.SystemStatus)

		// 系统配置管理
		// GET /api/v1/admin/system/config
		// PUT /api/v1/admin/system/config
		// 响应码: 0-成功, 10004-权限不足, 10002-参数错误
		// adminGroup.GET("/system/config", rm.adminHandler.GetSystemConfig)
		// adminGroup.PUT("/system/config", rm.adminHandler.UpdateSystemConfig)

		// ========== 用户管理 CRUD ==========

		// 创建用户
		// POST /api/v1/admin/users
		// 响应码: 0-成功, 10002-参数错误, 20002-用户已存在
		adminGroup.POST("/users", rm.userHandler.CreateUser)

		// 获取用户列表（管理员视图）
		// GET /api/v1/admin/users?page=1&page_size=20&role=student&keyword=张三&is_active=true
		// 响应码: 0-成功, 10002-参数错误
		adminGroup.GET("/users", rm.userHandler.ListUsers)

		// 更新用户信息
		// PUT /api/v1/admin/users/{id}
		// 响应码: 0-成功, 10002-参数错误, 20001-用户不存在
		adminGroup.PUT("/users/:id", rm.userHandler.UpdateUser)

		// 删除用户
		// DELETE /api/v1/admin/users/{id}
		// 响应码: 0-成功, 10002-参数错误, 20001-用户不存在
		adminGroup.DELETE("/users/:id", rm.userHandler.DeleteUser)

		// 激活用户
		// PUT /api/v1/admin/users/{id}/activate
		// 响应码: 0-成功, 10002-参数错误, 20001-用户不存在
		adminGroup.PUT("/users/:id/activate", rm.userHandler.ActivateUser)

		// 停用用户
		// PUT /api/v1/admin/users/{id}/deactivate
		// 响应码: 0-成功, 10002-参数错误, 20001-用户不存在
		adminGroup.PUT("/users/:id/deactivate", rm.userHandler.DeactivateUser)

		// 批量导入用户
		// POST /api/v1/admin/users/import
		// 响应码: 0-成功, 10002-参数错误
		// adminGroup.POST("/users/import", rm.adminHandler.ImportUsers)

		// 重置用户密码
		// PUT /api/v1/admin/users/{id}/reset-password
		// 响应码: 0-成功, 10002-参数错误, 20001-用户不存在
		// adminGroup.PUT("/users/:id/reset-password", rm.adminHandler.ResetUserPassword)

		// ========== 题目管理 ==========

		// 获取所有题目（管理员视图）
		// GET /api/v1/admin/problems?page=1&page_size=20&is_public=false
		// 响应码: 0-成功, 10002-参数错误
		// adminGroup.GET("/problems", rm.adminHandler.ListAllProblems)

		// 题目审核
		// PUT /api/v1/admin/problems/{id}/approve
		// PUT /api/v1/admin/problems/{id}/reject
		// 响应码: 0-成功, 10002-参数错误, 30001-题目不存在
		// adminGroup.PUT("/problems/:id/approve", rm.adminHandler.ApproveProblem)
		// adminGroup.PUT("/problems/:id/reject", rm.adminHandler.RejectProblem)

		// ========== 系统数据统计 ==========

		// 用户统计
		// GET /api/v1/admin/stats/users
		// 响应码: 0-成功
		// adminGroup.GET("/stats/users", rm.adminHandler.GetUserStats)

		// 题目统计
		// GET /api/v1/admin/stats/problems
		// 响应码: 0-成功
		// adminGroup.GET("/stats/problems", rm.adminHandler.GetProblemStats)

		// 提交统计
		// GET /api/v1/admin/stats/submissions
		// 响应码: 0-成功
		// adminGroup.GET("/stats/submissions", rm.adminHandler.GetSubmissionStats)

		// 判题系统统计
		// GET /api/v1/admin/stats/judge
		// 响应码: 0-成功
		// adminGroup.GET("/stats/judge", rm.adminHandler.GetJudgeStats)

		// ========== 系统日志 ==========

		// 获取系统日志
		// GET /api/v1/admin/logs?level=error&start_time=2024-01-01&end_time=2024-01-31
		// 响应码: 0-成功, 10002-参数错误
		// adminGroup.GET("/logs", rm.adminHandler.GetSystemLogs)

		// 获取操作审计日志
		// GET /api/v1/admin/audit-logs?user_id=xxx&action=create&resource=user
		// 响应码: 0-成功, 10002-参数错误
		// adminGroup.GET("/audit-logs", rm.adminHandler.GetAuditLogs)
	}
}
