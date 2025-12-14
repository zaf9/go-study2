package quiz_test

import (
	"os"
	"path/filepath"
)

func findQuizDataPath() string {
	candidates := []string{
		filepath.Join("backend", "quiz_data"),
		filepath.Join("backend", "backend", "quiz_data"),
		filepath.Join("..", "backend", "quiz_data"),
		filepath.Join("..", "..", "backend", "quiz_data"),
		filepath.Join("..", "..", "..", "backend", "quiz_data"),
	}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c
		}
	}
	return filepath.Join("backend", "backend", "quiz_data")
}
