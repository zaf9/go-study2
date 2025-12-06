package unit

import (
	"go-study2/internal/config"
	"testing"
)

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  config.Config
		wantErr bool
	}{
		{
			name: "有效配置",
			config: config.Config{
				Server: config.ServerConfig{
					Host: "127.0.0.1",
					Port: 8080,
				},
			},
			wantErr: false,
		},
		{
			name: "缺少Host",
			config: config.Config{
				Server: config.ServerConfig{
					Port: 8080,
				},
			},
			wantErr: true,
		},
		{
			name: "缺少Port",
			config: config.Config{
				Server: config.ServerConfig{
					Host: "127.0.0.1",
				},
			},
			wantErr: true,
		},
		{
			name: "Port超出范围(0)",
			config: config.Config{
				Server: config.ServerConfig{
					Host: "127.0.0.1",
					Port: 0,
				},
			},
			wantErr: true,
		},
		{
			name: "Port超出范围(65536)",
			config: config.Config{
				Server: config.ServerConfig{
					Host: "127.0.0.1",
					Port: 65536,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := config.Validate(&tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfigLoad(t *testing.T) {
	// 测试从项目根目录的 config.yaml 加载
	// 注意：这将依赖于 T001 创建的 config.yaml
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if cfg == nil {
		t.Fatal("Load() returned nil config")
	}

	if cfg.Server.Host != "127.0.0.1" {
		t.Errorf("Expected host 127.0.0.1, got %s", cfg.Server.Host)
	}

	if cfg.Server.Port != 8080 {
		t.Errorf("Expected port 8080, got %d", cfg.Server.Port)
	}
}
