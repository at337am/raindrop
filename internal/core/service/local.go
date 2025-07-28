package service

import (
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
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

	// addFile 是一个辅助函数, 用于向 content 中添加文件信息, 避免代码重复
	addFile := func(path string, info fs.FileInfo) {
		// 检查文件名是否已存在, 如果存在则记录警告并跳过, 防止映射冲突
		if _, exists := content.FileMap[info.Name()]; exists {
			slog.Warn("文件名冲突, 已跳过", "fileName", info.Name(), "path", path)
			return
		}
		fileInfo := handler.FileInfo{
			FileName: info.Name(),
			FileSize: humanize.IBytes(uint64(info.Size())),
		}
		content.Files = append(content.Files, fileInfo)
		content.FileMap[info.Name()] = path
	}

	// 1. 处理共享路径列表 (可以是文件或目录)
	if len(s.cfg.SharedPaths) > 0 {
		for _, path := range s.cfg.SharedPaths {
			info, err := os.Stat(path)
			if err != nil {
				slog.Error("无法获取共享路径信息", "path", path, "error", err)
				continue // 跳过无效路径
			}

			if info.IsDir() {
				// 使用 os.ReadDir 只遍历一级目录中的文件, 不再递归
				entries, err := os.ReadDir(path)
				if err != nil {
					slog.Error("无法读取目录", "path", path, "error", err)
					continue
				}

				for _, entry := range entries {
					// 跳过所有子目录
					if entry.IsDir() {
						continue
					}

					// 获取文件信息并添加
					fileInfo, statErr := entry.Info()
					if statErr != nil {
						slog.Error("无法获取文件信息", "path", filepath.Join(path, entry.Name()), "error", statErr)
						continue
					}
					addFile(filepath.Join(path, entry.Name()), fileInfo)
				}
			} else {
				// 如果是文件, 直接添加
				addFile(path, info)
			}
		}
	}

	// 2. 处理内容文件
	// 如果提供了 ContentFilePath, 则尝试读取
	if s.cfg.ContentPath != "" {
		raw, err := os.ReadFile(s.cfg.ContentPath)
		if err != nil {
			// 如果读取失败, 仅记录日志并继续
			slog.Error("无法读取内容文件", "path", s.cfg.ContentPath, "error", err)
		} else {
			content.Snippet = string(raw)
		}
	}

	return content
}
