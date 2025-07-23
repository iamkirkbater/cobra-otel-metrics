package main

import (
	"fmt"
	"log"

	metrics "github.com/iamkirkbater/cobra-otel-metrics"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
)

func main() {
	var rootCmd = metrics.Command{
		Command: cobra.Command{
			Use:   "minimal",
			Short: "Example minimal CLI application with OpenTelemetry metrics",
			Run:   runExample,
		},
	}

	stdoutExporter, _ := stdoutmetric.New()

	rootCmd.SetupMetrics(
		metrics.WithExporter(stdoutExporter),
	)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func runExample(cmd *cobra.Command, args []string) {
	fmt.Println("Hello World")
}
