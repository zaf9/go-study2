package variables_test

import (
	"strings"
	"testing"

	"go-study2/src/learning/variables"
)

// 测试结构化元素的零值描述是否存在。
func TestStructuredElementsZeroValue(t *testing.T) {
	content, err := variables.LoadContent(variables.TopicStorage)
	if err != nil {
		t.Fatalf("加载内容失败: %v", err)
	}
	found := false
	for _, detail := range content.Details {
		if strings.Contains(detail, "元素") && strings.Contains(detail, "零值") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("应包含结构化元素零值的说明")
	}
	if !strings.Contains(content.Snippet, "复合字面量") && !strings.Contains(content.Snippet, "结构体") {
		t.Fatalf("代码片段应展示结构体或复合字面量取址示例")
	}
}
