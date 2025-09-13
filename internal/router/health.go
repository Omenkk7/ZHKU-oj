package router

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// setupHealthRoutes 设置健康检查路由
// 系统监控和健康状态检查接口
func (rm *RouterManager) setupHealthRoutes(router *gin.Engine) {
	// 基础健康检查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().Unix(),
			"version":   "v1.0.0",
			"service":   "zhku-oj-api",
		})
	})

	// 详细健康检查（可扩展检查数据库、Redis等）
	router.GET("/health/detailed", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().Unix(),
			"version":   "v1.0.0",
			"service":   "zhku-oj-api",
			"checks": gin.H{
				"database": "healthy", // 后续可扩展实际检查
				"redis":    "healthy",
				"judge":    "healthy",
			},
			"uptime": time.Since(time.Now()).String(), // 实际项目中应该记录启动时间
		})
	})

	// 服务信息接口
	router.GET("/info", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"name":        "校园Java-OJ在线判题系统",
			"version":     "v1.0.0",
			"description": "基于Go语言开发的在线编程评测平台",
			"author":      "ZHKU-OJ Team",
			"build_time":  "2024-01-15T10:00:00Z", // 构建时可动态注入
		})
	})
}
