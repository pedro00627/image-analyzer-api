package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/pedro00627/image-analyzer-api/internal/application/usecase"
	"github.com/pedro00627/image-analyzer-api/internal/infrastructure/ai_analyzer"
	"github.com/pedro00627/image-analyzer-api/internal/infrastructure/config"
	"github.com/pedro00627/image-analyzer-api/internal/infrastructure/validator"
)

// testContainer implements bootstrap.DependencyContainer for testing
type testContainer struct {
	cfg     config.Config
	useCase usecase.AnalyzeImageUseCaseInterface
}

func (tc *testContainer) GetConfig() config.Config {
	return tc.cfg
}

func (tc *testContainer) GetAnalyzeImageUseCase() usecase.AnalyzeImageUseCaseInterface {
	return tc.useCase
}

// testConfig implements config.Config for testing
type testConfig struct{}

func (tc *testConfig) GetProvider() string {
	return "mock"
}

func (tc *testConfig) GetSecret() string {
	return "secret"
}

func (tc *testConfig) GetAllowedMimes() []string {
	return []string{"image/jpeg", "image/png"}
}

func (tc *testConfig) GetMaxSize() int64 {
	return 10485760
}

func (tc *testConfig) GetRateLimitPerMinute() int {
	return 4
}

func (tc *testConfig) GetRateLimitBurst() int {
	return 4
}

func (tc *testConfig) GetAllowedOrigins() []string {
	return []string{"http://localhost:4200"}
}

// TestNewServer tests the New function creates a valid server
func TestNewServer(t *testing.T) {
	gin.SetMode(gin.TestMode)

	analyzer := ai_analyzer.NewMockAnalyzer()
	validator := validator.NewImageValidator(&testConfig{})
	uc := usecase.NewAnalyzeImageUseCase(validator, analyzer)

	container := &testContainer{
		cfg:     &testConfig{},
		useCase: uc,
	}

	server := New(container)

	if server == nil {
		t.Errorf("New() returned nil")
	}

	if _, ok := server.(*Server); !ok {
		t.Errorf("New() did not return a *Server")
	}
}

// TestHealthCheckEndpoint tests the health check endpoint
func TestHealthCheckEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)

	analyzer := ai_analyzer.NewMockAnalyzer()
	validator := validator.NewImageValidator(&testConfig{})
	uc := usecase.NewAnalyzeImageUseCase(validator, analyzer)

	container := &testContainer{
		cfg:     &testConfig{},
		useCase: uc,
	}

	srv := New(container).(*Server)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)

	srv.engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Health endpoint returned status %d, want 200", w.Code)
	}

	if w.Body.String() != `{"status":"healthy"}` {
		t.Errorf("Health endpoint returned %s, want {\"status\":\"healthy\"}", w.Body.String())
	}
}
