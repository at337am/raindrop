package cli

import (
	"fmt"
	"net"
	"raindrop/assets"
	"raindrop/internal/config"
	"raindrop/internal/core/handler"
	"raindrop/internal/core/service"
	"raindrop/internal/router"
	"strconv"
)

// Runner 存储选项参数
type Runner struct {
	Config *config.Config
}

// NewRunner 构造函数, 创建API处理器实例并注入服务依赖
func NewRunner() *Runner {
	return &Runner{
		// 初始化 Config 字段, 避免在绑定命令行参数时发生空指针错误
		Config: &config.Config{},
	}
}

// Run 执行核心逻辑
func (r *Runner) Run() error {
	// 显示启动信息
	r.printServerInfo()

	// 依赖注入和初始化
	s := service.NewLocalService(r.Config)
	h := handler.NewAPIHandler(s)
	router := router.SetupRouter(h)

	// 拼接端口号
	addr := ":" + strconv.Itoa(r.Config.Port)

	// 启动服务
	if err := router.Run(addr); err != nil {
		return fmt.Errorf("服务启动失败: %w", err)
	}

	return nil
}

// orDefault 是一个辅助函数, 如果输入值为空字符串, 则返回指定的默认值。
func orDefault(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

// printServerInfo 打印服务器信息, 包括局域网 IP 和共享内容摘要
func (r *Runner) printServerInfo() {
	// 显示 logo
	successColor.Print(assets.Logo)

	// 打印共享内容信息
	warnColor.Printf("Sharing List:\n")
	fmt.Printf("  Message:          %s\n", orDefault(r.Config.Message, "<not set>"))
	fmt.Printf("  Content File:     %s\n", orDefault(r.Config.ContentFilePath, "<not set>"))
	fmt.Printf("  Shared Directory: %s\n", orDefault(r.Config.SharedDirPath, "<not set>"))
	fmt.Printf("  Shared Files (%d):\n", len(r.Config.SharedFilePaths))
	if len(r.Config.SharedFilePaths) > 0 {
		for i, path := range r.Config.SharedFilePaths {
			fmt.Printf("    %d. %s\n", i+1, path)
		}
	}

	port := strconv.Itoa(r.Config.Port)

	warnColor.Printf("\nAccess URLs:\n")
	fmt.Printf("  Local:   ")
	successColor.Printf("http://127.0.0.1:%s\n", port)

	interfaces, err := net.Interfaces()
	if err != nil {
		// 如果获取网络接口失败, 不打印错误, 仅跳过局域网地址的显示
		return
	}
	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			// 过滤掉回环地址和非 IPv4 地址
			if ip != nil && !ip.IsLoopback() && ip.To4() != nil {
				fmt.Printf("  Network: ")
				successColor.Printf("http://%s:%s\n", ip.String(), port)
			}
		}
	}

	fmt.Printf("\nPress Ctrl+C to stop the server\n")
}
