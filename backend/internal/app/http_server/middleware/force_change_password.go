package middleware

import (
	"net/http"

	"go-study2/internal/infrastructure/audit"
	"go-study2/internal/infrastructure/database"
	"go-study2/internal/infrastructure/repository"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
)

// ForceChangePassword 拦截仍需改密的用户，限制访问除改密与资料查询外的接口。
func ForceChangePassword(r *ghttp.Request) {
	userID := r.GetCtxVar("user_id").Int64()
	if userID <= 0 {
		r.Middleware.Next()
		return
	}

	if allowWithoutChange(r.URL.Path) {
		r.Middleware.Next()
		return
	}

	db := database.Default()
	if db == nil {
		r.Middleware.Next()
		return
	}
	repo := repository.NewUserRepository(db)

	user, err := repo.FindByID(gctx.New(), userID)
	if err != nil {
		g.Log().Error(gctx.New(), err)
		writeNeedChangePassword(r)
		return
	}
	if user == nil {
		writeNeedChangePassword(r)
		return
	}
	if user.MustChangePassword {
		audit.Record(r.GetCtx(), "access_blocked_need_change", userID, "blocked", r.URL.Path)
		writeNeedChangePassword(r)
		return
	}

	r.Middleware.Next()
}

func allowWithoutChange(path string) bool {
	if path == "/api/v1/auth/change-password" {
		return true
	}
	if path == "/api/v1/auth/profile" {
		return true
	}
	if path == "/api/v1/auth/logout" {
		return true
	}
	return false
}

func writeNeedChangePassword(r *ghttp.Request) {
	r.Response.WriteStatus(http.StatusForbidden)
	r.Response.ClearBuffer()
	r.Response.WriteJson(g.Map{
		"code":    40011,
		"message": "需要先修改密码",
		"data":    nil,
	})
	r.ExitAll()
}
