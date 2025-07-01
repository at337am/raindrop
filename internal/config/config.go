package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

// AppConfig 存储应用的配置信息, 经过校验后的可靠数据
type AppConfig struct {
	SharedFilePath  string
	Message         string
	ContentFilePath string
	Port            string
}

// Validate 专注于应用配置的有效性
func (cfg *AppConfig) Validate() error {
	// --- 端口校验 ---
	portNum, err := strconv.Atoi(cfg.Port)
	if err != nil {
		return fmt.Errorf("端口号格式不正确, 必须为数字: %w", err)
	}
	if portNum < 1 || portNum > 65535 {
		return fmt.Errorf("端口号 %d 无效, 必须在 1-65535 之间", portNum)
	}

	// --- 共享文件路径校验 ---
	if cfg.SharedFilePath != "" {
		cfg.SharedFilePath, err = isValidFilePath(cfg.SharedFilePath)
		if err != nil {
			return err
		}
	}

	// --- 内容文件路径校验 ---
	if cfg.ContentFilePath != "" {
		cfg.ContentFilePath, err = isValidFilePath(cfg.ContentFilePath)
		if err != nil {
			return err
		}
	}

	return nil
}

// isValidFilePath 检查路径是否有效, 是否为文件, 返回绝对路径
func isValidFilePath(path string) (string, error) {
	if path == "" {
		// 虽然上层已经判断了不为空, 但作为健壮的私有函数, 加上这个检查更好
		return "", fmt.Errorf("路径不能为空")
	}

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
