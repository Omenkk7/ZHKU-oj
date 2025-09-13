package router

import (
	"zhku-oj/internal/middleware"

	"github.com/gin-gonic/gin"
)

// setupUserRoutes 设置用户相关路由
// 用户信息管理、个人资料等功能
func (rm *RouterManager) setupUserRoutes(v1 *gin.RouterGroup) {
	userGroup := v1.Group("/users")
	userGroup.Use(middleware.AuthRequired()) // 所有用户接口都需要认证
	{
		// ========== 当前用户相关接口 ==========

		// 获取当前用户信息
		// GET /api/v1/users/profile
		// 响应码: 0-成功, 10002-参数错误, 20001-用户不存在
		userGroup.GET("/profile", rm.userHandler.GetProfile)

		// 更新当前用户信息
		// PUT /api/v1/users/profile
		// 响应码: 0-成功, 10002-参数错误, 20001-用户不存在
		userGroup.PUT("/profile", rm.userHandler.UpdateProfile)

		// ========== 用户查询接口 ==========

		// 获取指定用户信息
		// GET /api/v1/users/{id}
		// 响应码: 0-成功, 10002-参数错误, 20001-用户不存在
		userGroup.GET("/:id", rm.userHandler.GetUser)

		// 获取用户统计信息
		// GET /api/v1/users/{id}/stats
		// 响应码: 0-成功, 10002-参数错误, 20001-用户不存在
		userGroup.GET("/:id/stats", rm.userHandler.GetUserStats)

		// 获取用户提交历史（分页）
		// GET /api/v1/users/{id}/submissions?page=1&page_size=20
		// 响应码: 0-成功, 10002-参数错误, 20001-用户不存在
		// userGroup.GET("/:id/submissions", rm.submissionHandler.GetUserSubmissions)

		// 获取用户解题记录
		// GET /api/v1/users/{id}/solved
		// 响应码: 0-成功, 10002-参数错误, 20001-用户不存在
		// userGroup.GET("/:id/solved", rm.userHandler.GetUserSolvedProblems)

		// ========== 用户排行榜 ==========

		// 获取用户排行榜
		// GET /api/v1/users/ranking?page=1&page_size=50&class=软件工程1班
		// 响应码: 0-成功, 10002-参数错误
		// userGroup.GET("/ranking", rm.userHandler.GetUserRanking)
	}
}
