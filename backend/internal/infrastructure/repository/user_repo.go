package repository

import (
	"context"

	"go-study2/internal/domain/user"

	"github.com/gogf/gf/v2/database/gdb"
)

// UserRepository 使用 GoFrame gdb 实现用户仓储。
type UserRepository struct {
	db gdb.DB
}

// NewUserRepository 创建仓储实例。
func NewUserRepository(db gdb.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create 写入新用户并返回自增 ID。
func (r *UserRepository) Create(ctx context.Context, entity *user.User) (int64, error) {
	res, err := r.db.Insert(ctx, "users", map[string]interface{}{
		"username":      entity.Username,
		"password_hash": entity.PasswordHash,
	})
	if err != nil {
		return 0, err
	}
	newID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return newID, nil
}

// FindByUsername 按用户名查询用户。
func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*user.User, error) {
	record, err := r.db.Model("users").Where("username = ?", username).One(ctx)
	if err != nil {
		return nil, err
	}
	if record == nil || len(record.Map()) == 0 {
		return nil, nil
	}
	var entity user.User
	if err := record.Struct(&entity); err != nil {
		return nil, err
	}
	return &entity, nil
}

// FindByID 按 ID 查询用户。
func (r *UserRepository) FindByID(ctx context.Context, id int64) (*user.User, error) {
	record, err := r.db.Model("users").Where("id = ?", id).One(ctx)
	if err != nil {
		return nil, err
	}
	if record == nil || len(record.Map()) == 0 {
		return nil, nil
	}
	var entity user.User
	if err := record.Struct(&entity); err != nil {
		return nil, err
	}
	return &entity, nil
}

// SaveRefreshToken 持久化刷新令牌哈希。
func (r *UserRepository) SaveRefreshToken(ctx context.Context, token user.RefreshToken) error {
	_, err := r.db.Insert(ctx, "refresh_tokens", map[string]interface{}{
		"user_id":    token.UserID,
		"token_hash": token.TokenHash,
		"expires_at": token.ExpiresAt,
	})
	return err
}

// DeleteRefreshTokensByUser 删除指定用户的刷新令牌。
func (r *UserRepository) DeleteRefreshTokensByUser(ctx context.Context, userID int64) error {
	_, err := r.db.Exec(ctx, "DELETE FROM refresh_tokens WHERE user_id = ?", userID)
	return err
}

// FindRefreshToken 通过哈希查询刷新令牌记录。
func (r *UserRepository) FindRefreshToken(ctx context.Context, tokenHash string) (*user.RefreshToken, error) {
	record, err := r.db.Model("refresh_tokens").Where("token_hash = ?", tokenHash).One(ctx)
	if err != nil {
		return nil, err
	}
	if record == nil || len(record.Map()) == 0 {
		return nil, nil
	}
	var entity user.RefreshToken
	if err := record.Struct(&entity); err != nil {
		return nil, err
	}
	return &entity, nil
}
