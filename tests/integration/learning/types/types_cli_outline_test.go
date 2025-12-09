package types_test

import (
	"bytes"
	"strings"
	"testing"

	"go-study2/src/learning/types/cli"
)

// T035: CLI 提纲导出集成测试
func TestTypesCLIOutline(t *testing.T) {
	inputs := strings.Join([]string{
		"o",
		"q",
	}, "\n") + "\n"

	stdin := bytes.NewBufferString(inputs)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cli.DisplayMenu(stdin, &stdout, &stderr)
	out := stdout.String()
	if !strings.Contains(out, "提纲") {
		t.Fatalf("未输出提纲: %s", out)
	}
}
