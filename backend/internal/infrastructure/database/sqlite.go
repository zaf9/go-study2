package database

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go-study2/internal/config"
	"go-study2/internal/infrastructure/db_logging"

	_ "github.com/gogf/gf/contrib/drivers/sqlite/v2"
	_ "github.com/mattn/go-sqlite3"

	"github.com/gogf/gf/v2/database/gdb"
)

var defaultDB gdb.DB

// Init 建立 SQLite 连接、设置 PRAGMA 并执行迁移。
func Init(ctx context.Context, cfg config.DatabaseConfig) (gdb.DB, error) {
	if cfg.Path == "" {
		return nil, fmt.Errorf("数据库路径未配置")
	}

	dir := filepath.Dir(cfg.Path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("创建数据目录失败: %w", err)
	}

	absPath, err := filepath.Abs(cfg.Path)
	if err != nil {
		return nil, fmt.Errorf("解析数据库路径失败: %w", err)
	}

	dbType := cfg.Type
	if dbType == "" {
		dbType = "sqlite"
	}
	if strings.HasPrefix(dbType, "sqlite3") {
		dbType = "sqlite"
	}
	link := fmt.Sprintf("sqlite::@file(%s)", filepath.ToSlash(absPath))

	gdb.SetConfig(gdb.Config{
		"default": gdb.ConfigGroup{
			{
				Type: dbType,
				Link: link,
			},
		},
	})

	db, err := gdb.New(gdb.ConfigNode{
		Type: dbType,
		Link: link,
	})
	if err != nil {
		return nil, fmt.Errorf("创建数据库连接失败: %w", err)
	}
	defaultDB = db

	for _, pragma := range cfg.Pragmas {
		stmt := strings.TrimSpace(pragma)
		if !strings.HasPrefix(strings.ToUpper(stmt), "PRAGMA") {
			stmt = "PRAGMA " + stmt
		}
		if _, err := db.Exec(ctx, stmt); err != nil {
			return nil, fmt.Errorf("设置 PRAGMA 失败: %w", err)
		}
	}

	if err := Migrate(ctx, db); err != nil {
		return nil, fmt.Errorf("数据库迁移失败: %w", err)
	}

	if err := db.PingMaster(); err != nil {
		return nil, fmt.Errorf("数据库连接验证失败: %w", err)
	}

	// 注册数据库日志处理器
	slowThreshold := time.Duration(cfg.SlowThreshold) * time.Millisecond
	if slowThreshold <= 0 {
		slowThreshold = 100 * time.Millisecond // 默认100毫秒
	}
	if err := db_logging.RegisterDBLogging(db, slowThreshold); err != nil {
		return nil, fmt.Errorf("注册数据库日志处理器失败: %w", err)
	}

	return db, nil
}

// Default 返回通过 Init 初始化的默认数据库实例。
func Default() gdb.DB {
	return defaultDB
}
