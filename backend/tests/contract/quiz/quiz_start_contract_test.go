package contract

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
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

// T022: 抽题API契约测试 - 验证返回结构满足 contracts/api-spec.md
func TestQuizStartContract_Run(t *testing.T) {
	ctx := gctx.New()
	tmp := t.TempDir()
	dbPath := filepath.ToSlash(filepath.Join(tmp, "quiz_contract.db"))
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

	// Clear tables
	if _, err := db.Exec(ctx, "DELETE FROM quiz_questions"); err != nil {
		t.Fatalf("clear quiz_questions failed: %v", err)
	}
	if _, err := db.Exec(ctx, "DELETE FROM quiz_sessions"); err != nil {
		t.Fatalf("clear quiz_sessions failed: %v", err)
	}
	now := time.Now()
	// Insert a test user
	if _, err := db.Model("users").Data(map[string]interface{}{
		"id": 1, "username": "testuser", "email": "test@example.com", "password_hash": "hash", "created_at": now, "updated_at": now,
	}).Insert(); err != nil {
		t.Fatalf("insert user failed: %v", err)
	}
	// Clear and seed a quiz question
	seed := map[string]interface{}{
		"topic":           "constants",
		"chapter":         "boolean",
		"type":            quizdom.QuestionTypeSingle,
		"difficulty":      quizdom.DifficultyEasy,
		"question":        "布尔常量有哪些？",
		"options":         `["true","false"]`,
		"correct_answers": `["A"]`,
		"explanation":     "示例解析",
		"created_at":      now,
		"updated_at":      now,
	}
	if _, err := db.Model("quiz_questions").Data(seed).Insert(); err != nil {
		t.Fatalf("insert seed failed: %v", err)
	}

	// Start handler and inject service
	h := handler.New()
	repoImpl := infrarepo.NewQuizRepository(db)
	svc := appquiz.NewService(repoImpl)
	// sanity call
	if _, err := svc.GetQuizQuestions(gctx.New(), 1, "constants", "boolean"); err != nil {
		t.Fatalf("service GetQuizQuestions failed: %v", err)
	}

	// inject service if handler supports it
	type quizSetter interface{ SetQuizService(*appquiz.Service) }
	if s, ok := interface{}(h).(quizSetter); ok {
		s.SetQuizService(svc)
	}

	s := ghttp.GetServer(fmt.Sprintf("quiz-contract-%d", time.Now().UnixNano()))
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

	req := httptest.NewRequest(http.MethodGet, "/api/v1/quiz/constants/boolean", nil)
	req.Header.Set("X-User-ID", "1")
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("expected 200 got %d body=%s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal resp failed: %v body=%s", err, w.Body.String())
	}

	// Basic contract assertions
	if _, ok := resp["code"]; !ok {
		t.Fatalf("response missing 'code' field: %v", resp)
	}
	if _, ok := resp["message"]; !ok {
		t.Fatalf("response missing 'message' field: %v", resp)
	}
	data, ok := resp["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("response 'data' missing or wrong type: %v", resp)
	}
	questions, ok := data["questions"].([]interface{})
	if !ok {
		t.Fatalf("data.questions missing or wrong type: %v", data)
	}
	if len(questions) == 0 {
		t.Fatalf("questions array is empty: %v", data)
	}
	// check fields of first question
	q0, _ := questions[0].(map[string]interface{})
	// Accept either 'stem' or 'question' for question text
	if _, ok := q0["stem"]; !ok {
		if _, ok2 := q0["question"]; !ok2 {
			t.Fatalf("question missing 'stem' or 'question' field: %v", q0)
		}
	}
	// Required basic fields
	basic := []string{"id", "type", "difficulty", "options"}
	for _, f := range basic {
		if _, ok := q0[f]; !ok {
			t.Fatalf("question missing field %s: %v", f, q0)
		}
	}
	// options can be []string or []object with label
	switch opts := q0["options"].(type) {
	case []interface{}:
		if len(opts) == 0 {
			t.Fatalf("options empty: %v", opts)
		}
	default:
		t.Fatalf("options has unexpected type: %T", opts)
	}
}
