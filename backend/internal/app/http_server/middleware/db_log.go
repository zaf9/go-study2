package middleware

import (
	"context"
	"time"

	"go-study2/internal/infrastructure/db_logging"
	"go-study2/internal/infrastructure/logger"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// DBLogHandler 数据库日志 Handler - 符合 GoFrame Hook 规范
// dbHookInput is a minimal interface representing the subset of GoFrame's
// hook input required by our DBLogHandler. We use a local interface to avoid
// depending on a concrete type that may vary between GoFrame versions.
type dbHookInput interface {
	Sql() string
	Args() []interface{}
	Next(ctx context.Context) (gdb.Result, error)
}

// DBLogHandler 数据库日志 Handler - 符合 GoFrame Hook 规范
// Accepts any value implementing the minimal dbHookInput interface so tests
// and different GoFrame hook input types are supported.
func DBLogHandler(ctx context.Context, in dbHookInput) (result gdb.Result, err error) {
	startTime := time.Now()

	// 执行 SQL
	result, err = in.Next(ctx)

	// 计算耗时
	duration := time.Since(startTime)

	// 构建日志数据
	logData := g.Map{
		"sql":      in.Sql(),
		"args":     in.Args(),
		"duration": duration.Milliseconds(),
		"rows":     0,
	}

	if result != nil {
		logData["rows"] = result.Len()
	}

	// 使用 logger 记录数据库操作
	if err != nil {
		logData["error"] = err.Error()
		logger.LogError(ctx, err, "Database operation failed")
	} else {
		logger.LogBiz(ctx, "DB_OPERATION", logData, "completed", duration)

		// 检查是否为慢查询
		slowThreshold := 1 * time.Second // 可以从配置中读取
		if duration >= slowThreshold {
			logger.LogSlow(ctx, "Database slow query", duration, slowThreshold)
		}
	}

	return result, err
}

// RegisterDBLogging 注册数据库日志处理器到GoFrame数据库
func RegisterDBLogging(db gdb.DB, slowThreshold time.Duration) error {
	return db_logging.RegisterDBLogging(db, slowThreshold)
}
