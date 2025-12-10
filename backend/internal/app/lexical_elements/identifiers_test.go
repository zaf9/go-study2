package lexical_elements

import "testing"

func TestDisplayIdentifiers(t *testing.T) {
	output := captureOutput(DisplayIdentifiers)
	if output == "" {
		t.Error("DisplayIdentifiers() produced no output")
	}
}
