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

// T022/T037: 检索接口契约测试
func TestTypesSearchContract(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.Server("test-types-search")
		s.SetPort(0)
		s.SetAccessLogEnabled(false)

		s.Group("/api/v1", func(group *ghttp.RouterGroup) {
			h := handler.New()
			group.Middleware(middleware.Format)
			group.ALL("/topic/types/search", h.SearchTypes)
		})

		s.Start()
		defer s.Shutdown()

		port := s.GetListenedPort()
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", port))

		resp, err := client.Get(nil, "/api/v1/topic/types/search?keyword=map%20key&format=json")
		t.AssertNil(err)
		defer resp.Close()
		t.Assert(resp.StatusCode, 200)
		var result handler.Response
		err = json.Unmarshal(resp.ReadAll(), &result)
		t.AssertNil(err)
		t.Assert(result.Code, 0)
		data := result.Data.(map[string]interface{})
		results := data["results"].([]interface{})
		t.AssertGT(len(results), 0)
	})
}
