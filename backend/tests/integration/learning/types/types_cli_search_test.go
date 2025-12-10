package types_test

import (
	"bytes"
	"strings"
	"testing"

	"go-study2/src/learning/types/cli"
)

// T023: CLI 搜索集成测试
func TestTypesCLISearch(t *testing.T) {
	inputs := strings.Join([]string{
		"search map key",
		"q",
	}, "\n") + "\n"

	stdin := bytes.NewBufferString(inputs)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cli.DisplayMenu(stdin, &stdout, &stderr)
	out := stdout.String()
	if !strings.Contains(out, "map") {
		t.Fatalf("未输出搜索结果: %s", out)
	}
}
