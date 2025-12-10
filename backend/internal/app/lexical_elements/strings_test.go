package lexical_elements

import "testing"

func TestDisplayStrings(t *testing.T) {
	output := captureOutput(DisplayStrings)
	if output == "" {
		t.Error("DisplayStrings() produced no output")
	}
}
