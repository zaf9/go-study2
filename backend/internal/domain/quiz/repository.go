package quiz

// repository.go 可扩展的仓储方法（简单包装）

// Topics 返回所有主题名称
func (r *QuizRepository) Topics() []string {
	out := []string{}
	for t := range r.banks {
		out = append(out, t)
	}
	return out
}

// Chapters 返回指定主题的章节列表
func (r *QuizRepository) Chapters(topic string) []string {
	out := []string{}
	if chs, ok := r.banks[topic]; ok {
		for c := range chs {
			out = append(out, c)
		}
	}
	return out
}
