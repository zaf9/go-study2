package jwt

import (
	"errors"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
)

// Claims 表示业务需要的 JWT 声明。
type Claims struct {
	UserID int64 `json:"uid"`
	jwtlib.RegisteredClaims
}

// Options 配置签名密钥与过期时间。
type Options struct {
	Secret             string
	Issuer             string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
}

var (
	secretKey       []byte
	issuer          = "go-study2"
	accessTokenTTL  = 7 * 24 * time.Hour
	refreshTokenTTL = 7 * 24 * time.Hour
)

// Configure 配置签名参数。
func Configure(opts Options) error {
	if len(opts.Secret) < 16 {
		return errors.New("JWT 密钥长度至少 16 字符")
	}
	secretKey = []byte(opts.Secret)
	if opts.Issuer != "" {
		issuer = opts.Issuer
	}
	if opts.AccessTokenExpiry > 0 {
		accessTokenTTL = opts.AccessTokenExpiry
	}
	if opts.RefreshTokenExpiry > 0 {
		refreshTokenTTL = opts.RefreshTokenExpiry
	}
	return nil
}

// GenerateAccessToken 生成访问令牌。
func GenerateAccessToken(userID int64) (string, error) {
	return signToken(userID, accessTokenTTL)
}

// GenerateRefreshToken 生成刷新令牌。
func GenerateRefreshToken(userID int64) (string, error) {
	return signToken(userID, refreshTokenTTL)
}

// VerifyToken 验证令牌并返回声明。
func VerifyToken(tokenString string) (*Claims, error) {
	if len(secretKey) == 0 {
		return nil, errors.New("JWT 尚未配置密钥")
	}

	token, err := jwtlib.ParseWithClaims(tokenString, &Claims{}, func(token *jwtlib.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("令牌无效")
	}
	return claims, nil
}

func signToken(userID int64, ttl time.Duration) (string, error) {
	if len(secretKey) == 0 {
		return "", errors.New("JWT 尚未配置密钥")
	}
	now := time.Now()
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwtlib.RegisteredClaims{
			ExpiresAt: jwtlib.NewNumericDate(now.Add(ttl)),
			IssuedAt:  jwtlib.NewNumericDate(now),
			Issuer:    issuer,
		},
	}
	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// AccessTokenTTL 返回当前的访问令牌有效期。
func AccessTokenTTL() time.Duration {
	return accessTokenTTL
}

// RefreshTokenTTL 返回当前的刷新令牌有效期。
func RefreshTokenTTL() time.Duration {
	return refreshTokenTTL
}
