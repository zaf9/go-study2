package types

import "strings"

// US2 综合测验数据（身份/可比较性/接口类型集判定）。
var comprehensiveQuiz = []QuizItem{
	{
		ID:          "q-all-1",
		ConceptID:   "interface_impl",
		Stem:        "指针方法集与值方法集的关系是？",
		Options:     []string{"指针方法集包含值方法集", "值方法集包含指针方法集"},
		Answer:      "A",
		Explanation: "指针方法集包含值接收者方法，反之不成立。",
		RuleRef:     "TR-IFACE-IMPL",
		Difficulty:  "medium",
	},
	{
		ID:          "q-all-2",
		ConceptID:   "map",
		Stem:        "下列哪种类型不能作为 map 键？",
		Options:     []string{"string", "[]int"},
		Answer:      "B",
		Explanation: "切片不可比较，不能作为键。",
		RuleRef:     "TR-MAP-KEY",
		Difficulty:  "medium",
	},
	{
		ID:          "q-all-3",
		ConceptID:   "array",
		Stem:        "数组长度是否影响赋值兼容性？",
		Options:     []string{"是，长度属于类型", "否，只看元素类型"},
		Answer:      "A",
		Explanation: "数组长度是类型一部分，长度不同不兼容。",
		RuleRef:     "TR-ARRAY-LEN",
		Difficulty:  "easy",
	},
	{
		ID:          "q-all-4",
		ConceptID:   "interface_general",
		Stem:        "~int 在类型集中表示什么？",
		Options:     []string{"底层类型为 int 的命名类型也匹配", "仅匹配 int"},
		Answer:      "A",
		Explanation: "~T 约束匹配底层类型为 T 的所有命名类型。",
		RuleRef:     "TR-IFACE-TYPESET",
		Difficulty:  "medium",
	},
	{
		ID:          "q-all-5",
		ConceptID:   "function",
		Stem:        "函数值可以直接比较吗？",
		Options:     []string{"只能与 nil 比较", "可以与任意函数比较"},
		Answer:      "A",
		Explanation: "函数类型不可比较，仅能与 nil 比较。",
		RuleRef:     "TR-FUNC-NIL",
		Difficulty:  "easy",
	},
}

// LoadComprehensiveQuiz 返回综合测验数据。
func LoadComprehensiveQuiz() []QuizItem {
	return comprehensiveQuiz
}

// EvaluateComprehensiveQuiz 对综合测验评分。
func EvaluateComprehensiveQuiz(answers map[string]string) (QuizResult, error) {
	items := LoadComprehensiveQuiz()
	score := 0
	var details []QuizAnswerFeedback
	for _, item := range items {
		choice := strings.TrimSpace(answers[item.ID])
		correct := strings.EqualFold(choice, strings.TrimSpace(item.Answer))
		if correct {
			score++
		}
		details = append(details, QuizAnswerFeedback{
			ID:          item.ID,
			Correct:     correct,
			Answer:      item.Answer,
			Explanation: item.Explanation,
			RuleRef:     item.RuleRef,
		})
	}
	return QuizResult{
		Score:   score,
		Total:   len(items),
		Details: details,
	}, nil
}
