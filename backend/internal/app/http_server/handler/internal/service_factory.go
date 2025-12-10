package internal

import (
	"errors"

	"go-study2/internal/domain/progress"
	"go-study2/internal/domain/quiz"
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

// BuildProgressService 基于全局依赖构建学习进度服务。
func BuildProgressService() (*progress.Service, error) {
	db := database.Default()
	if db == nil {
		return nil, errors.New("数据库未初始化")
	}
	return progress.NewService(repository.NewProgressRepository(db)), nil
}

// BuildQuizService 基于全局依赖构建测验服务。
func BuildQuizService() (*quiz.Service, error) {
	db := database.Default()
	if db == nil {
		return nil, errors.New("数据库未初始化")
	}
	return quiz.NewService(repository.NewQuizRepository(db)), nil
}
