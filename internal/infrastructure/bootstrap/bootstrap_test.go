package bootstrap

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pedro00627/image-analyzer-api/internal/infrastructure/config/mocks"
)

// TestBootstrapWithConfig creates a container with mocked config
func TestBootstrapWithConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCfg := mocks.NewMockConfig(ctrl)
	mockCfg.EXPECT().GetProvider().Return("mock").AnyTimes()
	mockCfg.EXPECT().GetSecret().Return("secret").AnyTimes()
	mockCfg.EXPECT().GetAllowedMimes().Return([]string{"image/jpeg", "image/png"}).AnyTimes()
	mockCfg.EXPECT().GetMaxSize().Return(int64(10485760)).AnyTimes()

	ctx := context.Background()
	container, err := BootstrapWithConfig(ctx, mockCfg)

	if err != nil {
		t.Errorf("BootstrapWithConfig() error = %v, want nil", err)
	}

	if container == nil {
		t.Errorf("BootstrapWithConfig() returned nil container")
	}
}

// TestContainerImplementsDependencyContainer verifies interface implementation
func TestContainerImplementsDependencyContainer(t *testing.T) {
	var _ DependencyContainer = (*Container)(nil)
}

// TestContainer_GetConfig tests the GetConfig method
func TestContainer_GetConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCfg := mocks.NewMockConfig(ctrl)
	mockCfg.EXPECT().GetProvider().Return("mock").AnyTimes()
	mockCfg.EXPECT().GetSecret().Return("secret").AnyTimes()
	mockCfg.EXPECT().GetAllowedMimes().Return([]string{"image/jpeg"}).AnyTimes()
	mockCfg.EXPECT().GetMaxSize().Return(int64(10485760)).AnyTimes()

	ctx := context.Background()
	container, err := BootstrapWithConfig(ctx, mockCfg)
	if err != nil {
		t.Fatalf("BootstrapWithConfig() error = %v", err)
	}

	cfg := container.GetConfig()
	if cfg == nil {
		t.Errorf("GetConfig() returned nil")
	}

	if cfg != mockCfg {
		t.Errorf("GetConfig() returned different config")
	}
}

// TestContainer_GetAnalyzeImageUseCase tests the GetAnalyzeImageUseCase method
func TestContainer_GetAnalyzeImageUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCfg := mocks.NewMockConfig(ctrl)
	mockCfg.EXPECT().GetProvider().Return("mock").AnyTimes()
	mockCfg.EXPECT().GetSecret().Return("secret").AnyTimes()
	mockCfg.EXPECT().GetAllowedMimes().Return([]string{"image/jpeg"}).AnyTimes()
	mockCfg.EXPECT().GetMaxSize().Return(int64(10485760)).AnyTimes()

	ctx := context.Background()
	container, err := BootstrapWithConfig(ctx, mockCfg)
	if err != nil {
		t.Fatalf("BootstrapWithConfig() error = %v", err)
	}

	uc := container.GetAnalyzeImageUseCase()
	if uc == nil {
		t.Errorf("GetAnalyzeImageUseCase() returned nil")
	}
}

// TestContainer_GetConfig_ReturnsSameInstance tests idempotency
func TestContainer_GetConfig_ReturnsSameInstance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCfg := mocks.NewMockConfig(ctrl)
	mockCfg.EXPECT().GetProvider().Return("mock").AnyTimes()
	mockCfg.EXPECT().GetSecret().Return("secret").AnyTimes()
	mockCfg.EXPECT().GetAllowedMimes().Return([]string{"image/jpeg"}).AnyTimes()
	mockCfg.EXPECT().GetMaxSize().Return(int64(10485760)).AnyTimes()

	ctx := context.Background()
	container, err := BootstrapWithConfig(ctx, mockCfg)
	if err != nil {
		t.Fatalf("BootstrapWithConfig() error = %v", err)
	}

	cfg1 := container.GetConfig()
	cfg2 := container.GetConfig()

	if cfg1 != cfg2 {
		t.Errorf("GetConfig() returned different instances")
	}
}

// TestInitializeAIAnalyzer_FallsBackToMock tests that initializeAIAnalyzer falls back to mock
func TestInitializeAIAnalyzer_FallsBackToMock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCfg := mocks.NewMockConfig(ctrl)
	mockCfg.EXPECT().GetProvider().Return("invalid").AnyTimes()
	mockCfg.EXPECT().GetSecret().Return("secret").AnyTimes()
	mockCfg.EXPECT().GetAllowedMimes().Return([]string{"image/jpeg"}).AnyTimes()
	mockCfg.EXPECT().GetMaxSize().Return(int64(10485760)).AnyTimes()

	ctx := context.Background()
	analyzer := initializeAIAnalyzer(ctx, mockCfg)
	if analyzer == nil {
		t.Errorf("initializeAIAnalyzer() returned nil analyzer")
	}
}


