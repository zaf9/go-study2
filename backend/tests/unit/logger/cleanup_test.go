package logger_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	lg "go-study2/internal/infrastructure/logger"
)

func TestCleanOldLogs(t *testing.T) {
	dir, err := os.MkdirTemp("", "logcleanup_")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	// recent file (should remain)
	recent := filepath.Join(dir, "recent.log")
	if err := os.WriteFile(recent, []byte("ok"), 0o644); err != nil {
		t.Fatalf("write recent: %v", err)
	}

	// old file (should be removed)
	old := filepath.Join(dir, "old.log")
	if err := os.WriteFile(old, []byte("old"), 0o644); err != nil {
		t.Fatalf("write old: %v", err)
	}
	// set modtime to past
	past := time.Now().Add(-48 * time.Hour)
	if err := os.Chtimes(old, past, past); err != nil {
		t.Fatalf("chtimes: %v", err)
	}

	// run cleanup with 24h expire
	if err := lg.CleanOldLogs(dir, 24*time.Hour); err != nil {
		t.Fatalf("CleanOldLogs failed: %v", err)
	}

	if _, err := os.Stat(old); !os.IsNotExist(err) {
		t.Fatalf("expected old file removed")
	}
	if _, err := os.Stat(recent); err != nil {
		t.Fatalf("expected recent file to remain: %v", err)
	}
}
