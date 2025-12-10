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

// T008: 菜单与内容契约（JSON/HTML）
func TestTypesMenuAndContentContract(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.Server("test-types-contract")
		s.SetPort(0)
		s.SetAccessLogEnabled(false)

		s.Group("/api/v1", func(group *ghttp.RouterGroup) {
			h := handler.New()
			group.Middleware(middleware.Format)
			group.ALL("/topic/types", h.GetTypesMenu)
			group.ALL("/topic/types/:subtopic", h.GetTypesContent)
		})

		s.Start()
		defer s.Shutdown()

		port := s.GetListenedPort()
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", port))

		// JSON 菜单
		resp, err := client.Get(nil, "/api/v1/topic/types?format=json")
		t.AssertNil(err)
		defer resp.Close()
		t.Assert(resp.StatusCode, 200)

		var result handler.Response
		err = json.Unmarshal(resp.ReadAll(), &result)
		t.AssertNil(err)
		t.Assert(result.Code, 20000)
		items := result.Data.(map[string]interface{})["items"].([]interface{})
		t.Assert(len(items) >= 12, true)
		foundIface := false
		for _, it := range items {
			m := it.(map[string]interface{})
			if m["name"] == "interface_basic" {
				foundIface = true
			}
		}
		t.Assert(foundIface, true)

		// JSON 内容
		resp2, err := client.Get(nil, "/api/v1/topic/types/boolean?format=json")
		t.AssertNil(err)
		defer resp2.Close()
		t.Assert(resp2.StatusCode, 200)

		var contentResp handler.Response
		err = json.Unmarshal(resp2.ReadAll(), &contentResp)
		t.AssertNil(err)
		t.Assert(contentResp.Code, 20000)
		dataMap := contentResp.Data.(map[string]interface{})
		t.AssertNE(dataMap["content"], nil)
		t.AssertNE(dataMap["quiz"], nil)

		// 接口子主题 JSON
		respIface, err := client.Get(nil, "/api/v1/topic/types/interface_basic?format=json")
		t.AssertNil(err)
		defer respIface.Close()
		t.Assert(respIface.StatusCode, 200)

		// HTML 菜单
		resp3, err := client.Get(nil, "/api/v1/topic/types?format=html")
		t.AssertNil(err)
		defer resp3.Close()
		t.Assert(resp3.StatusCode, 200)
		body := resp3.ReadAllString()
		t.AssertIN("Types Learning", body)
		t.AssertIN("Boolean", body)
		t.AssertIN("Interface", body)

		// 未知子主题 404 JSON
		resp404, err := client.Get(nil, "/api/v1/topic/types/unknown?format=json")
		t.AssertNil(err)
		defer resp404.Close()
		t.Assert(resp404.StatusCode, 404)
		var notFound handler.Response
		err = json.Unmarshal(resp404.ReadAll(), &notFound)
		t.AssertNil(err)
		t.Assert(notFound.Code, 404)

		// 未知子主题 404 HTML
		resp404h, err := client.Get(nil, "/api/v1/topic/types/unknown?format=html")
		t.AssertNil(err)
		defer resp404h.Close()
		t.Assert(resp404h.StatusCode, 404)
		body404 := resp404h.ReadAllString()
		t.AssertIN("未知的 Types 子主题", body404)
	})
}
