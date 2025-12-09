package types_test

import (
	"testing"

	"go-study2/src/learning/types"
)

// T015: 测验评分与重做逻辑单测
func TestEvaluateComprehensiveQuiz(t *testing.T) {
	quiz := types.LoadComprehensiveQuiz()
	if len(quiz) < 5 {
		t.Fatalf("综合测验题目不足，实际 %d", len(quiz))
	}

	answers := map[string]string{}
	for _, q := range quiz {
		answers[q.ID] = q.Answer
	}
	result, err := types.EvaluateComprehensiveQuiz(answers)
	if err != nil {
		t.Fatalf("评分失败: %v", err)
	}
	if result.Score != result.Total {
		t.Fatalf("预期满分 %d，实际 %d", result.Total, result.Score)
	}
}
