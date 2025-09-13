package router

import (
	"zhku-oj/internal/middleware"

	"github.com/gin-gonic/gin"
)

// setupAuthRoutes 设置认证相关路由
// 用户注册、登录、登出等认证功能
func (rm *RouterManager) setupAuthRoutes(v1 *gin.RouterGroup) {
	authGroup := v1.Group("/auth")
	{
		// 用户注册
		// POST /api/v1/auth/register
		// 响应码: 0-成功, 10002-参数错误, 20002-用户已存在
		authGroup.POST("/register", rm.authHandler.Register)

		// 用户登录
		// POST /api/v1/auth/login
		// 响应码: 0-成功, 10002-参数错误, 20010-登录失败
		authGroup.POST("/login", rm.authHandler.Login)

		// 用户登出（需要认证）
		// POST /api/v1/auth/logout
		// 响应码: 0-成功, 10003-未授权
		authGroup.POST("/logout", middleware.AuthRequired(), rm.authHandler.Logout)

		// 刷新Token（需要认证）
		// POST /api/v1/auth/refresh
		// 响应码: 0-成功, 10003-未授权, 10009-Token无效
		// authGroup.POST("/refresh", middleware.AuthRequired(), rm.authHandler.RefreshToken)

		// 修改密码（需要认证）
		// PUT /api/v1/auth/password
		// 响应码: 0-成功, 10002-参数错误, 20013-旧密码不正确
		authGroup.PUT("/password", middleware.AuthRequired(), rm.userHandler.ChangePassword)

		// 验证Token状态
		// GET /api/v1/auth/verify
		// 响应码: 0-成功, 10003-未授权, 10009-Token无效
		authGroup.GET("/verify", middleware.AuthRequired(), func(c *gin.Context) {
			userID := middleware.GetUserID(c)
			username := middleware.GetUsername(c)
			role := middleware.GetUserRole(c)

			c.JSON(200, gin.H{
				"code":    0,
				"message": "Token有效",
				"data": gin.H{
					"user_id":  userID,
					"username": username,
					"role":     role,
				},
			})
		})
	}
}
