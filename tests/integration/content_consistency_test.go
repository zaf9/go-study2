package integration

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"go-study2/internal/app/http_server"
	"go-study2/internal/app/lexical_elements"
	"go-study2/internal/config"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

func TestContentConsistency(t *testing.T) {
	// 1. Setup Server
	cfg := &config.Config{
		Server: config.ServerConfig{
			Host: "127.0.0.1",
			Port: 0, // Random port
		},
	}
	s := http_server.NewServer(cfg, "consistency-test")
	s.SetAccessLogEnabled(false)
	s.Start()
	defer s.Shutdown()

	port := s.GetListenedPort()
	baseUrl := fmt.Sprintf("http://127.0.0.1:%d/api/v1", port)
	client := g.Client()
	client.SetPrefix(baseUrl)
	ctx := gctx.New()

	time.Sleep(100 * time.Millisecond)

	// Regex to match pointer addresses
	rePtr := regexp.MustCompile(`0x[0-9a-fA-F]+`)

	// Define test cases mapping API chapter ID to Direct Function
	tests := []struct {
		chapterID   string
		contentFunc func() string
		sanitizer   func(string) string
	}{
		{"comments", lexical_elements.GetCommentsContent, nil},
		{"tokens", lexical_elements.GetTokensContent, nil},
		{"semicolons", lexical_elements.GetSemicolonsContent, nil},
		{"identifiers", lexical_elements.GetIdentifiersContent, nil},
		{"keywords", lexical_elements.GetKeywordsContent, nil},
		{"operators", lexical_elements.GetOperatorsContent, func(s string) string {
			return rePtr.ReplaceAllString(s, "PTR")
		}},
		{"integers", lexical_elements.GetIntegersContent, nil},
		{"floats", lexical_elements.GetFloatsContent, nil},
		{"imaginary", lexical_elements.GetImaginaryContent, nil},
		{"runes", lexical_elements.GetRunesContent, nil},
		{"strings", lexical_elements.GetStringsContent, nil},
	}

	for _, tt := range tests {
		t.Run(tt.chapterID, func(t *testing.T) {
			// 1. Get Expected Content (CLI Mode logic)
			expected := tt.contentFunc()

			// 2. Get Actual Content (HTTP Mode)
			resp, err := client.Post(ctx, "/topic/lexical_elements/"+tt.chapterID)
			if err != nil {
				t.Fatalf("Failed to request: %v", err)
			}
			defer resp.Close()

			if resp.StatusCode != 200 {
				t.Errorf("Expected 200, got %d", resp.StatusCode)
			}

			// Parse JSON response which structure is {code:0, message:"OK", data: {title: "...", content: "..."}}
			var res struct {
				Code int `json:"code"`
				Data struct {
					Content string `json:"content"`
				} `json:"data"`
			}

			err = g.NewVar(resp.ReadAll()).Scan(&res)
			if err != nil {
				t.Fatalf("Failed to parse response: %v", err)
			}

			actual := res.Data.Content

			// Apply sanitizer if present
			if tt.sanitizer != nil {
				expected = tt.sanitizer(expected)
				actual = tt.sanitizer(actual)
			}

			// 3. Compare
			if actual != expected {
				t.Errorf("Content mismatch for %s.\nLength diff: %d vs %d",
					tt.chapterID, len(expected), len(actual))
			}
		})
	}
}
