package variables_test

import (
	"strings"
	"testing"

	"go-study2/src/learning/variables"
)

// 测试零值主题内容与测验。
func TestZeroValueContentAndQuiz(t *testing.T) {
	content, err := variables.LoadContent(variables.TopicZero)
	if err != nil {
		t.Fatalf("加载内容失败: %v", err)
	}
	if !strings.Contains(content.Summary, "零值") {
		t.Fatalf("摘要应提及零值")
	}

	items, err := variables.LoadQuiz(variables.TopicZero)
	if err != nil {
		t.Fatalf("加载测验失败: %v", err)
	}
	if len(items) == 0 {
		t.Fatalf("零值测验题目为空")
	}
	answers := map[string]string{}
	for _, item := range items {
		answers[item.ID] = item.Answer
	}
	result, err := variables.EvaluateQuiz(variables.TopicZero, answers)
	if err != nil {
		t.Fatalf("评分失败: %v", err)
	}
	if result.Score != result.Total {
		t.Fatalf("应全部答对: %d/%d", result.Score, result.Total)
	}
}
