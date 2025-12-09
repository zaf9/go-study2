package types_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"go-study2/internal/app/http_server/handler"
	"go-study2/internal/app/http_server/middleware"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/test/gtest"
)

// T016: 测验提交契约测试
func TestTypesQuizContract(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.Server("test-types-quiz")
		s.SetPort(0)
		s.SetAccessLogEnabled(false)

		s.Group("/api/v1", func(group *ghttp.RouterGroup) {
			h := handler.New()
			group.Middleware(middleware.Format)
			group.ALL("/topic/types/quiz/submit", h.SubmitTypesQuiz)
		})

		s.Start()
		defer s.Shutdown()

		port := s.GetListenedPort()
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", port))

		payload := []byte(`{"answers":[{"id":"q-all-1","choice":"A"},{"id":"q-all-2","choice":"B"},{"id":"q-all-3","choice":"A"},{"id":"q-all-4","choice":"A"},{"id":"q-all-5","choice":"A"}]}`)
		resp, err := client.Post(nil, "/api/v1/topic/types/quiz/submit?format=json", bytes.NewReader(payload))
		t.AssertNil(err)
		defer resp.Close()
		t.Assert(resp.StatusCode, 200)

		var result handler.Response
		err = json.Unmarshal(resp.ReadAll(), &result)
		t.AssertNil(err)
		t.Assert(result.Code, 0)
		data := result.Data.(map[string]interface{})
		t.AssertGT(data["score"], 0)
	})
}
