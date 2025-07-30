package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// Validate 校验所有配置参数
func (r *Runner) Validate() error {
	// --- 校验端口号 ---
	if r.Config.Port < 1 || r.Config.Port > 65535 {
		return fmt.Errorf("端口号 %d 无效, 必须在 1-65535 之间", r.Config.Port)
	}

	// --- 校验共享路径 (文件或目录) ---
	validatedPaths := make([]string, 0, len(r.Config.SharedPaths))
	for _, path := range r.Config.SharedPaths {
		absPath, err := filepath.Abs(path)
		if err != nil {
			return fmt.Errorf("无法获取绝对路径: %w", err)
		}

		// 只校验路径是否存在, 不区分是文件还是目录
		if _, err := os.Stat(absPath); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("路径不存在: '%s'", absPath)
			}
			return fmt.Errorf("访问路径 '%s' 时出错: %w", absPath, err)
		}
		validatedPaths = append(validatedPaths, absPath)
	}
	r.Config.SharedPaths = validatedPaths

	// --- 校验内容文件 ---
	if r.Config.ContentPath != "" {
		absPath, err := filepath.Abs(r.Config.ContentPath)
		if err != nil {
			return fmt.Errorf("无法获取绝对路径: %w", err)
		}

		info, err := os.Stat(absPath)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("内容文件路径不存在: '%s'", absPath)
			}
			return fmt.Errorf("访问内容文件路径 '%s' 时出错: %w", absPath, err)
		}

		// 确保内容文件不是一个目录
		if info.IsDir() {
			return fmt.Errorf("路径 '%s' 是一个目录, 请为内容文件提供一个文件路径", absPath)
		}
		r.Config.ContentPath = absPath
	}

	return nil
}
