package middleware

import (
	"net/http"
	"strings"

	"go-study2/internal/infrastructure/database"
	appjwt "go-study2/internal/pkg/jwt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
)

// Auth 验证 Bearer Token 的中间件。
func Auth(r *ghttp.Request) {
	// 优先从 Authorization header 获取 token
	tokenString := ""

	authHeader := r.Header.Get("Authorization")
	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		tokenString = strings.TrimPrefix(authHeader, "Bearer ")
	}

	// 如果 header 中没有 token，尝试从 URL 参数获取（用于 WebSocket）
	if tokenString == "" {
		tokenString = r.Get("token").String()
	}

	if tokenString == "" {
		r.Response.WriteStatus(http.StatusUnauthorized)
		r.Response.ClearBuffer()
		r.Response.WriteJson(g.Map{
			"code":    40001,
			"message": "未提供认证令牌",
			"data":    nil,
		})
		r.ExitAll()
		return
	}

	claims, err := appjwt.VerifyToken(tokenString)
	if err != nil {
		r.Response.WriteStatus(http.StatusUnauthorized)
		r.Response.ClearBuffer()
		r.Response.WriteJson(g.Map{
			"code":    40002,
			"message": "令牌无效或已过期",
			"data":    nil,
		})
		r.ExitAll()
		return
	}

	if db := database.Default(); db != nil {
		if count, _ := db.Model("refresh_tokens").Where("user_id", claims.UserID).Count(gctx.New()); count == 0 {
			r.Response.WriteStatus(http.StatusUnauthorized)
			r.Response.ClearBuffer()
			r.Response.WriteJson(g.Map{
				"code":    40002,
				"message": "令牌无效或已过期",
				"data":    nil,
			})
			r.ExitAll()
			return
		}
	}

	r.SetCtxVar("user_id", claims.UserID)
	r.Middleware.Next()
}
