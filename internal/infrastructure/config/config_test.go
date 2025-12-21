package config

import (
	"os"
	"testing"
	"time"
)

func TestLoadConfig_Success(t *testing.T) {
	// Setup test environment
	os.Setenv("AI_PROVIDER", "google_vision")
	os.Setenv("GOOGLE_VISION_API_KEY", "test_key_123")
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
		t.Fatalf("LoadConfig() error = %v", err)
	}

	// Verify values
	if cfg.AIProvider != "google_vision" {
		t.Errorf("AIProvider = %v, want google_vision", cfg.AIProvider)
	}
	if cfg.GoogleVisionAPIKey != "test_key_123" {
		t.Errorf("GoogleVisionAPIKey = %v, want test_key_123", cfg.GoogleVisionAPIKey)
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
	// Missing GOOGLE_VISION_API_KEY
	defer cleanEnv()

	_, err := LoadConfig()
	if err == nil {
		t.Error("LoadConfig() expected error for missing API key, got nil")
	}
}

func TestLoadConfig_InvalidAIProvider(t *testing.T) {
	os.Setenv("AI_PROVIDER", "invalid_provider")
	os.Setenv("GOOGLE_VISION_API_KEY", "test_key")
	defer cleanEnv()

	_, err := LoadConfig()
	if err == nil {
		t.Error("LoadConfig() expected error for invalid AI_PROVIDER, got nil")
	}
}

func TestLoadConfig_InvalidPort(t *testing.T) {
	os.Setenv("AI_PROVIDER", "google_vision")
	os.Setenv("GOOGLE_VISION_API_KEY", "test_key")
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

func TestLoadConfig_Defaults(t *testing.T) {
	os.Setenv("AI_PROVIDER", "google_vision")
	os.Setenv("GOOGLE_VISION_API_KEY", "test_key")
	// Only required vars, others should use defaults
	defer cleanEnv()

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if cfg.Port != 8080 {
		t.Errorf("Port default = %v, want 8080", cfg.Port)
	}
	if cfg.MaxFileSize != 10485760 {
		t.Errorf("MaxFileSize default = %v, want 10485760", cfg.MaxFileSize)
	}
}

func TestLoadConfig_ImaggaProvider(t *testing.T) {
	os.Setenv("AI_PROVIDER", "imagga")
	os.Setenv("IMAGGA_API_KEY", "imagga_key")
	os.Setenv("IMAGGA_API_SECRET", "imagga_secret")
	defer cleanEnv()

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if cfg.AIProvider != "imagga" {
		t.Errorf("AIProvider = %v, want imagga", cfg.AIProvider)
	}
	if cfg.ImaggaAPIKey != "imagga_key" {
		t.Errorf("ImaggaAPIKey = %v", cfg.ImaggaAPIKey)
	}
}

func TestLoadConfig_OpenAIProvider(t *testing.T) {
	os.Setenv("AI_PROVIDER", "openai")
	os.Setenv("OPENAI_API_KEY", "openai_key")
	defer cleanEnv()

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if cfg.AIProvider != "openai" {
		t.Errorf("AIProvider = %v, want openai", cfg.AIProvider)
	}
}

// cleanEnv clears all test environment variables
func cleanEnv() {
	os.Unsetenv("AI_PROVIDER")
	os.Unsetenv("GOOGLE_VISION_API_KEY")
	os.Unsetenv("IMAGGA_API_KEY")
	os.Unsetenv("IMAGGA_API_SECRET")
	os.Unsetenv("OPENAI_API_KEY")
	os.Unsetenv("PORT")
	os.Unsetenv("ALLOWED_ORIGINS")
	os.Unsetenv("MAX_FILE_SIZE")
	os.Unsetenv("ALLOWED_FILE_TYPES")
	os.Unsetenv("AI_SERVICE_TIMEOUT")
	os.Unsetenv("HTTP_READ_TIMEOUT")
	os.Unsetenv("HTTP_WRITE_TIMEOUT")
	os.Unsetenv("HTTP_IDLE_TIMEOUT")
}
