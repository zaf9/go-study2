package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/natefinch/lumberjack"
)

var (
	instances    = make(map[string]*glog.Logger)
	writers      = make(map[string]*lumberjack.Logger) // Track writers for cleanup
	mu           sync.RWMutex
	initialized  bool
	globalConfig *LoggerConfig
)

// Level definitions mapping to slog
const (
	LevelCritical = slog.LevelError + 4
)

// TimeFormat returns the configured time format for logs.
func TimeFormat() string {
	if globalConfig != nil && globalConfig.TimeFormat != "" {
		return globalConfig.TimeFormat
	}
	return "2006-01-02T15:04:05.000Z07:00"
}

// Reset resets the logger state and closes all file handles.
func Reset() {
	mu.Lock()
	defer mu.Unlock()

	// Close all lumberjack writers to release file handles
	for _, w := range writers {
		if w != nil {
			_ = w.Close()
		}
	}

	instances = make(map[string]*glog.Logger)
	writers = make(map[string]*lumberjack.Logger)
	initialized = false
	globalConfig = nil
}

// Initialize initializes the logging system.
// It uses slog + lumberjack under the hood to handle rotation stably,
// while providing glog interfaces for compatibility.
func Initialize(config *LoggerConfig) error {
	if initialized {
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

	for name, instCfg := range config.Instances {
		logger := glog.New()

		// 1. Create stable rotation backend with lumberjack
		lj := &lumberjack.Logger{
			Filename:   filepath.Join(instCfg.Path, instCfg.File),
			MaxSize:    parseSizeToMB(instCfg.RotateSize),
			MaxBackups: instCfg.RotateBackupLimit,
			MaxAge:     parseAgeToDays(instCfg.RotateExpire),
			LocalTime:  true,
			Compress:   true,
		}
		// Special handling for GoFrame placeholder filenames: lumberjack doesn't support them.
		// We use the base name as the primary log file name.
		if strings.Contains(instCfg.File, "{") {
			baseName := instCfg.File
			// Strip common GoFrame placeholders to get a stable base filename for lumberjack
			baseName = strings.ReplaceAll(baseName, "{Y-m-d}", "")
			baseName = strings.ReplaceAll(baseName, "{Ymd}", "")
			baseName = strings.Trim(baseName, "-._")
			if baseName == "" || baseName == ".log" {
				baseName = name + ".log"
			}
			if !strings.HasSuffix(baseName, ".log") {
				baseName += ".log"
			}
			lj.Filename = filepath.Join(instCfg.Path, baseName)
		}

		mu.Lock()
		writers[name] = lj
		mu.Unlock()

		// 2. Setup slog handler
		var handler slog.Handler
		opts := &slog.HandlerOptions{
			Level: parseSlogLevel(instCfg.Level),
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.TimeKey {
					return slog.String(slog.TimeKey, a.Value.Time().Format(TimeFormat()))
				}
				return a
			},
		}

		var out io.Writer = lj
		if instCfg.StdoutPrint || config.Stdout {
			out = io.MultiWriter(lj, os.Stdout)
		}

		if instCfg.Format == "json" {
			handler = slog.NewJSONHandler(out, opts)
		} else {
			handler = slog.NewTextHandler(out, opts)
		}

		slogger := slog.New(handler)

		// 3. Configure glog to redirect everything to our slog handler
		// We disable rotation in glog to avoid the panic bug in v2.9.5/v2.9.6.
		logger.SetConfig(glog.Config{
			Level:       glog.LEVEL_ALL, // Ensure all logs pass to our handler
			Path:        "",             // Disable glog internal file writing
			StdoutPrint: false,          // Disable glog internal stdout printing
			RotateSize:  0,              // IMPORTANT: Disable glog rotation to avoid panic
			HeaderPrint: false,          // Disable glog headers (handled by slog)
		})

		// Redirection handler
		logger.SetHandlers(func(ctx context.Context, input *glog.HandlerInput) {
			level := mapGlogLevelToSlog(input.Level)

			// Use the pre-formatted content if available, otherwise join values
			msg := input.Content
			if msg == "" && len(input.Values) > 0 {
				msg = fmt.Sprint(input.Values...)
			}

			slogger.Log(ctx, level, msg)
		})

		mu.Lock()
		instances[name] = logger
		mu.Unlock()

		if name == "app" {
			glog.SetDefaultLogger(logger)
		}
	}

	initialized = true
	return nil
}

// GetInstance returns a logger instance by name.
func GetInstance(name string) *glog.Logger {
	mu.RLock()
	defer mu.RUnlock()
	if logger, ok := instances[name]; ok {
		return logger
	}
	return g.Log()
}

// GetInstanceWithContext returns a logger instance by name with context support.
func GetInstanceWithContext(name string, ctx context.Context) *glog.Logger {
	return GetInstance(name)
}

// parseSizeToMB converts size strings like "100M" or "1G" to MB integers for lumberjack.
func parseSizeToMB(s string) int {
	if s == "" {
		return 100 // Default 100MB
	}
	s = strings.ToUpper(s)
	if strings.HasSuffix(s, "M") || strings.HasSuffix(s, "MB") {
		val, _ := strconv.Atoi(strings.TrimSuffix(strings.TrimSuffix(s, "MB"), "M"))
		return val
	}
	if strings.HasSuffix(s, "G") || strings.HasSuffix(s, "GB") {
		val, _ := strconv.Atoi(strings.TrimSuffix(strings.TrimSuffix(s, "GB"), "G"))
		return val * 1024
	}
	if strings.HasSuffix(s, "K") || strings.HasSuffix(s, "KB") {
		val, _ := strconv.Atoi(strings.TrimSuffix(strings.TrimSuffix(s, "KB"), "K"))
		return val / 1024
	}
	val, _ := strconv.Atoi(s)
	if val > 0 {
		return val / (1024 * 1024)
	}
	return 100
}

// parseAgeToDays converts duration strings like "30d" to day integers for lumberjack.
func parseAgeToDays(s string) int {
	if s == "" {
		return 30 // Default 30 days
	}
	if strings.HasSuffix(s, "d") {
		val, err := strconv.Atoi(strings.TrimSuffix(s, "d"))
		if err == nil {
			return val
		}
	}
	// Fallback to time.Duration for h, m, s
	d, err := time.ParseDuration(s)
	if err == nil {
		return int(d.Hours() / 24)
	}
	return 30
}

func parseSlogLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "all", "debug", "dev":
		return slog.LevelDebug
	case "info", "prod":
		return slog.LevelInfo
	case "notice":
		return slog.LevelInfo // Slog doesn't have Notice, map to Info
	case "warn", "warning":
		return slog.LevelWarn
	case "error", "err":
		return slog.LevelError
	case "crit", "critical":
		return LevelCritical
	default:
		return slog.LevelInfo
	}
}

func mapGlogLevelToSlog(glevel int) slog.Level {
	switch glevel {
	case glog.LEVEL_DEBU:
		return slog.LevelDebug
	case glog.LEVEL_INFO, glog.LEVEL_NOTI:
		return slog.LevelInfo
	case glog.LEVEL_WARN:
		return slog.LevelWarn
	case glog.LEVEL_ERRO:
		return slog.LevelError
	case glog.LEVEL_CRIT:
		return LevelCritical
	default:
		return slog.LevelInfo
	}
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
