//go:build ignore
// +build ignore

package main

// 生成器: 根据清单自动生成占位题库YAML文件，供人工后续替换为高质量题目。
// 使用方法：在仓库根目录运行 `go run backend/scripts/generate_quiz_yaml.go`

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

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

type QuizBank struct {
	Questions []QuizQuestion `yaml:"questions"`
}

func main() {
	rand.Seed(time.Now().UnixNano())

	manifest := map[string]map[string]int{
		"lexical_elements": {
			"comments": 35, "tokens": 40, "semicolons": 30, "identifiers": 40,
			"keywords": 30, "operators": 45, "integers": 40, "floats": 35,
			"imaginary": 30, "runes": 40, "strings": 45,
		},
		"constants": {
			"boolean": 30, "rune": 35, "integer": 40, "floating_point": 40,
			"complex": 35, "string": 35, "expressions": 45, "typed_untyped": 40,
			"conversions": 40, "builtin_functions": 35, "iota": 45, "implementation_restrictions": 30,
		},
		"variables": {
			"storage": 40, "static": 35, "dynamic": 40, "zero": 45,
		},
		"types": {
			"boolean": 30, "numeric": 50, "string": 40, "array": 45, "slice": 50,
			"struct": 50, "pointer": 40, "function": 45, "interface_basic": 45,
			"interface_embedded": 40, "interface_general": 45, "interface_impl": 40,
			"map": 45, "channel": 40,
		},
	}

	root := filepath.Join("backend", "quiz_data")
	for topic, chapters := range manifest {
		for chapter, count := range chapters {
			dir := filepath.Join(root, topic)
			if err := os.MkdirAll(dir, 0o755); err != nil {
				fmt.Fprintf(os.Stderr, "创建目录失败: %v\n", err)
				os.Exit(2)
			}

			bank := QuizBank{Questions: make([]QuizQuestion, 0, count)}
			for i := 1; i <= count; i++ {
				id := fmt.Sprintf("%s-%s-%03d", topicPrefix(topic), chapter, i)
				typ := sampleType(i)
				diff := sampleDifficulty(i, count)
				optCount := optionCountForType(typ)
				opts := make([]string, optCount)
				labels := []string{"A", "B", "C", "D", "E"}
				for j := 0; j < optCount; j++ {
					opts[j] = fmt.Sprintf("%s: 示例选项 %s", labels[j], labels[j])
				}
				ans := sampleAnswer(typ, optCount)
				q := QuizQuestion{
					ID:          id,
					Type:        typ,
					Difficulty:  diff,
					Stem:        fmt.Sprintf("【占位】%s - %s 第 %d 题：请替换为真实题干。", topic, chapter, i),
					Options:     opts,
					Answer:      ans,
					Explanation: "【占位解析】示例解析，请人工替换为详实中文解析。",
					Topic:       topic,
					Chapter:     chapter,
				}
				bank.Questions = append(bank.Questions, q)
			}

			outPath := filepath.Join(dir, chapter+".yaml")
			f, err := os.Create(outPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "创建文件失败 %s: %v\n", outPath, err)
				os.Exit(3)
			}
			enc := yaml.NewEncoder(f)
			enc.SetIndent(2)
			if err := enc.Encode(bank); err != nil {
				fmt.Fprintf(os.Stderr, "写入YAML失败 %s: %v\n", outPath, err)
				f.Close()
				os.Exit(4)
			}
			f.Close()
			fmt.Printf("生成: %s (%d 题)\n", outPath, count)
		}
	}

	// 生成 README 占位文件
	readmePath := filepath.Join(root, "README.md")
	if err := os.WriteFile(readmePath, []byte(readmeContent()), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "写入 README 失败: %v\n", err)
		os.Exit(5)
	}
	fmt.Println("题库占位文件生成完成。请人工审核并替换占位题目为高质量题目。")
}

func topicPrefix(topic string) string {
	switch topic {
	case "lexical_elements":
		return "lexical"
	case "constants":
		return "const"
	case "variables":
		return "var"
	case "types":
		return "type"
	default:
		return topic
	}
}

func sampleType(i int) string {
	if i%2 == 0 {
		return "single"
	}
	return "multiple"
}

func sampleDifficulty(i, total int) string {
	// 简单约40%，中等40%，困难20%
	r := float64(i) / float64(total)
	if r <= 0.4 {
		return "easy"
	}
	if r <= 0.8 {
		return "medium"
	}
	return "hard"
}

func optionCountForType(typ string) int {
	if typ == "single" {
		// 单选2-4
		return 2 + rand.Intn(3) // 2-4
	}
	// 多选3-5
	return 3 + rand.Intn(3) // 3-5
}

func sampleAnswer(typ string, optCount int) string {
	labels := []string{"A", "B", "C", "D", "E"}
	if typ == "single" {
		return labels[rand.Intn(optCount)]
	}
	// 多选: 随机选择2-4不同字母，按字母升序返回
	k := 2 + rand.Intn(3)
	if k > optCount {
		k = optCount
	}
	idxs := rand.Perm(optCount)[:k]
	sel := make([]string, 0, k)
	for _, id := range idxs {
		sel = append(sel, labels[id])
	}
	// 排序
	for i := 0; i < len(sel)-1; i++ {
		for j := i + 1; j < len(sel); j++ {
			if sel[i] > sel[j] {
				sel[i], sel[j] = sel[j], sel[i]
			}
		}
	}
	ans := ""
	for _, s := range sel {
		ans += s
	}
	return ans
}

func readmeContent() string {
	return `# 题库数据目录 (占位生成)

本目录由脚本 ` + "backend/scripts/generate_quiz_yaml.go" + ` 生成，包含按主题分组的章节 YAML 文件。

注意：此脚本生成的是占位题目，供人工审核和替换为高质量题目。生成的题目示例中包含中文占位题干和解析，请在发布前进行人工校验并完善题目内容。

生成规则：
- 每个章节文件为 YAML 格式，根节点为 'questions'。
- 每题包含字段：id, type, difficulty, stem, options, answer, explanation, topic, chapter
- 请参照 specs/013-quiz-question-bank/data-model.md 中的数据模型与约束进行最终校验。

生成命令：

    go run backend/scripts/generate_quiz_yaml.go

`
}
