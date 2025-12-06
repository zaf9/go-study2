package middleware

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// Logger 记录请求日志中间件
func Logger(r *ghttp.Request) {
	ctx := r.Context()
	startTime := time.Now()

	// 继续执行后续处理
	r.Middleware.Next()

	// 记录请求耗时和状态
	duration := time.Since(startTime)
	status := r.Response.Status

	g.Log().Infof(ctx,
		"Status: %d | Method: %s | URL: %s | ClientIP: %s | Duration: %v",
		status,
		r.Method,
		r.URL.String(),
		r.GetClientIp(),
		duration,
	)
}
