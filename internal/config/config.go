package config

// Config 聚合了应用的所有配置项
type Config struct {
	SharedFilePaths []string // 支持多个共享文件路径
	Message         string
	ContentFilePath string
	Port            int
}
