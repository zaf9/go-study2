package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/test/gtest"
)

func TestFormat_DefaultJSON(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.Server("test-format-default")
		s.SetPort(0)
		s.SetAccessLogEnabled(false)

		s.Group("/test", func(group *ghttp.RouterGroup) {
			group.Middleware(Format)
			group.GET("/data", func(r *ghttp.Request) {
				format := r.GetCtxVar("format").String()
				r.Response.WriteJson(map[string]interface{}{
					"format": format,
				})
			})
		})

		s.Start()
		defer s.Shutdown()

		port := s.GetListenedPort()
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", port))

		// 测试默认格式（不提供 format 参数）
		resp, err := client.Get(context.TODO(), "/test/data")
		t.AssertNil(err)
		defer resp.Close()
		t.Assert(resp.StatusCode, 200)

		var result map[string]interface{}
		err = json.Unmarshal(resp.ReadAll(), &result)
		t.AssertNil(err)
		t.Assert(result["format"], "json")
	})
}

func TestFormat_ExplicitJSON(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.Server("test-format-json")
		s.SetPort(0)
		s.SetAccessLogEnabled(false)

		s.Group("/test", func(group *ghttp.RouterGroup) {
			group.Middleware(Format)
			group.GET("/data", func(r *ghttp.Request) {
				format := r.GetCtxVar("format").String()
				r.Response.WriteJson(map[string]interface{}{
					"format": format,
				})
			})
		})

		s.Start()
		defer s.Shutdown()

		port := s.GetListenedPort()
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", port))

		// 测试显式指定 JSON 格式
		resp, err := client.Get(context.TODO(), "/test/data?format=json")
		t.AssertNil(err)
		defer resp.Close()
		t.Assert(resp.StatusCode, 200)

		var result map[string]interface{}
		err = json.Unmarshal(resp.ReadAll(), &result)
		t.AssertNil(err)
		t.Assert(result["format"], "json")
	})
}

func TestFormat_HTML(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.Server("test-format-html")
		s.SetPort(0)
		s.SetAccessLogEnabled(false)

		s.Group("/test", func(group *ghttp.RouterGroup) {
			group.Middleware(Format)
			group.GET("/data", func(r *ghttp.Request) {
				format := r.GetCtxVar("format").String()
				r.Response.Write(format)
			})
		})

		s.Start()
		defer s.Shutdown()

		port := s.GetListenedPort()
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", port))

		// 测试 HTML 格式
		resp, err := client.Get(context.TODO(), "/test/data?format=html")
		t.AssertNil(err)
		defer resp.Close()
		t.Assert(resp.StatusCode, 200)
		t.Assert(resp.ReadAllString(), "html")
	})
}

func TestFormat_InvalidFormat(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.Server("test-format-invalid")
		s.SetPort(0)
		s.SetAccessLogEnabled(false)

		s.Group("/test", func(group *ghttp.RouterGroup) {
			group.Middleware(Format)
			group.GET("/data", func(r *ghttp.Request) {
				r.Response.Write("should not reach here")
			})
		})

		s.Start()
		defer s.Shutdown()

		port := s.GetListenedPort()
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", port))

		// 测试无效格式
		resp, err := client.Get(context.TODO(), "/test/data?format=xml")
		t.AssertNil(err)
		defer resp.Close()
		t.Assert(resp.StatusCode, 400)

		body := resp.ReadAll()
		if len(body) > 0 {
			var result map[string]interface{}
			err = json.Unmarshal(body, &result)
			if err == nil {
				t.Assert(result["code"], 400)
				if msg, ok := result["message"].(string); ok {
					t.AssertIN("Invalid format", msg)
				}
			}
		}
	})
}

func TestFormat_EmptyFormat(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.Server("test-format-empty")
		s.SetPort(0)
		s.SetAccessLogEnabled(false)

		s.Group("/test", func(group *ghttp.RouterGroup) {
			group.Middleware(Format)
			group.GET("/data", func(r *ghttp.Request) {
				format := r.GetCtxVar("format").String()
				r.Response.WriteJson(map[string]interface{}{
					"format": format,
				})
			})
		})

		s.Start()
		defer s.Shutdown()

		port := s.GetListenedPort()
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", port))

		// 测试空格式参数（应该默认为 json）
		resp, err := client.Get(context.TODO(), "/test/data?format=")
		t.AssertNil(err)
		defer resp.Close()
		t.Assert(resp.StatusCode, 200)

		var result map[string]interface{}
		err = json.Unmarshal(resp.ReadAll(), &result)
		t.AssertNil(err)
		t.Assert(result["format"], "json")
	})
}
