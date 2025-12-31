//go:build ignore
// +build ignore

package main

// 生成器: 为单个章节生成 YAML 题库骨架，包含元数据和题目占位符
// 使用方法: go run generate_quiz_template.go <topic> <chapter> <question_count>
// 示例: go run generate_quiz_template.go lexical_elements keywords 40

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// QuizMetadata 题库元数据
type QuizMetadata struct {
	Topic          string `yaml:"topic"`
	Chapter        string `yaml:"chapter"`
	TotalQuestions int    `yaml:"total_questions"`
	LastUpdated    string `yaml:"last_updated"` // YYYY-MM-DD
	Description    string `yaml:"description"`
}

// QuizQuestion 题目结构
type QuizQuestion struct {
	ID          string   `yaml:"id"`
	Type        string   `yaml:"type"`
	Difficulty  string   `yaml:"difficulty"`
	Stem        string   `yaml:"stem"`
	Options     []string `yaml:"options"`
	Answer      string   `yaml:"answer"`
	Explanation string   `yaml:"explanation"`
	Topic       string   `yaml:"topic"`
	Chapter     string   `yaml:"chapter"`
}

// QuizBank 题库结构
type QuizBank struct {
	Metadata  QuizMetadata   `yaml:"metadata"`
	Questions []QuizQuestion `yaml:"questions"`
}

func main() {
	if len(os.Args) < 4 {
		fmt.Fprintf(os.Stderr, "用法: go run generate_quiz_template.go <topic> <chapter> <question_count>\n")
		fmt.Fprintf(os.Stderr, "示例: go run generate_quiz_template.go lexical_elements keywords 40\n")
		os.Exit(1)
	}

	topic := os.Args[1]
	chapter := os.Args[2]
	var count int
	if _, err := fmt.Sscanf(os.Args[3], "%d", &count); err != nil || count <= 0 {
		fmt.Fprintf(os.Stderr, "错误: 题目数量必须为正整数\n")
		os.Exit(1)
	}

	// 验证 topic
	validTopics := map[string]bool{
		"lexical_elements": true,
		"constants":        true,
		"variables":        true,
		"types":            true,
	}
	if !validTopics[topic] {
		fmt.Fprintf(os.Stderr, "错误: 无效的主题 '%s'，有效值: lexical_elements, constants, variables, types\n", topic)
		os.Exit(1)
	}

	// 创建目录
	dir := filepath.Join("backend", "quiz_data", topic)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "创建目录失败: %v\n", err)
		os.Exit(2)
	}

	// 生成题库
	bank := QuizBank{
		Metadata: QuizMetadata{
			Topic:          topic,
			Chapter:        chapter,
			TotalQuestions: count,
			LastUpdated:    time.Now().Format("2006-01-02"),
			Description:    fmt.Sprintf("【待补充】%s - %s 章节测验题库", getTopicName(topic), chapter),
		},
		Questions: make([]QuizQuestion, 0, count),
	}

	// 生成题目占位符
	// 题型分布: 40% single_choice, 30% multiple_choice, 20% code_output, 10% code_fix
	// 难度分布: 60% easy, 30% medium, 10% hard
	typeDistribution := []string{
		"single_choice", "single_choice", "single_choice", "single_choice", // 40%
		"multiple_choice", "multiple_choice", "multiple_choice",            // 30%
		"code_output", "code_output",                                       // 20%
		"code_fix", // 10%
	}

	difficultyDistribution := []string{
		"easy", "easy", "easy", "easy", "easy", "easy", // 60%
		"medium", "medium", "medium",                   // 30%
		"hard", // 10%
	}

	for i := 1; i <= count; i++ {
		questionType := typeDistribution[i%len(typeDistribution)]
		difficulty := difficultyDistribution[i%len(difficultyDistribution)]

		options := generateOptions(questionType)
		answer := generateAnswer(questionType, len(options))

		question := QuizQuestion{
			ID:          fmt.Sprintf("%s-%03d", chapter, i),
			Type:        questionType,
			Difficulty:  difficulty,
			Stem:        fmt.Sprintf("【占位】%s 第 %d 题：请替换为真实题干，描述清晰，引用 Go 1.24 规范相关章节", chapter, i),
			Options:     options,
			Answer:      answer,
			Explanation: "【占位】请补充详细解析：1. 解释正确答案原理；2. 引用 Go 规范章节（如 §3.1）；3. 提供代码示例（可选）；4. 说明常见错误",
			Topic:       topic,
			Chapter:     chapter,
		}
		bank.Questions = append(bank.Questions, question)
	}

	// 写入文件
	outPath := filepath.Join(dir, chapter+".yaml")
	f, err := os.Create(outPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "创建文件失败 %s: %v\n", outPath, err)
		os.Exit(3)
	}
	defer f.Close()

	enc := yaml.NewEncoder(f)
	enc.SetIndent(2)
	if err := enc.Encode(bank); err != nil {
		fmt.Fprintf(os.Stderr, "写入 YAML 失败 %s: %v\n", outPath, err)
		os.Exit(4)
	}

	fmt.Printf("✅ 成功生成: %s (%d 题)\n", outPath, count)
	fmt.Printf("📝 请编辑文件，将占位内容替换为高质量题目\n")
	fmt.Printf("🔍 完成后运行验证: go run validate_quiz_yaml.go %s\n", outPath)
}

// generateOptions 根据题型生成选项
func generateOptions(questionType string) []string {
	switch questionType {
	case "single_choice":
		return []string{
			"A: 【占位】选项 A",
			"B: 【占位】选项 B",
			"C: 【占位】选项 C",
			"D: 【占位】选项 D",
		}
	case "multiple_choice":
		return []string{
			"A: 【占位】选项 A",
			"B: 【占位】选项 B",
			"C: 【占位】选项 C",
			"D: 【占位】选项 D",
			"E: 【占位】选项 E（可选）",
		}
	case "code_output":
		return []string{
			"A: 【占位】输出结果 A",
			"B: 【占位】输出结果 B",
			"C: 【占位】编译错误",
			"D: 【占位】运行时错误",
		}
	case "code_fix":
		return []string{
			"A: 【占位】修复方案 A",
			"B: 【占位】修复方案 B",
			"C: 【占位】修复方案 C",
			"D: 【占位】修复方案 D",
		}
	default:
		return []string{}
	}
}

// generateAnswer 根据题型生成答案占位符
func generateAnswer(questionType string, optionCount int) string {
	switch questionType {
	case "single_choice", "code_output", "code_fix":
		return "【占位】B" // 单选默认 B
	case "multiple_choice":
		return "【占位】AC" // 多选默认 AC
	default:
		return "【占位】"
	}
}

// getTopicName 获取主题中文名
func getTopicName(topic string) string {
	names := map[string]string{
		"lexical_elements": "词法元素",
		"constants":        "常量",
		"variables":        "变量",
		"types":            "类型系统",
	}
	if name, ok := names[topic]; ok {
		return name
	}
	return topic
}
