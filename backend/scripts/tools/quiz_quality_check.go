package main

// 简易题库质量检查工具：扫描 backend/quiz_data 下的 YAML 文件，调用 domain/quiz.ValidateBank

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	quizdom "go-study2/internal/domain/quiz"

	"gopkg.in/yaml.v3"
)

func main() {
	base := filepath.Join("backend", "quiz_data")
	var entries []string
	_ = filepath.WalkDir(base, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ".yaml" || filepath.Ext(path) == ".yml" {
			entries = append(entries, path)
		}
		return nil
	})
	if len(entries) == 0 {
		fmt.Println("no quiz yaml files found under", base)
		os.Exit(1)
	}
	failed := 0
	for _, f := range entries {
		data, err := os.ReadFile(f)
		if err != nil {
			fmt.Printf("read %s failed: %v\n", f, err)
			failed++
			continue
		}
		var bank quizdom.YAMLBank
		if err := yaml.Unmarshal(data, &bank); err != nil {
			fmt.Printf("yaml unmarshal %s failed: %v\n", f, err)
			failed++
			continue
		}
		// 使用带文件上下文的验证以获得更好错误信息
		if err := quizdom.ValidateBankWithSource(bank, f, nil); err != nil {
			fmt.Printf("validation failed for %s: %v\n", f, err)
			failed++
			continue
		}
		fmt.Printf("OK: %s (%d questions)\n", f, len(bank.Questions))
	}
	if failed > 0 {
		fmt.Printf("%d files failed validation\n", failed)
		os.Exit(2)
	}
	fmt.Println("All quiz files passed validation")
}
