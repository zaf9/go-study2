package http

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

// 覆盖错误路径与空答案的处理。
func TestHandlers_Errors(t *testing.T) {
	h := NewHandler()

	// 不支持的主题
	req := httptest.NewRequest(http.MethodGet, "/api/variables/content?topic=unknown", nil)
	rr := httptest.NewRecorder()
	h.Content(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("未知主题应返回 400，得到 %d", rr.Code)
	}

	// quiz 不支持的主题
	req = httptest.NewRequest(http.MethodGet, "/api/variables/quiz?topic=unknown", nil)
	rr = httptest.NewRecorder()
	h.Quiz(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("未知主题 quiz 应返回 400，得到 %d", rr.Code)
	}

	// submit 解析失败
	req = httptest.NewRequest(http.MethodPost, "/api/variables/quiz/submit?topic=storage", bytes.NewReader([]byte("not-json")))
	rr = httptest.NewRecorder()
	h.SubmitQuiz(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("无效 JSON 应返回 400，得到 %d", rr.Code)
	}

	// submit 空答案
	req = httptest.NewRequest(http.MethodPost, "/api/variables/quiz/submit?topic=storage", bytes.NewReader([]byte(`{"answers":{}}`)))
	rr = httptest.NewRecorder()
	h.SubmitQuiz(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("空答案应返回 400，得到 %d", rr.Code)
	}
}
