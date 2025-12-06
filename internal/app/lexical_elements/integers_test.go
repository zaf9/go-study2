package lexical_elements

import "testing"

func TestDisplayIntegers(t *testing.T) {
	output := captureOutput(DisplayIntegers)
	if output == "" {
		t.Error("DisplayIntegers() produced no output")
	}
}
