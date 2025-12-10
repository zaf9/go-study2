package quiz

import (
	"context"
	"encoding/json"
	"errors"
	"sort"
	"strings"
	"time"

	"go-study2/src/learning/types"
	"go-study2/src/learning/variables"
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
func (s *Service) GetQuestions(ctx context.Context, topic, chapter string) ([]Question, error) {
	topic = strings.TrimSpace(topic)
	chapter = strings.TrimSpace(chapter)
	if !IsSupportedTopic(topic) || chapter == "" {
		return nil, ErrInvalidInput
	}
	return loadQuestions(topic, chapter)
}

// Submit 提交答案并记录测验结果。
func (s *Service) Submit(ctx context.Context, userID int64, topic, chapter string, answers []SubmitAnswer, durationMs int64) (*Result, error) {
	topic = strings.TrimSpace(topic)
	chapter = strings.TrimSpace(chapter)
	if userID <= 0 || !IsSupportedTopic(topic) || chapter == "" || len(answers) == 0 {
		return nil, ErrInvalidInput
	}

	questions, err := loadQuestions(topic, chapter)
	if err != nil {
		if errors.Is(err, ErrQuizUnavailable) {
			return nil, err
		}
		return nil, err
	}
	if len(questions) == 0 {
		return nil, ErrQuizUnavailable
	}

	result := evaluate(questions, answers, durationMs)
	encoded, _ := json.Marshal(answers)
	record := &Record{
		UserID:     userID,
		Topic:      topic,
		Chapter:    chapter,
		Score:      result.Score,
		Total:      result.Total,
		DurationMs: result.DurationMs,
		Answers:    string(encoded),
	}

	if _, err := s.repo.SaveRecord(ctx, record); err != nil {
		return nil, err
	}
	return result, nil
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

func loadQuestions(topic, chapter string) ([]Question, error) {
	switch topic {
	case "variables":
		t := variables.NormalizeTopic(chapter)
		items, err := variables.LoadQuiz(t)
		if err != nil {
			if errors.Is(err, variables.ErrQuizUnavailable) {
				return []Question{}, ErrQuizUnavailable
			}
			return nil, err
		}
		return convertVariableQuiz(items), nil
	case "types":
		t := types.NormalizeTopic(chapter)
		items, err := types.LoadQuiz(t)
		if err != nil {
			if errors.Is(err, types.ErrQuizUnavailable) {
				return []Question{}, ErrQuizUnavailable
			}
			return nil, err
		}
		return convertTypeQuiz(items), nil
	default:
		return []Question{}, ErrQuizUnavailable
	}
}

func convertVariableQuiz(items []variables.QuizItem) []Question {
	list := make([]Question, 0, len(items))
	for _, item := range items {
		options := make([]Option, 0, len(item.Options))
		for idx, opt := range item.Options {
			options = append(options, Option{
				ID:    optionID(idx),
				Label: opt,
			})
		}
		list = append(list, Question{
			ID:          item.ID,
			Stem:        item.Stem,
			Options:     options,
			Multi:       len(strings.TrimSpace(item.Answer)) > 1,
			Answer:      splitAnswer(item.Answer),
			Explanation: item.Explanation,
		})
	}
	return list
}

func convertTypeQuiz(items []types.QuizItem) []Question {
	list := make([]Question, 0, len(items))
	for _, item := range items {
		options := make([]Option, 0, len(item.Options))
		for idx, opt := range item.Options {
			options = append(options, Option{
				ID:    optionID(idx),
				Label: opt,
			})
		}
		list = append(list, Question{
			ID:          item.ID,
			Stem:        item.Stem,
			Options:     options,
			Multi:       len(strings.TrimSpace(item.Answer)) > 1,
			Answer:      splitAnswer(item.Answer),
			Explanation: item.Explanation,
		})
	}
	return list
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

func optionID(idx int) string {
	return string(rune('A' + idx))
}

func splitAnswer(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, "")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.ToUpper(strings.TrimSpace(p))
		if p != "" {
			result = append(result, p)
		}
	}
	return result
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

