package constants

import (
	"bytes"
	"strings"
	"testing"
)

// TestDisplayMenu_ShowsAllOptions 测试菜单显示所有 12 个子主题选项
func TestDisplayMenu_ShowsAllOptions(t *testing.T) {
	// 准备输入: 立即输入 'q' 退出
	stdin := strings.NewReader("q\n")
	var stdout, stderr bytes.Buffer

	// 执行
	DisplayMenu(stdin, &stdout, &stderr)

	// 验证输出包含所有 12 个子主题
	output := stdout.String()

	expectedTopics := []string{
		"Boolean Constants (布尔常量)",
		"Rune Constants (符文常量)",
		"Integer Constants (整数常量)",
		"Floating-point Constants (浮点常量)",
		"Complex Constants (复数常量)",
		"String Constants (字符串常量)",
		"Constant Expressions (常量表达式)",
		"Typed and Untyped Constants (类型化/无类型化常量)",
		"Conversions (类型转换)",
		"Built-in Functions (内置函数)",
		"Iota (iota 特性)",
		"Implementation Restrictions (实现限制)",
	}

	for _, topic := range expectedTopics {
		if !strings.Contains(output, topic) {
			t.Errorf("菜单输出缺少主题: %s", topic)
		}
	}

	// 验证菜单标题
	if !strings.Contains(output, "Constants 学习菜单") {
		t.Error("菜单输出缺少标题")
	}

	// 验证退出选项
	if !strings.Contains(output, "q. 返回上级菜单") {
		t.Error("菜单输出缺少退出选项")
	}

	// 验证无错误输出
	if stderr.Len() > 0 {
		t.Errorf("存在意外的错误输出: %s", stderr.String())
	}
}

// TestDisplayMenu_InvalidInput 测试无效输入处理
func TestDisplayMenu_InvalidInput(t *testing.T) {
	// 准备输入: 先输入无效选项,再退出
	stdin := strings.NewReader("invalid\nq\n")
	var stdout, stderr bytes.Buffer

	// 执行
	DisplayMenu(stdin, &stdout, &stderr)

	// 验证输出包含错误提示
	output := stdout.String()
	if !strings.Contains(output, "无效的选择") {
		t.Error("无效输入时应显示错误提示")
	}
}

// TestDisplayMenu_SelectTopic 测试选择主题功能
func TestDisplayMenu_SelectTopic(t *testing.T) {
	// 准备输入: 选择主题 0 (Boolean),然后退出
	stdin := strings.NewReader("0\nq\n")
	var stdout, stderr bytes.Buffer

	// 执行
	DisplayMenu(stdin, &stdout, &stderr)

	// 验证输出包含 Boolean 主题内容
	output := stdout.String()
	if !strings.Contains(output, "布尔常量") {
		t.Error("选择主题 0 后应显示布尔常量内容")
	}
}
