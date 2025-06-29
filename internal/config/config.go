package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"raindrop/assets"
	"raindrop/pkg/pathutil"
	"strconv"
	"strings"
)

// AppConfig 存储应用的配置信息, 经过校验后的可靠数据
type AppConfig struct {
	SharedFileAbsPath  string
	Message            string
	ContentFileAbsPath string
	Port               string
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

// ValidateAndLoadConfig 解析命令行参数并进行校验,
func ValidateAndLoadConfig() (*AppConfig, error) {
	// --help 帮助信息
	flag.Usage = func() {
		helpContent := getHelpContentByLocale()
		if helpContent == "" {
			flag.PrintDefaults()
			return
		}
		fmt.Fprint(flag.CommandLine.Output(), helpContent)
	}

	var (
		sharedFile  string
		contentFile string
		message     string
		port        string
	)
	flag.StringVar(&sharedFile, "i", "", "要共享的单个文件路径")
	flag.StringVar(&contentFile, "I", "", "要作为纯文本发送的文件路径")
	flag.StringVar(&message, "m", "", "要发送的消息内容")
	flag.StringVar(&port, "p", "1130", "指定服务器运行的端口")
	flag.Parse()

	// 处理非选项参数
	nonFlagArgs := flag.Args()

	if sharedFile == "" && contentFile == "" && message == "" && len(nonFlagArgs) == 0 {
		return nil, fmt.Errorf("请提供要共享的文件 (-i), 要发送的文件 (-I) 或要发送的消息 (-m)")
	}

	if len(nonFlagArgs) > 1 {
		return nil, fmt.Errorf("提供了多个非选项参数, 请使用 -i, -I 或 -m 等选项明确指定")
	}

	if len(nonFlagArgs) == 1 {
		// 如果用户已经通过 -i, -I, -m 指定了输入, 则产生冲突
		if sharedFile != "" || contentFile != "" || message != "" {
			return nil, fmt.Errorf("不能同时使用选项 (-i, -I, -m) 和非选项参数")
		}
		// 否则, 将该独立参数视为 -i 的值
		sharedFile = nonFlagArgs[0]
	}

	// 校验 port, 无论来源是默认值还是用户输入, 都必须是有效的
	portNum, err := strconv.Atoi(port)
	if err != nil {
		return nil, fmt.Errorf("端口号格式不正确, 必须为数字: %w", err)
	}
	if portNum < 1 || portNum > 65535 {
		return nil, fmt.Errorf("端口号 %d 无效, 必须在 1-65535 之间", portNum)
	}

	// 校验路径正确, 且是否文件
	var sharedFileAbsPath string
	if sharedFile != "" {
		var err error
		sharedFileAbsPath, err = filepath.Abs(sharedFile)
		if err != nil {
			return nil, err
		}
		if err := pathutil.IsValidFilePath(sharedFileAbsPath); err != nil {
			return nil, fmt.Errorf("校验 -i 选项参数时出错: %w", err)
		}
	}

	var contentFileAbsPath string
	if contentFile != "" {
		var err error
		contentFileAbsPath, err = filepath.Abs(contentFile)
		if err != nil {
			return nil, err
		}
		if err := pathutil.IsValidFilePath(contentFileAbsPath); err != nil {
			return nil, fmt.Errorf("校验 -I 选项参数时出错: %w", err)
		}
	}

	// 所有校验通过, 构建并返回配置对象
	return &AppConfig{
		SharedFileAbsPath:  sharedFileAbsPath,
		Message:            message,
		ContentFileAbsPath: contentFileAbsPath,
		Port:               port,
	}, nil
}
