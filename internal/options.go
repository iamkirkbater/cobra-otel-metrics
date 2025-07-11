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

// func WithGrpcExporter(opts ...Option) Option {
// 	grpcCfg := &ExporterConfig{}
// 	return option(func(cfg *Config) error {
// 		cfg.Exporters.GRPC = grpcCfg
// 		return nil
// 	})
// }

func WithExporter(exporter metricSdk.Exporter) Option {
	return option(func(cfg *Config) error {
		if cfg.Exporters == nil {
			cfg.Exporters = []metricSdk.Exporter{}
		}
		cfg.Exporters = append(cfg.Exporters, exporter)
		return nil
	})
}

// // WithCollectorURL sets the URL for the metrics collector
// func WithCollectorURL(url string) Option {
// 	return option(func(cfg *Config) error {
// 		cfg.CollectorURL = url
// 		return nil
// 	})
// }
//
// // WithInsecure allows insecure connections to the collector
// func WithInsecure(insecure bool) Option {
// 	return option(func(cfg *Config) error {
// 		cfg.AllowInsecure = insecure
// 		return nil
// 	})
// }

// WithServiceName sets the service name for metrics
func WithServiceName(name string) Option {
	return option(func(cfg *Config) error {
		cfg.ServiceName = name
		return nil
	})
}

// WithStdoutPrint enables printing metrics to stdout for debugging
func WithStdoutPrint(enabled bool) Option {
	return option(func(cfg *Config) error {
		cfg.PrintToStdout = enabled
		return nil
	})
}
