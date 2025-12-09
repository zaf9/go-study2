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
				},
				Http: config.HttpConfig{
					Port: 8080,
				},
			},
			wantErr: false,
		},
		{
			name: "缺少Host",
			config: config.Config{
				Server: config.ServerConfig{
					Host: "",
				},
				Http: config.HttpConfig{
					Port: 8080,
				},
			},
			wantErr: true,
		},
		{
			name: "缺少HTTP端口",
			config: config.Config{
				Server: config.ServerConfig{
					Host: "127.0.0.1",
				},
				Http: config.HttpConfig{
					Port: 0,
				},
			},
			wantErr: true,
		},
		{
			name: "HTTP端口超出范围(0)",
			config: config.Config{
				Server: config.ServerConfig{
					Host: "127.0.0.1",
				},
				Http: config.HttpConfig{
					Port: 0,
				},
			},
			wantErr: true,
		},
		{
			name: "HTTP端口超出范围(65536)",
			config: config.Config{
				Server: config.ServerConfig{
					Host: "127.0.0.1",
				},
				Http: config.HttpConfig{
					Port: 65536,
				},
			},
			wantErr: true,
		},
		{
			name: "HTTPS启用缺少证书",
			config: config.Config{
				Server: config.ServerConfig{
					Host: "127.0.0.1",
				},
				Http: config.HttpConfig{
					Port: 8080,
				},
				Https: config.HttpsConfig{
					Enabled: true,
					Port:    8443,
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
	// 测试从 configs/config.yaml 加载
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

	if cfg.Http.Port != 8080 {
		t.Errorf("Expected port 8080, got %d", cfg.Http.Port)
	}
}
