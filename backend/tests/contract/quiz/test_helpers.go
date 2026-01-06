package contract

import (
	"reflect"
	"unsafe"

	"go-study2/internal/app/http_server/handler"
	appquiz "go-study2/internal/app/quiz"
	quizdom "go-study2/internal/domain/quiz"
)

// setQuizService 使用反射注入测验服务
func setQuizService(h *handler.Handler, svc *appquiz.Service) {
	val := reflect.ValueOf(h).Elem().FieldByName("quizService")
	ptr := unsafe.Pointer(val.UnsafeAddr())
	reflect.NewAt(val.Type(), ptr).Elem().Set(reflect.ValueOf(svc))
}

// SetupTestYAMLRepository 创建并初始化测试用的YAML仓储（8道题：4单选+4多选）
func SetupTestYAMLRepository(topic, chapter string) *quizdom.QuizRepository {
	yamlRepo := quizdom.NewRepository()
	testQuestions := []quizdom.YAMLQuestion{
		{ID: "test_q1", Type: "single", Difficulty: "easy", Stem: "Single choice 1", Options: []string{"A", "B"}, Answer: "A", Explanation: "exp1", Topic: topic, Chapter: chapter},
		{ID: "test_q2", Type: "single", Difficulty: "easy", Stem: "Single choice 2", Options: []string{"A", "B"}, Answer: "A", Explanation: "exp2", Topic: topic, Chapter: chapter},
		{ID: "test_q3", Type: "single", Difficulty: "easy", Stem: "Single choice 3", Options: []string{"A", "B"}, Answer: "A", Explanation: "exp3", Topic: topic, Chapter: chapter},
		{ID: "test_q4", Type: "single", Difficulty: "easy", Stem: "Single choice 4", Options: []string{"A", "B"}, Answer: "A", Explanation: "exp4", Topic: topic, Chapter: chapter},
		{ID: "test_q5", Type: "multiple", Difficulty: "medium", Stem: "Multiple 1", Options: []string{"A", "B", "C"}, Answer: "AB", Explanation: "exp5", Topic: topic, Chapter: chapter},
		{ID: "test_q6", Type: "multiple", Difficulty: "medium", Stem: "Multiple 2", Options: []string{"A", "B", "C"}, Answer: "AB", Explanation: "exp6", Topic: topic, Chapter: chapter},
		{ID: "test_q7", Type: "multiple", Difficulty: "medium", Stem: "Multiple 3", Options: []string{"A", "B", "C"}, Answer: "AB", Explanation: "exp7", Topic: topic, Chapter: chapter},
		{ID: "test_q8", Type: "multiple", Difficulty: "medium", Stem: "Multiple 4", Options: []string{"A", "B", "C"}, Answer: "AB", Explanation: "exp8", Topic: topic, Chapter: chapter},
	}
	yamlRepo.AddBank(topic, chapter, testQuestions)
	return yamlRepo
}
