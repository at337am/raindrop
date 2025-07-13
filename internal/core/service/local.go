package service

import (
	"log/slog"
	"os"
	"raindrop/internal/config"
	"raindrop/internal/core/handler"

	"github.com/dustin/go-humanize"
)

// LocalService 是 Service 接口的具体实现, 负责从本地文件系统读取内容
type LocalService struct {
	cfg *config.Config
}

// NewLocalService 创建并返回一个基于本地文件的新服务实例
func NewLocalService(c *config.Config) *LocalService {
	return &LocalService{
		cfg: c,
	}
}

// GetContent 从配置中指定的路径检索文件信息和内容, 填充到 SharedContent 对象中
func (s *LocalService) GetContent() *handler.SharedContent {
	// 创建一个基础的对象, 先填充不受文件影响的数据
	content := &handler.SharedContent{
		Message: s.cfg.Message,
	}

	// 1. 处理共享文件
	// 如果提供了 SharedFilePath, 则获取其信息
	if s.cfg.SharedFilePath != "" {
		info, err := os.Stat(s.cfg.SharedFilePath)
		if err != nil {
			slog.Error("无法获取共享文件信息", "path", s.cfg.SharedFilePath, "error", err)
		} else if info.IsDir() {
			slog.Warn("路径是一个目录，不能作为文件共享", "path", s.cfg.SharedFilePath)
		} else {
			// 只有当没有错误且不是目录时才填充文件信息
			content.FileName = info.Name()
			content.FileSize = humanize.IBytes(uint64(info.Size()))
			content.FilePath = s.cfg.SharedFilePath
		}
	}

	// 2. 处理内容文件
	// 如果提供了 ContentFilePath, 则尝试读取
	if s.cfg.ContentFilePath != "" {
		raw, err := os.ReadFile(s.cfg.ContentFilePath)
		if err != nil {
			// 如果读取失败, 仅记录日志并继续
			slog.Error("无法读取内容文件", "path", s.cfg.ContentFilePath, "error", err)
		} else {
			content.Snippet = string(raw)
		}
	}

	return content
}
