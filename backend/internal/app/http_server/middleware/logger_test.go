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

func TestLogger(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.Server("test-logger")
		s.SetPort(0)
		s.SetAccessLogEnabled(false)

		// 注册测试路由和中间件
		s.Group("/test", func(group *ghttp.RouterGroup) {
			group.Middleware(Logger)
			group.GET("/hello", func(r *ghttp.Request) {
				r.Response.Write("Hello World")
			})
		})

		s.Start()
		defer s.Shutdown()

		port := s.GetListenedPort()
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", port))

		// 发送请求
		startTime := time.Now()
		resp, err := client.Get(context.TODO(), "/test/hello")
		duration := time.Since(startTime)

		t.AssertNil(err)
		defer resp.Close()
		t.Assert(resp.StatusCode, 200)
		t.Assert(resp.ReadAllString(), "Hello World")

		// 验证请求已处理（中间件会记录日志，但这里主要验证功能正常）
		// 由于日志是异步的，我们主要验证请求成功处理
		t.Assert(duration < time.Second, true)
	})
}

func TestLoggerWithError(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.Server("test-logger-error")
		s.SetPort(0)
		s.SetAccessLogEnabled(false)

		s.Group("/test", func(group *ghttp.RouterGroup) {
			group.Middleware(Logger)
			group.GET("/error", func(r *ghttp.Request) {
				r.Response.WriteStatus(500)
				r.Response.Write("Internal Server Error")
			})
		})

		s.Start()
		defer s.Shutdown()

		port := s.GetListenedPort()
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", port))

		resp, err := client.Get(context.TODO(), "/test/error")
		t.AssertNil(err)
		defer resp.Close()
		t.Assert(resp.StatusCode, 500)
	})
}
