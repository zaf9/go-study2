package internal

import (
	"errors"

	progapp "go-study2/internal/app/progress"
	appquiz "go-study2/internal/app/quiz"
	quizdom "go-study2/internal/domain/quiz"
	"go-study2/internal/domain/user"
	infrarepo "go-study2/internal/infra/repository"
	"go-study2/internal/infrastructure/database"
	legacyrepo "go-study2/internal/infrastructure/repository"
	appjwt "go-study2/internal/pkg/jwt"
)

// BuildUserService 基于全局依赖构建默认用户服务。
func BuildUserService() (*user.Service, error) {
	db := database.Default()
	if db == nil {
		return nil, errors.New("数据库未初始化")
	}
	return user.NewService(legacyrepo.NewUserRepository(db), appjwt.AccessTokenTTL(), appjwt.RefreshTokenTTL()), nil
}

// BuildProgressService 基于全局依赖构建学习进度服务（新模型）。
func BuildProgressService() (*progapp.Service, error) {
	db := database.Default()
	if db == nil {
		return nil, errors.New("数据库未初始化")
	}
	calc := progapp.NewCalculator(nil, nil)
	return progapp.NewService(infrarepo.NewProgressRepository(db), calc), nil
}

// BuildQuizService 基于全局依赖构建测验服务。
func BuildQuizService() (*appquiz.Service, error) {
	db := database.Default()
	if db == nil {
		return nil, errors.New("数据库未初始化")
	}

	// YAML 题库由全局单例提供，在 main.go 中初始化
	yamlRepo := quizdom.GetGlobalRepository()
	if yamlRepo == nil {
		return nil, errors.New("YAML 题库未初始化")
	}

	// 会话和答题记录仍存储在数据库
	sessionRepo := infrarepo.NewQuizRepository(db)

	return appquiz.NewService(yamlRepo, sessionRepo), nil
}
