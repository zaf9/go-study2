package integration

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

	"go-study2/internal/app/http_server"
	"go-study2/internal/config"
	quizdom "go-study2/internal/domain/quiz"
	"go-study2/internal/domain/user"
	"go-study2/internal/infrastructure/database"
	appjwt "go-study2/internal/pkg/jwt"

	"github.com/gogf/gf/v2/os/gctx"
)

func TestQuizFlow_EndToEnd(t *testing.T) {
	ctx := gctx.New()
	_ = os.MkdirAll("testdata", 0o755)
	dbPath := filepath.ToSlash(filepath.Join("testdata", fmt.Sprintf("integration_quiz_%d.db", time.Now().UnixNano())))
	cfg := &config.Config{
		Server: config.ServerConfig{Host: "127.0.0.1"},
		Http:   config.HttpConfig{Port: 0},
		Database: config.DatabaseConfig{
			Type: "sqlite3",
			Path: dbPath,
			Pragmas: []string{
				"journal_mode=WAL",
				"busy_timeout=5000",
				"synchronous=NORMAL",
				"cache_size=-64000",
				"foreign_keys=ON",
			},
		},
		Jwt: config.JwtConfig{
			Secret:             "integration-secret-quiz-abcdef123456",
			AccessTokenExpiry:  3600,
			RefreshTokenExpiry: 604800,
			Issuer:             "go-study2",
		},
	}

	if _, err := database.Init(ctx, cfg.Database); err != nil {
		t.Fatalf("初始化数据库失败: %v", err)
	}
	if err := appjwt.Configure(appjwt.Options{
		Secret:             cfg.Jwt.Secret,
		Issuer:             cfg.Jwt.Issuer,
		AccessTokenExpiry:  time.Duration(cfg.Jwt.AccessTokenExpiry) * time.Second,
		RefreshTokenExpiry: time.Duration(cfg.Jwt.RefreshTokenExpiry) * time.Second,
	}); err != nil {
		t.Fatalf("配置 JWT 失败: %v", err)
	}

	// 初始化全局YAML题库（使用项目的quiz_data目录）
	quizDataPath := filepath.Join("..", "..", "quiz_data")
	if _, err := os.Stat(quizDataPath); os.IsNotExist(err) {
		// 如果默认路径不存在，尝试绝对路径
		quizDataPath = filepath.Join("backend", "quiz_data")
	}
	if err := quizdom.InitializeGlobalRepository(quizDataPath); err != nil {
		t.Fatalf("初始化题库失败: %v", err)
	}
	t.Logf("题库已从 %s 加载", quizDataPath)

	server, err := http_server.NewServer(cfg, "quiz-integration")
	if err != nil {
		t.Fatalf("创建服务器失败: %v", err)
	}
	server.SetPort(0)
	server.SetAccessLogEnabled(false)
	server.Start()
	defer server.Shutdown()
	time.Sleep(80 * time.Millisecond)

	baseURL := fmt.Sprintf("http://127.0.0.1:%d", server.GetListenedPort())
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	adminLogin := doIntegrationPost(t, client, baseURL+"/api/v1/auth/login", fmt.Sprintf(`{"username":"%s","password":"%s","rememberMe":true}`, user.DefaultAdminUsername, user.DefaultAdminPassword))
	if adminLogin.Code != 20000 {
		t.Fatalf("管理员登录失败: %s", adminLogin.Message)
	}
	var adminData struct {
		AccessToken        string `json:"accessToken"`
		NeedPasswordChange bool   `json:"needPasswordChange"`
	}
	_ = json.Unmarshal(adminLogin.Data, &adminData)

	adminPassword := user.DefaultAdminPassword
	if adminData.NeedPasswordChange {
		changeReqBody := fmt.Sprintf(`{"oldPassword":"%s","newPassword":"NewQuizPass123!"}`, adminPassword)
		changeReq, _ := http.NewRequest(http.MethodPost, baseURL+"/api/v1/auth/change-password", bytes.NewBufferString(changeReqBody))
		changeReq.Header.Set("Content-Type", "application/json")
		changeReq.Header.Set("Authorization", "Bearer "+adminData.AccessToken)
		changeResp := doIntegrationRequest(t, client, changeReq)
		if changeResp.Code != 20000 {
			t.Fatalf("改密失败: %s", changeResp.Message)
		}

		adminPassword = "NewQuizPass123!"
		adminLogin = doIntegrationPost(t, client, baseURL+"/api/v1/auth/login", fmt.Sprintf(`{"username":"%s","password":"%s","rememberMe":true}`, user.DefaultAdminUsername, adminPassword))
		if adminLogin.Code != 20000 {
			t.Fatalf("改密后管理员登录失败: %s", adminLogin.Message)
		}
		_ = json.Unmarshal(adminLogin.Data, &adminData)
	}
	adminAccess := adminData.AccessToken

	registerReq, _ := http.NewRequest(http.MethodPost, baseURL+"/api/v1/auth/register", bytes.NewBufferString(`{"username":"quiz_flow","password":"TestPass123!","remember":true}`))
	registerReq.Header.Set("Content-Type", "application/json")
	registerReq.Header.Set("Authorization", "Bearer "+adminAccess)
	register := doIntegrationRequest(t, client, registerReq)
	if register.Code != 20000 {
		t.Fatalf("注册失败: %s", register.Message)
	}

	userLogin := doIntegrationPost(t, client, baseURL+"/api/v1/auth/login", `{"username":"quiz_flow","password":"TestPass123!","remember":true}`)
	if userLogin.Code != 20000 {
		t.Fatalf("登录失败: %s", userLogin.Message)
	}
	var tokens map[string]interface{}
	_ = json.Unmarshal(userLogin.Data, &tokens)
	access := fmt.Sprintf("%v", tokens["accessToken"])

	reqQuiz, _ := http.NewRequest(http.MethodGet, baseURL+"/api/v1/quiz/variables/storage", nil)
	reqQuiz.Header.Set("Authorization", "Bearer "+access)
	quizResp := doIntegrationRequest(t, client, reqQuiz)
	if quizResp.Code != 20000 {
		t.Fatalf("获取题目失败: %s", quizResp.Message)
	}

	var quizData struct {
		SessionID string                   `json:"sessionId"`
		Questions []map[string]interface{} `json:"questions"`
	}
	if err := json.Unmarshal(quizResp.Data, &quizData); err != nil {
		t.Fatalf("解析题目响应失败: %v, 原始数据: %s", err, string(quizResp.Data))
	}
	if quizData.SessionID == "" {
		t.Fatalf("题目响应缺少 sessionId, 原始数据: %s", string(quizResp.Data))
	}
	if len(quizData.Questions) == 0 {
		t.Fatalf("题目列表为空")
	}
	first := quizData.Questions[0]

	// 题目ID可能是字符串（YAML题目的hash ID转为字符串）或数字
	var qidStr string
	switch v := first["id"].(type) {
	case string:
		qidStr = v
	case float64:
		qidStr = fmt.Sprintf("%.0f", v)
	case int:
		qidStr = fmt.Sprintf("%d", v)
	case int64:
		qidStr = fmt.Sprintf("%d", v)
	}
	if qidStr == "" || qidStr == "0" {
		t.Fatalf("题目信息不完整, ID为空或0, first question: %+v", first)
	}
	answerChoice := "A"

	submitBody := fmt.Sprintf(`{"sessionId":"%s","topic":"variables","chapter":"storage","durationMs":5000,"answers":[{"questionId":%s,"userAnswers":["%s"]}]}`, quizData.SessionID, qidStr, answerChoice)
	submitReq, _ := http.NewRequest(http.MethodPost, baseURL+"/api/v1/quiz/submit", bytes.NewBufferString(submitBody))
	submitReq.Header.Set("Content-Type", "application/json")
	submitReq.Header.Set("Authorization", "Bearer "+access)
	submitResp := doIntegrationRequest(t, client, submitReq)
	if submitResp.Code != 20000 {
		t.Fatalf("提交测验失败: %s", submitResp.Message)
	}

	historyReq, _ := http.NewRequest(http.MethodGet, baseURL+"/api/v1/quiz/history", nil)
	historyReq.Header.Set("Authorization", "Bearer "+access)
	historyResp := doIntegrationRequest(t, client, historyReq)
	if historyResp.Code != 20000 {
		t.Fatalf("查询历史失败: %s", historyResp.Message)
	}
}
