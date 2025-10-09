package main

import (
	"fmt"

	"github.com/kluzzebass/mcpipboy/internal/tools"
	"github.com/spf13/cobra"
)

var (
	ean13Operation string
	ean13Input     string
	ean13Count     int
)

// ean13Cmd represents the ean13 command
var ean13Cmd = &cobra.Command{
	Use:   "ean13",
	Short: "Generate and validate European Article Numbers (EAN-13)",
	Long: `Generate and validate European Article Numbers with checksum validation.

Examples:
  # Validate an EAN-13
  mcpipboy ean13 --operation validate --input "1234567890128"

  # Validate an EAN-13 with formatting
  mcpipboy ean13 --operation validate --input "123-456-789-012-8"

  # Generate a single EAN-13
  mcpipboy ean13 --operation generate

  # Generate multiple EAN-13s
  mcpipboy ean13 --operation generate --count 5`,
	RunE: runEAN13,
}

func init() {
	rootCmd.AddCommand(ean13Cmd)

	// Add flags
	ean13Cmd.Flags().StringVar(&ean13Operation, "operation", "validate", "Operation to perform: validate or generate")
	ean13Cmd.Flags().StringVar(&ean13Input, "input", "", "EAN-13 number to validate (required for validate operation)")
	ean13Cmd.Flags().IntVar(&ean13Count, "count", 1, "Number of EAN-13s to generate (1-100, default: 1)")

	// Set command group
	ean13Cmd.GroupID = "tools"
}

func runEAN13(cmd *cobra.Command, args []string) error {
	// Create the EAN-13 tool
	tool := tools.NewEAN13Tool()

	// Build parameters
	params := make(map[string]interface{})

	// Add operation
	if ean13Operation != "" {
		params["operation"] = ean13Operation
	}

	// Add input for validation
	if ean13Input != "" {
		params["input"] = ean13Input
	}

	// Add count (always add, even if 0, so validation can handle it)
	params["count"] = float64(ean13Count)

	// Execute the tool
	result, err := tool.Execute(params)
	if err != nil {
		return fmt.Errorf("EAN-13 tool execution failed: %v", err)
	}

	// Handle the result based on operation
	if ean13Operation == "validate" {
		if resultMap, ok := result.(map[string]interface{}); ok {
			if valid, ok := resultMap["valid"].(bool); ok {
				if valid {
					fmt.Printf("Valid EAN-13: %s\n", resultMap["ean13"])
				} else {
					fmt.Printf("Invalid EAN-13: %s\n", resultMap["error"])
					if input, ok := resultMap["input"].(string); ok {
						fmt.Printf("   Input: %s\n", input)
					}
				}
			} else {
				fmt.Printf("EAN-13 validation result: %v\n", result)
			}
		} else {
			fmt.Printf("EAN-13 validation result: %v\n", result)
		}
	} else if ean13Operation == "generate" {
		if ean13Count == 1 {
			// Single EAN-13
			if ean13, ok := result.(string); ok {
				fmt.Printf("Generated EAN-13: %s\n", ean13)
			} else {
				fmt.Printf("Generated EAN-13: %v\n", result)
			}
		} else {
			// Multiple EAN-13s
			if ean13s, ok := result.([]string); ok {
				fmt.Printf("Generated %d EAN-13s:\n", len(ean13s))
				for i, ean13 := range ean13s {
					fmt.Printf("  %d. %s\n", i+1, ean13)
				}
			} else {
				fmt.Printf("Generated EAN-13s: %v\n", result)
			}
		}
	} else {
		fmt.Printf("EAN-13 result: %v\n", result)
	}

	return nil
}
