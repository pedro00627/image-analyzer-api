// Package config handles application configuration from environment variables
package config

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	return LoadConfigWithReader(envReader{})
}

// LoadConfigWithReader loads configuration using a custom Reader (for testing)
func LoadConfigWithReader(reader Reader) (*Config, error) {
	cfg := &Config{}

	if err := loadAIProviderConfig(cfg, reader); err != nil {
		return nil, err
	}

	if err := loadServerConfig(cfg, reader); err != nil {
		return nil, err
	}

	if err := loadFileUploadConfig(cfg, reader); err != nil {
		return nil, err
	}

	if err := loadTimeoutConfig(cfg, reader); err != nil {
		return nil, err
	}

	return cfg, nil
}

type Config struct {
	// AI Provider Configuration
	AIProvider string
	APIKey     string // Contains the API key(s) for the selected provider
	APISecret  string // Optional, used by some providers like Imagga

	// Server Configuration
	Port           int
	AllowedOrigins []string

	// File Upload Limits
	MaxFileSize      int64
	AllowedFileTypes []string

	// Timeouts
	AIServiceTimeout time.Duration
	HTTPReadTimeout  time.Duration
	HTTPWriteTimeout time.Duration
	HTTPIdleTimeout  time.Duration
}

// loadAIProviderConfig loads AI provider configuration
func loadAIProviderConfig(cfg *Config, reader Reader) error {
	cfg.AIProvider = reader.Get("AI_PROVIDER", "")
	if cfg.AIProvider == "" {
		return fmt.Errorf("AI_PROVIDER is required")
	}

	cfg.APIKey = reader.Get("API_KEY", "")
	cfg.APISecret = reader.Get("API_SECRET", "")

	return nil
}

// loadServerConfig loads server-related configuration
func loadServerConfig(cfg *Config, reader Reader) error {
	port, err := strconv.Atoi(reader.Get("PORT", "8080"))
	if err != nil {
		return fmt.Errorf("invalid PORT: %w", err)
	}
	cfg.Port = port

	cfg.AllowedOrigins = parseCommaSeparated(reader.Get("ALLOWED_ORIGINS", "http://localhost:4200,http://localhost:3000"))

	return nil
}

// loadFileUploadConfig loads file upload limits and allowed types
func loadFileUploadConfig(cfg *Config, reader Reader) error {
	maxSize, err := strconv.ParseInt(reader.Get("MAX_FILE_SIZE", "10485760"), 10, 64)
	if err != nil {
		return fmt.Errorf("invalid MAX_FILE_SIZE: %w", err)
	}
	cfg.MaxFileSize = maxSize

	cfg.AllowedFileTypes = parseCommaSeparated(reader.Get("ALLOWED_FILE_TYPES", "image/jpeg,image/png,image/gif,image/webp"))

	return nil
}

// loadTimeoutConfig loads all timeout configurations
func loadTimeoutConfig(cfg *Config, reader Reader) error {
	var err error

	cfg.AIServiceTimeout, err = parseDuration(reader, "AI_SERVICE_TIMEOUT", "30s")
	if err != nil {
		return err
	}

	cfg.HTTPReadTimeout, err = parseDuration(reader, "HTTP_READ_TIMEOUT", "10s")
	if err != nil {
		return err
	}

	cfg.HTTPWriteTimeout, err = parseDuration(reader, "HTTP_WRITE_TIMEOUT", "10s")
	if err != nil {
		return err
	}

	cfg.HTTPIdleTimeout, err = parseDuration(reader, "HTTP_IDLE_TIMEOUT", "60s")
	if err != nil {
		return err
	}

	return nil
}

// ValidationConfig returns a validator.ValidationConfig based on app config
// This is used when constructing the ImageValidator
func (c *Config) ValidationConfig() ValidationConfig {
	return ValidationConfig{
		AllowedMimes: c.AllowedFileTypes,
		MaxSize:      c.MaxFileSize,
	}
}

// ValidationConfig is a DTO for validator configuration
// It's defined here to avoid circular imports
type ValidationConfig struct {
	AllowedMimes []string
	MaxSize      int64
}

// parseDuration parses a duration from Reader or default value
func parseDuration(reader Reader, key, defaultValue string) (time.Duration, error) {
	value := reader.Get(key, defaultValue)
	duration, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %w", key, err)
	}
	return duration, nil
}

// parseCommaSeparated splits and trims comma-separated values
func parseCommaSeparated(value string) []string {
	parts := strings.Split(value, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}
