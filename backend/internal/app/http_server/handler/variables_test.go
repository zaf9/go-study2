package handler

import (
	"fmt"
	"testing"

	"go-study2/internal/app/http_server/middleware"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/test/gtest"
)

func TestVariablesMenuAndContent(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.Server("test-variables")
		s.SetPort(0)
		s.SetAccessLogEnabled(false)

		h := New()
		s.Group("/api/v1", func(group *ghttp.RouterGroup) {
			group.Middleware(middleware.Format)
			group.ALL("/topic/variables", h.GetVariablesMenu)
			group.ALL("/topic/variables/:subtopic", h.GetVariableContent)
		})

		s.Start()
		defer s.Shutdown()

		port := s.GetListenedPort()
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", port))

		// Menu JSON
		resp, err := client.Get(nil, "/api/v1/topic/variables?format=json")
		t.AssertNil(err)
		defer resp.Close()
		t.Assert(resp.StatusCode, 200)

		// Content HTML
		resp, err = client.Get(nil, "/api/v1/topic/variables/storage?format=html")
		t.AssertNil(err)
		defer resp.Close()
		t.Assert(resp.StatusCode, 200)
		body := resp.ReadAllString()
		t.AssertIN("(storage)", body)
		t.AssertIN("零值", body)
	})
}
