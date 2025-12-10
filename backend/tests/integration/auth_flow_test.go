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

	register := doIntegrationPost(t, client, baseURL+"/api/v1/auth/register", `{"username":"flow_user","password":"TestPass123","rememberMe":true}`)
	if register.Code != 20000 {
		t.Fatalf("注册失败: %v", register.Message)
	}

	var regTokens map[string]interface{}
	_ = json.Unmarshal(register.Data, &regTokens)
	access := fmt.Sprintf("%v", regTokens["accessToken"])
	if access == "" {
		t.Fatalf("注册响应缺少 accessToken")
	}

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
}

func doIntegrationPost(t *testing.T, client *http.Client, url string, payload string) apiResponse {
	t.Helper()
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	return doIntegrationRequest(t, client, req)
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
