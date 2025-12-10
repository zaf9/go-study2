package http

import (
	"encoding/json"
	"net/http"

	"go-study2/src/learning/variables"
)

// Handler 提供变量章节的 HTTP 处理骨架。
type Handler struct{}

// NewHandler 创建 Handler。
func NewHandler() *Handler {
	return &Handler{}
}

// Content 返回指定主题的内容。
func (h *Handler) Content(w http.ResponseWriter, r *http.Request) {
	topic := variables.NormalizeTopic(r.URL.Query().Get("topic"))
	content, err := variables.FetchContent(topic)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, content)
}

// Quiz 返回指定主题的测验题目。
func (h *Handler) Quiz(w http.ResponseWriter, r *http.Request) {
	topic := variables.NormalizeTopic(r.URL.Query().Get("topic"))
	items, err := variables.LoadQuiz(topic)
	if err != nil && err != variables.ErrQuizUnavailable {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err == variables.ErrQuizUnavailable {
		writeError(w, http.StatusNotImplemented, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"topic": topic,
		"items": items,
	})
}

// SubmitQuiz 接收答题并返回评估结果。
func (h *Handler) SubmitQuiz(w http.ResponseWriter, r *http.Request) {
	topic := variables.NormalizeTopic(r.URL.Query().Get("topic"))
	var payload struct {
		Answers map[string]string `json:"answers"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, "请求体解析失败")
		return
	}
	result, err := variables.EvaluateQuiz(topic, payload.Answers)
	if err != nil && err != variables.ErrQuizUnavailable {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err == variables.ErrQuizUnavailable {
		writeError(w, http.StatusNotImplemented, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]interface{}{
		"error": msg,
	})
}
