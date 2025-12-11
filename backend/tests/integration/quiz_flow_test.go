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
	_ = json.Unmarshal(quizResp.Data, &quizData)
	if quizData.SessionID == "" {
		t.Fatalf("题目响应缺少 sessionId")
	}
	if len(quizData.Questions) == 0 {
		t.Fatalf("题目列表为空")
	}
	first := quizData.Questions[0]
	var qid int64
	switch v := first["id"].(type) {
	case float64:
		qid = int64(v)
	case int:
		qid = int64(v)
	case int64:
		qid = v
	}
	if qid == 0 {
		t.Fatalf("题目信息不完整")
	}
	answerChoice := "A"

	submitBody := fmt.Sprintf(`{"sessionId":"%s","topic":"variables","chapter":"storage","answers":[{"questionId":%d,"userAnswers":["%s"]}]}`, quizData.SessionID, qid, answerChoice)
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
