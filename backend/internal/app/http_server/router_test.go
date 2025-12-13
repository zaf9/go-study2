package http_server

import (
	"context"
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
)

func TestRegisterRoutes(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.Server("test-router")
		s.SetPort(0)
		s.SetAccessLogEnabled(false)

		// 注册路由
		RegisterRoutes(s)

		// 启动服务器
		s.Start()
		defer s.Shutdown()

		port := s.GetListenedPort()
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", port))

		// 测试 topics 路由
		resp, err := client.Get(context.TODO(), "/api/v1/topics?format=json")
		t.AssertNil(err)
		defer resp.Close()
		t.Assert(resp.StatusCode, 200)

		// 测试 lexical_elements 菜单路由
		resp2, err := client.Get(context.TODO(), "/api/v1/topic/lexical_elements?format=json")
		t.AssertNil(err)
		defer resp2.Close()
		t.Assert(resp2.StatusCode, 200)

		// 测试 lexical_elements 内容路由
		resp3, err := client.Get(context.TODO(), "/api/v1/topic/lexical_elements/comments?format=json")
		t.AssertNil(err)
		defer resp3.Close()
		t.Assert(resp3.StatusCode, 200)

		// 测试 constants 菜单路由
		resp4, err := client.Get(context.TODO(), "/api/v1/topic/constants?format=json")
		t.AssertNil(err)
		defer resp4.Close()
		t.Assert(resp4.StatusCode, 200)

		// 测试 constants 内容路由
		resp5, err := client.Get(context.TODO(), "/api/v1/topic/constants/boolean?format=json")
		t.AssertNil(err)
		defer resp5.Close()
		t.Assert(resp5.StatusCode, 200)
	})
}

func TestRegisterRoutesWithMiddleware(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.Server("test-router-middleware")
		s.SetPort(0)
		s.SetAccessLogEnabled(false)

		RegisterRoutes(s)
		s.Start()
		defer s.Shutdown()

		port := s.GetListenedPort()
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", port))

		// 测试格式中间件：无效格式
		resp, err := client.Get(context.TODO(), "/api/v1/topics?format=invalid")
		t.AssertNil(err)
		defer resp.Close()
		t.Assert(resp.StatusCode, 400)

		// 测试格式中间件：HTML 格式
		resp2, err := client.Get(context.TODO(), "/api/v1/topics?format=html")
		t.AssertNil(err)
		defer resp2.Close()
		t.Assert(resp2.StatusCode, 200)
		t.AssertIN("text/html", resp2.Header.Get("Content-Type"))
	})
}
