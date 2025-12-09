package cli

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"go-study2/src/learning/types"
)

// Menu 提供 Types 章节的 CLI 菜单交互。
type Menu struct{}

// NewMenu 创建 Types 菜单实例。
func NewMenu() *Menu {
	return &Menu{}
}

// DisplayMenu 提供子主题选择与基础测验体验。
func DisplayMenu(stdin io.Reader, stdout, stderr io.Writer) {
	reader := bufio.NewReader(stdin)
	subtopics := types.AllTopics()

	for {
		fmt.Fprintln(stdout, "\nTypes 学习菜单")
		fmt.Fprintln(stdout, "---------------------------------")
		for idx, topic := range subtopics {
			fmt.Fprintf(stdout, "%d. %s\n", idx, formatTitle(topic))
		}
		fmt.Fprintln(stdout, "o. 打印提纲")
		fmt.Fprintln(stdout, "quiz. 综合测验")
		fmt.Fprintln(stdout, "search <keyword>. 关键词检索")
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
		if strings.HasPrefix(strings.ToLower(choice), "search ") {
			keyword := strings.TrimSpace(choice[len("search "):])
			handleSearch(keyword, stdout, stderr)
			continue
		}
		if choice == "o" || strings.EqualFold(choice, "outline") {
			handleOutline(stdout)
			continue
		}
		if strings.EqualFold(choice, "quiz") {
			handleComprehensiveQuiz(reader, stdout, stderr)
			continue
		}
		idx, err := parseIndex(choice, len(subtopics))
		if err != nil {
			fmt.Fprintln(stdout, "无效的选择,请重试。")
			continue
		}
		handleTopic(subtopics[idx], reader, stdout, stderr)
	}
}

func handleTopic(topic types.Topic, reader *bufio.Reader, stdout, stderr io.Writer) {
	content, err := types.LoadContent(topic)
	if err != nil {
		fmt.Fprintf(stderr, "加载内容失败: %v\n", err)
		return
	}
	fmt.Fprintf(stdout, "\n[%s] %s\n", content.Concept.ID, content.Concept.Title)
	fmt.Fprintln(stdout, content.Concept.Summary)
	if len(content.Concept.Rules) > 0 {
		fmt.Fprintln(stdout, "规则:")
		for _, rule := range content.Concept.Rules {
			fmt.Fprintf(stdout, "- %s\n", rule)
		}
	}
	if len(content.Examples) > 0 {
		fmt.Fprintln(stdout, "示例:")
		for _, ex := range content.Examples {
			fmt.Fprintf(stdout, "• %s\n%s\n", ex.Title, ex.Code)
			if ex.ExpectedOutput != "" {
				fmt.Fprintf(stdout, "=> %s\n", ex.ExpectedOutput)
			}
		}
	}

	items, err := types.LoadQuiz(topic)
	if err == types.ErrQuizUnavailable {
		fmt.Fprintln(stdout, "当前主题暂无测验数据。")
		return
	}
	if err != nil {
		fmt.Fprintf(stderr, "获取测验失败: %v\n", err)
		return
	}

	fmt.Fprintln(stdout, "\n测验: 输入选项字母作答，输入 q 结束测验。")
	answerMap := map[string]string{}
	for _, item := range items {
		fmt.Fprintf(stdout, "%s: %s\n", item.ID, item.Stem)
		for idx, opt := range item.Options {
			fmt.Fprintf(stdout, "%c) %s\n", 'A'+idx, opt)
		}
		fmt.Fprint(stdout, "答案: ")
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(stderr, "读取答案失败: %v\n", err)
			return
		}
		ans := strings.TrimSpace(text)
		if ans == "q" {
			break
		}
		answerMap[item.ID] = ans
	}

	result, err := types.EvaluateQuiz(topic, answerMap)
	if err != nil {
		fmt.Fprintf(stderr, "评分失败: %v\n", err)
		return
	}
	fmt.Fprintf(stdout, "得分: %d/%d\n", result.Score, result.Total)
	for _, d := range result.Details {
		state := "错误"
		if d.Correct {
			state = "正确"
		}
		fmt.Fprintf(stdout, "- %s: %s (答案: %s) - %s\n", d.ID, state, d.Answer, d.Explanation)
	}
}

func parseIndex(raw string, max int) (int, error) {
	idx, err := strconv.Atoi(raw)
	if err != nil {
		return -1, fmt.Errorf("invalid index")
	}
	if idx < 0 || idx >= max {
		return -1, fmt.Errorf("invalid index")
	}
	return idx, nil
}

func formatTitle(topic types.Topic) string {
	switch topic {
	case types.TopicBoolean:
		return "Boolean (布尔)"
	case types.TopicNumeric:
		return "Numeric (数值)"
	case types.TopicString:
		return "String (字符串)"
	case types.TopicArray:
		return "Array (数组)"
	case types.TopicSlice:
		return "Slice (切片)"
	case types.TopicStruct:
		return "Struct (结构体)"
	case types.TopicPointer:
		return "Pointer (指针)"
	case types.TopicFunction:
		return "Function (函数)"
	case types.TopicInterfaceBasic:
		return "Interface Basic (接口基础)"
	case types.TopicInterfaceEmbedded:
		return "Interface Embedded (接口嵌入)"
	case types.TopicInterfaceGeneral:
		return "Interface General (类型集)"
	case types.TopicInterfaceImpl:
		return "Interface Impl (接口实现)"
	case types.TopicMap:
		return "Map (映射)"
	case types.TopicChannel:
		return "Channel (通道)"
	default:
		return string(topic)
	}
}

func handleOutline(stdout io.Writer) {
	fmt.Fprintln(stdout, "\nTypes 章节提纲：")
	text := types.RenderPrintableOutline()
	fmt.Fprintln(stdout, text)
}

func handleComprehensiveQuiz(reader *bufio.Reader, stdout, stderr io.Writer) {
	for {
		items := types.LoadComprehensiveQuiz()
		answerMap := map[string]string{}
		fmt.Fprintln(stdout, "\n综合测验（输入选项字母，输入 q 结束）：")
		for _, item := range items {
			fmt.Fprintf(stdout, "%s: %s\n", item.ID, item.Stem)
			for idx, opt := range item.Options {
				fmt.Fprintf(stdout, "%c) %s\n", 'A'+idx, opt)
			}
			fmt.Fprint(stdout, "答案: ")
			text, err := reader.ReadString('\n')
			if err != nil {
				fmt.Fprintf(stderr, "读取答案失败: %v\n", err)
				return
			}
			ans := strings.TrimSpace(text)
			if ans == "q" {
				return
			}
			answerMap[item.ID] = ans
		}
		result, err := types.EvaluateComprehensiveQuiz(answerMap)
		if err != nil {
			fmt.Fprintf(stderr, "评分失败: %v\n", err)
			return
		}
		fmt.Fprintf(stdout, "综合得分: %d/%d\n", result.Score, result.Total)
		for _, d := range result.Details {
			state := "错误"
			if d.Correct {
				state = "正确"
			}
			fmt.Fprintf(stdout, "- %s: %s (答案: %s) - %s\n", d.ID, state, d.Answer, d.Explanation)
		}
		fmt.Fprint(stdout, "是否重做？(y/n): ")
		redo, _ := reader.ReadString('\n')
		if strings.ToLower(strings.TrimSpace(redo)) != "y" {
			return
		}
	}
}

func handleSearch(keyword string, stdout, stderr io.Writer) {
	results, err := types.SearchReferences(keyword)
	if err != nil {
		fmt.Fprintf(stderr, "检索失败: %v\n", err)
		return
	}
	if len(results) == 0 {
		fmt.Fprintln(stdout, "未找到匹配关键词。")
		return
	}
	fmt.Fprintf(stdout, "关键词 \"%s\" 结果:\n", keyword)
	for _, r := range results {
		fmt.Fprintf(stdout, "- %s: %s (%s)\n", r.Keyword, r.Summary, r.Anchors["http"])
	}
}
