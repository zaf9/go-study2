package quiz

import (
	"fmt"
	"testing"
)

func TestSelectQuestions(t *testing.T) {
	qs := []YAMLQuestion{}
	for i := 0; i < 10; i++ {
		qs = append(qs, YAMLQuestion{ID: fmt.Sprintf("s%d", i), Type: "single", Difficulty: "easy", Stem: "x", Options: []string{"A", "B"}, Answer: "A"})
	}
	for i := 0; i < 8; i++ {
		qs = append(qs, YAMLQuestion{ID: fmt.Sprintf("m%d", i), Type: "multiple", Difficulty: "easy", Stem: "x", Options: []string{"A", "B", "C"}, Answer: "AB"})
	}
	sel, err := SelectQuestions(qs, 3, 3, nil)
	if err != nil {
		t.Fatalf("select error: %v", err)
	}
	if len(sel) != 6 {
		t.Fatalf("expected 6 selected, got %d", len(sel))
	}
}
