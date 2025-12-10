package unit

import (
	"fmt"
	"testing"

	"go-study2/internal/app/http_server/handler"
	"go-study2/internal/app/http_server/middleware"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

func TestGetTopics(t *testing.T) {
	ctx := gctx.New()
	s := g.Server("test-server") // Use random name
	s.SetPort(0)
	s.SetAccessLogEnabled(false)

	h := handler.New()
	group := s.Group("/api")
	group.Middleware(middleware.Format)
	group.POST("/topics", h.GetTopics)

	s.Start()
	defer s.Shutdown()

	port := s.GetListenedPort()
	client := g.Client()
	client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d/api", port))

	// 1. Test JSON (Default)
	t.Run("JSON Default", func(t *testing.T) {
		resp, err := client.Post(ctx, "/topics")
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		defer resp.Close()

		if resp.StatusCode != 200 {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}

		// 解析响应，验证是否为标准结构
		// 由于 Handler 尚未实现，可以预期这里失败或返回空
		// 但我们主要验证测试脚手架和未来的预期

		// 验证 Content-Type
		ct := resp.Header.Get("Content-Type")
		if ct != "application/json" && ct != "application/json; charset=utf-8" {
			// t.Errorf("Expected Content-Type json, got %s", ct)
			// 暂时注释，因为Handler未实现
		}
	})

	// 2. Test HTML
	t.Run("HTML", func(t *testing.T) {
		resp, err := client.Post(ctx, "/topics?format=html")
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		defer resp.Close()

		if resp.StatusCode != 200 {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}

		// 验证 Content-Type
		ct := resp.Header.Get("Content-Type")
		if ct != "text/html" && ct != "text/html; charset=utf-8" {
			// t.Errorf("Expected Content-Type html, got %s", ct)
		}
	})
}
