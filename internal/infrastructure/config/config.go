// Package config handles application configuration from environment variables
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds all application configuration
type Config struct {
	// AI Provider Configuration
	AIProvider         string
	GoogleVisionAPIKey string
	ImaggaAPIKey       string
	ImaggaAPISecret    string
	OpenAIAPIKey       string

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

// LoadConfig loads configuration from environment variables
// Returns error if required variables are missing or invalid
func LoadConfig() (*Config, error) {
	cfg := &Config{}

	// AI Provider (required)
	cfg.AIProvider = getEnv("AI_PROVIDER", "")
	if cfg.AIProvider == "" {
		return nil, fmt.Errorf("AI_PROVIDER is required")
	}

	// API Keys (at least one required based on provider)
	cfg.GoogleVisionAPIKey = getEnv("GOOGLE_VISION_API_KEY", "")
	cfg.ImaggaAPIKey = getEnv("IMAGGA_API_KEY", "")
	cfg.ImaggaAPISecret = getEnv("IMAGGA_API_SECRET", "")
	cfg.OpenAIAPIKey = getEnv("OPENAI_API_KEY", "")

	// Validate API key based on provider
	switch cfg.AIProvider {
	case "google_vision":
		if cfg.GoogleVisionAPIKey == "" {
			return nil, fmt.Errorf("GOOGLE_VISION_API_KEY is required for google_vision provider")
		}
	case "imagga":
		if cfg.ImaggaAPIKey == "" || cfg.ImaggaAPISecret == "" {
			return nil, fmt.Errorf("IMAGGA_API_KEY and IMAGGA_API_SECRET are required for imagga provider")
		}
	case "openai":
		if cfg.OpenAIAPIKey == "" {
			return nil, fmt.Errorf("OPENAI_API_KEY is required for openai provider")
		}
	default:
		return nil, fmt.Errorf("invalid AI_PROVIDER: %s (must be google_vision, imagga, or openai)", cfg.AIProvider)
	}

	// Server Port (default: 8080)
	port, err := strconv.Atoi(getEnv("PORT", "8080"))
	if err != nil {
		return nil, fmt.Errorf("invalid PORT: %w", err)
	}
	cfg.Port = port

	// Allowed Origins (default: localhost)
	originsStr := getEnv("ALLOWED_ORIGINS", "http://localhost:4200,http://localhost:3000")
	cfg.AllowedOrigins = strings.Split(originsStr, ",")
	for i := range cfg.AllowedOrigins {
		cfg.AllowedOrigins[i] = strings.TrimSpace(cfg.AllowedOrigins[i])
	}

	// Max File Size (default: 10MB)
	maxSize, err := strconv.ParseInt(getEnv("MAX_FILE_SIZE", "10485760"), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid MAX_FILE_SIZE: %w", err)
	}
	cfg.MaxFileSize = maxSize

	// Allowed File Types (default: common image types)
	typesStr := getEnv("ALLOWED_FILE_TYPES", "image/jpeg,image/png,image/gif,image/webp")
	cfg.AllowedFileTypes = strings.Split(typesStr, ",")
	for i := range cfg.AllowedFileTypes {
		cfg.AllowedFileTypes[i] = strings.TrimSpace(cfg.AllowedFileTypes[i])
	}

	// Timeouts
	cfg.AIServiceTimeout, err = time.ParseDuration(getEnv("AI_SERVICE_TIMEOUT", "30s"))
	if err != nil {
		return nil, fmt.Errorf("invalid AI_SERVICE_TIMEOUT: %w", err)
	}

	cfg.HTTPReadTimeout, err = time.ParseDuration(getEnv("HTTP_READ_TIMEOUT", "10s"))
	if err != nil {
		return nil, fmt.Errorf("invalid HTTP_READ_TIMEOUT: %w", err)
	}

	cfg.HTTPWriteTimeout, err = time.ParseDuration(getEnv("HTTP_WRITE_TIMEOUT", "10s"))
	if err != nil {
		return nil, fmt.Errorf("invalid HTTP_WRITE_TIMEOUT: %w", err)
	}

	cfg.HTTPIdleTimeout, err = time.ParseDuration(getEnv("HTTP_IDLE_TIMEOUT", "60s"))
	if err != nil {
		return nil, fmt.Errorf("invalid HTTP_IDLE_TIMEOUT: %w", err)
	}

	return cfg, nil
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
