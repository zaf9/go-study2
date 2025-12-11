package quiz

// 模块说明：测验服务负责抽题、创建会话、幂等判分与历史查询，是题库与接口之间的业务桥梁。

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"time"

	progressdom "go-study2/internal/domain/progress"
	quizdom "go-study2/internal/domain/quiz"
)

// ErrInvalidInput 表示请求参数不合法。
var ErrInvalidInput = errors.New("测验参数不合法")

// ErrQuizUnavailable 表示当前章节暂无测验。
var ErrQuizUnavailable = errors.New("当前章节暂无测验")

// ErrDuplicateSubmit 表示重复提交同一 session。
var ErrDuplicateSubmit = errors.New("重复提交会话")

// QuizRepository 抽象测验持久化操作，便于替换与测试。
type QuizRepository interface {
	GetQuestionsByChapter(ctx context.Context, topic, chapter string) ([]quizdom.QuizQuestion, error)
	CreateSession(ctx context.Context, session *quizdom.QuizSession) (string, error)
	SaveAttempts(ctx context.Context, attempts []quizdom.QuizAttempt) error
	GetHistory(ctx context.Context, userID int64, topic string, limit int) ([]quizdom.QuizSession, error)
	GetSession(ctx context.Context, sessionID string) (*quizdom.QuizSession, error)
	UpdateSessionResult(ctx context.Context, sessionID string, correct int, score int, passed bool) error
}

// Service 提供测验题目获取与提交判分。
type Service struct {
	repo       QuizRepository
	manager    *QuestionManager
	scorer     *ScoringEngine
	submitted  map[string]struct{}
	submitLock sync.Mutex
}

// NewService 创建测验服务。
func NewService(repo QuizRepository) *Service {
	return &Service{
		repo:      repo,
		manager:   NewQuestionManager(),
		scorer:    NewScoringEngine(),
		submitted: map[string]struct{}{},
	}
}

// QuizSessionPayload 返回题目与 session 信息。
type QuizSessionPayload struct {
	SessionID string        `json:"sessionId"`
	Topic     string        `json:"topic"`
	Chapter   string        `json:"chapter"`
	Questions []QuestionDTO `json:"questions"`
}

// SubmitResult 返回测验提交结果。
type SubmitResult struct {
	Score          int            `json:"score"`
	TotalQuestions int            `json:"total_questions"`
	CorrectAnswers int            `json:"correct_answers"`
	Passed         bool           `json:"passed"`
	Details        []AnswerDetail `json:"details"`
}

// GetQuizQuestions 抽取题目并创建测验 session。
func (s *Service) GetQuizQuestions(ctx context.Context, userID int64, topic, chapter string) (*QuizSessionPayload, error) {
	topic = strings.TrimSpace(topic)
	chapter = strings.TrimSpace(chapter)
	if userID <= 0 || !quizdom.IsSupportedTopic(topic) || chapter == "" {
		return nil, ErrInvalidInput
	}

	records, err := s.repo.GetQuestionsByChapter(ctx, topic, chapter)
	if err != nil {
		return nil, err
	}
	prepared, views, err := s.manager.Prepare(records)
	if err != nil {
		return nil, err
	}

	session := &quizdom.QuizSession{
		UserID:         userID,
		Topic:          topic,
		Chapter:        chapter,
		TotalQuestions: len(prepared),
		StartedAt:      time.Now(),
	}
	sessionID, err := s.repo.CreateSession(ctx, session)
	if err != nil {
		return nil, err
	}

	return &QuizSessionPayload{
		SessionID: sessionID,
		Topic:     topic,
		Chapter:   chapter,
		Questions: views,
	}, nil
}

// SubmitQuiz 评判答案并写入会话与作答记录。
func (s *Service) SubmitQuiz(ctx context.Context, userID int64, sessionID, topic, chapter string, answers []AnswerSubmission) (*SubmitResult, error) {
	topic = strings.TrimSpace(topic)
	chapter = strings.TrimSpace(chapter)
	if userID <= 0 || sessionID == "" || !quizdom.IsSupportedTopic(topic) || chapter == "" || len(answers) == 0 {
		return nil, ErrInvalidInput
	}

	s.submitLock.Lock()
	if _, ok := s.submitted[sessionID]; ok {
		s.submitLock.Unlock()
		return nil, ErrDuplicateSubmit
	}
	s.submitted[sessionID] = struct{}{}
	s.submitLock.Unlock()

	if existing, _ := s.repo.GetSession(ctx, sessionID); existing != nil && existing.CompletedAt != nil {
		return nil, ErrDuplicateSubmit
	}

	records, err := s.repo.GetQuestionsByChapter(ctx, topic, chapter)
	if err != nil {
		return nil, err
	}
	prepared, _, err := s.manager.Prepare(records)
	if err != nil {
		return nil, err
	}

	answerMap := map[int64][]string{}
	for _, a := range answers {
		answerMap[a.QuestionID] = normalizeChoices(a.UserAnswers)
	}
	score := s.scorer.Evaluate(prepared, answerMap)

	var attempts []quizdom.QuizAttempt
	for _, detail := range score.Details {
		rawAns, _ := json.Marshal(answerMap[detail.QuestionID])
		attempts = append(attempts, quizdom.QuizAttempt{
			SessionID:   sessionID,
			UserID:      userID,
			Topic:       topic,
			Chapter:     chapter,
			QuestionID:  detail.QuestionID,
			UserAnswers: string(rawAns),
			IsCorrect:   detail.IsCorrect,
			AttemptedAt: time.Now(),
		})
	}
	if err := s.repo.SaveAttempts(ctx, attempts); err != nil {
		return nil, err
	}
	if err := s.repo.UpdateSessionResult(ctx, sessionID, score.CorrectAnswers, score.Score, score.Passed); err != nil {
		return nil, err
	}

	return &SubmitResult{
		Score:          score.Score,
		TotalQuestions: score.TotalQuestions,
		CorrectAnswers: score.CorrectAnswers,
		Passed:         score.Passed,
		Details:        score.Details,
	}, nil
}

// GetQuizHistory 返回用户测验历史。
func (s *Service) GetQuizHistory(ctx context.Context, userID int64, topic string, limit int) ([]quizdom.QuizSession, error) {
	if userID <= 0 {
		return nil, ErrInvalidInput
	}
	if topic != "" && !quizdom.IsSupportedTopic(topic) {
		return nil, ErrInvalidInput
	}
	return s.repo.GetHistory(ctx, userID, strings.TrimSpace(topic), limit)
}

// ConvertProgressStatus 返回测验通过后的进度状态，供后续拓展。
func ConvertProgressStatus(passed bool) string {
	if passed {
		return progressdom.StatusCompleted
	}
	return progressdom.StatusTested
}
