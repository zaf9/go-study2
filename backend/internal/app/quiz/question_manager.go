package quiz

// 模块说明：题目管理器负责将数据库题库转换为前端展示与判分输入，确保题型/选项/答案格式统一。

import (
	"encoding/json"
	"errors"
	"strings"

	quizdom "go-study2/internal/domain/quiz"
)

// OptionDTO 表示前端可见的选项。
type OptionDTO struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

// QuestionDTO 表示前端可见的题目。
type QuestionDTO struct {
	ID          int64       `json:"id"`
	Type        string      `json:"type"`
	Difficulty  string      `json:"difficulty"`
	Question    string      `json:"question"`
	Options     []OptionDTO `json:"options"`
	CodeSnippet *string     `json:"codeSnippet,omitempty"`
}

// PreparedQuestion 保存判分所需的题目信息。
type PreparedQuestion struct {
	View          QuestionDTO
	CorrectAnswer []string
	Explanation   string
}

// QuestionManager 负责将题库记录转换为可用的题目与答案。
type QuestionManager struct{}

// NewQuestionManager 创建题目管理器。
func NewQuestionManager() *QuestionManager {
	return &QuestionManager{}
}

// Prepare 将数据库题目转换为判分结构与前端视图。
func (m *QuestionManager) Prepare(records []quizdom.QuizQuestion) ([]PreparedQuestion, []QuestionDTO, error) {
	if len(records) == 0 {
		return nil, nil, ErrQuizUnavailable
	}
	var prepared []PreparedQuestion
	var views []QuestionDTO
	for _, item := range records {
		opts, err := parseOptions(item.Options)
		if err != nil {
			return nil, nil, err
		}
		correct, err := parseAnswers(item.CorrectAnswers)
		if err != nil {
			return nil, nil, err
		}
		view := QuestionDTO{
			ID:          item.ID,
			Type:        item.Type,
			Difficulty:  item.Difficulty,
			Question:    item.Question,
			Options:     opts,
			CodeSnippet: item.CodeSnippet,
		}
		prepared = append(prepared, PreparedQuestion{
			View:          view,
			CorrectAnswer: correct,
			Explanation:   item.Explanation,
		})
		views = append(views, view)
	}
	return prepared, views, nil
}

func parseOptions(raw string) ([]OptionDTO, error) {
	var arr []string
	if err := json.Unmarshal([]byte(raw), &arr); err != nil {
		return nil, err
	}
	opts := make([]OptionDTO, 0, len(arr))
	for idx, val := range arr {
		id := optionID(idx)
		opts = append(opts, OptionDTO{ID: id, Label: val})
	}
	return opts, nil
}

func parseAnswers(raw string) ([]string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, errors.New("缺少正确答案")
	}
	var arr []interface{}
	if err := json.Unmarshal([]byte(raw), &arr); err != nil {
		return nil, err
	}
	result := make([]string, 0, len(arr))
	for _, v := range arr {
		switch val := v.(type) {
		case string:
			val = strings.ToUpper(strings.TrimSpace(val))
			if val != "" {
				result = append(result, val)
			}
		case float64:
			result = append(result, optionID(int(val)))
		default:
			// 忽略未知类型
		}
	}
	return result, nil
}

func optionID(idx int) string {
	return string(rune('A' + idx))
}
