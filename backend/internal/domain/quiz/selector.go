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
// 根据宪章标准，难度分布为：Easy 40%, Medium 40%, Hard 20%
// 此函数将按照难度分布从题目池中选择题目
// 如果某个难度分组的题目不足，会从其他难度中补充
func SelectQuestions(questions []YAMLQuestion, singleCount, multipleCount int, difficultyDist map[string]int) ([]YAMLQuestion, error) {
	logger.LogWithFields(context.Background(), "INFO", "quiz.select.start", map[string]interface{}{
		"total_questions": len(questions),
		"singleCount":     singleCount,
		"multipleCount":   multipleCount,
	})

	// 先按类型和难度分组
	typeDifficultyGroups := map[string]map[string][]YAMLQuestion{
		"single":   {"easy": {}, "medium": {}, "hard": {}},
		"multiple": {"easy": {}, "medium": {}, "hard": {}},
	}

	for _, q := range questions {
		if typeGroups, ok := typeDifficultyGroups[q.Type]; ok {
			if difficultyGroup, ok := typeGroups[q.Difficulty]; ok {
				typeDifficultyGroups[q.Type][q.Difficulty] = append(difficultyGroup, q)
			}
		}
	}

	// 为单选和多选题分配难度配额（Easy 40%, Medium 40%, Hard 20%）
	singleEasy := singleCount * 40 / 100
	singleMedium := singleCount * 40 / 100
	singleHard := singleCount - singleEasy - singleMedium

	multipleEasy := multipleCount * 40 / 100
	multipleMedium := multipleCount * 40 / 100
	multipleHard := multipleCount - multipleEasy - multipleMedium

	logger.LogWithFields(context.Background(), "INFO", "quiz.select.distribution", map[string]interface{}{
		"single_easy":     singleEasy,
		"single_medium":   singleMedium,
		"single_hard":     singleHard,
		"multiple_easy":   multipleEasy,
		"multiple_medium": multipleMedium,
		"multiple_hard":   multipleHard,
	})

	// 从各分组中随机抽取，如果某个分组不足，从其他难度补充
	selected := []YAMLQuestion{}

	// 单选题抽取
	singleSelected := selectFromGroup(typeDifficultyGroups["single"]["easy"], singleEasy)
	singleSelected = append(singleSelected, selectFromGroup(typeDifficultyGroups["single"]["medium"], singleMedium)...)
	singleSelected = append(singleSelected, selectFromGroup(typeDifficultyGroups["single"]["hard"], singleHard)...)

	// 如果单选题不足，从所有单选题中补充
	if len(singleSelected) < singleCount {
		remaining := singleCount - len(singleSelected)
		// 收集所有未被选中的单选题
		var allSingles []YAMLQuestion
		for _, diff := range []string{"easy", "medium", "hard"} {
			for _, q := range typeDifficultyGroups["single"][diff] {
				alreadySelected := false
				for _, s := range singleSelected {
					if s.ID == q.ID {
						alreadySelected = true
						break
					}
				}
				if !alreadySelected {
					allSingles = append(allSingles, q)
				}
			}
		}
		singleSelected = append(singleSelected, selectFromGroup(allSingles, remaining)...)
	}
	// 截取到需要的数量
	if len(singleSelected) > singleCount {
		singleSelected = singleSelected[:singleCount]
	}

	// 多选题抽取（同上逻辑）
	multipleSelected := selectFromGroup(typeDifficultyGroups["multiple"]["easy"], multipleEasy)
	multipleSelected = append(multipleSelected, selectFromGroup(typeDifficultyGroups["multiple"]["medium"], multipleMedium)...)
	multipleSelected = append(multipleSelected, selectFromGroup(typeDifficultyGroups["multiple"]["hard"], multipleHard)...)

	if len(multipleSelected) < multipleCount {
		remaining := multipleCount - len(multipleSelected)
		var allMultiples []YAMLQuestion
		for _, diff := range []string{"easy", "medium", "hard"} {
			for _, q := range typeDifficultyGroups["multiple"][diff] {
				alreadySelected := false
				for _, s := range multipleSelected {
					if s.ID == q.ID {
						alreadySelected = true
						break
					}
				}
				if !alreadySelected {
					allMultiples = append(allMultiples, q)
				}
			}
		}
		multipleSelected = append(multipleSelected, selectFromGroup(allMultiples, remaining)...)
	}
	if len(multipleSelected) > multipleCount {
		multipleSelected = multipleSelected[:multipleCount]
	}

	selected = append(selected, singleSelected...)
	selected = append(selected, multipleSelected...)

	// 最后再打乱顺序
	rand.Shuffle(len(selected), func(i, j int) { selected[i], selected[j] = selected[j], selected[i] })

	logger.LogWithFields(context.Background(), "INFO", "quiz.select.complete", map[string]interface{}{
		"selected": len(selected),
	})
	return selected, nil
}

// selectFromGroup 从题目组中随机抽取指定数量的题目
func selectFromGroup(group []YAMLQuestion, count int) []YAMLQuestion {
	if len(group) == 0 || count <= 0 {
		return []YAMLQuestion{}
	}
	// 打乱组内顺序
	rand.Shuffle(len(group), func(i, j int) { group[i], group[j] = group[j], group[i] })
	// 返回前count个，如果count超过组大小则返回全部
	if count > len(group) {
		return group
	}
	return group[:count]
}
