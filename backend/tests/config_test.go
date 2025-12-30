package tests

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SetupTestDB 初始化内存数据库用于测试
// 创建一个SQLite内存数据库实例,用于单元测试和集成测试
func SetupTestDB(t *testing.T) gdb.DB {
	var (
		ctx  = context.Background()
		link = "sqlite::memory:"
	)

	// 创建内存数据库配置
	config := gdb.ConfigNode{
		Link: link,
		Type: "sqlite",
	}

	// 设置数据库配置
	gdb.SetConfig(gdb.Config{
		"default": gdb.ConfigGroup{
			config,
		},
	})

	// 获取数据库实例
	db := g.DB()

	// 创建测试用表结构
	setupTestSchema(t, ctx, db)

	return db
}

// setupTestSchema 创建测试用表结构
func setupTestSchema(t *testing.T, ctx context.Context, db gdb.DB) {
	// 创建 learning_progress 表
	_, err := db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS learning_progress (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			topic VARCHAR(100) NOT NULL,
			chapter VARCHAR(100) NOT NULL,
			status VARCHAR(20) DEFAULT 'in_progress',
			started_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			completed_at DATETIME,
			last_accessed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		t.Fatalf("创建 learning_progress 表失败: %v", err)
	}

	// 创建 quiz_sessions 表
	_, err = db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS quiz_sessions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			session_id VARCHAR(100) UNIQUE NOT NULL,
			user_id INTEGER NOT NULL,
			topic VARCHAR(100) NOT NULL,
			chapter VARCHAR(100) NOT NULL,
			score INTEGER DEFAULT 0,
			total_questions INTEGER DEFAULT 0,
			correct_answers INTEGER DEFAULT 0,
			submitted_at DATETIME,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		t.Fatalf("创建 quiz_sessions 表失败: %v", err)
	}

	// 创建 quiz_attempts 表
	_, err = db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS quiz_attempts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			session_id VARCHAR(100) NOT NULL,
			question_id INTEGER NOT NULL,
			user_answer TEXT,
			is_correct BOOLEAN DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		t.Fatalf("创建 quiz_attempts 表失败: %v", err)
	}

	// 创建索引
	db.Exec(ctx, `CREATE INDEX IF NOT EXISTS idx_progress_user_topic ON learning_progress(user_id, topic)`)
	db.Exec(ctx, `CREATE INDEX IF NOT EXISTS idx_quiz_session_user ON quiz_sessions(user_id)`)
}

// CleanupTestDB 清理测试数据库
func CleanupTestDB(t *testing.T, db gdb.DB) {
	ctx := context.Background()

	// 清空所有测试表
	tables := []string{"learning_progress", "quiz_sessions", "quiz_attempts"}
	for _, table := range tables {
		_, err := db.Exec(ctx, "DELETE FROM "+table)
		if err != nil {
			t.Logf("清空表 %s 失败: %v", table, err)
		}
	}
}
