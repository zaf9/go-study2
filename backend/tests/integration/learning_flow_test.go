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

type learningResp struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

// TestLearningFlow_EndToEnd 覆盖学习主题列表、章节内容与进度记录的端到端流程。
func TestLearningFlow_EndToEnd(t *testing.T) {
	ctx := gctx.New()
	_ = os.MkdirAll("testdata", 0o755)
	dbPath := filepath.ToSlash(filepath.Join("testdata", fmt.Sprintf("integration_learning_%d.db", time.Now().UnixNano())))
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
			Secret:             "learning-flow-secret-123456",
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

	server, err := http_server.NewServer(cfg, "learning-integration")
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

	adminLogin := doLearningPost(t, client, baseURL+"/api/v1/auth/login", fmt.Sprintf(`{"username":"%s","password":"%s","remember":true}`, user.DefaultAdminUsername, user.DefaultAdminPassword))
	if adminLogin.Code != 20000 {
		t.Fatalf("管理员登录失败: %s", adminLogin.Message)
	}
	var adminTokens struct {
		AccessToken        string `json:"accessToken"`
		NeedPasswordChange bool   `json:"needPasswordChange"`
	}
	_ = json.Unmarshal(adminLogin.Data, &adminTokens)
	adminPassword := user.DefaultAdminPassword
	if adminTokens.NeedPasswordChange {
		changeReqBody := fmt.Sprintf(`{"oldPassword":"%s","newPassword":"LearningFlow123!"}`, adminPassword)
		changeReq, _ := http.NewRequest(http.MethodPost, baseURL+"/api/v1/auth/change-password", bytes.NewBufferString(changeReqBody))
		changeReq.Header.Set("Content-Type", "application/json")
		changeReq.Header.Set("Authorization", "Bearer "+adminTokens.AccessToken)
		changeResp := doLearningRequest(t, client, changeReq)
		if changeResp.Code != 20000 {
			t.Fatalf("改密失败: %s", changeResp.Message)
		}
		adminPassword = "LearningFlow123!"
		adminLogin = doLearningPost(t, client, baseURL+"/api/v1/auth/login", fmt.Sprintf(`{"username":"%s","password":"%s","remember":true}`, user.DefaultAdminUsername, adminPassword))
		if adminLogin.Code != 20000 {
			t.Fatalf("改密后管理员登录失败: %s", adminLogin.Message)
		}
		_ = json.Unmarshal(adminLogin.Data, &adminTokens)
	}
	adminAccess := adminTokens.AccessToken

	registerReq, _ := http.NewRequest(http.MethodPost, baseURL+"/api/v1/auth/register", bytes.NewBufferString(`{"username":"learning_user","password":"TestPass123!","remember":true}`))
	registerReq.Header.Set("Content-Type", "application/json")
	registerReq.Header.Set("Authorization", "Bearer "+adminAccess)
	register := doLearningRequest(t, client, registerReq)
	if register.Code != 20000 {
		t.Fatalf("注册失败: %s", register.Message)
	}

	userLogin := doLearningPost(t, client, baseURL+"/api/v1/auth/login", `{"username":"learning_user","password":"TestPass123!","remember":true}`)
	if userLogin.Code != 20000 {
		t.Fatalf("登录失败: %s", userLogin.Message)
	}
	var tokens map[string]interface{}
	_ = json.Unmarshal(userLogin.Data, &tokens)
	access := fmt.Sprintf("%v", tokens["accessToken"])
	if access == "" {
		t.Fatalf("登录响应缺少 accessToken")
	}

	// 2) 获取主题列表
	topicsResp := doLearningGet(t, client, baseURL+"/api/v1/topics?format=json")
	if topicsResp.Code != 20000 {
		t.Fatalf("获取主题列表失败: %s", topicsResp.Message)
	}
	var topicData struct {
		Topics []map[string]interface{} `json:"topics"`
	}
	_ = json.Unmarshal(topicsResp.Data, &topicData)
	if len(topicData.Topics) == 0 {
		t.Fatalf("主题列表为空，期望至少一个主题")
	}
	firstTopic := fmt.Sprintf("%v", topicData.Topics[0]["id"])

	// 3) 获取章节菜单
	menuResp := doLearningGet(t, client, fmt.Sprintf("%s/api/v1/topic/%s?format=json", baseURL, firstTopic))
	if menuResp.Code != 20000 {
		t.Fatalf("获取章节菜单失败: %s", menuResp.Message)
	}
	var menuData struct {
		Items []map[string]interface{} `json:"items"`
	}
	_ = json.Unmarshal(menuResp.Data, &menuData)
	if len(menuData.Items) == 0 {
		t.Fatalf("章节菜单为空")
	}
	firstChapter := fmt.Sprintf("%v", menuData.Items[0]["name"])

	// 4) 获取章节内容
	contentResp := doLearningGet(t, client, fmt.Sprintf("%s/api/v1/topic/%s/%s?format=json", baseURL, firstTopic, firstChapter))
	if contentResp.Code != 20000 {
		t.Fatalf("获取章节内容失败: %s", contentResp.Message)
	}
	var contentData map[string]string
	_ = json.Unmarshal(contentResp.Data, &contentData)
	if contentData["title"] == "" || contentData["content"] == "" {
		t.Fatalf("章节内容缺失 title 或 content")
	}

	// 5) 记录学习进度
	progressPayload := fmt.Sprintf(`{"topic":"%s","chapter":"%s","status":"in_progress","position":"{\"section\":1}"}`, firstTopic, firstChapter)
	saveReq, _ := http.NewRequest(http.MethodPost, baseURL+"/api/v1/progress", bytes.NewBufferString(progressPayload))
	saveReq.Header.Set("Content-Type", "application/json")
	saveReq.Header.Set("Authorization", "Bearer "+access)
	saveResp := doLearningRequest(t, client, saveReq)
	if saveResp.Code != 20000 {
		t.Fatalf("保存进度失败: %s", saveResp.Message)
	}

	// 6) 查询进度验证写入
	getReq, _ := http.NewRequest(http.MethodGet, baseURL+fmt.Sprintf("/api/v1/progress/%s", firstTopic), nil)
	getReq.Header.Set("Authorization", "Bearer "+access)
	progressResp := doLearningRequest(t, client, getReq)
	if progressResp.Code != 20000 {
		t.Fatalf("查询进度失败: %s", progressResp.Message)
	}
	var progressList []map[string]interface{}
	_ = json.Unmarshal(progressResp.Data, &progressList)
	if len(progressList) == 0 {
		t.Fatalf("进度列表为空，期望至少一条记录")
	}
	if fmt.Sprintf("%v", progressList[0]["chapter"]) == "" {
		t.Fatalf("进度记录缺少章节信息")
	}
}

func doLearningPost(t *testing.T, client *http.Client, url string, payload string) learningResp {
	t.Helper()
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	return doLearningRequest(t, client, req)
}

func doLearningGet(t *testing.T, client *http.Client, url string) learningResp {
	t.Helper()
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	return doLearningRequest(t, client, req)
}

func doLearningRequest(t *testing.T, client *http.Client, req *http.Request) learningResp {
	t.Helper()
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	var body learningResp
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}
	return body
}
