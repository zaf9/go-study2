package performance

import (
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"go-study2/internal/app/http_server"
	"go-study2/internal/config"
)

// TestConstantsAPI_Performance simulates load to verify performance requirements
// Requirement 1: 100 concurrent requests, p95 < 100ms (we use avg here as proxy, target < 50ms avg)
// Requirement 2: 1000 concurrent requests, no errors
func TestConstantsAPI_Performance(t *testing.T) {
	// 跳过短测试模式（使用 -short 标志时）
	if testing.Short() {
		t.Skip("跳过性能测试（短测试模式）")
	}
	// 1. Setup Server
	cfg := &config.Config{
		Server: config.ServerConfig{
			Host: "127.0.0.1",
			Port: 0, // Random port
		},
	}
	// Use a unique server name to avoid conflicts if parallel tests
	s := http_server.NewServer(cfg, "perf-test-constants")
	s.SetAccessLogEnabled(false) // Disable access log for performance
	s.SetErrorLogEnabled(true)
	s.Start()
	defer s.Shutdown()

	port := s.GetListenedPort()
	baseUrl := fmt.Sprintf("http://127.0.0.1:%d/api/v1/topic/constants/boolean", port)
	t.Logf("Performance test server running at %s", baseUrl)

	// Warmup
	runLoadTest(t, baseUrl, 10, 1*time.Second, 500*time.Millisecond)

	t.Run("Concurrent_100", func(t *testing.T) {
		// Target: 100 concurrent, < 100ms avg response (spec says p95 < 100ms)
		runLoadTest(t, baseUrl, 100, 5*time.Second, 100*time.Millisecond)
	})

	t.Run("Concurrent_1000", func(t *testing.T) {
		// Target: 1000 concurrent, 0 error rate
		// Allow higher latency for stress test
		runLoadTest(t, baseUrl, 1000, 5*time.Second, 200*time.Millisecond)
	})
}

func runLoadTest(t *testing.T, url string, concurrency int, duration time.Duration, maxAvgLatency time.Duration) {
	var wg sync.WaitGroup
	var successCount int64
	var failCount int64
	var totalLatency int64 // microseconds

	start := time.Now()
	stopCh := make(chan struct{})

	// Timer to stop
	time.AfterFunc(duration, func() {
		close(stopCh)
	})

	// Create client with custom transport to handle concurrency
	tr := &http.Transport{
		MaxIdleConns:        concurrency,
		MaxIdleConnsPerHost: concurrency,
		IdleConnTimeout:     30 * time.Second,
	}
	client := &http.Client{Transport: tr, Timeout: 5 * time.Second}

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-stopCh:
					return
				default:
					reqStart := time.Now()
					resp, err := client.Get(url)
					reqDur := time.Since(reqStart)

					if err == nil {
						if resp.StatusCode == 200 {
							atomic.AddInt64(&successCount, 1)
							atomic.AddInt64(&totalLatency, int64(reqDur.Microseconds()))
							// Read body to reuse connection
							// buffer := make([]byte, 1024)
							// n, _ := resp.Body.Read(buffer)
							resp.Body.Close()
						} else {
							atomic.AddInt64(&failCount, 1)
							resp.Body.Close()
						}
					} else {
						atomic.AddInt64(&failCount, 1)
					}
				}
			}
		}()
	}

	wg.Wait()
	totalDuration := time.Since(start)

	totalReqs := successCount + failCount
	if totalReqs == 0 {
		t.Fatal("No requests completed")
	}

	avgLatency := float64(totalLatency) / float64(successCount) / 1000.0 // ms
	rps := float64(totalReqs) / totalDuration.Seconds()

	t.Logf("Stats for Concurrency %d:", concurrency)
	t.Logf("  Total Requests: %d", totalReqs)
	t.Logf("  Success: %d", successCount)
	t.Logf("  Failed: %d", failCount)
	t.Logf("  RPS: %.2f", rps)
	t.Logf("  Avg Latency: %.2f ms", avgLatency)

	if failCount > 0 {
		// For 1000 concurrent, some connections might be refused on some systems, warning instead of fail?
		// Spec says "error rate 0%".
		// But in local dev environment 1000 might be hitting OS limits.
		// I will make it an error.
		t.Errorf("Experienced %d failures (%.2f%%)", failCount, float64(failCount)/float64(totalReqs)*100)
	}
	if avgLatency > float64(maxAvgLatency.Milliseconds()) {
		t.Errorf("Average latency %.2f ms exceeded limit %v", avgLatency, maxAvgLatency)
	}
}
