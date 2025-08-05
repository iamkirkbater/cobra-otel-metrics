# Cobra OpenTelemetry Metrics

A Go package that provides easy integration of OpenTelemetry metrics with Cobra CLI applications.

This package is currently very much in an alpha state. No guarantees on what may change in the future are made.

The goal is to give Cobra CLI maintainers an extremely simple way to start collecting usage metrics, while at the same time allowing experienced users the abilty to highly configure the options they need.

## Features

- **Easy Integration**: Extremely simple setup for new and existing Cobra CLI applications
- **Signal Handling**: Proper signal handling for graceful shutdowns
- **Call counter with flags**: By default, we will create a metric to count calls passing the command executed as well as the flags passed. We will explicitly NOT pass values to the flags.
- **GDPR/Privacy/Opt-Out**: This will ship with an opt-out option for end users. Users will be prompted to opt-out and that configuration will be saved.
    - Non-Interactive sessions will not be prompted and metrics will be shipped.
    - Non-Interactive sessions can be opted-out by either creating an opt-out file at the default file path or by 

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
