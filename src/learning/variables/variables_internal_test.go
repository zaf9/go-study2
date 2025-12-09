package variables

import "testing"

// 覆盖内容与测验的加载、校验与评分。
func TestContentAndQuizLifecycle(t *testing.T) {
	topics := []Topic{TopicStorage, TopicStatic, TopicDynamic, TopicZero}
	for _, tp := range topics {
		content, err := LoadContent(tp)
		if err != nil {
			t.Fatalf("LoadContent(%s) 失败: %v", tp, err)
		}
		if content.Topic != tp {
			t.Fatalf("主题不一致: %s", content.Topic)
		}
		items, err := LoadQuiz(tp)
		if err != nil {
			t.Fatalf("LoadQuiz(%s) 失败: %v", tp, err)
		}
		if err := ValidateQuizItems(items); err != nil {
			t.Fatalf("ValidateQuizItems(%s) 失败: %v", tp, err)
		}
		answers := map[string]string{}
		for _, item := range items {
			answers[item.ID] = item.Answer
		}
		result, err := EvaluateQuiz(tp, answers)
		if err != nil {
			t.Fatalf("EvaluateQuiz(%s) 失败: %v", tp, err)
		}
		if result.Score != result.Total {
			t.Fatalf("评分应全对: %d/%d", result.Score, result.Total)
		}
	}
}

// 覆盖异常主题路径。
func TestUnsupportedTopic(t *testing.T) {
	_, err := LoadContent("unknown")
	if err == nil {
		t.Fatalf("未知主题应报错")
	}
	_, err = LoadQuiz("unknown")
	if err == nil {
		t.Fatalf("未知主题应报错")
	}
	_, err = EvaluateQuiz("unknown", map[string]string{"id": "A"})
	if err == nil {
		t.Fatalf("未知主题评估应报错")
	}
}
