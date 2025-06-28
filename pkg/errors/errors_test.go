package errors

import (
	"testing"
)

func TestNewConfigError(t *testing.T) {
	err := NewConfigError("test message", nil)
	if err.Error() != "configuration error: test message" {
		t.Errorf("expected 'configuration error: test message', got '%s'", err.Error())
	}
}

func TestNewConnectionError(t *testing.T) {
	err := NewConnectionError("test message", nil)
	if err.Error() != "connection error: test message" {
		t.Errorf("expected 'connection error: test message', got '%s'", err.Error())
	}
}

func TestNewValidationError(t *testing.T) {
	err := NewValidationError("field", "test message")
	if err.Error() != "validation error for field 'field': test message" {
		t.Errorf("expected 'validation error for field 'field': test message', got '%s'", err.Error())
	}
}
