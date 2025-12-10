package variables_contract_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	varhttp "go-study2/src/learning/variables/http"
)

// 测试静态与动态类型 content 接口。
func TestTypesHTTPContent(t *testing.T) {
	h := varhttp.NewHandler()
	topics := []string{"static", "dynamic"}
	for _, tp := range topics {
		req := httptest.NewRequest(http.MethodGet, "/api/variables/content?topic="+tp, nil)
		rr := httptest.NewRecorder()
		h.Content(rr, req)
		if rr.Code != http.StatusOK {
			t.Fatalf("%s content 返回码错误: %d", tp, rr.Code)
		}
	}
}
