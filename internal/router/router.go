package router

import (
	"zhku-oj/internal/handler/admin"
	"zhku-oj/internal/handler/auth"
	"zhku-oj/internal/handler/problem"
	"zhku-oj/internal/handler/submission"
	"zhku-oj/internal/handler/user"
	"zhku-oj/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RouterManager 路由管理器
type RouterManager struct {
	// Handler依赖
	authHandler       *auth.AuthHandler
	userHandler       *user.UserHandler
	problemHandler    *problem.ProblemHandler
	submissionHandler *submission.SubmissionHandler
	adminHandler      *admin.AdminHandler
}

// NewRouterManager 创建路由管理器
func NewRouterManager(
	authHandler *auth.AuthHandler,
	userHandler *user.UserHandler,
	problemHandler *problem.ProblemHandler,
	submissionHandler *submission.SubmissionHandler,
	adminHandler *admin.AdminHandler,
) *RouterManager {
	return &RouterManager{
		authHandler:       authHandler,
		userHandler:       userHandler,
		problemHandler:    problemHandler,
		submissionHandler: submissionHandler,
		adminHandler:      adminHandler,
	}
}

// SetupRoutes 设置所有路由
func (rm *RouterManager) SetupRoutes(router *gin.Engine) {
	// 设置全局中间件
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())
	router.Use(middleware.CORS())

	// 健康检查路由
	rm.setupHealthRoutes(router)

	// API路由组
	v1 := router.Group("/api/v1")
	{
		// 认证相关路由
		rm.setupAuthRoutes(v1)

		// 用户相关路由
		rm.setupUserRoutes(v1)

		// 题目相关路由
		rm.setupProblemRoutes(v1)

		// 提交相关路由
		rm.setupSubmissionRoutes(v1)

		// 管理员路由
		rm.setupAdminRoutes(v1)
	}
}
