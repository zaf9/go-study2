package repository

import (
	"path/filepath"
	"testing"
	"time"

	"go-study2/internal/config"
	"go-study2/internal/domain/quiz"
	infrarepo "go-study2/internal/infra/repository"
	infradb "go-study2/internal/infrastructure/database"

	_ "github.com/gogf/gf/contrib/drivers/sqlite/v2"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/os/gctx"
)

func TestQuizRepository_Flow(t *testing.T) {
	ctx := gctx.New()
	db := newQuizTestDB(t)
	ensureUserTable(t, db)
	runFeatureMigrations(t, db)
	createUser(t, db, 1)

	seedQuestions(t, db)

	repo := infrarepo.NewQuizRepository(db)

	questions, err := repo.GetQuestionsByChapter(ctx, "variables", "storage")
	if err != nil {
		t.Fatalf("获取题目失败: %v", err)
	}
	if len(questions) != 2 {
		count, _ := db.Model("quiz_questions").Count(ctx)
		t.Fatalf("应返回 2 道题目，得到 %d，当前题目总数 %d", len(questions), count)
	}

	session := &quiz.QuizSession{
		UserID:         1,
		Topic:          "variables",
		Chapter:        "storage",
		TotalQuestions: 2,
		CorrectAnswers: 1,
		Score:          50,
		Passed:         false,
	}
	sid, err := repo.CreateSession(ctx, session)
	if err != nil {
		t.Fatalf("创建测验会话失败: %v", err)
	}
	if sid == "" {
		t.Fatalf("sessionId 不应为空")
	}

	attempts := []quiz.QuizAttempt{
		{
			SessionID:   sid,
			UserID:      1,
			Topic:       "variables",
			Chapter:     "storage",
			QuestionID:  questions[0].ID,
			UserAnswers: `["A"]`,
			IsCorrect:   true,
		},
		{
			SessionID:   sid,
			UserID:      1,
			Topic:       "variables",
			Chapter:     "storage",
			QuestionID:  questions[1].ID,
			UserAnswers: `["B"]`,
			IsCorrect:   false,
		},
	}
	if err := repo.SaveAttempts(ctx, attempts); err != nil {
		t.Fatalf("保存作答失败: %v", err)
	}

	history, err := repo.GetHistory(ctx, 1, "variables", 10)
	if err != nil {
		t.Fatalf("获取历史失败: %v", err)
	}
	if len(history) != 1 {
		cnt, _ := db.Model("quiz_sessions").Count(ctx)
		t.Fatalf("历史记录数量应为 1，得到 %d，当前会话数 %d", len(history), cnt)
	}
	if history[0].Score != 50 || history[0].SessionID == "" {
		t.Fatalf("历史记录内容不正确: %+v", history[0])
	}
}

func newQuizTestDB(t *testing.T) gdb.DB {
	t.Helper()
	path := filepath.ToSlash(filepath.Join(t.TempDir(), "quiz_repo.db"))
	cfg := config.DatabaseConfig{
		Type: "sqlite3",
		Path: path,
		Pragmas: []string{
			"journal_mode=WAL",
			"busy_timeout=5000",
			"synchronous=NORMAL",
			"foreign_keys=ON",
		},
	}
	db, err := infradb.Init(gctx.New(), cfg)
	if err != nil {
		t.Fatalf("初始化测试数据库失败: %v", err)
	}
	t.Cleanup(func() {
		_ = db.Close(gctx.New())
	})
	return db
}

func seedQuestions(t *testing.T, db gdb.DB) {
	t.Helper()
	ctx := gctx.New()
	if _, err := db.Exec(ctx, "DELETE FROM quiz_questions"); err != nil {
		t.Fatalf("清理默认题库失败: %v", err)
	}
	now := time.Now()
	data := []gdb.Map{
		{
			"topic":           "variables",
			"chapter":         "storage",
			"type":            quiz.QuestionTypeSingle,
			"difficulty":      quiz.DifficultyEasy,
			"question":        "变量声明使用哪个关键字？",
			"options":         `["var","let","const","val"]`,
			"correct_answers": `["A"]`,
			"explanation":     "Go 使用 var 声明变量。",
			"created_at":      now,
			"updated_at":      now,
		},
		{
			"topic":           "variables",
			"chapter":         "storage",
			"type":            quiz.QuestionTypeTrueFalse,
			"difficulty":      quiz.DifficultyMedium,
			"question":        "短变量声明可用于包级作用域。",
			"options":         `["true","false"]`,
			"correct_answers": `["B"]`,
			"explanation":     "短变量声明只能在函数内使用。",
			"created_at":      now,
			"updated_at":      now,
		},
	}
	for _, row := range data {
		if _, err := db.Model("quiz_questions").Data(row).Insert(); err != nil {
			t.Fatalf("插入题目失败: %v", err)
		}
	}
}
