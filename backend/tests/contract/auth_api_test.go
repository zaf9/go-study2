package contract

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

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
)

type contractResp struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

func TestAuthAPI_Contract(t *testing.T) {
	baseURL, client, shutdown := startContractServer(t)
	defer shutdown()

	register := doContractPost(t, client, baseURL+"/api/v1/auth/register", `{"username":"contract_user","password":"TestPass123","rememberMe":true}`)
	if register.Code != 20000 {
		t.Fatalf("注册返回码错误: %d, msg=%s", register.Code, register.Message)
	}

	var regData map[string]interface{}
	if err := json.Unmarshal(register.Data, &regData); err != nil {
		t.Fatalf("解析注册 data 失败: %v", err)
	}
	if regData["accessToken"] == "" || regData["expiresIn"] == nil {
		t.Fatalf("注册响应缺少 accessToken 或 expiresIn")
	}

	invalidLogin := doContractPost(t, client, baseURL+"/api/v1/auth/login", `{"username":"","password":""}`)
	if invalidLogin.Code != 40004 {
		t.Fatalf("无效登录参数应返回 40004，实际 %d", invalidLogin.Code)
	}

	refresh := doContractPost(t, client, baseURL+"/api/v1/auth/refresh", `{}`)
	if refresh.Code != 20000 && refresh.Code != 40002 {
		t.Fatalf("刷新接口应返回成功或 40002，实际 %d", refresh.Code)
	}
}

func startContractServer(t *testing.T) (string, *http.Client, func()) {
	t.Helper()
	ensureConfigPath()
	ctx := gctx.New()
	_ = os.MkdirAll("testdata", 0o755)
	dbPath := filepath.ToSlash(filepath.Join("testdata", fmt.Sprintf("contract_auth_%d.db", time.Now().UnixNano())))
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
			Secret:             "fedcba9876543210fedcba9876543210",
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

	server, err := http_server.NewServer(cfg, "contract-auth")
	if err != nil {
		t.Fatalf("创建服务器失败: %v", err)
	}
	server.SetPort(0)
	server.SetAccessLogEnabled(false)
	server.Start()

	time.Sleep(80 * time.Millisecond)
	baseURL := fmt.Sprintf("http://127.0.0.1:%d", server.GetListenedPort())

	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	cleanup := func() {
		server.Shutdown()
	}

	return baseURL, client, cleanup
}

func ensureConfigPath() {
	adapter, ok := g.Cfg().GetAdapter().(*gcfg.AdapterFile)
	if !ok {
		return
	}
	if configFile, err := gfile.Search("configs/config.yaml"); err == nil && configFile != "" {
		_ = adapter.SetPath(filepath.Dir(configFile))
		adapter.SetFileName("config")
	}
}

func doContractPost(t *testing.T, client *http.Client, url string, payload string) contractResp {
	t.Helper()
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	return doContractRequest(t, client, req)
}

func doContractRequest(t *testing.T, client *http.Client, req *http.Request) contractResp {
	t.Helper()
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	var body contractResp
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}
	return body
}
