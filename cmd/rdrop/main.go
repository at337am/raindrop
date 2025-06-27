package main

import (
	"fmt"
	"net"
	"os"
	"raindrop/assets"
	"raindrop/internal/app/handler"
	"raindrop/internal/app/service"
	"raindrop/internal/config"
	"raindrop/internal/router"
	"raindrop/pkg/fmtutil"
)

func main() {
	// 1. 加载并校验配置
	appCfg, err := config.ValidateAndLoadConfig()
	if err != nil {
		fmtutil.PrintError(fmt.Sprint(err))
		fmt.Println("使用 -h --help 查看用法")
		os.Exit(1)
	}

	// 2. 显示 logo
	logoContent := assets.TTYLogo
	if logoContent == "" {
		fmtutil.PrintError("无法获取嵌入的 logo 内容")
	} else {
		// 输出带有高亮颜色的 logo
		fmt.Fprintln(os.Stdout, "\x1b[32m"+logoContent+"\x1b[0m")
	}

	// 3. 依赖注入和初始化
	apiService := service.NewAPIService(appCfg)
	apiHandler := handler.NewAPIHandler(apiService)
	router := router.SetupRouter(apiHandler)

	// 4. 启动服务并打印访问地址
	addr := ":" + appCfg.Port
	printServerInfo(appCfg.Port)
	if err := router.Run(addr); err != nil {
		fmtutil.PrintError(fmt.Sprintf("启动服务器失败: %v", err))
		os.Exit(1)
	}
}

// printServerInfo 打印服务器信息，包括局域网 IP
func printServerInfo(port string) {
	fmt.Printf("Starting raindrop server...\n")
	fmt.Printf("Access URLs:\n")
	fmt.Printf("   Local:   http://127.0.0.1:%s\n", port)

	interfaces, err := net.Interfaces()
	if err != nil {
		// 如果获取网络接口失败，不打印错误，仅跳过局域网地址的显示
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
				fmt.Printf("   Network: http://%s:%s\n", ip.String(), port)
			}
		}
	}

	fmt.Printf("Press Ctrl+C to stop the server.\n")
}
