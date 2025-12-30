package domain

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	appquiz "go-study2/internal/app/quiz"
	quizdom "go-study2/internal/domain/quiz"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockQuizRepository implements appquiz.QuizRepository for testing
type mockQuizRepository struct {
	questions        []quizdom.QuizQuestion
	sessions         map[string]*quizdom.QuizSession
	attempts         []quizdom.QuizAttempt
	history          []quizdom.QuizSession
	getQuestionsErr  error
	createSessionErr error
	saveAttemptsErr  error
	getSessionErr    error
	updateResultErr  error
}

func newMockQuizRepository() *mockQuizRepository {
	return &mockQuizRepository{
		sessions: make(map[string]*quizdom.QuizSession),
		attempts: []quizdom.QuizAttempt{},
		history:  []quizdom.QuizSession{},
	}
}

func (m *mockQuizRepository) GetQuestionsByChapter(ctx context.Context, topic, chapter string) ([]quizdom.QuizQuestion, error) {
	if m.getQuestionsErr != nil {
		return nil, m.getQuestionsErr
	}
	return m.questions, nil
}

func (m *mockQuizRepository) CreateSession(ctx context.Context, session *quizdom.QuizSession) (string, error) {
	if m.createSessionErr != nil {
		return "", m.createSessionErr
	}
	// Generate a session ID if not provided
	if session.SessionID == "" {
		session.SessionID = "test-session-" + time.Now().Format("20060102150405")
	}
	session.CreatedAt = time.Now()
	m.sessions[session.SessionID] = session
	return session.SessionID, nil
}

func (m *mockQuizRepository) SaveAttempts(ctx context.Context, attempts []quizdom.QuizAttempt) error {
	if m.saveAttemptsErr != nil {
		return m.saveAttemptsErr
	}
	m.attempts = append(m.attempts, attempts...)
	return nil
}

func (m *mockQuizRepository) GetHistory(ctx context.Context, userID int64, topic string, limit int) ([]quizdom.QuizSession, error) {
	return m.history, nil
}

func (m *mockQuizRepository) GetSession(ctx context.Context, sessionID string) (*quizdom.QuizSession, error) {
	if m.getSessionErr != nil {
		return nil, m.getSessionErr
	}
	session, ok := m.sessions[sessionID]
	if !ok {
		return nil, nil
	}
	return session, nil
}

func (m *mockQuizRepository) UpdateSessionResult(ctx context.Context, sessionID string, correct int, score int, passed bool) error {
	if m.updateResultErr != nil {
		return m.updateResultErr
	}
	session, ok := m.sessions[sessionID]
	if !ok {
		return errors.New("session not found")
	}
	now := time.Now()
	session.CorrectAnswers = correct
	session.Score = score
	session.Passed = passed
	session.CompletedAt = &now
	return nil
}

func (m *mockQuizRepository) GetAttemptsBySession(ctx context.Context, sessionID string) ([]quizdom.QuizAttempt, error) {
	var result []quizdom.QuizAttempt
	for _, attempt := range m.attempts {
		if attempt.SessionID == sessionID {
			result = append(result, attempt)
		}
	}
	return result, nil
}

// Helper function to create test questions
func createTestQuestions() []quizdom.QuizQuestion {
	now := time.Now()
	return []quizdom.QuizQuestion{
		{
			ID:             1,
			Topic:          "variables",
			Chapter:        "storage",
			Type:           "single",
			Difficulty:     "easy",
			Question:       "What is the default value of int?",
			Options:        `["0","1","nil","undefined"]`,
			CorrectAnswers: `["A"]`,
			Explanation:    "The default value is 0",
			CreatedAt:      now,
			UpdatedAt:      now,
		},
		{
			ID:             2,
			Topic:          "variables",
			Chapter:        "storage",
			Type:           "multiple",
			Difficulty:     "medium",
			Question:       "Which are value types?",
			Options:        `["int","string","slice","array"]`,
			CorrectAnswers: `["A","D"]`,
			Explanation:    "int and array are value types",
			CreatedAt:      now,
			UpdatedAt:      now,
		},
		{
			ID:             3,
			Topic:          "variables",
			Chapter:        "storage",
			Type:           "single",
			Difficulty:     "hard",
			Question:       "What happens to unused variables?",
			Options:        `["compile error","runtime error","ignored","warning"]`,
			CorrectAnswers: `["A"]`,
			Explanation:    "Go does not allow unused variables",
			CreatedAt:      now,
			UpdatedAt:      now,
		},
		{
			ID:             4,
			Topic:          "variables",
			Chapter:        "storage",
			Type:           "single",
			Difficulty:     "easy",
			Question:       "Can you redeclare variables?",
			Options:        `["yes","no","sometimes","depends"]`,
			CorrectAnswers: `["B"]`,
			Explanation:    "Variables cannot be redeclared in the same scope",
			CreatedAt:      now,
			UpdatedAt:      now,
		},
		{
			ID:             5,
			Topic:          "variables",
			Chapter:        "storage",
			Type:           "multiple",
			Difficulty:     "medium",
			Question:       "Which are reference types?",
			Options:        `["slice","map","channel","struct"]`,
			CorrectAnswers: `["A","B","C"]`,
			Explanation:    "Slice, map, and channel are reference types",
			CreatedAt:      now,
			UpdatedAt:      now,
		},
		{
			ID:             6,
			Topic:          "variables",
			Chapter:        "storage",
			Type:           "single",
			Difficulty:     "easy",
			Question:       "What is the zero value for bool?",
			Options:        `["false","true","0","nil"]`,
			CorrectAnswers: `["A"]`,
			Explanation:    "The zero value for bool is false",
			CreatedAt:      now,
			UpdatedAt:      now,
		},
	}
}

// T059: Test GetOrCreateSession (via GetQuizQuestions)
func TestGetQuizQuestions_Success(t *testing.T) {
	ctx := context.Background()
	repo := newMockQuizRepository()
	repo.questions = createTestQuestions()

	svc := appquiz.NewService(repo)

	// Test getting quiz questions
	result, err := svc.GetQuizQuestions(ctx, 1, "variables", "storage")

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "variables", result.Topic)
	assert.Equal(t, "storage", result.Chapter)
	assert.NotEmpty(t, result.SessionID)
	assert.NotEmpty(t, result.Questions)

	// Session should be created
	assert.Contains(t, repo.sessions, result.SessionID)
	session := repo.sessions[result.SessionID]
	assert.Equal(t, int64(1), session.UserID)
	assert.Equal(t, "variables", session.Topic)
	assert.Equal(t, "storage", session.Chapter)
	assert.Greater(t, session.TotalQuestions, 0)
}

func TestGetQuizQuestions_InvalidInput(t *testing.T) {
	ctx := context.Background()
	repo := newMockQuizRepository()
	svc := appquiz.NewService(repo)

	tests := []struct {
		name    string
		userID  int64
		topic   string
		chapter string
	}{
		{"零用户ID", 0, "variables", "storage"},
		{"负用户ID", -1, "variables", "storage"},
		{"不支持的主题", 1, "unsupported", "chapter"},
		{"空章节", 1, "variables", ""},
		{"空主题", 1, "", "storage"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := svc.GetQuizQuestions(ctx, tt.userID, tt.topic, tt.chapter)
			assert.Error(t, err)
			assert.Nil(t, result)
			assert.Equal(t, appquiz.ErrInvalidInput, err)
		})
	}
}

func TestGetQuizQuestions_NoQuestionsAvailable(t *testing.T) {
	ctx := context.Background()
	repo := newMockQuizRepository()
	repo.questions = []quizdom.QuizQuestion{} // No questions

	svc := appquiz.NewService(repo)

	result, err := svc.GetQuizQuestions(ctx, 1, "variables", "storage")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, appquiz.ErrQuizUnavailable, err)
}

// T060: Test loadQuestions (indirectly via GetQuizQuestions)
func TestLoadQuestions_RepositoryError(t *testing.T) {
	ctx := context.Background()
	repo := newMockQuizRepository()
	repo.getQuestionsErr = errors.New("database error")

	svc := appquiz.NewService(repo)

	result, err := svc.GetQuizQuestions(ctx, 1, "variables", "storage")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "database error", err.Error())
}

// T061: Test SubmitAnswers (via SubmitQuiz)
func TestSubmitQuiz_AllCorrect(t *testing.T) {
	ctx := context.Background()
	repo := newMockQuizRepository()
	repo.questions = createTestQuestions()

	svc := appquiz.NewService(repo)

	// First, get questions to create a session
	payload, err := svc.GetQuizQuestions(ctx, 1, "variables", "storage")
	require.NoError(t, err)
	sessionID := payload.SessionID

	// Prepare correct answers based on the actual questions from repository
	// Parse correct answers from JSON
	correctAnswerMap := make(map[int64][]string)
	for _, q := range repo.questions {
		var answers []string
		err := json.Unmarshal([]byte(q.CorrectAnswers), &answers)
		require.NoError(t, err)
		correctAnswerMap[q.ID] = answers
	}

	// Create answers for the questions in the payload
	answers := make([]appquiz.AnswerSubmission, 0, len(payload.Questions))
	for _, q := range payload.Questions {
		if correctAns, ok := correctAnswerMap[q.ID]; ok {
			answers = append(answers, appquiz.AnswerSubmission{
				QuestionID:  q.ID,
				UserAnswers: correctAns,
			})
		}
	}

	// Submit quiz
	result, err := svc.SubmitQuiz(ctx, 1, sessionID, "variables", "storage", answers)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, len(payload.Questions), result.TotalQuestions)
	assert.Equal(t, len(payload.Questions), result.CorrectAnswers)
	assert.Equal(t, 100, result.Score)
	assert.True(t, result.Passed)

	// Verify session updated
	session := repo.sessions[sessionID]
	assert.NotNil(t, session.CompletedAt)
	assert.Equal(t, len(payload.Questions), session.CorrectAnswers)
	assert.True(t, session.Passed)
}

func TestSubmitQuiz_PartialCorrect(t *testing.T) {
	ctx := context.Background()
	repo := newMockQuizRepository()
	repo.questions = createTestQuestions()

	svc := appquiz.NewService(repo)

	// Get questions
	payload, err := svc.GetQuizQuestions(ctx, 1, "variables", "storage")
	require.NoError(t, err)
	sessionID := payload.SessionID

	// Prepare correct answers map
	correctAnswerMap := make(map[int64][]string)
	for _, q := range repo.questions {
		var answers []string
		err := json.Unmarshal([]byte(q.CorrectAnswers), &answers)
		require.NoError(t, err)
		correctAnswerMap[q.ID] = answers
	}

	// Prepare answers: half correct, half wrong
	answers := make([]appquiz.AnswerSubmission, 0, len(payload.Questions))
	for i, q := range payload.Questions {
		if i%2 == 0 {
			// Correct answer
			if correctAns, ok := correctAnswerMap[q.ID]; ok {
				answers = append(answers, appquiz.AnswerSubmission{
					QuestionID:  q.ID,
					UserAnswers: correctAns,
				})
			}
		} else {
			// Wrong answer
			answers = append(answers, appquiz.AnswerSubmission{
				QuestionID:  q.ID,
				UserAnswers: []string{"Z"}, // Wrong answer
			})
		}
	}

	// Submit quiz
	result, err := svc.SubmitQuiz(ctx, 1, sessionID, "variables", "storage", answers)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, len(payload.Questions), result.TotalQuestions)
	assert.Less(t, result.CorrectAnswers, result.TotalQuestions)
	assert.Greater(t, result.CorrectAnswers, 0)
}

func TestSubmitQuiz_DuplicateSubmit(t *testing.T) {
	ctx := context.Background()
	repo := newMockQuizRepository()
	repo.questions = createTestQuestions()

	svc := appquiz.NewService(repo)

	// Get questions
	payload, err := svc.GetQuizQuestions(ctx, 1, "variables", "storage")
	require.NoError(t, err)
	sessionID := payload.SessionID

	// Prepare correct answers
	correctAnswerMap := make(map[int64][]string)
	for _, q := range repo.questions {
		var answers []string
		err := json.Unmarshal([]byte(q.CorrectAnswers), &answers)
		require.NoError(t, err)
		correctAnswerMap[q.ID] = answers
	}

	answers := make([]appquiz.AnswerSubmission, 0, len(payload.Questions))
	for _, q := range payload.Questions {
		if correctAns, ok := correctAnswerMap[q.ID]; ok {
			answers = append(answers, appquiz.AnswerSubmission{
				QuestionID:  q.ID,
				UserAnswers: correctAns,
			})
		}
	}

	// First submission - should succeed
	result1, err1 := svc.SubmitQuiz(ctx, 1, sessionID, "variables", "storage", answers)
	require.NoError(t, err1)
	assert.NotNil(t, result1)

	// Second submission - should fail with duplicate error
	result2, err2 := svc.SubmitQuiz(ctx, 1, sessionID, "variables", "storage", answers)
	assert.Error(t, err2)
	assert.Nil(t, result2)
	assert.Equal(t, appquiz.ErrDuplicateSubmit, err2)
}

func TestSubmitQuiz_InvalidInput(t *testing.T) {
	ctx := context.Background()
	repo := newMockQuizRepository()
	repo.questions = createTestQuestions()

	svc := appquiz.NewService(repo)

	// Get a valid session first
	payload, err := svc.GetQuizQuestions(ctx, 1, "variables", "storage")
	require.NoError(t, err)
	sessionID := payload.SessionID

	answers := []appquiz.AnswerSubmission{
		{QuestionID: 1, UserAnswers: []string{"A"}},
	}

	tests := []struct {
		name      string
		userID    int64
		sessionID string
		topic     string
		chapter   string
	}{
		{"零用户ID", 0, sessionID, "variables", "storage"},
		{"空sessionID", 1, "", "variables", "storage"},
		{"不支持的主题", 1, sessionID, "unsupported", "storage"},
		{"空章节", 1, sessionID, "variables", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := svc.SubmitQuiz(ctx, tt.userID, tt.sessionID, tt.topic, tt.chapter, answers)
			assert.Error(t, err)
			assert.Nil(t, result)
		})
	}
}

func TestSubmitQuiz_SessionNotFound(t *testing.T) {
	ctx := context.Background()
	repo := newMockQuizRepository()
	svc := appquiz.NewService(repo)

	answers := []appquiz.AnswerSubmission{
		{QuestionID: 1, UserAnswers: []string{"A"}},
	}

	result, err := svc.SubmitQuiz(ctx, 1, "nonexistent-session", "variables", "storage", answers)

	assert.Error(t, err)
	assert.Nil(t, result)
}
