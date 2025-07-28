package config

// Config 聚合了应用的所有配置项
type Config struct {
	SharedPaths []string // 多个共享路径, 包含文件或目录路径
	Message     string   // 消息内容
	ContentPath string   // 内容文件的路径
	Port        int
}
