package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	// Blockchain configuration
	BlockchainRPCURL  string
	ContractAddress   string
	StartBlock        uint64
	ReconnectInterval int // seconds

	// AWS SNS configuration
	SNSEnabled  bool
	SNSTopicARN string
	AWSRegion   string

	// Logging configuration
	LogLevel  string
	LogFormat string // json or console
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Try to load .env file, but don't fail if it doesn't exist
	_ = godotenv.Load()

	cfg := &Config{
		BlockchainRPCURL:  getEnv("BLOCKCHAIN_RPC_URL", ""),
		ContractAddress:   getEnv("CONTRACT_ADDRESS", ""),
		StartBlock:        getEnvAsUint64("START_BLOCK", 0),
		ReconnectInterval: getEnvAsInt("RECONNECT_INTERVAL", 5),
		SNSEnabled:        getEnvAsBool("SNS_ENABLED", false),
		SNSTopicARN:       getEnv("SNS_TOPIC_ARN", ""),
		AWSRegion:         getEnv("AWS_REGION", "us-east-1"),
		LogLevel:          getEnv("LOG_LEVEL", "info"),
		LogFormat:         getEnv("LOG_FORMAT", "console"),
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.BlockchainRPCURL == "" {
		return fmt.Errorf("BLOCKCHAIN_RPC_URL is required")
	}

	if c.ContractAddress == "" {
		return fmt.Errorf("CONTRACT_ADDRESS is required")
	}

	// Validate SNS configuration if enabled
	if c.SNSEnabled && c.SNSTopicARN == "" {
		return fmt.Errorf("SNS_TOPIC_ARN is required when SNS_ENABLED is true")
	}

	// Validate log level
	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}

	if !validLogLevels[c.LogLevel] {
		return fmt.Errorf("invalid LOG_LEVEL: %s (valid values: debug, info, warn, error)", c.LogLevel)
	}

	// Validate log format
	validLogFormats := map[string]bool{
		"json":    true,
		"console": true,
	}

	if !validLogFormats[c.LogFormat] {
		return fmt.Errorf("invalid LOG_FORMAT: %s (valid values: json, console)", c.LogFormat)
	}

	return nil
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvAsInt retrieves an environment variable as int or returns a default value
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}

// getEnvAsUint64 retrieves an environment variable as uint64 or returns a default value
func getEnvAsUint64(key string, defaultValue uint64) uint64 {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.ParseUint(valueStr, 10, 64)
	if err != nil {
		return defaultValue
	}

	return value
}

// getEnvAsBool retrieves an environment variable as bool or returns a default value
func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}
