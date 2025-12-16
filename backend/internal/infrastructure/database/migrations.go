package database

import (
	"context"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
)

// Migrate 执行数据库迁移，创建核心表结构。
func Migrate(ctx context.Context, db gdb.DB) error {
	migrations := []string{
		createUsersTableSQL,
		createLearningProgressTableSQL,
		createQuizQuestionsTableSQL,
		createQuizSessionsTableSQL,
		createQuizAttemptsTableSQL,
		createQuizRecordsTableSQL,
		createRefreshTokensTableSQL,
		createAuditEventsTableSQL,
	}

	for _, stmt := range migrations {
		parts := strings.Split(stmt, ";")
		for _, sql := range parts {
			sql = strings.TrimSpace(sql)
			if sql == "" {
				continue
			}
			if _, err := db.Exec(ctx, sql); err != nil {
				return err
			}
		}
	}

	if err := ensureUserColumns(ctx, db); err != nil {
		return err
	}
	if err := ensureLearningProgressColumns(ctx, db); err != nil {
		return err
	}

	if err := seedDefaultQuizQuestions(ctx, db); err != nil {
		return err
	}

	return nil
}

const createUsersTableSQL = `
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_username ON users(username);
`

func ensureUserColumns(ctx context.Context, db gdb.DB) error {
	columns, err := db.GetAll(ctx, "PRAGMA table_info(users)")
	if err != nil {
		return err
	}
	has := func(name string) bool {
		for _, col := range columns {
			if strings.EqualFold(col["name"].String(), name) {
				return true
			}
		}
		return false
	}

	type columnDef struct {
		name string
		def  string
	}

	additions := []columnDef{
		{name: "is_admin", def: "INTEGER NOT NULL DEFAULT 0"},
		{name: "status", def: "TEXT NOT NULL DEFAULT 'active'"},
		{name: "must_change_password", def: "INTEGER NOT NULL DEFAULT 0"},
	}

	for _, col := range additions {
		if has(col.name) {
			continue
		}
		if _, err := db.Exec(ctx, fmt.Sprintf("ALTER TABLE users ADD COLUMN %s %s", col.name, col.def)); err != nil {
			return err
		}
	}
	return nil
}

func ensureLearningProgressColumns(ctx context.Context, db gdb.DB) error {
	columns, err := db.GetAll(ctx, "PRAGMA table_info(learning_progress)")
	if err != nil {
		return err
	}
	has := func(name string) bool {
		for _, col := range columns {
			if strings.EqualFold(col["name"].String(), name) {
				return true
			}
		}
		return false
	}

	type columnDef struct {
		name string
		def  string
	}

	additions := []columnDef{
		{name: "read_duration", def: "INTEGER NOT NULL DEFAULT 0"},
		{name: "scroll_progress", def: "INTEGER NOT NULL DEFAULT 0"},
		{name: "quiz_score", def: "INTEGER"},
		{name: "quiz_passed", def: "INTEGER NOT NULL DEFAULT 0"},
		{name: "first_visit_at", def: "DATETIME DEFAULT '1970-01-01 00:00:00'"},
		{name: "completed_at", def: "DATETIME"},
		{name: "updated_at", def: "DATETIME DEFAULT '1970-01-01 00:00:00'"},
	}

	// Rename last_visit to last_visit_at if exists
	if has("last_visit") && !has("last_visit_at") {
		if _, err := db.Exec(ctx, "ALTER TABLE learning_progress RENAME COLUMN last_visit TO last_visit_at"); err != nil {
			return err
		}
	}

	for _, col := range additions {
		if has(col.name) {
			continue
		}
		if _, err := db.Exec(ctx, fmt.Sprintf("ALTER TABLE learning_progress ADD COLUMN %s %s", col.name, col.def)); err != nil {
			return err
		}
	}
	return nil
}

func seedDefaultQuizQuestions(ctx context.Context, db gdb.DB) error {
	count, err := db.Model("quiz_questions").Count(ctx)
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	seed := []gdb.Map{
		{
			"topic":           "variables",
			"chapter":         "storage",
			"type":            "single",
			"difficulty":      "easy",
			"question":        "变量存储类型是？",
			"options":         `["栈","堆"]`,
			"correct_answers": `["A"]`,
			"explanation":     "示例题，校验接口链路与评分流程",
		},
	}
	_, err = db.Model("quiz_questions").Data(seed).Insert()
	return err
}

const createLearningProgressTableSQL = `
CREATE TABLE IF NOT EXISTS learning_progress (
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
    first_visit_at DATETIME DEFAULT '1970-01-01 00:00:00',
    last_visit_at DATETIME DEFAULT '1970-01-01 00:00:00',
    completed_at DATETIME,
    created_at DATETIME DEFAULT '1970-01-01 00:00:00',
    updated_at DATETIME DEFAULT '1970-01-01 00:00:00',
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(user_id, topic, chapter)
);
CREATE INDEX IF NOT EXISTS idx_progress_user_topic ON learning_progress(user_id, topic);
CREATE INDEX IF NOT EXISTS idx_progress_user_status ON learning_progress(user_id, status);
CREATE INDEX IF NOT EXISTS idx_progress_last_visit ON learning_progress(user_id, last_visit_at DESC);
`

const createQuizQuestionsTableSQL = `
CREATE TABLE IF NOT EXISTS quiz_questions (
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
);
CREATE INDEX IF NOT EXISTS idx_quiz_questions_topic_chapter ON quiz_questions(topic, chapter);
CREATE INDEX IF NOT EXISTS idx_quiz_questions_difficulty ON quiz_questions(difficulty);
`

const createQuizSessionsTableSQL = `
CREATE TABLE IF NOT EXISTS quiz_sessions (
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
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_quiz_sessions_user_completed ON quiz_sessions(user_id, completed_at DESC);
`

const createQuizAttemptsTableSQL = `
CREATE TABLE IF NOT EXISTS quiz_attempts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id TEXT NOT NULL,
    user_id INTEGER NOT NULL,
    topic TEXT NOT NULL,
    chapter TEXT NOT NULL,
    question_id INTEGER NOT NULL,
    user_answers TEXT NOT NULL,
    is_correct INTEGER NOT NULL DEFAULT 0 CHECK(is_correct IN (0,1)),
    attempted_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (session_id) REFERENCES quiz_sessions(session_id) ON DELETE CASCADE,
    FOREIGN KEY (question_id) REFERENCES quiz_questions(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_quiz_attempts_user ON quiz_attempts(user_id);
CREATE INDEX IF NOT EXISTS idx_quiz_attempts_question ON quiz_attempts(question_id);
`

const createQuizRecordsTableSQL = `
CREATE TABLE IF NOT EXISTS quiz_records (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    topic TEXT NOT NULL,
    chapter TEXT,
    score INTEGER NOT NULL,
    total INTEGER NOT NULL,
    duration_ms INTEGER NOT NULL,
    answers TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_quiz_user_topic ON quiz_records(user_id, topic);
CREATE INDEX IF NOT EXISTS idx_quiz_created_at ON quiz_records(user_id, created_at DESC);
`

const createRefreshTokensTableSQL = `
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    token_hash TEXT NOT NULL UNIQUE,
    expires_at DATETIME NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user ON refresh_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_expires ON refresh_tokens(expires_at);
`

const createAuditEventsTableSQL = `
CREATE TABLE IF NOT EXISTS audit_events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    event_type TEXT NOT NULL,
    user_id INTEGER,
    result TEXT NOT NULL,
    metadata TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_audit_events_type ON audit_events(event_type);
CREATE INDEX IF NOT EXISTS idx_audit_events_user ON audit_events(user_id);
CREATE INDEX IF NOT EXISTS idx_audit_events_created ON audit_events(created_at DESC);
`
