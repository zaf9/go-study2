package integration

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"go-study2/internal/app/http_server"
	"go-study2/internal/config"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

func TestHttpMode(t *testing.T) {
	// 1. Setup Server
	cfg := &config.Config{
		Server: config.ServerConfig{
			Host: "127.0.0.1",
			Port: 0, // Random port
		},
	}
	s := http_server.NewServer(cfg, "http-mode-test")

	s.SetAccessLogEnabled(false)
	s.Start()
	defer s.Shutdown()

	port := s.GetListenedPort()
	baseUrl := fmt.Sprintf("http://127.0.0.1:%d/api/v1", port)
	client := g.Client()
	client.SetPrefix(baseUrl)
	ctx := gctx.New()

	// Wait for server start
	time.Sleep(100 * time.Millisecond)

	// 2. Test Get Topics (User Story 1/2 bridge)
	t.Run("GetTopics", func(t *testing.T) {
		resp, err := client.Post(ctx, "/topics")
		if err != nil {
			t.Fatalf("Failed to request: %v", err)
		}
		defer resp.Close()

		if resp.StatusCode != 200 {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}

		body := resp.ReadAllString()
		if len(body) == 0 {
			t.Error("Response body is empty")
		}

		if !strings.Contains(body, "topics") {
			t.Error("Response body should contain 'topics'")
		}

		// 验证包含两个主题
		if !strings.Contains(body, "lexical_elements") {
			t.Error("Response should contain 'lexical_elements'")
		}
		if !strings.Contains(body, "constants") {
			t.Error("Response should contain 'constants'")
		}
	})

	// 3. Test Lexical Menu (User Story 2)
	t.Run("GetLexicalMenu", func(t *testing.T) {
		resp, err := client.Post(ctx, "/topic/lexical_elements")
		if err != nil {
			t.Fatalf("Failed to request: %v", err)
		}
		defer resp.Close()

		if resp.StatusCode != 200 {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}

		body := resp.ReadAllString()
		// Check for lexical items count (11 items)
		// We expect 11 occurrences of "id" or "name" or specific titles
		// Let's just check for a few known titles
		if !strings.Contains(body, "Comments") {
			t.Error("Menu should contain 'Comments'")
		}
		if !strings.Contains(body, "Strings") {
			t.Error("Menu should contain 'Strings'")
		}
	})

	// 4. Test Token Content (User Story 2)
	t.Run("GetTokenContent", func(t *testing.T) {
		resp, err := client.Post(ctx, "/topic/lexical_elements/tokens")
		if err != nil {
			t.Fatalf("Failed to request: %v", err)
		}
		defer resp.Close()

		if resp.StatusCode != 200 {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}

		body := resp.ReadAllString()

		if !strings.Contains(body, "关键字") {
			t.Error("Content should contain '关键字'")
		}
	})

	// 5. Test HTML Format
	t.Run("GetContentHTML", func(t *testing.T) {
		resp, err := client.Post(ctx, "/topic/lexical_elements/comments?format=html")
		if err != nil {
			t.Fatalf("Failed to request: %v", err)
		}
		defer resp.Close()

		if resp.StatusCode != 200 {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}

		ct := resp.Header.Get("Content-Type")
		if ct == "" {
			t.Error("Content-Type should not be empty")
		}
		// Allow any charset
		if !strings.Contains(ct, "text/html") {
			t.Errorf("Expected Content-Type text/html, got %s", ct)
		}

		body := resp.ReadAllString()
		if !strings.Contains(body, "<!DOCTYPE html>") {
			t.Error("Response should contain DOCTYPE html")
		}
	})

	// 6. Test 404 for invalid chapter
	t.Run("InvalidChapter", func(t *testing.T) {
		resp, err := client.Post(ctx, "/topic/lexical_elements/invalid_chapter")
		if err != nil {
			t.Fatalf("Failed to request: %v", err)
		}
		defer resp.Close()

		if resp.StatusCode != 404 {
			t.Errorf("Expected 404, got %d", resp.StatusCode)
		}
	})
}
