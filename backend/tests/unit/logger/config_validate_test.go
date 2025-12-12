package logger_test

import (
	"strings"
	"testing"

	lg "go-study2/internal/infrastructure/logger"
)

func TestValidateReturnsAggregatedErrors(t *testing.T) {
	cfg := &lg.LoggerConfig{
		Level: "invalid-level",
		Instances: map[string]lg.InstanceConfig{
			"bad": {
				Path: "", // missing
				File: "", // missing
			},
		},
	}

	err := cfg.Validate()
	if err == nil {
		t.Fatalf("expected validation error but got nil")
	}
	msg := err.Error()
	// Ensure error contains both messages
	if !(strings.Contains(msg, "无效的日志级别") && strings.Contains(msg, "实例 bad: path 不能为空")) {
		t.Fatalf("expected aggregated messages in error: %s", msg)
	}
}
