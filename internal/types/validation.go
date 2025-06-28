package types

import (
	"time"

	"github.com/yourusername/k8s-controller-tutorial/pkg/errors"
)

// ListOptions represents options for listing resources
type ListOptions struct {
	Namespace string        `json:"namespace"`
	Timeout   time.Duration `json:"timeout"`
}

// Validate validates ListOptions
func (o *ListOptions) Validate() error {
	// Validate namespace
	if o.Namespace == "" {
		return errors.NewValidationError("namespace", "cannot be empty")
	}

	// Validate timeout
	if o.Timeout < time.Second || o.Timeout > 5*time.Minute {
		return errors.NewValidationError("timeout", "must be between 1s and 5m")
	}

	return nil
}

// SetDefaults sets default values for ListOptions
func (o *ListOptions) SetDefaults() {
	if o.Namespace == "" {
		o.Namespace = "default"
	}
	if o.Timeout == 0 {
		o.Timeout = 30 * time.Second
	}
}

// ClientConfig represents Kubernetes client configuration
type ClientConfig struct {
	KubeconfigPath string        `json:"kubeconfigPath"`
	Timeout        time.Duration `json:"timeout"`
}

// Validate validates ClientConfig
func (c *ClientConfig) Validate() error {
	// Validate timeout
	if c.Timeout < time.Second || c.Timeout > 5*time.Minute {
		return errors.NewValidationError("timeout", "must be between 1s and 5m")
	}

	return nil
}

// SetDefaults sets default values for ClientConfig
func (c *ClientConfig) SetDefaults() {
	if c.Timeout == 0 {
		c.Timeout = 30 * time.Second
	}
}
