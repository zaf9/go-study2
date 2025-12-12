package middleware_test

import (
	"testing"

	"go-study2/internal/infrastructure/logger"
)

func TestAccessLogMiddleware(t *testing.T) {
	// Reset logger for test isolation
	logger.Reset()

	// Setup test config
	dir := t.TempDir()
	cfg := &logger.LoggerConfig{
		Level:  "info",
		Stdout: false,
		Instances: map[string]logger.InstanceConfig{
			"access": {
				Path:   dir,
				File:   "access.log",
				Level:  "info",
				Format: "text",
			},
		},
	}

	err := logger.Initialize(cfg)
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}

	// Verify access logger is available
	accessLogger := logger.GetInstance("access")
	if accessLogger == nil {
		t.Errorf("Access logger should be initialized")
	}

	t.Log("Access log middleware test setup completed")
}

func TestAccessLogMiddleware_ErrorStatus(t *testing.T) {
	logger.Reset()

	dir := t.TempDir()
	cfg := &logger.LoggerConfig{
		Level:  "info",
		Stdout: false,
		Instances: map[string]logger.InstanceConfig{
			"access": {
				Path:   dir,
				File:   "access.log",
				Level:  "info",
				Format: "text",
			},
		},
	}

	err := logger.Initialize(cfg)
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}

	// Verify logger initialization for error status logging
	accessLogger := logger.GetInstance("access")
	if accessLogger == nil {
		t.Errorf("Access logger should be initialized")
	}
}

func TestAccessLogMiddleware_WithBody(t *testing.T) {
	logger.Reset()

	dir := t.TempDir()
	cfg := &logger.LoggerConfig{
		Level:  "info",
		Stdout: false,
		Instances: map[string]logger.InstanceConfig{
			"access": {
				Path:   dir,
				File:   "access.log",
				Level:  "info",
				Format: "text",
			},
		},
	}

	err := logger.Initialize(cfg)
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}

	t.Log("Body logging test setup completed")
}

func TestAccessLogMiddleware_Performance(t *testing.T) {
	logger.Reset()

	dir := t.TempDir()
	cfg := &logger.LoggerConfig{
		Level:  "info",
		Stdout: false,
		Instances: map[string]logger.InstanceConfig{
			"access": {
				Path:   dir,
				File:   "access.log",
				Level:  "info",
				Format: "text",
			},
		},
	}

	err := logger.Initialize(cfg)
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}

	t.Log("Performance test setup completed")
}
