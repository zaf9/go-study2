package quiz_test

import (
	"testing"

	appquiz "go-study2/internal/app/quiz"
)

func Test_MultipleAnswerOrderIndependence(t *testing.T) {
	engine := appquiz.NewScoringEngine()

	// 构造一道多选题，正确答案 ABC
	pq := appquiz.PreparedQuestion{
		View:          appquiz.QuestionDTO{ID: 1, Type: "multiple"},
		CorrectAnswer: []string{"A", "B", "C"},
	}

	questions := []appquiz.PreparedQuestion{pq}

	// 用户答案不同顺序但相同集合
	answers1 := map[int64][]string{1: {"A", "B", "C"}}
	answers2 := map[int64][]string{1: {"C", "B", "A"}}

	res1 := engine.Evaluate(questions, answers1)
	res2 := engine.Evaluate(questions, answers2)

	if res1.Score != res2.Score || res1.CorrectAnswers != res2.CorrectAnswers {
		t.Fatalf("评分应与答案顺序无关: %v vs %v", res1, res2)
	}
}
