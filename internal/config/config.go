package config

// Config 聚合了应用的所有配置项
type Config struct {
	SharedFilePaths []string // 多个共享文件路径
	SharedDirPath   string   // 共享目录的路径
	Message         string   // 消息内容
	ContentFilePath string   // 内容文件的路径
	Port            int
}
