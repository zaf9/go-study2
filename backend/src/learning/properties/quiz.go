package properties

// quiz.go 提供综合测验功能。
// 各子主题的测验在对应文件中注册，本文件用于综合测验的聚合。

// GetComprehensiveQuiz 获取综合测验，覆盖所有子主题。
func GetComprehensiveQuiz() ([]QuizItem, error) {
	var allItems []QuizItem
	for _, topic := range AllTopics() {
		items, err := LoadQuiz(topic)
		if err == nil {
			allItems = append(allItems, items...)
		}
	}
	if len(allItems) == 0 {
		return nil, ErrQuizUnavailable
	}
	return allItems, nil
}
