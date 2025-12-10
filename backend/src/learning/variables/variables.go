package variables

import (
	"errors"
	"fmt"
	"strings"
)

// Topic 表示变量章节的子主题。
type Topic string

const (
	TopicStorage Topic = "storage"
	TopicStatic  Topic = "static"
	TopicDynamic Topic = "dynamic"
	TopicZero    Topic = "zero"
)

var (
	ErrUnsupportedTopic = errors.New("不支持的变量主题")
	ErrQuizUnavailable  = errors.New("当前主题暂无测验数据")
)

// Example 用于展示代码示例与预期输出。
type Example struct {
	Title  string   `json:"title"`
	Code   string   `json:"code"`
	Output string   `json:"output"`
	Notes  []string `json:"notes,omitempty"`
}

// Content 描述某个主题的教学内容。
type Content struct {
	Topic    Topic     `json:"topic"`
	Title    string    `json:"title"`
	Summary  string    `json:"summary"`
	Details  []string  `json:"details,omitempty"`
	Examples []Example `json:"examples,omitempty"`
	Snippet  string    `json:"snippet,omitempty"`
}

// QuizItem 定义单个测验题目。
type QuizItem struct {
	ID          string   `json:"id"`
	Topic       Topic    `json:"topic"`
	Stem        string   `json:"stem"`
	Options     []string `json:"options"`
	Answer      string   `json:"answer"`
	Explanation string   `json:"explanation"`
}

// QuizAnswerFeedback 反馈答题结果。
type QuizAnswerFeedback struct {
	ID          string `json:"id"`
	Correct     bool   `json:"correct"`
	Explanation string `json:"explanation"`
}

// QuizResult 汇总测验得分与明细。
type QuizResult struct {
	Score   int                  `json:"score"`
	Total   int                  `json:"total"`
	Details []QuizAnswerFeedback `json:"details"`
}

// AllTopics 返回章节支持的所有主题。
func AllTopics() []Topic {
	return []Topic{TopicStorage, TopicStatic, TopicDynamic, TopicZero}
}

// NormalizeTopic 归一化主题字符串。
func NormalizeTopic(raw string) Topic {
	return Topic(strings.ToLower(strings.TrimSpace(raw)))
}

// IsSupportedTopic 判断主题是否受支持。
func IsSupportedTopic(topic Topic) bool {
	for _, t := range AllTopics() {
		if t == topic {
			return true
		}
	}
	return false
}

// LoadContent 返回指定主题的内容。
func LoadContent(topic Topic) (Content, error) {
	normalized := NormalizeTopic(string(topic))
	if !IsSupportedTopic(normalized) {
		return Content{}, ErrUnsupportedTopic
	}
	return FetchContent(normalized)
}

// LoadQuiz 返回指定主题的测验。
func LoadQuiz(topic Topic) ([]QuizItem, error) {
	normalized := NormalizeTopic(string(topic))
	if !IsSupportedTopic(normalized) {
		return nil, ErrUnsupportedTopic
	}
	return FetchQuiz(normalized)
}

// EvaluateQuiz 按提交答案评估得分。
func EvaluateQuiz(topic Topic, answers map[string]string) (QuizResult, error) {
	if len(answers) == 0 {
		return QuizResult{}, fmt.Errorf("未提供答案")
	}
	items, err := LoadQuiz(topic)
	if err != nil {
		return QuizResult{}, err
	}
	total := len(items)
	score := 0
	var details []QuizAnswerFeedback
	for _, item := range items {
		choice := answers[item.ID]
		correct := strings.EqualFold(strings.TrimSpace(choice), strings.TrimSpace(item.Answer))
		if correct {
			score++
		}
		details = append(details, QuizAnswerFeedback{
			ID:          item.ID,
			Correct:     correct,
			Explanation: item.Explanation,
		})
	}
	return QuizResult{
		Score:   score,
		Total:   total,
		Details: details,
	}, nil
}

// ValidateQuizItems 校验测验题目的基础结构。
func ValidateQuizItems(items []QuizItem) error {
	for _, item := range items {
		if item.ID == "" {
			return fmt.Errorf("测验题目缺少 ID")
		}
		if !IsSupportedTopic(item.Topic) {
			return fmt.Errorf("测验题目主题无效: %s", item.Topic)
		}
		if item.Stem == "" {
			return fmt.Errorf("测验题目缺少题干: %s", item.ID)
		}
		if len(item.Options) < 2 {
			return fmt.Errorf("测验题目选项不足: %s", item.ID)
		}
		if item.Answer == "" {
			return fmt.Errorf("测验题目缺少答案: %s", item.ID)
		}
	}
	return nil
}
