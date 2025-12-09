package variables_contract_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	varhttp "go-study2/src/learning/variables/http"
)

// 测试零值主题的 content 与 submit。
func TestZeroHTTP(t *testing.T) {
	h := varhttp.NewHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/variables/content?topic=zero", nil)
	rr := httptest.NewRecorder()
	h.Content(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("content 返回码错误: %d", rr.Code)
	}

	answers := map[string]string{
		"q-zero-1": "C",
		"q-zero-2": "B",
	}
	body, _ := json.Marshal(map[string]interface{}{
		"answers": answers,
	})
	req = httptest.NewRequest(http.MethodPost, "/api/variables/quiz/submit?topic=zero", bytes.NewReader(body))
	rr = httptest.NewRecorder()
	h.SubmitQuiz(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("submit 返回码错误: %d", rr.Code)
	}
}
