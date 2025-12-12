package logger

import (
	"context"
	"crypto/rand"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/os/gtime"
)

// TraceIDKey is the context key for trace ID
type TraceIDKey struct{}

// GenerateTraceID generates a new trace ID
func GenerateTraceID() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to timestamp-based ID if crypto rand fails
		return fmt.Sprintf("trace_%d", gtime.TimestampNano())
	}

	// Format as hex string
	return fmt.Sprintf("%x", bytes)
}

// ExtractTraceID extracts trace ID from context
func ExtractTraceID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	if traceID, ok := ctx.Value(TraceIDKey{}).(string); ok {
		return traceID
	}

	return ""
}

// EnsureTraceID ensures a trace ID exists in the context, generating one if needed
func EnsureTraceID(ctx context.Context) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	// Check if trace ID already exists
	if ExtractTraceID(ctx) != "" {
		return ctx
	}

	// Generate and inject new trace ID
	traceID := GenerateTraceID()
	return context.WithValue(ctx, TraceIDKey{}, traceID)
}

// GetTraceIDFromHeaders extracts trace ID from HTTP headers
func GetTraceIDFromHeaders(headers map[string]string) string {
	// Check common trace ID header names
	traceHeaders := []string{
		"X-Trace-Id",
		"X-Request-Id",
		"Trace-Id",
		"Request-Id",
	}

	for _, header := range traceHeaders {
		if traceID := headers[header]; traceID != "" {
			return strings.TrimSpace(traceID)
		}
	}

	return ""
}

// InjectTraceIDToContext injects trace ID from headers into context
func InjectTraceIDToContext(ctx context.Context, headers map[string]string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	// Check if trace ID already exists in context
	if ExtractTraceID(ctx) != "" {
		return ctx
	}

	// Try to get trace ID from headers
	if traceID := GetTraceIDFromHeaders(headers); traceID != "" {
		return context.WithValue(ctx, TraceIDKey{}, traceID)
	}

	// Generate new trace ID if none found
	return EnsureTraceID(ctx)
}
