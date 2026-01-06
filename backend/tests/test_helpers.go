package tests

import (
	quizdom "go-study2/internal/domain/quiz"
)

// CreateTestYAMLQuestions 创建测试用的YAML格式题目
// 返回8道题目：4道单选 + 4道多选，满足随机抽题的最低要求
func CreateTestYAMLQuestions(topic, chapter string) []quizdom.YAMLQuestion {
	return []quizdom.YAMLQuestion{
		{ID: "test_q1", Type: "single", Difficulty: "easy", Stem: "Single choice question 1", Options: []string{"Option A", "Option B"}, Answer: "A", Explanation: "Explanation 1", Topic: topic, Chapter: chapter},
		{ID: "test_q2", Type: "single", Difficulty: "easy", Stem: "Single choice question 2", Options: []string{"Option A", "Option B"}, Answer: "A", Explanation: "Explanation 2", Topic: topic, Chapter: chapter},
		{ID: "test_q3", Type: "single", Difficulty: "easy", Stem: "Single choice question 3", Options: []string{"Option A", "Option B"}, Answer: "A", Explanation: "Explanation 3", Topic: topic, Chapter: chapter},
		{ID: "test_q4", Type: "single", Difficulty: "easy", Stem: "Single choice question 4", Options: []string{"Option A", "Option B"}, Answer: "A", Explanation: "Explanation 4", Topic: topic, Chapter: chapter},
		{ID: "test_q5", Type: "multiple", Difficulty: "medium", Stem: "Multiple choice question 1", Options: []string{"Option A", "Option B", "Option C"}, Answer: "AB", Explanation: "Explanation 5", Topic: topic, Chapter: chapter},
		{ID: "test_q6", Type: "multiple", Difficulty: "medium", Stem: "Multiple choice question 2", Options: []string{"Option A", "Option B", "Option C"}, Answer: "AB", Explanation: "Explanation 6", Topic: topic, Chapter: chapter},
		{ID: "test_q7", Type: "multiple", Difficulty: "medium", Stem: "Multiple choice question 3", Options: []string{"Option A", "Option B", "Option C"}, Answer: "AB", Explanation: "Explanation 7", Topic: topic, Chapter: chapter},
		{ID: "test_q8", Type: "multiple", Difficulty: "medium", Stem: "Multiple choice question 4", Options: []string{"Option A", "Option B", "Option C"}, Answer: "AB", Explanation: "Explanation 8", Topic: topic, Chapter: chapter},
	}
}

// SetupTestYAMLRepository 创建并初始化测试用的YAML仓储
func SetupTestYAMLRepository(topic, chapter string) *quizdom.QuizRepository {
	yamlRepo := quizdom.NewRepository()
	testQuestions := CreateTestYAMLQuestions(topic, chapter)
	yamlRepo.AddBank(topic, chapter, testQuestions)
	return yamlRepo
}
