package config

// Config 聚合了应用的所有配置项
type Config struct {
	SharedFilePath  string
	Message         string
	ContentFilePath string
	Port            int
}
