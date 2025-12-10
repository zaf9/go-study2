package middleware

import (
	"github.com/gogf/gf/v2/net/ghttp"
)

// Format 处理响应格式中间件
// 验证 format 参数 (json/html)，默认为 json
// 将格式存储在上下文中供处理程序使用
func Format(r *ghttp.Request) {
	format := r.Get("format").String()

	if format == "" {
		format = "json"
	}

	if format != "json" && format != "html" {
		r.Response.WriteStatus(400)
		r.Response.WriteJson(map[string]interface{}{
			"code":    400,
			"message": "Invalid format parameter. Supported values: json, html",
		})
		r.Exit() // 停止后续处理
		return
	}

	// 将验证后的格式存储在上下文中，键名为 "format"
	r.SetCtxVar("format", format)

	r.Middleware.Next()
}
