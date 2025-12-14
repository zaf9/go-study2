package quiz

import (
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
	sample := `questions:
  - id: test-001
    type: single
    difficulty: easy
    stem: "示例题"
    options: ["A","B"]
    answer: "A"
    explanation: "解析"
    topic: "lexical_elements"
    chapter: "comments"
`
	fpath := filepath.Join(topicDir, "comments.yaml")
	if err := ioutil.WriteFile(fpath, []byte(sample), 0o644); err != nil {
		t.Fatal(err)
	}

	repo := NewRepository()
	if err := LoadAllBanks(dir, repo); err != nil {
		t.Fatalf("LoadAllBanks failed: %v", err)
	}
	qs, ok := repo.GetBank("lexical_elements", "comments")
	if !ok || len(qs) != 1 {
		t.Fatalf("expected 1 question, got %v ok=%v", len(qs), ok)
	}
}
