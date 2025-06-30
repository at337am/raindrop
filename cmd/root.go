package cmd

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"raindrop/assets"
	"raindrop/internal/app/handler"
	"raindrop/internal/app/service"
	"raindrop/internal/config"
	"raindrop/internal/router"
	"raindrop/pkg/fmtcolor"
	"strings"
)

// commandLineArgs 封装了从命令行解析出的所有参数
type commandLineArgs struct {
	sharedFile  string
	contentFile string
	message     string
	port        string
	nonFlagArgs []string
}

// Execute 程序的入口函数
func Execute() {
	// --help 帮助信息
	flag.Usage = func() {
		helpContent := getHelpContentByLocale()
		if helpContent == "" {
			flag.PrintDefaults()
			return
		}
		fmt.Fprint(flag.CommandLine.Output(), helpContent)
	}

	args := commandLineArgs{}

	flag.StringVar(&args.sharedFile, "i", "", "要共享的单个文件路径")
	flag.StringVar(&args.contentFile, "I", "", "要作为纯文本发送的文件路径")
	flag.StringVar(&args.message, "m", "", "要发送的消息内容")
	flag.StringVar(&args.port, "p", "1130", "指定服务器运行的端口")

	flag.Parse()

	// 处理非选项参数
	args.nonFlagArgs = flag.Args()

	// 校验选项参数
	if err := args.validateSyntax(); err != nil {
		fmtcolor.Error(err)
		fmtcolor.Error(errors.New("使用 -h --help 查看用法"))
		os.Exit(2)
	}

	// 使用 args 结构体中的值构建 AppConfig
	cfg := &config.AppConfig{
		SharedFilePath:  args.sharedFile,
		Message:         args.message,
		ContentFilePath: args.contentFile,
		Port:            args.port,
	}

	// 校验配置
	if err := cfg.Validate(); err != nil {
		fmtcolor.Error(err)
		fmtcolor.Error(fmt.Errorf("使用 -h --help 查看用法"))
		os.Exit(2)
	}

	// 显示启动信息
	printServerInfo(cfg.Port)

	// 调用 run() 启动服务
	if err := run(cfg); err != nil {
		fmtcolor.Error(err)
		os.Exit(1)
	}
}

// validateSyntax 专注于命令行语法校验
func (args *commandLineArgs) validateSyntax() error {
	if args.sharedFile == "" && args.contentFile == "" && args.message == "" && len(args.nonFlagArgs) == 0 {
		return fmt.Errorf("请提供要共享的文件 (-i), 要发送的文件 (-I) 或要发送的消息 (-m)")
	}

	if len(args.nonFlagArgs) > 1 {
		return fmt.Errorf("提供了多个非选项参数, 请使用 -i, -I 或 -m 等选项明确指定")
	}

	if len(args.nonFlagArgs) == 1 {
		// 如果用户已经通过 -i, -I, -m 指定了输入, 则产生冲突
		if args.sharedFile != "" || args.contentFile != "" || args.message != "" {
			return fmt.Errorf("不能同时使用选项 (-i, -I, -m) 和非选项参数")
		}
		// 否则, 将该独立参数视为 -i 的值
		args.sharedFile = args.nonFlagArgs[0]
	}

	return nil
}

// run 处理依赖注入, 并启动服务
func run(cfg *config.AppConfig) error {
	// 依赖注入和初始化
	svc := service.NewLocalFileService(cfg)
	h := handler.NewAPIHandler(svc)
	router := router.SetupRouter(h)

	// 启动服务
	addr := ":" + cfg.Port
	if err := router.Run(addr); err != nil {
		return fmt.Errorf("服务启动失败: %w", err)
	}
	return nil
}

// printServerInfo 打印服务器信息, 包括局域网 IP
func printServerInfo(port string) {
	// 2. 显示 logo
	logoContent := assets.TTYLogo
	if logoContent == "" {
		fmtcolor.Error(fmt.Errorf("无法获取嵌入的 logo 内容, 服务将继续启动, 但不会显示 logo"))
	} else {
		fmtcolor.Success(logoContent)
	}
	fmtcolor.Info("Starting raindrop server...")
	fmtcolor.Info("\nAccess URLs:")
	fmtcolor.Success(fmt.Sprintf("   Local:   http://127.0.0.1:%s", port))

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
				fmtcolor.Success(fmt.Sprintf("   Network: http://%s:%s", ip.String(), port))
			}
		}
	}

	fmt.Printf("\nPress Ctrl+C to stop the server.\n")
}

// getLocaleLanguage 获取当前系统的语言环境, 例如 "en" 或 "zh"
func getLocaleLanguage() string {
	lang := os.Getenv("LANG")
	if lang == "" {
		lang = os.Getenv("LC_ALL")
	}
	if lang != "" {
		// 提取主要语言部分, 例如 "zh_CN.UTF-8" -> "zh"
		if idx := strings.Index(lang, "_"); idx != -1 {
			return strings.ToLower(lang[:idx])
		}
		if idx := strings.Index(lang, "."); idx != -1 {
			return strings.ToLower(lang[:idx])
		}
		return strings.ToLower(lang)
	}
	return "en" // 默认语言
}

// getHelpContentByLocale 根据语言获取帮助内容
func getHelpContentByLocale() string {
	lang := getLocaleLanguage()
	switch lang {
	case "zh":
		return assets.HelpZH
	case "en":
		return assets.HelpEN
	default:
		// 如果没有对应的语言, 默认使用英文
		return assets.HelpEN
	}
}
