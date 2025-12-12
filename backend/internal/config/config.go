package config

import (
	"crypto/tls"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
)

// Config 应用配置结构
type Config struct {
	Http     HttpConfig     `json:"http"`
	Https    HttpsConfig    `json:"https"`
	Server   ServerConfig   `json:"server"`
	Logger   LoggerConfig   `json:"logger"`
	Database DatabaseConfig `json:"database"`
	Jwt      JwtConfig      `json:"jwt"`
	Static   StaticConfig   `json:"static"`
}

// HttpConfig HTTP 配置
type HttpConfig struct {
	Port int `json:"port"`
}

// HttpsConfig HTTPS 配置
type HttpsConfig struct {
	Enabled            bool   `json:"enabled"`
	Port               int    `json:"port"`
	CertFile           string `json:"certFile"`
	KeyFile            string `json:"keyFile"`
	CaFile             string `json:"caFile"`
	InsecureSkipVerify bool   `json:"insecureSkipVerify"`
}

// ServerConfig HTTP/HTTPS 通用服务器配置
type ServerConfig struct {
	Host            string `json:"host"`
	ShutdownTimeout int    `json:"shutdownTimeout"`
}

// LoggerConfig 日志配置
type LoggerConfig struct {
	Level  string `json:"level"`
	Path   string `json:"path"`
	Stdout bool   `json:"stdout"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Type            string   `json:"type"`
	Path            string   `json:"path"`
	MaxOpenConns    int      `json:"maxOpenConns"`
	MaxIdleConns    int      `json:"maxIdleConns"`
	ConnMaxLifetime int      `json:"connMaxLifetime"`
	SlowThreshold   int      `json:"slowThreshold"`
	Pragmas         []string `json:"pragmas"`
}

// JwtConfig JWT 配置
type JwtConfig struct {
	Secret             string `json:"secret"`
	AccessTokenExpiry  int64  `json:"accessTokenExpiry"`
	RefreshTokenExpiry int64  `json:"refreshTokenExpiry"`
	Issuer             string `json:"issuer"`
}

// StaticConfig 静态资源配置
type StaticConfig struct {
	Enabled     bool   `json:"enabled"`
	Path        string `json:"path"`
	SpaFallback bool   `json:"spaFallback"`
}

// Load 加载配置文件（默认读取 configs/config.yaml）
func Load() (*Config, error) {
	ctx := gctx.New()
	var cfg Config

	if err := setConfigPath(); err != nil {
		return nil, err
	}

	if err := g.Cfg().MustGet(ctx, ".").Scan(&cfg); err != nil {
		return nil, fmt.Errorf("加载配置文件失败: %w", err)
	}

	if err := Validate(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Validate 验证配置
func Validate(cfg *Config) error {
	if cfg.Server.Host == "" {
		return fmt.Errorf("配置项 server.host 为必填项，请在configs/config.yaml中设置")
	}

	if err := validatePort(cfg.Http.Port, "http.port", !cfg.Https.Enabled); err != nil {
		return err
	}

	if cfg.Https.Enabled {
		if err := validatePort(cfg.Https.Port, "https.port", true); err != nil {
			return err
		}
		if cfg.Https.CertFile == "" {
			return fmt.Errorf("配置项 https.certFile 为必填项（当 https.enabled = true 时）")
		}
		if cfg.Https.KeyFile == "" {
			return fmt.Errorf("配置项 https.keyFile 为必填项（当 https.enabled = true 时）")
		}
		resolvedCert, err := resolvePath(cfg.Https.CertFile)
		if err != nil {
			return err
		}
		resolvedKey, err := resolvePath(cfg.Https.KeyFile)
		if err != nil {
			return err
		}
		cfg.Https.CertFile = resolvedCert
		cfg.Https.KeyFile = resolvedKey

		if err := ensureFile(cfg.Https.CertFile, "证书文件不存在: %s"); err != nil {
			return err
		}
		if err := ensureFile(cfg.Https.KeyFile, "私钥文件不存在: %s"); err != nil {
			return err
		}
		if cfg.Https.CaFile != "" {
			resolvedCA, err := resolvePath(cfg.Https.CaFile)
			if err != nil {
				return err
			}
			cfg.Https.CaFile = resolvedCA
			if err := ensureFile(cfg.Https.CaFile, "CA 证书文件不存在: %s"); err != nil {
				return err
			}
		}
		if err := validateCertKeyPair(cfg.Https.CertFile, cfg.Https.KeyFile); err != nil {
			return err
		}
	}

	if cfg.Database.Type != "" || cfg.Database.Path != "" {
		if cfg.Database.Path == "" {
			return fmt.Errorf("配置项 database.path 为必填项，请在configs/config.yaml中设置")
		}
		if cfg.Database.Type == "" {
			cfg.Database.Type = "sqlite3"
		}
	}

	if cfg.Static.Enabled && cfg.Static.Path == "" {
		return fmt.Errorf("配置项 static.path 为必填项，请在configs/config.yaml中设置")
	}

	return nil
}

// setConfigPath 将配置适配器路径指向 configs 目录
func setConfigPath() error {
	adapter, ok := g.Cfg().GetAdapter().(*gcfg.AdapterFile)
	if !ok {
		return errors.New("配置适配器不是文件类型，无法设置路径")
	}

	configFile, err := gfile.Search("configs/config.yaml")
	if err != nil {
		return fmt.Errorf("查找配置文件失败: %w", err)
	}
	if configFile == "" {
		return fmt.Errorf("未找到配置文件 configs/config.yaml，请确认文件存在")
	}

	configPath := filepath.Dir(configFile)
	if err := adapter.SetPath(configPath); err != nil {
		return fmt.Errorf("设置配置目录失败: %w", err)
	}
	adapter.SetFileName("config")
	return nil
}

// validatePort 校验端口范围
func validatePort(port int, field string, required bool) error {
	if port == 0 {
		if required {
			return fmt.Errorf("配置项 %s 为必填项，请在configs/config.yaml中设置", field)
		}
		return nil
	}
	if port < 1 || port > 65535 {
		return fmt.Errorf("配置项 %s 必须在1-65535范围内", field)
	}
	return nil
}

// ensureFile 校验文件存在性与可读性
func ensureFile(path string, notFoundFmt string) error {
	if path == "" {
		return fmt.Errorf(notFoundFmt, path)
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("解析文件路径失败 %s: %w", path, err)
	}
	info, statErr := os.Stat(absPath)
	if errors.Is(statErr, os.ErrNotExist) {
		return fmt.Errorf(notFoundFmt, absPath)
	}
	if errors.Is(statErr, os.ErrPermission) {
		return fmt.Errorf("无法读取文件（权限不足）: %s", absPath)
	}
	if statErr != nil {
		return fmt.Errorf("读取文件失败 %s: %w", absPath, statErr)
	}
	if info.IsDir() {
		return fmt.Errorf("路径应为文件而非目录: %s", absPath)
	}
	f, openErr := os.Open(absPath)
	if openErr != nil {
		if errors.Is(openErr, os.ErrPermission) {
			return fmt.Errorf("无法读取文件（权限不足）: %s", absPath)
		}
		return fmt.Errorf("无法读取文件: %s, %v", absPath, openErr)
	}
	_ = f.Close()
	return nil
}

// resolvePath 支持相对路径与绝对路径，基于工作目录解析
func resolvePath(p string) (string, error) {
	if filepath.IsAbs(p) {
		return p, nil
	}
	cwd := gfile.Pwd()
	joined := filepath.Join(cwd, p)
	return joined, nil
}

// validateCertKeyPair 校验证书与私钥匹配
func validateCertKeyPair(certFile, keyFile string) error {
	if _, err := tls.LoadX509KeyPair(certFile, keyFile); err != nil {
		return fmt.Errorf("证书与私钥不匹配或无效: %v", err)
	}
	return nil
}
