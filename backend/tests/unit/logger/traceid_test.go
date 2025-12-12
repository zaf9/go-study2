package logger_test

import (
	"context"
	"testing"

	"go-study2/internal/infrastructure/logger"
)

func TestGenerateTraceID(t *testing.T) {
	id1 := logger.GenerateTraceID()
	id2 := logger.GenerateTraceID()

	// Check that IDs are generated
	if id1 == "" {
		t.Error("GenerateTraceID returned empty string")
	}
	if id2 == "" {
		t.Error("GenerateTraceID returned empty string")
	}

	// Check that IDs are different
	if id1 == id2 {
		t.Error("GenerateTraceID should generate unique IDs")
	}

	// Check that ID is reasonable length (hex encoded 16 bytes = 32 chars)
	if len(id1) != 32 {
		t.Errorf("Expected trace ID length 32, got %d", len(id1))
	}
}

func TestExtractTraceID(t *testing.T) {
	tests := []struct {
		name     string
		ctx      context.Context
		expected string
	}{
		{
			name:     "nil context",
			ctx:      nil,
			expected: "",
		},
		{
			name:     "empty context",
			ctx:      context.Background(),
			expected: "",
		},
		{
			name:     "context with trace ID",
			ctx:      context.WithValue(context.Background(), logger.TraceIDKey{}, "test-trace-id"),
			expected: "test-trace-id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := logger.ExtractTraceID(tt.ctx)
			if result != tt.expected {
				t.Errorf("ExtractTraceID() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestEnsureTraceID(t *testing.T) {
	// Test with empty context
	ctx1 := logger.EnsureTraceID(context.Background())
	traceID1 := logger.ExtractTraceID(ctx1)
	if traceID1 == "" {
		t.Error("EnsureTraceID should generate trace ID for empty context")
	}

	// Test with existing trace ID
	ctx2 := context.WithValue(context.Background(), logger.TraceIDKey{}, "existing-trace-id")
	ctx3 := logger.EnsureTraceID(ctx2)
	traceID2 := logger.ExtractTraceID(ctx3)
	if traceID2 != "existing-trace-id" {
		t.Errorf("EnsureTraceID should preserve existing trace ID, got %s", traceID2)
	}
}

func TestGetTraceIDFromHeaders(t *testing.T) {
	tests := []struct {
		name     string
		headers  map[string]string
		expected string
	}{
		{
			name:     "empty headers",
			headers:  map[string]string{},
			expected: "",
		},
		{
			name: "X-Trace-Id header",
			headers: map[string]string{
				"X-Trace-Id": "trace-from-x-trace-id",
			},
			expected: "trace-from-x-trace-id",
		},
		{
			name: "X-Request-Id header",
			headers: map[string]string{
				"X-Request-Id": "trace-from-x-request-id",
			},
			expected: "trace-from-x-request-id",
		},
		{
			name: "Trace-Id header",
			headers: map[string]string{
				"Trace-Id": "trace-from-trace-id",
			},
			expected: "trace-from-trace-id",
		},
		{
			name: "Request-Id header",
			headers: map[string]string{
				"Request-Id": "trace-from-request-id",
			},
			expected: "trace-from-request-id",
		},
		{
			name: "multiple headers - first one wins",
			headers: map[string]string{
				"X-Trace-Id":   "trace-x-trace-id",
				"X-Request-Id": "trace-x-request-id",
			},
			expected: "trace-x-trace-id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := logger.GetTraceIDFromHeaders(tt.headers)
			if result != tt.expected {
				t.Errorf("GetTraceIDFromHeaders() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestInjectTraceIDToContext(t *testing.T) {
	tests := []struct {
		name        string
		ctx         context.Context
		headers     map[string]string
		expectTrace bool
	}{
		{
			name:        "nil context with headers",
			ctx:         nil,
			headers:     map[string]string{"X-Trace-Id": "header-trace-id"},
			expectTrace: true,
		},
		{
			name:        "empty context with headers",
			ctx:         context.Background(),
			headers:     map[string]string{"X-Trace-Id": "header-trace-id"},
			expectTrace: true,
		},
		{
			name:        "context with existing trace ID",
			ctx:         context.WithValue(context.Background(), logger.TraceIDKey{}, "existing-trace-id"),
			headers:     map[string]string{"X-Trace-Id": "header-trace-id"},
			expectTrace: true, // existing should be preserved
		},
		{
			name:        "empty context without headers",
			ctx:         context.Background(),
			headers:     map[string]string{},
			expectTrace: true, // should generate new one
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultCtx := logger.InjectTraceIDToContext(tt.ctx, tt.headers)
			traceID := logger.ExtractTraceID(resultCtx)

			if tt.expectTrace && traceID == "" {
				t.Error("Expected trace ID to be present")
			}

			if !tt.expectTrace && traceID != "" {
				t.Error("Expected no trace ID")
			}

			// Check header extraction only when context doesn't have existing trace ID
			if headerTraceID := logger.GetTraceIDFromHeaders(tt.headers); headerTraceID != "" {
				existingTraceID := logger.ExtractTraceID(tt.ctx)
				if existingTraceID == "" && traceID != headerTraceID {
					t.Errorf("Expected trace ID from headers %s, got %s", headerTraceID, traceID)
				}
				// If context already has trace ID, it should be preserved
				if existingTraceID != "" && traceID != existingTraceID {
					t.Errorf("Expected existing trace ID %s to be preserved, got %s", existingTraceID, traceID)
				}
			}
		})
	}
}
