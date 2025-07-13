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

// Validate 校验参数
func (r *Runner) Validate() error {
	if r.Config.Port < 1 || r.Config.Port > 65535 {
		return fmt.Errorf("端口号 %d 无效, 必须在 1-65535 之间", r.Config.Port)
	}

	// isValidFilePath 检查路径是否有效, 是否为文件, 返回绝对路径
	isValidFilePath := func(path string) (string, error) {
		absPath, err := filepath.Abs(path)
		if err != nil {
			return "", fmt.Errorf("无法获取绝对路径: %w", err)
		}

		info, err := os.Stat(absPath)
		if err != nil {
			if os.IsNotExist(err) {
				return "", fmt.Errorf("文件不存在: '%s'", absPath)
			} else if os.IsPermission(err) {
				return "", fmt.Errorf("权限不足, 无法访问: %q", absPath)
			} else {
				return "", fmt.Errorf("访问路径 '%s' 时发生未知错误: %w", absPath, err)
			}
		}
		if info.IsDir() {
			return "", fmt.Errorf("路径 '%s' 是一个目录, 请提供文件路径", absPath)
		}

		return absPath, nil
	}

	// --- 共享文件路径校验 ---
	if len(r.Config.SharedFilePaths) > 0 {
		validatedPaths := make([]string, 0, len(r.Config.SharedFilePaths))
		for _, path := range r.Config.SharedFilePaths {
			absPath, err := isValidFilePath(path)
			if err != nil {
				return err
			}
			validatedPaths = append(validatedPaths, absPath)
		}
		r.Config.SharedFilePaths = validatedPaths
	}

	// --- 内容文件路径校验 ---
	if r.Config.ContentFilePath != "" {
		absPath, err := isValidFilePath(r.Config.ContentFilePath)
		if err != nil {
			return err
		}
		r.Config.ContentFilePath = absPath
	}

	return nil
}

// Run 执行核心逻辑
func (r *Runner) Run() error {
	// 依赖注入和初始化
	s := service.NewLocalService(r.Config)
	h := handler.NewAPIHandler(s)
	router := router.SetupRouter(h)

	// 拼接端口号
	addr := ":" + strconv.Itoa(r.Config.Port)

	// 显示启动信息
	printServerInfo(addr)

	// 启动服务
	if err := router.Run(addr); err != nil {
		return fmt.Errorf("服务启动失败: %w", err)
	}

	return nil
}

// printServerInfo 打印服务器信息, 包括局域网 IP
func printServerInfo(port string) {
	// 显示 logo
	successColor.Print(assets.Logo)

	fmt.Println("Starting raindrop server...")
	fmt.Println("\nAccess URLs:")
	successColor.Printf("   Local:   http://127.0.0.1%s\n", port)

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
				successColor.Printf("   Network: http://%s:%s\n", ip.String(), port)
			}
		}
	}

	fmt.Printf("\nPress Ctrl+C to stop the server.\n")
}
