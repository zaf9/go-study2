package variables_integration_test

import (
	"testing"

	"go-study2/src/learning/variables/cli"
)

// 测试 CLI 零值主题。
func TestZeroCLIMenu(t *testing.T) {
	menu := cli.NewMenu()
	content, err := menu.ShowContent("zero")
	if err != nil {
		t.Fatalf("ShowContent 失败: %v", err)
	}
	if content.Title == "" {
		t.Fatalf("内容标题缺失")
	}
	items, err := menu.StartQuiz("zero")
	if err != nil {
		t.Fatalf("StartQuiz 失败: %v", err)
	}
	if len(items) == 0 {
		t.Fatalf("题目为空")
	}
	answers := map[string]string{}
	for _, item := range items {
		answers[item.ID] = item.Answer
	}
	result, err := menu.SubmitQuiz("zero", answers)
	if err != nil {
		t.Fatalf("SubmitQuiz 失败: %v", err)
	}
	if result.Score != result.Total {
		t.Fatalf("应全部答对: %d/%d", result.Score, result.Total)
	}
}
