package lexical_elements

import "testing"

func TestDisplaySemicolons(t *testing.T) {
	output := captureOutput(DisplaySemicolons)
	if output == "" {
		t.Error("DisplaySemicolons() produced no output")
	}
}
