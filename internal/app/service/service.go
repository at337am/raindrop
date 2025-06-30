package service

import (
	"log"
	"os"
	"raindrop/internal/config"

	"github.com/dustin/go-humanize"
)

// Service 定义了应用的核心业务逻辑接口, 用于解耦和测试
type Service interface {
	GetContent() *SharedContent
}

// SharedContent 是一个领域模型, 封装了所有要共享的内容信息
type SharedContent struct {
	FileName string
	FileSize string
	FilePath string
	Message  string
	Snippet  string
}

// LocalFileService 是 Service 接口的具体实现, 负责从本地文件系统读取内容
type LocalFileService struct {
	cfg *config.AppConfig
}

// NewLocalFileService 创建并返回一个基于本地文件的新服务实例
func NewLocalFileService(c *config.AppConfig) Service {
	return &LocalFileService{cfg: c}
}

// GetContent 从配置中指定的路径检索文件信息和内容, 填充到 SharedContent 对象中
func (s *LocalFileService) GetContent() *SharedContent {
	// 创建一个基础的对象, 先填充不受文件影响的数据
	content := &SharedContent{
		Message: s.cfg.Message,
	}

	// 1. 处理共享文件
	// 如果提供了 SharedFilePath, 则获取其信息
	if s.cfg.SharedFilePath != "" {
		info, err := os.Stat(s.cfg.SharedFilePath)
		// 只有当没有错误且不是目录时才填充文件信息
		if err != nil {
			log.Printf("获取共享文件 '%s' 信息时出错: %v", s.cfg.SharedFilePath, err)
		} else if info.IsDir() {
			log.Printf("路径 '%s' 是一个目录, 无法作为共享文件", s.cfg.SharedFilePath)
		} else {
			// 既没有错误也不是目录, 才会填充
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
			log.Printf("读取内容文件 '%s' 时出错: %v", s.cfg.ContentFilePath, err)
		} else {
			content.Snippet = string(raw)
		}
	}

	return content
}
