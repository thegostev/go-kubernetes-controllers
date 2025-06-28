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
