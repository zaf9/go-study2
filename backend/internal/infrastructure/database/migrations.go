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

const createLearningProgressTableSQL = `
CREATE TABLE IF NOT EXISTS learning_progress (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    topic TEXT NOT NULL,
    chapter TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'not_started',
    last_visit DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_position TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(user_id, topic, chapter)
);
CREATE INDEX IF NOT EXISTS idx_progress_user_topic ON learning_progress(user_id, topic);
CREATE INDEX IF NOT EXISTS idx_progress_last_visit ON learning_progress(user_id, last_visit DESC);
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
