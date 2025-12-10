package handler

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

	"go-study2/internal/app/http_server/middleware"
	"go-study2/internal/config"
	"go-study2/internal/infrastructure/database"
	appjwt "go-study2/internal/pkg/jwt"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/guid"
)

type quizResp struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

func TestQuizHandlers_Flow(t *testing.T) {
	ctx := gctx.New()
	_ = os.MkdirAll("testdata", 0o755)
	dbPath := filepath.ToSlash(filepath.Join("testdata", fmt.Sprintf("quiz_handler_%d.db", time.Now().UnixNano())))
	cfg := config.DatabaseConfig{
		Type: "sqlite3",
		Path: dbPath,
		Pragmas: []string{
			"journal_mode=WAL",
			"busy_timeout=5000",
			"synchronous=NORMAL",
			"cache_size=-64000",
			"foreign_keys=ON",
		},
	}
	if _, err := database.Init(ctx, cfg); err != nil {
		t.Fatalf("初始化数据库失败: %v", err)
	}
	if err := appjwt.Configure(appjwt.Options{
		Secret:             "abcdefabcdefabcdefabcdefabcdef12",
		AccessTokenExpiry:  time.Hour,
		RefreshTokenExpiry: 24 * time.Hour,
	}); err != nil {
		t.Fatalf("配置 JWT 失败: %v", err)
	}

	server := ghttp.GetServer(guid.S())
	server.SetPort(0)
	server.SetAccessLogEnabled(false)
	h := New()
	server.Group("/api/v1", func(group *ghttp.RouterGroup) {
		group.Middleware(middleware.Format)
		group.POST("/auth/register", h.Register)
		group.POST("/auth/login", h.Login)
		group.POST("/auth/refresh", h.RefreshToken)
		group.Group("/", func(authGroup *ghttp.RouterGroup) {
			authGroup.Middleware(middleware.Auth)
			authGroup.GET("/quiz/:topic/:chapter", h.GetQuiz)
			authGroup.POST("/quiz/submit", h.SubmitQuiz)
			authGroup.GET("/quiz/history", h.GetQuizHistory)
		})
	})
	server.Start()
	defer server.Shutdown()
	time.Sleep(60 * time.Millisecond)
	baseURL := fmt.Sprintf("http://127.0.0.1:%d", server.GetListenedPort())

	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	register := doQuizPost(t, client, baseURL+"/api/v1/auth/register", `{"username":"quiz_user","password":"TestPass123","remember":true}`)
	if register.Code != 20000 {
		t.Fatalf("注册失败: %s", register.Message)
	}
	var tokens map[string]interface{}
	_ = json.Unmarshal(register.Data, &tokens)
	access := fmt.Sprintf("%v", tokens["accessToken"])

	getQuizReq, _ := http.NewRequest(http.MethodGet, baseURL+"/api/v1/quiz/variables/storage", nil)
	getQuizReq.Header.Set("Authorization", "Bearer "+access)
	quizListResp := doQuizRequest(t, client, getQuizReq)
	if quizListResp.Code != 20000 {
		t.Fatalf("获取题目失败: %s", quizListResp.Message)
	}

	var questions []map[string]interface{}
	_ = json.Unmarshal(quizListResp.Data, &questions)
	if len(questions) == 0 {
		t.Fatalf("题目列表为空")
	}

	first := questions[0]
	qid, _ := first["id"].(string)
	ans := ""
	if arr, ok := first["answer"].([]interface{}); ok && len(arr) > 0 {
		if a, ok := arr[0].(string); ok {
			ans = a
		}
	}
	if qid == "" || ans == "" {
		t.Fatalf("题目信息缺失")
	}

	submitBody := fmt.Sprintf(`{"topic":"variables","chapter":"storage","answers":[{"id":"%s","choices":["%s"]}]}`, qid, ans)
	submitReq, _ := http.NewRequest(http.MethodPost, baseURL+"/api/v1/quiz/submit", bytes.NewBufferString(submitBody))
	submitReq.Header.Set("Content-Type", "application/json")
	submitReq.Header.Set("Authorization", "Bearer "+access)
	submitResp := doQuizRequest(t, client, submitReq)
	if submitResp.Code != 20000 {
		t.Fatalf("提交测验失败: %s", submitResp.Message)
	}

	historyReq, _ := http.NewRequest(http.MethodGet, baseURL+"/api/v1/quiz/history", nil)
	historyReq.Header.Set("Authorization", "Bearer "+access)
	historyResp := doQuizRequest(t, client, historyReq)
	if historyResp.Code != 20000 {
		t.Fatalf("查询历史失败: %s", historyResp.Message)
	}
}

func doQuizPost(t *testing.T, client *http.Client, url string, payload string) quizResp {
	t.Helper()
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	return doQuizRequest(t, client, req)
}

func doQuizRequest(t *testing.T, client *http.Client, req *http.Request) quizResp {
	t.Helper()
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	var body quizResp
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}
	return body
}

