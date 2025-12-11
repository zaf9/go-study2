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

type apiResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

func TestAuthFlow_EndToEnd(t *testing.T) {
	ctx := gctx.New()
	_ = os.MkdirAll("testdata", 0o755)
	dbPath := filepath.ToSlash(filepath.Join("testdata", fmt.Sprintf("integration_auth_%d.db", time.Now().UnixNano())))
	cfg := &config.Config{
		Server: config.ServerConfig{
			Host: "127.0.0.1",
		},
		Http: config.HttpConfig{
			Port: 0,
		},
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
			Secret:             "abcdef1234567890abcdef1234567890",
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

	server, err := http_server.NewServer(cfg, "auth-integration")
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
		t.Fatalf("管理员登录失败: %v", adminLogin.Message)
	}
	var adminLoginData struct {
		AccessToken        string `json:"accessToken"`
		NeedPasswordChange bool   `json:"needPasswordChange"`
	}
	_ = json.Unmarshal(adminLogin.Data, &adminLoginData)
	adminPassword := user.DefaultAdminPassword
	if adminLoginData.NeedPasswordChange {
		changeReqBody := fmt.Sprintf(`{"oldPassword":"%s","newPassword":"AdminFlow123!"}`, adminPassword)
		changeReq, _ := http.NewRequest(http.MethodPost, baseURL+"/api/v1/auth/change-password", bytes.NewBufferString(changeReqBody))
		changeReq.Header.Set("Content-Type", "application/json")
		changeReq.Header.Set("Authorization", "Bearer "+adminLoginData.AccessToken)
		changeResp := doIntegrationRequest(t, client, changeReq)
		if changeResp.Code != 20000 {
			t.Fatalf("改密失败: %v", changeResp.Message)
		}
		adminPassword = "AdminFlow123!"
		adminLogin = doIntegrationPost(t, client, baseURL+"/api/v1/auth/login", fmt.Sprintf(`{"username":"%s","password":"%s","rememberMe":true}`, user.DefaultAdminUsername, adminPassword))
		if adminLogin.Code != 20000 {
			t.Fatalf("改密后管理员登录失败: %v", adminLogin.Message)
		}
		_ = json.Unmarshal(adminLogin.Data, &adminLoginData)
	}
	adminAccess := adminLoginData.AccessToken

	registerReq, _ := http.NewRequest(http.MethodPost, baseURL+"/api/v1/auth/register", bytes.NewBufferString(`{"username":"flow_user","password":"TestPass123!","rememberMe":true}`))
	registerReq.Header.Set("Content-Type", "application/json")
	registerReq.Header.Set("Authorization", "Bearer "+adminAccess)
	register := doIntegrationRequest(t, client, registerReq)
	if register.Code != 20000 {
		t.Fatalf("注册失败: %v", register.Message)
	}

	userLogin := doIntegrationPost(t, client, baseURL+"/api/v1/auth/login", `{"username":"flow_user","password":"TestPass123!","rememberMe":true}`)
	if userLogin.Code != 20000 {
		t.Fatalf("用户登录失败: %v", userLogin.Message)
	}
	var loginTokens map[string]interface{}
	_ = json.Unmarshal(userLogin.Data, &loginTokens)
	access := fmt.Sprintf("%v", loginTokens["accessToken"])

	refresh := doIntegrationPost(t, client, baseURL+"/api/v1/auth/refresh", `{}`)
	if refresh.Code != 20000 {
		t.Fatalf("刷新失败: %v", refresh.Message)
	}

	reqProfile, _ := http.NewRequest(http.MethodGet, baseURL+"/api/v1/auth/profile", nil)
	reqProfile.Header.Set("Authorization", "Bearer "+access)
	profile := doIntegrationRequest(t, client, reqProfile)
	if profile.Code != 20000 {
		t.Fatalf("获取 profile 失败: %v", profile.Message)
	}

	logoutReq, _ := http.NewRequest(http.MethodPost, baseURL+"/api/v1/auth/logout", nil)
	logoutReq.Header.Set("Authorization", "Bearer "+access)
	logout := doIntegrationRequest(t, client, logoutReq)
	if logout.Code != 20000 {
		t.Fatalf("退出失败: %v", logout.Message)
	}

	profileAfter, _ := http.NewRequest(http.MethodGet, baseURL+"/api/v1/auth/profile", nil)
	profileAfter.Header.Set("Authorization", "Bearer "+access)
	afterResp := doIntegrationRequest(t, client, profileAfter)
	if afterResp.Code == 20000 {
		t.Fatalf("退出后仍能访问 profile，不符合预期")
	}

	count := countAuditEvents(t, "register_success")
	if count == 0 {
		t.Fatalf("期望存在注册成功审计事件")
	}
}

func TestAuthFlow_ForceChangePassword(t *testing.T) {
	ctx := gctx.New()
	_ = os.MkdirAll("testdata", 0o755)
	dbPath := filepath.ToSlash(filepath.Join("testdata", fmt.Sprintf("integration_force_change_%d.db", time.Now().UnixNano())))

	cfg := &config.Config{
		Server: config.ServerConfig{
			Host: "127.0.0.1",
		},
		Http: config.HttpConfig{
			Port: 0,
		},
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
			Secret:             "force-change-secret-abcdef123456",
			AccessTokenExpiry:  1800,
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

	server, err := http_server.NewServer(cfg, "auth-force-change")
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

	loginResp := doIntegrationPost(t, client, baseURL+"/api/v1/auth/login", fmt.Sprintf(`{"username":"%s","password":"%s","rememberMe":true}`, user.DefaultAdminUsername, user.DefaultAdminPassword))
	if loginResp.Code != 20000 {
		t.Fatalf("默认管理员登录失败: %v", loginResp.Message)
	}
	var loginData struct {
		AccessToken        string `json:"accessToken"`
		NeedPasswordChange bool   `json:"needPasswordChange"`
	}
	_ = json.Unmarshal(loginResp.Data, &loginData)
	if !loginData.NeedPasswordChange {
		t.Fatalf("默认管理员未返回需改密标记")
	}

	protectedReq, _ := http.NewRequest(http.MethodGet, baseURL+"/api/v1/progress", nil)
	protectedReq.Header.Set("Authorization", "Bearer "+loginData.AccessToken)
	blocked := doIntegrationRequest(t, client, protectedReq)
	if blocked.Code != 40011 {
		t.Fatalf("需改密用户访问业务接口应被阻断，得到 code=%d", blocked.Code)
	}

	changeReqBody := fmt.Sprintf(`{"oldPassword":"%s","newPassword":"NewPass123!"}`, user.DefaultAdminPassword)
	changeReq, _ := http.NewRequest(http.MethodPost, baseURL+"/api/v1/auth/change-password", bytes.NewBufferString(changeReqBody))
	changeReq.Header.Set("Content-Type", "application/json")
	changeReq.Header.Set("Authorization", "Bearer "+loginData.AccessToken)
	changeResp := doIntegrationRequest(t, client, changeReq)
	if changeResp.Code != 20000 {
		t.Fatalf("改密失败: %v", changeResp.Message)
	}

	profileReq, _ := http.NewRequest(http.MethodGet, baseURL+"/api/v1/auth/profile", nil)
	profileReq.Header.Set("Authorization", "Bearer "+loginData.AccessToken)
	afterChange := doIntegrationRequest(t, client, profileReq)
	if afterChange.Code == 20000 {
		t.Fatalf("旧令牌不应在改密后继续访问成功")
	}

	relogin := doIntegrationPost(t, client, baseURL+"/api/v1/auth/login", `{"username":"admin","password":"NewPass123!","rememberMe":true}`)
	if relogin.Code != 20000 {
		t.Fatalf("改密后重新登录失败: %v", relogin.Message)
	}
	var reloginData struct {
		AccessToken        string `json:"accessToken"`
		NeedPasswordChange bool   `json:"needPasswordChange"`
	}
	_ = json.Unmarshal(relogin.Data, &reloginData)
	if reloginData.NeedPasswordChange {
		t.Fatalf("改密后不应仍要求改密")
	}

	progressReq, _ := http.NewRequest(http.MethodGet, baseURL+"/api/v1/progress", nil)
	progressReq.Header.Set("Authorization", "Bearer "+reloginData.AccessToken)
	progressResp := doIntegrationRequest(t, client, progressReq)
	if progressResp.Code != 20000 {
		t.Fatalf("改密后访问业务接口失败: %v", progressResp.Message)
	}

	changedCount := countAuditEvents(t, "password_changed")
	if changedCount == 0 {
		t.Fatalf("改密应记录审计事件")
	}
	blockedCount := countAuditEvents(t, "access_blocked_need_change")
	if blockedCount == 0 {
		t.Fatalf("需改密拦截应记录审计事件")
	}
}

func doIntegrationPost(t *testing.T, client *http.Client, url string, payload string) apiResponse {
	t.Helper()
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	return doIntegrationRequest(t, client, req)
}

func countAuditEvents(t *testing.T, eventType string) int {
	t.Helper()
	db := database.Default()
	if db == nil {
		t.Fatalf("数据库未初始化")
	}
	ctx := gctx.New()
	cnt, err := db.Model("audit_events").Where("event_type", eventType).Count(ctx)
	if err != nil {
		t.Fatalf("统计审计事件失败: %v", err)
	}
	return cnt
}

func doIntegrationRequest(t *testing.T, client *http.Client, req *http.Request) apiResponse {
	t.Helper()
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	var body apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}
	return body
}
