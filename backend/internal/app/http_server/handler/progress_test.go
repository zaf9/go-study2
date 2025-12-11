package handler

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

	"go-study2/internal/app/http_server/middleware"
	"go-study2/internal/config"
	"go-study2/internal/infrastructure/database"
	appjwt "go-study2/internal/pkg/jwt"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/guid"
)

type progressResp struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

func TestProgressHandlers_Flow(t *testing.T) {
	ctx := gctx.New()
	_ = os.MkdirAll("testdata", 0o755)
	dbPath := filepath.ToSlash(filepath.Join("testdata", fmt.Sprintf("progress_handler_%d.db", time.Now().UnixNano())))
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
	if _, err := database.Init(ctx, cfg); err != nil {
		t.Fatalf("初始化数据库失败: %v", err)
	}

	if err := appjwt.Configure(appjwt.Options{
		Secret:             "1234567890abcdef1234567890abcdef",
		AccessTokenExpiry:  time.Hour,
		RefreshTokenExpiry: 24 * time.Hour,
	}); err != nil {
		t.Fatalf("配置 JWT 失败: %v", err)
	}

	server := ghttp.GetServer(guid.S())
	server.SetPort(0)
	server.SetAccessLogEnabled(false)
	h := New()
	server.Group("/api/v1", func(group *ghttp.RouterGroup) {
		group.Middleware(middleware.Format)
		group.POST("/auth/register", h.Register)
		group.POST("/auth/login", h.Login)
		group.POST("/auth/refresh", h.RefreshToken)
		group.Group("/", func(authGroup *ghttp.RouterGroup) {
			authGroup.Middleware(middleware.Auth)
			authGroup.POST("/progress", h.SaveProgress)
			authGroup.GET("/progress", h.GetAllProgress)
		})
	})
	server.Start()
	defer server.Shutdown()
	time.Sleep(60 * time.Millisecond)

	baseURL := fmt.Sprintf("http://127.0.0.1:%d", server.GetListenedPort())
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	register := doProgressPost(t, client, baseURL+"/api/v1/auth/register", `{"username":"progress_user","password":"TestPass123","remember":true}`)
	if register.Code != 20000 {
		t.Fatalf("注册失败: %s", register.Message)
	}
	var tokens map[string]interface{}
	_ = json.Unmarshal(register.Data, &tokens)
	access := fmt.Sprintf("%v", tokens["accessToken"])

	payload := `{"topic":"variables","chapter":"storage","status":"in_progress","position":"{\"scroll\":100}"}`
	req, _ := http.NewRequest(http.MethodPost, baseURL+"/api/v1/progress", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+access)
	saveResp := doProgressRequest(t, client, req)
	if saveResp.Code != 20000 {
		t.Fatalf("保存进度失败: %s", saveResp.Message)
	}

	getReq, _ := http.NewRequest(http.MethodGet, baseURL+"/api/v1/progress", nil)
	getReq.Header.Set("Authorization", "Bearer "+access)
	listResp := doProgressRequest(t, client, getReq)
	if listResp.Code != 20000 {
		t.Fatalf("查询进度失败: %s", listResp.Message)
	}

	var list []map[string]interface{}
	_ = json.Unmarshal(listResp.Data, &list)
	if len(list) == 0 {
		t.Fatalf("进度列表为空，不符合预期")
	}
}

func doProgressPost(t *testing.T, client *http.Client, url string, payload string) progressResp {
	t.Helper()
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	return doProgressRequest(t, client, req)
}

func doProgressRequest(t *testing.T, client *http.Client, req *http.Request) progressResp {
	t.Helper()
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	var body progressResp
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}
	return body
}
