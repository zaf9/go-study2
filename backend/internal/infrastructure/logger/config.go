package logger

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"gopkg.in/yaml.v3"
)

// LoggerConfig 表示整个日志系统的配置
type LoggerConfig struct {
	Level      string                    `yaml:"level" json:"level"`
	Stdout     bool                      `yaml:"stdout" json:"stdout"`
	TimeFormat string                    `yaml:"timeFormat" json:"timeFormat"`
	CtxKeys    []string                  `yaml:"ctxKeys" json:"ctxKeys"`
	Instances  map[string]InstanceConfig `yaml:",inline" json:",inline"`
}

// InstanceConfig 表示单个日志实例的配置
type InstanceConfig struct {
	Path                string `yaml:"path" json:"path"`
	File                string `yaml:"file" json:"file"`
	Format              string `yaml:"format" json:"format"`
	Level               string `yaml:"level" json:"level"`
	RotateSize          string `yaml:"rotateSize" json:"rotateSize"`
	RotateExpire        string `yaml:"rotateExpire" json:"rotateExpire"`
	RotateBackupLimit   int    `yaml:"rotateBackupLimit" json:"rotateBackupLimit"`
	RotateBackupExpire  string `yaml:"rotateBackupExpire" json:"rotateBackupExpire"`
	RotateCheckInterval string `yaml:"rotateCheckInterval" json:"rotateCheckInterval"`
	StdoutPrint         bool   `yaml:"stdoutPrint" json:"stdoutPrint"`
}

var ValidLevels = []string{"all", "dev", "prod", "debug", "info", "notice", "warning", "error", "critical"}

// LoadConfig 从 YAML 文件加载 LoggerConfig.
// 如果没有传入路径,会尝试在若干候选路径中查找配置文件并加载.
func LoadConfig(paths ...string) (*LoggerConfig, error) {
	var target string
	if len(paths) > 0 && paths[0] != "" {
		target = paths[0]
	} else {
		candidates := []string{
			"backend/configs/config.yaml",
			"backend/configs/config.dev.yaml",
			"configs/config.yaml",
			"configs/config.dev.yaml",
		}
		for _, p := range candidates {
			if _, err := os.Stat(p); err == nil {
				target = p
				break
			}
		}
		if target == "" {
			// 尝试从 GoFrame 配置读取
			var cfg2 LoggerConfig
			if err := g.Cfg().MustGet(context.Background(), "logger").Scan(&cfg2); err == nil {
				// set defaults
				if cfg2.TimeFormat == "" {
					cfg2.TimeFormat = "2006-01-02T15:04:05.000Z07:00"
				}
				if len(cfg2.CtxKeys) == 0 {
					cfg2.CtxKeys = []string{"TraceId"}
				}
				return &cfg2, nil
			}
			return nil, fmt.Errorf("未找到配置文件,请在 backend/configs/config.yaml 中创建")
		}
	}

	b, err := ioutil.ReadFile(target)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}
	var cfg struct {
		Logger LoggerConfig `yaml:"logger"`
	}
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}
	// set defaults
	if cfg.Logger.TimeFormat == "" {
		cfg.Logger.TimeFormat = "2006-01-02T15:04:05.000Z07:00"
	}
	if len(cfg.Logger.CtxKeys) == 0 {
		cfg.Logger.CtxKeys = []string{"TraceId"}
	}
	return &cfg.Logger, nil
}

// Validate 验证 LoggerConfig
func (c *LoggerConfig) Validate() error {
	if c == nil {
		return errors.New("LoggerConfig 为空")
	}
	var errs []string
	if !contains(ValidLevels, strings.ToLower(c.Level)) {
		errs = append(errs, fmt.Sprintf("无效的日志级别: %s", c.Level))
	}
	if len(c.Instances) == 0 {
		errs = append(errs, "至少需要配置一个日志实例")
	}
	for name, ins := range c.Instances {
		if err := ins.Validate(name); err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("配置校验失败:\n%s", strings.Join(errs, "\n"))
	}
	return nil
}

// Validate 验证 InstanceConfig
func (ic *InstanceConfig) Validate(name string) error {
	if ic == nil {
		return fmt.Errorf("实例 %s: 配置为空", name)
	}
	if ic.Path == "" {
		return fmt.Errorf("实例 %s: path 不能为空", name)
	}
	if ic.File == "" {
		return fmt.Errorf("实例 %s: file 不能为空", name)
	}
	if ic.Format != "" && ic.Format != "json" && ic.Format != "text" {
		return fmt.Errorf("实例 %s: 无效的格式 %s, 有效值: json, text", name, ic.Format)
	}
	if err := checkDirectoryPermission(ic.Path); err != nil {
		return fmt.Errorf("实例 %s: %v", name, err)
	}
	return nil
}

// checkDirectoryPermission 确保目录存在且可写
func checkDirectoryPermission(path string) error {
	abs := path
	if !filepath.IsAbs(path) {
		wd, _ := os.Getwd()
		abs = filepath.Join(wd, path)
	}
	fi, err := os.Stat(abs)
	if err != nil {
		if os.IsNotExist(err) {
			// 尝试创建目录
			if err := os.MkdirAll(abs, 0o755); err != nil {
				return fmt.Errorf("创建日志目录失败: %w", err)
			}
			return nil
		}
		return fmt.Errorf("无法访问目录: %w", err)
	}
	if !fi.IsDir() {
		return fmt.Errorf("路径不是目录: %s", abs)
	}
	// 尝试创建临时文件以检测写权限
	f, err := os.OpenFile(filepath.Join(abs, ".__perm_test"), os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("目录不可写 (%s): %w", abs, err)
	}
	f.Close()
	_ = os.Remove(filepath.Join(abs, ".__perm_test"))
	return nil
}

func contains(list []string, v string) bool {
	for _, s := range list {
		if s == v {
			return true
		}
	}
	return false
}
