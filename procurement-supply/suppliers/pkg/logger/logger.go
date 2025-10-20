package logger

import (
	"fmt"
	"os"
	"time"
)

// Logger represents a contextual logger with a specific prefix
type Logger struct {
	prefix string
}

// New creates a new logger with the specified prefix (e.g., "KAFKA", "HTTP", "APP")
func New(prefix string) *Logger {
	return &Logger{
		prefix: prefix,
	}
}

// formatMessage formats a log message in Gin style
func (l *Logger) formatMessage(level Level, format string, args ...interface{}) string {
	timestamp := time.Now().Format("2006/01/02 - 15:04:05")
	message := fmt.Sprintf(format, args...)

	if level != "" {
		return fmt.Sprintf("[%s] %s | %s | %s", l.prefix, timestamp, level, message)
	}
	return fmt.Sprintf("[%s] %s | %s", l.prefix, timestamp, message)
}

// Debug logs a debug level message
func (l *Logger) Debug(format string, args ...interface{}) {
	fmt.Println(l.formatMessage(Debug, format, args...))
}

// Info logs an info level message
func (l *Logger) Info(format string, args ...interface{}) {
	fmt.Println(l.formatMessage(Info, format, args...))
}

// Warn logs a warning level message
func (l *Logger) Warn(format string, args ...interface{}) {
	fmt.Println(l.formatMessage(Warn, format, args...))
}

// Error logs an error level message
func (l *Logger) Error(format string, args ...interface{}) {
	fmt.Println(l.formatMessage(Error, format, args...))
}

// Fatal logs a fatal level message and exits the program with status code 1
func (l *Logger) Fatal(format string, args ...interface{}) {
	fmt.Println(l.formatMessage(Fatal, format, args...))
	os.Exit(1)
}

// Log logs a message without a level (for simple logging)
func (l *Logger) Log(format string, args ...interface{}) {
	fmt.Println(l.formatMessage("", format, args...))
}
