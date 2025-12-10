package database

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"go-study2/internal/config"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// Init 建立 SQLite 连接、设置 PRAGMA 并执行迁移。
func Init(ctx context.Context, cfg config.DatabaseConfig) (gdb.DB, error) {
	if cfg.Path == "" {
		return nil, fmt.Errorf("数据库路径未配置")
	}

	dir := filepath.Dir(cfg.Path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("创建数据目录失败: %w", err)
	}

	dbType := cfg.Type
	if dbType == "" {
		dbType = "sqlite3"
	}

	gdb.SetConfig(gdb.Config{
		"default": gdb.ConfigGroup{
			{
				Type: dbType,
				Link: cfg.Path,
			},
		},
	})

	db := g.DB()

	for _, pragma := range cfg.Pragmas {
		if _, err := db.Exec(ctx, pragma); err != nil {
			return nil, fmt.Errorf("设置 PRAGMA 失败: %w", err)
		}
	}

	if err := Migrate(ctx, db); err != nil {
		return nil, fmt.Errorf("数据库迁移失败: %w", err)
	}

	if err := db.PingMaster(); err != nil {
		return nil, fmt.Errorf("数据库连接验证失败: %w", err)
	}

	return db, nil
}
