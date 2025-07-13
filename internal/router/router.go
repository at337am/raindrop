package router

import (
	"net/http"
	"raindrop/assets"
	"raindrop/internal/middleware"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	HandleGetInfo(c *gin.Context)
	HandleDownload(c *gin.Context)
}

func SetupRouter(h Handler) *gin.Engine {

	gin.SetMode(gin.ReleaseMode) // 设置为 ReleaseMode 发布模式, 减少日志输出和性能优化

	// r := gin.Default() // 默认包含了 Logger 和 Recovery 中间件, 会输出请求日志

	r := gin.New() // 创建一个“空白”的 Gin 实例, 不包含任何默认中间件

	r.Use(gin.Recovery()) // 中间件: 用于捕获处理请求过程中发生的 panic, 并返回 500 错误, 防止程序崩溃

	// 使用中间件, 记录访问日志
	r.Use(middleware.AccessLogger())

	// 静态资源路由: 将 /static 映射到嵌入的模板文件系统
	r.StaticFS("/static", http.FS(assets.FS))

	// API 路由: 定义所有 API 接口, 并应用 NoCache 中间件
	apiGroup := r.Group("/api")
	apiGroup.Use(middleware.NoCache()) // 禁止浏览器缓存
	{
		apiGroup.GET("/info", h.HandleGetInfo)
		apiGroup.GET("/download", h.HandleDownload)
	}

	// 根路由: 提供应用程序的入口页面 (index.html)
	r.GET("/", func(c *gin.Context) {
		fileBytes, err := assets.FS.ReadFile("templates/index.html")
		if err != nil {
			c.String(http.StatusInternalServerError, "Internal Server Error: index.html not found")
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", fileBytes)
	})

	return r
}
