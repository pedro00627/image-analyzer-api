package config

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

//go:generate mockgen -source=./config.go -destination=./mocks/mock_config.go -package=mocks

// Config defines the application configuration interface
type Config interface {
	// AI Provider
	GetProvider() string
	GetSecret() string

	// Validator
	GetAllowedMimes() []string
	GetMaxSize() int64

	// Rate Limiting
	GetRateLimitPerMinute() int
	GetRateLimitBurst() int

	// Worker Pool
	GetWorkerPoolSize() int
	GetWorkerPoolQueueSize() int

	// Server
	GetAllowedOrigins() []string
}

// ConfigImpl is the concrete implementation of Config
type ConfigImpl struct {
	// AI Provider Configuration
	AIProvider string
	APIKey     string
	APISecret  string

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

	// Rate Limiting
	RateLimitPerMinute int
	RateLimitBurst     int

	// Worker Pool
	WorkerPoolSize      int
	WorkerPoolQueueSize int
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (Config, error) {
	return LoadConfigWithReader(envReader{})
}

// LoadConfigWithReader loads configuration using a custom Reader (for testing)
func LoadConfigWithReader(reader Reader) (Config, error) {
	cfg := &ConfigImpl{}

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

	if err := loadRateLimitConfig(cfg, reader); err != nil {
		return nil, err
	}

	if err := loadWorkerPoolConfig(cfg, reader); err != nil {
		return nil, err
	}

	return cfg, nil
}

// GetProvider implements Config
func (c *ConfigImpl) GetProvider() string {
	return c.AIProvider
}

// GetSecret implements Config
func (c *ConfigImpl) GetSecret() string {
	return c.APISecret
}

// GetAllowedMimes implements Config
func (c *ConfigImpl) GetAllowedMimes() []string {
	return c.AllowedFileTypes
}

// GetMaxSize implements Config
func (c *ConfigImpl) GetMaxSize() int64 {
	return c.MaxFileSize
}

// GetRateLimitPerMinute implements Config
func (c *ConfigImpl) GetRateLimitPerMinute() int {
	return c.RateLimitPerMinute
}

// GetRateLimitBurst implements Config
func (c *ConfigImpl) GetRateLimitBurst() int {
	return c.RateLimitBurst
}

// GetAllowedOrigins implements Config
func (c *ConfigImpl) GetAllowedOrigins() []string {
	return c.AllowedOrigins
}

// GetWorkerPoolSize implements Config
func (c *ConfigImpl) GetWorkerPoolSize() int {
	return c.WorkerPoolSize
}

// GetWorkerPoolQueueSize implements Config
func (c *ConfigImpl) GetWorkerPoolQueueSize() int {
	return c.WorkerPoolQueueSize
}

// loadAIProviderConfig loads AI provider configuration
func loadAIProviderConfig(cfg *ConfigImpl, reader Reader) error {
	cfg.AIProvider = reader.Get("AI_PROVIDER", "")
	if cfg.AIProvider == "" {
		return fmt.Errorf("AI_PROVIDER is required")
	}

	cfg.APIKey = reader.Get("API_KEY", "")
	cfg.APISecret = reader.Get("API_SECRET", "")

	return nil
}

// loadServerConfig loads server-related configuration
func loadServerConfig(cfg *ConfigImpl, reader Reader) error {
	port, err := strconv.Atoi(reader.Get("PORT", "8080"))
	if err != nil {
		return fmt.Errorf("invalid PORT: %w", err)
	}
	cfg.Port = port

	cfg.AllowedOrigins = parseCommaSeparated(reader.Get("ALLOWED_ORIGINS", "http://localhost:4200,http://localhost:3000"))

	return nil
}

// loadFileUploadConfig loads file upload limits and allowed types
func loadFileUploadConfig(cfg *ConfigImpl, reader Reader) error {
	maxSize, err := strconv.ParseInt(reader.Get("MAX_FILE_SIZE", "10485760"), 10, 64)
	if err != nil {
		return fmt.Errorf("invalid MAX_FILE_SIZE: %w", err)
	}
	cfg.MaxFileSize = maxSize

	cfg.AllowedFileTypes = parseCommaSeparated(reader.Get("ALLOWED_FILE_TYPES", "image/jpeg,image/png,image/gif,image/webp"))

	return nil
}

// loadTimeoutConfig loads all timeout configurations from environment variables
func loadTimeoutConfig(cfg *ConfigImpl, reader Reader) error {
	var err error

	cfg.AIServiceTimeout, err = parseDuration(reader, "AI_SERVICE_TIMEOUT", "")
	if err != nil {
		return err
	}

	cfg.HTTPReadTimeout, err = parseDuration(reader, "HTTP_READ_TIMEOUT", "")
	if err != nil {
		return err
	}

	cfg.HTTPWriteTimeout, err = parseDuration(reader, "HTTP_WRITE_TIMEOUT", "")
	if err != nil {
		return err
	}

	cfg.HTTPIdleTimeout, err = parseDuration(reader, "HTTP_IDLE_TIMEOUT", "")
	if err != nil {
		return err
	}

	return nil
}

// parseDuration parses a duration from Reader or validates it was provided
func parseDuration(reader Reader, key, defaultValue string) (time.Duration, error) {
	value := reader.Get(key, defaultValue)
	if value == "" {
		return 0, fmt.Errorf("%s is required", key)
	}
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

// loadRateLimitConfig loads rate limiting configuration
func loadRateLimitConfig(cfg *ConfigImpl, reader Reader) error {
	perMinute, err := strconv.Atoi(reader.Get("RATE_LIMIT_PER_MINUTE", "4"))
	if err != nil {
		return fmt.Errorf("invalid RATE_LIMIT_PER_MINUTE: %w", err)
	}
	if perMinute <= 0 {
		return fmt.Errorf("RATE_LIMIT_PER_MINUTE must be positive")
	}
	cfg.RateLimitPerMinute = perMinute

	burst, err := strconv.Atoi(reader.Get("RATE_LIMIT_BURST", "4"))
	if err != nil {
		return fmt.Errorf("invalid RATE_LIMIT_BURST: %w", err)
	}
	if burst <= 0 {
		return fmt.Errorf("RATE_LIMIT_BURST must be positive")
	}
	cfg.RateLimitBurst = burst

	return nil
}