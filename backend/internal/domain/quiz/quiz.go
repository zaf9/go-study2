package quiz

import "time"

// QuizQuestion 描述题库中的单个题目。
type QuizQuestion struct {
	ID             int64     `json:"id" orm:"id"`
	Topic          string    `json:"topic" orm:"topic"`
	Chapter        string    `json:"chapter" orm:"chapter"`
	Type           string    `json:"type" orm:"type"`
	Difficulty     string    `json:"difficulty" orm:"difficulty"`
	Question       string    `json:"question" orm:"question"`
	Options        string    `json:"options" orm:"options"`
	CorrectAnswers string    `json:"correctAnswers" orm:"correct_answers"`
	Explanation    string    `json:"explanation" orm:"explanation"`
	CodeSnippet    *string   `json:"codeSnippet" orm:"code_snippet"`
	CreatedAt      time.Time `json:"createdAt" orm:"created_at"`
	UpdatedAt      time.Time `json:"updatedAt" orm:"updated_at"`
}

// QuizSession 记录一次测验的摘要信息。
type QuizSession struct {
	ID             int64      `json:"id" orm:"id"`
	SessionID      string     `json:"sessionId" orm:"session_id"`
	UserID         int64      `json:"userId" orm:"user_id"`
	Topic          string     `json:"topic" orm:"topic"`
	Chapter        string     `json:"chapter" orm:"chapter"`
	TotalQuestions int        `json:"totalQuestions" orm:"total_questions"`
	CorrectAnswers int        `json:"correctAnswers" orm:"correct_answers"`
	Score          int        `json:"score" orm:"score"`
	Passed         bool       `json:"passed" orm:"passed"`
	StartedAt      time.Time  `json:"startedAt" orm:"started_at"`
	CompletedAt    *time.Time `json:"completedAt" orm:"completed_at"`
	CreatedAt      time.Time  `json:"createdAt" orm:"created_at"`
}

// QuizAttempt 记录用户对单题的作答。
type QuizAttempt struct {
	ID          int64     `json:"id" orm:"id"`
	SessionID   string    `json:"sessionId" orm:"session_id"`
	UserID      int64     `json:"userId" orm:"user_id"`
	Topic       string    `json:"topic" orm:"topic"`
	Chapter     string    `json:"chapter" orm:"chapter"`
	QuestionID  int64     `json:"questionId" orm:"question_id"`
	UserAnswers string    `json:"userAnswers" orm:"user_answers"`
	IsCorrect   bool      `json:"isCorrect" orm:"is_correct"`
	AttemptedAt time.Time `json:"attemptedAt" orm:"attempted_at"`
}

// 题型与难度枚举。
const (
	QuestionTypeSingle         = "single"
	QuestionTypeMultiple       = "multiple"
	QuestionTypeTrueFalse      = "truefalse"
	QuestionTypeCodeOutput     = "code_output"
	QuestionTypeCodeCorrection = "code_correction"

	DifficultyEasy   = "easy"
	DifficultyMedium = "medium"
	DifficultyHard   = "hard"
)
