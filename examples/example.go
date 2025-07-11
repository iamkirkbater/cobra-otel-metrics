package main

import (
	"context"
	"fmt"
	"log"
	"time"

	metrics "github.com/iamkirkbater/cobra-otel-metrics"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/metric"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "example",
		Short: "Example CLI application with OpenTelemetry metrics",
		RunE:  runExample,
	}

	stdoutExp, err := stdoutmetric.New()
	if err != nil {
		log.Fatal("Failed to setup stdoutExporter:", err)
	}

	// Setup metrics for the Cobra command
	err = metrics.SetupCobraMetrics(
		rootCmd,
		metrics.WithServiceName("example-service"), // Service name
		metrics.WithExporter(stdoutExp),
	)
	if err != nil {
		log.Fatal("Failed to setup metrics:", err)
	}

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func runExample(cmd *cobra.Command, args []string) error {
	// Get the meter for creating instruments
	meter := metrics.GetMeter()

	// Create a counter metric
	counter, err := meter.Int64Counter(
		"example_requests_total",
		metric.WithDescription("Total number of requests"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return fmt.Errorf("failed to create counter: %w", err)
	}

	// Create a histogram metric
	histogram, err := meter.Float64Histogram(
		"example_request_duration_seconds",
		metric.WithDescription("Request duration in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return fmt.Errorf("failed to create histogram: %w", err)
	}

	// Create a gauge metric
	gauge, err := meter.Int64UpDownCounter(
		"example_active_connections",
		metric.WithDescription("Number of active connections"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return fmt.Errorf("failed to create gauge: %w", err)
	}

	// Simulate some work and record metrics
	fmt.Println("Starting example work...")

	for i := 0; i < 10; i++ {
		start := time.Now()

		// Simulate some work
		time.Sleep(100 * time.Millisecond)

		// Record metrics
		counter.Add(context.Background(), 1, metric.WithAttributes(
			attribute.String("method", "GET"),
			attribute.String("status", "200"),
		))

		duration := time.Since(start).Seconds()
		histogram.Record(context.Background(), duration, metric.WithAttributes(
			attribute.String("method", "GET"),
		))

		gauge.Add(context.Background(), 1)

		fmt.Printf("Processed request %d (duration: %.3fs)\n", i+1, duration)
	}

	// Simulate decreasing active connections
	for i := 0; i < 10; i++ {
		gauge.Add(context.Background(), -1)
	}

	fmt.Println("Example work completed. Metrics have been recorded.")
	return nil
}
