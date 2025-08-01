package internal

import (
	"context"

	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

type mockExporter struct{}

func (m *mockExporter) Export(ctx context.Context, rm *metricdata.ResourceMetrics) error {
	return nil
}

func (m *mockExporter) ForceFlush(ctx context.Context) error {
	return nil
}

func (m *mockExporter) Shutdown(ctx context.Context) error {
	return nil
}

func (m *mockExporter) Temporality(kind metric.InstrumentKind) metricdata.Temporality {
	return metricdata.CumulativeTemporality
}

func (m *mockExporter) Aggregation(kind metric.InstrumentKind) metric.Aggregation {
	return metric.DefaultAggregationSelector(kind)
}
