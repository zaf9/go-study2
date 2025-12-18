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
	"go-study2/internal/app/http_server/middleware"
	appquiz "go-study2/internal/app/quiz"
	"go-study2/internal/config"
	quizdom "go-study2/internal/domain/quiz"
	infrarepo "go-study2/internal/infra/repository"
	"go-study2/internal/infrastructure/database"
	appjwt "go-study2/internal/pkg/jwt"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
)

// T024: GET /quiz/history/{sessionId} API 契约测试
func TestQuizReviewContract(t *testing.T) {
	ctx := gctx.New()
	tmp := t.TempDir()
	dbPath := filepath.ToSlash(filepath.Join(tmp, "quiz_review_contract.db"))
	cfgDB := config.DatabaseConfig{
		Type:    "sqlite3",
		Path:    dbPath,
		Pragmas: []string{"journal_mode=WAL", "busy_timeout=5000", "synchronous=NORMAL", "foreign_keys=ON"},
	}
	db, err := database.Init(ctx, cfgDB)
	if err != nil {
		t.Fatalf("初始化数据库失败: %v", err)
	}
	defer db.Close(ctx)

	// 清理表
	if _, err := db.Exec(ctx, "DELETE FROM quiz_attempts"); err != nil {
		t.Fatalf("清理 quiz_attempts 失败: %v", err)
	}
	if _, err := db.Exec(ctx, "DELETE FROM quiz_sessions"); err != nil {
		t.Fatalf("清理 quiz_sessions 失败: %v", err)
	}
	if _, err := db.Exec(ctx, "DELETE FROM quiz_questions"); err != nil {
		t.Fatalf("清理 quiz_questions 失败: %v", err)
	}

	// 插入测试用户
	now := time.Now()
	if _, err := db.Model("users").Data(map[string]interface{}{
		"id": 1, "username": "testuser", "email": "test@example.com",
		"password_hash": "hash", "created_at": now, "updated_at": now,
	}).Insert(); err != nil {
		t.Fatalf("插入用户失败: %v", err)
	}

	// 插入测试题目
	questionSeed := map[string]interface{}{
		"topic":           "constants",
		"chapter":         "boolean",
		"type":            quizdom.QuestionTypeSingle,
		"difficulty":      quizdom.DifficultyEasy,
		"question":        "布尔常量的零值是？",
		"options":         `["false","true"]`,
		"correct_answers": `["A"]`,
		"explanation":     "布尔类型的零值是 false",
		"created_at":      now,
		"updated_at":      now,
	}
	result, err := db.Model("quiz_questions").Data(questionSeed).Insert()
	if err != nil {
		t.Fatalf("插入题目失败: %v", err)
	}
	questionID, _ := result.LastInsertId()

	// 插入测试会话
	completedAt := now.Add(10 * time.Minute)
	sessionID := "session-review-1"
	sessionData := map[string]interface{}{
		"session_id":      sessionID,
		"user_id":         1,
		"topic":           "constants",
		"chapter":         "boolean",
		"total_questions": 1,
		"correct_answers": 1,
		"score":           100,
		"passed":          true,
		"started_at":      now,
		"completed_at":    completedAt,
		"created_at":      now,
	}
	if _, err := db.Model("quiz_sessions").Data(sessionData).Insert(); err != nil {
		t.Fatalf("插入会话失败: %v", err)
	}

	// 插入测试答题记录
	attemptData := map[string]interface{}{
		"session_id":   sessionID,
		"user_id":      1,
		"topic":        "constants",
		"chapter":      "boolean",
		"question_id":  questionID,
		"user_answers": `["false"]`,
		"is_correct":   true,
		"attempted_at": now.Add(5 * time.Minute),
	}
	if _, err := db.Model("quiz_attempts").Data(attemptData).Insert(); err != nil {
		t.Fatalf("插入答题记录失败: %v", err)
	}

	// 插入 refresh_token 记录（用于 Auth 中间件验证）
	if _, err := db.Model("refresh_tokens").Data(map[string]interface{}{
		"user_id":    1,
		"token_hash": "test-hash",
		"expires_at": now.Add(24 * time.Hour),
		"created_at": now,
	}).Insert(); err != nil {
		t.Fatalf("插入 refresh_token 失败: %v", err)
	}

	// 配置 JWT
	if err := appjwt.Configure(appjwt.Options{
		Secret:             "abcdefabcdefabcdefabcdefabcdef12",
		AccessTokenExpiry:  time.Hour,
		RefreshTokenExpiry: 24 * time.Hour,
	}); err != nil {
		t.Fatalf("配置 JWT 失败: %v", err)
	}

	// 创建 handler 和服务
	h := handler.New()
	repoImpl := infrarepo.NewQuizRepository(db)
	svc := appquiz.NewService(repoImpl)
	setQuizService(h, svc)

	// 创建测试服务器
	s := ghttp.GetServer(fmt.Sprintf("quiz-review-contract-%d", time.Now().UnixNano()))
	s.BindMiddlewareDefault(func(r *ghttp.Request) {
		r.SetCtxVar("user_id", int64(1))
		r.Middleware.Next()
	})
	s.Group("/api/v1", func(group *ghttp.RouterGroup) {
		group.Middleware(middleware.Format)
		group.Group("/", func(authGroup *ghttp.RouterGroup) {
			authGroup.Middleware(middleware.Auth)
			authGroup.GET("/quiz/history/:sessionId", h.GetQuizReview)
		})
	})
	s.SetPort(0)
	go s.Start()
	defer s.Shutdown()
	time.Sleep(30 * time.Millisecond)

	// 生成 token
	token, err := appjwt.GenerateAccessToken(1)
	if err != nil {
		t.Fatalf("生成 token 失败: %v", err)
	}

	// 测试 GET /quiz/history/{sessionId}
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/quiz/history/%s", sessionID), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("期望状态码 200，但得到 %d，响应体: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("解析响应失败: %v，响应体: %s", err, w.Body.String())
	}

	// 验证响应结构
	if code, ok := resp["code"].(float64); !ok || int(code) != 20000 {
		t.Fatalf("响应 code 字段缺失或值不正确: %v", resp)
	}

	data, ok := resp["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("响应 data 字段缺失或类型不正确，期望对象: %v", resp)
	}

	// 验证 meta 字段
	meta, ok := data["meta"].(map[string]interface{})
	if !ok {
		t.Fatalf("data.meta 字段缺失或类型不正确: %v", data)
	}

	// 验证 meta 必需字段
	metaRequiredFields := []string{"sessionId", "topic", "chapter", "score", "passed", "completedAt"}
	for _, field := range metaRequiredFields {
		if _, ok := meta[field]; !ok {
			t.Fatalf("meta 缺少必需字段: %s，meta: %v", field, meta)
		}
	}

	// 验证 items 字段
	items, ok := data["items"].([]interface{})
	if !ok {
		t.Fatalf("data.items 字段缺失或类型不正确，期望数组: %v", data)
	}

	if len(items) == 0 {
		t.Fatalf("items 数组为空")
	}

	// 验证第一条题目的结构
	firstItem, ok := items[0].(map[string]interface{})
	if !ok {
		t.Fatalf("第一条题目类型不正确: %v", items[0])
	}

	// 验证题目必需字段
	itemRequiredFields := []string{"questionId", "stem", "options", "userChoice", "correctChoice", "isCorrect"}
	for _, field := range itemRequiredFields {
		if _, ok := firstItem[field]; !ok {
			t.Fatalf("题目缺少必需字段: %s，题目: %v", field, firstItem)
		}
	}

	// 验证字段类型
	if _, ok := firstItem["questionId"].(string); !ok {
		if _, ok2 := firstItem["questionId"].(float64); !ok2 {
			t.Fatalf("questionId 字段类型不正确，期望 string 或 number")
		}
	}
	if _, ok := firstItem["stem"].(string); !ok {
		t.Fatalf("stem 字段类型不正确，期望 string")
	}
	if _, ok := firstItem["options"].([]interface{}); !ok {
		t.Fatalf("options 字段类型不正确，期望数组")
	}
	if _, ok := firstItem["isCorrect"].(bool); !ok {
		t.Fatalf("isCorrect 字段类型不正确，期望 boolean")
	}
}
