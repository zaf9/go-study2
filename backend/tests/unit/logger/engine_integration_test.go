package logger_test

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"go-study2/internal/infrastructure/logger"
)

func TestSlogLumberjackIntegration(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "engine_int_")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)
	defer logger.Reset()

	logFileName := "integration.log"
	cfg := &logger.LoggerConfig{
		Level: "all",
		Instances: map[string]logger.InstanceConfig{
			"app": {
				Path:        tmpDir,
				File:        logFileName,
				Format:      "json", // Test JSON redirection
				Level:       "all",
				StdoutPrint: false,
			},
		},
	}

	logger.Reset()
	if err := logger.Initialize(cfg); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	testMsg := "Test redirection message"
	traceID := "test-trace-123"
	ctx := context.WithValue(context.Background(), logger.TraceIDKey{}, traceID)

	// Use LogInfo which calls GetInstance("app").Info()
	logger.LogInfo(ctx, testMsg)

	// Lumberjack might need a tiny bit of time or we just close the logger to flush
	logger.Reset()

	logFilePath := filepath.Join(tmpDir, logFileName)
	content, err := ioutil.ReadFile(logFilePath)
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}

	if len(content) == 0 {
		t.Fatalf("log file is empty")
	}

	// Verify JSON format
	var logMap map[string]interface{}
	if err := json.Unmarshal(content, &logMap); err != nil {
		t.Errorf("log content is not valid JSON: %v, content: %s", err, string(content))
	}

	// Verify fields
	if !strings.Contains(logMap["msg"].(string), testMsg) {
		t.Errorf("expected msg to contain %s, got %v", testMsg, logMap["msg"])
	}
	if !strings.Contains(logMap["msg"].(string), traceID) {
		t.Errorf("expected msg to contain trace ID %s, got %v", traceID, logMap["msg"])
	}
	if logMap["level"] != "INFO" {
		t.Errorf("expected level INFO, got %v", logMap["level"])
	}
}

func TestSlogTextFormatIntegration(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "engine_text_")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)
	defer logger.Reset()

	logFileName := "text.log"
	cfg := &logger.LoggerConfig{
		Level: "all",
		Instances: map[string]logger.InstanceConfig{
			"app": {
				Path:        tmpDir,
				File:        logFileName,
				Format:      "text",
				Level:       "all",
				StdoutPrint: false,
			},
		},
	}

	logger.Reset()
	if err := logger.Initialize(cfg); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	testMsg := "Test text message"
	logger.LogInfo(context.Background(), testMsg)

	logger.Reset()

	logFilePath := filepath.Join(tmpDir, logFileName)
	content, err := ioutil.ReadFile(logFilePath)
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}

	if !strings.Contains(string(content), "level=INFO") {
		t.Errorf("expected log to contain level=INFO, got %s", string(content))
	}
	if !strings.Contains(string(content), testMsg) {
		t.Errorf("expected log to contain %s, got %s", testMsg, string(content))
	}
}
