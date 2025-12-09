package types_test

import (
	"testing"

	"go-study2/src/learning/types"
)

// T007: 内容覆盖与结构校验
func TestContentRegistry(t *testing.T) {
	topics := types.AllTopics()
	if len(topics) < 12 {
		t.Fatalf("预期至少 12 个子主题，实际 %d", len(topics))
	}

	for _, topic := range topics {
		content, err := types.LoadContent(topic)
		if err != nil {
			t.Fatalf("加载主题 %s 失败: %v", topic, err)
		}
		if content.Concept.Title == "" || content.Concept.ID == "" {
			t.Fatalf("主题 %s 缺少标题或 ID", topic)
		}
		if len(content.Examples) == 0 {
			t.Fatalf("主题 %s 缺少示例", topic)
		}
		quiz, err := types.LoadQuiz(topic)
		if err != nil {
			t.Fatalf("主题 %s 缺少测验: %v", topic, err)
		}
		if len(quiz) == 0 {
			t.Fatalf("主题 %s 测验为空", topic)
		}
	}
}

func TestSearchReferences(t *testing.T) {
	results, err := types.SearchReferences("map")
	if err != nil {
		t.Fatalf("检索失败: %v", err)
	}
	if len(results) == 0 {
		t.Fatalf("检索 map 结果为空")
	}
}
