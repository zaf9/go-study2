package logger

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// MaxQueryMatches limits the number of matched entries returned by query
// functions to avoid excessive memory usage. It can be tuned in tests or
// production initialization.
var MaxQueryMatches = 10000

// LogEntry represents a parsed log entry
type LogEntry struct {
	Timestamp time.Time
	Level     string
	TraceID   string
	Message   string
	Raw       string
}

// QueryResult contains the results of a log query
type QueryResult struct {
	Entries []LogEntry
	Total   int
	Matched int
}

// ReadLogFile reads and parses a log file
func ReadLogFile(filePath string) ([]LogEntry, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	var entries []LogEntry
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if entry := parseLogEntry(line); entry != nil {
			entries = append(entries, *entry)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading log file: %w", err)
	}

	return entries, nil
}

// QueryByTraceID searches for log entries by trace ID
func QueryByTraceID(logPath, traceID string) (*QueryResult, error) {
	file, err := os.Open(logPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	total := 0
	var matched []LogEntry
	for scanner.Scan() {
		total++
		line := scanner.Text()
		entry := parseLogEntry(line)
		if entry == nil {
			continue
		}
		if strings.Contains(entry.TraceID, traceID) {
			matched = append(matched, *entry)
			if MaxQueryMatches > 0 && len(matched) >= MaxQueryMatches {
				break
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading log file: %w", err)
	}
	return &QueryResult{Entries: matched, Total: total, Matched: len(matched)}, nil
}

// QueryByTimeRange searches for log entries within a time range
func QueryByTimeRange(logPath string, start, end time.Time) (*QueryResult, error) {
	file, err := os.Open(logPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	total := 0
	var matched []LogEntry
	for scanner.Scan() {
		total++
		line := scanner.Text()
		entry := parseLogEntry(line)
		if entry == nil {
			continue
		}
		if (entry.Timestamp.Equal(start) || entry.Timestamp.After(start)) && entry.Timestamp.Before(end) {
			matched = append(matched, *entry)
			if MaxQueryMatches > 0 && len(matched) >= MaxQueryMatches {
				break
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading log file: %w", err)
	}
	return &QueryResult{Entries: matched, Total: total, Matched: len(matched)}, nil
}

// QueryByLevel searches for log entries by log level
func QueryByLevel(logPath, level string) (*QueryResult, error) {
	file, err := os.Open(logPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	total := 0
	var matched []LogEntry
	for scanner.Scan() {
		total++
		line := scanner.Text()
		entry := parseLogEntry(line)
		if entry == nil {
			continue
		}
		if strings.EqualFold(entry.Level, level) {
			matched = append(matched, *entry)
			if MaxQueryMatches > 0 && len(matched) >= MaxQueryMatches {
				break
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading log file: %w", err)
	}
	return &QueryResult{Entries: matched, Total: total, Matched: len(matched)}, nil
}

// QueryByKeyword searches for log entries containing a keyword
func QueryByKeyword(logPath, keyword string) (*QueryResult, error) {
	file, err := os.Open(logPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	total := 0
	var matched []LogEntry
	for scanner.Scan() {
		total++
		line := scanner.Text()
		entry := parseLogEntry(line)
		if entry == nil {
			continue
		}
		if strings.Contains(entry.Message, keyword) || strings.Contains(entry.Raw, keyword) {
			matched = append(matched, *entry)
			if MaxQueryMatches > 0 && len(matched) >= MaxQueryMatches {
				break
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading log file: %w", err)
	}
	return &QueryResult{Entries: matched, Total: total, Matched: len(matched)}, nil
}

// parseLogEntry parses a single log line into a LogEntry
func parseLogEntry(line string) *LogEntry {
	// Simple parsing logic - can be enhanced based on actual log format
	// Expected format: [timestamp] [level] [TraceID:xxx] message

	entry := &LogEntry{Raw: line}

	// Extract trace ID
	traceIDRegex := regexp.MustCompile(`\[TraceID:([^\]]+)\]`)
	if matches := traceIDRegex.FindStringSubmatch(line); len(matches) > 1 {
		entry.TraceID = matches[1]
	}

	// Extract level (INFO, ERROR, etc.)
	levelRegex := regexp.MustCompile(`\b(INFO|ERROR|WARN|WARNING|DEBUG|FATAL|CRITICAL)\b`)
	if matches := levelRegex.FindStringSubmatch(line); len(matches) > 1 {
		entry.Level = matches[1]
	}

	// Extract timestamp (attempt to parse first bracketed timestamp using
	// the same format produced by formatAccessLog: 02/Jan/2006:15:04:05 -0700)
	// e.g. [13/Dec/2025:10:15:30 +0800]
	// Attempt parsing with configured time format first (if available).
	timestampRegex := regexp.MustCompile(`\[([^\]]+)\]`)
	if tsMatches := timestampRegex.FindStringSubmatch(line); len(tsMatches) > 1 {
		ts := tsMatches[1]
		// Try configured format
		if globalConfig != nil && globalConfig.TimeFormat != "" {
			if t, err := time.Parse(globalConfig.TimeFormat, ts); err == nil {
				entry.Timestamp = t
				goto PARSED
			}
		}
		// Try common access log format
		if t, err := time.Parse("02/Jan/2006:15:04:05 -0700", ts); err == nil {
			entry.Timestamp = t
			goto PARSED
		}
		// Fallback to RFC3339
		if t2, err2 := time.Parse(time.RFC3339, ts); err2 == nil {
			entry.Timestamp = t2
			goto PARSED
		}
	}
	// Last resort: now
	entry.Timestamp = time.Now()
PARSED:

	// Extract message (everything after TraceID)
	if strings.Contains(line, "] ") {
		parts := strings.SplitN(line, "] ", 2)
		if len(parts) > 1 {
			entry.Message = parts[1]
		}
	}

	return entry
}

// ListLogFiles returns all log files in the logs directory
func ListLogFiles(logsDir string) ([]string, error) {
	var files []string

	err := filepath.Walk(logsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".log") {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}
