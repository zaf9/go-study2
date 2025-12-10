package handler

import (
	"strings"
	"testing"
)

func TestGetHtmlPage(t *testing.T) {
	title := "Test Page"
	content := "<p>Test Content</p>"

	result := getHtmlPage(title, content)

	// 验证包含基本HTML结构
	if !strings.Contains(result, "<!DOCTYPE html>") {
		t.Error("HTML should contain DOCTYPE declaration")
	}

	if !strings.Contains(result, "<html>") {
		t.Error("HTML should contain html tag")
	}

	if !strings.Contains(result, "<head>") {
		t.Error("HTML should contain head tag")
	}

	if !strings.Contains(result, "<body>") {
		t.Error("HTML should contain body tag")
	}

	// 验证包含标题
	if !strings.Contains(result, title) {
		t.Errorf("HTML should contain title: %s", title)
	}

	// 验证包含内容
	if !strings.Contains(result, content) {
		t.Errorf("HTML should contain content: %s", content)
	}

	// 验证包含样式
	if !strings.Contains(result, "<style>") {
		t.Error("HTML should contain style tag")
	}

	if !strings.Contains(result, "font-family") {
		t.Error("HTML should contain CSS styles")
	}
}

func TestGetHtmlPage_EmptyContent(t *testing.T) {
	title := "Empty Page"
	content := ""

	result := getHtmlPage(title, content)

	// 即使内容为空，也应该生成有效的HTML结构
	if !strings.Contains(result, "<!DOCTYPE html>") {
		t.Error("HTML should contain DOCTYPE declaration even with empty content")
	}

	if !strings.Contains(result, title) {
		t.Errorf("HTML should contain title: %s", title)
	}
}

func TestGetHtmlPage_SpecialCharacters(t *testing.T) {
	title := "Test & <Special> Characters"
	content := "<p>Content with &amp; entities</p>"

	result := getHtmlPage(title, content)

	// 验证特殊字符被正确处理（虽然这里没有转义，但至少应该包含）
	if !strings.Contains(result, title) {
		t.Errorf("HTML should contain title with special characters: %s", title)
	}

	if !strings.Contains(result, content) {
		t.Errorf("HTML should contain content: %s", content)
	}
}
