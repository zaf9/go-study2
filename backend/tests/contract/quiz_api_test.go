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

	// 清空并初始化测试数据
	if _, err := db.Exec(ctx, "DELETE FROM quiz_questions"); err != nil {
		t.Fatalf("clear quiz_questions failed: %v", err)
	}
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

	// 插入测试题目
	question := map[string]interface{}{
		"topic":           "constants",
		"chapter":         "boolean",
		"type":            quizdom.QuestionTypeSingle,
		"difficulty":      quizdom.DifficultyEasy,
		"question":        "Go语言中布尔类型的零值是什么?",
		"options":         `["true","false","nil","0"]`,
		"correct_answers": `["B"]`,
		"explanation":     "布尔类型的零值是false",
		"created_at":      now,
		"updated_at":      now,
	}
	if _, err := db.Model("quiz_questions").Data(question).Insert(); err != nil {
		t.Fatalf("insert question failed: %v", err)
	}

	// 创建handler和service
	h := handler.New()
	repoImpl := infrarepo.NewQuizRepository(db)
	svc := appquiz.NewService(repoImpl)

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
	for _, table := range []string{"quiz_questions", "quiz_sessions", "quiz_attempts", "users"} {
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

	// 插入测试题目 - 正确答案为B
	question := map[string]interface{}{
		"topic":           "constants",
		"chapter":         "boolean",
		"type":            quizdom.QuestionTypeSingle,
		"difficulty":      quizdom.DifficultyEasy,
		"question":        "Go语言中布尔类型的零值是什么?",
		"options":         `["true","false"]`,
		"correct_answers": `["B"]`,
		"explanation":     "布尔类型的零值是false",
		"created_at":      now,
		"updated_at":      now,
	}
	result, err := db.Model("quiz_questions").Data(question).Insert()
	if err != nil {
		t.Fatalf("insert question failed: %v", err)
	}
	questionID, _ := result.LastInsertId()

	// 创建handler和service
	h := handler.New()
	repoImpl := infrarepo.NewQuizRepository(db)
	svc := appquiz.NewService(repoImpl)

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

	// 提交正确答案
	submitPayload := fmt.Sprintf(`{
		"sessionId": "%s",
		"topic": "constants",
		"chapter": "boolean",
		"answers": [{"questionId": %d, "userAnswers": ["B"]}],
		"durationMs": 5000
	}`, sessionID, questionID)

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
	for _, table := range []string{"quiz_questions", "quiz_sessions", "quiz_attempts", "users"} {
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

	question := map[string]interface{}{
		"topic":           "constants",
		"chapter":         "boolean",
		"type":            quizdom.QuestionTypeSingle,
		"difficulty":      quizdom.DifficultyEasy,
		"question":        "测试题目",
		"options":         `["选项A"]`,
		"correct_answers": `["A"]`,
		"explanation":     "解析",
		"created_at":      now,
		"updated_at":      now,
	}
	result, err := db.Model("quiz_questions").Data(question).Insert()
	if err != nil {
		t.Fatalf("insert question failed: %v", err)
	}
	questionID, _ := result.LastInsertId()

	h := handler.New()
	repoImpl := infrarepo.NewQuizRepository(db)
	svc := appquiz.NewService(repoImpl)

	type quizSetter interface{ SetQuizService(*appquiz.Service) }
	if s, ok := interface{}(h).(quizSetter); ok {
		s.SetQuizService(svc)
	}

	payload, err := svc.GetQuizQuestions(ctx, 1, "constants", "boolean")
	if err != nil {
		t.Fatalf("GetQuizQuestions failed: %v", err)
	}
	sessionID := payload.SessionID

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
	}`, sessionID, questionID)

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
