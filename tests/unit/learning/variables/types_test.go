package variables_test

import (
	"testing"

	"go-study2/src/learning/variables"
)

// 测试静态与动态类型内容与测验。
func TestStaticAndDynamicContent(t *testing.T) {
	staticContent, err := variables.LoadContent(variables.TopicStatic)
	if err != nil {
		t.Fatalf("加载静态类型内容失败: %v", err)
	}
	if staticContent.Title == "" || staticContent.Summary == "" {
		t.Fatalf("静态类型内容应包含标题与摘要")
	}

	dynamicContent, err := variables.LoadContent(variables.TopicDynamic)
	if err != nil {
		t.Fatalf("加载动态类型内容失败: %v", err)
	}
	if dynamicContent.Title == "" || dynamicContent.Summary == "" {
		t.Fatalf("动态类型内容应包含标题与摘要")
	}
}

// 测试静态与动态类型测验评分。
func TestStaticAndDynamicQuiz(t *testing.T) {
	topics := []variables.Topic{variables.TopicStatic, variables.TopicDynamic}
	for _, topic := range topics {
		items, err := variables.LoadQuiz(topic)
		if err != nil {
			t.Fatalf("加载测验失败: %v", err)
		}
		if err := variables.ValidateQuizItems(items); err != nil {
			t.Fatalf("测验校验失败: %v", err)
		}
		answers := map[string]string{}
		for _, item := range items {
			answers[item.ID] = item.Answer
		}
		result, err := variables.EvaluateQuiz(topic, answers)
		if err != nil {
			t.Fatalf("评分失败: %v", err)
		}
		if result.Score != result.Total {
			t.Fatalf("应全部答对: %d/%d", result.Score, result.Total)
		}
	}
}
