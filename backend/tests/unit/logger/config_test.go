package logger_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"go-study2/internal/infrastructure/logger"
)

func TestLoadAndValidateConfig(t *testing.T) {
	dir, err := ioutil.TempDir("", "logcfgtest")
	if err != nil {
		t.Fatalf("临时目录创建失败: %v", err)
	}
	defer os.RemoveAll(dir)

	cfgYml := `logger:
  level: "all"
  stdout: true
  timeFormat: "2006-01-02T15:04:05.000Z07:00"
  ctxKeys: ["TraceId"]
  default:
    path: "./backend/tests/unit/logger/tmp_logs"
    file: "{Y-m-d}.log"
    format: "json"
`

	f := filepath.Join(dir, "config.yaml")
	if err := ioutil.WriteFile(f, []byte(cfgYml), 0644); err != nil {
		t.Fatalf("写入配置文件失败: %v", err)
	}

	cfg, err := logger.LoadConfig(f)
	if err != nil {
		t.Fatalf("LoadConfig 失败: %v", err)
	}
	if cfg.Level != "all" {
		t.Fatalf("期望 level=all, got=%s", cfg.Level)
	}

	// validate should create the path under repo root; call Validate
	if err := cfg.Validate(); err != nil {
		t.Fatalf("配置验证失败: %v", err)
	}
}

func TestInvalidLevel(t *testing.T) {
	cfg := &logger.LoggerConfig{Level: "invalid", Instances: map[string]logger.InstanceConfig{"default": {Path: "./logs", File: "a.log"}}}
	if err := cfg.Validate(); err == nil {
		t.Fatalf("期望无效级别报错")
	}
}
