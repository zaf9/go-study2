package user

import "context"

// Repository 定义用户领域所需的持久化接口。
type Repository interface {
	Create(ctx context.Context, user *User) (int64, error)
	FindByUsername(ctx context.Context, username string) (*User, error)
	FindByID(ctx context.Context, id int64) (*User, error)
	SaveRefreshToken(ctx context.Context, token RefreshToken) error
	DeleteRefreshTokensByUser(ctx context.Context, userID int64) error
	FindRefreshToken(ctx context.Context, tokenHash string) (*RefreshToken, error)
	UpdatePasswordAndFlag(ctx context.Context, userID int64, passwordHash string, mustChange bool) error
}
