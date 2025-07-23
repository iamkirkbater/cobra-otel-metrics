package main

import (
	"fmt"
	"log"

	metrics "github.com/iamkirkbater/cobra-otel-metrics"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
)

var testFlag bool
var stringFlag string

func main() {
	var rootCmd = metrics.Command{
		Command: cobra.Command{
			Use:   "sub",
			Short: "Example CLI application with OpenTelemetry metrics",
			Run: func(cmd *cobra.Command, args []string) {
				cmd.Usage()
			},
		},
	}

	stdoutExporter, _ := stdoutmetric.New()

	rootCmd.SetupMetrics(
		metrics.WithExporter(stdoutExporter),
	)

	subCmd := &cobra.Command{
		Use:  "my-subcommand",
		Args: cobra.NoArgs,
		Run:  runSubCmd,
	}
	flags := subCmd.Flags()
	flags.BoolVarP(&testFlag, "test", "t", false, "test flag for command")
	flags.StringVarP(&stringFlag, "string", "s", "", "a random string to pass")

	anotherChildCmd := &cobra.Command{
		Use:  "child",
		Args: cobra.NoArgs,
		Run:  runSubCmd,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			fmt.Println("In the child cmd persistent pre-run")
		},
	}

	grandChildCommand := &cobra.Command{
		Use:  "grandchild",
		Args: cobra.NoArgs,
		Run:  runSubCmd,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("In the grandchild command pre-run")
			return nil
		},
	}

	grandChildFlags := grandChildCommand.Flags()
	grandChildFlags.StringVarP(&stringFlag, "string", "s", "", "a random string to pass")
	anotherChildCmd.AddCommand(grandChildCommand)

	rootCmd.AddCommand(subCmd)
	rootCmd.AddCommand(anotherChildCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func runSubCmd(cmd *cobra.Command, args []string) {
	fmt.Println("Flag Values")
	fmt.Printf(" - testFlag: %t\n", testFlag)
	fmt.Printf(" - stringFlag: %s\n", stringFlag)
}
