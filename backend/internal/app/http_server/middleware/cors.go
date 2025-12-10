package middleware

import "github.com/gogf/gf/v2/net/ghttp"

// Cors 处理跨域请求，仅用于开发环境。
func Cors(r *ghttp.Request) {
	r.Response.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	r.Response.Header().Set("Access-Control-Allow-Credentials", "true")
	r.Response.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	r.Response.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

	if r.Method == "OPTIONS" {
		r.Response.WriteStatus(204)
		r.ExitAll()
		return
	}

	r.Middleware.Next()
}
