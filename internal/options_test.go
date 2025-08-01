package internal

import (
	"testing"
)

func TestWithExporter(t *testing.T) {
	mockExporter := &mockExporter{}
	config := &Config{}

	opt := WithExporter(mockExporter)
	err := opt.apply(config)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(config.Exporters) != 1 {
		t.Errorf("Expected 1 exporter, got %d", len(config.Exporters))
	}

	if config.Exporters[0] != mockExporter {
		t.Error("Expected exporter to match")
	}
}

func TestWithExporterMultiple(t *testing.T) {
	exporter1 := &mockExporter{}
	exporter2 := &mockExporter{}
	config := &Config{}

	opt1 := WithExporter(exporter1)
	opt2 := WithExporter(exporter2)

	err := opt1.apply(config)
	if err != nil {
		t.Errorf("Unexpected error applying first exporter: %v", err)
	}

	err = opt2.apply(config)
	if err != nil {
		t.Errorf("Unexpected error applying second exporter: %v", err)
	}

	if len(config.Exporters) != 2 {
		t.Errorf("Expected 2 exporters, got %d", len(config.Exporters))
	}

	if config.Exporters[0] != exporter1 {
		t.Error("Expected first exporter to match")
	}

	if config.Exporters[1] != exporter2 {
		t.Error("Expected second exporter to match")
	}
}

func TestWithServiceName(t *testing.T) {
	tests := []struct {
		name        string
		serviceName string
		expected    string
	}{
		{
			name:        "valid service name",
			serviceName: "my-service",
			expected:    "my-service",
		},
		{
			name:        "empty service name",
			serviceName: "",
			expected:    "",
		},
		{
			name:        "service name with special characters",
			serviceName: "my-service-123_test",
			expected:    "my-service-123_test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{ServiceName: "initial"}
			opt := WithServiceName(tt.serviceName)

			err := opt.apply(config)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if config.ServiceName != tt.expected {
				t.Errorf("Expected service name %s, got %s", tt.expected, config.ServiceName)
			}
		})
	}
}

func TestOptionInterface(t *testing.T) {
	config := &Config{}

	var opt Option = WithServiceName("test-service")

	err := opt.apply(config)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if config.ServiceName != "test-service" {
		t.Errorf("Expected service name 'test-service', got %s", config.ServiceName)
	}
}

func TestOptionFunctionType(t *testing.T) {
	config := &Config{}

	var opt option = func(cfg *Config) error {
		cfg.ServiceName = "test-from-function"
		return nil
	}

	err := opt.apply(config)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if config.ServiceName != "test-from-function" {
		t.Errorf("Expected service name 'test-from-function', got %s", config.ServiceName)
	}
}

func TestApplierInterface(t *testing.T) {
	config := &Config{}

	var applier applier = WithServiceName("test-applier")

	err := applier.apply(config)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if config.ServiceName != "test-applier" {
		t.Errorf("Expected service name 'test-applier', got %s", config.ServiceName)
	}
}
