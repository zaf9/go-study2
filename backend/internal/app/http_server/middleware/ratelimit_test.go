package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/test/gtest"
)

func TestLoginRateLimit_Basic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// 创建专用的速率限制器（使用较小的配置便于测试）
		testConfig := RateLimiterConfig{
			MaxAttempts:    3,                      // 3 次尝试
			WindowDuration: 100 * time.Millisecond, // 100ms 窗口
			BlockDuration:  200 * time.Millisecond, // 200ms 封禁
		}
		testLimiter := NewRateLimiter(testConfig)

		// 保存全局限制器，测试后恢复
		originalLimiter := GlobalLoginRateLimiter
		GlobalLoginRateLimiter = testLimiter
		defer func() { GlobalLoginRateLimiter = originalLimiter }()

		s := g.Server("test-ratelimit-basic")
		s.SetPort(0)
		s.SetAccessLogEnabled(false)

		s.Group("/api/v1/auth", func(group *ghttp.RouterGroup) {
			group.Middleware(LoginRateLimit)
			group.POST("/login", func(r *ghttp.Request) {
				r.Response.WriteJson(g.Map{
					"code":    20000,
					"message": "登录成功",
				})
			})
		})

		s.Start()
		defer s.Shutdown()

		port := s.GetListenedPort()
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", port))

		// 测试 1: 前 3 次请求应该成功
		for i := 0; i < 3; i++ {
			resp, err := client.Post(context.Background(), "/api/v1/auth/login", `{"username":"test","password":"wrong"}`)
			t.AssertNil(err)
			defer resp.Close()

			t.Assert(resp.StatusCode, 200)

			// 验证响应体包含预期的内容
			body := resp.ReadAll()
			t.Assert(len(body) > 0, true)

			var result map[string]interface{}
			err = json.Unmarshal(body, &result)
			t.AssertNil(err)

			// 验证响应码字段存在且正确
			codeVal, exists := result["code"]
			t.Assert(exists, true)
			t.AssertNE(codeVal, nil)
			t.Assert(codeVal, float64(20000))
		}

		// 测试 2: 第 4 次请求应该被限流
		resp, err := client.Post(context.Background(), "/api/v1/auth/login", `{"username":"test","password":"wrong"}`)
		t.AssertNil(err)
		defer resp.Close()

		t.Assert(resp.StatusCode, 429)

		// 读取响应体并解析 JSON
		body := resp.ReadAll()
		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		t.AssertNil(err)

		// 验证响应包含正确的错误码和消息
		t.AssertNE(result["code"], nil)
		t.Assert(strings.Contains(fmt.Sprint(result["message"]), "登录尝试过多"), true)

		// 测试 3: 等待封禁期结束，再次请求应该成功
		time.Sleep(testConfig.BlockDuration + 50*time.Millisecond)
		resp2, err := client.Post(context.Background(), "/api/v1/auth/login", `{"username":"test","password":"wrong"}`)
		t.AssertNil(err)
		defer resp2.Close()

		t.Assert(resp2.StatusCode, 200)
		body2 := resp2.ReadAll()
		var result2 map[string]interface{}
		err = json.Unmarshal(body2, &result2)
		t.AssertNil(err)
		t.Assert(result2["code"], float64(20000))
	})
}

func TestLoginRateLimit_DifferentPaths(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.Server("test-ratelimit-paths")
		s.SetPort(0)
		s.SetAccessLogEnabled(false)

		s.Group("/api/v1", func(group *ghttp.RouterGroup) {
			// 只对登录路径应用限流
			group.POST("/auth/login", LoginRateLimit, func(r *ghttp.Request) {
				r.Response.WriteJson(g.Map{"code": 20000})
			})
			// 其他路径不限流
			group.POST("/auth/refresh", func(r *ghttp.Request) {
				r.Response.WriteJson(g.Map{"code": 20000})
			})
			group.GET("/topics", func(r *ghttp.Request) {
				r.Response.WriteJson(g.Map{"code": 20000})
			})
		})

		s.Start()
		defer s.Shutdown()

		port := s.GetListenedPort()
		baseURL := fmt.Sprintf("http://127.0.0.1:%d", port)

		// 测试：多次请求非登录路径不应该被限流
		for i := 0; i < 10; i++ {
			resp, err := http.Post(baseURL+"/api/v1/auth/refresh", "application/json", strings.NewReader("{}"))
			t.AssertNil(err)
			defer resp.Body.Close()
			t.Assert(resp.StatusCode, 200)
		}

		// GET 请求也不应该被限流
		resp, err := http.Get(baseURL + "/api/v1/topics")
		t.AssertNil(err)
		defer resp.Body.Close()
		t.Assert(resp.StatusCode, 200)
	})
}

func TestRateLimiter_Allow(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		config := RateLimiterConfig{
			MaxAttempts:    2,
			WindowDuration: 100 * time.Millisecond,
			BlockDuration:  200 * time.Millisecond,
		}
		limiter := NewRateLimiter(config)

		ip := "192.168.1.100"

		// 测试 1: 首次请求应该允许
		allowed, _ := limiter.Allow(ip)
		t.Assert(allowed, true)

		// 测试 2: 第二次请求应该允许
		allowed, _ = limiter.Allow(ip)
		t.Assert(allowed, true)

		// 测试 3: 第三次请求应该被拒绝
		allowed, retryAfter := limiter.Allow(ip)
		t.Assert(allowed, false)
		t.Assert(retryAfter > 0, true)

		// 测试 4: 等待窗口过期后，应该重置
		time.Sleep(config.WindowDuration + 50*time.Millisecond)
		allowed, _ = limiter.Allow(ip)
		t.Assert(allowed, true)
	})
}

func TestRateLimiter_MultipleIPs(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		config := RateLimiterConfig{
			MaxAttempts:    2,
			WindowDuration: 100 * time.Millisecond,
			BlockDuration:  200 * time.Millisecond,
		}
		limiter := NewRateLimiter(config)

		ip1 := "192.168.1.100"
		ip2 := "192.168.1.101"

		// IP1 请求 2 次
		limiter.Allow(ip1)
		limiter.Allow(ip1)

		// IP1 第 3 次应该被拒绝
		allowed, _ := limiter.Allow(ip1)
		t.Assert(allowed, false)

		// IP2 应该不受 IP1 影响
		allowed, _ = limiter.Allow(ip2)
		t.Assert(allowed, true)

		allowed, _ = limiter.Allow(ip2)
		t.Assert(allowed, true)

		allowed, _ = limiter.Allow(ip2)
		t.Assert(allowed, false)
	})
}

func TestRateLimiter_Reset(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		config := RateLimiterConfig{
			MaxAttempts:    2,
			WindowDuration: 50 * time.Millisecond,
			BlockDuration:  100 * time.Millisecond,
		}
		limiter := NewRateLimiter(config)

		ip := "192.168.1.200"

		// 触发限流
		limiter.Allow(ip)
		limiter.Allow(ip)
		limiter.Allow(ip) // 被封禁

		// 验证被封禁
		allowed, _ := limiter.Allow(ip)
		t.Assert(allowed, false)

		// 等待封禁期和窗口期都过期
		time.Sleep(200 * time.Millisecond)

		// 调用 Reset 清理过期记录
		limiter.Reset()

		// 再次请求应该从零开始
		allowed, _ = limiter.Allow(ip)
		t.Assert(allowed, true)
	})
}

func TestRateLimiter_Stats(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		config := RateLimiterConfig{
			MaxAttempts:    2,
			WindowDuration: 100 * time.Millisecond,
			BlockDuration:  200 * time.Millisecond,
		}
		limiter := NewRateLimiter(config)

		// 添加一些活动
		limiter.Allow("192.168.1.1")
		limiter.Allow("192.168.1.1")
		limiter.Allow("192.168.1.1") // 触发封禁

		limiter.Allow("192.168.1.2")

		stats := limiter.Stats()

		t.Assert(stats["total_ips"], 2)
		t.Assert(stats["max_attempts"], 2)
		t.Assert(stats["window_duration"], "100ms")
		t.Assert(stats["block_duration"], "200ms")
	})
}

func TestLoginRateLimitWithProxyHeaders(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.Server("test-ratelimit-proxy")
		s.SetPort(0)
		s.SetAccessLogEnabled(false)

		s.Group("/api/v1/auth", func(group *ghttp.RouterGroup) {
			group.Middleware(LoginRateLimit)
			group.POST("/login", func(r *ghttp.Request) {
				r.Response.WriteJson(g.Map{"code": 20000})
			})
		})

		s.Start()
		defer s.Shutdown()

		port := s.GetListenedPort()
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", port))

		// 从同一个 IP 发送多次请求（模拟通过代理）
		ip := "203.0.113.195" // 测试用 IP

		for i := 0; i < 6; i++ {
			// GoFrame 客户端设置自定义头
			resp, err := client.Post(context.Background(), "/api/v1/auth/login",
				`{"username":"test","password":"wrong"}`,
				http.Header{
					"X-Forwarded-For": []string{ip},
					"X-Real-IP":      []string{ip},
				})

			t.AssertNil(err)
			defer resp.Close()

			if i < 5 {
				t.Assert(resp.StatusCode, 200)
			} else {
				t.Assert(resp.StatusCode, 429)
			}
		}
	})
}
