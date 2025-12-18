package quiz

// 模块说明：测验服务负责抽题、创建会话、幂等判分与历史查询，是题库与接口之间的业务桥梁。

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"math/big"
	"strings"
	"sync"
	"time"

	progressdom "go-study2/internal/domain/progress"
	quizdom "go-study2/internal/domain/quiz"
	"go-study2/internal/infrastructure/logger"
)

func init() {
}

// ErrInvalidInput 表示请求参数不合法。
var ErrInvalidInput = errors.New("测验参数不合法")

// ErrQuizUnavailable 表示当前章节暂无测验。
var ErrQuizUnavailable = errors.New("当前章节暂无测验")

// ErrDuplicateSubmit 表示重复提交同一 session。
var ErrDuplicateSubmit = errors.New("重复提交会话")

// QuizRepository 抽象测验持久化操作，便于替换与测试。
type QuizRepository interface {
	GetQuestionsByChapter(ctx context.Context, topic, chapter string) ([]quizdom.QuizQuestion, error)
	CreateSession(ctx context.Context, session *quizdom.QuizSession) (string, error)
	SaveAttempts(ctx context.Context, attempts []quizdom.QuizAttempt) error
	GetHistory(ctx context.Context, userID int64, topic string, limit int) ([]quizdom.QuizSession, error)
	GetSession(ctx context.Context, sessionID string) (*quizdom.QuizSession, error)
	UpdateSessionResult(ctx context.Context, sessionID string, correct int, score int, passed bool) error
	GetAttemptsBySession(ctx context.Context, sessionID string) ([]quizdom.QuizAttempt, error)
}

// Service 提供测验题目获取与提交判分。
type Service struct {
	repo       QuizRepository
	manager    *QuestionManager
	scorer     *ScoringEngine
	submitted  map[string]struct{}
	submitLock sync.Mutex
}

// NewService 创建测验服务。
func NewService(repo QuizRepository) *Service {
	return &Service{
		repo:      repo,
		manager:   NewQuestionManager(),
		scorer:    NewScoringEngine(),
		submitted: map[string]struct{}{},
	}
}

// QuizSessionPayload 返回题目与 session 信息。
type QuizSessionPayload struct {
	SessionID string        `json:"sessionId"`
	Topic     string        `json:"topic"`
	Chapter   string        `json:"chapter"`
	Questions []QuestionDTO `json:"questions"`
}

// QuizStats 表示章节统计信息（总量/按题型/按难度分布）
type QuizStats struct {
	Total        int            `json:"total"`
	ByType       map[string]int `json:"byType"`
	ByDifficulty map[string]int `json:"byDifficulty"`
}

// SubmitResult 返回测验提交结果。
type SubmitResult struct {
	Score          int            `json:"score"`
	TotalQuestions int            `json:"total_questions"`
	CorrectAnswers int            `json:"correct_answers"`
	Passed         bool           `json:"passed"`
	Details        []AnswerDetail `json:"details"`
}

// GetQuizQuestions 抽取题目并创建测验 session。
func (s *Service) GetQuizQuestions(ctx context.Context, userID int64, topic, chapter string) (*QuizSessionPayload, error) {
	topic = strings.TrimSpace(topic)
	chapter = strings.TrimSpace(chapter)
	if userID <= 0 || !quizdom.IsSupportedTopic(topic) || chapter == "" {
		return nil, ErrInvalidInput
	}

	records, err := s.repo.GetQuestionsByChapter(ctx, topic, chapter)
	if err != nil {
		return nil, err
	}
	prepared, _, err := s.manager.Prepare(records)
	if err != nil {
		return nil, err
	}

	// Shuffle the entire prepared pool to increase randomness
	for i := len(prepared) - 1; i > 0; i-- {
		j, _ := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		jInt := int(j.Int64())
		prepared[i], prepared[jInt] = prepared[jInt], prepared[i]
	}

	// select a randomized subset: pick a random first question from the whole pool to
	// increase diversity of the first item, then fill remaining slots respecting
	// single/multiple counts (3-5 each).
	var singles []PreparedQuestion
	var multiples []PreparedQuestion
	for idx, p := range prepared {
		// we will pick a random first later; collect lists with indices
		if p.View.Type == "single" {
			singles = append(singles, p)
		} else if p.View.Type == "multiple" {
			multiples = append(multiples, p)
		}
		_ = idx
	}
	totalPool := len(prepared)
	if totalPool == 0 {
		return nil, ErrQuizUnavailable
	}
	// choose a random first question from full prepared pool
	firstBig, _ := rand.Int(rand.Reader, big.NewInt(int64(totalPool)))
	firstIdx := int(firstBig.Int64())
	first := prepared[firstIdx]

	// remove first from the singles/multiples selection pools
	var remSingles []PreparedQuestion
	var remMultiples []PreparedQuestion
	for _, p := range prepared {
		if p.View.ID == first.View.ID {
			continue
		}
		if p.View.Type == "single" {
			remSingles = append(remSingles, p)
		} else if p.View.Type == "multiple" {
			remMultiples = append(remMultiples, p)
		}
	}

	singleBig, _ := rand.Int(rand.Reader, big.NewInt(3))
	multipleBig, _ := rand.Int(rand.Reader, big.NewInt(3))
	singleCount := 3 + int(singleBig.Int64())     // 3..5
	multipleCount := 3 + int(multipleBig.Int64()) // 3..5
	// adjust counts if first occupies one slot
	if first.View.Type == "single" {
		if singleCount > 0 {
			singleCount--
		}
	} else if first.View.Type == "multiple" {
		if multipleCount > 0 {
			multipleCount--
		}
	}
	if singleCount > len(remSingles) {
		singleCount = len(remSingles)
	}
	if multipleCount > len(remMultiples) {
		multipleCount = len(remMultiples)
	}
	// Shuffle remSingles and remMultiples using crypto/rand
	for i := len(remSingles) - 1; i > 0; i-- {
		j, _ := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		jInt := int(j.Int64())
		remSingles[i], remSingles[jInt] = remSingles[jInt], remSingles[i]
	}
	for i := len(remMultiples) - 1; i > 0; i-- {
		j, _ := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		jInt := int(j.Int64())
		remMultiples[i], remMultiples[jInt] = remMultiples[jInt], remMultiples[i]
	}

	selectedViews := []QuestionDTO{first.View}
	for i := 0; i < singleCount; i++ {
		selectedViews = append(selectedViews, remSingles[i].View)
	}
	for i := 0; i < multipleCount; i++ {
		selectedViews = append(selectedViews, remMultiples[i].View)
	}

	// final shuffle (keep the chosen first question at index 0, shuffle the rest)
	if len(selectedViews) > 1 {
		rest := selectedViews[1:]
		for i := len(rest) - 1; i > 0; i-- {
			j, _ := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
			jInt := int(j.Int64())
			rest[i], rest[jInt] = rest[jInt], rest[i]
		}
		// write back
		for i := 1; i < len(selectedViews); i++ {
			selectedViews[i] = rest[i-1]
		}
	}

	session := &quizdom.QuizSession{
		UserID:         userID,
		Topic:          topic,
		Chapter:        chapter,
		TotalQuestions: len(selectedViews),
		StartedAt:      time.Now(),
	}
	sessionID, err := s.repo.CreateSession(ctx, session)
	if err != nil {
		return nil, err
	}

	// 在返回给前端前，对每题的选项顺序进行随机打乱，避免答案位置规律
	// 使用 crypto/rand 保证并发与随机性的强度
	for i := range selectedViews {
		shuffleOptionsCrypto(&selectedViews[i].Options)
	}

	// 记录抽题操作到结构化日志
	logger.LogWithFields(ctx, "INFO", "quiz.selection", map[string]interface{}{
		"topic":      topic,
		"chapter":    chapter,
		"session_id": sessionID,
		"count":      len(selectedViews),
	})

	return &QuizSessionPayload{
		SessionID: sessionID,
		Topic:     topic,
		Chapter:   chapter,
		Questions: selectedViews,
	}, nil
}

// shuffleOptionsCrypto 使用 crypto/rand 对选项切片做 Fisher-Yates 洗牌
func shuffleOptionsCrypto(opts *[]OptionDTO) {
	if opts == nil || len(*opts) <= 1 {
		return
	}
	n := len(*opts)
	for i := n - 1; i > 0; i-- {
		jBig, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			// 若 crypto/rand 失败，退回到简单交换以保证功能不中断
			j := time.Now().UnixNano() % int64(i+1)
			(*opts)[i], (*opts)[j] = (*opts)[j], (*opts)[i]
			continue
		}
		j := int(jBig.Int64())
		(*opts)[i], (*opts)[j] = (*opts)[j], (*opts)[i]
	}
}

// SubmitQuiz 评判答案并写入会话与作答记录。
func (s *Service) SubmitQuiz(ctx context.Context, userID int64, sessionID, topic, chapter string, answers []AnswerSubmission) (*SubmitResult, error) {
	topic = strings.TrimSpace(topic)
	chapter = strings.TrimSpace(chapter)
	if userID <= 0 || sessionID == "" || !quizdom.IsSupportedTopic(topic) || chapter == "" {
		logger.LogWithFields(ctx, "WARNING", "quiz.submit.invalid_input", map[string]interface{}{
			"user_id":    userID,
			"session_id": sessionID,
			"topic":      topic,
			"chapter":    chapter,
			"reason":     "参数验证失败",
		})
		return nil, ErrInvalidInput
	}

	// 先检查数据库中的 session 状态，避免内存检查的竞态条件
	existing, err := s.repo.GetSession(ctx, sessionID)
	if err != nil {
		logger.LogWithFields(ctx, "ERROR", "quiz.submit.get_session_failed", map[string]interface{}{
			"user_id":    userID,
			"session_id": sessionID,
			"topic":      topic,
			"chapter":    chapter,
			"error":      err.Error(),
		})
		return nil, err
	}
	if existing == nil {
		logger.LogWithFields(ctx, "WARNING", "quiz.submit.session_not_found", map[string]interface{}{
			"user_id":    userID,
			"session_id": sessionID,
			"topic":      topic,
			"chapter":    chapter,
		})
		return nil, ErrInvalidInput
	}
	// 如果 session 已经完成，直接返回重复提交错误
	if existing.CompletedAt != nil {
		logger.LogWithFields(ctx, "WARNING", "quiz.submit.duplicate_submit", map[string]interface{}{
			"user_id":      userID,
			"session_id":   sessionID,
			"topic":        topic,
			"chapter":      chapter,
			"completed_at": existing.CompletedAt.Format(time.RFC3339),
			"reason":       "session已完成",
		})
		return nil, ErrDuplicateSubmit
	}

	// 使用内存锁防止并发重复提交
	s.submitLock.Lock()
	if _, ok := s.submitted[sessionID]; ok {
		s.submitLock.Unlock()
		logger.LogWithFields(ctx, "WARNING", "quiz.submit.duplicate_submit", map[string]interface{}{
			"user_id":    userID,
			"session_id": sessionID,
			"topic":      topic,
			"chapter":    chapter,
			"reason":     "内存中检测到重复提交",
		})
		return nil, ErrDuplicateSubmit
	}
	s.submitted[sessionID] = struct{}{}
	s.submitLock.Unlock()

	// 确保在处理完成后释放内存锁，防止内存泄漏和误判
	defer func() {
		s.submitLock.Lock()
		delete(s.submitted, sessionID)
		s.submitLock.Unlock()
	}()

	// 获取实际测试的题目数量
	actualTotalQuestions := existing.TotalQuestions
	if actualTotalQuestions <= 0 {
		logger.LogWithFields(ctx, "WARNING", "quiz.submit.invalid_total_questions", map[string]interface{}{
			"user_id":         userID,
			"session_id":      sessionID,
			"topic":           topic,
			"chapter":         chapter,
			"total_questions": actualTotalQuestions,
		})
		return nil, ErrInvalidInput
	}

	// 构建用户答案映射
	answerMap := map[int64][]string{}
	questionIDs := make(map[int64]struct{})
	for _, a := range answers {
		answerMap[a.QuestionID] = normalizeChoices(a.UserAnswers)
		questionIDs[a.QuestionID] = struct{}{}
	}

	// 验证提交的答案数量是否与实际测试题目数量一致
	// 如果不一致，说明有些题目没有提交答案（这些题目将得0分）
	submittedCount := len(answers)
	if submittedCount > actualTotalQuestions {
		logger.LogWithFields(ctx, "WARNING", "quiz.submit.answer_count_mismatch", map[string]interface{}{
			"user_id":         userID,
			"session_id":      sessionID,
			"topic":           topic,
			"chapter":         chapter,
			"submitted_count": submittedCount,
			"expected_count":  actualTotalQuestions,
		})
		return nil, ErrInvalidInput
	}

	// 只获取用户提交答案对应的题目，而不是章节的所有题目
	records, err := s.repo.GetQuestionsByChapter(ctx, topic, chapter)
	if err != nil {
		logger.LogWithFields(ctx, "ERROR", "quiz.submit.get_questions_failed", map[string]interface{}{
			"user_id":    userID,
			"session_id": sessionID,
			"topic":      topic,
			"chapter":    chapter,
			"error":      err.Error(),
		})
		return nil, err
	}
	// 筛选出实际测试的题目（用户提交答案的题目）
	var testRecords []quizdom.QuizQuestion
	for _, r := range records {
		if _, ok := questionIDs[r.ID]; ok {
			testRecords = append(testRecords, r)
		}
	}
	if len(testRecords) == 0 {
		logger.LogWithFields(ctx, "WARNING", "quiz.submit.no_test_records", map[string]interface{}{
			"user_id":         userID,
			"session_id":      sessionID,
			"topic":           topic,
			"chapter":         chapter,
			"submitted_count": submittedCount,
			"total_records":   len(records),
		})
		return nil, ErrInvalidInput
	}

	prepared, _, err := s.manager.Prepare(testRecords)
	if err != nil {
		logger.LogWithFields(ctx, "ERROR", "quiz.submit.prepare_questions_failed", map[string]interface{}{
			"user_id":    userID,
			"session_id": sessionID,
			"topic":      topic,
			"chapter":    chapter,
			"error":      err.Error(),
		})
		return nil, err
	}

	// 使用实际测试的题目数量进行评分
	// 如果用户提交的答案数量少于实际测试题目数量，未提交的题目将得0分
	score := s.scorer.EvaluateWithTotal(prepared, answerMap, actualTotalQuestions)

	logger.LogWithFields(ctx, "INFO", "quiz.submit.scoring_completed", map[string]interface{}{
		"user_id":         userID,
		"session_id":      sessionID,
		"topic":           topic,
		"chapter":         chapter,
		"score":           score.Score,
		"total_questions": score.TotalQuestions,
		"correct_answers": score.CorrectAnswers,
		"submitted_count": submittedCount,
		"actual_total":    actualTotalQuestions,
	})

	var attempts []quizdom.QuizAttempt
	for _, detail := range score.Details {
		rawAns, _ := json.Marshal(answerMap[detail.QuestionID])
		attempts = append(attempts, quizdom.QuizAttempt{
			SessionID:   sessionID,
			UserID:      userID,
			Topic:       topic,
			Chapter:     chapter,
			QuestionID:  detail.QuestionID,
			UserAnswers: string(rawAns),
			IsCorrect:   detail.IsCorrect,
			AttemptedAt: time.Now(),
		})
	}
	if err := s.repo.SaveAttempts(ctx, attempts); err != nil {
		logger.LogWithFields(ctx, "ERROR", "quiz.submit.save_attempts_failed", map[string]interface{}{
			"user_id":    userID,
			"session_id": sessionID,
			"topic":      topic,
			"chapter":    chapter,
			"attempts":   len(attempts),
			"error":      err.Error(),
		})
		return nil, err
	}
	if err := s.repo.UpdateSessionResult(ctx, sessionID, score.CorrectAnswers, score.Score, score.Passed); err != nil {
		logger.LogWithFields(ctx, "ERROR", "quiz.submit.update_session_failed", map[string]interface{}{
			"user_id":         userID,
			"session_id":      sessionID,
			"topic":           topic,
			"chapter":         chapter,
			"score":           score.Score,
			"correct_answers": score.CorrectAnswers,
			"error":           err.Error(),
		})
		return nil, err
	}

	logger.LogWithFields(ctx, "INFO", "quiz.submit.success", map[string]interface{}{
		"user_id":         userID,
		"session_id":      sessionID,
		"topic":           topic,
		"chapter":         chapter,
		"score":           score.Score,
		"total_questions": score.TotalQuestions,
		"correct_answers": score.CorrectAnswers,
		"passed":          score.Passed,
	})

	return &SubmitResult{
		Score:          score.Score,
		TotalQuestions: score.TotalQuestions,
		CorrectAnswers: score.CorrectAnswers,
		Passed:         score.Passed,
		Details:        score.Details,
	}, nil
}

// GetQuizHistory 返回用户测验历史。
func (s *Service) GetQuizHistory(ctx context.Context, userID int64, topic string, limit int) ([]quizdom.QuizSession, error) {
	if userID <= 0 {
		return nil, ErrInvalidInput
	}
	if topic != "" && !quizdom.IsSupportedTopic(topic) {
		return nil, ErrInvalidInput
	}
	return s.repo.GetHistory(ctx, userID, strings.TrimSpace(topic), limit)
}

// GetStats 返回指定 topic/chapter 的题库统计信息
func (s *Service) GetStats(ctx context.Context, topic, chapter string) (*QuizStats, error) {
	topic = strings.TrimSpace(topic)
	chapter = strings.TrimSpace(chapter)
	if !quizdom.IsSupportedTopic(topic) || chapter == "" {
		return nil, ErrInvalidInput
	}
	records, err := s.repo.GetQuestionsByChapter(ctx, topic, chapter)
	if err != nil {
		return nil, err
	}
	stats := &QuizStats{
		Total:        len(records),
		ByType:       map[string]int{},
		ByDifficulty: map[string]int{},
	}
	for _, r := range records {
		stats.ByType[r.Type]++
		stats.ByDifficulty[r.Difficulty]++
	}
	return stats, nil
}

// QuizReviewDetail 表示测验回顾详情。
type QuizReviewDetail struct {
	Meta  *QuizReviewMeta  `json:"meta"`
	Items []QuizReviewItem `json:"items"`
}

// QuizReviewMeta 表示回顾的会话元信息。
type QuizReviewMeta struct {
	SessionID   string     `json:"sessionId"`
	Topic       string     `json:"topic"`
	Chapter     string     `json:"chapter"`
	Score       int        `json:"score"`
	Passed      bool       `json:"passed"`
	CompletedAt *time.Time `json:"completedAt"`
}

// QuizReviewItem 表示回顾中的单题详情。
type QuizReviewItem struct {
	QuestionID    int64    `json:"questionId"`
	Stem          string   `json:"stem"`
	Options       []string `json:"options"`
	UserChoice    string   `json:"userChoice"`
	CorrectChoice string   `json:"correctChoice"`
	IsCorrect     bool     `json:"isCorrect"`
	Explanation   string   `json:"explanation"`
}

// GetQuizReview 获取指定会话的回顾详情。
func (s *Service) GetQuizReview(ctx context.Context, userID int64, sessionID string) (*QuizReviewDetail, error) {
	if userID <= 0 || sessionID == "" {
		return nil, ErrInvalidInput
	}

	// 获取会话信息
	session, err := s.repo.GetSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, ErrInvalidInput
	}

	// 验证会话属于当前用户
	if session.UserID != userID {
		return nil, ErrInvalidInput
	}

	// 获取答题记录
	attempts, err := s.repo.GetAttemptsBySession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	// 获取题目信息
	records, err := s.repo.GetQuestionsByChapter(ctx, session.Topic, session.Chapter)
	if err != nil {
		return nil, err
	}

	// 构建题目映射
	questionMap := make(map[int64]quizdom.QuizQuestion)
	for _, r := range records {
		questionMap[r.ID] = r
	}

	// 构建回顾详情
	meta := &QuizReviewMeta{
		SessionID:   session.SessionID,
		Topic:       session.Topic,
		Chapter:     session.Chapter,
		Score:       session.Score,
		Passed:      session.Passed,
		CompletedAt: session.CompletedAt,
	}

	var items = make([]QuizReviewItem, 0)
	for _, attempt := range attempts {
		question, ok := questionMap[attempt.QuestionID]
		if !ok {
			continue
		}

		// 解析选项
		var options []string
		if err := json.Unmarshal([]byte(question.Options), &options); err != nil {
			continue
		}

		// 解析正确答案
		var correctAnswers []string
		if err := json.Unmarshal([]byte(question.CorrectAnswers), &correctAnswers); err != nil {
			continue
		}

		// 解析用户答案
		var userAnswers []string
		if err := json.Unmarshal([]byte(attempt.UserAnswers), &userAnswers); err != nil {
			continue
		}

		// 将答案标签转换为选项内容
		userChoice := ""
		if len(userAnswers) > 0 {
			userChoiceContents := make([]string, 0, len(userAnswers))
			for _, ans := range userAnswers {
				ans = strings.ToUpper(strings.TrimSpace(ans))
				// 尝试将标签（如 "A", "B"）转换为选项内容
				if idx := strings.Index("ABCDEFGHIJKLMNOPQRSTUVWXYZ", ans); idx >= 0 && idx < len(options) {
					userChoiceContents = append(userChoiceContents, options[idx])
				} else {
					// 如果已经是内容，直接使用
					userChoiceContents = append(userChoiceContents, ans)
				}
			}
			userChoice = strings.Join(userChoiceContents, ",")
		}

		correctChoice := ""
		if len(correctAnswers) > 0 {
			// correctAnswers 存储的是选项标签（如 "A", "B"）
			correctContents := make([]string, 0, len(correctAnswers))
			for _, ans := range correctAnswers {
				ans = strings.ToUpper(strings.TrimSpace(ans))
				if idx := strings.Index("ABCDEFGHIJKLMNOPQRSTUVWXYZ", ans); idx >= 0 && idx < len(options) {
					correctContents = append(correctContents, options[idx])
				}
			}
			correctChoice = strings.Join(correctContents, ",")
		}

		items = append(items, QuizReviewItem{
			QuestionID:    question.ID,
			Stem:          question.Question,
			Options:       options,
			UserChoice:    userChoice,
			CorrectChoice: correctChoice,
			IsCorrect:     attempt.IsCorrect,
			Explanation:   question.Explanation,
		})
	}

	return &QuizReviewDetail{
		Meta:  meta,
		Items: items,
	}, nil
}

// ConvertProgressStatus 返回测验通过后的进度状态，供后续拓展。
func ConvertProgressStatus(passed bool) string {
	if passed {
		return progressdom.StatusCompleted
	}
	return progressdom.StatusTested
}
