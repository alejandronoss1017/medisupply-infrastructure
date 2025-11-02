package logger

import (
	"os"
	"strings"

	"go.uber.org/zap/zapcore"
)

// Config holds logger configuration
type Config struct {
	// Level is the minimum enabled logging level
	Level zapcore.Level
	// Environment determines output format (development or production)
	Environment string
	// Encoding sets the logger's encoding (json or console)
	Encoding string
	// EnableCaller adds caller information (file:line)
	EnableCaller bool
	// EnableStacktrace enables stack trace on error and above
	EnableStacktrace bool
}

// DefaultConfig returns a default logger configuration
func DefaultConfig() *Config {
	return &Config{
		Level:            zapcore.InfoLevel,
		Environment:      "development",
		Encoding:         "console",
		EnableCaller:     true,
		EnableStacktrace: true,
	}
}

// LoadConfigFromEnv loads configuration from environment variables
func LoadConfigFromEnv() *Config {
	config := DefaultConfig()

	// LOG_LEVEL: debug, info, warn, error, fatal
	if levelStr := os.Getenv("LOG_LEVEL"); levelStr != "" {
		if level, err := zapcore.ParseLevel(strings.ToLower(levelStr)); err == nil {
			config.Level = level
		}
	}

	// LOG_ENV: development, production
	if env := os.Getenv("LOG_ENV"); env != "" {
		config.Environment = strings.ToLower(env)
	}

	// LOG_ENCODING: json, console
	if encoding := os.Getenv("LOG_ENCODING"); encoding != "" {
		config.Encoding = strings.ToLower(encoding)
	}

	// Auto-configure based on environment
	if config.Environment == "production" {
		// Production: JSON, less verbose
		if config.Encoding == "console" {
			config.Encoding = "json"
		}
		config.EnableStacktrace = false // Only on errors
	}

	return config
}
