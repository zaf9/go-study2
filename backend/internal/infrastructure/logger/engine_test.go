package logger

import (
	"log/slog"
	"testing"
	"time"

	"github.com/gogf/gf/v2/os/glog"
)

func TestParseSizeToMB(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"100M", 100},
		{"100MB", 100},
		{"1G", 1024},
		{"1GB", 1024},
		{"1024K", 1},
		{"1024KB", 1},
		{"2097152", 2}, // bytes
		{"", 100},      // default
		{"invalid", 100},
	}

	for _, tt := range tests {
		result := parseSizeToMB(tt.input)
		if result != tt.expected {
			t.Errorf("parseSizeToMB(%s) = %d; want %d", tt.input, result, tt.expected)
		}
	}
}

func TestParseAgeToDays(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"30d", 30},
		{"1d", 1},
		{"72h", 3},
		{"", 30}, // default
		{"invalid", 30},
	}

	for _, tt := range tests {
		result := parseAgeToDays(tt.input)
		if result != tt.expected {
			t.Errorf("parseAgeToDays(%s) = %d; want %d", tt.input, result, tt.expected)
		}
	}
}

func TestMapGlogLevelToSlog(t *testing.T) {
	tests := []struct {
		input    int
		expected slog.Level
	}{
		{glog.LEVEL_DEBU, slog.LevelDebug},
		{glog.LEVEL_INFO, slog.LevelInfo},
		{glog.LEVEL_NOTI, slog.LevelInfo},
		{glog.LEVEL_WARN, slog.LevelWarn},
		{glog.LEVEL_ERRO, slog.LevelError},
		{glog.LEVEL_CRIT, LevelCritical},
		{999, slog.LevelInfo}, // default
	}

	for _, tt := range tests {
		result := mapGlogLevelToSlog(tt.input)
		if result != tt.expected {
			t.Errorf("mapGlogLevelToSlog(%d) = %v; want %v", tt.input, result, tt.expected)
		}
	}
}

func TestCleanOldLogsMock(t *testing.T) {
	// CleanOldLogs is now a no-op or restored to original logic.
	// We already have a test for the restored logic in tests/unit/logger/cleanup_test.go.
	// This just ensures it doesn't crash.
	err := CleanOldLogs(".", time.Hour)
	if err != nil {
		t.Errorf("CleanOldLogs failed: %v", err)
	}
}
