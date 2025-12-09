package types_test

import (
	"bytes"
	"strings"
	"testing"

	"go-study2/src/learning/types/cli"
)

// T017: CLI 综合测验评分与重做
func TestTypesCLIQuiz(t *testing.T) {
	inputs := strings.Join([]string{
		"quiz",
		"A", "B", "A", "A", "A",
		"n",
		"q",
	}, "\n") + "\n"

	stdin := bytes.NewBufferString(inputs)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cli.DisplayMenu(stdin, &stdout, &stderr)
	out := stdout.String()
	if !strings.Contains(out, "综合得分") {
		t.Fatalf("未输出综合测验得分: %s", out)
	}
}
