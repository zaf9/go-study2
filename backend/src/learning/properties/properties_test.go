package properties

import (
	"testing"
)

func TestAllTopics(t *testing.T) {
	topics := AllTopics()
	if len(topics) != 7 {
		t.Errorf("Expected 7 topics, got %d", len(topics))
	}

	expectedTopics := []Topic{
		TopicRepresentation,
		TopicUnderlyingType,
		TopicCoreType,
		TopicTypeIdentity,
		TopicAssignability,
		TopicRepresentability,
		TopicMethodSet,
	}

	for i, topic := range topics {
		if topic != expectedTopics[i] {
			t.Errorf("Topic %d: expected %v, got %v", i, expectedTopics[i], topic)
		}
	}
}

func TestNormalizeTopic(t *testing.T) {
	tests := []struct {
		input    string
		expected Topic
	}{
		{"representation", TopicRepresentation},
		{"Representation", TopicRepresentation},
		{"  Representation  ", TopicRepresentation},
		{"underlying_type", TopicUnderlyingType},
		{"core_type", TopicCoreType},
		{"type_identity", TopicTypeIdentity},
		{"assignability", TopicAssignability},
		{"representability", TopicRepresentability},
		{"method_set", TopicMethodSet},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := NormalizeTopic(tt.input)
			if result != tt.expected {
				t.Errorf("NormalizeTopic(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsSupportedTopic(t *testing.T) {
	tests := []struct {
		topic Topic
		want  bool
	}{
		{TopicRepresentation, true},
		{TopicUnderlyingType, true},
		{TopicCoreType, true},
		{TopicTypeIdentity, true},
		{TopicAssignability, true},
		{TopicRepresentability, true},
		{TopicMethodSet, true},
		{Topic("invalid"), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.topic), func(t *testing.T) {
			got := IsSupportedTopic(tt.topic)
			if got != tt.want {
				t.Errorf("IsSupportedTopic(%v) = %v, want %v", tt.topic, got, tt.want)
			}
		})
	}
}

func TestLoadContent(t *testing.T) {
	topics := AllTopics()
	for _, topic := range topics {
		t.Run(string(topic), func(t *testing.T) {
			content, err := LoadContent(topic)
			if err != nil {
				t.Fatalf("LoadContent(%v) error = %v", topic, err)
			}

			if content.Concept.ID == "" {
				t.Error("Concept ID is empty")
			}
			if content.Concept.Title == "" {
				t.Error("Concept Title is empty")
			}
			if content.Concept.Summary == "" {
				t.Error("Concept Summary is empty")
			}
		})
	}
}

func TestLoadQuiz(t *testing.T) {
	topics := AllTopics()
	for _, topic := range topics {
		t.Run(string(topic), func(t *testing.T) {
			quiz, err := LoadQuiz(topic)
			if err != nil {
				t.Fatalf("LoadQuiz(%v) error = %v", topic, err)
			}

			if len(quiz) == 0 {
				t.Error("Quiz is empty")
			}

			for _, item := range quiz {
				if item.ID == "" {
					t.Error("Quiz item ID is empty")
				}
				if item.Stem == "" {
					t.Error("Quiz item Stem is empty")
				}
				if item.Answer == "" {
					t.Error("Quiz item Answer is empty")
				}
				if item.Explanation == "" {
					t.Error("Quiz item Explanation is empty")
				}
			}
		})
	}
}

func TestEvaluateQuiz(t *testing.T) {
	topic := TopicRepresentation
	quiz, err := LoadQuiz(topic)
	if err != nil {
		t.Fatalf("LoadQuiz() error = %v", err)
	}

	if len(quiz) == 0 {
		t.Skip("No quiz items to test")
	}

	// Test with correct answers
	answers := make(map[string]string)
	for _, item := range quiz {
		answers[item.ID] = item.Answer
	}

	result, err := EvaluateQuiz(topic, answers)
	if err != nil {
		t.Fatalf("EvaluateQuiz() error = %v", err)
	}

	if result.Score != result.Total {
		t.Errorf("Expected all correct, got %d/%d", result.Score, result.Total)
	}

	for _, detail := range result.Details {
		if !detail.Correct {
			t.Errorf("Item %s should be correct", detail.ID)
		}
	}
}

func TestSearchReferences(t *testing.T) {
	tests := []struct {
		keyword     string
		expectFound bool
	}{
		{"底层类型", true},
		{"可赋值性", true},
		{"方法集", true},
		{"nonexistent", false},
	}

	for _, tt := range tests {
		t.Run(tt.keyword, func(t *testing.T) {
			results, err := SearchReferences(tt.keyword)
			if err != nil {
				t.Fatalf("SearchReferences(%q) error = %v", tt.keyword, err)
			}

			found := len(results) > 0
			if found != tt.expectFound {
				t.Errorf("SearchReferences(%q) found = %v, want %v", tt.keyword, found, tt.expectFound)
			}

			if found {
				for _, result := range results {
					if result.Keyword == "" {
						t.Error("Result keyword is empty")
					}
					if result.Summary == "" {
						t.Error("Result summary is empty")
					}
				}
			}
		})
	}
}

func TestGetOverview(t *testing.T) {
	overview := GetOverview()

	if overview.Title == "" {
		t.Error("Overview title is empty")
	}
	if overview.Version == "" {
		t.Error("Overview version is empty")
	}
	if len(overview.Concepts) != 7 {
		t.Errorf("Expected 7 concepts, got %d", len(overview.Concepts))
	}
	if len(overview.Printable) == 0 {
		t.Error("Printable outline is empty")
	}
}

func TestLearningProgress(t *testing.T) {
	progress := NewProgress("user1")

	if progress.UserID != "user1" {
		t.Errorf("UserID = %v, want user1", progress.UserID)
	}
	if len(progress.CompletedConcepts) != 0 {
		t.Errorf("CompletedConcepts = %v, want empty", progress.CompletedConcepts)
	}

	// Test marking as completed
	progress.MarkCompleted(TopicRepresentation)
	if len(progress.CompletedConcepts) != 1 {
		t.Errorf("CompletedConcepts length = %v, want 1", len(progress.CompletedConcepts))
	}

	// Test recording quiz score
	result := QuizResult{Score: 3, Total: 3}
	progress.RecordQuizScore(TopicRepresentation, result)

	if progress.QuizScores == nil {
		t.Error("QuizScores is nil")
	}
}

func TestGetComprehensiveQuiz(t *testing.T) {
	quiz, err := GetComprehensiveQuiz()
	if err != nil {
		t.Fatalf("GetComprehensiveQuiz() error = %v", err)
	}

	if len(quiz) == 0 {
		t.Error("Comprehensive quiz is empty")
	}

	// Check that all questions are unique
	ids := make(map[string]bool)
	for _, item := range quiz {
		if ids[item.ID] {
			t.Errorf("Duplicate quiz item ID: %s", item.ID)
		}
		ids[item.ID] = true
	}
}
