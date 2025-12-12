package logger_test

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"go-study2/internal/infrastructure/logger"
)

func TestLogInfo(t *testing.T) {
	// Setup test logger
	logger.Reset()
	// Use manual temp dir so we can control removal after logger Reset to avoid
	// file-handle races on Windows.
	dir, err := os.MkdirTemp("", "TestLogBiz")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() {
		logger.Reset()
		time.Sleep(100 * time.Millisecond)
		_ = os.RemoveAll(dir)
	}()
	// Debug: list directory before initialization
	if entries, err := os.ReadDir(dir); err == nil {
		var names []string
		for _, e := range entries {
			names = append(names, e.Name())
		}
		t.Logf("dir before init: %v", names)
	}

	// Setup logger configuration with unique file names
	cfg := &logger.LoggerConfig{
		Level:  "info",
		Stdout: true, // Use stdout to avoid file locking
		Instances: map[string]logger.InstanceConfig{
			"app": {
				Path:   dir,
				File:   "app_test.log",
				Level:  "info",
				Format: "text",
			},
		},
	}

	err = logger.Initialize(cfg)
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}

	// Test LogInfo with trace ID
	ctx := context.WithValue(context.Background(), logger.TraceIDKey{}, "test-trace-id")
	logger.LogInfo(ctx, "Test info message")

	// Test LogInfo without trace ID
	logger.LogInfo(context.Background(), "Test info message without trace")

	// Test LogInfo with formatted message
	logger.LogInfo(ctx, "Test formatted message: %s %d", "hello", 123)
}

func TestLogError(t *testing.T) {
	// Setup test logger
	logger.Reset()
	dir := t.TempDir()
	cfg := &logger.LoggerConfig{
		Level: "info",
		// Use stdout in unit tests to avoid file handle pooling issues on Windows.
		Stdout: true,
		Instances: map[string]logger.InstanceConfig{
			"error": {
				Path:   dir,
				File:   "error_test.log",
				Level:  "info",
				Format: "text",
			},
			"app": {
				Path:   dir,
				File:   "app_test.log",
				Level:  "info",
				Format: "text",
			},
		},
	}

	err := logger.Initialize(cfg)
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}

	testErr := errors.New("test error")

	// Test LogError with trace ID and error
	ctx := context.WithValue(context.Background(), logger.TraceIDKey{}, "test-trace-id")
	logger.LogError(ctx, testErr, "Test error message")

	// Test LogError without error
	logger.LogError(ctx, nil, "Test error message without error")

	// Test LogError with formatted message
	logger.LogError(ctx, testErr, "Test formatted error: %s %d", "error", 456)
}

func TestLogSlow(t *testing.T) {
	// Setup test logger
	logger.Reset()
	dir := t.TempDir()
	cfg := &logger.LoggerConfig{
		Level: "info",
		// Use stdout to avoid file locking in unit tests on Windows.
		Stdout: true,
		Instances: map[string]logger.InstanceConfig{
			"slow": {
				Path:   dir,
				File:   "slow_test.log",
				Level:  "info",
				Format: "text",
			},
			"app": {
				Path:   dir,
				File:   "app_test.log",
				Level:  "info",
				Format: "text",
			},
		},
	}

	err := logger.Initialize(cfg)
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}

	ctx := context.WithValue(context.Background(), logger.TraceIDKey{}, "test-trace-id")

	// Test LogSlow above threshold
	logger.LogSlow(ctx, "test operation", 200*time.Millisecond, 100*time.Millisecond)

	// Test LogSlow below threshold (should not log)
	logger.LogSlow(ctx, "test operation", 50*time.Millisecond, 100*time.Millisecond)

	// Test LogSlow at threshold (should not log)
	logger.LogSlow(ctx, "test operation", 100*time.Millisecond, 100*time.Millisecond)
}

func TestLogBiz(t *testing.T) {
	// Setup test logger
	logger.Reset()
	// Use manual temp dir so we can control cleanup timing and avoid races on Windows.
	dir, err := os.MkdirTemp("", "TestLogBiz")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() {
		logger.Reset()
		time.Sleep(100 * time.Millisecond)
		_ = os.RemoveAll(dir)
	}()

	cfg := &logger.LoggerConfig{
		Level: "info",
		// Use stdout to avoid file handle pooling issues on Windows unit tests.
		Stdout: true,
		Instances: map[string]logger.InstanceConfig{
			"biz": {
				Path:   dir,
				File:   "biz_test.log",
				Level:  "info",
				Format: "text",
			},
			"app": {
				Path:   dir,
				File:   "app_test.log",
				Level:  "info",
				Format: "text",
			},
		},
	}

	err = logger.Initialize(cfg)
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}

	ctx := context.WithValue(context.Background(), logger.TraceIDKey{}, "test-trace-id")

	// Test LogBiz with params and result
	params := map[string]interface{}{
		"user_id": 123,
		"action":  "login",
	}
	result := map[string]interface{}{
		"success": true,
		"token":   "abc123",
	}

	logger.LogBiz(ctx, "user_login", params, result, 50*time.Millisecond)

	// Test LogBiz with nil params and result
	logger.LogBiz(ctx, "simple_operation", nil, nil, 10*time.Millisecond)

	// Test LogBiz with empty params
	logger.LogBiz(ctx, "empty_params", map[string]interface{}{}, "success", 20*time.Millisecond)

	// Debug: list directory after logging
	if entries, err := os.ReadDir(dir); err == nil {
		var names []string
		for _, e := range entries {
			names = append(names, e.Name())
		}
		t.Logf("dir after logging: %v", names)
	}
}

func TestLogWithFields(t *testing.T) {
	// Setup test logger
	logger.Reset()
	dir := t.TempDir()
	cfg := &logger.LoggerConfig{
		Level:  "info",
		Stdout: true,
		Instances: map[string]logger.InstanceConfig{
			"biz": {
				Path:   dir,
				File:   "biz_test.log",
				Level:  "info",
				Format: "text",
			},
		},
	}

	err := logger.Initialize(cfg)
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}

	ctx := context.WithValue(context.Background(), logger.TraceIDKey{}, "test-trace-id")

	// Test LogWithFields with various levels
	fields := map[string]interface{}{
		"user_id": 123,
		"action":  "test",
		"status":  "success",
	}

	logger.LogWithFields(ctx, "INFO", "Test message with fields", fields)
	logger.LogWithFields(ctx, "ERROR", "Error message with fields", fields)
	logger.LogWithFields(ctx, "DEBUG", "Debug message with fields", fields)
	logger.LogWithFields(ctx, "WARNING", "Warning message with fields", fields)

	// Test LogWithFields with empty fields
	logger.LogWithFields(ctx, "INFO", "Message with no fields", map[string]interface{}{})
}
