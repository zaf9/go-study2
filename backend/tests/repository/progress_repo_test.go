package repository

import (
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"go-study2/internal/config"
	"go-study2/internal/domain/progress"
	infrarepo "go-study2/internal/infra/repository"
	infradb "go-study2/internal/infrastructure/database"

	_ "github.com/gogf/gf/contrib/drivers/sqlite/v2"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/os/gctx"
)

func TestProgressRepository_CreateOrUpdate(t *testing.T) {
	ctx := gctx.New()
	db := newTestDB(t)
	ensureUserTable(t, db)
	runFeatureMigrations(t, db)
	createUser(t, db, 1)

	repo := infrarepo.NewProgressRepository(db)
	record := &progress.LearningProgress{
		UserID:         1,
		Topic:          "variables",
		Chapter:        "storage",
		Status:         progress.StatusInProgress,
		ReadDuration:   120,
		ScrollProgress: 70,
		LastPosition:   "100",
	}
	if err := repo.CreateOrUpdate(ctx, record); err != nil {
		t.Fatalf("首写入失败: %v", err)
	}

	stored, err := repo.Get(ctx, 1, "variables", "storage")
	if err != nil {
		t.Fatalf("查询失败: %v", err)
	}
	if stored == nil || stored.Status != progress.StatusInProgress || stored.ScrollProgress != 70 {
		t.Fatalf("查询结果不符合预期: %+v", stored)
	}

	updated := &progress.LearningProgress{
		UserID:         1,
		Topic:          "variables",
		Chapter:        "storage",
		Status:         progress.StatusCompleted,
		ReadDuration:   30,
		ScrollProgress: 95,
		LastPosition:   "500",
		QuizScore:      90,
		QuizPassed:     true,
	}
	if err := repo.CreateOrUpdate(ctx, updated); err != nil {
		t.Fatalf("二次更新失败: %v", err)
	}

	stored, err = repo.Get(ctx, 1, "variables", "storage")
	if err != nil {
		t.Fatalf("查询失败: %v", err)
	}
	if stored.ReadDuration != 150 {
		t.Fatalf("时长应累加为 150，得到 %d", stored.ReadDuration)
	}
	if stored.Status != progress.StatusCompleted {
		t.Fatalf("状态应为 completed，得到 %s", stored.Status)
	}
	if !stored.QuizPassed || stored.QuizScore != 90 {
		t.Fatalf("测验结果未写入: %+v", stored)
	}
	if stored.CompletedAt == nil {
		t.Fatalf("完成时间应被设置")
	}

	all, err := repo.GetByUser(ctx, 1)
	if err != nil || len(all) != 1 {
		count, _ := db.Model("learning_progress").Where("user_id", 1).Count(ctx)
		t.Fatalf("GetByUser 返回异常: len=%d count=%d err=%v", len(all), count, err)
	}
	topicList, err := repo.GetByTopic(ctx, 1, "variables")
	if err != nil || len(topicList) != 1 {
		t.Fatalf("GetByTopic 返回异常: len=%d err=%v", len(topicList), err)
	}
}

func newTestDB(t *testing.T) gdb.DB {
	t.Helper()
	path := filepath.ToSlash(filepath.Join(t.TempDir(), "progress_repo.db"))
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

func ensureUserTable(t *testing.T, db gdb.DB) {
	t.Helper()
	_, err := db.Exec(gctx.New(), `
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);`)
	if err != nil {
		t.Fatalf("创建 users 表失败: %v", err)
	}
}

func runFeatureMigrations(t *testing.T, db gdb.DB) {
	t.Helper()
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS learning_progress (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    topic TEXT NOT NULL,
    chapter TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'not_started' CHECK(status IN ('not_started','in_progress','completed','tested')),
    read_duration INTEGER NOT NULL DEFAULT 0 CHECK(read_duration >= 0),
    scroll_progress INTEGER NOT NULL DEFAULT 0 CHECK(scroll_progress >= 0 AND scroll_progress <= 100),
    last_position INTEGER NOT NULL DEFAULT 0 CHECK(last_position >= 0),
    quiz_score INTEGER,
    quiz_passed INTEGER NOT NULL DEFAULT 0 CHECK(quiz_passed IN (0,1)),
    first_visit_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_visit_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at DATETIME,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, topic, chapter)
);`,
		`CREATE INDEX IF NOT EXISTS idx_learning_progress_user ON learning_progress(user_id);`,
		`CREATE INDEX IF NOT EXISTS idx_learning_progress_user_status ON learning_progress(user_id, status);`,
		`CREATE INDEX IF NOT EXISTS idx_learning_progress_user_last_visit ON learning_progress(user_id, last_visit_at DESC);`,
		`CREATE TABLE IF NOT EXISTS quiz_questions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    topic TEXT NOT NULL,
    chapter TEXT NOT NULL,
    type TEXT NOT NULL CHECK(type IN ('single','multiple','truefalse','code_output','code_correction')),
    difficulty TEXT NOT NULL CHECK(difficulty IN ('easy','medium','hard')),
    question TEXT NOT NULL,
    options TEXT NOT NULL,
    correct_answers TEXT NOT NULL,
    explanation TEXT NOT NULL,
    code_snippet TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);`,
		`CREATE INDEX IF NOT EXISTS idx_quiz_questions_topic_chapter ON quiz_questions(topic, chapter);`,
		`CREATE INDEX IF NOT EXISTS idx_quiz_questions_difficulty ON quiz_questions(difficulty);`,
		`CREATE TABLE IF NOT EXISTS quiz_sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id TEXT NOT NULL UNIQUE,
    user_id INTEGER NOT NULL,
    topic TEXT NOT NULL,
    chapter TEXT NOT NULL,
    total_questions INTEGER NOT NULL DEFAULT 0 CHECK(total_questions >= 0),
    correct_answers INTEGER NOT NULL DEFAULT 0 CHECK(correct_answers >= 0),
    score INTEGER NOT NULL DEFAULT 0 CHECK(score >= 0 AND score <= 100),
    passed INTEGER NOT NULL DEFAULT 0 CHECK(passed IN (0,1)),
    started_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at DATETIME,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);`,
		`CREATE INDEX IF NOT EXISTS idx_quiz_sessions_user_completed ON quiz_sessions(user_id, completed_at DESC);`,
		`CREATE TABLE IF NOT EXISTS quiz_attempts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id TEXT NOT NULL,
    user_id INTEGER NOT NULL,
    topic TEXT NOT NULL,
    chapter TEXT NOT NULL,
    question_id INTEGER NOT NULL,
    user_answers TEXT NOT NULL,
    is_correct INTEGER NOT NULL DEFAULT 0 CHECK(is_correct IN (0,1)),
    attempted_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);`,
		`CREATE INDEX IF NOT EXISTS idx_quiz_attempts_user ON quiz_attempts(user_id);`,
		`CREATE INDEX IF NOT EXISTS idx_quiz_attempts_question ON quiz_attempts(question_id);`,
	}
	for _, stmt := range stmts {
		if _, err := db.Exec(gctx.New(), stmt); err != nil {
			t.Fatalf("执行迁移失败: %v", err)
		}
	}
}

func createUser(t *testing.T, db gdb.DB, id int64) {
	t.Helper()
	_, err := db.Model("users").Data(gdb.Map{
		"id":            id,
		"username":      fmt.Sprintf("user_%d", id),
		"password_hash": "hash",
		"created_at":    time.Now(),
		"updated_at":    time.Now(),
	}).Insert()
	if err != nil {
		t.Fatalf("创建测试用户失败: %v", err)
	}
}
