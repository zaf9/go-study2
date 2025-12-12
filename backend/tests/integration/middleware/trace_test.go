package middleware_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"go-study2/internal/infrastructure/logger"

	"go-study2/internal/app/http_server/middleware"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func TestHTTPTraceIntegration(t *testing.T) {
	// Reset logger for test isolation
	logger.Reset()

	// Create temporary directory for logs. Use os.MkdirTemp instead of t.TempDir
	// so we can control deletion timing (call RemoveAll after logger reset) and
	// avoid race between logger releasing file handles and the testing framework
	// attempting to remove the temp dir.
	dir, err := os.MkdirTemp("", "TestHTTPTraceIntegration")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Setup logger configuration
	cfg := &logger.LoggerConfig{
		Level:  "info",
		Stdout: false,
		Instances: map[string]logger.InstanceConfig{
			"access": {
				Path:   dir,
				File:   "access.log",
				Level:  "info",
				Format: "text",
			},
			"app": {
				Path:   dir,
				File:   "app.log",
				Level:  "info",
				Format: "text",
			},
		},
	}

	err = logger.Initialize(cfg)
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}
	// Ensure logger resources (file handles, goroutines) are cleaned up after test
	// Add a small delay after Reset to give the OS time to release file handles on Windows.
	defer func() {
		// Reset logger and give OS a moment to release file handles, then
		// remove the temporary directory we created above.
		logger.Reset()
		time.Sleep(100 * time.Millisecond)
		_ = os.RemoveAll(dir)
	}()

	// Create test server
	s := ghttp.GetServer()
	s.BindHandler("/test", func(r *ghttp.Request) {
		// Simulate some processing
		time.Sleep(10 * time.Millisecond)
		r.Response.WriteJson(g.Map{
			"message":  "test response",
			"trace_id": logger.ExtractTraceID(r.Context()),
		})
	})

	// Add access log middleware
	s.Use(middleware.AccessLog)

	// Start server on random port
	s.SetPort(0)
	s.Start()
	defer s.Shutdown()

	// Wait for server to start
	time.Sleep(100 * time.Millisecond)

	// Get server address
	addr := s.GetListenedAddress()
	if addr == "" {
		t.Fatal("Server did not start properly")
	}

	// Test 1: Request without trace ID header
	t.Run("RequestWithoutTraceID", func(t *testing.T) {
		resp, err := http.Get("http://" + addr + "/test")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response: %v", err)
		}

		// Check response contains trace ID
		if !bytes.Contains(body, []byte("trace_id")) {
			t.Error("Response should contain trace_id")
		}

		// Wait for logs to be written and locate the actual access.log file (it may be in a subdir)
		accessLogPath, err := findAccessLogFile(dir, 500*time.Millisecond)
		if err != nil {
			t.Fatalf("Failed to locate access log: %v", err)
		}
		content, err := os.ReadFile(accessLogPath)
		if err != nil {
			t.Fatalf("Failed to read access log: %v", err)
		}

		logContent := string(content)
		if !strings.Contains(logContent, "[START]") {
			t.Error("Access log should contain request start")
		}
		if !strings.Contains(logContent, "/test") {
			t.Error("Access log should contain request path")
		}
	})

	// Test 2: Request with trace ID header
	t.Run("RequestWithTraceID", func(t *testing.T) {
		req, err := http.NewRequest("GET", "http://"+addr+"/test", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		testTraceID := "test-trace-id-12345"
		req.Header.Set("X-Trace-Id", testTraceID)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response: %v", err)
		}

		// Check response contains the injected trace ID
		if !bytes.Contains(body, []byte(testTraceID)) {
			t.Errorf("Response should contain injected trace ID %s", testTraceID)
		}

		// Wait for logs to be written and locate the actual access.log file
		accessLogPath, err := findAccessLogFile(dir, 500*time.Millisecond)
		if err != nil {
			t.Fatalf("Failed to locate access log: %v", err)
		}
		content, err := os.ReadFile(accessLogPath)
		if err != nil {
			t.Fatalf("Failed to read access log: %v", err)
		}

		logContent := string(content)
		// The trace ID should be in the context used for logging
		// Since we inject it into the request context, it should appear in logs
		if !strings.Contains(logContent, testTraceID) {
			t.Errorf("Access log should contain trace ID %s", testTraceID)
		}
	})

	// Test 3: Multiple requests should have different trace IDs
	t.Run("MultipleRequestsDifferentTraceIDs", func(t *testing.T) {
		// Make first request
		resp1, err := http.Get("http://" + addr + "/test")
		if err != nil {
			t.Fatalf("Failed to make first request: %v", err)
		}
		defer resp1.Body.Close()

		body1, err := io.ReadAll(resp1.Body)
		if err != nil {
			t.Fatalf("Failed to read first response: %v", err)
		}

		// Make second request
		resp2, err := http.Get("http://" + addr + "/test")
		if err != nil {
			t.Fatalf("Failed to make second request: %v", err)
		}
		defer resp2.Body.Close()

		body2, err := io.ReadAll(resp2.Body)
		if err != nil {
			t.Fatalf("Failed to read second response: %v", err)
		}

		// Extract trace IDs from responses
		traceID1 := extractTraceIDFromResponse(string(body1))
		traceID2 := extractTraceIDFromResponse(string(body2))

		if traceID1 == "" || traceID2 == "" {
			t.Error("Both responses should contain trace IDs")
		}

		if traceID1 == traceID2 {
			t.Error("Different requests should have different trace IDs")
		}
	})

	// Test TraceID propagation through context
	t.Run("TraceIDPropagation", func(t *testing.T) {
		// Create a request with trace ID
		testTraceID := "propagation-test-123"
		req, err := http.NewRequest("GET", "http://"+addr+"/test", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("X-Trace-Id", testTraceID)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		// The middleware should have processed the request and injected trace ID
		// We can't directly test context propagation in integration test
		// but we can verify the request was processed
		if resp.StatusCode != 200 {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})
}

// extractTraceIDFromResponse extracts trace_id from JSON response
func extractTraceIDFromResponse(response string) string {
	// Simple extraction for test purposes
	start := strings.Index(response, `"trace_id":"`)
	if start == -1 {
		return ""
	}
	start += len(`"trace_id":"`)
	end := strings.Index(response[start:], `"`)
	if end == -1 {
		return ""
	}
	return response[start : start+end]
}

// findAccessLogFile searches recursively under dir for a file named access.log.
// It retries until timeout to handle log write races.
func findAccessLogFile(dir string, timeout time.Duration) (string, error) {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		var found string
		_ = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			if !d.IsDir() && strings.EqualFold(d.Name(), "access.log") {
				found = path
				return fmt.Errorf("found")
			}
			return nil
		})
		if found != "" {
			return found, nil
		}
		time.Sleep(50 * time.Millisecond)
	}
	return "", fmt.Errorf("access.log not found under %s", dir)
}

func TestTraceIDRecovery(t *testing.T) {
	// Test that trace ID is recovered when interrupted
	logger.Reset()

	dir := t.TempDir()
	cfg := &logger.LoggerConfig{
		Level:  "info",
		Stdout: false,
		Instances: map[string]logger.InstanceConfig{
			"access": {
				Path:   dir,
				File:   "access.log",
				Level:  "info",
				Format: "text",
			},
		},
	}

	err := logger.Initialize(cfg)
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}

	// Test EnsureTraceID recovery
	ctx := context.Background()

	// Initially no trace ID
	traceID1 := logger.ExtractTraceID(ctx)
	if traceID1 != "" {
		t.Error("Initial context should not have trace ID")
	}

	// Ensure trace ID generates one
	ctx = logger.EnsureTraceID(ctx)
	traceID2 := logger.ExtractTraceID(ctx)
	if traceID2 == "" {
		t.Error("EnsureTraceID should generate trace ID")
	}

	// Ensure again should preserve existing
	ctx = logger.EnsureTraceID(ctx)
	traceID3 := logger.ExtractTraceID(ctx)
	if traceID3 != traceID2 {
		t.Error("EnsureTraceID should preserve existing trace ID")
	}
}
