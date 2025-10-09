// Package main provides the random command for mcpipboy
package main

import (
	"fmt"
	"io"
	"os"

	"github.com/kluzzebass/mcpipboy/internal/tools"
	"github.com/spf13/cobra"
)

var randomCmd = &cobra.Command{
	Use:   "random [flags]",
	Short: "Generate random numbers with various types and distributions",
	Long: `Random number generator provides comprehensive random number functionality including:

- Integer generation with min/max range
- Float generation with min/max range and precision control
- Boolean generation
- Batch generation with count parameter

Examples:
  mcpipboy random --type integer --min 1 --max 100
  mcpipboy random --type float --min 0 --max 1 --precision 3 --count 5
  mcpipboy random --type boolean --count 10`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runRandom(cmd, args, os.Stdout)
	},
}

var (
	randomType      string
	randomCount     int
	randomMin       float64
	randomMax       float64
	randomPrecision int
)

func init() {
	// Set group ID for random command
	randomCmd.GroupID = "tools"

	randomCmd.Flags().StringVar(&randomType, "type", "integer", "Type of random number: integer, float, boolean")
	randomCmd.Flags().IntVar(&randomCount, "count", 1, "Number of random values to generate (1-1000)")
	randomCmd.Flags().Float64Var(&randomMin, "min", 0, "Minimum value (for integer/float types)")
	randomCmd.Flags().Float64Var(&randomMax, "max", 100, "Maximum value (for integer/float types)")
	randomCmd.Flags().IntVar(&randomPrecision, "precision", 2, "Decimal places for float values (0-10)")

	// Add command to root
	rootCmd.AddCommand(randomCmd)
}

func runRandom(cmd *cobra.Command, args []string, out io.Writer) error {
	// Build parameters map
	params := make(map[string]interface{})

	// Add parameters if they were set
	if randomType != "" {
		params["type"] = randomType
	}
	params["count"] = float64(randomCount)
	if randomMin != 0 || randomMax != 100 {
		params["min"] = randomMin
		params["max"] = randomMax
	}
	if randomPrecision != 2 {
		params["precision"] = float64(randomPrecision)
	}

	// Create and execute the random tool
	tool := tools.NewRandomTool()

	// Validate parameters
	if err := tool.ValidateParams(params); err != nil {
		return fmt.Errorf("parameter validation failed: %v", err)
	}

	// Execute the tool
	result, err := tool.Execute(params)
	if err != nil {
		return fmt.Errorf("random tool execution failed: %v", err)
	}

	// Print the result
	fmt.Fprintln(out, result)
	return nil
}
