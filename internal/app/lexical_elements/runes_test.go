package lexical_elements

import "testing"

func TestDisplayRunes(t *testing.T) {
	output := captureOutput(DisplayRunes)
	if output == "" {
		t.Error("DisplayRunes() produced no output")
	}
}
