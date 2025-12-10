package lexical_elements

import "testing"

func TestDisplayKeywords(t *testing.T) {
	output := captureOutput(DisplayKeywords)
	if output == "" {
		t.Error("DisplayKeywords() produced no output")
	}
}
