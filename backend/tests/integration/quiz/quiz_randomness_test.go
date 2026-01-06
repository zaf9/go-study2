package quiz

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"os"
	"path/filepath"
	"testing"
	"time"

	"go-study2/internal/app/http_server/handler"
	middleware "go-study2/internal/app/http_server/middleware"
	appquiz "go-study2/internal/app/quiz"
	"go-study2/internal/config"
	quizdom "go-study2/internal/domain/quiz"
	infrarepo "go-study2/internal/infra/repository"
	"go-study2/internal/infrastructure/database"
	appjwt "go-study2/internal/pkg/jwt"
	"go-study2/internal/pkg/password"
	"io"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/guid"
)

type quizResp struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

func Test_QuizRandomness(t *testing.T) {
	ctx := gctx.New()
	_ = os.MkdirAll("testdata", 0o755)
	dbPath := filepath.ToSlash(filepath.Join("testdata", fmt.Sprintf("quiz_random_%d.db", time.Now().UnixNano())))
	cfg := config.DatabaseConfig{Type: "sqlite3", Path: dbPath}
	if _, err := database.Init(ctx, cfg); err != nil {
		t.Fatalf("初始化数据库失败: %v", err)
	}

	// 创建YAML仓储并添加足够的测试题目（120题，确保随机性）
	yamlRepo := quizdom.NewRepository()
	testYAMLQuestions := make([]quizdom.YAMLQuestion, 120)
	for i := 0; i < 120; i++ {
		diff := "easy"
		if i%3 == 1 {
			diff = "medium"
		} else if i%3 == 2 {
			diff = "hard"
		}
		typev := "single"
		if i%4 == 0 {
			typev = "multiple"
		}
		answer := "A"
		if typev == "multiple" {
			answer = "AB"
		}
		testYAMLQuestions[i] = quizdom.YAMLQuestion{
			ID:          fmt.Sprintf("rand%d", i+1),
			Type:        typev,
			Difficulty:  diff,
			Stem:        fmt.Sprintf("Seed question %d", i),
			Options:     []string{"A", "B", "C", "D"},
			Answer:      answer,
			Explanation: "seed",
			Topic:       "variables",
			Chapter:     "storage",
		}
	}
	yamlRepo.AddBank("variables", "storage", testYAMLQuestions)

	repoImpl := infrarepo.NewQuizRepository(database.Default())
	svc := appquiz.NewService(yamlRepo, repoImpl)

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

	// 注入quiz service到handler
	type quizSetter interface{ SetQuizService(*appquiz.Service) }
	if s, ok := interface{}(h).(quizSetter); ok {
		s.SetQuizService(svc)
	}

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
	client := &http.Client{Jar: jar}

	// create admin user
	adminPwd, _ := password.Hash("Admin123!")
	if _, err := database.Default().Insert(ctx, "users", map[string]interface{}{
		"username":      "admin_rand",
		"password_hash": adminPwd,
		"is_admin":      1,
		"status":        "active",
	}); err != nil {
		t.Fatalf("创建管理员失败: %v", err)
	}

	// login admin (sanity)
	adminLogin := doPost(t, client, baseURL+"/api/v1/auth/login", `{"username":"admin_rand","password":"Admin123!","remember":true}`)
	if adminLogin.Code != 20000 {
		t.Fatalf("管理员登录失败: %s", adminLogin.Message)
	}

	// create normal user directly in DB (bypass register endpoint)
	userPwd, _ := password.Hash("TestPass123!")
	if _, err := database.Default().Insert(ctx, "users", map[string]interface{}{
		"username":      "rand_user",
		"password_hash": userPwd,
		"is_admin":      0,
		"status":        "active",
	}); err != nil {
		t.Fatalf("创建普通用户失败: %v", err)
	}

	login := doPost(t, client, baseURL+"/api/v1/auth/login", `{"username":"rand_user","password":"TestPass123!","remember":true}`)
	if login.Code != 20000 {
		t.Fatalf("登录失败: %s", login.Message)
	}
	var tokens map[string]interface{}
	_ = json.Unmarshal(login.Data, &tokens)
	access := fmt.Sprintf("%v", tokens["accessToken"])

	// call quiz endpoint multiple times and collect first-question ids
	distinct := map[string]struct{}{}
	iterations := 100
	for i := 0; i < iterations; i++ {
		req, _ := http.NewRequest(http.MethodGet, baseURL+"/api/v1/quiz/variables/storage", nil)
		req.Header.Set("Authorization", "Bearer "+access)
		resp := doRequest(t, client, req)
		if resp.Code != 20000 {
			t.Fatalf("获取题目失败: %s", resp.Message)
		}
		var data struct {
			SessionID string                   `json:"sessionId"`
			Questions []map[string]interface{} `json:"questions"`
		}
		_ = json.Unmarshal(resp.Data, &data)
		if len(data.Questions) == 0 {
			t.Fatalf("题目为空")
		}
		first := data.Questions[0]
		id := fmt.Sprintf("%v", first["id"])
		distinct[id] = struct{}{}
	}
	if len(distinct) < iterations/2 {
		t.Fatalf("随机性不足: %d 个不同的首题, 需要至少 %d", len(distinct), iterations/2)
	}
}

// helpers
func bytesReader(s string) *bytes.Reader { return bytes.NewReader([]byte(s)) }

func doPost(t *testing.T, client *http.Client, url, payload string) quizResp {
	t.Helper()
	req, _ := http.NewRequest(http.MethodPost, url, bytesReader(payload))
	req.Header.Set("Content-Type", "application/json")
	return doRequest(t, client, req)
}

func doRequest(t *testing.T, client *http.Client, req *http.Request) quizResp {
	t.Helper()
	t.Logf("Sending HTTP %s %s", req.Method, req.URL.String())
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()
	var body quizResp
	data, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(data, &body); err != nil {
		t.Logf("HTTP %s %s -> Status=%s", req.Method, req.URL.String(), resp.Status)
		t.Fatalf("解析响应失败: %v; body=%s", err, string(data))
	}
	return body
}
