package logger_test

import (
	"os"
	"testing"

	lg "go-study2/internal/infrastructure/logger"
)

func TestInvalidConfigShouldFailInitialize(t *testing.T) {
	// create a temporary config file with missing instances
	f, err := os.CreateTemp("", "logger_invalid_*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(f.Name())
	content := `logger:
  level: "all"
  # no instances configured
`
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}
	_ = f.Close()

	cfg, err := lg.LoadConfig(f.Name())
	if err != nil {
		t.Fatalf("LoadConfig should parse file even if instances missing: %v", err)
	}

	// Initialize should fail due to validation (no instances)
	lg.Reset()
	if err := lg.Initialize(cfg); err == nil {
		t.Fatalf("Initialize should fail for invalid config without instances")
	}
}
