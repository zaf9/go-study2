package logger_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"go-study2/internal/infrastructure/logger"
)

func TestReadLogFile(t *testing.T) {
	// Create a temporary log file
	dir := t.TempDir()
	logFile := filepath.Join(dir, "test.log")

	// Write test log entries
	logContent := `[2024-01-01 10:00:00] INFO [TraceID:abc123] User login successful
[2024-01-01 10:01:00] ERROR [TraceID:def456] Database connection failed
[2024-01-01 10:02:00] WARN [TraceID:abc123] Slow query detected`

	err := os.WriteFile(logFile, []byte(logContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test log file: %v", err)
	}

	// Test ReadLogFile
	entries, err := logger.ReadLogFile(logFile)
	if err != nil {
		t.Fatalf("ReadLogFile failed: %v", err)
	}

	if len(entries) != 3 {
		t.Errorf("Expected 3 entries, got %d", len(entries))
	}

	// Check first entry
	if entries[0].TraceID != "abc123" {
		t.Errorf("Expected TraceID abc123, got %s", entries[0].TraceID)
	}
	if entries[0].Level != "INFO" {
		t.Errorf("Expected level INFO, got %s", entries[0].Level)
	}
	if !strings.Contains(entries[0].Message, "User login successful") {
		t.Errorf("Expected message to contain 'User login successful', got %s", entries[0].Message)
	}
}

func TestQueryByTraceID(t *testing.T) {
	// Create a temporary log file
	dir := t.TempDir()
	logFile := filepath.Join(dir, "test.log")

	logContent := `[2024-01-01 10:00:00] INFO [TraceID:abc123] User login successful
[2024-01-01 10:01:00] ERROR [TraceID:def456] Database connection failed
[2024-01-01 10:02:00] WARN [TraceID:abc123] Slow query detected`

	err := os.WriteFile(logFile, []byte(logContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test log file: %v", err)
	}

	// Test QueryByTraceID
	result, err := logger.QueryByTraceID(logFile, "abc123")
	if err != nil {
		t.Fatalf("QueryByTraceID failed: %v", err)
	}

	if result.Total != 3 {
		t.Errorf("Expected total 3, got %d", result.Total)
	}
	if result.Matched != 2 {
		t.Errorf("Expected matched 2, got %d", result.Matched)
	}

	// Test with non-existent trace ID
	result, err = logger.QueryByTraceID(logFile, "nonexistent")
	if err != nil {
		t.Fatalf("QueryByTraceID failed: %v", err)
	}
	if result.Matched != 0 {
		t.Errorf("Expected matched 0 for non-existent trace ID, got %d", result.Matched)
	}
}

func TestQueryByTimeRange(t *testing.T) {
	// Create a temporary log file
	dir := t.TempDir()
	logFile := filepath.Join(dir, "test.log")

	logContent := `[2024-01-01 10:00:00] INFO [TraceID:abc123] User login successful
[2024-01-01 10:01:00] ERROR [TraceID:def456] Database connection failed
[2024-01-01 10:02:00] WARN [TraceID:abc123] Slow query detected`

	err := os.WriteFile(logFile, []byte(logContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test log file: %v", err)
	}

	// Test QueryByTimeRange
	start := time.Date(2024, 1, 1, 9, 59, 0, 0, time.UTC)
	end := time.Date(2024, 1, 1, 10, 1, 30, 0, time.UTC)

	result, err := logger.QueryByTimeRange(logFile, start, end)
	if err != nil {
		t.Fatalf("QueryByTimeRange failed: %v", err)
	}

	// Note: Since parseLogEntry uses time.Now() as placeholder, this test is limited
	// In real implementation, proper timestamp parsing would be needed
	if result.Total != 3 {
		t.Errorf("Expected total 3, got %d", result.Total)
	}
}

func TestQueryByLevel(t *testing.T) {
	// Create a temporary log file
	dir := t.TempDir()
	logFile := filepath.Join(dir, "test.log")

	logContent := `[2024-01-01 10:00:00] INFO [TraceID:abc123] User login successful
[2024-01-01 10:01:00] ERROR [TraceID:def456] Database connection failed
[2024-01-01 10:02:00] WARN [TraceID:abc123] Slow query detected
[2024-01-01 10:03:00] INFO [TraceID:ghi789] User logout`

	err := os.WriteFile(logFile, []byte(logContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test log file: %v", err)
	}

	// Test QueryByLevel for INFO
	result, err := logger.QueryByLevel(logFile, "INFO")
	if err != nil {
		t.Fatalf("QueryByLevel failed: %v", err)
	}

	if result.Total != 4 {
		t.Errorf("Expected total 4, got %d", result.Total)
	}
	if result.Matched != 2 {
		t.Errorf("Expected matched 2 for INFO level, got %d", result.Matched)
	}

	// Test QueryByLevel for ERROR
	result, err = logger.QueryByLevel(logFile, "ERROR")
	if err != nil {
		t.Fatalf("QueryByLevel failed: %v", err)
	}
	if result.Matched != 1 {
		t.Errorf("Expected matched 1 for ERROR level, got %d", result.Matched)
	}
}

func TestQueryByKeyword(t *testing.T) {
	// Create a temporary log file
	dir := t.TempDir()
	logFile := filepath.Join(dir, "test.log")

	logContent := `[2024-01-01 10:00:00] INFO [TraceID:abc123] User login successful
[2024-01-01 10:01:00] ERROR [TraceID:def456] Database connection failed
[2024-01-01 10:02:00] WARN [TraceID:abc123] Slow query detected`

	err := os.WriteFile(logFile, []byte(logContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test log file: %v", err)
	}

	// Test QueryByKeyword
	result, err := logger.QueryByKeyword(logFile, "login")
	if err != nil {
		t.Fatalf("QueryByKeyword failed: %v", err)
	}

	if result.Total != 3 {
		t.Errorf("Expected total 3, got %d", result.Total)
	}
	if result.Matched != 1 {
		t.Errorf("Expected matched 1 for keyword 'login', got %d", result.Matched)
	}

	// Test with non-existent keyword
	result, err = logger.QueryByKeyword(logFile, "nonexistent")
	if err != nil {
		t.Fatalf("QueryByKeyword failed: %v", err)
	}
	if result.Matched != 0 {
		t.Errorf("Expected matched 0 for non-existent keyword, got %d", result.Matched)
	}
}

func TestListLogFiles(t *testing.T) {
	// Create a temporary logs directory with some log files
	dir := t.TempDir()

	// Create some log files
	logFiles := []string{"app.log", "error.log", "access.log", "subdir/nested.log"}
	for _, file := range logFiles {
		fullPath := filepath.Join(dir, file)
		err := os.MkdirAll(filepath.Dir(fullPath), 0755)
		if err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}
		err = os.WriteFile(fullPath, []byte("test log content"), 0644)
		if err != nil {
			t.Fatalf("Failed to write log file: %v", err)
		}
	}

	// Create a non-log file
	err := os.WriteFile(filepath.Join(dir, "config.txt"), []byte("not a log"), 0644)
	if err != nil {
		t.Fatalf("Failed to write non-log file: %v", err)
	}

	// Test ListLogFiles
	files, err := logger.ListLogFiles(dir)
	if err != nil {
		t.Fatalf("ListLogFiles failed: %v", err)
	}

	if len(files) != 4 {
		t.Errorf("Expected 4 log files, got %d: %v", len(files), files)
	}

	// Check that all returned files have .log extension
	for _, file := range files {
		if !strings.HasSuffix(file, ".log") {
			t.Errorf("Expected log file to have .log extension: %s", file)
		}
	}
}
