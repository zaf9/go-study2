package user

import (
	"context"
	"errors"
	"testing"
	"time"

	appjwt "go-study2/internal/pkg/jwt"
)

type mockRepo struct {
	users       map[string]*User
	usersByID   map[int64]*User
	refreshData map[string]RefreshToken
	autoID      int64
}

func newMockRepo() *mockRepo {
	return &mockRepo{
		users:       make(map[string]*User),
		usersByID:   make(map[int64]*User),
		refreshData: make(map[string]RefreshToken),
		autoID:      1,
	}
}

func (m *mockRepo) Create(_ context.Context, u *User) (int64, error) {
	id := m.autoID
	m.autoID++
	clone := *u
	clone.ID = id
	m.users[u.Username] = &clone
	m.usersByID[id] = &clone
	return id, nil
}

func (m *mockRepo) FindByUsername(_ context.Context, username string) (*User, error) {
	if u, ok := m.users[username]; ok {
		clone := *u
		return &clone, nil
	}
	return nil, nil
}

func (m *mockRepo) FindByID(_ context.Context, id int64) (*User, error) {
	if u, ok := m.usersByID[id]; ok {
		clone := *u
		return &clone, nil
	}
	return nil, nil
}

func (m *mockRepo) SaveRefreshToken(_ context.Context, token RefreshToken) error {
	m.refreshData[token.TokenHash] = token
	return nil
}

func (m *mockRepo) DeleteRefreshTokensByUser(_ context.Context, userID int64) error {
	for hash, token := range m.refreshData {
		if token.UserID == userID {
			delete(m.refreshData, hash)
		}
	}
	return nil
}

func (m *mockRepo) FindRefreshToken(_ context.Context, tokenHash string) (*RefreshToken, error) {
	if t, ok := m.refreshData[tokenHash]; ok {
		clone := t
		return &clone, nil
	}
	return nil, nil
}

func TestService_RegisterAndLogin(t *testing.T) {
	_ = appjwt.Configure(appjwt.Options{
		Secret:             "0123456789abcdef",
		AccessTokenExpiry:  time.Hour,
		RefreshTokenExpiry: 24 * time.Hour,
	})

	repo := newMockRepo()
	svc := NewService(repo, time.Hour, 24*time.Hour)
	ctx := context.Background()

	result, err := svc.Register(ctx, "tester_1", "TestPass123")
	if err != nil {
		t.Fatalf("注册失败: %v", err)
	}
	if result.User.ID == 0 || result.Tokens.AccessToken == "" || result.Tokens.RefreshToken == "" {
		t.Fatalf("注册返回数据不完整")
	}

	_, err = svc.Register(ctx, "tester_1", "TestPass123")
	if !errors.Is(err, ErrUserExists) {
		t.Fatalf("重复注册应返回 ErrUserExists")
	}

	login, err := svc.Login(ctx, "tester_1", "TestPass123")
	if err != nil {
		t.Fatalf("登录失败: %v", err)
	}
	if login.Tokens.AccessToken == "" || login.Tokens.RefreshToken == "" {
		t.Fatalf("登录返回的令牌为空")
	}

	_, err = svc.Login(ctx, "tester_1", "WrongPass123")
	if !errors.Is(err, ErrInvalidCredential) {
		t.Fatalf("错误密码应返回 ErrInvalidCredential")
	}
}

func TestService_RefreshAndLogout(t *testing.T) {
	_ = appjwt.Configure(appjwt.Options{
		Secret:             "abcdef0123456789",
		AccessTokenExpiry:  time.Hour,
		RefreshTokenExpiry: time.Hour,
	})

	repo := newMockRepo()
	svc := NewService(repo, time.Hour, time.Hour)
	ctx := context.Background()

	register, err := svc.Register(ctx, "tester_2", "TestPass123")
	if err != nil {
		t.Fatalf("注册失败: %v", err)
	}

	refreshed, err := svc.Refresh(ctx, register.Tokens.RefreshToken)
	if err != nil {
		t.Fatalf("刷新失败: %v", err)
	}
	if refreshed.Tokens.AccessToken == "" || refreshed.Tokens.RefreshToken == "" {
		t.Fatalf("刷新结果缺少令牌")
	}

	_ = repo.SaveRefreshToken(ctx, RefreshToken{
		UserID:    register.User.ID,
		TokenHash: hashToken("expired"),
		ExpiresAt: time.Now().Add(-time.Hour),
	})
	_, err = svc.Refresh(ctx, "expired")
	if !errors.Is(err, ErrRefreshTokenExpired) {
		t.Fatalf("过期刷新令牌应返回 ErrRefreshTokenExpired")
	}

	if err := svc.Logout(ctx, register.User.ID); err != nil {
		t.Fatalf("退出失败: %v", err)
	}
	if len(repo.refreshData) != 0 {
		t.Fatalf("退出后应清空刷新令牌")
	}
}
