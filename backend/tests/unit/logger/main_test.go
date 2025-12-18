package logger_test

import (
	"os"
	"testing"
	"time"

	"go-study2/internal/infrastructure/logger"
)

// TestMain provides test lifecycle management for the logger package tests.
// It ensures that all resources are properly cleaned up after tests complete,
// preventing hangs caused by lingering goroutines.
func TestMain(m *testing.M) {
	// Run all tests
	code := m.Run()

	// Final cleanup
	logger.Reset()

	// Give a brief moment for cleanup to complete
	time.Sleep(200 * time.Millisecond)

	// Force exit after a timeout to prevent hanging
	// This is a safety net in case some goroutine doesn't exit properly
	go func() {
		time.Sleep(5 * time.Second)
		os.Exit(code)
	}()

	os.Exit(code)
}
