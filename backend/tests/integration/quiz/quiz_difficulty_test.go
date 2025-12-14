package quiz

import (
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
	"go-study2/internal/config"
	"go-study2/internal/infrastructure/database"
	appjwt "go-study2/internal/pkg/jwt"
	"go-study2/internal/pkg/password"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/guid"
)

func Test_QuizDifficultyDistribution(t *testing.T) {
	ctx := gctx.New()
	_ = os.MkdirAll("testdata", 0o755)
	dbPath := filepath.ToSlash(filepath.Join("testdata", fmt.Sprintf("quiz_diff_%d.db", time.Now().UnixNano())))
	cfg := config.DatabaseConfig{Type: "sqlite3", Path: dbPath}
	if _, err := database.Init(ctx, cfg); err != nil {
		t.Fatalf("初始化数据库失败: %v", err)
	}

	// seed quiz questions for difficulty distribution test
	now := time.Now()
	for i := 0; i < 300; i++ {
		d := "easy"
		if i%5 == 0 {
			d = "hard"
		} else if i%2 == 0 {
			d = "medium"
		}
		q := map[string]interface{}{
			"topic":           "variables",
			"chapter":         "storage",
			"type":            "single",
			"difficulty":      d,
			"question":        fmt.Sprintf("Diff seed %d", i),
			"options":         `["A","B","C","D"]`,
			"correct_answers": `[["A"]]`,
			"explanation":     "seed",
			"created_at":      now,
			"updated_at":      now,
		}
		if _, err := database.Default().Insert(ctx, "quiz_questions", q); err != nil {
			t.Fatalf("插入难度种子题目失败: %v", err)
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
	client := &http.Client{Jar: jar}

	// insert admin and register user
	adminPwd, _ := password.Hash("Admin123!")
	if _, err := database.Default().Insert(ctx, "users", map[string]interface{}{
		"username":      "admin_diff",
		"password_hash": adminPwd,
		"is_admin":      1,
		"status":        "active",
	}); err != nil {
		t.Fatalf("创建管理员失败: %v", err)
	}
	adminLogin := doPost(t, client, baseURL+"/api/v1/auth/login", `{"username":"admin_diff","password":"Admin123!","remember":true}`)
	if adminLogin.Code != 20000 {
		t.Fatalf("管理员登录失败: %s", adminLogin.Message)
	}

	// create normal user directly in DB (bypass register endpoint)
	userPwd, _ := password.Hash("TestPass123!")
	if _, err := database.Default().Insert(ctx, "users", map[string]interface{}{
		"username":      "diff_user",
		"password_hash": userPwd,
		"is_admin":      0,
		"status":        "active",
	}); err != nil {
		t.Fatalf("创建普通用户失败: %v", err)
	}

	login := doPost(t, client, baseURL+"/api/v1/auth/login", `{"username":"diff_user","password":"TestPass123!","remember":true}`)
	if login.Code != 20000 {
		t.Fatalf("登录失败: %s", login.Message)
	}
	var tokens map[string]interface{}
	_ = json.Unmarshal(login.Data, &tokens)
	access := fmt.Sprintf("%v", tokens["accessToken"])

	iterations := 200
	counts := map[string]int{"easy": 0, "medium": 0, "hard": 0}
	totalQuestions := 0
	for i := 0; i < iterations; i++ {
		req, _ := http.NewRequest(http.MethodGet, baseURL+"/api/v1/quiz/variables/storage", nil)
		req.Header.Set("Authorization", "Bearer "+access)
		resp := doRequest(t, client, req)
		if resp.Code != 20000 {
			t.Fatalf("获取题目失败: %s", resp.Message)
		}
		var data struct {
			Questions []map[string]interface{} `json:"questions"`
		}
		_ = json.Unmarshal(resp.Data, &data)
		for _, q := range data.Questions {
			if d, ok := q["difficulty"].(string); ok {
				switch d {
				case "easy":
					counts["easy"]++
				case "medium":
					counts["medium"]++
				case "hard":
					counts["hard"]++
				}
			}
			totalQuestions++
		}
	}
	if totalQuestions == 0 {
		t.Fatalf("没有收集到题目")
	}
	easyPct := float64(counts["easy"]) / float64(totalQuestions)
	medPct := float64(counts["medium"]) / float64(totalQuestions)
	hardPct := float64(counts["hard"]) / float64(totalQuestions)

	if !(easyPct >= 0.30 && easyPct <= 0.50) {
		t.Fatalf("easy 比例超出范围: %.2f", easyPct)
	}
	if !(medPct >= 0.30 && medPct <= 0.50) {
		t.Fatalf("medium 比例超出范围: %.2f", medPct)
	}
	if !(hardPct >= 0.10 && hardPct <= 0.30) {
		t.Fatalf("hard 比例超出范围: %.2f", hardPct)
	}
}
