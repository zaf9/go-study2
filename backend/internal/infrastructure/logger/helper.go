package logger

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

// LogInfo 记录信息级别日志，自动提取TraceID
func LogInfo(ctx context.Context, message string, args ...interface{}) {
	logger := GetInstance("app")
	if logger == nil {
		return
	}

	// 格式化消息
	formattedMsg := formatLogMessage(ctx, "INFO", message, args...)
	logger.Info(ctx, formattedMsg)
}

// LogError 记录错误级别日志，自动提取TraceID和堆栈信息
func LogError(ctx context.Context, err error, message string, args ...interface{}) {
	logger := GetInstance("error")
	if logger == nil {
		logger = GetInstance("app")
		if logger == nil {
			return
		}
	}

	// 获取堆栈信息
	_, file, line, ok := runtime.Caller(1)
	stackInfo := ""
	if ok {
		stackInfo = fmt.Sprintf(" [%s:%d]", file, line)
	}

	// 格式化消息
	formattedMsg := formatLogMessage(ctx, "ERROR", message+stackInfo, args...)
	if err != nil {
		formattedMsg += fmt.Sprintf(" Error: %v", err)
	}

	logger.Error(ctx, formattedMsg)
}

// LogSlow 记录慢操作日志，包含执行时间
func LogSlow(ctx context.Context, operation string, duration time.Duration, threshold time.Duration) {
	logger := GetInstance("slow")
	if logger == nil {
		logger = GetInstance("app")
		if logger == nil {
			return
		}
	}

	// 只有超过阈值才记录
	if duration < threshold {
		return
	}

	formattedMsg := formatLogMessage(ctx, "SLOW", fmt.Sprintf("Slow operation: %s took %v (threshold: %v)",
		operation, duration, threshold))

	logger.Warning(ctx, formattedMsg)
}

// LogBiz 记录业务操作日志，包含操作详情
func LogBiz(ctx context.Context, operation string, params map[string]interface{}, result interface{}, duration time.Duration) {
	logger := GetInstance("biz")
	if logger == nil {
		logger = GetInstance("app")
		if logger == nil {
			return
		}
	}

	// 格式化业务日志
	formattedMsg := formatLogMessage(ctx, "BIZ", fmt.Sprintf("Business operation: %s, Duration: %v",
		operation, duration))

	// 添加参数信息
	if len(params) > 0 {
		formattedMsg += fmt.Sprintf(", Params: %+v", params)
	}

	// 添加结果信息
	if result != nil {
		formattedMsg += fmt.Sprintf(", Result: %+v", result)
	}

	logger.Info(ctx, formattedMsg)
}

// formatLogMessage 格式化日志消息，自动添加TraceID前缀
func formatLogMessage(ctx context.Context, level, message string, args ...interface{}) string {
	traceID := ExtractTraceID(ctx)
	if traceID == "" {
		traceID = "no-trace"
	}

	// 格式化消息
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}

	// 添加TraceID前缀
	return fmt.Sprintf("[TraceID:%s] %s", traceID, message)
}

// LogWithFields 记录带字段的结构化日志
func LogWithFields(ctx context.Context, level string, message string, fields map[string]interface{}) {
	logger := GetInstance("app")
	if logger == nil {
		// 如果获取不到 logger，使用默认 logger
		logger = g.Log()
	}
	if logger == nil {
		// 如果默认 logger 也为 nil，直接返回
		return
	}

	// 构建带字段的消息
	fieldStr := ""
	for k, v := range fields {
		fieldStr += fmt.Sprintf(" %s=%v", k, v)
	}

	formattedMsg := formatLogMessage(ctx, level, message+fieldStr)

	// 根据级别记录
	switch level {
	case "DEBUG":
		logger.Debug(ctx, formattedMsg)
	case "INFO":
		logger.Info(ctx, formattedMsg)
	case "WARNING", "WARN":
		logger.Warning(ctx, formattedMsg)
	case "ERROR":
		logger.Error(ctx, formattedMsg)
	case "CRITICAL", "FATAL":
		logger.Critical(ctx, formattedMsg)
	default:
		logger.Info(ctx, formattedMsg)
	}
}
