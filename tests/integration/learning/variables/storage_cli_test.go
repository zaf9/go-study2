package variables_integration_test

import (
	"testing"

	"go-study2/src/learning/variables"
	"go-study2/src/learning/variables/cli"
)

// 测试 CLI 菜单与测验流程。
func TestStorageCLIMenu(t *testing.T) {
	menu := cli.NewMenu()

	topics := menu.ListTopics()
	if len(topics) == 0 {
		t.Fatalf("应返回至少一个主题")
	}

	content, err := menu.ShowContent("storage")
	if err != nil {
		t.Fatalf("ShowContent 失败: %v", err)
	}
	if content.Topic != variables.TopicStorage {
		t.Fatalf("返回主题错误: %s", content.Topic)
	}

	items, err := menu.StartQuiz("storage")
	if err != nil {
		t.Fatalf("StartQuiz 失败: %v", err)
	}
	if len(items) == 0 {
		t.Fatalf("应返回题目")
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
		t.Fatalf("应全部答对: %d/%d", result.Score, result.Total)
	}
}
