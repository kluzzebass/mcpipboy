package main

import (
	"fmt"
	"io"
	"os"

	"github.com/kluzzebass/mcpipboy/internal/tools"
	"github.com/spf13/cobra"
)

var (
	imoOperation string
	imoInput     string
	imoCount     int
)

// imoCmd represents the imo command
var imoCmd = &cobra.Command{
	Use:   "imo",
	Short: "Generate and validate IMO numbers",
	Long: `Generate and validate International Maritime Organization (IMO) numbers.

IMO numbers are 7-digit numbers with a check digit calculated using a weighted sum algorithm.
The check digit is calculated as: (7×d1 + 6×d2 + 5×d3 + 4×d4 + 3×d5 + 2×d6) mod 10

Examples:
  mcpipboy imo --operation validate --input "1234567"
  mcpipboy imo --operation generate --count 5
  mcpipboy imo --operation generate`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runIMO(cmd, args, os.Stdout)
	},
}

func init() {
	rootCmd.AddCommand(imoCmd)
	imoCmd.GroupID = "tools"

	// Add flags
	imoCmd.Flags().StringVar(&imoOperation, "operation", "validate", "Operation to perform: 'validate' or 'generate'")
	imoCmd.Flags().StringVar(&imoInput, "input", "", "IMO number to validate (required for validation)")
	imoCmd.Flags().IntVar(&imoCount, "count", 1, "Number of IMO numbers to generate (max: 100)")

	// Mark input as required only for validation
	// We'll handle this in the runIMO function
}

func runIMO(cmd *cobra.Command, args []string, out io.Writer) error {
	// Create the IMO tool
	tool := tools.NewIMOTool()

	// Build parameters
	params := map[string]interface{}{
		"operation": imoOperation,
		"count":     imoCount,
	}

	// Add input for validation
	if imoOperation == "validate" {
		if imoInput == "" {
			return fmt.Errorf("input is required for validation operation")
		}
		params["input"] = imoInput
	}

	// Validate parameters
	if err := tool.ValidateParams(params); err != nil {
		return fmt.Errorf("invalid parameters: %v", err)
	}

	// Execute the tool
	result, err := tool.Execute(params)
	if err != nil {
		return fmt.Errorf("IMO tool execution failed: %v", err)
	}

	// Handle different result types
	switch v := result.(type) {
	case string:
		// Single IMO number
		fmt.Fprintln(out, v)
	case []string:
		// Multiple IMO numbers
		for _, imo := range v {
			fmt.Fprintln(out, imo)
		}
	case map[string]interface{}:
		// Validation result
		if valid, ok := v["valid"].(bool); ok {
			if valid {
				fmt.Fprintf(out, "Valid IMO: %s\n", v["imo"])
			} else {
				fmt.Fprintf(out, "Invalid IMO: %s\n", v["error"])
				if input, ok := v["input"].(string); ok {
					fmt.Fprintf(out, "   Input: %s\n", input)
				}
			}
		} else {
			// Fallback for unexpected result format
			fmt.Fprintf(out, "%+v\n", v)
		}
	default:
		return fmt.Errorf("unexpected result type: %T", result)
	}

	return nil
}
