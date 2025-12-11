package http

import (
	"net/http"

	"github.com/gogf/gf/v2/net/ghttp"
)

// writeSuccess 返回标准成功响应。
func writeSuccess(r *ghttp.Request, data interface{}) {
	r.Response.WriteHeader(http.StatusOK)
	r.Response.WriteJson(map[string]interface{}{
		"code":    0,
		"message": "success",
		"data":    data,
	})
}

// writeError 返回标准错误响应。
func writeError(r *ghttp.Request, status int, code int, message string) {
	r.Response.WriteHeader(status)
	r.Response.WriteJson(map[string]interface{}{
		"code":    code,
		"message": message,
	})
}
