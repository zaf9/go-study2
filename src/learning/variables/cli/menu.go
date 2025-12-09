package cli

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"

	"go-study2/src/learning/variables"
)

// Menu 提供变量章节的 CLI 菜单骨架。
type Menu struct{}

// NewMenu 创建变量章节菜单。
func NewMenu() *Menu {
	return &Menu{}
}

// DisplayMenu 仿照 Constants 章节的交互方式，提供编号选择的子菜单。
func DisplayMenu(stdin io.Reader, stdout, stderr io.Writer) {
	reader := bufio.NewReader(stdin)

	subMenu := map[string]variables.Topic{
		"0": variables.TopicStorage,
		"1": variables.TopicStatic,
		"2": variables.TopicDynamic,
		"3": variables.TopicZero,
	}
	subMenuDesc := map[string]string{
		"0": "Storage (存储与取址)",
		"1": "Static (静态类型与可赋值性)",
		"2": "Dynamic (接口动态类型与 nil)",
		"3": "Zero (零值与取值规则)",
	}

	for {
		fmt.Fprintln(stdout, "\nVariables 学习菜单")
		fmt.Fprintln(stdout, "---------------------------------")
		fmt.Fprintln(stdout, "请选择要学习的主题:")
		for i := 0; i <= 3; i++ {
			key := fmt.Sprintf("%d", i)
			fmt.Fprintf(stdout, "%s. %s\n", key, subMenuDesc[key])
		}
		fmt.Fprintln(stdout, "q. 返回上级菜单")
		fmt.Fprint(stdout, "\n请输入您的选择: ")

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(stderr, "读取输入失败: %v\n", err)
			return
		}
		choice := strings.TrimSpace(input)
		if choice == "q" {
			return
		}
		topic, ok := subMenu[choice]
		if !ok {
			fmt.Fprintln(stdout, "无效的选择,请重试。")
			continue
		}
		displayTopic(topic, stdout, stderr)
	}
}

// ListTopics 返回可用主题。
func (m *Menu) ListTopics() []variables.Topic {
	return variables.AllTopics()
}

// ShowContent 展示指定主题内容，当前为占位实现。
func (m *Menu) ShowContent(topic string) (variables.Content, error) {
	normalized := variables.NormalizeTopic(topic)
	if !variables.IsSupportedTopic(normalized) {
		return variables.Content{}, variables.ErrUnsupportedTopic
	}
	return variables.FetchContent(normalized)
}

// StartQuiz 获取测验数据，当前为占位实现。
func (m *Menu) StartQuiz(topic string) ([]variables.QuizItem, error) {
	normalized := variables.NormalizeTopic(topic)
	return variables.LoadQuiz(normalized)
}

// SubmitQuiz 评估答案，当前为占位实现。
func (m *Menu) SubmitQuiz(topic string, answers map[string]string) (variables.QuizResult, error) {
	if len(answers) == 0 {
		return variables.QuizResult{}, errors.New("未提供答案")
	}
	normalized := variables.NormalizeTopic(topic)
	return variables.EvaluateQuiz(normalized, answers)
}

// displayTopic 输出主题内容与测验列表，便于与 Constants 章节交互一致。
func displayTopic(topic variables.Topic, stdout, stderr io.Writer) {
	content, err := variables.LoadContent(topic)
	if err != nil {
		fmt.Fprintf(stderr, "获取内容失败: %v\n", err)
		return
	}
	fmt.Fprintf(stdout, "\n[%s] %s\n", content.Topic, content.Title)
	fmt.Fprintln(stdout, content.Summary)
	for _, d := range content.Details {
		fmt.Fprintf(stdout, "- %s\n", d)
	}
	for _, ex := range content.Examples {
		fmt.Fprintf(stdout, "示例: %s\n%s\n=> %s\n", ex.Title, ex.Code, ex.Output)
	}

	items, err := variables.LoadQuiz(topic)
	if err == variables.ErrQuizUnavailable {
		fmt.Fprintln(stdout, "当前主题暂无测验数据。")
		return
	}
	if err != nil {
		fmt.Fprintf(stderr, "获取测验失败: %v\n", err)
		return
	}
	fmt.Fprintln(stdout, "\n测验题目:")
	for _, item := range items {
		fmt.Fprintf(stdout, "[%s] %s\n", item.ID, item.Stem)
		for idx, opt := range item.Options {
			fmt.Fprintf(stdout, "%c) %s\n", 'A'+idx, opt)
		}
		fmt.Fprintf(stdout, "答案: %s，解析: %s\n", item.Answer, item.Explanation)
	}
}

// FormatTopicList 将主题列表格式化为可读文本。
func FormatTopicList(topics []variables.Topic) string {
	if len(topics) == 0 {
		return "暂无可用主题"
	}
	return fmt.Sprintf("可用主题: %v", topics)
}
