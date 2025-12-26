package ai_analyzer

import (
	"context"
	"testing"

	"github.com/pedro00627/image-analyzer-api/internal/infrastructure/config"
	"github.com/stretchr/testify/assert"
)

// mockConfig for integration test
type mockConfigForIntegration struct{}

func (m *mockConfigForIntegration) GetProvider() string         { return "google_vision" }
func (m *mockConfigForIntegration) GetSecret() string           { return "" }
func (m *mockConfigForIntegration) GetAllowedMimes() []string   { return []string{} }
func (m *mockConfigForIntegration) GetMaxSize() int64           { return 0 }
func (m *mockConfigForIntegration) GetRateLimitPerMinute() int  { return 0 }
func (m *mockConfigForIntegration) GetRateLimitBurst() int      { return 0 }
func (m *mockConfigForIntegration) GetWorkerPoolSize() int      { return 5 }
func (m *mockConfigForIntegration) GetWorkerPoolQueueSize() int { return 100 }
func (m *mockConfigForIntegration) GetAllowedOrigins() []string { return []string{} }

var _ config.Config = (*mockConfigForIntegration)(nil)

// TestAnalyzerWorkerPool_CloseWithRealAnalyzer tests closing a pool with a real analyzer that implements Close
func TestAnalyzerWorkerPool_CloseWithRealAnalyzer(t *testing.T) {
	ctx := context.Background()
	cfg := &mockConfigForIntegration{}

	// Create a real analyzer (will use mock internally if no credentials)
	baseAnalyzer, err := NewGoogleVisionAnalyzer(ctx, cfg)
	if err != nil {
		// If we can't create the analyzer (no credentials), use mock with Close
		t.Skip("Skipping integration test - no Google Vision credentials")
	}

	// Wrap it in a worker pool
	pool, err := NewAnalyzerWorkerPool(baseAnalyzer, 2, 5)
	assert.NoError(t, err)
	assert.NotNil(t, pool)

	// Close the pool - this should call Close on the underlying analyzer
	err = pool.Close()
	assert.NoError(t, err)

	// Verify we can't analyze after close
	result, err := pool.Analyze(ctx, []byte("test"))
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "shutting down")
}
