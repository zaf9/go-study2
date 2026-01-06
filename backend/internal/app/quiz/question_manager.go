package quiz

// 模块说明：题目管理器负责将题库（YAML或数据库）转换为前端展示与判分输入，确保题型/选项/答案格式统一。

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
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
	ID          int64       `json:"-"`  // 内部使用，不序列化到JSON
	IDString    string      `json:"id"` // JSON序列化时使用字符串格式，避免JavaScript大整数精度问题
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
			ID:          item.ID,                    // 内部使用 int64
			IDString:    fmt.Sprintf("%d", item.ID), // JSON序列化为字符串
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

// PrepareFromYAML 将YAML题目转换为准备好的题目（用于YAML数据源）
func (m *QuestionManager) PrepareFromYAML(yamlQuestions []quizdom.YAMLQuestion) ([]PreparedQuestion, []QuestionDTO, error) {
	if len(yamlQuestions) == 0 {
		return nil, nil, ErrQuizUnavailable
	}
	var prepared []PreparedQuestion
	var views []QuestionDTO

	for _, yq := range yamlQuestions {
		// 构建选项
		opts := make([]OptionDTO, 0, len(yq.Options))
		for idx, opt := range yq.Options {
			opts = append(opts, OptionDTO{
				ID:    optionID(idx), // A, B, C, D
				Label: opt,
			})
		}

		// 解析答案（支持单选"A"或多选"ABC"）
		correctAnswer := parseYAMLAnswer(yq.Answer)

		// 使用字符串ID的哈希值作为int64 ID（保持前端兼容）
		id := hashQuestionID(yq.ID)

		view := QuestionDTO{
			ID:         id,                    // 内部使用 int64
			IDString:   fmt.Sprintf("%d", id), // JSON序列化为字符串
			Type:       yq.Type,
			Difficulty: yq.Difficulty,
			Question:   yq.Stem,
			Options:    opts,
		}

		prepared = append(prepared, PreparedQuestion{
			View:          view,
			CorrectAnswer: correctAnswer,
			Explanation:   yq.Explanation,
		})
		views = append(views, view)
	}

	return prepared, views, nil
}

// parseYAMLAnswer 解析YAML答案格式（"A" 或 "ABC" 等）
func parseYAMLAnswer(answer string) []string {
	answer = strings.ToUpper(strings.TrimSpace(answer))
	result := make([]string, 0, len(answer))
	for _, ch := range answer {
		if ch >= 'A' && ch <= 'Z' {
			result = append(result, string(ch))
		}
	}
	return result
}

// hashQuestionID 将字符串ID转为int64（使用哈希保证唯一性）
func hashQuestionID(id string) int64 {
	h := sha256.Sum256([]byte(id))
	// 取前8字节转为int64
	hashStr := hex.EncodeToString(h[:8])
	var result int64
	for i := 0; i < len(hashStr) && i < 16; i++ {
		var val int64
		if hashStr[i] >= '0' && hashStr[i] <= '9' {
			val = int64(hashStr[i] - '0')
		} else {
			val = int64(hashStr[i]-'a') + 10
		}
		result = result*16 + val
	}
	// 确保为正数
	if result < 0 {
		result = -result
	}
	return result
}
