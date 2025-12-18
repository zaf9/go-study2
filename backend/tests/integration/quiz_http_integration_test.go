package integration

import (
	"bytes"
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
	"reflect"
	"unsafe"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
)

// 简单集成测试：在真实数据库上插入题目，通过 http_server.NewServer 启动服务器并调用 /api/v1/quiz/:topic/:chapter
func TestQuizHTTP_StartSubmitHistory(t *testing.T) {
	ctx := gctx.New()
	// 初始化测试数据库（Temp in tests environment）
	// 使用临时文件作为 SQLite 数据库
	tmp := t.TempDir()
	dbPath := filepath.ToSlash(filepath.Join(tmp, "quiz_integration.db"))
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

	// 确保迁移已执行（Init 已经 run Migrate）

	now := time.Now()
	// 创建测试用户（外键约束需要 users 表存在该用户）
	if _, err := db.Exec(ctx, "DELETE FROM users"); err != nil {
		t.Fatalf("clear users failed: %v", err)
	}
	if _, err := db.Model("users").Data(map[string]interface{}{
		"id":            1,
		"username":      "testuser",
		"password_hash": "x",
		"created_at":    now,
		"updated_at":    now,
	}).Insert(); err != nil {
		t.Fatalf("insert test user failed: %v", err)
	}

	// 清理并插入自定义题目
	if _, err := db.Exec(ctx, "DELETE FROM quiz_questions"); err != nil {
		t.Fatalf("clear quiz_questions failed: %v", err)
	}

	seed := []map[string]interface{}{
		{
			"topic":           "variables",
			"chapter":         "storage",
			"type":            quizdom.QuestionTypeSingle,
			"difficulty":      quizdom.DifficultyEasy,
			"question":        "变量存储类型是？",
			"options":         `["栈","堆"]`,
			"correct_answers": `["A"]`,
			"explanation":     "示例",
			"created_at":      now,
			"updated_at":      now,
		},
	}
	for _, r := range seed {
		if _, err := db.Model("quiz_questions").Data(r).Insert(); err != nil {
			t.Fatalf("insert seed failed: %v", err)
		}
	}

	// 使用独立 ghttp.Server 并注入真实服务到 handler（避免全局中间件/认证影响）
	h := handler.New()
	// 构造 DB-backed repository 并创建服务
	repoImpl := infrarepo.NewQuizRepository(db)
	svc := appquiz.NewService(repoImpl)
	// 预先调用服务验证抽题流程是否可用（便于定位错误）
	if _, err := svc.GetQuizQuestions(gctx.New(), 1, "variables", "storage"); err != nil {
		t.Fatalf("service GetQuizQuestions failed: %v", err)
	}
	setQuizService(h, svc)

	s := ghttp.GetServer(fmt.Sprintf("quiz-integration-%d", time.Now().UnixNano()))
	// 绑定中间件注入 user_id（测试环境简单模拟）
	s.BindMiddlewareDefault(func(r *ghttp.Request) {
		if uid := r.Header.Get("X-User-ID"); uid != "" {
			r.SetCtxVar("user_id", 1)
		}
		r.Middleware.Next()
	})

	s.Group("/api/v1", func(group *ghttp.RouterGroup) {
		group.GET("/quiz/{topic}/{chapter}", h.GetQuiz)
		group.POST("/quiz/submit", h.SubmitQuiz)
		group.GET("/quiz/history", h.GetQuizHistory)
	})

	go s.Start()
	defer s.Shutdown()
	time.Sleep(30 * time.Millisecond)

	// GET /api/v1/quiz/variables/storage
	req := httptest.NewRequest(http.MethodGet, "/api/v1/quiz/variables/storage", nil)
	req.Header.Set("X-User-ID", "1")
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("expected 200 got %d body=%s", w.Code, w.Body.String())
	}

	// 解析返回并提交
	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal resp failed: %v body=%s", err, w.Body.String())
	}
	// data.sessionId 在通用响应结构里
	data, _ := resp["data"].(map[string]interface{})
	sid, _ := data["sessionId"].(string)
	if sid == "" {
		t.Fatalf("missing sessionId in response: %v", resp)
	}
	questions, _ := data["questions"].([]interface{})
	if len(questions) == 0 {
		t.Fatalf("no questions in response: %v", resp)
	}
	firstQuestion, _ := questions[0].(map[string]interface{})
	questionId, _ := firstQuestion["id"].(float64)

	// 提交答案
	body := fmt.Sprintf(`{"sessionId":"%s","topic":"variables","chapter":"storage","durationMs":5000,"answers":[{"questionId":%d,"userAnswers":["A"]}]}`, sid, int(questionId))
	req2 := httptest.NewRequest(http.MethodPost, "/api/v1/quiz/submit", bytes.NewBufferString(body))
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("X-User-ID", "1")
	w2 := httptest.NewRecorder()
	s.ServeHTTP(w2, req2)
	if w2.Code != 200 {
		t.Fatalf("submit expected 200 got %d body=%s", w2.Code, w2.Body.String())
	}

	// 历史
	req3 := httptest.NewRequest(http.MethodGet, "/api/v1/quiz/history", nil)
	req3.Header.Set("X-User-ID", "1")
	w3 := httptest.NewRecorder()
	s.ServeHTTP(w3, req3)
	if w3.Code != 200 {
		t.Fatalf("history expected 200 got %d body=%s", w3.Code, w3.Body.String())
	}
}

// setQuizService 使用反射注入测验服务，便于测试。
func setQuizService(h *handler.Handler, svc *appquiz.Service) {
	// 通过与 tests/interfaces 中相同的技巧注入私有字段
	// 为避免重复实现 reflection 细节，这里采用简单反射写法
	// 注意：此处假设 Handler 类型包含名为 quizService 的字段
	val := reflect.ValueOf(h).Elem().FieldByName("quizService")
	ptr := unsafe.Pointer(val.UnsafeAddr())
	reflect.NewAt(val.Type(), ptr).Elem().Set(reflect.ValueOf(svc))
}
