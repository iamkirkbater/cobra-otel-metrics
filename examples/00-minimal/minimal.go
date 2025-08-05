package main

import (
	"context"
	"fmt"
	"log"

	metrics "github.com/iamkirkbater/cobra-otel-metrics"
	"github.com/spf13/cobra"
	http "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
)

func main() {
	var rootCmd = metrics.Command{
		Command: cobra.Command{
			Use:   "minimal",
			Short: "Example minimal CLI application with OpenTelemetry metrics",
			Run:   runExample,
		},
	}

	httpExporter, err := http.New(
		context.Background(),
		http.WithInsecure(),
	)
	if err != nil {
		log.Fatal(err)
	}

	rootCmd.SetupMetrics(
		metrics.WithExporter(httpExporter),
	)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func runExample(cmd *cobra.Command, args []string) {
	fmt.Println("Hello World")
}
