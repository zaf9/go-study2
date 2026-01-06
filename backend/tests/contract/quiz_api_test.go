package contract

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"go-study2/internal/app/http_server/handler"
	appquiz "go-study2/internal/app/quiz"
	"go-study2/internal/config"
	quizdom "go-study2/internal/domain/quiz"
	infrarepo "go-study2/internal/infra/repository"
	"go-study2/internal/infrastructure/database"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
)

// T053: 测试GET /api/v1/quiz/:topic/:chapter返回正确的题目结构
func TestQuizAPI_GetQuiz_ReturnsCorrectStructure(t *testing.T) {
	ctx := gctx.New()
	tmp := t.TempDir()
	dbPath := filepath.ToSlash(filepath.Join(tmp, "quiz_api_test.db"))
	cfgDB := config.DatabaseConfig{
		Type:    "sqlite3",
		Path:    dbPath,
		Pragmas: []string{"journal_mode=WAL", "busy_timeout=5000", "synchronous=NORMAL", "foreign_keys=ON"},
	}
	db, err := database.Init(ctx, cfgDB)
	if err != nil {
		t.Fatalf("init db failed: %v", err)
	}
	defer db.Close(ctx)

	// 清空并初始化测试数据（quiz_questions表已废弃，题目从YAML加载）
	if _, err := db.Exec(ctx, "DELETE FROM quiz_sessions"); err != nil {
		t.Fatalf("clear quiz_sessions failed: %v", err)
	}
	if _, err := db.Exec(ctx, "DELETE FROM quiz_attempts"); err != nil {
		t.Fatalf("clear quiz_attempts failed: %v", err)
	}

	now := time.Now()
	// 插入测试用户
	if _, err := db.Model("users").Data(map[string]interface{}{
		"id": 1, "username": "testuser", "email": "test@example.com",
		"password_hash": "hash", "created_at": now, "updated_at": now,
	}).Insert(); err != nil {
		t.Fatalf("insert user failed: %v", err)
	}

	// 创建handler和service
	h := handler.New()

	// 创建YAML仓储并添加测试题目（至少4单选+4多选）
	yamlRepo := quizdom.NewRepository()
	testYAMLQuestions := []quizdom.YAMLQuestion{
		{ID: "ct1", Type: "single", Difficulty: "easy", Stem: "Go语言中布尔类型的零值是什么?", Options: []string{"true", "false", "nil", "0"}, Answer: "B", Explanation: "布尔类型的零值是false", Topic: "constants", Chapter: "boolean"},
		{ID: "ct2", Type: "single", Difficulty: "easy", Stem: "Single 2", Options: []string{"A", "B"}, Answer: "A", Explanation: "exp2", Topic: "constants", Chapter: "boolean"},
		{ID: "ct3", Type: "single", Difficulty: "easy", Stem: "Single 3", Options: []string{"A", "B"}, Answer: "A", Explanation: "exp3", Topic: "constants", Chapter: "boolean"},
		{ID: "ct4", Type: "single", Difficulty: "easy", Stem: "Single 4", Options: []string{"A", "B"}, Answer: "A", Explanation: "exp4", Topic: "constants", Chapter: "boolean"},
		{ID: "ct5", Type: "multiple", Difficulty: "medium", Stem: "Multiple 1", Options: []string{"A", "B", "C"}, Answer: "AB", Explanation: "exp5", Topic: "constants", Chapter: "boolean"},
		{ID: "ct6", Type: "multiple", Difficulty: "medium", Stem: "Multiple 2", Options: []string{"A", "B", "C"}, Answer: "AB", Explanation: "exp6", Topic: "constants", Chapter: "boolean"},
		{ID: "ct7", Type: "multiple", Difficulty: "medium", Stem: "Multiple 3", Options: []string{"A", "B", "C"}, Answer: "AB", Explanation: "exp7", Topic: "constants", Chapter: "boolean"},
		{ID: "ct8", Type: "multiple", Difficulty: "medium", Stem: "Multiple 4", Options: []string{"A", "B", "C"}, Answer: "AB", Explanation: "exp8", Topic: "constants", Chapter: "boolean"},
	}
	yamlRepo.AddBank("constants", "boolean", testYAMLQuestions)

	repoImpl := infrarepo.NewQuizRepository(db)
	svc := appquiz.NewService(yamlRepo, repoImpl)

	type quizSetter interface{ SetQuizService(*appquiz.Service) }
	if s, ok := interface{}(h).(quizSetter); ok {
		s.SetQuizService(svc)
	}

	// 启动测试服务器
	s := ghttp.GetServer(fmt.Sprintf("quiz-api-test-%d", time.Now().UnixNano()))
	s.BindMiddlewareDefault(func(r *ghttp.Request) {
		if uid := r.Header.Get("X-User-ID"); uid != "" {
			r.SetCtxVar("user_id", 1)
		}
		r.Middleware.Next()
	})
	s.Group("/api/v1", func(group *ghttp.RouterGroup) {
		group.GET("/quiz/{topic}/{chapter}", h.GetQuiz)
	})
	go s.Start()
	defer s.Shutdown()
	time.Sleep(30 * time.Millisecond)

	// 发送请求
	req := httptest.NewRequest(http.MethodGet, "/api/v1/quiz/constants/boolean", nil)
	req.Header.Set("X-User-ID", "1")
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)

	// 验证响应
	if w.Code != 200 {
		t.Fatalf("expected 200 got %d body=%s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal resp failed: %v", err)
	}

	// 验证契约结构
	if _, ok := resp["code"]; !ok {
		t.Fatalf("response missing 'code' field")
	}
	if _, ok := resp["message"]; !ok {
		t.Fatalf("response missing 'message' field")
	}
	data, ok := resp["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("response 'data' missing or wrong type")
	}

	// 验证sessionId存在
	if sessionID, ok := data["sessionId"].(string); !ok || sessionID == "" {
		t.Fatalf("sessionId missing or empty")
	}

	// 验证questions数组
	questions, ok := data["questions"].([]interface{})
	if !ok {
		t.Fatalf("questions missing or wrong type")
	}
	if len(questions) == 0 {
		t.Fatalf("questions array is empty")
	}

	// 验证第一个题目的结构
	q0, _ := questions[0].(map[string]interface{})
	requiredFields := []string{"id", "type", "difficulty", "question", "options"}
	for _, field := range requiredFields {
		if _, ok := q0[field]; !ok {
			t.Fatalf("question missing field %s", field)
		}
	}
}

// T054: 测试POST /api/v1/quiz/submit返回正确的评分结果
func TestQuizAPI_Submit_ReturnsCorrectScore(t *testing.T) {
	ctx := gctx.New()
	tmp := t.TempDir()
	dbPath := filepath.ToSlash(filepath.Join(tmp, "quiz_submit_test.db"))
	cfgDB := config.DatabaseConfig{
		Type:    "sqlite3",
		Path:    dbPath,
		Pragmas: []string{"journal_mode=WAL", "busy_timeout=5000", "synchronous=NORMAL", "foreign_keys=ON"},
	}
	db, err := database.Init(ctx, cfgDB)
	if err != nil {
		t.Fatalf("init db failed: %v", err)
	}
	defer db.Close(ctx)

	// 清空表
	for _, table := range []string{"quiz_sessions", "quiz_attempts", "users"} {
		if _, err := db.Exec(ctx, fmt.Sprintf("DELETE FROM %s", table)); err != nil {
			t.Fatalf("clear %s failed: %v", table, err)
		}
	}

	now := time.Now()
	// 插入测试用户
	if _, err := db.Model("users").Data(map[string]interface{}{
		"id": 1, "username": "testuser", "email": "test@example.com",
		"password_hash": "hash", "created_at": now, "updated_at": now,
	}).Insert(); err != nil {
		t.Fatalf("insert user failed: %v", err)
	}

	// 创建handler和服务
	h := handler.New()

	// 创建YAML仓储并添加测试题目（至少4单选+4多选）
	yamlRepo := quizdom.NewRepository()
	testYAMLQuestions := []quizdom.YAMLQuestion{
		{ID: "ct1", Type: "single", Difficulty: "easy", Stem: "Go语言中布尔类型的零值是什么?", Options: []string{"true", "false", "nil", "0"}, Answer: "B", Explanation: "布尔类型的零值是false", Topic: "constants", Chapter: "boolean"},
		{ID: "ct2", Type: "single", Difficulty: "easy", Stem: "Single 2", Options: []string{"A", "B"}, Answer: "A", Explanation: "exp2", Topic: "constants", Chapter: "boolean"},
		{ID: "ct3", Type: "single", Difficulty: "easy", Stem: "Single 3", Options: []string{"A", "B"}, Answer: "A", Explanation: "exp3", Topic: "constants", Chapter: "boolean"},
		{ID: "ct4", Type: "single", Difficulty: "easy", Stem: "Single 4", Options: []string{"A", "B"}, Answer: "A", Explanation: "exp4", Topic: "constants", Chapter: "boolean"},
		{ID: "ct5", Type: "multiple", Difficulty: "medium", Stem: "Multiple 1", Options: []string{"A", "B", "C"}, Answer: "AB", Explanation: "exp5", Topic: "constants", Chapter: "boolean"},
		{ID: "ct6", Type: "multiple", Difficulty: "medium", Stem: "Multiple 2", Options: []string{"A", "B", "C"}, Answer: "AB", Explanation: "exp6", Topic: "constants", Chapter: "boolean"},
		{ID: "ct7", Type: "multiple", Difficulty: "medium", Stem: "Multiple 3", Options: []string{"A", "B", "C"}, Answer: "AB", Explanation: "exp7", Topic: "constants", Chapter: "boolean"},
		{ID: "ct8", Type: "multiple", Difficulty: "medium", Stem: "Multiple 4", Options: []string{"A", "B", "C"}, Answer: "AB", Explanation: "exp8", Topic: "constants", Chapter: "boolean"},
	}
	yamlRepo.AddBank("constants", "boolean", testYAMLQuestions)

	repoImpl := infrarepo.NewQuizRepository(db)
	svc := appquiz.NewService(yamlRepo, repoImpl)

	type quizSetter interface{ SetQuizService(*appquiz.Service) }
	if s, ok := interface{}(h).(quizSetter); ok {
		s.SetQuizService(svc)
	}

	// 先获取quiz创建session
	payload, err := svc.GetQuizQuestions(ctx, 1, "constants", "boolean")
	if err != nil {
		t.Fatalf("GetQuizQuestions failed: %v", err)
	}
	sessionID := payload.SessionID
	if len(payload.Questions) == 0 {
		t.Fatalf("no questions returned")
	}

	// 构建正确答案映射（从YAML题目）
	correctAnswerMap := make(map[string]string)
	for _, yq := range testYAMLQuestions {
		correctAnswerMap[yq.Stem] = yq.Answer
	}

	// 构建所有题目的正确答案
	var answers []string
	for _, q := range payload.Questions {
		correctAnswer := correctAnswerMap[q.Question]
		// 将答案字符串转换为字符数组（例如 "AB" -> ["A", "B"]）
		userAnswers := make([]string, 0, len(correctAnswer))
		for _, ch := range correctAnswer {
			userAnswers = append(userAnswers, fmt.Sprintf(`"%s"`, string(ch)))
		}
		answers = append(answers, fmt.Sprintf(`{"questionId":%d,"userAnswers":[%s]}`, q.ID, strings.Join(userAnswers, ",")))
	}
	answersJSON := fmt.Sprintf("[%s]", strings.Join(answers, ","))

	// 启动测试服务器
	s := ghttp.GetServer(fmt.Sprintf("quiz-submit-test-%d", time.Now().UnixNano()))
	s.BindMiddlewareDefault(func(r *ghttp.Request) {
		if uid := r.Header.Get("X-User-ID"); uid != "" {
			r.SetCtxVar("user_id", 1)
		}
		r.Middleware.Next()
	})
	s.Group("/api/v1", func(group *ghttp.RouterGroup) {
		group.POST("/quiz/submit", h.SubmitQuiz)
	})
	go s.Start()
	defer s.Shutdown()
	time.Sleep(30 * time.Millisecond)

	// 提交所有题目的正确答案
	submitPayload := fmt.Sprintf(`{
		"sessionId": "%s",
		"topic": "constants",
		"chapter": "boolean",
		"answers": %s,
		"durationMs": 5000
	}`, sessionID, answersJSON)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/quiz/submit", strings.NewReader(submitPayload))
	req.Header.Set("X-User-ID", "1")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected 200 got %d body=%s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal resp failed: %v", err)
	}

	data, ok := resp["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("response 'data' missing")
	}

	// 验证得分 - 答对了应该是100分
	score, ok := data["score"].(float64)
	if !ok {
		t.Fatalf("score missing or wrong type")
	}
	if score != 100 {
		t.Errorf("expected score 100, got %.0f", score)
	}

	// 验证通过状态
	passed, ok := data["passed"].(bool)
	if !ok || !passed {
		t.Errorf("expected passed=true, got %v", passed)
	}
}

// T055: 测试重复提交被拒绝(幂等性)
func TestQuizAPI_Submit_RejectsDuplicateSubmission(t *testing.T) {
	ctx := gctx.New()
	tmp := t.TempDir()
	dbPath := filepath.ToSlash(filepath.Join(tmp, "quiz_duplicate_test.db"))
	cfgDB := config.DatabaseConfig{
		Type:    "sqlite3",
		Path:    dbPath,
		Pragmas: []string{"journal_mode=WAL", "busy_timeout=5000", "synchronous=NORMAL", "foreign_keys=ON"},
	}
	db, err := database.Init(ctx, cfgDB)
	if err != nil {
		t.Fatalf("init db failed: %v", err)
	}
	defer db.Close(ctx)

	// 清空表
	for _, table := range []string{"quiz_sessions", "quiz_attempts", "users"} {
		if _, err := db.Exec(ctx, fmt.Sprintf("DELETE FROM %s", table)); err != nil {
			t.Fatalf("clear %s failed: %v", table, err)
		}
	}

	now := time.Now()
	if _, err := db.Model("users").Data(map[string]interface{}{
		"id": 1, "username": "testuser", "email": "test@example.com",
		"password_hash": "hash", "created_at": now, "updated_at": now,
	}).Insert(); err != nil {
		t.Fatalf("insert user failed: %v", err)
	}

	h := handler.New()

	// 创建YAML仓储并添加测试题目（至少4单选+4多选）
	yamlRepo := quizdom.NewRepository()
	testYAMLQuestions := []quizdom.YAMLQuestion{
		{ID: "ct1", Type: "single", Difficulty: "easy", Stem: "测试题目", Options: []string{"选项A", "选项B"}, Answer: "A", Explanation: "解析", Topic: "constants", Chapter: "boolean"},
		{ID: "ct2", Type: "single", Difficulty: "easy", Stem: "Single 2", Options: []string{"A", "B"}, Answer: "A", Explanation: "exp2", Topic: "constants", Chapter: "boolean"},
		{ID: "ct3", Type: "single", Difficulty: "easy", Stem: "Single 3", Options: []string{"A", "B"}, Answer: "A", Explanation: "exp3", Topic: "constants", Chapter: "boolean"},
		{ID: "ct4", Type: "single", Difficulty: "easy", Stem: "Single 4", Options: []string{"A", "B"}, Answer: "A", Explanation: "exp4", Topic: "constants", Chapter: "boolean"},
		{ID: "ct5", Type: "multiple", Difficulty: "medium", Stem: "Multiple 1", Options: []string{"A", "B", "C"}, Answer: "AB", Explanation: "exp5", Topic: "constants", Chapter: "boolean"},
		{ID: "ct6", Type: "multiple", Difficulty: "medium", Stem: "Multiple 2", Options: []string{"A", "B", "C"}, Answer: "AB", Explanation: "exp6", Topic: "constants", Chapter: "boolean"},
		{ID: "ct7", Type: "multiple", Difficulty: "medium", Stem: "Multiple 3", Options: []string{"A", "B", "C"}, Answer: "AB", Explanation: "exp7", Topic: "constants", Chapter: "boolean"},
		{ID: "ct8", Type: "multiple", Difficulty: "medium", Stem: "Multiple 4", Options: []string{"A", "B", "C"}, Answer: "AB", Explanation: "exp8", Topic: "constants", Chapter: "boolean"},
	}
	yamlRepo.AddBank("constants", "boolean", testYAMLQuestions)

	repoImpl := infrarepo.NewQuizRepository(db)
	svc := appquiz.NewService(yamlRepo, repoImpl)

	type quizSetter interface{ SetQuizService(*appquiz.Service) }
	if s, ok := interface{}(h).(quizSetter); ok {
		s.SetQuizService(svc)
	}

	payload, err := svc.GetQuizQuestions(ctx, 1, "constants", "boolean")
	if err != nil {
		t.Fatalf("GetQuizQuestions failed: %v", err)
	}
	sessionID := payload.SessionID
	if len(payload.Questions) == 0 {
		t.Fatalf("no questions returned")
	}
	// 获取第一个题目的ID（从字符串转换为int64）
	firstQuestionID := payload.Questions[0].ID

	s := ghttp.GetServer(fmt.Sprintf("quiz-duplicate-test-%d", time.Now().UnixNano()))
	s.BindMiddlewareDefault(func(r *ghttp.Request) {
		if uid := r.Header.Get("X-User-ID"); uid != "" {
			r.SetCtxVar("user_id", 1)
		}
		r.Middleware.Next()
	})
	s.Group("/api/v1", func(group *ghttp.RouterGroup) {
		group.POST("/quiz/submit", h.SubmitQuiz)
	})
	go s.Start()
	defer s.Shutdown()
	time.Sleep(30 * time.Millisecond)

	submitPayload := fmt.Sprintf(`{
		"sessionId": "%s",
		"topic": "constants",
		"chapter": "boolean",
		"answers": [{"questionId": %d, "userAnswers": ["A"]}],
		"durationMs": 5000
	}`, sessionID, firstQuestionID)

	// 第一次提交 - 应该成功
	req1 := httptest.NewRequest(http.MethodPost, "/api/v1/quiz/submit", strings.NewReader(submitPayload))
	req1.Header.Set("X-User-ID", "1")
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	s.ServeHTTP(w1, req1)

	if w1.Code != 200 {
		t.Fatalf("first submission failed: %d %s", w1.Code, w1.Body.String())
	}

	// 第二次提交 - 应该被拒绝
	req2 := httptest.NewRequest(http.MethodPost, "/api/v1/quiz/submit", strings.NewReader(submitPayload))
	req2.Header.Set("X-User-ID", "1")
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	s.ServeHTTP(w2, req2)

	// 应该返回错误(可能是409 Conflict或400 Bad Request)
	if w2.Code == 200 {
		t.Errorf("expected error on duplicate submission, got 200 success")
	}

	// 验证错误消息包含"重复"相关提示
	var resp map[string]interface{}
	if err := json.Unmarshal(w2.Body.Bytes(), &resp); err == nil {
		if msg, ok := resp["message"].(string); ok {
			if !strings.Contains(msg, "重复") && !strings.Contains(msg, "已提交") {
				t.Logf("Warning: error message doesn't mention duplicate: %s", msg)
			}
		}
	}
}
