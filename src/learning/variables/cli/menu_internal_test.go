package cli

import "testing"

// 覆盖错误主题与空答案路径。
func TestMenu_Errors(t *testing.T) {
	menu := NewMenu()

	if _, err := menu.ShowContent("unknown"); err == nil {
		t.Fatalf("未知主题应报错")
	}
	if _, err := menu.StartQuiz("unknown"); err == nil {
		t.Fatalf("未知主题应报错")
	}
	if _, err := menu.SubmitQuiz("storage", map[string]string{}); err == nil {
		t.Fatalf("空答案应报错")
	}
}
