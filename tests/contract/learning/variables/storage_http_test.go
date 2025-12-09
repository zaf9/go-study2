package variables_contract_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-study2/src/learning/variables"
	varhttp "go-study2/src/learning/variables/http"
)

// 测试 content 与 quiz 接口契约。
func TestStorageHTTPContract(t *testing.T) {
	h := varhttp.NewHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/variables/content?topic=storage", nil)
	rr := httptest.NewRecorder()
	h.Content(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("content 返回码错误: %d", rr.Code)
	}
	var content variables.Content
	if err := json.Unmarshal(rr.Body.Bytes(), &content); err != nil {
		t.Fatalf("解析 content 失败: %v", err)
	}
	if content.Topic != variables.TopicStorage {
		t.Fatalf("content topic 错误: %s", content.Topic)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/variables/quiz?topic=storage", nil)
	rr = httptest.NewRecorder()
	h.Quiz(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("quiz 返回码错误: %d", rr.Code)
	}
	var quizResp struct {
		Topic variables.Topic      `json:"topic"`
		Items []variables.QuizItem `json:"items"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &quizResp); err != nil {
		t.Fatalf("解析 quiz 失败: %v", err)
	}
	if quizResp.Topic != variables.TopicStorage {
		t.Fatalf("quiz topic 错误: %s", quizResp.Topic)
	}
	if len(quizResp.Items) == 0 {
		t.Fatalf("quiz items 为空")
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
	var result variables.QuizResult
	if err := json.Unmarshal(rr.Body.Bytes(), &result); err != nil {
		t.Fatalf("解析 submit 失败: %v", err)
	}
	if result.Score != result.Total {
		t.Fatalf("应全部答对，得分 %d/%d", result.Score, result.Total)
	}
}
