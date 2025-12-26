package ai_analyzer

import (
	"context"
	"fmt"
	"sync"

	"github.com/pedro00627/image-analyzer-api/internal/domain/entity"
	"github.com/pedro00627/image-analyzer-api/internal/domain/service"
)

// AnalysisJob represents a single image analysis job
type AnalysisJob struct {
	ctx       context.Context
	imageData []byte
	resultCh  chan<- *JobResult
}

// JobResult represents the result of an analysis job
type JobResult struct {
	Result *entity.AnalysisResult
	Error  error
}

// AnalyzerWorkerPool wraps an AIAnalyzer with a worker pool pattern
// to limit concurrent calls to the underlying AI service
type AnalyzerWorkerPool struct {
	analyzer  service.AIAnalyzer
	workers   int
	queueSize int
	jobsCh    chan *AnalysisJob
	wg        sync.WaitGroup
	once      sync.Once
	stopCh    chan struct{}
}

var _ service.AIAnalyzer = (*AnalyzerWorkerPool)(nil)

// NewAnalyzerWorkerPool creates a new worker pool for image analysis
// workers: number of concurrent workers
// queueSize: size of the job queue buffer
func NewAnalyzerWorkerPool(analyzer service.AIAnalyzer, workers, queueSize int) (*AnalyzerWorkerPool, error) {
	if analyzer == nil {
		return nil, fmt.Errorf("analyzer cannot be nil")
	}
	if workers <= 0 {
		return nil, fmt.Errorf("workers must be positive, got: %d", workers)
	}
	if queueSize <= 0 {
		return nil, fmt.Errorf("queueSize must be positive, got: %d", queueSize)
	}

	pool := &AnalyzerWorkerPool{
		analyzer:  analyzer,
		workers:   workers,
		queueSize: queueSize,
		jobsCh:    make(chan *AnalysisJob, queueSize),
		stopCh:    make(chan struct{}),
	}

	pool.start()
	return pool, nil
}

// start initializes the worker pool
func (p *AnalyzerWorkerPool) start() {
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go p.worker(i)
	}
}

// worker processes jobs from the queue
func (p *AnalyzerWorkerPool) worker(id int) {
	defer p.wg.Done()

	for {
		select {
		case <-p.stopCh:
			return
		case job, ok := <-p.jobsCh:
			if !ok {
				return
			}
			p.processJob(job)
		}
	}
}

// processJob executes a single analysis job
func (p *AnalyzerWorkerPool) processJob(job *AnalysisJob) {
	result, err := p.analyzer.Analyze(job.ctx, job.imageData)
	job.resultCh <- &JobResult{
		Result: result,
		Error:  err,
	}
}

// Analyze implements service.AIAnalyzer by delegating to the worker pool
func (p *AnalyzerWorkerPool) Analyze(ctx context.Context, imageData []byte) (*entity.AnalysisResult, error) {
	// Check if pool is shutting down first (before attempting to send)
	select {
	case <-p.stopCh:
		return nil, fmt.Errorf("worker pool is shutting down")
	default:
		// Pool is still active, continue
	}

	resultCh := make(chan *JobResult, 1)

	job := &AnalysisJob{
		ctx:       ctx,
		imageData: imageData,
		resultCh:  resultCh,
	}

	// Try to enqueue the job
	select {
	case p.jobsCh <- job:
		// Job enqueued successfully
	case <-ctx.Done():
		return nil, fmt.Errorf("context cancelled while enqueuing job: %w", ctx.Err())
	case <-p.stopCh:
		return nil, fmt.Errorf("worker pool is shutting down")
	}

	// Wait for the result
	select {
	case result := <-resultCh:
		return result.Result, result.Error
	case <-ctx.Done():
		return nil, fmt.Errorf("context cancelled while waiting for result: %w", ctx.Err())
	case <-p.stopCh:
		return nil, fmt.Errorf("worker pool is shutting down")
	}
}

// Close gracefully shuts down the worker pool
func (p *AnalyzerWorkerPool) Close() error {
	var closeErr error

	p.once.Do(func() {
		// Signal workers to stop by closing stopCh
		// This will cause all workers to exit their select loops
		close(p.stopCh)

		// Wait for all workers to finish processing current jobs
		p.wg.Wait()

		// Note: We intentionally do NOT close jobsCh here because:
		// 1. It could cause panics if someone tries to send to it
		// 2. The stopCh signal is sufficient to prevent new work
		// 3. The Go GC will clean up the channel when the pool is garbage collected

		// Close the underlying analyzer if it implements Close
		if closer, ok := p.analyzer.(interface{ Close() error }); ok {
			closeErr = closer.Close()
		}
	})

	return closeErr
}

// Stats returns statistics about the worker pool
func (p *AnalyzerWorkerPool) Stats() PoolStats {
	return PoolStats{
		Workers:       p.workers,
		QueueSize:     p.queueSize,
		QueuedJobs:    len(p.jobsCh),
		QueueCapacity: cap(p.jobsCh),
	}
}

// PoolStats contains statistics about the worker pool
type PoolStats struct {
	Workers       int
	QueueSize     int
	QueuedJobs    int
	QueueCapacity int
}
