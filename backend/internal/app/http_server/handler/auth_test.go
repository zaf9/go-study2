package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
	"go-study2/internal/pkg/password"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/util/guid"
)

type authResp struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

func TestAuthHandlers_Flow(t *testing.T) {
	baseURL, client, shutdown := startAuthServer(t)
	defer shutdown()

	adminPwd, _ := password.Hash("Admin123!")
	if _, err := database.Default().Insert(gctx.New(), "users", map[string]interface{}{
		"username":             "admin",
		"password_hash":        adminPwd,
		"is_admin":             1,
		"status":               "active",
		"must_change_password": 0,
	}); err != nil {
		t.Fatalf("创建管理员失败: %v", err)
	}

	adminLoginPayload := `{"username":"admin","password":"Admin123!","rememberMe":true}`
	adminLogin := doPost(t, client, baseURL+"/api/v1/auth/login", adminLoginPayload)
	if adminLogin.Code != 20000 {
		t.Fatalf("管理员登录失败: %d, msg=%s", adminLogin.Code, adminLogin.Message)
	}
	var adminTokens authResponse
	if err := json.Unmarshal(adminLogin.Data, &adminTokens); err != nil {
		t.Fatalf("解析管理员登录响应失败: %v", err)
	}
	if adminTokens.AccessToken == "" {
		t.Fatalf("管理员登录未返回 accessToken")
	}

	registerReq, _ := http.NewRequest(http.MethodPost, baseURL+"/api/v1/auth/register", bytes.NewBufferString(`{"username":"auth_user","password":"TestPass123!","rememberMe":true}`))
	registerReq.Header.Set("Content-Type", "application/json")
	registerReq.Header.Set("Authorization", "Bearer "+adminTokens.AccessToken)
	registerResp := doRequest(t, client, registerReq)
	if registerResp.Code != 20000 {
		t.Fatalf("注册返回错误码: %d, msg=%s", registerResp.Code, registerResp.Message)
	}

	loginResp := doPost(t, client, baseURL+"/api/v1/auth/login", `{"username":"auth_user","password":"TestPass123!","rememberMe":true}`)
	if loginResp.Code != 20000 {
		t.Fatalf("登录返回错误码: %d, msg=%s", loginResp.Code, loginResp.Message)
	}

	var tokens authResponse
	if err := json.Unmarshal(loginResp.Data, &tokens); err != nil {
		t.Fatalf("解析登录响应失败: %v", err)
	}
	if tokens.AccessToken == "" {
		t.Fatalf("登录响应缺少 accessToken")
	}

	refreshResp := doPost(t, client, baseURL+"/api/v1/auth/refresh", `{}`)
	if refreshResp.Code != 20000 {
		t.Fatalf("刷新返回错误码: %d", refreshResp.Code)
	}

	reqProfile, _ := http.NewRequest(http.MethodGet, baseURL+"/api/v1/auth/profile", nil)
	reqProfile.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
	profileResp := doRequest(t, client, reqProfile)
	if profileResp.Code != 20000 {
		t.Fatalf("Profile 返回错误码: %d", profileResp.Code)
	}

	logoutReq, _ := http.NewRequest(http.MethodPost, baseURL+"/api/v1/auth/logout", nil)
	logoutReq.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
	logoutResp := doRequest(t, client, logoutReq)
	if logoutResp.Code != 20000 {
		t.Fatalf("退出返回错误码: %d", logoutResp.Code)
	}

	profileAfterLogout := doRequest(t, client, reqProfile)
	if profileAfterLogout.Code == 20000 {
		t.Fatalf("退出后访问 profile 不应成功")
	}
}

func startAuthServer(t *testing.T) (string, *http.Client, func()) {
	t.Helper()
	ensureConfigPath()
	_ = os.MkdirAll("testdata", 0o755)
	dbPath := filepath.ToSlash(filepath.Join("testdata", fmt.Sprintf("auth_handler_%d.db", time.Now().UnixNano())))
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
	db, err := database.Init(gctx.New(), cfg)
	if err != nil {
		t.Fatalf("初始化数据库失败: %v", err)
	}
	_, _ = db.Tables(gctx.New()) // 确保迁移成功
	_, _ = db.TableFields(gctx.New(), "users")

	if err := appjwt.Configure(appjwt.Options{
		Secret:             "1234567890abcdef",
		AccessTokenExpiry:  time.Hour,
		RefreshTokenExpiry: 24 * time.Hour,
	}); err != nil {
		t.Fatalf("配置 JWT 失败: %v", err)
	}

	server := g.Server(guid.S())
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
			authGroup.GET("/auth/profile", h.GetProfile)
			authGroup.POST("/auth/logout", h.Logout)
		})
	})
	server.Start()

	time.Sleep(50 * time.Millisecond)
	port := server.GetListenedPort()
	baseURL := fmt.Sprintf("http://127.0.0.1:%d", port)

	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}

	cleanup := func() {
		server.Shutdown()
	}

	return baseURL, client, cleanup
}

func doPost(t *testing.T, client *http.Client, url string, payload string) authResp {
	t.Helper()
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	return doRequest(t, client, req)
}

func doRequest(t *testing.T, client *http.Client, req *http.Request) authResp {
	t.Helper()
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	var body authResp
	if err := json.Unmarshal(bodyBytes, &body); err != nil {
		t.Fatalf("解析响应失败: %v, body=%s", err, string(bodyBytes))
	}
	return body
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
