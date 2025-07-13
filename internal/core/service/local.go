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
		Files:   []handler.FileInfo{},
		FileMap: make(map[string]string),
	}

	// 1. 处理共享文件列表
	// 如果提供了 SharedFilePaths, 则获取每个文件的信息
	if len(s.cfg.SharedFilePaths) > 0 {
		for _, path := range s.cfg.SharedFilePaths {
			info, err := os.Stat(path)
			if err != nil {
				slog.Error("无法获取共享文件信息", "path", path, "error", err)
				continue // 跳过无效文件
			}
			if info.IsDir() {
				slog.Warn("路径是一个目录, 已跳过", "path", path)
				continue // 跳过目录
			}

			// 填充文件信息
			fileInfo := handler.FileInfo{
				FileName: info.Name(),
				FileSize: humanize.IBytes(uint64(info.Size())),
			}
			content.Files = append(content.Files, fileInfo)
			content.FileMap[info.Name()] = path
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
