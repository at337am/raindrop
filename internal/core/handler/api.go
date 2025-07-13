package handler

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SharedContent 存储返回 JSON 的数据
type SharedContent struct {
	FileName string
	FileSize string
	FilePath string
	Message  string
	Snippet  string
}

// Service 接口
type Service interface {
	GetContent() *SharedContent
}

type APIHandler struct {
	svc Service
}

func NewAPIHandler(s Service) *APIHandler {
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
		slog.Warn("所有共享内容字段均为空")
	}

	c.JSON(http.StatusOK, resp)
}

// HandleDownload 将领域对象适配为文件下载响应
func (h *APIHandler) HandleDownload(c *gin.Context) {
	sharedContent := h.svc.GetContent()

	if sharedContent.FilePath == "" {
		slog.Warn(
			"文件下载失败: 未找到可共享文件或文件路径无效",
			"clientIP", c.ClientIP(),
		)
		c.Status(http.StatusNotFound)
		return
	}

	// 记录下载开始的日志
	if c.GetHeader("Range") == "" {
		slog.Info(
			"开始下载文件",
			"fileName", sharedContent.FileName,
			"clientIP", c.ClientIP(),
		)
	}

	// 将指定文件作为附件发送给客户端, 浏览器将提示用户下载该文件
	c.FileAttachment(sharedContent.FilePath, sharedContent.FileName)
}
