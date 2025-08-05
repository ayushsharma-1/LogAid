package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

var (
	InfoColor    = color.New(color.FgCyan)
	WarnColor    = color.New(color.FgYellow)
	ErrorColor   = color.New(color.FgRed)
	SuccessColor = color.New(color.FgGreen)
	DebugColor   = color.New(color.FgMagenta)
)

type Logger struct {
	level    string
	file     *os.File
	logger   *log.Logger
	colorful bool
}

var AppLogger *Logger

// Init initializes the logger
func Init() error {
	level := os.Getenv("LOG_LEVEL")
	if level == "" {
		level = "info"
	}

	logFile := os.Getenv("LOG_FILE")
	if logFile == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			logFile = ".logaid/logs/logaid.log"
		} else {
			logFile = filepath.Join(homeDir, ".logaid", "logs", "logaid.log")
		}
	}

	// Create log directory if it doesn't exist
	logDir := filepath.Dir(logFile)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	// Open log file
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	AppLogger = &Logger{
		level:    strings.ToLower(level),
		file:     file,
		logger:   log.New(file, "", log.LstdFlags),
		colorful: os.Getenv("ENABLE_COLORS") != "false",
	}

	return nil
}

// Close closes the logger
func (l *Logger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

// Debug logs a debug message
func (l *Logger) Debug(msg string) {
	if l.shouldLog("debug") {
		l.logger.Printf("[DEBUG] %s", msg)
		if l.colorful {
			DebugColor.Printf("[DEBUG] %s\n", msg)
		} else {
			fmt.Printf("[DEBUG] %s\n", msg)
		}
	}
}

// Info logs an info message
func (l *Logger) Info(msg string) {
	if l.shouldLog("info") {
		l.logger.Printf("[INFO] %s", msg)
		if l.colorful {
			InfoColor.Printf("[INFO] %s\n", msg)
		} else {
			fmt.Printf("[INFO] %s\n", msg)
		}
	}
}

// Warn logs a warning message
func (l *Logger) Warn(msg string) {
	if l.shouldLog("warn") {
		l.logger.Printf("[WARN] %s", msg)
		if l.colorful {
			WarnColor.Printf("[WARN] %s\n", msg)
		} else {
			fmt.Printf("[WARN] %s\n", msg)
		}
	}
}

// Error logs an error message
func (l *Logger) Error(msg string) {
	if l.shouldLog("error") {
		l.logger.Printf("[ERROR] %s", msg)
		if l.colorful {
			ErrorColor.Printf("[ERROR] %s\n", msg)
		} else {
			fmt.Printf("[ERROR] %s\n", msg)
		}
	}
}

// Success logs a success message
func (l *Logger) Success(msg string) {
	l.logger.Printf("[SUCCESS] %s", msg)
	if l.colorful {
		SuccessColor.Printf("✓ %s\n", msg)
	} else {
		fmt.Printf("✓ %s\n", msg)
	}
}

func (l *Logger) shouldLog(level string) bool {
	levels := map[string]int{
		"debug": 0,
		"info":  1,
		"warn":  2,
		"error": 3,
	}

	currentLevel, exists := levels[l.level]
	if !exists {
		currentLevel = 1 // default to info
	}

	msgLevel, exists := levels[level]
	if !exists {
		return false
	}

	return msgLevel >= currentLevel
}

// Global logging functions for convenience
func Debug(msg string) {
	if AppLogger != nil {
		AppLogger.Debug(msg)
	}
}

func Info(msg string) {
	if AppLogger != nil {
		AppLogger.Info(msg)
	}
}

func Warn(msg string) {
	if AppLogger != nil {
		AppLogger.Warn(msg)
	}
}

func Error(msg string) {
	if AppLogger != nil {
		AppLogger.Error(msg)
	}
}

func Success(msg string) {
	if AppLogger != nil {
		AppLogger.Success(msg)
	}
}
