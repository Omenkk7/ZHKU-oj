package middleware

import (
	"time"
	"zhku-oj/pkg/logger"

	"github.com/gin-gonic/gin"
)

// Logger 请求日志中间件
func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		logger.Info("请求日志",
			"status", param.StatusCode,
			"method", param.Method,
			"path", param.Path,
			"ip", param.ClientIP,
			"user_agent", param.Request.UserAgent(),
			"latency", param.Latency,
			"time", param.TimeStamp.Format(time.RFC3339),
		)
		return ""
	})
}

// Recovery 异常恢复中间件
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		logger.Error("请求panic恢复",
			"error", recovered,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"ip", c.ClientIP(),
		)

		c.JSON(500, gin.H{
			"code":    500,
			"message": "服务器内部错误",
		})
	})
}
