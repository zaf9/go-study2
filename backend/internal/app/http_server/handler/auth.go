package handler

import (
	"net/http"
	"time"

	"go-study2/internal/app/http_server/handler/internal"
	"go-study2/internal/domain/user"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type authRequest struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	RememberMe bool   `json:"rememberMe"`
	Remember   bool   `json:"remember"`
}

type authResponse struct {
	AccessToken        string `json:"accessToken"`
	ExpiresIn          int64  `json:"expiresIn"`
	NeedPasswordChange bool   `json:"needPasswordChange"`
	IsAdmin            bool   `json:"isAdmin"`
}

type profileResponse struct {
	ID                 int64  `json:"id"`
	Username           string `json:"username"`
	IsAdmin            bool   `json:"isAdmin"`
	MustChangePassword bool   `json:"mustChangePassword"`
}

type changePasswordRequest struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

// Register 处理用户注册。
func (h *Handler) Register(r *ghttp.Request) {
	svc, err := h.ensureUserService()
	if err != nil {
		writeError(r, http.StatusInternalServerError, 50001, "认证服务不可用")
		return
	}

	req, ok := parseAuthRequest(r)
	if !ok {
		return
	}

	operatorID := r.GetCtxVar("user_id").Int64()
	if operatorID <= 0 {
		writeError(r, http.StatusUnauthorized, 40001, "缺少管理员凭证")
		return
	}

	result, err := svc.Register(r.GetCtx(), operatorID, req.Username, req.Password)
	if err != nil {
		writeAuthError(r, err)
		return
	}

	if operatorID == result.User.ID {
		h.setRefreshCookie(r, svc, result.Tokens.RefreshToken, req.isRemember())
	}
	writeSuccess(r, "注册成功", authResponse{
		AccessToken:        result.Tokens.AccessToken,
		ExpiresIn:          result.Tokens.AccessExpiresIn,
		NeedPasswordChange: result.User.MustChangePassword,
		IsAdmin:            result.User.IsAdmin,
	})
}

// Login 处理用户登录。
func (h *Handler) Login(r *ghttp.Request) {
	svc, err := h.ensureUserService()
	if err != nil {
		writeError(r, http.StatusInternalServerError, 50001, "认证服务不可用")
		return
	}

	req, ok := parseAuthRequest(r)
	if !ok {
		return
	}

	result, err := svc.Login(r.GetCtx(), req.Username, req.Password)
	if err != nil {
		writeAuthError(r, err)
		return
	}

	h.setRefreshCookie(r, svc, result.Tokens.RefreshToken, req.isRemember())
	writeSuccess(r, "登录成功", authResponse{
		AccessToken:        result.Tokens.AccessToken,
		ExpiresIn:          result.Tokens.AccessExpiresIn,
		NeedPasswordChange: result.User.MustChangePassword,
		IsAdmin:            result.User.IsAdmin,
	})
}

// Logout 处理退出登录。
func (h *Handler) Logout(r *ghttp.Request) {
	svc, err := h.ensureUserService()
	if err != nil {
		writeError(r, http.StatusInternalServerError, 50001, "认证服务不可用")
		return
	}

	userID := r.GetCtxVar("user_id").Int64()
	if userID <= 0 {
		writeError(r, http.StatusUnauthorized, 40001, "认证信息缺失")
		return
	}

	if err := svc.Logout(r.GetCtx(), userID); err != nil {
		writeAuthError(r, err)
		return
	}

	h.clearRefreshCookie(r)
	writeSuccess(r, "退出成功", nil)
}

// ChangePassword 处理已登录用户修改密码，完成后要求重新登录。
func (h *Handler) ChangePassword(r *ghttp.Request) {
	svc, err := h.ensureUserService()
	if err != nil {
		writeError(r, http.StatusInternalServerError, 50001, "认证服务不可用")
		return
	}

	userID := r.GetCtxVar("user_id").Int64()
	if userID <= 0 {
		writeError(r, http.StatusUnauthorized, 40001, "认证信息缺失")
		return
	}

	var req changePasswordRequest
	if err := r.Parse(&req); err != nil || req.OldPassword == "" || req.NewPassword == "" {
		writeError(r, http.StatusBadRequest, 40004, "请求参数无效")
		return
	}

	if err := svc.ChangePassword(r.GetCtx(), userID, req.OldPassword, req.NewPassword); err != nil {
		writeAuthError(r, err)
		return
	}

	h.clearRefreshCookie(r)
	writeSuccess(r, "密码修改成功，请重新登录", nil)
}

// RefreshToken 通过刷新令牌换取新的访问令牌。
func (h *Handler) RefreshToken(r *ghttp.Request) {
	svc, err := h.ensureUserService()
	if err != nil {
		writeError(r, http.StatusInternalServerError, 50001, "认证服务不可用")
		return
	}

	token := r.Cookie.Get("refresh_token").String()
	if token == "" {
		writeError(r, http.StatusUnauthorized, 40002, "未找到刷新令牌")
		return
	}

	remember := r.Cookie.Get("remember_me").Bool()

	result, err := svc.Refresh(r.GetCtx(), token)
	if err != nil {
		writeAuthError(r, err)
		return
	}

	h.setRefreshCookie(r, svc, result.Tokens.RefreshToken, remember)
	writeSuccess(r, "刷新成功", authResponse{
		AccessToken: result.Tokens.AccessToken,
		ExpiresIn:   result.Tokens.AccessExpiresIn,
	})
}

// GetProfile 返回当前登录用户信息。
func (h *Handler) GetProfile(r *ghttp.Request) {
	svc, err := h.ensureUserService()
	if err != nil {
		writeError(r, http.StatusInternalServerError, 50001, "认证服务不可用")
		return
	}

	userID := r.GetCtxVar("user_id").Int64()
	if userID <= 0 {
		writeError(r, http.StatusUnauthorized, 40001, "认证信息缺失")
		return
	}

	info, err := svc.Profile(r.GetCtx(), userID)
	if err != nil {
		writeAuthError(r, err)
		return
	}

	writeSuccess(r, "success", profileResponse{
		ID:                 info.ID,
		Username:           info.Username,
		IsAdmin:            info.IsAdmin,
		MustChangePassword: info.MustChangePassword,
	})
}

func parseAuthRequest(r *ghttp.Request) (*authRequest, bool) {
	var req authRequest
	if err := r.Parse(&req); err != nil || req.Username == "" || req.Password == "" {
		writeError(r, http.StatusBadRequest, 40004, "请求参数无效")
		return nil, false
	}
	return &req, true
}

func (req *authRequest) isRemember() bool {
	return req.RememberMe || req.Remember
}

func writeSuccess(r *ghttp.Request, message string, data interface{}) {
	r.Response.WriteJson(Response{
		Code:    20000,
		Message: message,
		Data:    data,
	})
}

func writeError(r *ghttp.Request, status int, code int, message string) {
	r.Response.WriteStatus(status)
	r.Response.ClearBuffer()
	r.Response.WriteJson(Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

func writeAuthError(r *ghttp.Request, err error) {
	switch err {
	case user.ErrInvalidInput:
		writeError(r, http.StatusBadRequest, 40004, "请求参数无效")
	case user.ErrUserExists:
		writeError(r, http.StatusConflict, 40009, "用户名已存在")
	case user.ErrPermissionDenied:
		writeError(r, http.StatusForbidden, 40010, "需要管理员权限")
	case user.ErrMustChangePassword:
		writeError(r, http.StatusForbidden, 40011, "需要先修改密码")
	case user.ErrInvalidCredential:
		writeError(r, http.StatusUnauthorized, 40001, "用户名或密码错误")
	case user.ErrRefreshTokenInvalid, user.ErrRefreshTokenExpired:
		writeError(r, http.StatusUnauthorized, 40002, "刷新令牌无效或已过期")
	case user.ErrUserNotFound:
		// 为避免泄露用户名是否存在（信息泄露），对用户未找到的情况
		// 返回与凭证错误相同的响应：401 + 通用提示。
		writeError(r, http.StatusUnauthorized, 40001, "用户名或密码错误")
	default:
		g.Log().Error(r.GetCtx(), err)
		writeError(r, http.StatusInternalServerError, 50001, "服务器繁忙，请稍后再试")
	}
}

func (h *Handler) setRefreshCookie(r *ghttp.Request, svc *user.Service, token string, remember bool) {
	if token == "" {
		return
	}
	maxAge := time.Duration(0)
	if remember {
		maxAge = svc.RefreshTTL()
	}

	r.Cookie.SetCookie("refresh_token", token, "", "/", maxAge, ghttp.CookieOptions{
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   r.TLS != nil,
	})

	flag := "0"
	if remember {
		flag = "1"
	}
	r.Cookie.SetCookie("remember_me", flag, "", "/", maxAge, ghttp.CookieOptions{
		HttpOnly: false,
		SameSite: http.SameSiteLaxMode,
		Secure:   r.TLS != nil,
	})
}

func (h *Handler) clearRefreshCookie(r *ghttp.Request) {
	expireAge := -1 * time.Hour
	r.Cookie.SetCookie("refresh_token", "", "", "/", expireAge, ghttp.CookieOptions{
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   r.TLS != nil,
	})
	r.Cookie.SetCookie("remember_me", "", "", "/", expireAge, ghttp.CookieOptions{
		HttpOnly: false,
		SameSite: http.SameSiteLaxMode,
		Secure:   r.TLS != nil,
	})
}

func (h *Handler) ensureUserService() (*user.Service, error) {
	if h.userService != nil {
		return h.userService, nil
	}
	svc, err := internal.BuildUserService()
	if err != nil {
		return nil, err
	}
	h.userService = svc
	return svc, nil
}
