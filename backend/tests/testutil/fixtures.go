package testutil

import (
	"context"
	"testing"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
)

// CreateTestUser 创建测试用户数据
func CreateTestUser(t *testing.T, db gdb.DB, userID int64) {
	ctx := context.Background()

	// 简化版用户数据,仅用于测试关联
	_, err := db.Insert(ctx, "users", gdb.Map{
		"id":         userID,
		"username":   "testuser",
		"email":      "test@example.com",
		"created_at": time.Now(),
	})

	// 忽略主键冲突错误(用户可能已存在)
	if err != nil && !gdb.IsError(err, "UNIQUE constraint failed") {
		t.Logf("创建测试用户失败(可忽略): %v", err)
	}
}

// CreateTestProgress 创建测试学习进度数据
func CreateTestProgress(t *testing.T, db gdb.DB, userID int64, topic, chapter, status string) int64 {
	ctx := context.Background()

	data := gdb.Map{
		"user_id":          userID,
		"topic":            topic,
		"chapter":          chapter,
		"status":           status,
		"started_at":       time.Now().Add(-time.Hour),
		"last_accessed_at": time.Now(),
		"created_at":       time.Now(),
		"updated_at":       time.Now(),
	}

	if status == "completed" {
		completedAt := time.Now()
		data["completed_at"] = completedAt
	}

	result, err := db.Insert(ctx, "learning_progress", data)
	if err != nil {
		t.Fatalf("创建测试进度数据失败: %v", err)
	}

	id, _ := result.LastInsertId()
	return id
}

// CreateTestQuizSession 创建测试测验会话数据
func CreateTestQuizSession(t *testing.T, db gdb.DB, sessionID string, userID int64, topic, chapter string, submitted bool) int64 {
	ctx := context.Background()

	data := gdb.Map{
		"session_id":      sessionID,
		"user_id":         userID,
		"topic":           topic,
		"chapter":         chapter,
		"total_questions": 10,
		"correct_answers": 0,
		"score":           0,
		"created_at":      time.Now(),
		"updated_at":      time.Now(),
	}

	if submitted {
		submittedAt := time.Now()
		data["submitted_at"] = submittedAt
		data["score"] = 80
		data["correct_answers"] = 8
	}

	result, err := db.Insert(ctx, "quiz_sessions", data)
	if err != nil {
		t.Fatalf("创建测试测验会话失败: %v", err)
	}

	id, _ := result.LastInsertId()
	return id
}

// CreateTestQuizAttempt 创建测试答题记录数据
func CreateTestQuizAttempt(t *testing.T, db gdb.DB, sessionID string, questionID int, userAnswer string, isCorrect bool) {
	ctx := context.Background()

	_, err := db.Insert(ctx, "quiz_attempts", gdb.Map{
		"session_id":  sessionID,
		"question_id": questionID,
		"user_answer": userAnswer,
		"is_correct":  isCorrect,
		"created_at":  time.Now(),
	})

	if err != nil {
		t.Fatalf("创建测试答题记录失败: %v", err)
	}
}

// CleanTestData 清理指定用户的测试数据
func CleanTestData(t *testing.T, db gdb.DB, userID int64) {
	ctx := context.Background()

	// 清理学习进度
	_, err := db.Delete(ctx, "learning_progress", "user_id = ?", userID)
	if err != nil {
		t.Logf("清理learning_progress失败: %v", err)
	}

	// 清理测验会话
	_, err = db.Delete(ctx, "quiz_sessions", "user_id = ?", userID)
	if err != nil {
		t.Logf("清理quiz_sessions失败: %v", err)
	}

	// 注意: quiz_attempts 通过 session_id 关联,会话删除后自动清理
}
