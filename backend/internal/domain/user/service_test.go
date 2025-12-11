package user

import (
	"context"
	"errors"
	"testing"
	"time"

	appjwt "go-study2/internal/pkg/jwt"
	"go-study2/internal/pkg/password"
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
	if clone.Status == "" {
		clone.Status = defaultUserStatus
	}
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

func (m *mockRepo) UpdatePasswordAndFlag(_ context.Context, userID int64, passwordHash string, mustChange bool) error {
	if u, ok := m.usersByID[userID]; ok {
		clone := *u
		clone.PasswordHash = passwordHash
		clone.MustChangePassword = mustChange
		m.usersByID[userID] = &clone
		m.users[u.Username] = &clone
		return nil
	}
	return errors.New("user not found")
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

	adminID, _ := repo.Create(ctx, &User{
		Username: "admin",
		IsAdmin:  true,
		Status:   "active",
	})

	result, err := svc.Register(ctx, adminID, "tester_1", "TestPass123!")
	if err != nil {
		t.Fatalf("注册失败: %v", err)
	}
	if result.User.ID == 0 || result.Tokens.AccessToken == "" || result.Tokens.RefreshToken == "" {
		t.Fatalf("注册返回数据不完整")
	}

	_, err = svc.Register(ctx, adminID, "tester_1", "TestPass123!")
	if !errors.Is(err, ErrUserExists) {
		t.Fatalf("重复注册应返回 ErrUserExists")
	}

	login, err := svc.Login(ctx, "tester_1", "TestPass123!")
	if err != nil {
		t.Fatalf("登录失败: %v", err)
	}
	if login.Tokens.AccessToken == "" || login.Tokens.RefreshToken == "" {
		t.Fatalf("登录返回的令牌为空")
	}

	_, err = svc.Login(ctx, "tester_1", "WrongPass123!")
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

	adminID, _ := repo.Create(ctx, &User{
		Username: "admin",
		IsAdmin:  true,
		Status:   "active",
	})

	register, err := svc.Register(ctx, adminID, "tester_2", "TestPass123!")
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

func TestService_RejectsWeakPassword(t *testing.T) {
	_ = appjwt.Configure(appjwt.Options{
		Secret:             "abcdef0123456789",
		AccessTokenExpiry:  time.Hour,
		RefreshTokenExpiry: time.Hour,
	})

	repo := newMockRepo()
	svc := NewService(repo, time.Hour, time.Hour)
	ctx := context.Background()

	adminID, _ := repo.Create(ctx, &User{
		Username: "admin",
		IsAdmin:  true,
		Status:   "active",
	})

	_, err := svc.Register(ctx, adminID, "weakuser", "Weakpass1")
	if !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("弱口令应返回 ErrInvalidInput")
	}
}

func TestService_EnsureDefaultAdmin_CreatesWhenMissing(t *testing.T) {
	_ = appjwt.Configure(appjwt.Options{
		Secret:             "abcdef0123456789",
		AccessTokenExpiry:  time.Hour,
		RefreshTokenExpiry: time.Hour,
	})

	repo := newMockRepo()
	svc := NewService(repo, time.Hour, time.Hour)
	ctx := context.Background()

	err := svc.EnsureDefaultAdmin(ctx)
	if err != nil {
		t.Fatalf("创建默认管理员失败: %v", err)
	}
	created, _ := repo.FindByUsername(ctx, DefaultAdminUsername)
	if created == nil || !created.IsAdmin || !created.MustChangePassword {
		t.Fatalf("默认管理员字段不符合预期")
	}
}

func TestService_EnsureDefaultAdmin_Idempotent(t *testing.T) {
	_ = appjwt.Configure(appjwt.Options{
		Secret:             "abcdef0123456789",
		AccessTokenExpiry:  time.Hour,
		RefreshTokenExpiry: time.Hour,
	})

	repo := newMockRepo()
	svc := NewService(repo, time.Hour, time.Hour)
	ctx := context.Background()

	hashed, _ := password.Hash("Admin123!")
	_, _ = repo.Create(ctx, &User{
		Username:           DefaultAdminUsername,
		PasswordHash:       hashed,
		IsAdmin:            true,
		Status:             defaultUserStatus,
		MustChangePassword: false,
	})

	err := svc.EnsureDefaultAdmin(ctx)
	if err != nil {
		t.Fatalf("幂等检查失败: %v", err)
	}
	after, _ := repo.FindByUsername(ctx, DefaultAdminUsername)
	if after == nil || after.PasswordHash != hashed || after.MustChangePassword {
		t.Fatalf("已有管理员不应被覆盖")
	}
}

func TestService_Register_PermissionDenied(t *testing.T) {
	_ = appjwt.Configure(appjwt.Options{
		Secret:             "abcdef0123456789",
		AccessTokenExpiry:  time.Hour,
		RefreshTokenExpiry: time.Hour,
	})

	repo := newMockRepo()
	svc := NewService(repo, time.Hour, time.Hour)
	ctx := context.Background()

	nonAdminID, _ := repo.Create(ctx, &User{
		Username: "user1",
		Status:   "active",
	})

	_, err := svc.Register(ctx, nonAdminID, "tester_x", "TestPass123!")
	if !errors.Is(err, ErrPermissionDenied) {
		t.Fatalf("非管理员注册应拒绝，得到: %v", err)
	}
}

func TestService_ChangePassword_Success(t *testing.T) {
	_ = appjwt.Configure(appjwt.Options{
		Secret:             "abcdef0123456789",
		AccessTokenExpiry:  time.Hour,
		RefreshTokenExpiry: time.Hour,
	})

	repo := newMockRepo()
	svc := NewService(repo, time.Hour, time.Hour)
	ctx := context.Background()

	oldHash := hashOrFail(t, "OldPass123!")
	userID, _ := repo.Create(ctx, &User{
		Username:           "need_change",
		Status:             defaultUserStatus,
		PasswordHash:       oldHash,
		MustChangePassword: true,
	})
	_ = repo.SaveRefreshToken(ctx, RefreshToken{
		UserID:    userID,
		TokenHash: hashToken("old_refresh"),
		ExpiresAt: time.Now().Add(time.Hour),
	})

	if err := svc.ChangePassword(ctx, userID, "OldPass123!", "NewPass123!"); err != nil {
		t.Fatalf("改密失败: %v", err)
	}

	updated, _ := repo.FindByID(ctx, userID)
	if updated == nil || updated.MustChangePassword {
		t.Fatalf("改密后需改密标记应为 false")
	}
	if updated.PasswordHash == "" || updated.PasswordHash == oldHash {
		t.Fatalf("密码未更新")
	}
	if len(repo.refreshData) != 0 {
		t.Fatalf("改密后应清理刷新令牌")
	}
}

func TestService_ChangePassword_InvalidOldPassword(t *testing.T) {
	_ = appjwt.Configure(appjwt.Options{
		Secret:             "abcdef0123456789",
		AccessTokenExpiry:  time.Hour,
		RefreshTokenExpiry: time.Hour,
	})

	repo := newMockRepo()
	svc := NewService(repo, time.Hour, time.Hour)
	ctx := context.Background()

	userID, _ := repo.Create(ctx, &User{
		Username:           "need_change",
		Status:             defaultUserStatus,
		PasswordHash:       hashOrFail(t, "OldPass123!"),
		MustChangePassword: true,
	})

	err := svc.ChangePassword(ctx, userID, "wrongOld!", "NewPass123!")
	if !errors.Is(err, ErrInvalidCredential) {
		t.Fatalf("错误旧密码应返回 ErrInvalidCredential")
	}
}

func hashOrFail(t *testing.T, raw string) string {
	t.Helper()
	hashed, err := password.Hash(raw)
	if err != nil {
		t.Fatalf("hash 失败: %v", err)
	}
	return hashed
}
