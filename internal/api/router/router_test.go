package router

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/pedro00627/image-analyzer-api/internal/application/usecase"
	"github.com/pedro00627/image-analyzer-api/internal/infrastructure/ai_analyzer"
	"github.com/pedro00627/image-analyzer-api/internal/infrastructure/config"
	"github.com/pedro00627/image-analyzer-api/internal/infrastructure/config/mocks"
	"github.com/pedro00627/image-analyzer-api/internal/infrastructure/validator"
)

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

func TestSetupRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCfg := mocks.NewMockConfig(ctrl)
	mockCfg.EXPECT().GetAllowedMimes().Return([]string{"image/jpeg"}).AnyTimes()
	mockCfg.EXPECT().GetMaxSize().Return(int64(10485760)).AnyTimes()
	mockCfg.EXPECT().GetAllowedOrigins().Return([]string{"http://localhost:4200"}).AnyTimes()
	mockCfg.EXPECT().GetRateLimitPerMinute().Return(4).AnyTimes()
	mockCfg.EXPECT().GetRateLimitBurst().Return(4).AnyTimes()

	analyzer := ai_analyzer.NewMockAnalyzer()
	imgValidator := validator.NewImageValidator(mockCfg)
	uc := usecase.NewAnalyzeImageUseCase(imgValidator, analyzer)

	container := &testContainer{
		cfg:     mockCfg,
		useCase: uc,
	}

	engine := gin.New()
	
	// Should not panic
	SetupRoutes(engine, container)

	// Verify routes are registered
	routes := engine.Routes()
	
	foundAnalyzeRoute := false
	for _, route := range routes {
		if route.Path == "/api/analyze" && route.Method == "POST" {
			foundAnalyzeRoute = true
		}
	}

	if !foundAnalyzeRoute {
		t.Error("POST /api/analyze route not registered")
	}
}
