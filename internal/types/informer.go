package types

import (
	"time"

	"github.com/thegostev/go-kubernetes-controllers/pkg/errors"
)

// InformerConfig represents informer configuration
type InformerConfig struct {
	Namespace       string        `json:"namespace"`
	ResyncPeriod    time.Duration `json:"resyncPeriod"`
	MaxCacheSize    int           `json:"maxCacheSize"`
	MaxConnections  int           `json:"maxConnections"`
	EventBufferSize int           `json:"eventBufferSize"`
	Workers         int           `json:"workers"`
}

// Validate validates InformerConfig
func (c *InformerConfig) Validate() error {
	if c.ResyncPeriod < time.Second || c.ResyncPeriod > 30*time.Minute {
		return errors.NewValidationError("resyncPeriod", "must be between 1s and 30m")
	}

	if c.MaxCacheSize <= 0 {
		return errors.NewValidationError("maxCacheSize", "must be positive")
	}

	if c.MaxConnections <= 0 {
		return errors.NewValidationError("maxConnections", "must be positive")
	}

	if c.EventBufferSize <= 0 {
		return errors.NewValidationError("eventBufferSize", "must be positive")
	}

	if c.Workers <= 0 {
		return errors.NewValidationError("workers", "must be positive")
	}

	return nil
}

// SetDefaults sets default values for InformerConfig
func (c *InformerConfig) SetDefaults() {
	if c.ResyncPeriod == 0 {
		c.ResyncPeriod = 10 * time.Minute
	}
	if c.MaxCacheSize == 0 {
		c.MaxCacheSize = 1000
	}
	if c.MaxConnections == 0 {
		c.MaxConnections = 10
	}
	if c.EventBufferSize == 0 {
		c.EventBufferSize = 100
	}
	if c.Workers == 0 {
		c.Workers = 2
	}
}

// InformerHealth represents informer health status
type InformerHealth struct {
	IsHealthy bool      `json:"isHealthy"`
	LastSync  time.Time `json:"lastSync"`
	Error     error     `json:"error,omitempty"`
	CacheSize int       `json:"cacheSize"`
	Workers   int       `json:"workers"`
}

// Event represents a deployment event
type Event struct {
	Type      string      `json:"type"` // "add", "update", "delete"
	Namespace string      `json:"namespace"`
	Name      string      `json:"name"`
	Timestamp time.Time   `json:"timestamp"`
	Object    interface{} `json:"object,omitempty"`
}
