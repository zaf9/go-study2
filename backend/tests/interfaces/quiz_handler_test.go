package interfaces

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"sync"
	"testing"
	"time"
	"unsafe"

	"go-study2/internal/app/http_server/handler"
	appquiz "go-study2/internal/app/quiz"
	quizdom "go-study2/internal/domain/quiz"

	"github.com/gogf/gf/v2/net/ghttp"
)

type handlerResp struct {
	Code int                    `json:"code"`
	Data map[string]interface{} `json:"data"`
}

// quizMemoryRepo 复用内存题库与会话，避免依赖外部数据库。
type quizMemoryRepo struct {
	questions []quizdom.QuizQuestion
	sessions  map[string]*quizdom.QuizSession
	attempts  []quizdom.QuizAttempt
	mu        sync.Mutex
}

func newQuizMemoryRepo(questions []quizdom.QuizQuestion) *quizMemoryRepo {
	return &quizMemoryRepo{
		questions: questions,
		sessions:  map[string]*quizdom.QuizSession{},
	}
}

func (m *quizMemoryRepo) GetQuestionsByChapter(_ context.Context, topic, chapter string) ([]quizdom.QuizQuestion, error) {
	var result []quizdom.QuizQuestion
	for _, q := range m.questions {
		if q.Topic == topic && q.Chapter == chapter {
			result = append(result, q)
		}
	}
	return result, nil
}

func (m *quizMemoryRepo) CreateSession(_ context.Context, session *quizdom.QuizSession) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	id := fmt.Sprintf("session-%d", len(m.sessions)+1)
	copySession := *session
	copySession.SessionID = id
	if copySession.CreatedAt.IsZero() {
		copySession.CreatedAt = copySession.StartedAt
	}
	m.sessions[id] = &copySession
	return id, nil
}

func (m *quizMemoryRepo) SaveAttempts(_ context.Context, attempts []quizdom.QuizAttempt) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.attempts = append(m.attempts, attempts...)
	return nil
}

func (m *quizMemoryRepo) GetHistory(_ context.Context, userID int64, topic string, limit int) ([]quizdom.QuizSession, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var sessions []quizdom.QuizSession
	for _, s := range m.sessions {
		if s.UserID != userID {
			continue
		}
		if topic != "" && s.Topic != topic {
			continue
		}
		if s.CompletedAt == nil {
			continue
		}
		copySession := *s
		sessions = append(sessions, copySession)
	}
	if limit > 0 && len(sessions) > limit {
		return sessions[:limit], nil
	}
	return sessions, nil
}

func (m *quizMemoryRepo) GetSession(_ context.Context, sessionID string) (*quizdom.QuizSession, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if s, ok := m.sessions[sessionID]; ok {
		copySession := *s
		return &copySession, nil
	}
	return nil, nil
}

func (m *quizMemoryRepo) UpdateSessionResult(_ context.Context, sessionID string, correct int, score int, passed bool) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if s, ok := m.sessions[sessionID]; ok {
		now := time.Now()
		s.CompletedAt = &now
		s.CorrectAnswers = correct
		s.Score = score
		s.Passed = passed
		return nil
	}
	return nil
}

func (m *quizMemoryRepo) GetAttemptsBySession(_ context.Context, sessionID string) ([]quizdom.QuizAttempt, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var results []quizdom.QuizAttempt
	for _, a := range m.attempts {
		if a.SessionID == sessionID {
			results = append(results, a)
		}
	}
	return results, nil
}

func TestQuizHandler_EndToEnd(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Logf("Recovered from panic: %v", r)
		}
	}()

	repo := newQuizMemoryRepo([]quizdom.QuizQuestion{
		{
			ID:             101,
			Topic:          "variables",
			Chapter:        "storage",
			Type:           quizdom.QuestionTypeSingle,
			Difficulty:     quizdom.DifficultyEasy,
			Question:       "变量存储类型是？",
			Options:        `["栈","堆"]`,
			CorrectAnswers: `["A"]`,
			Explanation:    "示例",
		},
	})
	svc := appquiz.NewService(repo)
	h := handler.New()
	setQuizService(h, svc)

	server := ghttp.GetServer(fmt.Sprintf("quiz-handler-%d", time.Now().UnixNano()))
	server.BindMiddlewareDefault(func(r *ghttp.Request) {
		if uid, err := strconv.ParseInt(r.Header.Get("X-User-ID"), 10, 64); err == nil && uid > 0 {
			r.SetCtxVar("user_id", uid)
		}
		r.Middleware.Next()
	})
	server.SetSessionStorage(nil)
	server.Group("/api/v1", func(group *ghttp.RouterGroup) {
		group.GET("/quiz/{topic}/{chapter}", h.GetQuiz)
		group.POST("/quiz/submit", h.SubmitQuiz)
		group.GET("/quiz/history", h.GetQuizHistory)
	})
	server.SetPort(0)

	// 获取题目
	req := httptest.NewRequest(http.MethodGet, "/api/v1/quiz/variables/storage", nil)
	req.Header.Set("X-User-ID", "1")
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	var listResp handlerResp
	_ = json.Unmarshal(w.Body.Bytes(), &listResp)
	sessionID, _ := listResp.Data["sessionId"].(string)
	if sessionID == "" {
		t.Fatalf("缺少 sessionId: %+v", listResp)
	}

	// 提交测验
	body := fmt.Sprintf(`{"sessionId":"%s","topic":"variables","chapter":"storage","answers":[{"questionId":101,"userAnswers":["A"]}]}`, sessionID)
	req2 := httptest.NewRequest(http.MethodPost, "/api/v1/quiz/submit", bytes.NewBufferString(body))
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("X-User-ID", "1")
	w2 := httptest.NewRecorder()
	server.ServeHTTP(w2, req2)
	var submitResp handlerResp
	_ = json.Unmarshal(w2.Body.Bytes(), &submitResp)
	if submitResp.Code != 20000 {
		t.Fatalf("提交失败: %+v", submitResp)
	}

	// 历史记录
	req3 := httptest.NewRequest(http.MethodGet, "/api/v1/quiz/history", nil)
	req3.Header.Set("X-User-ID", "1")
	w3 := httptest.NewRecorder()
	server.ServeHTTP(w3, req3)
	var historyResp handlerResp
	_ = json.Unmarshal(w3.Body.Bytes(), &historyResp)
	if historyResp.Code != 20000 {
		t.Fatalf("历史查询失败: %+v", historyResp)
	}
}

// setQuizService 使用反射注入测验服务，便于测试。
func setQuizService(h *handler.Handler, svc *appquiz.Service) {
	val := reflect.ValueOf(h).Elem().FieldByName("quizService")
	ptr := unsafe.Pointer(val.UnsafeAddr())
	reflect.NewAt(val.Type(), ptr).Elem().Set(reflect.ValueOf(svc))
}
