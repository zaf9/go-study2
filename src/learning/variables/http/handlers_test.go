package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// 覆盖 HTTP handler 的 content、quiz、submit。
func TestHandlers(t *testing.T) {
	h := NewHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/variables/content?topic=storage", nil)
	rr := httptest.NewRecorder()
	h.Content(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("content 返回码错误: %d", rr.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/variables/quiz?topic=storage", nil)
	rr = httptest.NewRecorder()
	h.Quiz(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("quiz 返回码错误: %d", rr.Code)
	}
	var quizResp struct {
		Items []struct {
			ID     string `json:"id"`
			Answer string `json:"answer"`
		} `json:"items"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &quizResp); err != nil {
		t.Fatalf("解析 quiz 失败: %v", err)
	}

	answers := map[string]string{}
	for _, item := range quizResp.Items {
		answers[item.ID] = item.Answer
	}
	body, _ := json.Marshal(map[string]interface{}{
		"answers": answers,
	})
	req = httptest.NewRequest(http.MethodPost, "/api/variables/quiz/submit?topic=storage", bytes.NewReader(body))
	rr = httptest.NewRecorder()
	h.SubmitQuiz(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("submit 返回码错误: %d", rr.Code)
	}
}
