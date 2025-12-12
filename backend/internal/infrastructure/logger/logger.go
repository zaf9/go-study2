package logger

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfpool"
	"github.com/gogf/gf/v2/os/glog"
)

var (
	instances    = make(map[string]*glog.Logger)
	mu           sync.RWMutex
	initialized  bool
	globalConfig *LoggerConfig
	// cleanupStop is used to signal background cleanup goroutines to stop.
	cleanupStop chan struct{}
)

// TimeFormat returns the configured time format for logs. If no global
// configuration is available, it falls back to a Common Log Format-like
// timestamp used by access logs.
func TimeFormat() string {
	if globalConfig != nil && globalConfig.TimeFormat != "" {
		return globalConfig.TimeFormat
	}
	// Preserve legacy access-log format when no configuration is provided.
	return "02/Jan/2006:15:04:05 -0700"
}

// parseSize parses a size string like "100M" into bytes
func parseSize(s string) (int64, error) {
	if s == "" {
		return 0, nil
	}
	if strings.HasSuffix(s, "M") {
		num, err := strconv.ParseInt(strings.TrimSuffix(s, "M"), 10, 64)
		if err != nil {
			return 0, err
		}
		return num * 1024 * 1024, nil
	}
	if strings.HasSuffix(s, "K") {
		num, err := strconv.ParseInt(strings.TrimSuffix(s, "K"), 10, 64)
		if err != nil {
			return 0, err
		}
		return num * 1024, nil
	}
	// assume bytes if no suffix
	return strconv.ParseInt(s, 10, 64)
}

// Reset resets the logger state for testing purposes.
// This function should only be used in tests.
func Reset() {
	mu.Lock()
	defer mu.Unlock()
	// Signal any background cleanup goroutines to stop first.
	if cleanupStop != nil {
		close(cleanupStop)
		cleanupStop = nil
	}

	// Attempt to gracefully close existing logger instances if they expose a Close method.
	for _, l := range instances {
		if l == nil {
			continue
		}
		// Use reflection to call Close() if available to avoid compile-time
		// dependency on specific glog API versions.
		// This is a best-effort cleanup to release file handles on Windows
		// so that temporary directories can be removed in tests.
		func() {
			defer func() {
				_ = recover()
			}()
			rv := reflect.ValueOf(l)
			// Try a list of common close/release method names used by different versions
			// of glog or custom wrappers. This is best-effort.
			methodNames := []string{"Close", "CloseLogger", "CloseFile", "CloseFiles", "CloseAll", "Destroy", "Stop", "Release", "Shutdown"}
			for _, name := range methodNames {
				m := rv.MethodByName(name)
				if m.IsValid() && m.Type().NumIn() == 0 {
					m.Call(nil)
					break
				}
			}
		}()
	}

	// Best-effort: close any file pointers held by gfpool for configured
	// logger instances. This iterates through configured instance paths and
	// closes pooled file pointers to help Windows release file handles so
	// temporary directories can be removed in tests.
	if globalConfig != nil {
		for _, ins := range globalConfig.Instances {
			abs := ins.Path
			if abs == "" {
				continue
			}
			if !filepath.IsAbs(abs) {
				wd, _ := os.Getwd()
				abs = filepath.Join(wd, abs)
			}
			_ = filepath.WalkDir(abs, func(path string, d os.DirEntry, err error) error {
				if err != nil || d.IsDir() {
					return nil
				}
				// Attempt to get pooled file and close underlying file handle.
				f := gfpool.Get(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
				if f != nil {
					// Close underlying os.File to release handle.
					_ = f.Close(true)
				}
				return nil
			})
		}
	}

	instances = make(map[string]*glog.Logger)
	initialized = false
	globalConfig = nil
	glog.SetDefaultLogger(glog.New())

	// Additionally, attempt to close pooled file pointers under the OS temp
	// directory. This is a broad best-effort sweep to help Windows release
	// lingering file handles created by gfpool during tests.
	if tmp := os.TempDir(); tmp != "" {
		_ = filepath.WalkDir(tmp, func(path string, d os.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return nil
			}
			f := gfpool.Get(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
			if f != nil {
				_ = f.Close(true)
			}
			return nil
		})
	}
}

// Initialize initializes the logging system with the provided configuration.
// It creates logger instances for each configured instance and sets up global defaults.
func Initialize(config *LoggerConfig) error {
	if initialized {
		// If the same config is used for initialization again, return an error
		// to match historical behavior and existing tests that expect a
		// duplicate-initialization to fail. If a different config is provided,
		// reset the existing logger state and proceed with new configuration.
		if reflect.DeepEqual(config, globalConfig) {
			return fmt.Errorf("logger already initialized")
		}
		Reset()
	}
	if config == nil {
		return fmt.Errorf("config cannot be nil")
	}
	if err := config.Validate(); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	globalConfig = config
	instances = make(map[string]*glog.Logger)

	// Initialize cleanupStop to control background cleanup goroutines.
	cleanupStop = make(chan struct{})

	for name, instanceCfg := range config.Instances {
		logger := glog.New()

		// Configure glog with instance configuration
		if err := configureGLog(name, &instanceCfg, config.Stdout, logger); err != nil {
			return fmt.Errorf("failed to configure logger %s: %w", name, err)
		}

		mu.Lock()
		instances[name] = logger
		mu.Unlock()

		// Start cleanup routine for rotateBackupExpire if configured
		if instanceCfg.RotateBackupExpire != "" {
			if d, err := time.ParseDuration(instanceCfg.RotateBackupExpire); err == nil {
				// use rotate check interval if set
				interval := 1 * time.Hour
				if instanceCfg.RotateCheckInterval != "" {
					if iv, err := time.ParseDuration(instanceCfg.RotateCheckInterval); err == nil {
						interval = iv
					}
				}
				// launch goroutine to periodically clean old logs
				go func(p string, expire time.Duration, iv time.Duration) {
					ticker := time.NewTicker(iv)
					defer ticker.Stop()
					for {
						select {
						case <-ticker.C:
							if err := CleanOldLogs(p, expire); err != nil {
								// don't fail initialization for cleanup errors, just log
								log.Printf("logger cleanup error for %s: %v", p, err)
							}
						case <-cleanupStop:
							return
						}
					}
				}(instanceCfg.Path, d, interval)
			}
		}
	}

	// Set default logger if 'app' instance exists
	if logger, ok := instances["app"]; ok {
		glog.SetDefaultLogger(logger)
	}

	initialized = true
	return nil
}

// GetInstance returns a logger instance by name.
// If the instance doesn't exist, it returns the default logger.
func GetInstance(name string) *glog.Logger {
	mu.RLock()
	defer mu.RUnlock()
	if logger, ok := instances[name]; ok {
		return logger
	}
	return g.Log()
}

// GetInstanceWithContext returns a logger instance by name with context support.
// The context is used for TraceID injection if configured.
func GetInstanceWithContext(name string, ctx context.Context) *glog.Logger {
	logger := GetInstance(name)
	// Note: Context is used by the caller when invoking logging methods like logger.InfoCtx(ctx, "message")
	// The logger instance itself doesn't store context, it's passed per log call
	return logger
}

// parseLevel converts string level to glog level
func parseLevel(level string) int {
	switch strings.ToLower(level) {
	case "all":
		return glog.LEVEL_ALL
	case "dev", "debug":
		return glog.LEVEL_DEBU
	case "info":
		return glog.LEVEL_INFO
	case "notic", "notice":
		return glog.LEVEL_NOTI
	case "warn", "warning":
		return glog.LEVEL_WARN
	case "error", "err":
		return glog.LEVEL_ERRO
	case "crit", "critical":
		return glog.LEVEL_CRIT
	default:
		return glog.LEVEL_INFO
	}
}

// configureGLog applies instance configuration to a glog.Logger
func configureGLog(name string, instanceCfg *InstanceConfig, globalStdout bool, logger *glog.Logger) error {
	level := parseLevel(instanceCfg.Level)

	rotateSize, err := parseSize(instanceCfg.RotateSize)
	if err != nil {
		return fmt.Errorf("invalid rotateSize: %w", err)
	}

	cfg := glog.Config{
		Path:                instanceCfg.Path,
		File:                instanceCfg.File,
		Level:               level,
		StdoutPrint:         globalStdout || instanceCfg.StdoutPrint,
		HeaderPrint:         true,
		RotateSize:          rotateSize,
		RotateExpire:        0,
		RotateBackupLimit:   instanceCfg.RotateBackupLimit,
		RotateCheckInterval: 60 * time.Second,
		Flags:               glog.F_TIME_DATE | glog.F_TIME_TIME | glog.F_FILE_LONG,
	}

	// If stdout is enabled globally or for this instance, avoid opening file handles
	// by clearing Path/File so that the underlying glog does not create file writers
	// that can keep temp files locked on Windows during tests.
	if cfg.StdoutPrint {
		cfg.Path = ""
		cfg.File = ""
	}

	// Note: in the past we forced stdout logging on Windows to avoid file handle
	// locking during tests. That had the side-effect of preventing integration
	// tests from validating file-based logging. We now prefer to respect the
	// provided configuration (instanceCfg.StdoutPrint / globalStdout) and rely
	// on Reset() to close file handles during tests when needed.

	if instanceCfg.RotateExpire != "" {
		if d, err := time.ParseDuration(instanceCfg.RotateExpire); err == nil {
			cfg.RotateExpire = d
		}
	}

	if instanceCfg.RotateCheckInterval != "" {
		if iv, err := time.ParseDuration(instanceCfg.RotateCheckInterval); err == nil {
			cfg.RotateCheckInterval = iv
		}
	}

	if strings.ToLower(instanceCfg.Format) == "json" {
		cfg.Flags |= glog.F_FILE_SHORT
	}

	if err := logger.SetConfig(cfg); err != nil {
		return err
	}
	// Note: intentionally not printing per-instance config to avoid noisy test output.
	return nil
}

// CleanOldLogs removes files under dir older than expire duration.
// It ignores directories and only removes regular files.
func CleanOldLogs(dir string, expire time.Duration) error {
	abs := dir
	if !filepath.IsAbs(dir) {
		wd, _ := os.Getwd()
		abs = filepath.Join(wd, dir)
	}
	entries, err := os.ReadDir(abs)
	if err != nil {
		return fmt.Errorf("read dir %s failed: %w", abs, err)
	}
	now := time.Now()
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		info, err := e.Info()
		if err != nil {
			// skip unreadable
			continue
		}
		if now.Sub(info.ModTime()) > expire {
			p := filepath.Join(abs, e.Name())
			_ = os.Remove(p) // best-effort
		}
	}
	return nil
}
