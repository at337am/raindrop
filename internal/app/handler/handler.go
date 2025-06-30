package handler

import (
	"log"
	"net/http"
	"raindrop/internal/app/service"

	"github.com/gin-gonic/gin"
)

// APIHandler 封装HTTP请求处理器, 通过svc接口调用业务逻辑
type APIHandler struct {
	svc service.Service
}

// NewAPIHandler 构造函数, 创建API处理器实例并注入服务依赖
func NewAPIHandler(s service.Service) *APIHandler {
	return &APIHandler{svc: s}
}

// HandleGetInfo 将领域对象适配为 /api/info 的 JSON 响应
func (h *APIHandler) HandleGetInfo(c *gin.Context) {
	sharedContent := h.svc.GetContent()

	// 定义响应结构
	type pageInfoResponse struct {
		FileName    string `json:"fileName"`
		FileSize    string `json:"fileSize"`
		Description string `json:"description"`
		Snippet     string `json:"snippet"`
		IsEmpty     bool   `json:"isEmpty"`
	}

	// 从领域对象映射到响应
	resp := pageInfoResponse{
		FileName:    sharedContent.FileName,
		FileSize:    sharedContent.FileSize,
		Description: sharedContent.Message,
		Snippet:     sharedContent.Snippet,
	}

	// 判断所有内容字段是否都为空
	if resp.FileName == "" && resp.FileSize == "" && resp.Description == "" && resp.Snippet == "" {
		resp.IsEmpty = true
		log.Printf("所有的共享内容字段都为空")
	}

	c.JSON(http.StatusOK, resp)
}

// HandleDownload 将领域对象适配为文件下载响应
func (h *APIHandler) HandleDownload(c *gin.Context) {
	sharedContent := h.svc.GetContent()

	if sharedContent.FilePath == "" {
		log.Printf("文件下载失败: 未找到可共享文件或文件路径无效, 来自 -> %s\n", c.ClientIP())
		c.Status(http.StatusNotFound)
		return
	}

	// 记录下载开始的日志
	if c.GetHeader("Range") == "" {
		log.Printf("开始下载文件: %s, 来自 -> %s\n", sharedContent.FileName, c.ClientIP())
	}

	// 将指定文件作为附件发送给客户端, 浏览器将提示用户下载该文件
	c.FileAttachment(sharedContent.FilePath, sharedContent.FileName)
}
