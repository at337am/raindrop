package handler

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FileInfo struct {
	FileName string
	FileSize string
}

// SharedContent 存储返回 JSON 的数据
type SharedContent struct {
	Files   []FileInfo        // 文件信息列表
	FileMap map[string]string // 用于按文件名快速查找文件路径的映射
	Message string
	Snippet string
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

	// 定义文件信息的响应结构
	type fileInfo struct {
		FileName string `json:"fileName"`
		FileSize string `json:"fileSize"`
	}

	// 定义响应结构
	type pageInfoResponse struct {
		Files       []fileInfo `json:"files"` // 文件列表
		Description string     `json:"description"`
		Snippet     string     `json:"snippet"`
		IsEmpty     bool       `json:"isEmpty"`
	}

	// 从领域对象映射到响应
	resp := pageInfoResponse{
		Files:       make([]fileInfo, 0, len(sharedContent.Files)), // 初始化切片
		Description: sharedContent.Message,
		Snippet:     sharedContent.Snippet,
	}

	// 填充文件列表
	for _, f := range sharedContent.Files {
		resp.Files = append(resp.Files, fileInfo{
			FileName: f.FileName,
			FileSize: f.FileSize,
		})
	}

	// 判断所有内容字段是否都为空
	if len(resp.Files) == 0 && resp.Description == "" && resp.Snippet == "" {
		resp.IsEmpty = true
		slog.Warn("所有共享内容字段均为空")
	}

	c.JSON(http.StatusOK, resp)
}

// HandleDownload 将领域对象适配为文件下载响应
func (h *APIHandler) HandleDownload(c *gin.Context) {
	// 从查询参数中获取请求的文件名
	fileName := c.Query("file")
	if fileName == "" {
		c.String(http.StatusBadRequest, "文件名参数 'file' 不能为空")
		return
	}

	sharedContent := h.svc.GetContent()

	filePath, ok := sharedContent.FileMap[fileName]

	if !ok || filePath == "" {
		slog.Warn(
			"文件下载失败: 未找到指定文件或路径无效",
			"requestedFile", fileName,
			"clientIP", c.ClientIP(),
		)
		c.Status(http.StatusNotFound)
		return
	}

	// 记录下载开始的日志
	// 记录下载开始的日志
	if c.GetHeader("Range") == "" {
		slog.Info(
			"开始下载文件",
			"fileName", fileName,
			"clientIP", c.ClientIP(),
		)
	}

	// 将指定文件作为附件发送给客户端
	c.FileAttachment(filePath, fileName)
}
