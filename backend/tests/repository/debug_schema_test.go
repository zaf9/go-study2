package repository

import (
	"context"
	"testing"
)

func TestDebugTableSchema(t *testing.T) {
	db := newQuizTestDB(t)
	ensureUserTable(t, db)

	// Drop table if exists
	if _, err := db.Exec(context.Background(), "DROP TABLE IF EXISTS quiz_sessions"); err != nil {
		t.Fatalf("Failed to drop table: %v", err)
	}

	// Manually check what SQL is being executed
	stmt := `CREATE TABLE quiz_sessions (
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
    submitted_at DATETIME,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);`
	t.Logf("Executing SQL: %s", stmt)
	if _, err := db.Exec(context.Background(), stmt); err != nil {
		t.Fatalf("Failed to execute SQL: %v", err)
	}

	// Query table schema
	result, err := db.GetAll(context.Background(), "PRAGMA table_info(quiz_sessions)")
	if err != nil {
		t.Fatalf("Failed to get schema: %v", err)
	}

	t.Log("quiz_sessions schema:")
	for _, row := range result {
		t.Logf("  Column: %v", row.Map())
	}
}
