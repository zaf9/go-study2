package handler

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/test/gtest"
)

var ctx = gctx.New()

// TestGetConstantsMenu 测试获取 Constants 菜单
func TestGetConstantsMenu(t *testing.T) {
	s := g.Server("test-constants-menu")
	s.SetPort(0)
	s.SetAccessLogEnabled(false)

	s.Group("/", func(group *ghttp.RouterGroup) {
		h := New()
		group.GET("/constants", h.GetConstantsMenu)
	})
	s.Start()
	defer s.Shutdown()

	port := s.GetListenedPort()
	client := g.Client()
	client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", port))

	gtest.C(t, func(t *gtest.T) {
		// Test JSON response
		resp, err := client.Get(ctx, "/constants")
		t.AssertNil(err)
		defer resp.Close()

		var result Response
		err = json.Unmarshal(resp.ReadAll(), &result)
		t.AssertNil(err)
		t.Assert(result.Code, 0)

		// 验证数据结构
		dataMap := result.Data.(map[string]interface{})
		items := dataMap["items"].([]interface{})
		t.Assert(len(items), 12) // 应有 12 个子主题

		// 验证包含特定主题
		foundBoolean := false
		for _, item := range items {
			itemMap := item.(map[string]interface{})
			if itemMap["name"] == "boolean" {
				foundBoolean = true
				break
			}
		}
		t.Assert(foundBoolean, true)
	})
}
