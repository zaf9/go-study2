package quiz

import "testing"

func TestValidateQuestion(t *testing.T) {
	q := YAMLQuestion{
		ID:         "q1",
		Type:       "single",
		Difficulty: "easy",
		Stem:       "题干",
		Options:    []string{"A", "B"},
		Answer:     "A",
	}
	if err := ValidateQuestion(q); err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}
}

func TestValidateQuestionBad(t *testing.T) {
	q := YAMLQuestion{ID: "", Type: "foo", Difficulty: "x", Stem: "", Options: []string{"A"}, Answer: ""}
	if err := ValidateQuestion(q); err == nil {
		t.Fatalf("expected validation error")
	}
}
