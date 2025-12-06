package lexical_elements

import "testing"

func TestDisplayFloats(t *testing.T) {
	output := captureOutput(DisplayFloats)
	if output == "" {
		t.Error("DisplayFloats() produced no output")
	}
}
