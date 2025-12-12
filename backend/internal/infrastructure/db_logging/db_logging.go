package db_logging

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"

	"go-study2/internal/infrastructure/logger"

	"github.com/gogf/gf/v2/database/gdb"
)

// DBLogHandler 数据库操作日志处理器
type DBLogHandler struct {
	slowThreshold time.Duration
}

// NewDBLogHandler 创建新的数据库日志处理器
func NewDBLogHandler(slowThreshold time.Duration) *DBLogHandler {
	return &DBLogHandler{
		slowThreshold: slowThreshold,
	}
}

// RegisterDBLogging 注册数据库日志处理器到GoFrame数据库
func RegisterDBLogging(db gdb.DB, slowThreshold time.Duration) error {
	if db == nil {
		return fmt.Errorf("db is nil")
	}

	handler := NewDBLogHandler(slowThreshold)

	// Register hooks on the default model so that common operations go through
	// our logging handler. Use the Model().Hook API which accepts a HookHandler
	// with specific function signatures for select/insert/update/delete.
	db.Model().Hook(gdb.HookHandler{
		Select: func(ctx context.Context, in *gdb.HookSelectInput) (gdb.Result, error) {
			start := time.Now()
			res, err := in.Next(ctx)
			duration := time.Since(start)
			handler.LogDBOperation(ctx, "SELECT", in.Sql, in.Args, duration, err)
			return res, err
		},
		Insert: func(ctx context.Context, in *gdb.HookInsertInput) (sql.Result, error) {
			start := time.Now()
			res, err := in.Next(ctx)
			duration := time.Since(start)
			// Insert hook does not expose raw SQL string in the same way; try to
			// log meaningful info. Data is available in in.Data.
			handler.LogDBOperation(ctx, "INSERT", "<insert>", nil, duration, err)
			return res, err
		},
		Update: func(ctx context.Context, in *gdb.HookUpdateInput) (sql.Result, error) {
			start := time.Now()
			res, err := in.Next(ctx)
			duration := time.Since(start)
			handler.LogDBOperation(ctx, "UPDATE", "<update>", in.Args, duration, err)
			return res, err
		},
		Delete: func(ctx context.Context, in *gdb.HookDeleteInput) (sql.Result, error) {
			start := time.Now()
			res, err := in.Next(ctx)
			duration := time.Since(start)
			handler.LogDBOperation(ctx, "DELETE", "<delete>", in.Args, duration, err)
			return res, err
		},
	})

	return nil
}

// LogDBOperation 记录数据库操作
func (h *DBLogHandler) LogDBOperation(ctx context.Context, operation string, sql string, args []interface{}, duration time.Duration, err error) {
	params := map[string]interface{}{
		"operation": operation,
		"sql":       sql,
		"duration":  duration.String(),
		"args":      args,
	}

	if err != nil {
		logger.LogError(ctx, err, fmt.Sprintf("Database operation failed: %s", operation), params)
	} else {
		logger.LogBiz(ctx, operation, params, "completed", duration)

		// 检查是否为慢查询
		if duration >= h.slowThreshold {
			logger.LogSlow(ctx, fmt.Sprintf("Database slow query: %s", sql), duration, h.slowThreshold)
		}
	}
}

// Handle 实现 database/sql/driver.Conn 接口，用于拦截数据库操作
func (h *DBLogHandler) Handle(ctx context.Context, operation string, sql string, args []driver.NamedValue, result driver.Result, err error) {
	// Try to obtain a start time from the context if upstream hook provides it.
	var duration time.Duration
	if v := ctx.Value("dbStartTime"); v != nil {
		if t, ok := v.(time.Time); ok {
			duration = time.Since(t)
		}
	}
	if duration == 0 {
		// best-effort fallback: if result implements driver.Result with unknown metadata,
		// we don't have a reliable start time — keep duration as zero to avoid negative values.
		duration = 0
	}
	h.LogDBOperation(ctx, operation, sql, convertNamedValues(args), duration, err)
}

// HandleQuery 专门处理查询操作
func (h *DBLogHandler) HandleQuery(ctx context.Context, sql string, args []driver.NamedValue, rows *sql.Rows, err error) {
	var duration time.Duration
	if v := ctx.Value("dbStartTime"); v != nil {
		if t, ok := v.(time.Time); ok {
			duration = time.Since(t)
		}
	}
	h.LogDBOperation(ctx, "QUERY", sql, convertNamedValues(args), duration, err)
}

// HandleExec 专门处理执行操作
func (h *DBLogHandler) HandleExec(ctx context.Context, sql string, args []driver.NamedValue, result sql.Result, err error) {
	var duration time.Duration
	if v := ctx.Value("dbStartTime"); v != nil {
		if t, ok := v.(time.Time); ok {
			duration = time.Since(t)
		}
	}
	h.LogDBOperation(ctx, "EXEC", sql, convertNamedValues(args), duration, err)
}

// convertNamedValues 将 driver.NamedValue 转换为 interface{} 切片
func convertNamedValues(namedArgs []driver.NamedValue) []interface{} {
	args := make([]interface{}, len(namedArgs))
	for i, arg := range namedArgs {
		args[i] = arg.Value
	}
	return args
}
