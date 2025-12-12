package middleware

import (
	"context"
	"fmt"
	"runtime/debug"

	"go-study2/internal/infrastructure/logger"

	"github.com/gogf/gf/v2/net/ghttp"
)

// PanicRecovery returns a middleware that recovers from panics and logs them
func PanicRecovery(r *ghttp.Request) {
	defer func() {
		if err := recover(); err != nil {
			// Get trace ID from context
			traceID := logger.ExtractTraceID(r.Context())
			if traceID == "" {
				traceID = "unknown"
			}

			// Get stack trace
			stackTrace := string(debug.Stack())

			// Log panic with stack trace
			logPanic(traceID, err, stackTrace, r)

			// Return 500 Internal Server Error
			r.Response.WriteStatusExit(500, "Internal Server Error")
		}
	}()

	// Continue processing
	r.Middleware.Next()
}

// logPanic logs panic information with stack trace
func logPanic(traceID string, err interface{}, stackTrace string, r *ghttp.Request) {
	errorLogger := logger.GetInstance("error")
	if errorLogger == nil {
		// Fallback to default logger if error logger not configured
		errorLogger = logger.GetInstance("app")
		if errorLogger == nil {
			return
		}
	}

	ctx := context.WithValue(context.Background(), logger.TraceIDKey{}, traceID)

	// Format panic log message
	logMsg := fmt.Sprintf("[PANIC] TraceID: %s, Method: %s, Path: %s, Error: %v\nStack Trace:\n%s",
		traceID,
		r.Method,
		r.URL.Path,
		err,
		stackTrace)

	errorLogger.Error(ctx, logMsg)
}
