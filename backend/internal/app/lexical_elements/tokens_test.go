package lexical_elements

import "testing"

func TestDisplayTokens(t *testing.T) {
	output := captureOutput(DisplayTokens)
	if output == "" {
		t.Error("DisplayTokens() produced no output")
	}
}
