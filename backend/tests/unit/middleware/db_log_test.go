package middleware_test

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"go-study2/internal/app/http_server/middleware"
	"go-study2/internal/infrastructure/logger"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// hookInputMock implements gdb.HookInput for testing
type hookInputMock struct {
	sql      string
	args     []interface{}
	nextFunc func(ctx context.Context) (gdb.Result, error)
}

func (m *hookInputMock) Sql() string {
	return m.sql
}

func (m *hookInputMock) Args() []interface{} {
	return m.args
}

func (m *hookInputMock) Next(ctx context.Context) (gdb.Result, error) {
	return m.nextFunc(ctx)
}

func TestDBLogHandler(t *testing.T) {
	// Setup test logger using a manual temp dir so we can control cleanup order on Windows
	dir, err := os.MkdirTemp("", "TestDBLogHandler")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	cfg := &logger.LoggerConfig{
		Level:  "info",
		Stdout: false,
		Instances: map[string]logger.InstanceConfig{
			"biz": {
				Path:   dir,
				File:   "biz.log",
				Level:  "info",
				Format: "text",
			},
			"slow": {
				Path:   dir,
				File:   "slow.log",
				Level:  "info",
				Format: "text",
			},
			"error": {
				Path:   dir,
				File:   "error.log",
				Level:  "info",
				Format: "text",
			},
			"app": {
				Path:   dir,
				File:   "app.log",
				Level:  "info",
				Format: "text",
			},
		},
	}

	err = logger.Initialize(cfg)
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}
	// Ensure logger resources are released and then remove the temp dir (Windows-safe)
	defer func() {
		logger.Reset()
		// give OS a moment to release file handles
		time.Sleep(50 * time.Millisecond)
		_ = os.RemoveAll(dir)
	}()

	ctx := context.WithValue(context.Background(), logger.TraceIDKey{}, "test-trace-id")

	// Test successful DB operation
	t.Run("SuccessfulDBOperation", func(t *testing.T) {
		input := &hookInputMock{
			sql:  "SELECT * FROM users WHERE id = ?",
			args: g.Array{1},
			nextFunc: func(ctx context.Context) (gdb.Result, error) {
				// Simulate successful query
				time.Sleep(10 * time.Millisecond)
				return gdb.Result{}, nil
			},
		}

		result, err := middleware.DBLogHandler(ctx, input)
		if err != nil {
			t.Errorf("DBLogHandler should not return error for successful operation, got %v", err)
		}
		if result == nil {
			t.Error("DBLogHandler should return result")
		}
	})

	// Test DB operation with error
	t.Run("DBOperationWithError", func(t *testing.T) {
		testErr := errors.New("database connection failed")
		input := &hookInputMock{
			sql:  "SELECT * FROM users",
			args: g.Array{},
			nextFunc: func(ctx context.Context) (gdb.Result, error) {
				// Simulate failed query
				time.Sleep(5 * time.Millisecond)
				return nil, testErr
			},
		}

		result, err := middleware.DBLogHandler(ctx, input)
		if err != testErr {
			t.Errorf("DBLogHandler should return the original error, expected %v, got %v", testErr, err)
		}
		if result != nil {
			t.Error("DBLogHandler should return nil result on error")
		}
	})

	// Test slow query detection
	t.Run("SlowQueryDetection", func(t *testing.T) {
		input := &hookInputMock{
			sql:  "SELECT * FROM large_table",
			args: g.Array{},
			nextFunc: func(ctx context.Context) (gdb.Result, error) {
				// Simulate slow query (longer than 1 second threshold)
				time.Sleep(1100 * time.Millisecond)
				return gdb.Result{}, nil
			},
		}

		result, err := middleware.DBLogHandler(ctx, input)
		if err != nil {
			t.Errorf("DBLogHandler should not return error for slow operation, got %v", err)
		}
		if result == nil {
			t.Error("DBLogHandler should return result")
		}
	})

	// Test with result rows
	t.Run("DBOperationWithRows", func(t *testing.T) {
		// Create a mock result with rows
		mockResult := gdb.Result{}
		// Note: In real implementation, this would have actual rows
		// For testing purposes, we just ensure the handler processes it

		input := &hookInputMock{
			sql:  "SELECT id, name FROM users",
			args: g.Array{},
			nextFunc: func(ctx context.Context) (gdb.Result, error) {
				time.Sleep(5 * time.Millisecond)
				return mockResult, nil
			},
		}

		result, err := middleware.DBLogHandler(ctx, input)
		if err != nil {
			t.Errorf("DBLogHandler should not return error, got %v", err)
		}
		if result == nil {
			t.Error("DBLogHandler should return result")
		}
	})
}

func TestDBLogHandlerLogging(t *testing.T) {
	// Setup test logger using a manual temp dir so we can control cleanup order on Windows
	dir, err := os.MkdirTemp("", "TestDBLogHandlerLogging")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	cfg := &logger.LoggerConfig{
		Level:  "info",
		Stdout: false,
		Instances: map[string]logger.InstanceConfig{
			"biz": {
				Path:   dir,
				File:   "biz.log",
				Level:  "info",
				Format: "text",
			},
		},
	}

	err = logger.Initialize(cfg)
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}
	// Ensure logger resources are released and then remove the temp dir (Windows-safe)
	defer func() {
		logger.Reset()
		time.Sleep(50 * time.Millisecond)
		_ = os.RemoveAll(dir)
	}()

	ctx := context.WithValue(context.Background(), logger.TraceIDKey{}, "test-trace-id")

	// Test that logging functions are called (we can't easily verify the log content in unit tests)
	// but we can ensure the handler completes without panicking
	// Reuse hookInputMock to exercise the handler with a concrete input type.
	input := &hookInputMock{
		sql:  "SELECT COUNT(*) FROM users",
		args: g.Array{},
		nextFunc: func(ctx context.Context) (gdb.Result, error) {
			return gdb.Result{}, nil
		},
	}

	_, err = middleware.DBLogHandler(ctx, input)
	if err != nil {
		t.Errorf("DBLogHandler should complete without error, got %v", err)
	}
}
