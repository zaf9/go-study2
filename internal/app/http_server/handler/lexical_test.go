package handler

import (
	"encoding/json"
	"fmt"
	"testing"

	"go-study2/internal/app/http_server/middleware"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/test/gtest"
)

func TestGetLexicalMenu_JSON(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.Server("test-lexical-menu-json")
		s.SetPort(0)
		s.SetAccessLogEnabled(false)

		s.Group("/api/v1", func(group *ghttp.RouterGroup) {
			h := New()
			group.Middleware(middleware.Format)
			group.ALL("/topic/lexical_elements", h.GetLexicalMenu)
		})

		s.Start()
		defer s.Shutdown()

		port := s.GetListenedPort()
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", port))

		resp, err := client.Get(nil, "/api/v1/topic/lexical_elements?format=json")
		t.AssertNil(err)
		defer resp.Close()
		t.Assert(resp.StatusCode, 200)

		var result Response
		err = json.Unmarshal(resp.ReadAll(), &result)
		t.AssertNil(err)
		t.Assert(result.Code, 0)
		t.Assert(result.Message, "OK")

		// 验证数据结构
		dataMap := result.Data.(map[string]interface{})
		items := dataMap["items"].([]interface{})
		t.Assert(len(items), 11)

		// 验证包含特定章节
		foundComments := false
		for _, item := range items {
			itemMap := item.(map[string]interface{})
			if itemMap["name"] == "comments" {
				foundComments = true
				break
			}
		}
		t.Assert(foundComments, true)
	})
}

func TestGetLexicalMenu_HTML(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.Server("test-lexical-menu-html")
		s.SetPort(0)
		s.SetAccessLogEnabled(false)

		s.Group("/api/v1", func(group *ghttp.RouterGroup) {
			h := New()
			group.Middleware(middleware.Format)
			group.ALL("/topic/lexical_elements", h.GetLexicalMenu)
		})

		s.Start()
		defer s.Shutdown()

		port := s.GetListenedPort()
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", port))

		resp, err := client.Get(nil, "/api/v1/topic/lexical_elements?format=html")
		t.AssertNil(err)
		defer resp.Close()
		t.Assert(resp.StatusCode, 200)

		body := resp.ReadAllString()
		t.AssertIN("Lexical Elements", body)
		t.AssertIN("<!DOCTYPE html>", body)
		t.AssertIN("Back to Topics", body)
	})
}

func TestGetLexicalContent_ValidChapter_JSON(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.Server("test-lexical-content-json")
		s.SetPort(0)
		s.SetAccessLogEnabled(false)

		s.Group("/api/v1", func(group *ghttp.RouterGroup) {
			h := New()
			group.Middleware(middleware.Format)
			group.ALL("/topic/lexical_elements/:chapter", h.GetLexicalContent)
		})

		s.Start()
		defer s.Shutdown()

		port := s.GetListenedPort()
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", port))

		resp, err := client.Get(nil, "/api/v1/topic/lexical_elements/comments?format=json")
		t.AssertNil(err)
		defer resp.Close()
		t.Assert(resp.StatusCode, 200)

		var result Response
		err = json.Unmarshal(resp.ReadAll(), &result)
		t.AssertNil(err)
		t.Assert(result.Code, 0)
		t.Assert(result.Message, "OK")

		// 验证内容数据
		dataMap := result.Data.(map[string]interface{})
		t.AssertNE(dataMap["title"], nil)
		t.AssertNE(dataMap["content"], nil)
	})
}

func TestGetLexicalContent_ValidChapter_HTML(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.Server("test-lexical-content-html")
		s.SetPort(0)
		s.SetAccessLogEnabled(false)

		s.Group("/api/v1", func(group *ghttp.RouterGroup) {
			h := New()
			group.Middleware(middleware.Format)
			group.ALL("/topic/lexical_elements/:chapter", h.GetLexicalContent)
		})

		s.Start()
		defer s.Shutdown()

		port := s.GetListenedPort()
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", port))

		resp, err := client.Get(nil, "/api/v1/topic/lexical_elements/comments?format=html")
		t.AssertNil(err)
		defer resp.Close()
		t.Assert(resp.StatusCode, 200)

		body := resp.ReadAllString()
		t.AssertIN("<!DOCTYPE html>", body)
		t.AssertIN("Back to Menu", body)
	})
}

func TestGetLexicalContent_InvalidChapter_JSON(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.Server("test-lexical-content-invalid-json")
		s.SetPort(0)
		s.SetAccessLogEnabled(false)

		s.Group("/api/v1", func(group *ghttp.RouterGroup) {
			h := New()
			group.Middleware(middleware.Format)
			group.ALL("/topic/lexical_elements/:chapter", h.GetLexicalContent)
		})

		s.Start()
		defer s.Shutdown()

		port := s.GetListenedPort()
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", port))

		resp, err := client.Get(nil, "/api/v1/topic/lexical_elements/invalid_chapter?format=json")
		t.AssertNil(err)
		defer resp.Close()
		t.Assert(resp.StatusCode, 404)

		body := resp.ReadAll()
		if len(body) > 0 {
			var result Response
			err = json.Unmarshal(body, &result)
			if err == nil {
				t.Assert(result.Code, 404)
				t.AssertIN("not found", result.Message)
			}
		}
	})
}

func TestGetLexicalContent_InvalidChapter_HTML(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.Server("test-lexical-content-invalid-html")
		s.SetPort(0)
		s.SetAccessLogEnabled(false)

		s.Group("/api/v1", func(group *ghttp.RouterGroup) {
			h := New()
			group.Middleware(middleware.Format)
			group.ALL("/topic/lexical_elements/:chapter", h.GetLexicalContent)
		})

		s.Start()
		defer s.Shutdown()

		port := s.GetListenedPort()
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", port))

		resp, err := client.Get(nil, "/api/v1/topic/lexical_elements/invalid_chapter?format=html")
		t.AssertNil(err)
		defer resp.Close()
		t.Assert(resp.StatusCode, 404)

		body := resp.ReadAllString()
		t.AssertIN("Chapter not found", body)
	})
}

func TestSendLexicalMenuJSON(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.Server("test-send-lexical-menu-json")
		s.SetPort(0)
		s.SetAccessLogEnabled(false)

		h := New()
		var capturedResponse Response

		s.Group("/test", func(group *ghttp.RouterGroup) {
			group.GET("/menu", func(r *ghttp.Request) {
				items := []LexicalMenuItem{
					{ID: 0, Title: "Test Chapter", Name: "test"},
				}
				h.sendLexicalMenuJSON(r, items)
			})
		})

		s.Start()
		defer s.Shutdown()

		port := s.GetListenedPort()
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", port))

		resp, err := client.Get(nil, "/test/menu")
		t.AssertNil(err)
		defer resp.Close()

		err = json.Unmarshal(resp.ReadAll(), &capturedResponse)
		t.AssertNil(err)
		t.Assert(capturedResponse.Code, 0)
		t.Assert(capturedResponse.Message, "OK")
	})
}

func TestSendLexicalMenuHTML(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.Server("test-send-lexical-menu-html")
		s.SetPort(0)
		s.SetAccessLogEnabled(false)

		h := New()

		s.Group("/test", func(group *ghttp.RouterGroup) {
			group.GET("/menu", func(r *ghttp.Request) {
				items := []LexicalMenuItem{
					{ID: 0, Title: "Test Chapter", Name: "test"},
				}
				h.sendLexicalMenuHTML(r, items)
			})
		})

		s.Start()
		defer s.Shutdown()

		port := s.GetListenedPort()
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", port))

		resp, err := client.Get(nil, "/test/menu")
		t.AssertNil(err)
		defer resp.Close()

		body := resp.ReadAllString()
		t.AssertIN("<!DOCTYPE html>", body)
		t.AssertIN("Test Chapter", body)
	})
}

