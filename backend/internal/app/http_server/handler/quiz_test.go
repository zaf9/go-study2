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
	quizdom "go-study2/internal/domain/quiz"
	"go-study2/internal/infrastructure/database"
	appjwt "go-study2/internal/pkg/jwt"
	"go-study2/internal/pkg/password"

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

	// 初始化 YAML 仓库并添加测试数据
	// 注意：使用全局仓库单例
	yamlRepo := quizdom.GetGlobalRepository()
	testQuestions := []quizdom.YAMLQuestion{
		{
			ID:          "1",
			Type:        "single",
			Difficulty:  "easy",
			Stem:        "What is the default value of int?",
			Options:     []string{"0", "1", "nil", "undefined"},
			Answer:      "A",
			Explanation: "The default value is 0",
			Topic:       "variables",
			Chapter:     "storage",
		},
		{
			ID:          "2",
			Type:        "single",
			Difficulty:  "easy",
			Stem:        "Which keyword declares a variable?",
			Options:     []string{"var", "let", "const", "val"},
			Answer:      "A",
			Explanation: "Go uses var to declare variables",
			Topic:       "variables",
			Chapter:     "storage",
		},
		{
			ID:          "3",
			Type:        "single",
			Difficulty:  "easy",
			Stem:        "Which is a value type?",
			Options:     []string{"int", "slice", "map", "channel"},
			Answer:      "A",
			Explanation: "int is a value type",
			Topic:       "variables",
			Chapter:     "storage",
		},
		{
			ID:          "4",
			Type:        "single",
			Difficulty:  "easy",
			Stem:        "What is zero value for string?",
			Options:     []string{`""`, "nil", "0", "false"},
			Answer:      "A",
			Explanation: "Empty string is zero value",
			Topic:       "variables",
			Chapter:     "storage",
		},
		{
			ID:          "5",
			Type:        "multiple",
			Difficulty:  "medium",
			Stem:        "Which are value types?",
			Options:     []string{"int", "string", "slice", "array"},
			Answer:      "A,B,D",
			Explanation: "int, string, and array are value types",
			Topic:       "variables",
			Chapter:     "storage",
		},
		{
			ID:          "6",
			Type:        "multiple",
			Difficulty:  "medium",
			Stem:        "Which keywords declare variables?",
			Options:     []string{"var", "const", ":=", "let"},
			Answer:      "A,B,C",
			Explanation: "var, const, and := are used in Go",
			Topic:       "variables",
			Chapter:     "storage",
		},
		{
			ID:          "7",
			Type:        "multiple",
			Difficulty:  "medium",
			Stem:        "Which are reference types?",
			Options:     []string{"slice", "map", "channel", "int"},
			Answer:      "A,B,C",
			Explanation: "slice, map, and channel are reference types",
			Topic:       "variables",
			Chapter:     "storage",
		},
		{
			ID:          "8",
			Type:        "multiple",
			Difficulty:  "medium",
			Stem:        "Which can store multiple values?",
			Options:     []string{"array", "slice", "map", "int"},
			Answer:      "A,B,C",
			Explanation: "array, slice, and map store multiple values",
			Topic:       "variables",
			Chapter:     "storage",
		},
	}
	yamlRepo.AddBank("variables", "storage", testQuestions)

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
		group.POST("/auth/login", h.Login)
		group.POST("/auth/refresh", h.RefreshToken)
		group.Group("/", func(authGroup *ghttp.RouterGroup) {
			authGroup.Middleware(middleware.Auth)
			authGroup.POST("/auth/register", h.Register)
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

	adminPwd, _ := password.Hash("Admin123!")
	if _, err := database.Default().Insert(ctx, "users", map[string]interface{}{
		"username":             "admin",
		"password_hash":        adminPwd,
		"is_admin":             1,
		"status":               "active",
		"must_change_password": 0,
	}); err != nil {
		t.Fatalf("创建管理员失败: %v", err)
	}

	adminLogin := doQuizPost(t, client, baseURL+"/api/v1/auth/login", `{"username":"admin","password":"Admin123!","remember":true}`)
	if adminLogin.Code != 20000 {
		t.Fatalf("管理员登录失败: %s", adminLogin.Message)
	}
	var adminTokens map[string]interface{}
	_ = json.Unmarshal(adminLogin.Data, &adminTokens)
	adminAccess := fmt.Sprintf("%v", adminTokens["accessToken"])

	registerReq, _ := http.NewRequest(http.MethodPost, baseURL+"/api/v1/auth/register", bytes.NewBufferString(`{"username":"quiz_user","password":"TestPass123!","remember":true}`))
	registerReq.Header.Set("Content-Type", "application/json")
	registerReq.Header.Set("Authorization", "Bearer "+adminAccess)
	register := doQuizRequest(t, client, registerReq)
	if register.Code != 20000 {
		t.Fatalf("注册失败: %s", register.Message)
	}

	login := doQuizPost(t, client, baseURL+"/api/v1/auth/login", `{"username":"quiz_user","password":"TestPass123!","remember":true}`)
	if login.Code != 20000 {
		t.Fatalf("登录失败: %s", login.Message)
	}
	var tokens map[string]interface{}
	_ = json.Unmarshal(login.Data, &tokens)
	access := fmt.Sprintf("%v", tokens["accessToken"])

	getQuizReq, _ := http.NewRequest(http.MethodGet, baseURL+"/api/v1/quiz/variables/storage", nil)
	getQuizReq.Header.Set("Authorization", "Bearer "+access)
	quizListResp := doQuizRequest(t, client, getQuizReq)
	if quizListResp.Code != 20000 {
		t.Fatalf("获取题目失败: %s", quizListResp.Message)
	}

	var quizData struct {
		SessionID string                   `json:"sessionId"`
		Questions []map[string]interface{} `json:"questions"`
	}
	_ = json.Unmarshal(quizListResp.Data, &quizData)
	if quizData.SessionID == "" {
		t.Fatalf("题目响应缺少 sessionId")
	}
	if len(quizData.Questions) == 0 {
		t.Fatalf("题目列表为空")
	}

	first := quizData.Questions[0]
	// 题目 ID 以字符串形式返回（避免 JavaScript 大整数精度问题）
	idStr, ok := first["id"].(string)
	if !ok || idStr == "" {
		t.Fatalf("题目 ID 缺失或格式错误: %+v", first)
	}
	// 提取数字部分（从 "1" 转换为 1）
	var qid int64
	if _, err := fmt.Sscanf(idStr, "%d", &qid); err != nil || qid == 0 {
		t.Fatalf("题目 ID 无法解析为整数: %s", idStr)
	}
	if qid == 0 {
		t.Fatalf("题目信息缺失")
	}

	submitBody := fmt.Sprintf(`{"sessionId":"%s","topic":"variables","chapter":"storage","durationMs":5000,"answers":[{"questionId":%d,"userAnswers":["A"]}]}`, quizData.SessionID, qid)
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
