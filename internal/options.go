package internal

import metricSdk "go.opentelemetry.io/otel/sdk/metric"

type Option interface {
	applier
}

type option func(*Config) error

func (f option) apply(cfg *Config) error {
	return f(cfg)
}

var _ Option = (option)(nil)

type applier interface {
	apply(*Config) error
}

func WithExporter(exporter metricSdk.Exporter) Option {
	return option(func(cfg *Config) error {
		if cfg.Exporters == nil {
			cfg.Exporters = []metricSdk.Exporter{}
		}
		cfg.Exporters = append(cfg.Exporters, exporter)
		return nil
	})
}

// WithServiceName sets the service name for metrics
func WithServiceName(name string) Option {
	return option(func(cfg *Config) error {
		cfg.ServiceName = name
		return nil
	})
}
