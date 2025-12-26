package config

import (
	"fmt"
	"strconv"
)

// loadWorkerPoolConfig loads worker pool configuration
func loadWorkerPoolConfig(cfg *ConfigImpl, reader Reader) error {
	poolSize, err := strconv.Atoi(reader.Get("WORKER_POOL_SIZE", "5"))
	if err != nil {
		return fmt.Errorf("invalid WORKER_POOL_SIZE: %w", err)
	}
	if poolSize <= 0 {
		return fmt.Errorf("WORKER_POOL_SIZE must be positive")
	}
	cfg.WorkerPoolSize = poolSize

	queueSize, err := strconv.Atoi(reader.Get("WORKER_POOL_QUEUE_SIZE", "100"))
	if err != nil {
		return fmt.Errorf("invalid WORKER_POOL_QUEUE_SIZE: %w", err)
	}
	if queueSize <= 0 {
		return fmt.Errorf("WORKER_POOL_QUEUE_SIZE must be positive")
	}
	cfg.WorkerPoolQueueSize = queueSize

	return nil
}
