package lexical_elements

import "testing"

func TestDisplayImaginary(t *testing.T) {
	output := captureOutput(DisplayImaginary)
	if output == "" {
		t.Error("DisplayImaginary() produced no output")
	}
}
