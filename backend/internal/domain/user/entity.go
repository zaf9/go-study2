package user

import "time"

// User 表示系统中的用户实体。
type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// RefreshToken 表示刷新令牌的持久化记录。
type RefreshToken struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"userId"`
	TokenHash string    `json:"tokenHash"`
	ExpiresAt time.Time `json:"expiresAt"`
	CreatedAt time.Time `json:"createdAt"`
}

// TokenPair 封装一对访问令牌与刷新令牌。
type TokenPair struct {
	AccessToken      string    `json:"accessToken"`
	AccessExpiresIn  int64     `json:"expiresIn"`
	RefreshToken     string    `json:"refreshToken"`
	RefreshExpiresAt time.Time `json:"refreshExpiresAt"`
}

// AuthResult 返回用户基础信息与令牌对。
type AuthResult struct {
	User   *User     `json:"user"`
	Tokens TokenPair `json:"tokens"`
}
