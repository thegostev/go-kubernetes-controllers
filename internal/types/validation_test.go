package types

import (
	"testing"
	"time"
)

func TestListOptionsValidate(t *testing.T) {
	tests := []struct {
		name    string
		options ListOptions
		wantErr bool
	}{
		{
			name: "valid options",
			options: ListOptions{
				Namespace: "default",
				Timeout:   30 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "empty namespace",
			options: ListOptions{
				Namespace: "",
				Timeout:   30 * time.Second,
			},
			wantErr: true,
		},
		{
			name: "invalid timeout",
			options: ListOptions{
				Namespace: "default",
				Timeout:   6 * time.Minute,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.options.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ListOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestListOptionsSetDefaults(t *testing.T) {
	options := &ListOptions{}
	options.SetDefaults()

	if options.Namespace != "default" {
		t.Errorf("expected namespace to be 'default', got '%s'", options.Namespace)
	}

	if options.Timeout != 30*time.Second {
		t.Errorf("expected timeout to be 30s, got %v", options.Timeout)
	}
}
