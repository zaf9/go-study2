package types_test

import (
	"bytes"
	"strings"
	"testing"

	"go-study2/src/learning/types/cli"
)

// T009: CLI 菜单与展示集成测试
func TestTypesCLIMenuFlow(t *testing.T) {
	inputs := strings.Join([]string{
		"0", // 选择第一个主题
		"A", // 测验答案
		"8", // 选择接口基础主题
		"A",
		"99", // 非法编号
		"q",  // 返回菜单
		"q",  // 退出
	}, "\n") + "\n"

	stdin := bytes.NewBufferString(inputs)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cli.DisplayMenu(stdin, &stdout, &stderr)

	out := stdout.String()
	if !strings.Contains(out, "Types 学习菜单") {
		t.Fatalf("未找到菜单输出: %s", out)
	}
	if !strings.Contains(out, "布尔类型") && !strings.Contains(out, "Boolean") {
		t.Fatalf("未展示布尔主题内容: %s", out)
	}
	if !strings.Contains(out, "得分") {
		t.Fatalf("未展示测验得分: %s", out)
	}
	if !strings.Contains(out, "无效的选择") {
		t.Fatalf("未提示非法选择: %s", out)
	}
}
