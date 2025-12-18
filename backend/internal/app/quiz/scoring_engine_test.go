package quiz

import (
	"testing"
)

func TestScoringEngine_EvaluateWithTotal(t *testing.T) {
	engine := NewScoringEngine()

	questions := []PreparedQuestion{
		{
			View: QuestionDTO{
				ID:   1,
				Type: "single",
			},
			CorrectAnswer: []string{"A"},
		},
		{
			View: QuestionDTO{
				ID:   2,
				Type: "multiple",
			},
			CorrectAnswer: []string{"A", "B"},
		},
	}

	// 场景 1: 全部正确
	answers := map[int64][]string{
		1: {"A"},
		2: {"A", "B"},
	}
	result := engine.EvaluateWithTotal(questions, answers, 2)
	if result.Score != 100 {
		t.Errorf("expected score 100, got %d", result.Score)
	}
	if !result.Passed {
		t.Error("expected passed to be true")
	}
	if result.CorrectAnswers != 2 {
		t.Errorf("expected 2 correct answers, got %d", result.CorrectAnswers)
	}

	// 场景 2: 部分正确 (多选题半对)
	// partialScore 为 1/2 = 0.5
	// 总分为 (1 + 0.5) / 2 * 100 = 75
	answers = map[int64][]string{
		1: {"A"},
		2: {"A"},
	}
	result = engine.EvaluateWithTotal(questions, answers, 2)
	if result.Score != 75 {
		t.Errorf("expected score 75, got %d", result.Score)
	}
	if !result.Passed {
		t.Error("expected passed to be true")
	}
	if result.CorrectAnswers != 1 {
		t.Errorf("expected 1 correct answer (single only), got %d", result.CorrectAnswers)
	}

	// 场景 3: 不及格
	answers = map[int64][]string{
		1: {"B"},
		2: {"A"},
	}
	result = engine.EvaluateWithTotal(questions, answers, 2)
	if result.Score != 25 {
		t.Errorf("expected score 25, got %d", result.Score)
	}
	if result.Passed {
		t.Error("expected passed to be false")
	}
}
