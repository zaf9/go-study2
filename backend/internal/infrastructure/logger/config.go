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
	Instances  map[string]InstanceConfig `yaml:"-" json:"instances"`
}

// UnmarshalYAML 自定义 YAML 解析，将 app、access、error、slow 等键解析到 Instances map 中
func (c *LoggerConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// 定义已知的顶级字段
	knownFields := map[string]bool{
		"level":      true,
		"stdout":     true,
		"timeFormat": true,
		"ctxKeys":    true,
	}

	// 先解析已知字段
	type loggerConfigAlias struct {
		Level      string   `yaml:"level"`
		Stdout     bool     `yaml:"stdout"`
		TimeFormat string   `yaml:"timeFormat"`
		CtxKeys    []string `yaml:"ctxKeys"`
	}
	var alias loggerConfigAlias
	if err := unmarshal(&alias); err != nil {
		return err
	}
	c.Level = alias.Level
	c.Stdout = alias.Stdout
	c.TimeFormat = alias.TimeFormat
	c.CtxKeys = alias.CtxKeys

	// 解析所有字段到 map，找出实例配置
	var raw map[string]interface{}
	if err := unmarshal(&raw); err != nil {
		return err
	}

	// 初始化 Instances map
	c.Instances = make(map[string]InstanceConfig)

	// 遍历所有字段，将非已知字段作为实例配置
	for key, value := range raw {
		if !knownFields[key] {
			// 将 value 转换为 InstanceConfig
			var instanceCfg InstanceConfig
			data, err := yaml.Marshal(value)
			if err != nil {
				return fmt.Errorf("解析实例 %s 配置失败: %w", key, err)
			}
			if err := yaml.Unmarshal(data, &instanceCfg); err != nil {
				return fmt.Errorf("解析实例 %s 配置失败: %w", key, err)
			}
			c.Instances[key] = instanceCfg
		}
	}

	return nil
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
// 优先查找 logger.yaml 文件，如果不存在则从 config.yaml 中的 logger 键读取。
func LoadConfig(paths ...string) (*LoggerConfig, error) {
	var target string
	if len(paths) > 0 && paths[0] != "" {
		target = paths[0]
	} else {
		// 优先查找 logger.yaml 文件
		loggerCandidates := []string{
			"backend/configs/logger.yaml",
			"configs/logger.yaml",
		}
		for _, p := range loggerCandidates {
			if _, err := os.Stat(p); err == nil {
				target = p
				break
			}
		}
		// 如果 logger.yaml 不存在，尝试从 config.yaml 中读取
		if target == "" {
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
			return nil, fmt.Errorf("未找到配置文件,请在 backend/configs/logger.yaml 或 backend/configs/config.yaml 中创建")
		}
	}

	b, err := ioutil.ReadFile(target)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 无论是 logger.yaml 还是 config.yaml，都需要从 logger 键中提取配置
	var configFile struct {
		Logger LoggerConfig `yaml:"logger"`
	}
	if err := yaml.Unmarshal(b, &configFile); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}
	cfg := configFile.Logger

	// set defaults
	if cfg.TimeFormat == "" {
		cfg.TimeFormat = "2006-01-02T15:04:05.000Z07:00"
	}
	if len(cfg.CtxKeys) == 0 {
		cfg.CtxKeys = []string{"TraceId"}
	}
	return &cfg, nil
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
