package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// LogLevel represents the logging level
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

// LogEntry represents a structured log entry
type LogEntry struct {
	Timestamp     time.Time `json:"timestamp"`
	Level         string    `json:"level"`
	Message       string    `json:"message"`
	Command       string    `json:"command,omitempty"`
	Error         string    `json:"error,omitempty"`
	Suggestion    string    `json:"suggestion,omitempty"`
	Explanation   string    `json:"explanation,omitempty"`
	UserApproved  bool      `json:"user_approved,omitempty"`
	Outcome       string    `json:"outcome,omitempty"`
	Plugin        string    `json:"plugin,omitempty"`
	AIProvider    string    `json:"ai_provider,omitempty"`
	ProcessingTime int64     `json:"processing_time_ms,omitempty"`
}

// Logger handles structured logging for LogAid
type Logger struct {
	level   LogLevel
	logFile *os.File
}

// New creates a new logger instance
func New(levelStr, logPath string) (*Logger, error) {
	level, err := parseLogLevel(levelStr)
	if err != nil {
		return nil, fmt.Errorf("invalid log level: %w", err)
	}

	// Ensure log directory exists
	if err := os.MkdirAll(filepath.Dir(logPath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	// Open log file
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	return &Logger{
		level:   level,
		logFile: logFile,
	}, nil
}

// Close closes the logger
func (l *Logger) Close() error {
	if l.logFile != nil {
		return l.logFile.Close()
	}
	return nil
}

// Debug logs a debug message
func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(DEBUG, fmt.Sprintf(format, args...))
}

// Info logs an info message
func (l *Logger) Info(format string, args ...interface{}) {
	l.log(INFO, fmt.Sprintf(format, args...))
}

// Warn logs a warning message
func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(WARN, fmt.Sprintf(format, args...))
}

// Error logs an error message
func (l *Logger) Error(format string, args ...interface{}) {
	l.log(ERROR, fmt.Sprintf(format, args...))
}

// LogCommand logs a command execution
func (l *Logger) LogCommand(command, output, stderr string, exitCode int) {
	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     "COMMAND",
		Command:   command,
		Message:   fmt.Sprintf("Command executed with exit code %d", exitCode),
	}

	if stderr != "" {
		entry.Error = stderr
	}

	l.writeEntry(entry)
}

// LogSuggestion logs an AI suggestion
func (l *Logger) LogSuggestion(command, error, suggestion, explanation, plugin, aiProvider string, userApproved bool, outcome string, processingTime time.Duration) {
	entry := LogEntry{
		Timestamp:      time.Now(),
		Level:          "SUGGESTION",
		Message:        "AI suggestion generated",
		Command:        command,
		Error:          error,
		Suggestion:     suggestion,
		Explanation:    explanation,
		UserApproved:   userApproved,
		Outcome:        outcome,
		Plugin:         plugin,
		AIProvider:     aiProvider,
		ProcessingTime: processingTime.Milliseconds(),
	}

	l.writeEntry(entry)
}

// log writes a log entry
func (l *Logger) log(level LogLevel, message string) {
	if level < l.level {
		return
	}

	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     levelToString(level),
		Message:   message,
	}

	l.writeEntry(entry)

	// Also write to stderr for immediate feedback
	if level >= WARN {
		fmt.Fprintf(os.Stderr, "[%s] %s: %s\n", entry.Timestamp.Format("15:04:05"), entry.Level, message)
	}
}

// writeEntry writes a log entry to the log file
func (l *Logger) writeEntry(entry LogEntry) {
	if l.logFile == nil {
		return
	}

	data, err := json.Marshal(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to marshal log entry: %v\n", err)
		return
	}

	data = append(data, '\n')
	if _, err := l.logFile.Write(data); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write log entry: %v\n", err)
	}
}

// parseLogLevel parses a log level string
func parseLogLevel(levelStr string) (LogLevel, error) {
	switch strings.ToUpper(levelStr) {
	case "DEBUG":
		return DEBUG, nil
	case "INFO":
		return INFO, nil
	case "WARN", "WARNING":
		return WARN, nil
	case "ERROR":
		return ERROR, nil
	default:
		return INFO, fmt.Errorf("unknown log level: %s", levelStr)
	}
}

// levelToString converts a log level to string
func levelToString(level LogLevel) string {
	switch level {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}
