package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Restore() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("获取用户目录失败: %v", err)
	}

	extensionsDir := filepath.Join(home, ".vscode", "extensions")
	pattern := filepath.Join(extensionsDir, "eamodio.gitlens-*", "dist", "*.backup")

	matches, _ := filepath.Glob(pattern)
	for _, backup := range matches {
		original := strings.TrimSuffix(backup, ".backup")
		if err := os.Rename(backup, original); err != nil {
			return fmt.Errorf("还原文件失败: %v", err)
		}
	}

	return nil
}
