package benchmark

import (
	"testing"
	"time"

	appquiz "go-study2/internal/app/quiz"
	quizdom "go-study2/internal/domain/quiz"
)

func BenchmarkQuestionPrepare(b *testing.B) {
	// 准备一个较大的题目集合
	var records []quizdom.QuizQuestion
	for i := 0; i < 1000; i++ {
		records = append(records, quizdom.QuizQuestion{
			ID:             int64(i + 1),
			Topic:          "types",
			Chapter:        "array",
			Type:           quizdom.QuestionTypeSingle,
			Difficulty:     quizdom.DifficultyMedium,
			Question:       "示例题",
			Options:        `["A","B","C","D"]`,
			CorrectAnswers: `["A"]`,
			Explanation:    "示例",
		})
	}
	mgr := appquiz.NewQuestionManager()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		start := time.Now()
		_, _, err := mgr.Prepare(records)
		if err != nil {
			b.Fatalf("Prepare 失败: %v", err)
		}
		_ = time.Since(start)
	}
}
