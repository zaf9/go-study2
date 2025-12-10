package repository

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"go-study2/internal/config"
	"go-study2/internal/domain/user"
	"go-study2/internal/infrastructure/database"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
)

func TestUserRepository_CRUD(t *testing.T) {
	ctx := gctx.New()
	db := setupRepoDB(t)
	repo := NewUserRepository(db)

	userID, err := repo.Create(ctx, &user.User{
		Username:     "repo_user",
		PasswordHash: "hash",
	})
	if err != nil {
		t.Fatalf("创建用户失败: %v", err)
	}
	if userID == 0 {
		t.Fatalf("返回的用户 ID 无效")
	}

	found, err := repo.FindByUsername(ctx, "repo_user")
	if err != nil {
		t.Fatalf("查询用户失败: %v", err)
	}
	if found == nil || found.ID != userID {
		t.Fatalf("查询结果与创建的用户不一致")
	}

	foundByID, err := repo.FindByID(ctx, userID)
	if err != nil {
		t.Fatalf("按 ID 查询失败: %v", err)
	}
	if foundByID == nil || foundByID.Username != "repo_user" {
		t.Fatalf("按 ID 查询返回为空或用户名不匹配")
	}

	expires := time.Now().Add(time.Hour)
	tokenHash := "token-hash"
	saveErr := repo.SaveRefreshToken(ctx, user.RefreshToken{
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: expires,
	})
	if saveErr != nil {
		t.Fatalf("保存刷新令牌失败: %v", saveErr)
	}

	refresh, err := repo.FindRefreshToken(ctx, tokenHash)
	if err != nil {
		t.Fatalf("查询刷新令牌失败: %v", err)
	}
	if refresh == nil || refresh.UserID != userID {
		t.Fatalf("刷新令牌查询结果不正确")
	}

	if err := repo.DeleteRefreshTokensByUser(ctx, userID); err != nil {
		t.Fatalf("删除刷新令牌失败: %v", err)
	}
	afterDelete, err := repo.FindRefreshToken(ctx, tokenHash)
	if err != nil {
		t.Fatalf("删除后查询刷新令牌失败: %v", err)
	}
	if afterDelete != nil {
		t.Fatalf("删除后刷新令牌仍然存在")
	}
}

func setupRepoDB(t *testing.T) gdb.DB {
	t.Helper()
	ensureConfigPath()
	_ = os.MkdirAll("testdata", 0o755)
	dbPath := filepath.ToSlash(filepath.Join("testdata", fmt.Sprintf("repo_%d.db", time.Now().UnixNano())))
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

func ensureConfigPath() {
	adapter, ok := g.Cfg().GetAdapter().(*gcfg.AdapterFile)
	if !ok {
		return
	}
	if configFile, err := gfile.Search("configs/config.yaml"); err == nil && configFile != "" {
		_ = adapter.SetPath(filepath.Dir(configFile))
		adapter.SetFileName("config")
	}
}
