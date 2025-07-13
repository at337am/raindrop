package middleware

import (
	"log/slog"

	"github.com/gin-gonic/gin"
)

// AccessLogger 是一个 Gin 中间件, 用于记录访问日志
func AccessLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 判断请求的路径
		if c.Request.URL.Path == "/api/info" {
			slog.Info("会话已建立", "clientIP", c.ClientIP())
		}

		c.Next()
	}
}
