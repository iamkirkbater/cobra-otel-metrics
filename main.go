package metrics

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/iamkirkbater/cobra-otel-metrics/internal"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

type (
	Option = internal.Option
)

// Global metrics provider instance
var globalProvider *internal.MetricsProvider

// initialize sets up the metrics provider with the given options
func initialize(ctx context.Context, options ...Option) (*internal.MetricsProvider, error) {
	config, err := internal.NewConfig(options...)
	if err != nil {
		return nil, fmt.Errorf("failed to create config: %w", err)
	}

	provider, err := internal.NewMetricsProvider(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create metrics provider: %w", err)
	}

	globalProvider = provider

	return provider, nil
}

// GetMeter returns the meter for creating instruments
func GetMeter() metric.Meter {
	return globalProvider.GetMeter()
}

// Shutdown gracefully shuts down the metrics provider
func Shutdown(ctx context.Context) error {
	return globalProvider.Shutdown(ctx)
}

// SetupCobraMetrics is a convenience function to set up metrics for a Cobra command
// It initializes the metrics provider and sets up a cleanup handler
func SetupCobraMetrics(cmd *cobra.Command, options ...Option) error {
	// Initialize metrics provider
	ctx := context.Background()
	provider, err := initialize(ctx, options...)
	if err != nil {
		return fmt.Errorf("failed to initialize metrics: %w", err)
	}

	cleanupFunc := func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		collectedMetrics := &metricdata.ResourceMetrics{}
		internal.Reader.Collect(ctx, collectedMetrics)

		for _, exporter := range internal.Exporters {
			exporter.Export(ctx, collectedMetrics)
		}

		if err := provider.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(os.Stderr, "Error shutting down metrics provider: %v\n", err)
		}
	}

	// Set up cleanup on command completion
	// This is set up so that we can trap and still report metrics if
	// the command is exited prematurely
	originalPreRun := cmd.PreRunE
	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		// Set up signal handling for graceful shutdown
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			cleanupFunc()
			os.Exit(1)
		}()

		// Run original PreRunE if it exists
		if originalPreRun != nil {
			return originalPreRun(cmd, args)
		}
		return nil
	}

	// Set up cleanup on command completion
	// This is needed in addition to the trap running above
	originalPostRun := cmd.PostRunE
	cmd.PostRunE = func(cmd *cobra.Command, args []string) error {
		var originalPostRunErr error
		// Run original PostRunE if it exists
		if originalPostRun != nil {
			if err := originalPostRun(cmd, args); err != nil {
				originalPostRunErr = err
			}
		}

		cleanupFunc()

		// return the original error if it exists
		return originalPostRunErr
	}

	return nil
}

// Convenience functions for creating options
var (
	WithServiceName = internal.WithServiceName
	WithStdoutPrint = internal.WithStdoutPrint
	WithExporter    = internal.WithExporter
)
