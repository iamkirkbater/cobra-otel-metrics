package internal

import (
	"context"
	"errors"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	metricSdk "go.opentelemetry.io/otel/sdk/metric"
	resourceSdk "go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

type Config struct {
	PrintToStdout bool
	ServiceName   string
	Exporters     []metricSdk.Exporter
}

type MetricsProvider struct {
	config   *Config
	Provider *metricSdk.MeterProvider
	Meter    metric.Meter
}

var (
	Reader    = metricSdk.NewManualReader()
	Exporters []metricSdk.Exporter
)

func NewConfig(opts ...Option) (*Config, error) {
	config := &Config{
		PrintToStdout: false,
	}

	for _, opt := range opts {
		if err := opt.apply(config); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	Exporters = config.Exporters

	return config, nil
}

func (c *Config) validate() error {
	var serviceNameEmptyErr error

	if c.ServiceName == "" {
		serviceNameEmptyErr = errors.New("Service Name cannot be empty")
	}

	return errors.Join(serviceNameEmptyErr)
}

func NewMetricsProvider(ctx context.Context, config *Config) (*MetricsProvider, error) {
	// Create resource
	res, err := resourceSdk.New(ctx,
		resourceSdk.WithAttributes(
			semconv.ServiceName(config.ServiceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// if config.Exporters.GRPC != nil {
	// 	grpcConfig := config.Exporters.GRPC
	// 	// Create OTLP exporter
	// 	var dialOpts []grpc.DialOption
	// 	if grpcConfig.AllowInsecure {
	// 		dialOpts = append(dialOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	// 	}

	// 	exporter, err := otlpmetricgrpc.New(ctx,
	// 		otlpmetricgrpc.WithEndpoint(config.Exporters.GRPC.CollectorURL),
	// 		otlpmetricgrpc.WithDialOption(dialOpts...),
	// 	)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("failed to create OTLP exporter: %w", err)
	// 	}

	// 	Exporters = append(Exporters, exporter)
	// }

	provider := metricSdk.NewMeterProvider(
		metricSdk.WithResource(res),
		metricSdk.WithReader(Reader),
	)

	// Set global meter provider
	otel.SetMeterProvider(provider)

	// Create meter
	meter := provider.Meter("cobra-otel-metrics")

	return &MetricsProvider{
		config:   config,
		Provider: provider,
		Meter:    meter,
	}, nil
}

func (mp *MetricsProvider) GetMeter() metric.Meter {
	return mp.Meter
}

func (mp *MetricsProvider) Shutdown(ctx context.Context) error {
	return mp.Provider.Shutdown(ctx)
}

func (mp *MetricsProvider) Config() *Config {
	return mp.config
}
