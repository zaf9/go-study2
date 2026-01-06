package repository

import (
	"context"
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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQuizRepository_Flow(t *testing.T) {
	ctx := gctx.New()
	db := newQuizTestDB(t)
	ensureUserTable(t, db)
	runFeatureMigrations(t, db)
	createUser(t, db, 1)

	// 注意：题目现在从 YAML 文件加载，不再需要 seedQuestions
	// seedQuestions(t, db)

	repo := infrarepo.NewQuizRepository(db)

	// 测试会话管理功能（题目从 YAML 仓库获取）
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
			QuestionID:  1, // 使用实际的题目 ID
			UserAnswers: `["A"]`,
			IsCorrect:   true,
		},
		{
			SessionID:   sid,
			UserID:      1,
			Topic:       "variables",
			Chapter:     "storage",
			QuestionID:  2, // 使用实际的题目 ID
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

// T056: 测试CreateSession()方法
func TestCreateSession(t *testing.T) {
	db := newQuizTestDB(t)
	ensureUserTable(t, db)
	runFeatureMigrations(t, db)
	repo := infrarepo.NewQuizRepository(db)
	ctx := context.Background()

	t.Run("成功创建会话", func(t *testing.T) {
		createUser(t, db, 1)
		session := &quiz.QuizSession{
			UserID:         1,
			Topic:          "constants",
			Chapter:        "boolean",
			TotalQuestions: 5,
		}

		sessionID, err := repo.CreateSession(ctx, session)

		require.NoError(t, err)
		assert.NotEmpty(t, sessionID)
		assert.Equal(t, sessionID, session.SessionID)

		// 验证数据库中存在该会话
		retrieved, err := repo.GetSession(ctx, sessionID)
		require.NoError(t, err)
		assert.NotNil(t, retrieved)
		assert.Equal(t, int64(1), retrieved.UserID)
		assert.Equal(t, "constants", retrieved.Topic)
		assert.Equal(t, "boolean", retrieved.Chapter)
		assert.Equal(t, 5, retrieved.TotalQuestions)
		assert.False(t, retrieved.StartedAt.IsZero())
	})

	t.Run("会话ID自动生成", func(t *testing.T) {
		createUser(t, db, 2)
		session := &quiz.QuizSession{
			UserID:         2,
			Topic:          "variables",
			Chapter:        "storage",
			TotalQuestions: 3,
		}

		sessionID, err := repo.CreateSession(ctx, session)

		require.NoError(t, err)
		assert.NotEmpty(t, sessionID)
		// UUID格式验证（36个字符，包含4个横线）
		assert.Len(t, sessionID, 36)
	})

	t.Run("传入nil会话返回错误", func(t *testing.T) {
		sessionID, err := repo.CreateSession(ctx, nil)

		assert.Error(t, err)
		assert.Empty(t, sessionID)
		assert.Contains(t, err.Error(), "nil")
	})
}

// T057: 测试GetActiveSession()查询24小时内的会话
func TestGetActiveSession(t *testing.T) {
	db := newQuizTestDB(t)
	ensureUserTable(t, db)
	runFeatureMigrations(t, db)
	repo := infrarepo.NewQuizRepository(db)
	ctx := context.Background()

	userID := int64(10)
	topic := "constants"
	chapter := "boolean"

	t.Run("获取24小时内的会话", func(t *testing.T) {
		createUser(t, db, userID)
		// 创建一个最近的会话
		recentSession := &quiz.QuizSession{
			UserID:         userID,
			Topic:          topic,
			Chapter:        chapter,
			TotalQuestions: 5,
			StartedAt:      time.Now().Add(-1 * time.Hour),
		}
		sessionID, err := repo.CreateSession(ctx, recentSession)
		require.NoError(t, err)

		// 获取活跃会话
		activeSession, err := repo.GetActiveSession(ctx, userID, topic, chapter)

		require.NoError(t, err)
		assert.NotNil(t, activeSession)
		assert.Equal(t, sessionID, activeSession.SessionID)
	})

	t.Run("不返回24小时之前的会话", func(t *testing.T) {
		createUser(t, db, userID+1)
		// 创建一个旧会话
		oldTime := time.Now().Add(-25 * time.Hour)
		oldSession := &quiz.QuizSession{
			SessionID:      "old-session-id-unique",
			UserID:         userID + 1,
			Topic:          topic,
			Chapter:        chapter,
			TotalQuestions: 3,
			StartedAt:      oldTime,
			CreatedAt:      oldTime,
		}
		_, err := db.Model("quiz_sessions").Data(oldSession).FieldsEx("id").Insert()
		require.NoError(t, err)

		// 尝试获取活跃会话
		activeSession, err := repo.GetActiveSession(ctx, userID+1, topic, chapter)

		require.NoError(t, err)
		assert.Nil(t, activeSession) // 不应该找到旧会话
	})

	t.Run("不返回已完成的会话", func(t *testing.T) {
		createUser(t, db, userID+2)
		// 创建并完成一个会话
		completedSession := &quiz.QuizSession{
			UserID:         userID + 2,
			Topic:          topic,
			Chapter:        chapter,
			TotalQuestions: 5,
		}
		sessionID, err := repo.CreateSession(ctx, completedSession)
		require.NoError(t, err)

		// 标记为已提交
		err = repo.MarkAsSubmitted(ctx, sessionID)
		require.NoError(t, err)

		// 尝试获取活跃会话
		activeSession, err := repo.GetActiveSession(ctx, userID+2, topic, chapter)

		require.NoError(t, err)
		assert.Nil(t, activeSession) // 已完成的会话不应该返回
	})
}

// T058: 测试MarkAsSubmitted()设置submitted_at字段
func TestMarkAsSubmitted(t *testing.T) {
	db := newQuizTestDB(t)
	ensureUserTable(t, db)
	runFeatureMigrations(t, db)
	repo := infrarepo.NewQuizRepository(db)
	ctx := context.Background()

	t.Run("成功标记会话为已提交", func(t *testing.T) {
		createUser(t, db, 100)
		// 创建会话
		session := &quiz.QuizSession{
			UserID:         100,
			Topic:          "constants",
			Chapter:        "boolean",
			TotalQuestions: 5,
		}
		sessionID, err := repo.CreateSession(ctx, session)
		require.NoError(t, err)

		// 验证submitted_at初始为空
		before, err := repo.GetSession(ctx, sessionID)
		require.NoError(t, err)
		assert.Nil(t, before.SubmittedAt)

		// 标记为已提交
		err = repo.MarkAsSubmitted(ctx, sessionID)
		require.NoError(t, err)

		// 验证submitted_at已设置
		after, err := repo.GetSession(ctx, sessionID)
		require.NoError(t, err)
		assert.NotNil(t, after.SubmittedAt)
		assert.False(t, after.SubmittedAt.IsZero())
	})

	t.Run("标记不存在的会话不报错", func(t *testing.T) {
		err := repo.MarkAsSubmitted(ctx, "non-existent-session-id")
		assert.NoError(t, err) // GoFrame的Update对0行不报错
	})
}
