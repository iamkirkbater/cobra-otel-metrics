# Cobra OpenTelemetry Metrics

A Go package that provides easy integration of OpenTelemetry metrics with Cobra CLI applications using the Functional Options pattern.

This package is currently very much in an alpha state. No guarantees on what may change in the future are made.

## Features

- **Easy Integration**: Extremely simple setup for new and existing Cobra CLI applications
- **Signal Handling**: Proper signal handling for graceful shutdowns

## Installation

```bash
go get github.com/iamkirkbater/cobra-otel-metrics
```

## Quick Start

### Basic Usage

Simply create one or many otel exporters, wrap your root command in a `metrics.Command` wrapper, and then run the metrics setup function passing in your exporter(s).

See the `examples` in the [examples](/examples) directory for usage examples.

## Signal Handling

When using `SetupCobraMetrics`, the package automatically handles SIGINT and SIGTERM signals for graceful shutdown, ensuring that metrics are properly flushed before the application exits.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request or open an Issue.
