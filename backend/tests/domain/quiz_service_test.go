package domain

import (
	"context"
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

	// 创建YAML仓储并添加测试题目（至少4单选+4多选以支持随机抽题）
	yamlRepo := quizdom.NewRepository()
	testYAMLQuestions := createTestYAMLQuestions()
	yamlRepo.AddBank("variables", "storage", testYAMLQuestions)

	svc := appquiz.NewService(yamlRepo, repo)

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
	yamlRepo := quizdom.NewRepository()
	svc := appquiz.NewService(yamlRepo, repo)

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
	yamlRepo := quizdom.NewRepository()
	// 不添加题目，模拟题库为空

	svc := appquiz.NewService(yamlRepo, repo)

	result, err := svc.GetQuizQuestions(ctx, 1, "variables", "storage")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, appquiz.ErrQuizUnavailable, err)
}

// T060: Test loadQuestions (indirectly via GetQuizQuestions)
func TestLoadQuestions_RepositoryError(t *testing.T) {
	ctx := context.Background()
	repo := newMockQuizRepository()
	yamlRepo := quizdom.NewRepository()

	svc := appquiz.NewService(yamlRepo, repo)
	// YAML repository doesn't have the topic/chapter, so it should return ErrQuizUnauthorized

	result, err := svc.GetQuizQuestions(ctx, 1, "variables", "storage")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, appquiz.ErrQuizUnavailable, err)
}

// T061: Test SubmitAnswers (via SubmitQuiz)
func TestSubmitQuiz_AllCorrect(t *testing.T) {
	ctx := context.Background()
	repo := newMockQuizRepository()

	// 创建YAML仓储并添加测试题目
	yamlRepo := quizdom.NewRepository()
	testYAMLQuestions := createTestYAMLQuestions()
	yamlRepo.AddBank("variables", "storage", testYAMLQuestions)

	svc := appquiz.NewService(yamlRepo, repo)

	// Get the actual questions from the service
	payload, err := svc.GetQuizQuestions(ctx, 1, "variables", "storage")
	require.NoError(t, err)
	sessionID := payload.SessionID

	// Build a map of the returned question IDs to their correct answers
	// Since we can't easily hash the YAML IDs, we'll use the Options field to identify questions
	correctAnswerMap := make(map[string]string)
	for _, yq := range testYAMLQuestions {
		// Create a unique key from options to identify the question
		key := yq.Stem
		correctAnswerMap[key] = yq.Answer
	}

	// Submit correct answers based on actual YAML correct answers
	answers := make([]appquiz.AnswerSubmission, 0, len(payload.Questions))
	for _, q := range payload.Questions {
		// Use the question text to look up the correct answer
		correctAnswer := correctAnswerMap[q.Question]
		// Convert answer string to individual characters (e.g., "AB" -> ["A", "B"])
		userAnswers := make([]string, 0, len(correctAnswer))
		for _, ch := range correctAnswer {
			userAnswers = append(userAnswers, string(ch))
		}
		answers = append(answers, appquiz.AnswerSubmission{
			QuestionID:  q.ID,
			UserAnswers: userAnswers,
		})
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

	yamlRepo := quizdom.NewRepository()
	testYAMLQuestions := createTestYAMLQuestions()
	yamlRepo.AddBank("variables", "storage", testYAMLQuestions)

	svc := appquiz.NewService(yamlRepo, repo)

	// Get questions
	payload, err := svc.GetQuizQuestions(ctx, 1, "variables", "storage")
	require.NoError(t, err)
	sessionID := payload.SessionID

	// Prepare answers: half correct, half wrong
	answers := make([]appquiz.AnswerSubmission, 0, len(payload.Questions))
	for i, q := range payload.Questions {
		if i%2 == 0 {
			// Correct answer
			if q.Type == "single" {
				answers = append(answers, appquiz.AnswerSubmission{
					QuestionID:  q.ID,
					UserAnswers: []string{"A"},
				})
			} else {
				answers = append(answers, appquiz.AnswerSubmission{
					QuestionID:  q.ID,
					UserAnswers: []string{"A", "D"},
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

	yamlRepo := quizdom.NewRepository()
	testYAMLQuestions := createTestYAMLQuestions()
	yamlRepo.AddBank("variables", "storage", testYAMLQuestions)

	svc := appquiz.NewService(yamlRepo, repo)

	// Get questions
	payload, err := svc.GetQuizQuestions(ctx, 1, "variables", "storage")
	require.NoError(t, err)
	sessionID := payload.SessionID

	// Prepare correct answers
	answers := make([]appquiz.AnswerSubmission, 0, len(payload.Questions))
	for _, q := range payload.Questions {
		if q.Type == "single" {
			answers = append(answers, appquiz.AnswerSubmission{
				QuestionID:  q.ID,
				UserAnswers: []string{"A"},
			})
		} else {
			answers = append(answers, appquiz.AnswerSubmission{
				QuestionID:  q.ID,
				UserAnswers: []string{"A", "D"},
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

	yamlRepo := quizdom.NewRepository()
	testYAMLQuestions := createTestYAMLQuestions()
	yamlRepo.AddBank("variables", "storage", testYAMLQuestions)

	svc := appquiz.NewService(yamlRepo, repo)

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
	yamlRepo := quizdom.NewRepository()
	svc := appquiz.NewService(yamlRepo, repo)

	answers := []appquiz.AnswerSubmission{
		{QuestionID: 1, UserAnswers: []string{"A"}},
	}

	result, err := svc.SubmitQuiz(ctx, 1, "nonexistent-session", "variables", "storage", answers)

	assert.Error(t, err)
	assert.Nil(t, result)
}

// createTestYAMLQuestions 创建YAML格式的测试题目
func createTestYAMLQuestions() []quizdom.YAMLQuestion {
	return []quizdom.YAMLQuestion{
		{ID: "yq1", Type: "single", Difficulty: "easy", Stem: "What is the default value of int?", Options: []string{"0", "1", "nil", "undefined"}, Answer: "A", Explanation: "The default value is 0", Topic: "variables", Chapter: "storage"},
		{ID: "yq2", Type: "single", Difficulty: "easy", Stem: "Single choice 2", Options: []string{"A1", "A2"}, Answer: "A", Explanation: "exp2", Topic: "variables", Chapter: "storage"},
		{ID: "yq3", Type: "single", Difficulty: "easy", Stem: "Single choice 3", Options: []string{"B1", "B2"}, Answer: "A", Explanation: "exp3", Topic: "variables", Chapter: "storage"},
		{ID: "yq4", Type: "single", Difficulty: "easy", Stem: "Single choice 4", Options: []string{"C1", "C2"}, Answer: "A", Explanation: "exp4", Topic: "variables", Chapter: "storage"},
		{ID: "yq5", Type: "multiple", Difficulty: "medium", Stem: "Which are value types?", Options: []string{"int", "string", "slice", "array"}, Answer: "AD", Explanation: "int and array are value types", Topic: "variables", Chapter: "storage"},
		{ID: "yq6", Type: "multiple", Difficulty: "medium", Stem: "Multiple choice 2", Options: []string{"A", "B", "C"}, Answer: "AB", Explanation: "exp6", Topic: "variables", Chapter: "storage"},
		{ID: "yq7", Type: "multiple", Difficulty: "medium", Stem: "Multiple choice 3", Options: []string{"D", "E", "F"}, Answer: "AB", Explanation: "exp7", Topic: "variables", Chapter: "storage"},
		{ID: "yq8", Type: "multiple", Difficulty: "medium", Stem: "Multiple choice 4", Options: []string{"G", "H", "I"}, Answer: "AB", Explanation: "exp8", Topic: "variables", Chapter: "storage"},
	}
}
