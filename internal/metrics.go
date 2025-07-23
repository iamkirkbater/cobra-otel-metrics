package internal

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	metricSdk "go.opentelemetry.io/otel/sdk/metric"
	resourceSdk "go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

type Config struct {
	ServiceName string
	Exporters   []metricSdk.Exporter
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

func NewConfig(cmd *cobra.Command, opts ...Option) (*Config, error) {
	config := &Config{
		ServiceName: GetRootCmdName(cmd),
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
