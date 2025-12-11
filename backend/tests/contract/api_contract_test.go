package contract

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
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

type contractAPIResp struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

// TestAPIContract_Responses 验证基础 API 响应结构与字段完整性。
func TestAPIContract_Responses(t *testing.T) {
	baseURL, client, shutdown := startGenericServer(t)
	defer shutdown()

	// 1) topics 响应结构
	topics := doContractRequest(t, client, mustNewRequest(t, http.MethodGet, baseURL+"/api/v1/topics?format=json", ""))
	if topics.Code != 20000 {
		t.Fatalf("topics 返回码错误: %d, msg=%s", topics.Code, topics.Message)
	}
	var topicData struct {
		Topics []map[string]interface{} `json:"topics"`
	}
	_ = json.Unmarshal(topics.Data, &topicData)
	if len(topicData.Topics) == 0 {
		t.Fatalf("topics 数据应至少包含一项")
	}
	if topicData.Topics[0]["id"] == nil || topicData.Topics[0]["title"] == nil {
		t.Fatalf("topics 数据缺少 id 或 title 字段")
	}

	// 2) 认证 + 进度契约
	register := doContractRequest(t, client, mustNewRequest(t, http.MethodPost, baseURL+"/api/v1/auth/register", `{"username":"contract_api_user","password":"TestPass123","remember":true}`))
	if register.Code != 20000 {
		t.Fatalf("注册失败: %s", register.Message)
	}
	var regData map[string]interface{}
	_ = json.Unmarshal(register.Data, &regData)
	access := fmt.Sprintf("%v", regData["accessToken"])
	if access == "" {
		t.Fatalf("注册响应缺少 accessToken")
	}

	// 3) 受保护接口返回 20000 且 data 为数组
	req, _ := http.NewRequest(http.MethodGet, baseURL+"/api/v1/progress", nil)
	req.Header.Set("Authorization", "Bearer "+access)
	progressResp := doContractRequest(t, client, req)
	if progressResp.Code != 20000 {
		t.Fatalf("进度接口返回码异常: %d, msg=%s", progressResp.Code, progressResp.Message)
	}
	var progressList []map[string]interface{}
	_ = json.Unmarshal(progressResp.Data, &progressList)
	if progressList == nil {
		t.Fatalf("进度接口 data 解析失败")
	}
}

func startGenericServer(t *testing.T) (string, *http.Client, func()) {
	t.Helper()
	ensureConfigPath()
	ctx := gctx.New()
	_ = gfile.Mkdir("testdata")
	dbPath := filepath.ToSlash(filepath.Join("testdata", fmt.Sprintf("contract_api_%d.db", time.Now().UnixNano())))

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
			Secret:             "contract-api-secret-123456",
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

	server, err := http_server.NewServer(cfg, "contract-api")
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

	return baseURL, client, func() { server.Shutdown() }
}

func doContractRequest(t *testing.T, client *http.Client, req *http.Request) contractAPIResp {
	t.Helper()
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	var body contractAPIResp
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}
	return body
}

func mustNewRequest(t *testing.T, method, url, payload string) *http.Request {
	t.Helper()
	req, err := http.NewRequest(method, url, bytes.NewBufferString(payload))
	if err != nil {
		t.Fatalf("构造请求失败: %v", err)
	}
	if payload != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	return req
}

