package cli

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"raindrop/assets"
	"raindrop/internal/config"
	"raindrop/internal/core/handler"
	"raindrop/internal/core/service"
	"raindrop/internal/router"
	"strconv"

	"github.com/fatih/color"
)

var (
	successColor = color.New(color.FgGreen)
	warnColor    = color.New(color.FgCyan)
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

// printServerInfo 打印服务器信息, 包括局域网 IP 和共享内容摘要
func (r *Runner) printServerInfo() {
	// orDefault 是一个辅助函数, 如果输入值为空字符串, 则返回指定的默认值。
	orDefault := func(value, defaultValue string) string {
		if value == "" {
			return defaultValue
		}
		return value
	}

	// 获取当前工作目录, 用于将绝对路径转换为相对路径进行显示
	wd, _ := os.Getwd()

	// toRelative 是一个辅助函数, 尝试将绝对路径转换为相对于当前工作目录的路径。
	toRelative := func(absPath string) string {
		// 如果路径为空或无法获取工作目录, 则返回原始路径
		if absPath == "" || wd == "" {
			return absPath
		}
		// 尝试进行路径转换, 如果失败则返回原始绝对路径
		if relPath, err := filepath.Rel(wd, absPath); err == nil {
			return relPath
		}
		return absPath
	}

	// 显示 logo
	successColor.Print(assets.Logo)

	// 打印共享内容信息
	warnColor.Printf("Sharing List:\n")
	fmt.Printf("  Message:          %s\n", orDefault(r.Config.Message, "<not set>"))
	fmt.Printf("  Content File:     %s\n", orDefault(toRelative(r.Config.ContentPath), "<not set>"))
	fmt.Printf("  Shared Paths (%d):\n", len(r.Config.SharedPaths))
	if len(r.Config.SharedPaths) > 0 {
		for i, path := range r.Config.SharedPaths {
			fmt.Printf("    %d. %s\n", i+1, toRelative(path))
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

	for _, iface := range interfaces {
		// 1. 首先过滤掉无效接口: 非活动状态, 环回接口
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP

			// 2. 获取纯粹的 net.IP 地址
			if ipnet, ok := addr.(*net.IPNet); ok {
				ip = ipnet.IP
			} else {
				// 如果需要处理其他类型, 可以在这里添加, 但*net.IPNet最常用
				continue
			}

			// 3. 确保是有效的全局单播 IPv4 地址
			if ip != nil && ip.IsGlobalUnicast() && ip.To4() != nil {
				fmt.Printf("  Network: ")
				successColor.Printf("http://%s:%s\n", ip.String(), port)
			}
		}
	}

	fmt.Printf("\nPress Ctrl+C to stop the server\n")
}
