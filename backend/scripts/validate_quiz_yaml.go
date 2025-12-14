//go:build ignore
// +build ignore

package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

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

var (
	labelRe = regexp.MustCompile(`^[A-E]: `)
)

func main() {
	rootPtr := flag.String("root", "backend/quiz_data", "quiz_data root directory")
	failOnError := flag.Bool("fail", false, "exit with non-zero on validation failures")
	flag.Parse()

	root := *rootPtr
	var files []string
	_ = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if strings.HasSuffix(d.Name(), ".yaml") {
			files = append(files, path)
		}
		return nil
	})

	if len(files) == 0 {
		fmt.Fprintln(os.Stderr, "未找到任何 YAML 文件于", root)
		os.Exit(1)
	}

	head := fmt.Sprintf("| %-40s | %-5s | %-6s | %-6s | %-6s | %-8s | %-7s |\n", "File", "Count", "Single", "Multi", "Easy%", "Medium%", "Hard%")
	sep := strings.Repeat("-", len(head)-1) + "\n"
	fmt.Print(head)
	fmt.Print(sep)

	anyFail := false
	for _, f := range files {
		ok, msg := validateFile(f)
		if !ok {
			anyFail = true
		}
		fmt.Println(msg)
	}

	if anyFail && *failOnError {
		os.Exit(2)
	}
}

func validateFile(path string) (bool, string) {
	b, err := os.ReadFile(path)
	if err != nil {
		return false, fmt.Sprintf("%s -> ERROR reading: %v", short(path), err)
	}
	var bank QuizBank
	if err := yaml.Unmarshal(b, &bank); err != nil {
		return false, fmt.Sprintf("%s -> ERROR parsing YAML: %v", short(path), err)
	}
	qcount := len(bank.Questions)
	if qcount == 0 {
		return false, fmt.Sprintf("%s -> FAIL: 题目为空", short(path))
	}
	if qcount < 30 || qcount > 50 {
		// fail
		return false, fmt.Sprintf("%s -> FAIL: 题目数量 %d 不在 30-50 范围", short(path), qcount)
	}

	// collect stats
	single := 0
	multi := 0
	easy := 0
	med := 0
	hard := 0
	ids := map[string]int{}
	optionLabels := map[string]struct{}{}
	fileTopic, fileChapter := inferTopicChapter(path)
	var errorsList []string

	for i, q := range bank.Questions {
		idx := i + 1
		if q.ID == "" {
			errorsList = append(errorsList, fmt.Sprintf("#%d: ID 为空", idx))
		} else {
			if prev, ok := ids[q.ID]; ok {
				errorsList = append(errorsList, fmt.Sprintf("ID 重复: %s (第 %d 和 %d)", q.ID, prev, idx))
			} else {
				ids[q.ID] = idx
			}
		}
		if q.Type != "single" && q.Type != "multiple" {
			errorsList = append(errorsList, fmt.Sprintf("#%d: 无效题型: %s", idx, q.Type))
		} else if q.Type == "single" {
			single++
		} else {
			multi++
		}
		if q.Difficulty == "easy" {
			easy++
		} else if q.Difficulty == "medium" {
			med++
		} else if q.Difficulty == "hard" {
			hard++
		} else {
			errorsList = append(errorsList, fmt.Sprintf("#%d: 无效难度: %s", idx, q.Difficulty))
		}
		if strings.TrimSpace(q.Stem) == "" {
			errorsList = append(errorsList, fmt.Sprintf("#%d: 题干为空", idx))
		}
		if len(q.Options) < 2 || len(q.Options) > 5 {
			errorsList = append(errorsList, fmt.Sprintf("#%d: 选项数量 %d 超出 2-5", idx, len(q.Options)))
		}
		// validate option labels
		labels := map[string]bool{}
		for _, opt := range q.Options {
			if !labelRe.MatchString(opt) {
				errorsList = append(errorsList, fmt.Sprintf("#%d: 选项格式错误: %s", idx, opt))
				continue
			}
			label := string(opt[0])
			labels[label] = true
			optionLabels[label] = struct{}{}
		}
		// answer format
		if q.Type == "single" {
			if len(q.Answer) != 1 {
				errorsList = append(errorsList, fmt.Sprintf("#%d: 单选题答案应为单个字母: %s", idx, q.Answer))
			} else {
				if !labels[q.Answer] {
					errorsList = append(errorsList, fmt.Sprintf("#%d: 答案 %s 不在选项范围内", idx, q.Answer))
				}
			}
		} else {
			// multiple
			if len(q.Answer) < 2 || len(q.Answer) > 4 {
				errorsList = append(errorsList, fmt.Sprintf("#%d: 多选题答案字母数应为2-4: %s", idx, q.Answer))
			}
			seen := map[rune]bool{}
			for _, ch := range q.Answer {
				label := string(ch)
				if !labels[label] {
					errorsList = append(errorsList, fmt.Sprintf("#%d: 答案字母 %s 不在选项范围内", idx, label))
				}
				if seen[ch] {
					errorsList = append(errorsList, fmt.Sprintf("#%d: 答案字母重复: %s", idx, label))
				}
				seen[ch] = true
			}
		}
		if strings.TrimSpace(q.Explanation) == "" {
			errorsList = append(errorsList, fmt.Sprintf("#%d: 解析为空", idx))
		}
		// topic/chapter consistency
		if fileTopic != "" && fileTopic != q.Topic {
			errorsList = append(errorsList, fmt.Sprintf("#%d: topic 与文件路径不一致: %s vs %s", idx, q.Topic, fileTopic))
		}
		if fileChapter != "" && fileChapter != q.Chapter {
			errorsList = append(errorsList, fmt.Sprintf("#%d: chapter 与文件名不一致: %s vs %s", idx, q.Chapter, fileChapter))
		}
	}

	// build summary
	easyPct := pct(easy, qcount)
	medPct := pct(med, qcount)
	hardPct := pct(hard, qcount)

	status := "PASS"
	if len(errorsList) > 0 {
		status = "FAIL"
	}
	line := fmt.Sprintf("| %-40s | %-5d | %-6d | %-6d | %-6s | %-8s | %-7s |", short(path), qcount, single, multi, fmt.Sprintf("%d%%", easyPct), fmt.Sprintf("%d%%", medPct), fmt.Sprintf("%d%%", hardPct))
	if status == "FAIL" {
		line = line + "  <-- " + strings.Join(errorsList, "; ")
	}
	return status == "PASS", line
}

func inferTopicChapter(path string) (string, string) {
	// expect path .../quiz_data/{topic}/{chapter}.yaml
	parts := strings.Split(filepath.ToSlash(path), "/")
	for i := len(parts) - 1; i >= 0; i-- {
		if parts[i] == "quiz_data" && i+2 < len(parts) {
			return parts[i+1], strings.TrimSuffix(parts[i+2], ".yaml")
		}
	}
	// alternative: backend/backend/quiz_data
	for i := len(parts) - 1; i >= 0; i-- {
		if parts[i] == "backend" && i+3 < len(parts) && parts[i+1] == "quiz_data" {
			return parts[i+2], strings.TrimSuffix(parts[i+3], ".yaml")
		}
	}
	// fallback
	if len(parts) >= 2 {
		return parts[len(parts)-2], strings.TrimSuffix(parts[len(parts)-1], ".yaml")
	}
	return "", ""
}

func pct(a, b int) int {
	if b == 0 {
		return 0
	}
	return int(float64(a) / float64(b) * 100)
}

func short(p string) string {
	wd, _ := os.Getwd()
	r := strings.Replace(p, wd+string(os.PathSeparator), "", 1)
	return r
}
