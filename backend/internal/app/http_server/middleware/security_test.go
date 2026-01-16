package middleware

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/test/gtest"
)

func TestSecurity(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.Server("test-security")
		s.SetPort(0)
		s.SetAccessLogEnabled(false)

		s.Group("/test", func(group *ghttp.RouterGroup) {
			group.Middleware(Security)
			group.GET("/hello", func(r *ghttp.Request) {
				r.Response.Write("Hello World")
			})
		})

		s.Start()
		defer s.Shutdown()

		port := s.GetListenedPort()
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", port))

		resp, err := client.Get(context.TODO(), "/test/hello")
		t.AssertNil(err)
		defer resp.Close()

		// 验证响应成功
		t.Assert(resp.StatusCode, 200)

		// 验证安全头存在（非 TLS 模式下 HSTS 不会设置）
		headers := resp.Header

		// CSP 头应该存在
		csp := headers.Get("Content-Security-Policy")
		t.Assert(csp != "", true)

		// X-Frame-Options 应该是 DENY
		frameOptions := headers.Get("X-Frame-Options")
		t.Assert(frameOptions, "DENY")

		// X-Content-Type-Options 应该是 nosniff
		ctOptions := headers.Get("X-Content-Type-Options")
		t.Assert(ctOptions, "nosniff")

		// X-XSS-Protection 应该存在
		xss := headers.Get("X-XSS-Protection")
		t.Assert(xss != "", true)

		// Referrer-Policy 应该存在
		referrer := headers.Get("Referrer-Policy")
		t.Assert(referrer != "", true)

		// Permissions-Policy 应该存在
		permissions := headers.Get("Permissions-Policy")
		t.Assert(permissions != "", true)
	})
}

func TestSecurityWithHTTPS(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.Server("test-security-https")
		s.SetPort(0)
		s.SetAccessLogEnabled(false)

		s.Group("/test", func(group *ghttp.RouterGroup) {
			group.Middleware(Security)
			group.GET("/hello", func(r *ghttp.Request) {
				r.Response.Write("Hello World")
			})
		})

		s.Start()
		defer s.Shutdown()

		port := s.GetListenedPort()
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", port))

		resp, err := client.Get(context.TODO(), "/test/hello")
		t.AssertNil(err)
		defer resp.Close()

		// 验证响应成功
		t.Assert(resp.StatusCode, 200)

		// 注意：由于测试环境使用 HTTP 而非 HTTPS，
		// HSTS 头不会设置（Security 中间件检查 r.TLS != nil）
		headers := resp.Header
		hsts := headers.Get("Strict-Transport-Security")
		t.Assert(hsts, "") // HTTP 模式下不应该有 HSTS
	})
}

func TestSecurityMiddlewareOrder(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.Server("test-security-order")
		s.SetPort(0)
		s.SetAccessLogEnabled(false)

		executedOrder := []string{}
		s.Group("/test", func(group *ghttp.RouterGroup) {
			group.Middleware(Security)
			group.Middleware(func(r *ghttp.Request) {
				executedOrder = append(executedOrder, "after-security")
				r.Middleware.Next()
			})
			group.GET("/check", func(r *ghttp.Request) {
				executedOrder = append(executedOrder, "handler")
				r.Response.Write("OK")
			})
		})

		s.Start()
		defer s.Shutdown()

		port := s.GetListenedPort()
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", port))

		startTime := time.Now()
		resp, err := client.Get(context.TODO(), "/test/check")
		duration := time.Since(startTime)

		t.AssertNil(err)
		defer resp.Close()
		t.Assert(resp.StatusCode, 200)

		// 验证中间件执行顺序正确
		t.Assert(len(executedOrder), 2)
		t.Assert(executedOrder[0], "after-security")
		t.Assert(executedOrder[1], "handler")

		// 验证性能：安全头中间件不应该显著增加延迟
		t.Assert(duration < 100*time.Millisecond, true)
	})
}

func TestSecurityHeadersValues(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.Server("test-security-values")
		s.SetPort(0)
		s.SetAccessLogEnabled(false)

		s.Group("/test", func(group *ghttp.RouterGroup) {
			group.Middleware(Security)
			group.GET("/headers", func(r *ghttp.Request) {
				r.Response.Write("OK")
			})
		})

		s.Start()
		defer s.Shutdown()

		port := s.GetListenedPort()
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", port))

		resp, err := client.Get(context.TODO(), "/test/headers")
		t.AssertNil(err)
		defer resp.Close()

		headers := resp.Header

		// 验证 X-XSS-Protection 的具体值
		xssProtection := headers.Get("X-XSS-Protection")
		t.Assert(xssProtection, "1; mode=block")

		// 验证 Referrer-Policy 的具体值
		referrerPolicy := headers.Get("Referrer-Policy")
		t.Assert(referrerPolicy, "strict-origin-when-cross-origin")

		// 验证 Permissions-Policy 包含必要的限制
		permissionsPolicy := headers.Get("Permissions-Policy")
		t.Assert(permissionsPolicy != "", true)
		// 应该禁用地理位置、摄像头、麦克风
		t.Assert(permissionsPolicy, "geolocation=(), microphone=(), camera=(), payment=()")

		// 验证 CSP 包含 default-src 'self'
		csp := headers.Get("Content-Security-Policy")
		t.Assert(csp != "", true)
		t.Assert(cspContains(csp, "default-src 'self'"), true)
		t.Assert(cspContains(csp, "frame-ancestors 'none'"), true)
	})
}

// cspContains 检查 CSP 是否包含指定指令
func cspContains(csp, directive string) bool {
	// 简单的字符串包含检查
	// 实际生产中可能需要更复杂的解析
	return len(csp) > 0 && contains(csp, directive)
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && findSubstring(s, substr) >= 0
}

func findSubstring(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
