package variables_integration_test

import (
	"testing"

	"go-study2/src/learning/variables/cli"
)

// 测试 CLI 对静态/动态类型的展示。
func TestTypesCLI(t *testing.T) {
	menu := cli.NewMenu()
	topics := []string{"static", "dynamic"}
	for _, tp := range topics {
		_, err := menu.ShowContent(tp)
		if err != nil {
			t.Fatalf("%s ShowContent 失败: %v", tp, err)
		}
		items, err := menu.StartQuiz(tp)
		if err != nil {
			t.Fatalf("%s StartQuiz 失败: %v", tp, err)
		}
		if len(items) == 0 {
			t.Fatalf("%s 题目为空", tp)
		}
		answers := map[string]string{}
		for _, item := range items {
			answers[item.ID] = item.Answer
		}
		if _, err := menu.SubmitQuiz(tp, answers); err != nil {
			t.Fatalf("%s SubmitQuiz 失败: %v", tp, err)
		}
	}
}
