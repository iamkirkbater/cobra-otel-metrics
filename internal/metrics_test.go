package internal

import (
	"context"
	"errors"
	"testing"

	"github.com/spf13/cobra"
)

var mp *MetricsProvider
var config = &Config{
	ServiceName: "test-service",
}

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name        string
		cmd         *cobra.Command
		opts        []Option
		expectError bool
		expectedSvc string
	}{
		{
			name:        "valid config with default service name",
			cmd:         &cobra.Command{Use: "myapp"},
			opts:        []Option{},
			expectError: false,
			expectedSvc: "myapp",
		},
		{
			name: "valid config with custom service name",
			cmd:  &cobra.Command{Use: "myapp"},
			opts: []Option{
				WithServiceName("custom-service"),
			},
			expectError: false,
			expectedSvc: "custom-service",
		},
		{
			name: "invalid config with empty service name",
			cmd:  &cobra.Command{Use: "myapp"},
			opts: []Option{
				WithServiceName(""),
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := NewConfig(tt.cmd, tt.opts...)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if config.ServiceName != tt.expectedSvc {
				t.Errorf("Expected service name %s, got %s", tt.expectedSvc, config.ServiceName)
			}
		})
	}
}

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		expectError bool
	}{
		{
			name: "valid config",
			config: &Config{
				ServiceName: "test-service",
			},
			expectError: false,
		},
		{
			name: "empty service name",
			config: &Config{
				ServiceName: "",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.validate()

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

// These two functions MUST be called in order.
// I hate it, but there's a global registration it seems that
// happens when you create a metrics provider. If we call this
// test and then TestmetricsProviderMethods later we get an error
// "did not register manual reader - duplicate reader registration"
// and while that's just a warning, the subsequent mp.Shutdown test
// in TestMetricsProviderMethods below will fail if we run
// defer mp.Shutdown(ctx) here like we should.
func TestNewMetricsProvider(t *testing.T) {
	ctx := context.Background()
	t.Run("valid config", func(t *testing.T) {
		var err error
		mp, err = NewMetricsProvider(ctx, config)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
			return
		}

		if mp == nil {
			t.Error("Expected MetricsProvider but got nil")
			return
		}

		if mp.Provider == nil {
			t.Error("Expected Provider but got nil")
		}

		if mp.Meter == nil {
			t.Error("Expected Meter but got nil")
		}

		if mp.config != config {
			t.Error("Expected config to match")
		}
	})
}

// This function MUST be called after TestNewMetricsProvider
// and also needs to be called AFTER we do any metrics testing
// if that's ever added, as the last method here runs shutdown
func TestMetricsProviderMethods(t *testing.T) {
	ctx := context.Background()

	t.Run("GetMeter", func(t *testing.T) {
		meter := mp.GetMeter()
		if meter == nil {
			t.Error("Expected meter but got nil")
		}
	})

	t.Run("Config", func(t *testing.T) {
		returnedConfig := mp.Config()
		if returnedConfig != config {
			t.Error("Expected config to match original")
		}
	})

	t.Run("Shutdown", func(t *testing.T) {
		err := mp.Shutdown(ctx)
		if err != nil {
			t.Errorf("Unexpected error during shutdown: %v", err)
		}
	})
}

func TestNewConfigWithExporter(t *testing.T) {
	cmd := &cobra.Command{Use: "myapp"}

	mockExporter := &mockExporter{}

	config, err := NewConfig(cmd, WithExporter(mockExporter))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(config.Exporters) != 1 {
		t.Errorf("Expected 1 exporter, got %d", len(config.Exporters))
	}

	if config.Exporters[0] != mockExporter {
		t.Error("Expected exporter to match")
	}
}

func TestConfigValidateMultipleErrors(t *testing.T) {
	config := &Config{
		ServiceName: "",
	}

	err := config.validate()
	if err == nil {
		t.Error("Expected error but got nil")
		return
	}

	if !errors.Is(err, ErrServiceNameEmpty) {
		t.Error("Expected service name error")
	}
}
