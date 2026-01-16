package middleware

import (
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// RateLimiterConfig 速率限制配置
type RateLimiterConfig struct {
	MaxAttempts    int           // 最大尝试次数
	WindowDuration time.Duration // 时间窗口
	BlockDuration  time.Duration // 封禁时长
}

// DefaultLoginRateLimitConfig 默认登录速率限制配置
var DefaultLoginRateLimitConfig = RateLimiterConfig{
	MaxAttempts:    5,               // 5 次
	WindowDuration: 1 * time.Minute, // 每分钟
	BlockDuration:  5 * time.Minute, // 封禁 5 分钟
}

// requestInfo 请求信息
type requestInfo struct {
	count      int       // 请求计数
	lastReset  time.Time // 上次重置时间
	blocked    bool      // 是否被封禁
	blockUntil time.Time // 封禁截止时间
}

// RateLimiter 速率限制器
type RateLimiter struct {
	mu     sync.RWMutex
	config RateLimiterConfig
	ips    map[string]*requestInfo // IP → 请求信息
}

// NewRateLimiter 创建速率限制器
func NewRateLimiter(config RateLimiterConfig) *RateLimiter {
	return &RateLimiter{
		config: config,
		ips:    make(map[string]*requestInfo),
	}
}

// GlobalLoginRateLimiter 全局登录速率限制器（单例）
var GlobalLoginRateLimiter = NewRateLimiter(DefaultLoginRateLimitConfig)

// LoginRateLimit 登录速率限制中间件
//
// 配置：每 IP 每分钟最多 5 次登录尝试，超过后封禁 5 分钟
// 只对 /api/v1/auth/login 路由生效
func LoginRateLimit(r *ghttp.Request) {
	// 只限制登录接口
	if r.URL.Path != "/api/v1/auth/login" || r.Method != "POST" {
		r.Middleware.Next()
		return
	}

	// 获取客户端真实 IP（考虑代理情况）
	ip := getClientIP(r)

	// 检查速率限制
	allowed, retryAfter := GlobalLoginRateLimiter.Allow(ip)
	if !allowed {
		// 设置状态码和响应头
		r.Response.WriteHeader(429)
		r.Response.Header().Set("Retry-After", formatRetryAfter(retryAfter))
		r.Response.WriteJson(g.Map{
			"code":    42901,
			"message": "登录尝试过多，请稍后再试",
			"data": g.Map{
				"retryAfter": int(retryAfter.Seconds()),
			},
		})
		r.Exit()
		return
	}

	r.Middleware.Next()
}

// Allow 检查是否允许请求
// 返回值：
//   - bool: true 允许，false 拒绝
//   - time.Duration: 重试等待时间
func (rl *RateLimiter) Allow(ip string) (bool, time.Duration) {
	now := time.Now()

	rl.mu.Lock()
	defer rl.mu.Unlock()

	info, exists := rl.ips[ip]

	// 情况 1: 首次请求或窗口过期
	if !exists || now.Sub(info.lastReset) > rl.config.WindowDuration {
		rl.ips[ip] = &requestInfo{
			count:     1,
			lastReset: now,
			blocked:   false,
		}
		return true, 0
	}

	// 情况 2: 正在封禁中
	if info.blocked {
		if now.Before(info.blockUntil) {
			// 仍在封禁期
			retryAfter := info.blockUntil.Sub(now)
			return false, retryAfter
		}
		// 封禁期已过，解封
		info.blocked = false
		info.count = 0
		info.lastReset = now
	}

	// 增加计数
	info.count++

	// 情况 3: 超过阈值，触发封禁
	if info.count > rl.config.MaxAttempts {
		info.blocked = true
		info.blockUntil = now.Add(rl.config.BlockDuration)
		return false, rl.config.BlockDuration
	}

	return true, 0
}

// isValidIP 验证 IP 地址格式
func isValidIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

// formatRetryAfter 格式化 Retry-After 头（秒）
func formatRetryAfter(d time.Duration) string {
	seconds := int(d.Seconds())
	if seconds < 1 {
		return "1"
	}
	// 使用 strconv.Itoa 将整数转换为字符串
	return strconv.Itoa(seconds)
}

// Reset 清理过期记录（可定期调用以释放内存）
func (rl *RateLimiter) Reset() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	for ip, info := range rl.ips {
		// 删除超过 10 分钟未活动的记录
		if now.Sub(info.lastReset) > 10*time.Minute && !info.blocked {
			delete(rl.ips, ip)
		}
		// 删除已过封禁期的记录
		if info.blocked && now.After(info.blockUntil) {
			delete(rl.ips, ip)
		}
	}
}

// Stats 获取统计信息
func (rl *RateLimiter) Stats() map[string]interface{} {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	now := time.Now()
	activeIPs := 0
	blockedIPs := 0

	for _, info := range rl.ips {
		if now.Sub(info.lastReset) < rl.config.WindowDuration {
			activeIPs++
		}
		if info.blocked && now.Before(info.blockUntil) {
			blockedIPs++
		}
	}

	return map[string]interface{}{
		"total_ips":       len(rl.ips),
		"active_ips":      activeIPs,
		"blocked_ips":     blockedIPs,
		"max_attempts":    rl.config.MaxAttempts,
		"window_duration": rl.config.WindowDuration.String(),
		"block_duration":  rl.config.BlockDuration.String(),
	}
}
