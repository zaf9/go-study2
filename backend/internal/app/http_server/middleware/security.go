package middleware

import (
	"github.com/gogf/gf/v2/net/ghttp"
)

// Security 添加安全相关的 HTTP 头，防御常见 Web 攻击
//
// 安全头说明：
//   - HSTS (Strict-Transport-Security): 强制浏览器使用 HTTPS，防止降级攻击
//   - CSP (Content-Security-Policy): 白名单机制控制资源加载，防止 XSS 攻击
//   - X-Frame-Options: 防止点击劫持攻击（页面被嵌入 iframe）
//   - X-Content-Type-Options: 防止 MIME 嗅探攻击
//   - X-XSS-Protection: 启用浏览器 XSS 过滤器
//   - Referrer-Policy: 控制 Referer 信息泄露
//   - Permissions-Policy: 限制浏览器功能访问（摄像头、麦克风等）
func Security(r *ghttp.Request) {
	// 1. HSTS: 强制浏览器使用 HTTPS（仅在生产环境且启用 TLS 时）
	// max-age=31536000: 浏览器记住 1 年
	// includeSubDomains: 包含所有子域名
	// preload: 允许加入 HSTS Preload 列表（浏览器内置强制 HTTPS）
	if r.TLS != nil {
		r.Response.Header().Set("Strict-Transport-Security",
			"max-age=31536000; includeSubDomains; preload")
	}

	// 2. CSP: 内容安全策略（开发环境使用宽松模式）
	// default-src 'self': 默认只允许加载同源资源
	// script-src: 允许同源脚本 + 内联脚本（开发环境需要）
	// style-src: 允许同源样式 + 内联样式
	// img-src: 允许图片、data: URL 和 HTTPS 图片
	// connect-src: 限制 API 请求目标为同源
	// font-src: 只允许同源字体
	// frame-ancestors 'none': 禁止页面被嵌入 iframe
	cspDirectives := []string{
		"default-src 'self'",
		"script-src 'self' 'unsafe-inline' 'unsafe-eval'", // 开发环境需要内联脚本
		"style-src 'self' 'unsafe-inline'",                // 开发环境需要内联样式
		"img-src 'self' data: https:",
		"connect-src 'self'",
		"font-src 'self'",
		"object-src 'none'",  // 禁止插件（Flash、Java 等）
		"frame-ancestors 'none'", // 禁止被嵌入 iframe
	}
	r.Response.Header().Set("Content-Security-Policy", joinDirectives(cspDirectives))

	// 3. 防止点击劫持（禁止页面被嵌入 iframe）
	// DENY: 完全禁止；SAMEORIGIN: 允许同源嵌入
	r.Response.Header().Set("X-Frame-Options", "DENY")

	// 4. 防止 MIME 嗅探攻击
	// 浏览器严格遵守声明的 Content-Type，不猜测文件类型
	r.Response.Header().Set("X-Content-Type-Options", "nosniff")

	// 5. XSS 防护（启用浏览器内置 XSS 过滤器）
	// 1: 启用；mode=block: 检测到 XSS 时阻止整个页面渲染
	r.Response.Header().Set("X-XSS-Protection", "1; mode=block")

	// 6. Referrer 策略（防止 URL 中的敏感信息泄露）
	// strict-origin-when-cross-origin: 同源完整发送，跨域只发送 origin
	r.Response.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

	// 7. 权限策略（限制浏览器功能访问）
	// 禁用地理位置、摄像头、麦克风等敏感功能（除非显式需要）
	r.Response.Header().Set("Permissions-Policy",
		"geolocation=(), microphone=(), camera=(), payment=()")

	// 继续执行后续中间件
	r.Middleware.Next()
}

// joinDirectives 将 CSP 指令数组拼接为字符串
func joinDirectives(directives []string) string {
	result := ""
	for i, d := range directives {
		if i > 0 {
			result += "; "
		}
		result += d
	}
	return result
}
