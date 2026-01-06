package quiz

import (
	"context"
	"errors"
	"sort"
	"strings"
	"time"
)

// ErrInvalidInput 表示请求参数不合法。
var ErrInvalidInput = errors.New("测验参数不合法")

// ErrQuizUnavailable 表示当前主题暂无测验。
var ErrQuizUnavailable = errors.New("当前主题暂无测验")

// Service 封装测验相关业务。
type Service struct {
	repo Repository
}

// NewService 创建测验服务。
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// GetQuestions 获取测验题目列表。
// 已废弃: 使用 app/quiz.Service.GetQuizQuestions 替代
func (s *Service) GetQuestions(ctx context.Context, topic, chapter string) ([]Question, error) {
	topic = strings.TrimSpace(topic)
	chapter = strings.TrimSpace(chapter)
	if !IsSupportedTopic(topic) || chapter == "" {
		return nil, ErrInvalidInput
	}
	return nil, ErrQuizUnavailable
}

// Submit 提交答案并记录测验结果。
// 已废弃: 使用 app/quiz.Service.SubmitQuiz 替代
func (s *Service) Submit(ctx context.Context, userID int64, topic, chapter string, answers []SubmitAnswer, durationMs int64) (*Result, error) {
	return nil, errors.New("此方法已废弃")
}

// History 返回测验历史记录。
func (s *Service) History(ctx context.Context, userID int64, topic string, from, to *time.Time) ([]HistoryItem, error) {
	if userID <= 0 {
		return nil, ErrInvalidInput
	}
	if topic != "" && !IsSupportedTopic(strings.TrimSpace(topic)) {
		return nil, ErrInvalidInput
	}
	records, err := s.repo.ListRecords(ctx, userID, strings.TrimSpace(topic), from, to)
	if err != nil {
		return nil, err
	}
	items := make([]HistoryItem, 0, len(records))
	for _, rec := range records {
		items = append(items, HistoryItem{
			ID:         rec.ID,
			Topic:      rec.Topic,
			Chapter:    rec.Chapter,
			Score:      rec.Score,
			Total:      rec.Total,
			DurationMs: rec.DurationMs,
			CreatedAt:  rec.CreatedAt,
		})
	}
	return items, nil
}

func evaluate(questions []Question, submitted []SubmitAnswer, durationMs int64) *Result {
	answerMap := map[string][]string{}
	for _, ans := range submitted {
		answerMap[ans.ID] = normalizeChoices(ans.Choices)
	}

	correct := make([]string, 0, len(questions))
	wrong := make([]string, 0, len(questions))
	score := 0

	for _, q := range questions {
		expected := normalizeChoices(q.Answer)
		given := answerMap[q.ID]
		if equalChoice(expected, given) {
			score++
			correct = append(correct, q.ID)
		} else {
			wrong = append(wrong, q.ID)
		}
	}

	return &Result{
		Score:       score,
		Total:       len(questions),
		CorrectIDs:  correct,
		WrongIDs:    wrong,
		SubmittedAt: time.Now(),
		DurationMs:  durationMs,
	}
}

func normalizeChoices(list []string) []string {
	var cleaned []string
	for _, item := range list {
		item = strings.ToUpper(strings.TrimSpace(item))
		if item != "" {
			cleaned = append(cleaned, item)
		}
	}
	sort.Strings(cleaned)
	return cleaned
}

func equalChoice(expected, given []string) bool {
	if len(expected) == 0 && len(given) == 0 {
		return false
	}
	return strings.Join(expected, ",") == strings.Join(given, ",")
}
