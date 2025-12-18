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
	infrarepo "go-study2/internal/infra/repository"
	"go-study2/internal/infrastructure/database"
	appjwt "go-study2/internal/pkg/jwt"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
)

// T023: GET /quiz/history API 契约测试
func TestQuizHistoryContract(t *testing.T) {
	ctx := gctx.New()
	tmp := t.TempDir()
	dbPath := filepath.ToSlash(filepath.Join(tmp, "quiz_history_contract.db"))
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
	if _, err := db.Exec(ctx, "DELETE FROM quiz_sessions"); err != nil {
		t.Fatalf("清理 quiz_sessions 失败: %v", err)
	}

	// 插入测试用户
	now := time.Now()
	if _, err := db.Model("users").Data(map[string]interface{}{
		"id": 1, "username": "testuser", "email": "test@example.com",
		"password_hash": "hash", "created_at": now, "updated_at": now,
	}).Insert(); err != nil {
		t.Fatalf("插入用户失败: %v", err)
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

	// 插入测试会话数据
	completedAt := now.Add(10 * time.Minute)
	sessions := []map[string]interface{}{
		{
			"session_id":      "session-1",
			"user_id":         1,
			"topic":           "constants",
			"chapter":         "boolean",
			"total_questions": 5,
			"correct_answers": 4,
			"score":           80,
			"passed":          true,
			"started_at":      now,
			"completed_at":    completedAt,
			"created_at":      now,
		},
		{
			"session_id":      "session-2",
			"user_id":         1,
			"topic":           "constants",
			"chapter":         "boolean",
			"total_questions": 5,
			"correct_answers": 2,
			"score":           40,
			"passed":          false,
			"started_at":      now.Add(-1 * time.Hour),
			"completed_at":    completedAt.Add(-1 * time.Hour),
			"created_at":      now.Add(-1 * time.Hour),
		},
	}
	for _, sess := range sessions {
		if _, err := db.Model("quiz_sessions").Data(sess).Insert(); err != nil {
			t.Fatalf("插入会话失败: %v", err)
		}
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
	s := ghttp.GetServer(fmt.Sprintf("quiz-history-contract-%d", time.Now().UnixNano()))
	s.BindMiddlewareDefault(func(r *ghttp.Request) {
		r.SetCtxVar("user_id", int64(1))
		r.Middleware.Next()
	})
	s.Group("/api/v1", func(group *ghttp.RouterGroup) {
		group.Middleware(middleware.Format)
		group.Group("/", func(authGroup *ghttp.RouterGroup) {
			authGroup.Middleware(middleware.Auth)
			authGroup.GET("/quiz/history", h.GetQuizHistory)
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

	// 测试 GET /quiz/history
	req := httptest.NewRequest(http.MethodGet, "/api/v1/quiz/history", nil)
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

	data, ok := resp["data"].([]interface{})
	if !ok {
		t.Fatalf("响应 data 字段缺失或类型不正确，期望数组: %v", resp)
	}

	if len(data) == 0 {
		t.Fatalf("历史记录列表为空")
	}

	// 验证第一条记录的结构
	firstItem, ok := data[0].(map[string]interface{})
	if !ok {
		t.Fatalf("第一条记录类型不正确: %v", data[0])
	}

	// 验证必需字段
	requiredFields := []string{"sessionId", "topic", "chapter", "score", "passed", "completedAt"}
	for _, field := range requiredFields {
		if _, ok := firstItem[field]; !ok {
			t.Fatalf("记录缺少必需字段: %s，记录: %v", field, firstItem)
		}
	}

	// 验证字段类型
	if _, ok := firstItem["sessionId"].(string); !ok {
		t.Fatalf("sessionId 字段类型不正确，期望 string")
	}
	if _, ok := firstItem["score"].(float64); !ok {
		t.Fatalf("score 字段类型不正确，期望 number")
	}
	if _, ok := firstItem["passed"].(bool); !ok {
		t.Fatalf("passed 字段类型不正确，期望 boolean")
	}
}
