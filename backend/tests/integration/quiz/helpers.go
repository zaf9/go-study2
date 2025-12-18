package quiz

import (
	"os"
	"path/filepath"
)

func findQuizDataPath() string {
	// 从当前目录开始向上寻找 go.mod
	dir, err := os.Getwd()
	if err != nil {
		return "quiz_data" // Fallback
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			// 找到根目录，构造路径
			p := filepath.Join(dir, "quiz_data")
			if _, err := os.Stat(p); err == nil {
				return p
			}
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			// 到达文件系统根目录
			break
		}
		dir = parent
	}

	// 针对异常情况的回退
	return "quiz_data"
}
