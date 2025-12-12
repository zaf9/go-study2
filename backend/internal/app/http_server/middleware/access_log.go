package middleware

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go-study2/internal/infrastructure/logger"

	"github.com/gogf/gf/v2/net/ghttp"
)

// AccessLog returns a middleware that logs HTTP access requests
func AccessLog(r *ghttp.Request) {
	start := time.Now()

	// Extract headers for trace ID
	headers := make(map[string]string)
	for key, values := range r.Header {
		if len(values) > 0 {
			headers[key] = values[0]
		}
	}

	// Inject trace ID into context
	ctx := logger.InjectTraceIDToContext(r.Context(), headers)
	r.SetCtx(ctx)

	// Get trace ID for logging
	traceID := logger.ExtractTraceID(ctx)

	// Get client IP
	clientIP := getClientIP(r)

	// Get user agent
	userAgent := r.Header.Get("User-Agent")
	if userAgent == "" {
		userAgent = "-"
	}

	// Get request method and path
	method := r.Method
	path := r.URL.Path

	// Log request start
	logRequestStart(traceID, method, path, clientIP, userAgent)

	// Process request
	r.Middleware.Next()

	// Calculate duration
	duration := time.Since(start)

	// Get status code from response
	statusCode := r.Response.Status

	// Log request end
	logRequestEnd(traceID, method, path, statusCode, duration, clientIP, userAgent)
}

// responseWriter wraps http.ResponseWriter to capture status code and response size
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(data []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(data)
	rw.size += size
	return size, err
}

// getClientIP extracts the real client IP from the request
func getClientIP(r *ghttp.Request) string {
	// Check X-Forwarded-For header first
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		// Take the first IP if multiple
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			ip := strings.TrimSpace(ips[0])
			if ip != "" {
				return ip
			}
		}
	}

	// Check X-Real-IP header
	xri := r.Header.Get("X-Real-IP")
	if xri != "" {
		return xri
	}

	// Fall back to remote address
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

// logRequestStart logs the start of an HTTP request
func logRequestStart(traceID, method, path, clientIP, userAgent string) {
	accessLogger := logger.GetInstance("access")
	if accessLogger == nil {
		// Fallback to default logger if access logger not configured
		accessLogger = logger.GetInstance("app")
		if accessLogger == nil {
			return
		}
	}

	ctx := context.WithValue(context.Background(), logger.TraceIDKey{}, traceID)
	// Include TraceID directly in the access log message to ensure downstream
	// consumers (and tests) can find it even if the underlying logger does not
	// automatically include context values in the formatted output.
	logMsg := fmt.Sprintf("[TraceID:%s] [START] %s %s from %s", traceID, method, path, clientIP)
	accessLogger.Info(ctx, logMsg)
}

// logRequestEnd logs the end of an HTTP request
func logRequestEnd(traceID, method, path string, status int, duration time.Duration, clientIP, userAgent string) {
	accessLogger := logger.GetInstance("access")
	if accessLogger == nil {
		// Fallback to default logger if access logger not configured
		accessLogger = logger.GetInstance("app")
		if accessLogger == nil {
			return
		}
	}

	ctx := context.WithValue(context.Background(), logger.TraceIDKey{}, traceID)
	// Prefix the access log entry with the trace ID so tests and log processors
	// can reliably extract it.
	logMsg := fmt.Sprintf("[TraceID:%s] %s", traceID, formatAccessLog(method, path, status, duration, clientIP, userAgent))

	// Log based on status code
	if status >= 400 {
		accessLogger.Error(ctx, logMsg)
	} else {
		accessLogger.Info(ctx, logMsg)
	}
}

// formatAccessLog formats the access log entry
func formatAccessLog(method, path string, status int, duration time.Duration, clientIP, userAgent string) string {
	// Common Log Format-like: %h %l %u %t "%r" %>s %b
	// clientIP - - [timestamp] "method path" status "userAgent" duration

	timestamp := time.Now().Format(logger.TimeFormat())
	durationMs := duration.Milliseconds()

	return strings.Join([]string{
		clientIP,
		"-", // remote logname (not implemented)
		"-", // remote user (not implemented)
		"[" + timestamp + "]",
		"\"" + method + " " + path + "\"",
		strconv.Itoa(status),
		"\"" + userAgent + "\"",
		fmt.Sprintf("%dms", durationMs),
	}, " ")
}
