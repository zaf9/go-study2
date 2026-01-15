package quiz

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadAllBanks(t *testing.T) {
	dir, err := ioutil.TempDir("", "quiztest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	// create topic dir and yaml file
	topicDir := filepath.Join(dir, "lexical_elements")
	os.MkdirAll(topicDir, 0o755)

	// 创建符合宪章标准的题库：
	// - 30题（最少）
	// - 单选15题(50%)，多选15题(50%)
	// - Easy 12题(40%), Medium 12题(40%), Hard 6题(20%)
	sample := `questions:
`
	// 单选简单题 6 道
	for i := 0; i < 6; i++ {
		sample += fmt.Sprintf(`  - id: test-single-easy-%03d
    type: single
    difficulty: easy
    stem: "单选简单题 %d"
    options: ["A","B","C","D"]
    answer: "A"
    explanation: "解析"
    topic: "lexical_elements"
    chapter: "comments"
`, i+1, i+1)
	}
	// 单选中等题 6 道
	for i := 0; i < 6; i++ {
		sample += fmt.Sprintf(`  - id: test-single-medium-%03d
    type: single
    difficulty: medium
    stem: "单选中等题 %d"
    options: ["A","B","C","D"]
    answer: "A"
    explanation: "解析"
    topic: "lexical_elements"
    chapter: "comments"
`, i+1, i+1)
	}
	// 单选困难题 3 道
	for i := 0; i < 3; i++ {
		sample += fmt.Sprintf(`  - id: test-single-hard-%03d
    type: single
    difficulty: hard
    stem: "单选困难题 %d"
    options: ["A","B","C","D"]
    answer: "A"
    explanation: "解析"
    topic: "lexical_elements"
    chapter: "comments"
`, i+1, i+1)
	}
	// 多选简单题 6 道
	for i := 0; i < 6; i++ {
		sample += fmt.Sprintf(`  - id: test-multiple-easy-%03d
    type: multiple
    difficulty: easy
    stem: "多选简单题 %d"
    options: ["A","B","C","D"]
    answer: "AB"
    explanation: "解析"
    topic: "lexical_elements"
    chapter: "comments"
`, i+1, i+1)
	}
	// 多选中等题 6 道
	for i := 0; i < 6; i++ {
		sample += fmt.Sprintf(`  - id: test-multiple-medium-%03d
    type: multiple
    difficulty: medium
    stem: "多选中等题 %d"
    options: ["A","B","C","D"]
    answer: "AB"
    explanation: "解析"
    topic: "lexical_elements"
    chapter: "comments"
`, i+1, i+1)
	}
	// 多选困难题 3 道
	for i := 0; i < 3; i++ {
		sample += fmt.Sprintf(`  - id: test-multiple-hard-%03d
    type: multiple
    difficulty: hard
    stem: "多选困难题 %d"
    options: ["A","B","C","D"]
    answer: "AB"
    explanation: "解析"
    topic: "lexical_elements"
    chapter: "comments"
`, i+1, i+1)
	}

	fpath := filepath.Join(topicDir, "comments.yaml")
	if err := ioutil.WriteFile(fpath, []byte(sample), 0o644); err != nil {
		t.Fatal(err)
	}

	repo := NewRepository()
	if err := LoadAllBanks(dir, repo); err != nil {
		t.Fatalf("LoadAllBanks failed: %v", err)
	}
	qs, ok := repo.GetBank("lexical_elements", "comments")
	if !ok || len(qs) != 30 {
		t.Fatalf("expected 30 questions, got %v ok=%v", len(qs), ok)
	}
}
