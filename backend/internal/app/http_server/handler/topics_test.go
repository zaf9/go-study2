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

func TestGetTopics_JSON(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.Server("test-topics-json")
		s.SetPort(0)
		s.SetAccessLogEnabled(false)

		s.Group("/api/v1", func(group *ghttp.RouterGroup) {
			h := New()
			group.Middleware(middleware.Format)
			group.ALL("/topics", h.GetTopics)
		})

		s.Start()
		defer s.Shutdown()

		port := s.GetListenedPort()
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", port))

		resp, err := client.Get(nil, "/api/v1/topics?format=json")
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
		topics := dataMap["topics"].([]interface{})
		t.Assert(len(topics), 4)

		// 验证包含特定主题
		foundLexical := false
		foundConstants := false
		foundVariables := false
		foundTypes := false
		for _, topic := range topics {
			topicMap := topic.(map[string]interface{})
			if topicMap["id"] == "lexical_elements" {
				foundLexical = true
			}
			if topicMap["id"] == "constants" {
				foundConstants = true
			}
			if topicMap["id"] == "variables" {
				foundVariables = true
			}
			if topicMap["id"] == "types" {
				foundTypes = true
			}
		}
		t.Assert(foundLexical, true)
		t.Assert(foundConstants, true)
		t.Assert(foundVariables, true)
		t.Assert(foundTypes, true)
	})
}

func TestGetTopics_HTML(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.Server("test-topics-html")
		s.SetPort(0)
		s.SetAccessLogEnabled(false)

		s.Group("/api/v1", func(group *ghttp.RouterGroup) {
			h := New()
			group.Middleware(middleware.Format)
			group.ALL("/topics", h.GetTopics)
		})

		s.Start()
		defer s.Shutdown()

		port := s.GetListenedPort()
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", port))

		resp, err := client.Get(nil, "/api/v1/topics?format=html")
		t.AssertNil(err)
		defer resp.Close()
		t.Assert(resp.StatusCode, 200)

		body := resp.ReadAllString()
		t.AssertIN("text/html", resp.Header.Get("Content-Type"))
		t.AssertIN("Available Learning Topics", body)
		t.AssertIN("Lexical Elements", body)
		t.AssertIN("Constants", body)
		t.AssertIN("Variables", body)
		t.AssertIN("Types", body)
		t.AssertIN("<!DOCTYPE html>", body)
	})
}

func TestGetTopics_DefaultFormat(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.Server("test-topics-default")
		s.SetPort(0)
		s.SetAccessLogEnabled(false)

		s.Group("/api/v1", func(group *ghttp.RouterGroup) {
			h := New()
			group.Middleware(middleware.Format)
			group.ALL("/topics", h.GetTopics)
		})

		s.Start()
		defer s.Shutdown()

		port := s.GetListenedPort()
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", port))

		// 不提供 format 参数，应该默认为 JSON
		resp, err := client.Get(nil, "/api/v1/topics")
		t.AssertNil(err)
		defer resp.Close()
		t.Assert(resp.StatusCode, 200)

		var result Response
		err = json.Unmarshal(resp.ReadAll(), &result)
		t.AssertNil(err)
		t.Assert(result.Code, 0)
	})
}

func TestSendTopicsJSON(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.Server("test-send-topics-json")
		s.SetPort(0)
		s.SetAccessLogEnabled(false)

		h := New()
		var capturedResponse Response

		s.Group("/test", func(group *ghttp.RouterGroup) {
			group.GET("/topics", func(r *ghttp.Request) {
				topics := []Topic{
					{ID: "test1", Title: "Test 1", Description: "Description 1"},
					{ID: "test2", Title: "Test 2", Description: "Description 2"},
				}
				h.sendTopicsJSON(r, topics)
			})
		})

		s.Start()
		defer s.Shutdown()

		port := s.GetListenedPort()
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", port))

		resp, err := client.Get(nil, "/test/topics")
		t.AssertNil(err)
		defer resp.Close()

		err = json.Unmarshal(resp.ReadAll(), &capturedResponse)
		t.AssertNil(err)
		t.Assert(capturedResponse.Code, 0)
		t.Assert(capturedResponse.Message, "OK")
	})
}

func TestSendTopicsHTML(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.Server("test-send-topics-html")
		s.SetPort(0)
		s.SetAccessLogEnabled(false)

		h := New()

		s.Group("/test", func(group *ghttp.RouterGroup) {
			group.GET("/topics", func(r *ghttp.Request) {
				topics := []Topic{
					{ID: "test1", Title: "Test 1", Description: "Description 1"},
				}
				h.sendTopicsHTML(r, topics)
			})
		})

		s.Start()
		defer s.Shutdown()

		port := s.GetListenedPort()
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", port))

		resp, err := client.Get(nil, "/test/topics")
		t.AssertNil(err)
		defer resp.Close()

		body := resp.ReadAllString()
		t.AssertIN("<!DOCTYPE html>", body)
		t.AssertIN("Test 1", body)
		t.AssertIN("Description 1", body)
	})
}
