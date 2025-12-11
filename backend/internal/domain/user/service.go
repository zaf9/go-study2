package user

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"regexp"
	"time"

	"go-study2/internal/infrastructure/audit"
	appjwt "go-study2/internal/pkg/jwt"
	"go-study2/internal/pkg/password"
)

// 域内错误定义，便于 handler 做精确映射。
var (
	ErrUserExists          = errors.New("用户名已存在")
	ErrInvalidCredential   = errors.New("用户名或密码错误")
	ErrInvalidInput        = errors.New("参数格式不正确")
	ErrRefreshTokenInvalid = errors.New("refresh token 无效")
	ErrRefreshTokenExpired = errors.New("refresh token 已过期")
	ErrUserNotFound        = errors.New("用户不存在")
	ErrPermissionDenied    = errors.New("权限不足")
	ErrMustChangePassword  = errors.New("需要先修改密码")
)

var (
	usernamePattern = regexp.MustCompile(`^[A-Za-z0-9_]{3,50}$`)
)

const (
	DefaultAdminUsername = "admin"
	DefaultAdminPassword = "GoStudy@123"
	defaultUserStatus    = "active"
)

type createUserOptions struct {
	isAdmin            bool
	mustChangePassword bool
	issueTokens        bool
}

// Service 封装用户注册、登录、刷新令牌等业务能力。
type Service struct {
	repo       Repository
	accessTTL  time.Duration
	refreshTTL time.Duration
}

// NewService 创建服务实例，需传入仓储实现与令牌过期时间。
func NewService(repo Repository, accessTTL, refreshTTL time.Duration) *Service {
	return &Service{
		repo:       repo,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

// Register 由管理员创建新用户，返回令牌对与用户信息。
func (s *Service) Register(ctx context.Context, operatorID int64, username, rawPassword string) (*AuthResult, error) {
	if operatorID <= 0 {
		audit.Record(ctx, "register_denied", 0, "permission_denied", "missing_operator")
		return nil, ErrPermissionDenied
	}
	operator, err := s.repo.FindByID(ctx, operatorID)
	if err != nil {
		return nil, err
	}
	if operator == nil || !operator.IsAdmin {
		audit.Record(ctx, "register_denied", operatorID, "permission_denied", "non_admin_operator")
		return nil, ErrPermissionDenied
	}
	result, regErr := s.registerUser(ctx, username, rawPassword, createUserOptions{
		isAdmin:            false,
		mustChangePassword: false,
		issueTokens:        true,
	})
	if regErr == nil {
		audit.Record(ctx, "register_success", operatorID, "ok", username)
	}
	return result, regErr
}

// EnsureDefaultAdmin 确保默认管理员存在，幂等且不覆盖已存在账户。
func (s *Service) EnsureDefaultAdmin(ctx context.Context) error {
	existing, err := s.repo.FindByUsername(ctx, DefaultAdminUsername)
	if err != nil {
		return err
	}
	if existing != nil {
		return nil
	}

	created, err := s.registerUser(ctx, DefaultAdminUsername, DefaultAdminPassword, createUserOptions{
		isAdmin:            true,
		mustChangePassword: true,
		issueTokens:        false,
	})
	if err == nil && created != nil && created.User != nil {
		audit.Record(ctx, "default_admin_created", created.User.ID, "ok", "")
	}
	return err
}

// Login 验证账号密码并返回新的令牌对。
func (s *Service) Login(ctx context.Context, username, rawPassword string) (*AuthResult, error) {
	if err := s.validateCredential(username, rawPassword); err != nil {
		return nil, ErrInvalidInput
	}

	existing, err := s.repo.FindByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, ErrInvalidCredential
	}

	if err := password.Verify(existing.PasswordHash, rawPassword); err != nil {
		return nil, ErrInvalidCredential
	}

	tokens, err := s.issueTokenPair(ctx, existing.ID)
	if err != nil {
		return nil, err
	}

	return &AuthResult{
		User:   existing,
		Tokens: *tokens,
	}, nil
}

// Refresh 根据刷新令牌换取新的令牌对。
func (s *Service) Refresh(ctx context.Context, refreshToken string) (*AuthResult, error) {
	if refreshToken == "" {
		return nil, ErrRefreshTokenInvalid
	}

	tokenHash := hashToken(refreshToken)
	record, err := s.repo.FindRefreshToken(ctx, tokenHash)
	if err != nil {
		return nil, err
	}
	if record == nil {
		return nil, ErrRefreshTokenInvalid
	}

	now := time.Now()
	if !record.ExpiresAt.IsZero() && record.ExpiresAt.Before(now) {
		return nil, ErrRefreshTokenExpired
	}

	access, err := appjwt.GenerateAccessToken(record.UserID)
	if err != nil {
		return nil, err
	}

	if err := s.repo.DeleteRefreshTokensByUser(ctx, record.UserID); err != nil {
		return nil, err
	}

	newRefresh, err := appjwt.GenerateRefreshToken(record.UserID)
	if err != nil {
		return nil, err
	}

	expiresAt := now.Add(s.refreshTTL)
	saveErr := s.repo.SaveRefreshToken(ctx, RefreshToken{
		UserID:    record.UserID,
		TokenHash: hashToken(newRefresh),
		ExpiresAt: expiresAt,
	})
	if saveErr != nil {
		return nil, saveErr
	}

	return &AuthResult{
		User: &User{
			ID:       record.UserID,
			Username: "",
		},
		Tokens: TokenPair{
			AccessToken:      access,
			AccessExpiresIn:  int64(s.accessTTL.Seconds()),
			RefreshToken:     newRefresh,
			RefreshExpiresAt: expiresAt,
		},
	}, nil
}

// Logout 移除用户关联的刷新令牌。
func (s *Service) Logout(ctx context.Context, userID int64) error {
	if userID <= 0 {
		return ErrInvalidInput
	}
	return s.repo.DeleteRefreshTokensByUser(ctx, userID)
}

// Profile 查询用户基础信息。
func (s *Service) Profile(ctx context.Context, userID int64) (*User, error) {
	if userID <= 0 {
		return nil, ErrInvalidInput
	}
	record, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if record == nil {
		return nil, ErrUserNotFound
	}
	return record, nil
}

// ChangePassword 修改密码并重置需改密标记，清理历史刷新令牌。
func (s *Service) ChangePassword(ctx context.Context, userID int64, oldPassword, newPassword string) error {
	if userID <= 0 {
		return ErrInvalidInput
	}

	record, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if record == nil {
		return ErrUserNotFound
	}

	if err := password.Verify(record.PasswordHash, oldPassword); err != nil {
		return ErrInvalidCredential
	}
	if err := password.Validate(newPassword); err != nil {
		return ErrInvalidInput
	}

	hashed, err := password.Hash(newPassword)
	if err != nil {
		return err
	}

	if err := s.repo.UpdatePasswordAndFlag(ctx, userID, hashed, false); err != nil {
		return err
	}

	audit.Record(ctx, "password_changed", userID, "ok", "")
	return s.repo.DeleteRefreshTokensByUser(ctx, userID)
}

// RefreshTTL 返回刷新令牌有效期。
func (s *Service) RefreshTTL() time.Duration {
	return s.refreshTTL
}

func (s *Service) validateCredential(username, rawPassword string) error {
	if !usernamePattern.MatchString(username) {
		return ErrInvalidInput
	}
	if err := password.Validate(rawPassword); err != nil {
		return ErrInvalidInput
	}
	return nil
}

func (s *Service) issueTokenPair(ctx context.Context, userID int64) (*TokenPair, error) {
	if err := s.repo.DeleteRefreshTokensByUser(ctx, userID); err != nil {
		return nil, err
	}

	access, err := appjwt.GenerateAccessToken(userID)
	if err != nil {
		return nil, err
	}
	refresh, err := appjwt.GenerateRefreshToken(userID)
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(s.refreshTTL)
	if err := s.repo.SaveRefreshToken(ctx, RefreshToken{
		UserID:    userID,
		TokenHash: hashToken(refresh),
		ExpiresAt: expiresAt,
	}); err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:      access,
		AccessExpiresIn:  int64(s.accessTTL.Seconds()),
		RefreshToken:     refresh,
		RefreshExpiresAt: expiresAt,
	}, nil
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func (s *Service) registerUser(ctx context.Context, username, rawPassword string, opts createUserOptions) (*AuthResult, error) {
	if err := s.validateCredential(username, rawPassword); err != nil {
		return nil, ErrInvalidInput
	}

	existing, err := s.repo.FindByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrUserExists
	}

	hashed, err := password.Hash(rawPassword)
	if err != nil {
		return nil, err
	}

	user := &User{
		Username:           username,
		PasswordHash:       hashed,
		IsAdmin:            opts.isAdmin,
		Status:             defaultUserStatus,
		MustChangePassword: opts.mustChangePassword,
	}

	userID, err := s.repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	user.ID = userID

	if !opts.issueTokens {
		return &AuthResult{
			User:   user,
			Tokens: TokenPair{},
		}, nil
	}

	tokens, err := s.issueTokenPair(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &AuthResult{
		User:   user,
		Tokens: *tokens,
	}, nil
}
