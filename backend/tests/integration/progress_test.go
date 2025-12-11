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

func TestProgressFlow_EndToEnd(t *testing.T) {
	ctx := gctx.New()
	_ = os.MkdirAll("testdata", 0o755)
	dbPath := filepath.ToSlash(filepath.Join("testdata", fmt.Sprintf("integration_progress_%d.db", time.Now().UnixNano())))
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
			Secret:             "integration-secret-1234567890abcd",
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

	server, err := http_server.NewServer(cfg, "progress-integration")
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

	register := doIntegrationPost(t, client, baseURL+"/api/v1/auth/register", `{"username":"progress_flow","password":"TestPass123","remember":true}`)
	if register.Code != 20000 {
		t.Fatalf("注册失败: %s", register.Message)
	}
	var tokens map[string]interface{}
	_ = json.Unmarshal(register.Data, &tokens)
	access := fmt.Sprintf("%v", tokens["accessToken"])

	progressPayload := `{"topic":"variables","chapter":"storage","status":"in_progress","position":"{\"anchor\":\"section-1\"}"}`
	req, _ := http.NewRequest(http.MethodPost, baseURL+"/api/v1/progress", bytes.NewBufferString(progressPayload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+access)
	saveResp := doIntegrationRequest(t, client, req)
	if saveResp.Code != 20000 {
		t.Fatalf("保存进度失败: %s", saveResp.Message)
	}

	getReq, _ := http.NewRequest(http.MethodGet, baseURL+"/api/v1/progress/variables", nil)
	getReq.Header.Set("Authorization", "Bearer "+access)
	listResp := doIntegrationRequest(t, client, getReq)
	if listResp.Code != 20000 {
		t.Fatalf("查询进度失败: %s", listResp.Message)
	}
}
