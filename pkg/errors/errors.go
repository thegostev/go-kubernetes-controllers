package errors

import "fmt"

// Error types for different failure scenarios
type (
	// ConfigError represents configuration-related errors
	ConfigError struct {
		Message string
		Err     error
	}

	// ConnectionError represents cluster connection errors
	ConnectionError struct {
		Message string
		Err     error
	}

	// ValidationError represents input validation errors
	ValidationError struct {
		Field   string
		Message string
	}

	// WatchError represents watch stream failures
	WatchError struct {
		Message string
		Err     error
	}

	// CacheError represents cache corruption or storage issues
	CacheError struct {
		Message string
		Err     error
	}

	// ResyncError represents resync operation failures
	ResyncError struct {
		Message string
		Err     error
	}
)

// Error implementations
func (e *ConfigError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("configuration error: %s: %v", e.Message, e.Err)
	}
	return fmt.Sprintf("configuration error: %s", e.Message)
}

func (e *ConnectionError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("connection error: %s: %v", e.Message, e.Err)
	}
	return fmt.Sprintf("connection error: %s", e.Message)
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error for field '%s': %s", e.Field, e.Message)
}

func (e *WatchError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("watch error: %s: %v", e.Message, e.Err)
	}
	return fmt.Sprintf("watch error: %s", e.Message)
}

func (e *CacheError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("cache error: %s: %v", e.Message, e.Err)
	}
	return fmt.Sprintf("cache error: %s", e.Message)
}

func (e *ResyncError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("resync error: %s: %v", e.Message, e.Err)
	}
	return fmt.Sprintf("resync error: %s", e.Message)
}

// Helper functions to create errors
func NewConfigError(message string, err error) *ConfigError {
	return &ConfigError{Message: message, Err: err}
}

func NewConnectionError(message string, err error) *ConnectionError {
	return &ConnectionError{Message: message, Err: err}
}

func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{Field: field, Message: message}
}

func NewWatchError(message string, err error) *WatchError {
	return &WatchError{Message: message, Err: err}
}

func NewCacheError(message string, err error) *CacheError {
	return &CacheError{Message: message, Err: err}
}

func NewResyncError(message string, err error) *ResyncError {
	return &ResyncError{Message: message, Err: err}
}
