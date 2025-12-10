package variables

import "fmt"

var contentTemplates = map[Topic]Content{}
var quizBank = map[Topic][]QuizItem{}

// RegisterContent 可在后续实现中为主题填充内容。
func RegisterContent(topic Topic, content Content) {
	if !IsSupportedTopic(topic) {
		return
	}
	contentTemplates[topic] = content
}

// FetchContent 从模板表中读取内容，如未填充则返回占位。
func FetchContent(topic Topic) (Content, error) {
	normalized := NormalizeTopic(string(topic))
	if !IsSupportedTopic(normalized) {
		return Content{}, ErrUnsupportedTopic
	}
	if content, ok := contentTemplates[normalized]; ok {
		return content, nil
	}
	return Content{}, fmt.Errorf("主题暂无内容: %s", normalized)
}

// RegisterQuiz 注册测验题目。
func RegisterQuiz(topic Topic, items []QuizItem) {
	if !IsSupportedTopic(topic) {
		return
	}
	quizBank[topic] = items
}

// FetchQuiz 获取测验题目，未注册则返回 ErrQuizUnavailable。
func FetchQuiz(topic Topic) ([]QuizItem, error) {
	normalized := NormalizeTopic(string(topic))
	if !IsSupportedTopic(normalized) {
		return nil, ErrUnsupportedTopic
	}
	items, ok := quizBank[normalized]
	if !ok || len(items) == 0 {
		return nil, ErrQuizUnavailable
	}
	return items, nil
}
