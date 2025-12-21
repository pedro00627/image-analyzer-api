package config

import (
	"os"
	"testing"
	"time"
)

const (
	loadConfigErrorMsg = "LoadConfig() error = %v"
)

func TestLoadConfig_Success(t *testing.T) {
	// Setup test environment
	os.Setenv("AI_PROVIDER", "google_vision")
	os.Setenv("API_KEY", "test_key_123")
	os.Setenv("PORT", "8080")
	os.Setenv("ALLOWED_ORIGINS", "http://localhost:4200,http://localhost:3000")
	os.Setenv("MAX_FILE_SIZE", "10485760")
	os.Setenv("ALLOWED_FILE_TYPES", "image/jpeg,image/png")
	os.Setenv("AI_SERVICE_TIMEOUT", "30s")
	os.Setenv("HTTP_READ_TIMEOUT", "10s")
	os.Setenv("HTTP_WRITE_TIMEOUT", "10s")
	os.Setenv("HTTP_IDLE_TIMEOUT", "60s")
	defer cleanEnv()

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf(loadConfigErrorMsg, err)
	}

	// Verify values
	if cfg.AIProvider != "google_vision" {
		t.Errorf("AIProvider = %v, want google_vision", cfg.AIProvider)
	}
	if cfg.APIKey != "test_key_123" {
		t.Errorf("APIKey = %v, want test_key_123", cfg.APIKey)
	}
	if cfg.Port != 8080 {
		t.Errorf("Port = %v, want 8080", cfg.Port)
	}
	if len(cfg.AllowedOrigins) != 2 {
		t.Errorf("AllowedOrigins length = %v, want 2", len(cfg.AllowedOrigins))
	}
	if cfg.MaxFileSize != 10485760 {
		t.Errorf("MaxFileSize = %v, want 10485760", cfg.MaxFileSize)
	}
	if cfg.AIServiceTimeout != 30*time.Second {
		t.Errorf("AIServiceTimeout = %v, want 30s", cfg.AIServiceTimeout)
	}
}

func TestLoadConfig_MissingAIProvider(t *testing.T) {
	cleanEnv()

	_, err := LoadConfig()
	if err == nil {
		t.Error("LoadConfig() expected error for missing AI_PROVIDER, got nil")
	}
}

func TestLoadConfig_MissingAPIKey(t *testing.T) {
	os.Setenv("AI_PROVIDER", "google_vision")
	// Missing API_KEY - no error, credential validation is done by adapters
	defer cleanEnv()

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() unexpected error: %v", err)
	}
	if cfg.AIProvider != "google_vision" {
		t.Errorf("AIProvider = %v, want google_vision", cfg.AIProvider)
	}
}

func TestLoadConfig_InvalidAIProvider(t *testing.T) {
	os.Setenv("AI_PROVIDER", "invalid_provider")
	os.Setenv("API_KEY", "test_key")
	defer cleanEnv()

	// No validation on provider value, just requires it to be set
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() unexpected error: %v", err)
	}
	if cfg.AIProvider != "invalid_provider" {
		t.Errorf("AIProvider = %v, want invalid_provider", cfg.AIProvider)
	}
}

func TestLoadConfig_InvalidPort(t *testing.T) {
	os.Setenv("AI_PROVIDER", "google_vision")
	os.Setenv("API_KEY", "test_key")
	os.Setenv("PORT", "invalid")
	defer cleanEnv()

	_, err := LoadConfig()
	if err == nil {
		t.Error("LoadConfig() expected error for invalid PORT, got nil")
	}
}

func TestLoadConfig_InvalidTimeout(t *testing.T) {
	os.Setenv("AI_PROVIDER", "google_vision")
	os.Setenv("GOOGLE_VISION_API_KEY", "test_key")
	os.Setenv("AI_SERVICE_TIMEOUT", "invalid")
	defer cleanEnv()

	_, err := LoadConfig()
	if err == nil {
		t.Error("LoadConfig() expected error for invalid timeout, got nil")
	}
}

func TestLoadConfig_InvalidReadTimeout(t *testing.T) {
	os.Setenv("AI_PROVIDER", "google_vision")
	os.Setenv("API_KEY", "test_key")
	os.Setenv("HTTP_READ_TIMEOUT", "invalid")
	defer cleanEnv()

	_, err := LoadConfig()
	if err == nil {
		t.Error("LoadConfig() expected error for invalid HTTP_READ_TIMEOUT, got nil")
	}
}

func TestLoadConfig_InvalidWriteTimeout(t *testing.T) {
	os.Setenv("AI_PROVIDER", "google_vision")
	os.Setenv("API_KEY", "test_key")
	os.Setenv("HTTP_WRITE_TIMEOUT", "invalid")
	defer cleanEnv()

	_, err := LoadConfig()
	if err == nil {
		t.Error("LoadConfig() expected error for invalid HTTP_WRITE_TIMEOUT, got nil")
	}
}

func TestLoadConfig_InvalidIdleTimeout(t *testing.T) {
	os.Setenv("AI_PROVIDER", "google_vision")
	os.Setenv("API_KEY", "test_key")
	os.Setenv("HTTP_IDLE_TIMEOUT", "invalid")
	defer cleanEnv()

	_, err := LoadConfig()
	if err == nil {
		t.Error("LoadConfig() expected error for invalid HTTP_IDLE_TIMEOUT, got nil")
	}
}

func TestLoadConfig_Defaults(t *testing.T) {
	os.Setenv("AI_PROVIDER", "google_vision")
	os.Setenv("GOOGLE_VISION_API_KEY", "test_key")
	// Only required vars, others should use defaults
	defer cleanEnv()

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf(loadConfigErrorMsg, err)
	}

	if cfg.Port != 8080 {
		t.Errorf("Port default = %v, want 8080", cfg.Port)
	}
	if cfg.MaxFileSize != 10485760 {
		t.Errorf("MaxFileSize default = %v, want 10485760", cfg.MaxFileSize)
	}
}

func TestLoadConfig_InvalidMaxFileSize(t *testing.T) {
	os.Setenv("AI_PROVIDER", "google_vision")
	os.Setenv("API_KEY", "test_key")
	os.Setenv("MAX_FILE_SIZE", "not-a-number")
	defer cleanEnv()

	_, err := LoadConfig()
	if err == nil {
		t.Error("LoadConfig() expected error for invalid MAX_FILE_SIZE, got nil")
	}
}

func TestConfigValidationConfigDTO(t *testing.T) {
	cfg := &Config{
		AllowedFileTypes: []string{"image/png"},
		MaxFileSize:      123,
	}

	dto := cfg.ValidationConfig()

	if len(dto.AllowedMimes) != 1 || dto.AllowedMimes[0] != "image/png" {
		t.Fatalf("ValidationConfig AllowedMimes = %v, want [image/png]", dto.AllowedMimes)
	}
	if dto.MaxSize != 123 {
		t.Fatalf("ValidationConfig MaxSize = %v, want 123", dto.MaxSize)
	}
}

func TestLoadConfig_ImaggaProvider(t *testing.T) {
	os.Setenv("AI_PROVIDER", "imagga")
	os.Setenv("API_KEY", "imagga_key")
	os.Setenv("API_SECRET", "imagga_secret")
	defer cleanEnv()

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf(loadConfigErrorMsg, err)
	}

	if cfg.AIProvider != "imagga" {
		t.Errorf("AIProvider = %v, want imagga", cfg.AIProvider)
	}
	if cfg.APIKey != "imagga_key" {
		t.Errorf("APIKey = %v, want imagga_key", cfg.APIKey)
	}
	if cfg.APISecret != "imagga_secret" {
		t.Errorf("APISecret = %v, want imagga_secret", cfg.APISecret)
	}
}

func TestLoadConfig_OpenAIProvider(t *testing.T) {
	os.Setenv("AI_PROVIDER", "openai")
	os.Setenv("API_KEY", "openai_key")
	defer cleanEnv()

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf(loadConfigErrorMsg, err)
	}

	if cfg.AIProvider != "openai" {
		t.Errorf("AIProvider = %v, want openai", cfg.AIProvider)
	}
}

// cleanEnv clears all test environment variables
func cleanEnv() {
	os.Unsetenv("AI_PROVIDER")
	os.Unsetenv("API_KEY")
	os.Unsetenv("API_SECRET")
	os.Unsetenv("PORT")
	os.Unsetenv("ALLOWED_ORIGINS")
	os.Unsetenv("MAX_FILE_SIZE")
	os.Unsetenv("ALLOWED_FILE_TYPES")
	os.Unsetenv("AI_SERVICE_TIMEOUT")
	os.Unsetenv("HTTP_READ_TIMEOUT")
	os.Unsetenv("HTTP_WRITE_TIMEOUT")
	os.Unsetenv("HTTP_IDLE_TIMEOUT")
}
