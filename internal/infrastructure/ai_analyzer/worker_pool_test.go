package ai_analyzer

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/pedro00627/image-analyzer-api/internal/domain/entity"
	"github.com/stretchr/testify/assert"
)

// mockAnalyzerForPool is a mock implementation of service.AIAnalyzer for testing the worker pool
type mockAnalyzerForPool struct {
	mu          sync.Mutex
	callCount   int
	callDelay   time.Duration
	shouldBlock bool
	shouldError bool
	errorMsg    string
}

func (m *mockAnalyzerForPool) Analyze(ctx context.Context, imageData []byte) (*entity.AnalysisResult, error) {
	m.mu.Lock()
	m.callCount++
	m.mu.Unlock()

	// Simulate processing delay
	if m.callDelay > 0 {
		select {
		case <-time.After(m.callDelay):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	// Block if configured
	if m.shouldBlock {
		<-ctx.Done()
		return nil, ctx.Err()
	}

	// Return error if configured
	if m.shouldError {
		return nil, errors.New(m.errorMsg)
	}

	// Return success
	tags := []entity.Tag{entity.NewTag("test", 0.95)}
	result := entity.NewAnalysisResult(tags)
	return &result, nil
}

func (m *mockAnalyzerForPool) GetCallCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.callCount
}

func TestNewAnalyzerWorkerPool(t *testing.T) {
	t.Run("valid pool creation", func(t *testing.T) {
		analyzer := &mockAnalyzerForPool{}
		pool, err := NewAnalyzerWorkerPool(analyzer, 3, 10)

		assert.NoError(t, err)
		assert.NotNil(t, pool)
		assert.Equal(t, 3, pool.workers)
		assert.Equal(t, 10, pool.queueSize)

		pool.Close()
	})

	t.Run("nil analyzer", func(t *testing.T) {
		pool, err := NewAnalyzerWorkerPool(nil, 3, 10)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "analyzer cannot be nil")
		assert.Nil(t, pool)
	})

	t.Run("zero workers", func(t *testing.T) {
		analyzer := &mockAnalyzerForPool{}
		pool, err := NewAnalyzerWorkerPool(analyzer, 0, 10)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "workers must be positive")
		assert.Nil(t, pool)
	})

	t.Run("negative workers", func(t *testing.T) {
		analyzer := &mockAnalyzerForPool{}
		pool, err := NewAnalyzerWorkerPool(analyzer, -1, 10)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "workers must be positive")
		assert.Nil(t, pool)
	})

	t.Run("zero queue size", func(t *testing.T) {
		analyzer := &mockAnalyzerForPool{}
		pool, err := NewAnalyzerWorkerPool(analyzer, 3, 0)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "queueSize must be positive")
		assert.Nil(t, pool)
	})

	t.Run("negative queue size", func(t *testing.T) {
		analyzer := &mockAnalyzerForPool{}
		pool, err := NewAnalyzerWorkerPool(analyzer, 3, -1)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "queueSize must be positive")
		assert.Nil(t, pool)
	})
}

func TestAnalyzerWorkerPool_Analyze_Success(t *testing.T) {
	mockAnalyzer := &mockAnalyzerForPool{}

	pool, err := NewAnalyzerWorkerPool(mockAnalyzer, 2, 5)
	assert.NoError(t, err)
	defer pool.Close()

	ctx := context.Background()
	imageData := []byte("test image data")

	result, err := pool.Analyze(ctx, imageData)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, len(result.Tags))
	assert.Equal(t, "test", result.Tags[0].Label)
}

func TestAnalyzerWorkerPool_Analyze_Error(t *testing.T) {
	mockAnalyzer := &mockAnalyzerForPool{
		shouldError: true,
		errorMsg:    "analysis failed",
	}

	pool, err := NewAnalyzerWorkerPool(mockAnalyzer, 2, 5)
	assert.NoError(t, err)
	defer pool.Close()

	ctx := context.Background()
	imageData := []byte("test image data")

	result, err := pool.Analyze(ctx, imageData)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "analysis failed")
}

func TestAnalyzerWorkerPool_Analyze_ContextCancellation(t *testing.T) {
	mockAnalyzer := &mockAnalyzerForPool{
		callDelay: 100 * time.Millisecond,
	}

	pool, err := NewAnalyzerWorkerPool(mockAnalyzer, 2, 5)
	assert.NoError(t, err)
	defer pool.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	imageData := []byte("test image data")

	result, err := pool.Analyze(ctx, imageData)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "context cancelled")
}

func TestAnalyzerWorkerPool_Analyze_ContextTimeout(t *testing.T) {
	mockAnalyzer := &mockAnalyzerForPool{
		callDelay: 200 * time.Millisecond,
	}

	pool, err := NewAnalyzerWorkerPool(mockAnalyzer, 2, 5)
	assert.NoError(t, err)
	defer pool.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	imageData := []byte("test image data")

	result, err := pool.Analyze(ctx, imageData)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "context cancelled")
}

func TestAnalyzerWorkerPool_ConcurrentAnalysis(t *testing.T) {
	mockAnalyzer := &mockAnalyzerForPool{
		callDelay: 50 * time.Millisecond,
	}

	workers := 3
	pool, err := NewAnalyzerWorkerPool(mockAnalyzer, workers, 10)
	assert.NoError(t, err)
	defer pool.Close()

	numJobs := 10
	var wg sync.WaitGroup
	errors := make(chan error, numJobs)
	results := make(chan *entity.AnalysisResult, numJobs)

	start := time.Now()

	for i := 0; i < numJobs; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			ctx := context.Background()
			imageData := []byte("test image data")

			result, err := pool.Analyze(ctx, imageData)
			if err != nil {
				errors <- err
			} else {
				results <- result
			}
		}(i)
	}

	wg.Wait()
	close(errors)
	close(results)
	duration := time.Since(start)

	// Check no errors
	assert.Equal(t, 0, len(errors))

	// Check all results received
	assert.Equal(t, numJobs, len(results))

	// With 3 workers and 50ms delay, 10 jobs should take roughly:
	// ceil(10/3) * 50ms = 4 * 50ms = 200ms
	// Allow some margin for overhead
	assert.Less(t, duration, 500*time.Millisecond, "Should complete faster with workers")
}

func TestAnalyzerWorkerPool_QueueLimit(t *testing.T) {
	mockAnalyzer := &mockAnalyzerForPool{
		callDelay: 500 * time.Millisecond, // Use delay instead of blocking
	}

	workers := 1
	queueSize := 2
	pool, err := NewAnalyzerWorkerPool(mockAnalyzer, workers, queueSize)
	assert.NoError(t, err)
	defer pool.Close()

	// Fill the worker and queue
	imageData := []byte("test data")

	// Launch jobs that will fill the worker and queue
	// 1 worker + 2 queue = 3 total capacity
	for i := 0; i < workers+queueSize; i++ {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			pool.Analyze(ctx, imageData)
		}()
	}

	// Wait for queue to fill
	time.Sleep(100 * time.Millisecond)

	// This one should block since queue is full, use short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	result, err := pool.Analyze(ctx, imageData)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "context cancelled")
}

func TestAnalyzerWorkerPool_Close(t *testing.T) {
	mockAnalyzer := &mockAnalyzerForPool{}

	pool, err := NewAnalyzerWorkerPool(mockAnalyzer, 2, 5)
	assert.NoError(t, err)

	// Close the pool
	err = pool.Close()
	assert.NoError(t, err)

	// Subsequent closes should not panic
	err = pool.Close()
	assert.NoError(t, err)

	// Try to analyze after close
	ctx := context.Background()
	imageData := []byte("test data")

	result, err := pool.Analyze(ctx, imageData)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "shutting down")
}

func TestAnalyzerWorkerPool_Stats(t *testing.T) {
	mockAnalyzer := &mockAnalyzerForPool{}

	workers := 3
	queueSize := 10
	pool, err := NewAnalyzerWorkerPool(mockAnalyzer, workers, queueSize)
	assert.NoError(t, err)
	defer pool.Close()

	stats := pool.Stats()

	assert.Equal(t, workers, stats.Workers)
	assert.Equal(t, queueSize, stats.QueueSize)
	assert.Equal(t, queueSize, stats.QueueCapacity)
	assert.Equal(t, 0, stats.QueuedJobs) // No jobs queued yet
}

func TestAnalyzerWorkerPool_WorkerConcurrency(t *testing.T) {
	mockAnalyzer := &mockAnalyzerForPool{
		callDelay: 100 * time.Millisecond,
	}

	workers := 5
	pool, err := NewAnalyzerWorkerPool(mockAnalyzer, workers, 20)
	assert.NoError(t, err)
	defer pool.Close()

	numJobs := 15
	var wg sync.WaitGroup

	start := time.Now()

	for i := 0; i < numJobs; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctx := context.Background()
			imageData := []byte("test data")
			pool.Analyze(ctx, imageData)
		}()
	}

	wg.Wait()
	duration := time.Since(start)

	// With 5 workers and 100ms delay, 15 jobs should take roughly:
	// ceil(15/5) * 100ms = 3 * 100ms = 300ms
	expectedDuration := time.Duration(((numJobs + workers - 1) / workers)) * mockAnalyzer.callDelay

	// Allow 200ms overhead
	assert.Less(t, duration, expectedDuration+200*time.Millisecond)

	// Verify we used concurrency (should be much faster than sequential)
	sequentialDuration := time.Duration(numJobs) * mockAnalyzer.callDelay
	assert.Less(t, duration, sequentialDuration)
}
