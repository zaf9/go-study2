package lexical_elements

import "testing"

func TestDisplayOperators(t *testing.T) {
	output := captureOutput(DisplayOperators)
	if output == "" {
		t.Error("DisplayOperators() produced no output")
	}
}
