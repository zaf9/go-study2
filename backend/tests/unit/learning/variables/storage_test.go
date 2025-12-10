package variables_test

import (
	"testing"

	"go-study2/src/learning/variables"
)

// 测试存储主题内容加载。
func TestStorageContent(t *testing.T) {
	content, err := variables.LoadContent(variables.TopicStorage)
	if err != nil {
		t.Fatalf("加载内容失败: %v", err)
	}
	if content.Topic != variables.TopicStorage {
		t.Fatalf("主题不匹配: %s", content.Topic)
	}
	if content.Summary == "" {
		t.Fatalf("摘要应存在")
	}
	if len(content.Examples) == 0 {
		t.Fatalf("应包含至少一个示例")
	}
}

// 测试存储主题测验与评分。
func TestStorageQuiz(t *testing.T) {
	items, err := variables.LoadQuiz(variables.TopicStorage)
	if err != nil {
		t.Fatalf("加载测验失败: %v", err)
	}
	if len(items) < 2 {
		t.Fatalf("测验题目数量不足，获得 %d", len(items))
	}
	if err := variables.ValidateQuizItems(items); err != nil {
		t.Fatalf("测验校验失败: %v", err)
	}
	answers := map[string]string{}
	for _, item := range items {
		answers[item.ID] = item.Answer
	}
	result, err := variables.EvaluateQuiz(variables.TopicStorage, answers)
	if err != nil {
		t.Fatalf("评分失败: %v", err)
	}
	if result.Score != result.Total {
		t.Fatalf("评分应全对: %d/%d", result.Score, result.Total)
	}
}
