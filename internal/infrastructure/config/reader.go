// Package config handles application configuration from environment variables
package config

import "os"

// Reader defines the interface for reading configuration values
type Reader interface {
	Get(key, defaultValue string) string
}

// envReader reads configuration from environment variables
type envReader struct{}

func (e envReader) Get(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
