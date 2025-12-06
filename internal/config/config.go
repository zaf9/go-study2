package config

import (
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

// Config 应用配置结构
type Config struct {
	Server ServerConfig `json:"server"`
	Logger LoggerConfig `json:"logger"`
}

// ServerConfig HTTP服务器配置
type ServerConfig struct {
	Host            string `json:"host"`
	Port            int    `json:"port"`
	ShutdownTimeout int    `json:"shutdownTimeout"`
}

// LoggerConfig 日志配置
type LoggerConfig struct {
	Level  string `json:"level"`
	Path   string `json:"path"`
	Stdout bool   `json:"stdout"`
}

// Load 加载配置文件
func Load() (*Config, error) {
	ctx := gctx.New()
	var cfg Config

	// 加载配置
	if err := g.Cfg().MustGet(ctx, ".").Scan(&cfg); err != nil {
		return nil, fmt.Errorf("加载配置文件失败: %w", err)
	}

	// 验证必填项
	if err := Validate(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Validate 验证配置
func Validate(cfg *Config) error {
	// 验证Host（必填）
	if cfg.Server.Host == "" {
		return fmt.Errorf("配置项 server.host 为必填项，请在config.yaml中设置")
	}

	// 验证Port（必填且范围检查）
	if cfg.Server.Port == 0 {
		return fmt.Errorf("配置项 server.port 为必填项，请在config.yaml中设置")
	}
	if cfg.Server.Port < 1 || cfg.Server.Port > 65535 {
		return fmt.Errorf("配置项 server.port 必须在1-65535范围内")
	}

	return nil
}
