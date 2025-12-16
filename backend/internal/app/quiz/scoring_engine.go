package quiz

// 模块说明：判分引擎负责按题型计算得分与细节，用于提交测验时的核心评分逻辑。

import (
	"math"
	"strings"
)

// AnswerSubmission 表示用户对单题的回答。
type AnswerSubmission struct {
	QuestionID  int64    `json:"questionId"`
	UserAnswers []string `json:"userAnswers"`
}

// AnswerDetail 返回判分细节。
type AnswerDetail struct {
	QuestionID     int64    `json:"question_id"`
	IsCorrect      bool     `json:"is_correct"`
	CorrectAnswers []string `json:"correct_answers"`
	Explanation    string   `json:"explanation"`
	ScorePart      float64  `json:"score_part"`
}

// ScoringResult 汇总判分结果。
type ScoringResult struct {
	Score          int
	TotalQuestions int
	CorrectAnswers int
	Passed         bool
	Details        []AnswerDetail
}

// ScoringEngine 负责根据题型计算得分。
type ScoringEngine struct{}

// NewScoringEngine 创建判分引擎。
func NewScoringEngine() *ScoringEngine {
	return &ScoringEngine{}
}

// Evaluate 对答题结果进行评分。
func (e *ScoringEngine) Evaluate(questions []PreparedQuestion, answers map[int64][]string) ScoringResult {
	total := len(questions)
	if total == 0 {
		return ScoringResult{}
	}
	var details []AnswerDetail
	var totalScore float64
	var correctCount int

	for _, q := range questions {
		userAns := normalizeChoices(answers[q.View.ID])
		correct := normalizeChoices(q.CorrectAnswer)
		part := e.scoreByType(q.View.Type, correct, userAns)
		if part >= 1 {
			correctCount++
		}
		totalScore += part
		details = append(details, AnswerDetail{
			QuestionID:     q.View.ID,
			IsCorrect:      part >= 1,
			CorrectAnswers: correct,
			Explanation:    q.Explanation,
			ScorePart:      part,
		})
	}

	score := int(math.Round((totalScore / float64(total)) * 100))
	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	return ScoringResult{
		Score:          score,
		TotalQuestions: total,
		CorrectAnswers: correctCount,
		Passed:         score >= 60,
		Details:        details,
	}
}

// EvaluateWithTotal 对答题结果进行评分，使用指定的总题数计算得分。
// 这用于当实际测试的题目数量与传入的题目数量不一致时（例如从session中获取的实际测试题目数）。
func (e *ScoringEngine) EvaluateWithTotal(questions []PreparedQuestion, answers map[int64][]string, actualTotalQuestions int) ScoringResult {
	if actualTotalQuestions <= 0 {
		return ScoringResult{}
	}
	var details []AnswerDetail
	var totalScore float64
	var correctCount int

	// 只对用户提交答案的题目进行评分
	for _, q := range questions {
		userAns, hasAnswer := answers[q.View.ID]
		if !hasAnswer {
			// 如果用户没有提交答案，该题得分为0
			details = append(details, AnswerDetail{
				QuestionID:     q.View.ID,
				IsCorrect:      false,
				CorrectAnswers: q.CorrectAnswer,
				Explanation:    q.Explanation,
				ScorePart:      0,
			})
			continue
		}

		userAns = normalizeChoices(userAns)
		correct := normalizeChoices(q.CorrectAnswer)
		part := e.scoreByType(q.View.Type, correct, userAns)
		if part >= 1 {
			correctCount++
		}
		totalScore += part
		details = append(details, AnswerDetail{
			QuestionID:     q.View.ID,
			IsCorrect:      part >= 1,
			CorrectAnswers: correct,
			Explanation:    q.Explanation,
			ScorePart:      part,
		})
	}

	// 使用实际测试的题目数量计算得分
	score := int(math.Round((totalScore / float64(actualTotalQuestions)) * 100))
	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	return ScoringResult{
		Score:          score,
		TotalQuestions: actualTotalQuestions,
		CorrectAnswers: correctCount,
		Passed:         score >= 60,
		Details:        details,
	}
}

func (e *ScoringEngine) scoreByType(qType string, correct, given []string) float64 {
	if len(correct) == 0 {
		return 0
	}
	switch qType {
	case "multiple", "code_correction":
		return partialScore(correct, given)
	default:
		if equalChoice(correct, given) {
			return 1
		}
		return 0
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
	sortStrings(cleaned)
	return cleaned
}

func equalChoice(expected, given []string) bool {
	if len(expected) != len(given) {
		return false
	}
	for i := range expected {
		if expected[i] != given[i] {
			return false
		}
	}
	return true
}

func partialScore(expected, given []string) float64 {
	if len(given) == 0 {
		return 0
	}
	expSet := map[string]struct{}{}
	for _, v := range expected {
		expSet[v] = struct{}{}
	}
	for _, g := range given {
		if _, ok := expSet[g]; !ok {
			return 0
		}
	}
	correctCount := 0
	for _, g := range given {
		if _, ok := expSet[g]; ok {
			correctCount++
		}
	}
	return float64(correctCount) / float64(len(expected))
}

func sortStrings(list []string) {
	for i := 0; i < len(list); i++ {
		for j := i + 1; j < len(list); j++ {
			if list[j] < list[i] {
				list[i], list[j] = list[j], list[i]
			}
		}
	}
}
