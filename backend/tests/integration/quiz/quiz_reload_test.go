package quiz_test

import (
	"os"
	"path/filepath"
	"testing"

	quiz "go-study2/internal/domain/quiz"
)

func TestLoadAllBanks_ReloadsUpdatedFile(t *testing.T) {
	root := t.TempDir()
	topicDir := filepath.Join(root, "testtopic")
	if err := os.MkdirAll(topicDir, 0o755); err != nil {
		t.Fatalf("创建目录失败: %v", err)
	}
	filePath := filepath.Join(topicDir, "chapter1.yaml")
	yaml := `questions:
  - id: q1
    type: single
    difficulty: easy
    stem: "样例题 1"
    options: ["A","B"]
    answer: A
    explanation: "示例"
    topic: testtopic
    chapter: chapter1
  - id: q2
    type: multiple
    difficulty: medium
    stem: "样例题 2"
    options: ["A","B","C"]
    answer: AB
    explanation: "示例"
    topic: testtopic
    chapter: chapter1
`
	if err := os.WriteFile(filePath, []byte(yaml), 0o644); err != nil {
		t.Fatalf("写入 yaml 失败: %v", err)
	}

	repo := quiz.NewRepository()
	if err := quiz.LoadAllBanks(root, repo); err != nil {
		t.Fatalf("LoadAllBanks 失败: %v", err)
	}
	qs, ok := repo.GetBank("testtopic", "chapter1")
	if !ok {
		t.Fatalf("未找到加载的题库")
	}
	if len(qs) != 2 {
		t.Fatalf("期望 2 道题, 但得到 %d", len(qs))
	}

	// 修改文件并再次加载
	yaml2 := `questions:
  - id: q1
    type: single
    difficulty: easy
    stem: "样例题 1 修改"
    options: ["A","B"]
    answer: A
    explanation: "示例"
    topic: testtopic
    chapter: chapter1
`
	if err := os.WriteFile(filePath, []byte(yaml2), 0o644); err != nil {
		t.Fatalf("写入 yaml2 失败: %v", err)
	}
	// 重新创建 repo 并加载，或清空后加载
	repo2 := quiz.NewRepository()
	if err := quiz.LoadAllBanks(root, repo2); err != nil {
		t.Fatalf("第二次 LoadAllBanks 失败: %v", err)
	}
	qs2, ok2 := repo2.GetBank("testtopic", "chapter1")
	if !ok2 || len(qs2) != 1 {
		t.Fatalf("期望 1 道题（修改后），但得到 %d", len(qs2))
	}
}
