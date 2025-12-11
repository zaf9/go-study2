package integration

import (
	"bytes"
	"context"
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
	"go-study2/internal/pkg/password"

	"github.com/gogf/gf/v2/os/gctx"
)

func TestDefaultAdminCreatedWhenMissing(t *testing.T) {
	ctx := gctx.New()
	_ = os.MkdirAll("testdata", 0o755)
	dbPath := filepath.ToSlash(filepath.Join("testdata", fmt.Sprintf("default_admin_%d.db", time.Now().UnixNano())))

	baseURL, shutdown := startServerWithDB(t, ctx, dbPath)
	defer shutdown()

	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	login := postJSON(t, client, baseURL+"/api/v1/auth/login", fmt.Sprintf(`{"username":"%s","password":"%s","rememberMe":true}`, user.DefaultAdminUsername, user.DefaultAdminPassword))
	if login.Code != 20000 {
		t.Fatalf("默认管理员登录失败: %v", login.Message)
	}

	record, err := database.Default().Model("users").Where("username", user.DefaultAdminUsername).One(ctx)
	if err != nil {
		t.Fatalf("查询默认管理员失败: %v", err)
	}
	var admin user.User
	if err := record.Struct(&admin); err != nil {
		t.Fatalf("结构体映射失败: %v", err)
	}
	if !admin.IsAdmin || !admin.MustChangePassword {
		t.Fatalf("默认管理员属性不符合预期: %+v", admin)
	}

	count, err := database.Default().Model("audit_events").Where("event_type", "default_admin_created").Count(ctx)
	if err != nil {
		t.Fatalf("查询审计事件失败: %v", err)
	}
	if count == 0 {
		t.Fatalf("应记录默认管理员创建审计事件")
	}
}

func TestDefaultAdminNotOverwritten(t *testing.T) {
	ctx := gctx.New()
	_ = os.MkdirAll("testdata", 0o755)
	dbPath := filepath.ToSlash(filepath.Join("testdata", fmt.Sprintf("default_admin_existing_%d.db", time.Now().UnixNano())))

	cfg := defaultConfig(dbPath)
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

	customPwd := "KeepPass1!"
	hashed, _ := password.Hash(customPwd)
	if _, err := database.Default().Insert(ctx, "users", map[string]interface{}{
		"username":             user.DefaultAdminUsername,
		"password_hash":        hashed,
		"is_admin":             1,
		"status":               "active",
		"must_change_password": 0,
	}); err != nil {
		t.Fatalf("预置管理员失败: %v", err)
	}

	server, err := http_server.NewServer(cfg, "default-admin-existing")
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

	login := postJSON(t, client, baseURL+"/api/v1/auth/login", fmt.Sprintf(`{"username":"%s","password":"%s","rememberMe":true}`, user.DefaultAdminUsername, customPwd))
	if login.Code != 20000 {
		t.Fatalf("已存在管理员登录失败: %v", login.Message)
	}

	record, err := database.Default().Model("users").Where("username", user.DefaultAdminUsername).One(ctx)
	if err != nil {
		t.Fatalf("查询管理员失败: %v", err)
	}
	var admin user.User
	if err := record.Struct(&admin); err != nil {
		t.Fatalf("结构体映射失败: %v", err)
	}
	if admin.MustChangePassword {
		t.Fatalf("已有管理员不应被重置 must_change_password")
	}
}

func startServerWithDB(t *testing.T, ctx context.Context, dbPath string) (string, func()) {
	t.Helper()
	cfg := defaultConfig(dbPath)

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

	server, err := http_server.NewServer(cfg, "default-admin")
	if err != nil {
		t.Fatalf("创建服务器失败: %v", err)
	}
	server.SetPort(0)
	server.SetAccessLogEnabled(false)
	server.Start()

	time.Sleep(80 * time.Millisecond)
	baseURL := fmt.Sprintf("http://127.0.0.1:%d", server.GetListenedPort())

	cleanup := func() {
		server.Shutdown()
	}
	return baseURL, cleanup
}

func defaultConfig(dbPath string) *config.Config {
	return &config.Config{
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
			Secret:             "default-admin-secret-abcdef123456",
			AccessTokenExpiry:  3600,
			RefreshTokenExpiry: 604800,
			Issuer:             "go-study2",
		},
	}
}

func postJSON(t *testing.T, client *http.Client, url string, payload string) apiResponse {
	t.Helper()
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	return doIntegrationRequest(t, client, req)
}
