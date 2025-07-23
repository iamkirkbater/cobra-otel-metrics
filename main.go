package metrics

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/iamkirkbater/cobra-otel-metrics/internal"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

type Option = internal.Option

// Exposed Options
var (
	// Optional configuration to set the service name
	// this will default to your root command name
	WithServiceName = internal.WithServiceName

	// Can be passed multiple times to add many exporters
	WithExporter = internal.WithExporter
)

// Extend the cobra.Command struct here to allow drop-in replacement
// of the root command while also giving us a place to store context
// and other things we may need
type Command struct {
	cobra.Command

	ctx context.Context
}

// SetupMetrics is a convenience function to set up metrics for a Cobra command
// It initializes the metrics provider and sets up a cleanup handler
func (c *Command) SetupMetrics(options ...Option) error {
	ctx := context.Background()

	_, err := initialize(ctx, &c.Command, options...)
	if err != nil {
		return fmt.Errorf("failed to initialize metrics: %w", err)
	}

	wrapPreRun(&c.Command, true)
	wrapPostRun(&c.Command, true)

	c.ctx = ctx
	return nil
}

func wrapPostRun(cmd *cobra.Command, force bool) {
	originalPostRunE := cmd.PersistentPostRunE
	originalPostRun := cmd.PersistentPostRun
	if originalPostRunE != nil || originalPostRun != nil || force {
		cmd.PersistentPostRunE = func(cmd *cobra.Command, args []string) error {
			var originalPostRunErr error
			// Run original PostRunE if it exists
			if originalPostRunE != nil {
				if err := originalPostRunE(cmd, args); err != nil {
					originalPostRunErr = err
				}
			}

			// return the original error if it exists
			return originalPostRunErr
		}
	}
}

func wrapPreRun(cmd *cobra.Command, force bool) {
	originalPreRunE := cmd.PersistentPreRunE
	originalPreRun := cmd.PersistentPreRun
	if originalPreRunE != nil || originalPreRun != nil || force {
		cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
			// create the initial metric around the command called
			err := createInvocationMetric(cmd)
			if err != nil {
				return err
			}

			// Run original PreRunE if it exists
			if originalPreRunE != nil {
				return originalPreRunE(cmd, args)
			}
			return nil
		}
	}
}

func wrapSubCommandRunHooks(cmd *cobra.Command) {
	wrapPreRun(cmd, false)
	wrapPostRun(cmd, false)
	wrapHelpFunc(cmd)
	children := cmd.Commands()
	for i := range children {
		wrapSubCommandRunHooks(children[i])
	}
}

func wrapHelpFunc(cmd *cobra.Command) {
	oldHelpFunc := cmd.HelpFunc()
	cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		createInvocationMetric(cmd)
		oldHelpFunc(cmd, args)
	})
}

func (c *Command) Execute() error {
	// catch trap signals and send metrics if we can
	c.trap()

	for _, command := range c.Command.Commands() {
		wrapSubCommandRunHooks(command)
	}

	// Run the command
	err := c.Command.Execute()

	// push metrics even if the command wasn't successful
	c.cleanup()

	return err
}

func (c *Command) trap() {
	// Trap Command Cancellations
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ch
		c.cleanup()
		os.Exit(1)
	}()

}

func (c *Command) cleanup() {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	collectedMetrics := &metricdata.ResourceMetrics{}
	internal.Reader.Collect(c.ctx, collectedMetrics)

	for _, exporter := range internal.Exporters {
		exporter.Export(c.ctx, collectedMetrics)
	}

	if err := shutdown(shutdownCtx); err != nil {
		fmt.Fprintf(os.Stderr, "Error shutting down metrics provider: %v\n", err)
	}
}

func createInvocationMetric(cmd *cobra.Command) error {
	// Create a counter metric
	counter, err := GetMeter().Int64Counter(
		internal.ParseCmdName(cmd),
		metric.WithDescription("Command Invocation"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return fmt.Errorf("failed to create counter: %w", err)
	}

	attributes := internal.ParseCmdFlagsToAttributes(cmd)

	isTTY := false
	if isatty.IsTerminal(os.Stdin.Fd()) {
		isTTY = true
	}
	attributes = append(attributes, attribute.Bool("tty", isTTY))

	attributeSet, _ := attribute.NewSetWithFiltered(attributes, nil)

	counter.Add(context.Background(), 1, metric.WithAttributeSet(attributeSet))

	return nil
}

// Global metrics provider instance
var globalProvider *internal.MetricsProvider

// GetMeter returns the meter for creating instruments
func GetMeter() metric.Meter {
	return globalProvider.GetMeter()
}

// shutdown gracefully shuts down the metrics provider
func shutdown(ctx context.Context) error {
	return globalProvider.Shutdown(ctx)
}

// initialize sets up the metrics provider with the given options
func initialize(ctx context.Context, cmd *cobra.Command, options ...Option) (*internal.MetricsProvider, error) {
	config, err := internal.NewConfig(cmd, options...)
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
