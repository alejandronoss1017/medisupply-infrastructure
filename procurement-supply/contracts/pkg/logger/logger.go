package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New creates a new logger with the given name/component
// This returns a sugared logger for easier use
func New(name string) *zap.SugaredLogger {
	config := LoadConfigFromEnv()
	return NewWithConfig(name, config)
}

// NewWithConfig creates a logger with a specific configuration
func NewWithConfig(name string, config *Config) *zap.SugaredLogger {
	logger := NewZapLogger(config)
	// Add a component name as a field to all logs
	return logger.Named(name).Sugar()
}

// NewZapLogger creates the base zap.Logger with configuration
func NewZapLogger(config *Config) *zap.Logger {
	// Encoder configuration
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Development mode: colorized, human-readable
	if config.Environment == "development" {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006/01/02 - 15:04:05")
	}

	// Create encoder
	var encoder zapcore.Encoder
	if config.Encoding == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// Core combines encoder, output, and level
	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stdout),
		config.Level,
	)

	// Logger options
	options := []zap.Option{
		zap.AddCaller(),
		zap.AddCallerSkip(1), // Skip one level to show actual caller
	}

	if config.EnableStacktrace {
		options = append(options, zap.AddStacktrace(zapcore.ErrorLevel))
	}

	return zap.New(core, options...)
}

// NewProduction creates a production-ready logger
// Outputs JSON, info level, with caller information
func NewProduction(name string) *zap.SugaredLogger {
	config := &Config{
		Level:            zapcore.InfoLevel,
		Environment:      "production",
		Encoding:         "json",
		EnableCaller:     true,
		EnableStacktrace: false,
	}
	return NewWithConfig(name, config)
}

// NewDevelopment creates a development logger
// Outputs colorized console, debug level, with caller information
func NewDevelopment(name string) *zap.SugaredLogger {
	config := &Config{
		Level:            zapcore.DebugLevel,
		Environment:      "development",
		Encoding:         "console",
		EnableCaller:     true,
		EnableStacktrace: true,
	}
	return NewWithConfig(name, config)
}
