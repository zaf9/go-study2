# JWT 工具包

提供访问令牌与刷新令牌的生成与验证，供 HTTP 认证与集成测试复用。

## 配置

在应用启动时调用 `jwt.Configure`：

```go
import appjwt "go-study2/internal/pkg/jwt"

err := appjwt.Configure(appjwt.Options{
    Secret:             os.Getenv("JWT_SECRET"),
    Issuer:             "go-study2",
    AccessTokenExpiry:  7 * 24 * time.Hour,
    RefreshTokenExpiry: 7 * 24 * time.Hour,
})
```

- `Secret` 必填，长度需 ≥16。
- `Issuer` 可选，默认 `go-study2`。
- 过期时间以秒级或 `time.Duration` 配置。

## 核心方法

- `GenerateAccessToken(userID int64) (string, error)`
- `GenerateRefreshToken(userID int64) (string, error)`
- `VerifyToken(token string) (*Claims, error)`
- `AccessTokenTTL()` / `RefreshTokenTTL()` 返回当前有效期。

## 使用建议

- 密钥从环境变量或配置文件读取，避免硬编码。
- 在测试中为每次运行生成独立密钥，隔离数据。
- 配合 `http_server/middleware/auth.go` 使用，统一返回码与错误消息。

