package quiz

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"go-study2/internal/app/http_server/handler"
	middleware "go-study2/internal/app/http_server/middleware"
	"go-study2/internal/config"
	"go-study2/internal/infrastructure/database"
	appjwt "go-study2/internal/pkg/jwt"
	"go-study2/internal/pkg/password"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/guid"
)

func Test_QuizConcurrency(t *testing.T) {
	ctx := gctx.New()
	_ = os.MkdirAll("testdata", 0o755)
	dbPath := filepath.ToSlash(filepath.Join("testdata", fmt.Sprintf("quiz_conc_%d.db", time.Now().UnixNano())))
	cfg := config.DatabaseConfig{Type: "sqlite3", Path: dbPath}
	if _, err := database.Init(ctx, cfg); err != nil {
		t.Fatalf("初始化数据库失败: %v", err)
	}

	// seed quiz questions for concurrency test
	now := time.Now()
	for i := 0; i < 30; i++ {
		q := map[string]interface{}{
			"topic":           "variables",
			"chapter":         "storage",
			"type":            "single",
			"difficulty":      "easy",
			"question":        fmt.Sprintf("Conc seed %d", i),
			"options":         `["A","B","C","D"]`,
			"correct_answers": `[["A"]]`,
			"explanation":     "seed",
			"created_at":      now,
			"updated_at":      now,
		}
		if _, err := database.Default().Insert(ctx, "quiz_questions", q); err != nil {
			t.Fatalf("插入并发种子题目失败: %v", err)
		}
	}

	// Configure JWT for tests
	if err := appjwt.Configure(appjwt.Options{
		Secret:             "testsecret0123456789012345678901",
		AccessTokenExpiry:  time.Hour,
		RefreshTokenExpiry: 24 * time.Hour,
	}); err != nil {
		t.Fatalf("配置 JWT 失败: %v", err)
	}

	server := ghttp.GetServer(guid.S())
	server.SetPort(0)
	server.SetAccessLogEnabled(false)
	h := handler.New()
	server.Group("/api/v1", func(group *ghttp.RouterGroup) {
		group.POST("/auth/login", h.Login)
		group.POST("/auth/register", middleware.Auth, h.Register)
		group.Group("/", func(authGroup *ghttp.RouterGroup) {
			authGroup.Middleware(middleware.Auth)
			authGroup.GET("/quiz/:topic/:chapter", h.GetQuiz)
		})
	})
	server.Start()
	defer server.Shutdown()
	time.Sleep(20 * time.Millisecond)
	baseURL := fmt.Sprintf("http://127.0.0.1:%d", server.GetListenedPort())

	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar, Timeout: 5 * time.Second}

	adminPwd, _ := password.Hash("Admin123!")
	if _, err := database.Default().Insert(ctx, "users", map[string]interface{}{
		"username":      "admin_conc",
		"password_hash": adminPwd,
		"is_admin":      1,
		"status":        "active",
	}); err != nil {
		t.Fatalf("创建管理员失败: %v", err)
	}
	adminLogin := doPost(t, client, baseURL+"/api/v1/auth/login", `{"username":"admin_conc","password":"Admin123!","remember":true}`)
	if adminLogin.Code != 20000 {
		t.Fatalf("管理员登录失败: %s", adminLogin.Message)
	}

	// create normal user directly in DB (bypass register endpoint)
	userPwd, _ := password.Hash("TestPass123!")
	if _, err := database.Default().Insert(ctx, "users", map[string]interface{}{
		"username":      "conc_user",
		"password_hash": userPwd,
		"is_admin":      0,
		"status":        "active",
	}); err != nil {
		t.Fatalf("创建普通用户失败: %v", err)
	}

	login := doPost(t, client, baseURL+"/api/v1/auth/login", `{"username":"conc_user","password":"TestPass123!","remember":true}`)
	if login.Code != 20000 {
		t.Fatalf("登录失败: %s", login.Message)
	}
	var tokens map[string]interface{}
	_ = json.Unmarshal(login.Data, &tokens)
	access := fmt.Sprintf("%v", tokens["accessToken"])

	concurrency := 100
	wg := sync.WaitGroup{}
	wg.Add(concurrency)
	errCh := make(chan error, concurrency)
	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			req, _ := http.NewRequest(http.MethodGet, baseURL+"/api/v1/quiz/variables/storage", nil)
			req.Header.Set("Authorization", "Bearer "+access)
			resp := doRequest(t, client, req)
			if resp.Code != 20000 {
				errCh <- fmt.Errorf("响应码错误: %d", resp.Code)
				return
			}
		}()
	}
	wg.Wait()
	close(errCh)
	if len(errCh) != 0 {
		t.Fatalf("并发请求中有错误: %v", <-errCh)
	}
}

// small helpers
// bytesReader helper is defined in quiz_randomness_test.go and reused by other tests in this package.
