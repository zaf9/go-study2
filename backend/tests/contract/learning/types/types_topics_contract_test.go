package types_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"go-study2/internal/app/http_server/handler"
	"go-study2/internal/app/http_server/middleware"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/test/gtest"
)

// T040: 主题列表契约测试，确保包含 types。
func TestTopicsContract(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.Server("test-topics-contract")
		s.SetPort(0)
		s.SetAccessLogEnabled(false)

		s.Group("/api/v1", func(group *ghttp.RouterGroup) {
			h := handler.New()
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

		var result handler.Response
		err = json.Unmarshal(resp.ReadAll(), &result)
		t.AssertNil(err)
		t.Assert(result.Code, 20000)
		data := result.Data.(map[string]interface{})
		topics := data["topics"].([]interface{})
		foundTypes := false
		for _, tp := range topics {
			m := tp.(map[string]interface{})
			if m["id"] == "types" {
				foundTypes = true
				break
			}
		}
		t.Assert(foundTypes, true)

		resp2, err := client.Get(nil, "/api/v1/topics?format=html")
		t.AssertNil(err)
		defer resp2.Close()
		t.Assert(resp2.StatusCode, 200)
		body := resp2.ReadAllString()
		t.AssertIN("Types", body)
	})
}
