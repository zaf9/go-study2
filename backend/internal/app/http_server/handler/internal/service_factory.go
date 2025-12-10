package internal

import (
	"errors"

	"go-study2/internal/domain/user"
	"go-study2/internal/infrastructure/database"
	"go-study2/internal/infrastructure/repository"
	appjwt "go-study2/internal/pkg/jwt"
)

// BuildUserService 基于全局依赖构建默认用户服务。
func BuildUserService() (*user.Service, error) {
	db := database.Default()
	if db == nil {
		return nil, errors.New("数据库未初始化")
	}
	return user.NewService(repository.NewUserRepository(db), appjwt.AccessTokenTTL(), appjwt.RefreshTokenTTL()), nil
}
