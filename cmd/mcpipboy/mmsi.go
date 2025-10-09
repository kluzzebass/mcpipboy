// Package main provides the mmsi command for mcpipboy
package main

import (
	"fmt"

	"github.com/kluzzebass/mcpipboy/internal/tools"
	"github.com/spf13/cobra"
)

var (
	mmsiOperation   string
	mmsiInput       string
	mmsiType        string
	mmsiCountryCode string
	mmsiCount       int
)

// mmsiCmd represents the mmsi command
var mmsiCmd = &cobra.Command{
	Use:   "mmsi",
	Short: "Generate and validate MMSI numbers",
	Long: `Generate and validate Maritime Mobile Service Identity (MMSI) numbers.

MMSI numbers are 9-digit identifiers used for maritime communication.
They can be validated for format and Maritime Identification Digit (MID), or generated with optional country code.

Examples:
  # Validate an MMSI number
  mcpipboy mmsi --operation validate --input "366123456"

  # Generate a regular ship MMSI for US
  mcpipboy mmsi --operation generate --type ship --country-code "US"

  # Generate a SAR aircraft MMSI
  mcpipboy mmsi --operation generate --type sar-aircraft

  # Generate US Federal MMSI
  mcpipboy mmsi --operation generate --type us-federal

  # Generate multiple coast stations for UK
  mcpipboy mmsi --operation generate --type coast --country-code "GB" --count 5`,
	RunE: runMMSI,
}

func init() {
	rootCmd.AddCommand(mmsiCmd)
	mmsiCmd.GroupID = "tools"

	// Add flags
	mmsiCmd.Flags().StringVar(&mmsiOperation, "operation", "validate", "Operation to perform: 'validate' or 'generate'")
	mmsiCmd.Flags().StringVar(&mmsiInput, "input", "", "MMSI number to validate (required for validation operation)")
	mmsiCmd.Flags().StringVar(&mmsiType, "type", "", "MMSI type to generate (optional for generation)")
	mmsiCmd.Flags().StringVar(&mmsiCountryCode, "country-code", "US", "Country code for generation (e.g., US, GB, DE, FR, etc.)")
	mmsiCmd.Flags().IntVar(&mmsiCount, "count", 1, "Number of MMSI numbers to generate (max: 100)")

	// Mark required flags
	mmsiCmd.MarkFlagRequired("input")
}

func runMMSI(cmd *cobra.Command, args []string) error {
	// Build parameters map
	params := make(map[string]interface{})

	// Add parameters if they were set
	if mmsiOperation != "" {
		params["operation"] = mmsiOperation
	}
	if mmsiInput != "" {
		params["input"] = mmsiInput
	}
	if mmsiType != "" {
		params["type"] = mmsiType
	}
	if mmsiCountryCode != "" {
		params["country-code"] = mmsiCountryCode
	}
	if mmsiCount > 0 {
		params["count"] = mmsiCount
	}

	// Create and execute the MMSI tool
	tool := tools.NewMMSITool()

	// Validate parameters
	if err := tool.ValidateParams(params); err != nil {
		return fmt.Errorf("parameter validation failed: %v", err)
	}

	// Execute the tool
	result, err := tool.Execute(params)
	if err != nil {
		return fmt.Errorf("MMSI tool execution failed: %v", err)
	}

	// Display results
	switch v := result.(type) {
	case string:
		// Single MMSI number
		fmt.Println(v)
	case []string:
		// Multiple MMSI numbers
		for _, mmsi := range v {
			fmt.Println(mmsi)
		}
	case map[string]interface{}:
		// Validation result
		if valid, ok := v["valid"].(bool); ok {
			if valid {
				fmt.Printf("Valid MMSI: %s\n", v["mmsi"])
				if countryName, ok := v["country_name"].(string); ok {
					fmt.Printf("   Country: %s\n", countryName)
				}
			} else {
				fmt.Printf("Invalid MMSI: %s\n", v["error"])
				if input, ok := v["input"].(string); ok {
					fmt.Printf("   Input: %s\n", input)
				}
			}
		}
	default:
		fmt.Printf("Result: %v\n", result)
	}

	return nil
}
