//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"io/fs"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type QuizQuestion struct {
	ID          string   `yaml:"id" json:"id"`
	Type        string   `yaml:"type" json:"type"`
	Difficulty  string   `yaml:"difficulty" json:"difficulty"`
	Stem        string   `yaml:"stem" json:"stem"`
	Options     []string `yaml:"options" json:"options"`
	Answer      string   `yaml:"answer" json:"answer"`
	Explanation string   `yaml:"explanation" json:"explanation"`
	Topic       string   `yaml:"topic" json:"topic"`
	Chapter     string   `yaml:"chapter" json:"chapter"`
}

type QuizBank struct {
	Questions []QuizQuestion `yaml:"questions" json:"questions"`
}

func main() {
	rand.Seed(time.Now().UnixNano())
	// try possible quiz_data locations to be robust to how generator was run
	candidates := []string{
		filepath.Join("backend", "backend", "quiz_data"),
		filepath.Join("backend", "quiz_data"),
		filepath.Join("quiz_data"),
	}
	var root string
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			root = c
			break
		}
	}
	if root == "" {
		fmt.Fprintln(os.Stderr, "未找到 quiz_data 目录，尝试的路径:", candidates)
		os.Exit(1)
	}
	count := 0
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(d.Name(), ".yaml") {
			return nil
		}
		// Skip file if it's comments.yaml already enhanced (we keep it)
		if strings.Contains(path, "lexical_elements") && strings.HasSuffix(path, "comments.yaml") {
			fmt.Println("跳过已增强文件:", path)
			return nil
		}

		fmt.Println("处理:", path)
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		var bank QuizBank
		if err := yaml.Unmarshal(b, &bank); err != nil {
			return fmt.Errorf("解析 YAML 失败 %s: %w", path, err)
		}

		for i := range bank.Questions {
			q := &bank.Questions[i]
			enhanceQuestion(q)
		}

		// ensure answer letters are sorted for multiple
		for i := range bank.Questions {
			if bank.Questions[i].Type == "multiple" {
				letters := strings.Split(bank.Questions[i].Answer, "")
				sort.Strings(letters)
				bank.Questions[i].Answer = strings.Join(letters, "")
			}
		}

		out, err := yaml.Marshal(&bank)
		if err != nil {
			return err
		}
		if err := os.WriteFile(path, out, 0o644); err != nil {
			return err
		}
		count++
		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "批量增强失败: %v\n", err)
		os.Exit(2)
	}
	fmt.Printf("完成：增强 %d 个文件\n", count)
}

func enhanceQuestion(q *QuizQuestion) {
	// normalize
	topic := q.Topic
	chapter := q.Chapter
	// generate stem based on difficulty
	switch q.Difficulty {
	case "easy":
		q.Stem = fmt.Sprintf("【单项基础】关于 %s.%s，下面哪个选项是正确的？", topic, chapter)
	case "medium":
		q.Stem = fmt.Sprintf("【应用题】在 %s.%s 的情境下，哪个说法更合适？", topic, chapter)
	case "hard":
		q.Stem = fmt.Sprintf("【思考题】针对 %s.%s，判断下列说法的正确性并选择最合适的选项。", topic, chapter)
	default:
		q.Stem = fmt.Sprintf("关于 %s.%s 的题目：请结合语境选择正确答案。", topic, chapter)
	}

	// decide option count
	var optCount int
	if q.Type == "single" {
		optCount = 4
	} else {
		optCount = 4
		// allow 3-5 but prefer 4
	}
	labels := []string{"A", "B", "C", "D", "E"}
	opts := make([]string, optCount)

	// create a plausible correct index
	correctIdx := rand.Intn(optCount)
	// for difficulty adjust preferred correct index
	switch q.Difficulty {
	case "easy":
		correctIdx = rand.Intn(2) // A or B
	case "medium":
		correctIdx = 1 + rand.Intn(2) // B or C
	case "hard":
		correctIdx = 2 + rand.Intn(max(1, optCount-2)) // later options
	}

	// Fill options with templated content
	for i := 0; i < optCount; i++ {
		label := labels[i]
		if i == correctIdx && q.Type == "single" {
			opts[i] = fmt.Sprintf("%s: 正确：关于%s.%s的关键点描述（正确选项）", label, chapter, topic)
			continue
		}
		// distractors
		d := rand.Intn(4)
		switch d {
		case 0:
			opts[i] = fmt.Sprintf("%s: 常见误解或反例（容易混淆）", label)
		case 1:
			opts[i] = fmt.Sprintf("%s: 部分正确但不完整的说法", label)
		case 2:
			opts[i] = fmt.Sprintf("%s: 语法/语义上不正确的陈述", label)
		default:
			opts[i] = fmt.Sprintf("%s: 与%s.%s 相关但不适用的选项", label, chapter, topic)
		}
	}

	// For multiple, pick 2-3 correct indices
	if q.Type == "multiple" {
		k := 2 + rand.Intn(2) // 2 or 3
		if k > optCount {
			k = optCount
		}
		idxs := rand.Perm(optCount)[:k]
		sel := make([]string, 0, k)
		for _, id := range idxs {
			// set correct phrasing
			opts[id] = fmt.Sprintf("%s: 正确项（关于%s.%s的重要结论）", labels[id], chapter, topic)
			sel = append(sel, labels[id])
		}
		sort.Strings(sel)
		q.Answer = strings.Join(sel, "")
		q.Options = opts
		q.Explanation = fmt.Sprintf("本题考查 %s.%s 的相关知识。正确答案：%s。解析：%s 的关键点在于...（请在校稿时补充具体细节）。", topic, chapter, q.Answer, chapter)
		return
	}

	// single
	q.Options = opts
	q.Answer = labels[correctIdx]
	q.Explanation = fmt.Sprintf("本题考查 %s.%s 的概念。正确答案为 %s。解析：%s 的关键点在于...（请在校稿时补充具体细节）。", topic, chapter, q.Answer, chapter)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
