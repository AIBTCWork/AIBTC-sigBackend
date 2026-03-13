package utils

import (
	"os"
	"path/filepath"
)

func GetRootDir() string {
	// 以包含 go.mod 的目录作为项目根目录，避免因执行目录不同导致路径错误。
	// 找不到时兜底返回当前工作目录。
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	for {
		if _, statErr := os.Stat(filepath.Join(dir, "go.mod")); statErr == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return dir
		}
		dir = parent
	}
}
