package quiz

import (
	"context"
	"math/rand"
	"time"

	"go-study2/internal/infrastructure/logger"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// SelectQuestions 从给定题目集中按类型和难度分布抽取指定数量
func SelectQuestions(questions []YAMLQuestion, singleCount, multipleCount int, difficultyDist map[string]int) ([]YAMLQuestion, error) {
	logger.LogWithFields(context.Background(), "INFO", "quiz.select.start", map[string]interface{}{
		"total_questions": len(questions),
		"singleCount":     singleCount,
		"multipleCount":   multipleCount,
	})
	// 简单实现：先按类型分组，再随机抽取数量
	var singles []YAMLQuestion
	var multiples []YAMLQuestion
	for _, q := range questions {
		if q.Type == "single" {
			singles = append(singles, q)
		} else if q.Type == "multiple" {
			multiples = append(multiples, q)
		}
	}
	if len(singles) < singleCount || len(multiples) < multipleCount {
		// 不足则返回全部可用
	}
	rand.Shuffle(len(singles), func(i, j int) { singles[i], singles[j] = singles[j], singles[i] })
	rand.Shuffle(len(multiples), func(i, j int) { multiples[i], multiples[j] = multiples[j], multiples[i] })

	sel := []YAMLQuestion{}
	if singleCount > len(singles) {
		sel = append(sel, singles...)
	} else {
		sel = append(sel, singles[:singleCount]...)
	}
	if multipleCount > len(multiples) {
		sel = append(sel, multiples...)
	} else {
		sel = append(sel, multiples[:multipleCount]...)
	}

	// 最后再打乱顺序
	rand.Shuffle(len(sel), func(i, j int) { sel[i], sel[j] = sel[j], sel[i] })
	logger.LogWithFields(context.Background(), "INFO", "quiz.select.complete", map[string]interface{}{
		"selected": len(sel),
	})
	return sel, nil
}
