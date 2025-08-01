package internal

import (
	"testing"

	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel/attribute"
)

func TestParseCmdName(t *testing.T) {
	tests := []struct {
		name     string
		cmd      *cobra.Command
		expected string
	}{
		{
			name:     "root command only",
			cmd:      &cobra.Command{Use: "myapp"},
			expected: "root",
		},
		{
			name: "single subcommand",
			cmd: func() *cobra.Command {
				root := &cobra.Command{Use: "myapp"}
				sub := &cobra.Command{Use: "subcmd"}
				root.AddCommand(sub)
				return sub
			}(),
			expected: "subcmd",
		},
		{
			name: "nested subcommands",
			cmd: func() *cobra.Command {
				root := &cobra.Command{Use: "myapp"}
				sub1 := &cobra.Command{Use: "sub1"}
				sub2 := &cobra.Command{Use: "sub2"}
				root.AddCommand(sub1)
				sub1.AddCommand(sub2)
				return sub2
			}(),
			expected: "sub1-sub2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseCmdName(tt.cmd)
			if result != tt.expected {
				t.Errorf("ParseCmdName() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetRootCmdName(t *testing.T) {
	tests := []struct {
		name     string
		cmd      *cobra.Command
		expected string
	}{
		{
			name:     "root command only",
			cmd:      &cobra.Command{Use: "myapp"},
			expected: "myapp",
		},
		{
			name: "with subcommand",
			cmd: func() *cobra.Command {
				root := &cobra.Command{Use: "myapp"}
				sub := &cobra.Command{Use: "subcmd"}
				root.AddCommand(sub)
				return sub
			}(),
			expected: "myapp",
		},
		{
			name: "nested subcommands",
			cmd: func() *cobra.Command {
				root := &cobra.Command{Use: "myapp"}
				sub1 := &cobra.Command{Use: "sub1"}
				sub2 := &cobra.Command{Use: "sub2"}
				root.AddCommand(sub1)
				sub1.AddCommand(sub2)
				return sub2
			}(),
			expected: "myapp",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetRootCmdName(tt.cmd)
			if result != tt.expected {
				t.Errorf("GetRootCmdName() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestParseCmdFlagsToAttributes(t *testing.T) {
	tests := []struct {
		name         string
		setupCmd     func() *cobra.Command
		expectedKeys []string
	}{
		{
			name: "no flags",
			setupCmd: func() *cobra.Command {
				return &cobra.Command{Use: "test"}
			},
			expectedKeys: []string{},
		},
		{
			name: "single flag",
			setupCmd: func() *cobra.Command {
				cmd := &cobra.Command{Use: "test"}
				cmd.Flags().String("verbose", "", "verbose output")
				cmd.SetArgs([]string{"--verbose", "true"})
				cmd.Execute()
				return cmd
			},
			expectedKeys: []string{"verbose"},
		},
		{
			name: "multiple flags",
			setupCmd: func() *cobra.Command {
				cmd := &cobra.Command{Use: "test"}
				cmd.Flags().String("output", "", "output format")
				cmd.Flags().Bool("debug", false, "debug mode")
				cmd.SetArgs([]string{"--output", "json", "--debug"})
				cmd.Execute()
				return cmd
			},
			expectedKeys: []string{"output", "debug"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.setupCmd()
			result := ParseCmdFlagsToAttributes(cmd)

			if len(result) != len(tt.expectedKeys) {
				t.Errorf("ParseCmdFlagsToAttributes() returned %d attributes, want %d", len(result), len(tt.expectedKeys))
			}

			resultKeys := make(map[string]bool)
			for _, attr := range result {
				resultKeys[string(attr.Key)] = true
				if attr.Value.AsInt64() != 1 {
					t.Errorf("Expected attribute value to be 1, got %v", attr.Value.AsInt64())
				}
			}

			for _, expectedKey := range tt.expectedKeys {
				if !resultKeys[expectedKey] {
					t.Errorf("Expected attribute key %s not found in result", expectedKey)
				}
			}
		})
	}
}

func TestParseCmdFlagsToAttributesValues(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().String("format", "", "output format")
	cmd.SetArgs([]string{"--format", "json"})
	cmd.Execute()

	result := ParseCmdFlagsToAttributes(cmd)

	if len(result) != 1 {
		t.Fatalf("Expected 1 attribute, got %d", len(result))
	}

	attr := result[0]
	if string(attr.Key) != "format" {
		t.Errorf("Expected key 'format', got %s", string(attr.Key))
	}

	if attr.Value.Type() != attribute.INT64 {
		t.Errorf("Expected INT64 type, got %v", attr.Value.Type())
	}

	if attr.Value.AsInt64() != 1 {
		t.Errorf("Expected value 1, got %v", attr.Value.AsInt64())
	}
}