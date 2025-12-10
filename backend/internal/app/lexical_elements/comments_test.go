package lexical_elements

import "testing"

func TestDisplayComments(t *testing.T) {
	output := captureOutput(DisplayComments)
	if output == "" {
		t.Error("DisplayComments() produced no output")
	}
}
