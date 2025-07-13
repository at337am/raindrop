package cli

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Validate 校验所有配置参数
func (r *Runner) Validate() error {
	if err := r.validatePort(); err != nil {
		return err
	}

	// 校验并规范化所有用户提供的路径
	if err := r.validateAndNormalizePaths(); err != nil {
		return err
	}

	// 如果指定了共享目录, 则遍历并添加文件
	if err := r.walkSharedDir(); err != nil {
		return err
	}

	return nil
}

// validatePath 校验路径是否存在, 并根据要求检查其是否为目录
func validatePath(path string, mustBeDir bool) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("无法获取绝对路径: %w", err)
	}

	info, err := os.Stat(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("路径不存在: '%s'", absPath)
		}
		return "", fmt.Errorf("访问路径 '%s' 时出错: %w", absPath, err)
	}

	isDir := info.IsDir()
	if mustBeDir && !isDir {
		return "", fmt.Errorf("路径 '%s' 不是一个目录", absPath)
	}
	if !mustBeDir && isDir {
		return "", fmt.Errorf("路径 '%s' 是一个目录, 请提供文件路径", absPath)
	}

	return absPath, nil
}

// validatePort 校验端口号
func (r *Runner) validatePort() error {
	if r.Config.Port < 1 || r.Config.Port > 65535 {
		return fmt.Errorf("端口号 %d 无效, 必须在 1-65535 之间", r.Config.Port)
	}
	return nil
}

// validateAndNormalizePaths 校验所有由用户直接提供的路径, 并将它们转换为绝对路径。
func (r *Runner) validateAndNormalizePaths() error {
	// --- 校验共享文件列表 ---
	validatedPaths := make([]string, 0, len(r.Config.SharedFilePaths))
	for _, path := range r.Config.SharedFilePaths {
		absPath, err := validatePath(path, false)
		if err != nil {
			return err
		}
		validatedPaths = append(validatedPaths, absPath)
	}
	r.Config.SharedFilePaths = validatedPaths

	// --- 校验共享目录 ---
	if r.Config.SharedDirPath != "" {
		absPath, err := validatePath(r.Config.SharedDirPath, true)
		if err != nil {
			return err
		}
		r.Config.SharedDirPath = absPath
	}

	// --- 校验内容文件 ---
	if r.Config.ContentFilePath != "" {
		absPath, err := validatePath(r.Config.ContentFilePath, false)
		if err != nil {
			return err
		}
		r.Config.ContentFilePath = absPath
	}

	return nil
}

// walkSharedDir 遍历共享目录, 并将找到的文件路径添加到配置中。
func (r *Runner) walkSharedDir() error {
	if r.Config.SharedDirPath == "" {
		return nil
	}

	walkErr := filepath.WalkDir(r.Config.SharedDirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// 检查文件名或目录名是否以 "." 开头
		if strings.HasPrefix(d.Name(), ".") {
			if d.IsDir() {
				return fs.SkipDir // 如果是目录, 则跳过整个目录
			}
			return nil // 如果是文件, 则仅跳过此文件
		}

		if !d.IsDir() {
			r.Config.SharedFilePaths = append(r.Config.SharedFilePaths, path)
		}
		return nil
	})

	if walkErr != nil {
		return fmt.Errorf("遍历目录 '%s' 时失败: %w", r.Config.SharedDirPath, walkErr)
	}

	return nil
}
