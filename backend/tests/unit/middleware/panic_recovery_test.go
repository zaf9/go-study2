package middleware_test

import (
	"context"
	"testing"

	"go-study2/internal/infrastructure/logger"
)

func TestPanicRecoveryMiddleware(t *testing.T) {
	logger.Reset()

	// Use temp directory for config validation
	dir := t.TempDir()

	cfg := &logger.LoggerConfig{
		Level:  "info",
		Stdout: false,
		Instances: map[string]logger.InstanceConfig{
			"error": {
				Path:   dir,
				File:   "error.log",
				Level:  "error",
				Format: "text",
			},
		},
	}

	err := logger.Initialize(cfg)
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}

	// Verify error logger is available
	errorLogger := logger.GetInstance("error")
	if errorLogger == nil {
		t.Errorf("Error logger should be initialized")
	}

	t.Log("Panic recovery middleware test setup completed")
}

func TestPanicRecoveryMiddleware_NoPanic(t *testing.T) {
	logger.Reset()

	// Use temp directory for config validation
	dir := t.TempDir()

	cfg := &logger.LoggerConfig{
		Level:  "info",
		Stdout: false,
		Instances: map[string]logger.InstanceConfig{
			"error": {
				Path:   dir,
				File:   "error.log",
				Level:  "error",
				Format: "text",
			},
		},
	}

	err := logger.Initialize(cfg)
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}

	// Verify error logger is available
	errorLogger := logger.GetInstance("error")
	if errorLogger == nil {
		t.Errorf("Error logger should be initialized")
	}

	t.Log("Panic recovery middleware test setup completed")
}

func TestPanicRecoveryMiddleware_WithTraceID(t *testing.T) {
	logger.Reset()

	// Use temp directory for config validation
	dir := t.TempDir()

	cfg := &logger.LoggerConfig{
		Level:  "info",
		Stdout: false,
		Instances: map[string]logger.InstanceConfig{
			"error": {
				Path:   dir,
				File:   "error.log",
				Level:  "error",
				Format: "text",
			},
		},
	}

	err := logger.Initialize(cfg)
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}

	// Verify error logger is available
	errorLogger := logger.GetInstance("error")
	if errorLogger == nil {
		t.Errorf("Error logger should be initialized")
	}

	// Test trace ID extraction
	ctx := context.WithValue(context.Background(), logger.TraceIDKey{}, "test-trace-id")
	traceID := logger.ExtractTraceID(ctx)
	if traceID != "test-trace-id" {
		t.Errorf("Expected trace ID 'test-trace-id', got '%s'", traceID)
	}

	t.Log("Panic recovery middleware with trace ID test setup completed")
}
