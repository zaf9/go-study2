package logger_test

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	logger "go-study2/internal/infrastructure/logger"
)

func TestInitializeLogger(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "glog_test_")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	cfg := &logger.LoggerConfig{
		Level: "all",
		Instances: map[string]logger.InstanceConfig{
			"app": {
				Path:       tmpDir,
				File:       "{Y-m-d}.log",
				Format:     "json",
				Level:      "all",
				RotateSize: "1M", // MB
			},
		},
	}

	// reset previous state
	logger.Reset()

	if err := logger.Initialize(cfg); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	inst := logger.GetInstance("app")
	if inst == nil {
		t.Fatalf("expected app instance to be available")
	}

	// ensure log file directory exists
	if _, err := os.Stat(filepath.Join(tmpDir)); err != nil {
		t.Fatalf("log path not created: %v", err)
	}
}

func TestMultiInstanceInitialization(t *testing.T) {
	tmpDir1, _ := os.MkdirTemp("", "glog_test1_")
	tmpDir2, _ := os.MkdirTemp("", "glog_test2_")
	defer os.RemoveAll(tmpDir1)
	defer os.RemoveAll(tmpDir2)

	cfg := &logger.LoggerConfig{
		Level: "all",
		Instances: map[string]logger.InstanceConfig{
			"access": {
				Path:   tmpDir1,
				File:   "access-{Y-m-d}.log",
				Format: "json",
				Level:  "info",
			},
			"error": {
				Path:   tmpDir2,
				File:   "error-{Y-m-d}.log",
				Format: "text",
				Level:  "error",
			},
		},
	}

	logger.Reset()
	if err := logger.Initialize(cfg); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	if logger.GetInstance("access") == nil || logger.GetInstance("error") == nil {
		t.Fatalf("expected both instances to be initialized")
	}
}

func TestInitializeBehavior(t *testing.T) {
	logger.Reset() // Reset for test isolation

	dir, err := ioutil.TempDir("", "loggertest")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(dir)

	// Create test config
	cfg := &logger.LoggerConfig{
		Level:   "info",
		Stdout:  false,
		CtxKeys: []string{"TraceId"},
		Instances: map[string]logger.InstanceConfig{
			"app": {
				Path:              filepath.Join(dir, "logs"),
				File:              "{Y-m-d}.log",
				Format:            "text",
				Level:             "info",
				RotateSize:        "100M",
				RotateExpire:      "24h",
				RotateBackupLimit: 7,
			},
			"error": {
				Path:              filepath.Join(dir, "logs"),
				File:              "error-{Y-m-d}.log",
				Format:            "json",
				Level:             "error",
				RotateSize:        "50M",
				RotateExpire:      "168h",
				RotateBackupLimit: 30,
			},
		},
	}

	// Test successful initialization
	if err := logger.Initialize(cfg); err != nil {
		t.Fatalf("初始化失败: %v", err)
	}

	// Test double initialization
	if err := logger.Initialize(cfg); err == nil {
		t.Fatalf("期望重复初始化报错")
	}

	// Test nil config
	if err := logger.Initialize(nil); err == nil {
		t.Fatalf("期望 nil 配置报错")
	}

	// Test invalid config
	invalidCfg := &logger.LoggerConfig{Level: "invalid"}
	if err := logger.Initialize(invalidCfg); err == nil {
		t.Fatalf("期望无效配置报错")
	}
}

func TestGetInstance(t *testing.T) {
	logger.Reset() // Reset for test isolation

	dir, err := ioutil.TempDir("", "loggertest")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(dir)

	cfg := &logger.LoggerConfig{
		Level: "info",
		Instances: map[string]logger.InstanceConfig{
			"app": {
				Path:  filepath.Join(dir, "logs"),
				File:  "{Y-m-d}.log",
				Level: "info",
			},
		},
	}

	if err := logger.Initialize(cfg); err != nil {
		t.Fatalf("初始化失败: %v", err)
	}

	// Test existing instance
	appLogger := logger.GetInstance("app")
	if appLogger == nil {
		t.Fatalf("期望获取到 app logger")
	}

	// Test non-existing instance (should return default)
	defaultLogger := logger.GetInstance("nonexistent")
	if defaultLogger == nil {
		t.Fatalf("期望获取到默认 logger")
	}
}

func TestGetInstanceWithContext(t *testing.T) {
	logger.Reset() // Reset for test isolation

	dir, err := ioutil.TempDir("", "loggertest")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(dir)

	cfg := &logger.LoggerConfig{
		Level:   "info",
		CtxKeys: []string{"TraceId", "RequestId"},
		Instances: map[string]logger.InstanceConfig{
			"app": {
				Path:  filepath.Join(dir, "logs"),
				File:  "{Y-m-d}.log",
				Level: "info",
			},
		},
	}

	if err := logger.Initialize(cfg); err != nil {
		t.Fatalf("初始化失败: %v", err)
	}

	ctx := context.WithValue(context.Background(), "TraceId", "test-trace-123")

	// Test with context
	ctxLogger := logger.GetInstanceWithContext("app", ctx)
	if ctxLogger == nil {
		t.Fatalf("期望获取到带上下文的 logger")
	}

	// Test with nil context
	nilCtxLogger := logger.GetInstanceWithContext("app", nil)
	if nilCtxLogger == nil {
		t.Fatalf("期望获取到 logger")
	}

	// Test with non-existing instance
	defaultCtxLogger := logger.GetInstanceWithContext("nonexistent", ctx)
	if defaultCtxLogger == nil {
		t.Fatalf("期望获取到默认 logger")
	}
}
