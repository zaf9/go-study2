package repository

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"go-study2/internal/config"
	"go-study2/internal/domain/progress"
	"go-study2/internal/infrastructure/database"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
)

func TestProgressRepository_UpsertAndQuery(t *testing.T) {
	ctx := gctx.New()
	db := setupProgressRepoDB(t)
	repo := NewProgressRepository(db)

	record := &progress.Progress{
		UserID:  1,
		Topic:   "variables",
		Chapter: "storage",
		Status:  progress.StatusInProgress,
	}
	if err := repo.Upsert(ctx, record); err != nil {
		t.Fatalf("写入进度失败: %v", err)
	}

	all, err := repo.ListByUser(ctx, 1)
	if err != nil {
		t.Fatalf("查询全部进度失败: %v", err)
	}
	if len(all) != 1 {
		t.Fatalf("期望 1 条进度记录，得到 %d", len(all))
	}

	record.Status = progress.StatusDone
	record.LastPosition = `{"scroll":500}`
	if err := repo.Upsert(ctx, record); err != nil {
		t.Fatalf("更新进度失败: %v", err)
	}

	topicProgress, err := repo.ListByTopic(ctx, 1, "variables")
	if err != nil {
		t.Fatalf("按主题查询失败: %v", err)
	}
	if len(topicProgress) != 1 {
		t.Fatalf("按主题返回数量不正确")
	}
	if topicProgress[0].Status != progress.StatusDone {
		t.Fatalf("状态未更新，得到 %s", topicProgress[0].Status)
	}
	if topicProgress[0].LastPosition == "" {
		t.Fatalf("LastPosition 应被保存")
	}
}

func setupProgressRepoDB(t *testing.T) gdb.DB {
	t.Helper()
	ensureRepoConfigPath()
	_ = os.MkdirAll("testdata", 0o755)
	dbPath := filepath.ToSlash(filepath.Join("testdata", fmt.Sprintf("progress_%d.db", time.Now().UnixNano())))
	cfg := config.DatabaseConfig{
		Type: "sqlite3",
		Path: dbPath,
		Pragmas: []string{
			"journal_mode=WAL",
			"busy_timeout=5000",
			"synchronous=NORMAL",
			"cache_size=-64000",
			"foreign_keys=ON",
		},
	}
	db, err := database.Init(gctx.New(), cfg)
	if err != nil {
		t.Fatalf("初始化测试数据库失败: %v", err)
	}
	return db
}

func ensureRepoConfigPath() {
	adapter, ok := g.Cfg().GetAdapter().(*gcfg.AdapterFile)
	if !ok {
		return
	}
	if configFile, err := gfile.Search("configs/config.yaml"); err == nil && configFile != "" {
		_ = adapter.SetPath(filepath.Dir(configFile))
		adapter.SetFileName("config")
	}
}

