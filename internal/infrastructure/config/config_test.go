package config

import (
	"os"
	"testing"
	"time"
)

const (
	loadConfigErrorMsg = "LoadConfig() error = %v"
)

// Helper to cast Config to *ConfigImpl for testing
func asImpl(t *testing.T, cfg Config) *ConfigImpl {
	t.Helper()
	impl, ok := cfg.(*ConfigImpl)
	if !ok {
		t.Fatal("Config is not *ConfigImpl")
	}
	return impl
}

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

	impl := asImpl(t, cfg)

	// Verify values
	if impl.AIProvider != "google_vision" {
		t.Errorf("AIProvider = %v, want google_vision", impl.AIProvider)
	}
	if impl.APIKey != "test_key_123" {
		t.Errorf("APIKey = %v, want test_key_123", impl.APIKey)
	}
	if impl.Port != 8080 {
		t.Errorf("Port = %v, want 8080", impl.Port)
	}
	if len(impl.AllowedOrigins) != 2 {
		t.Errorf("AllowedOrigins length = %v, want 2", len(impl.AllowedOrigins))
	}
	if impl.MaxFileSize != 10485760 {
		t.Errorf("MaxFileSize = %v, want 10485760", impl.MaxFileSize)
	}
	if impl.AIServiceTimeout != 30*time.Second {
		t.Errorf("AIServiceTimeout = %v, want 30s", impl.AIServiceTimeout)
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
	os.Setenv("AI_SERVICE_TIMEOUT", "30s")
	os.Setenv("HTTP_READ_TIMEOUT", "10s")
	os.Setenv("HTTP_WRITE_TIMEOUT", "10s")
	os.Setenv("HTTP_IDLE_TIMEOUT", "60s")
	defer cleanEnv()

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() unexpected error: %v", err)
	}
	impl := asImpl(t, cfg)
	if impl.AIProvider != "google_vision" {
		t.Errorf("AIProvider = %v, want google_vision", impl.AIProvider)
	}
}

func TestLoadConfig_InvalidAIProvider(t *testing.T) {
	os.Setenv("AI_PROVIDER", "invalid_provider")
	os.Setenv("API_KEY", "test_key")
	os.Setenv("AI_SERVICE_TIMEOUT", "30s")
	os.Setenv("HTTP_READ_TIMEOUT", "10s")
	os.Setenv("HTTP_WRITE_TIMEOUT", "10s")
	os.Setenv("HTTP_IDLE_TIMEOUT", "60s")
	defer cleanEnv()

	// No validation on provider value, just requires it to be set
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() unexpected error: %v", err)
	}
	impl := asImpl(t, cfg)
	if impl.AIProvider != "invalid_provider" {
		t.Errorf("AIProvider = %v, want invalid_provider", impl.AIProvider)
	}
}

func TestLoadConfig_InvalidPort(t *testing.T) {
	os.Setenv("AI_PROVIDER", "google_vision")
	os.Setenv("API_KEY", "test_key")
	os.Setenv("PORT", "invalid")
	os.Setenv("AI_SERVICE_TIMEOUT", "30s")
	os.Setenv("HTTP_READ_TIMEOUT", "10s")
	os.Setenv("HTTP_WRITE_TIMEOUT", "10s")
	os.Setenv("HTTP_IDLE_TIMEOUT", "60s")
	defer cleanEnv()

	_, err := LoadConfig()
	if err == nil {
		t.Error("LoadConfig() expected error for invalid PORT, got nil")
	}
}

func TestLoadConfig_InvalidTimeout(t *testing.T) {
	os.Setenv("AI_PROVIDER", "google_vision")
	os.Setenv("API_KEY", "test_key")
	os.Setenv("AI_SERVICE_TIMEOUT", "invalid")
	os.Setenv("HTTP_READ_TIMEOUT", "10s")
	os.Setenv("HTTP_WRITE_TIMEOUT", "10s")
	os.Setenv("HTTP_IDLE_TIMEOUT", "60s")
	defer cleanEnv()

	_, err := LoadConfig()
	if err == nil {
		t.Error("LoadConfig() expected error for invalid timeout, got nil")
	}
}

func TestLoadConfig_InvalidReadTimeout(t *testing.T) {
	os.Setenv("AI_PROVIDER", "google_vision")
	os.Setenv("API_KEY", "test_key")
	os.Setenv("AI_SERVICE_TIMEOUT", "30s")
	os.Setenv("HTTP_READ_TIMEOUT", "invalid")
	os.Setenv("HTTP_WRITE_TIMEOUT", "10s")
	os.Setenv("HTTP_IDLE_TIMEOUT", "60s")
	defer cleanEnv()

	_, err := LoadConfig()
	if err == nil {
		t.Error("LoadConfig() expected error for invalid HTTP_READ_TIMEOUT, got nil")
	}
}

func TestLoadConfig_InvalidWriteTimeout(t *testing.T) {
	os.Setenv("AI_PROVIDER", "google_vision")
	os.Setenv("API_KEY", "test_key")
	os.Setenv("AI_SERVICE_TIMEOUT", "30s")
	os.Setenv("HTTP_READ_TIMEOUT", "10s")
	os.Setenv("HTTP_WRITE_TIMEOUT", "invalid")
	os.Setenv("HTTP_IDLE_TIMEOUT", "60s")
	defer cleanEnv()

	_, err := LoadConfig()
	if err == nil {
		t.Error("LoadConfig() expected error for invalid HTTP_WRITE_TIMEOUT, got nil")
	}
}

func TestLoadConfig_InvalidIdleTimeout(t *testing.T) {
	os.Setenv("AI_PROVIDER", "google_vision")
	os.Setenv("API_KEY", "test_key")
	os.Setenv("AI_SERVICE_TIMEOUT", "30s")
	os.Setenv("HTTP_READ_TIMEOUT", "10s")
	os.Setenv("HTTP_WRITE_TIMEOUT", "10s")
	os.Setenv("HTTP_IDLE_TIMEOUT", "invalid")
	defer cleanEnv()

	_, err := LoadConfig()
	if err == nil {
		t.Error("LoadConfig() expected error for invalid HTTP_IDLE_TIMEOUT, got nil")
	}
}

func TestLoadConfig_Defaults(t *testing.T) {
	os.Setenv("AI_PROVIDER", "google_vision")
	os.Setenv("API_KEY", "test_key")
	os.Setenv("AI_SERVICE_TIMEOUT", "30s")
	os.Setenv("HTTP_READ_TIMEOUT", "10s")
	os.Setenv("HTTP_WRITE_TIMEOUT", "10s")
	os.Setenv("HTTP_IDLE_TIMEOUT", "60s")
	// Only required vars, others should use defaults
	defer cleanEnv()

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf(loadConfigErrorMsg, err)
	}

	impl := asImpl(t, cfg)
	if impl.Port != 8080 {
		t.Errorf("Port default = %v, want 8080", impl.Port)
	}
	if impl.MaxFileSize != 10485760 {
		t.Errorf("MaxFileSize default = %v, want 10485760", impl.MaxFileSize)
	}
}

func TestLoadConfig_InvalidMaxFileSize(t *testing.T) {
	os.Setenv("AI_PROVIDER", "google_vision")
	os.Setenv("API_KEY", "test_key")
	os.Setenv("MAX_FILE_SIZE", "not-a-number")
	os.Setenv("AI_SERVICE_TIMEOUT", "30s")
	os.Setenv("HTTP_READ_TIMEOUT", "10s")
	os.Setenv("HTTP_WRITE_TIMEOUT", "10s")
	os.Setenv("HTTP_IDLE_TIMEOUT", "60s")
	defer cleanEnv()

	_, err := LoadConfig()
	if err == nil {
		t.Error("LoadConfig() expected error for invalid MAX_FILE_SIZE, got nil")
	}
}

func TestConfigValidationMethods(t *testing.T) {
	cfg := &ConfigImpl{
		AllowedFileTypes: []string{"image/png"},
		MaxFileSize:      123,
	}

	if len(cfg.GetAllowedMimes()) != 1 || cfg.GetAllowedMimes()[0] != "image/png" {
		t.Fatalf("GetAllowedMimes = %v, want [image/png]", cfg.GetAllowedMimes())
	}
	if cfg.GetMaxSize() != 123 {
		t.Fatalf("GetMaxSize = %v, want 123", cfg.GetMaxSize())
	}
}

func TestLoadConfig_ImaggaProvider(t *testing.T) {
	os.Setenv("AI_PROVIDER", "imagga")
	os.Setenv("API_KEY", "imagga_key")
	os.Setenv("API_SECRET", "imagga_secret")
	os.Setenv("AI_SERVICE_TIMEOUT", "30s")
	os.Setenv("HTTP_READ_TIMEOUT", "10s")
	os.Setenv("HTTP_WRITE_TIMEOUT", "10s")
	os.Setenv("HTTP_IDLE_TIMEOUT", "60s")
	defer cleanEnv()

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf(loadConfigErrorMsg, err)
	}

	impl := asImpl(t, cfg)
	if impl.AIProvider != "imagga" {
		t.Errorf("AIProvider = %v, want imagga", impl.AIProvider)
	}
	if impl.APIKey != "imagga_key" {
		t.Errorf("APIKey = %v, want imagga_key", impl.APIKey)
	}
	if impl.APISecret != "imagga_secret" {
		t.Errorf("APISecret = %v, want imagga_secret", impl.APISecret)
	}
}

func TestLoadConfig_OpenAIProvider(t *testing.T) {
	os.Setenv("AI_PROVIDER", "openai")
	os.Setenv("API_KEY", "openai_key")
	os.Setenv("AI_SERVICE_TIMEOUT", "30s")
	os.Setenv("HTTP_READ_TIMEOUT", "10s")
	os.Setenv("HTTP_WRITE_TIMEOUT", "10s")
	os.Setenv("HTTP_IDLE_TIMEOUT", "60s")
	defer cleanEnv()

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf(loadConfigErrorMsg, err)
	}

	impl := asImpl(t, cfg)
	if impl.AIProvider != "openai" {
		t.Errorf("AIProvider = %v, want openai", impl.AIProvider)
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
