package internal

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"go.opentelemetry.io/otel/attribute"
)

// Returns the full subcommand path without the root command
// however if the root command is the only command called we
// return "root"
func ParseCmdName(cmd *cobra.Command) string {
	cmdString := cmd.CommandPath()
	stringArr := strings.Split(cmdString, " ")
	if len(stringArr) == 1 {
		return "root"
	}
	stringArr = stringArr[1:]
	return strings.Join(stringArr, "-")
}

// Returns the root command name
func GetRootCmdName(cmd *cobra.Command) string {
	cmdString := cmd.CommandPath()
	stringArr := strings.Split(cmdString, " ")
	return stringArr[0]
}

// Loops through all provided flags and converts ONLY the
// name of the flag (NOT THE VALUE) to an attribute for
// additional metric collection
func ParseCmdFlagsToAttributes(cmd *cobra.Command) attribute.Set {
	flags := []attribute.KeyValue{}
	parseFlag := func(f *pflag.Flag) {
		flags = append(flags, attribute.String(f.Name, "1"))
	}
	cmd.Flags().Visit(parseFlag)

	flagsAttrSet, _ := attribute.NewSetWithFiltered(flags, nil)
	return flagsAttrSet
}
