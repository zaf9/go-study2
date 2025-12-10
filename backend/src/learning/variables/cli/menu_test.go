package cli

import (
	"testing"
)

// 覆盖 CLI 菜单流程。
func TestMenuFlow(t *testing.T) {
	menu := NewMenu()
	if len(menu.ListTopics()) == 0 {
		t.Fatalf("应返回主题列表")
	}
	content, err := menu.ShowContent("storage")
	if err != nil {
		t.Fatalf("ShowContent 失败: %v", err)
	}
	if content.Topic == "" {
		t.Fatalf("内容缺少主题")
	}
	items, err := menu.StartQuiz("storage")
	if err != nil {
		t.Fatalf("StartQuiz 失败: %v", err)
	}
	answers := map[string]string{}
	for _, item := range items {
		answers[item.ID] = item.Answer
	}
	result, err := menu.SubmitQuiz("storage", answers)
	if err != nil {
		t.Fatalf("SubmitQuiz 失败: %v", err)
	}
	if result.Score != result.Total {
		t.Fatalf("评分应全对: %d/%d", result.Score, result.Total)
	}
}
