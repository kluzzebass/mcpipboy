package main

import (
	"fmt"

	"github.com/kluzzebass/mcpipboy/internal/tools"
	"github.com/spf13/cobra"
)

var (
	isbnOperation string
	isbnInput     string
	isbnFormat    string
	isbnCount     int
)

// isbnCmd represents the isbn command
var isbnCmd = &cobra.Command{
	Use:   "isbn",
	Short: "Generate and validate International Standard Book Numbers (ISBN-10 and ISBN-13)",
	Long: `Generate and validate International Standard Book Numbers with format support.

Examples:
  # Validate an ISBN-10
  mcpipboy isbn --operation validate --input "0-123456-78-9" --format "isbn10"

  # Validate an ISBN-13
  mcpipboy isbn --operation validate --input "978-0-123456-78-9" --format "isbn13"

  # Auto-detect format
  mcpipboy isbn --operation validate --input "9780123456789"

  # Generate a single ISBN-13 (default)
  mcpipboy isbn --operation generate

  # Generate multiple ISBN-10s
  mcpipboy isbn --operation generate --format "isbn10" --count 3

  # Generate multiple ISBN-13s
  mcpipboy isbn --operation generate --format "isbn13" --count 5`,
	RunE: runISBN,
}

func init() {
	rootCmd.AddCommand(isbnCmd)

	// Add flags
	isbnCmd.Flags().StringVar(&isbnOperation, "operation", "validate", "Operation to perform: validate or generate")
	isbnCmd.Flags().StringVar(&isbnInput, "input", "", "ISBN number to validate (required for validate operation)")
	isbnCmd.Flags().StringVar(&isbnFormat, "format", "", "ISBN format: isbn10, isbn13, or auto (default: auto for validation, isbn13 for generation)")
	isbnCmd.Flags().IntVar(&isbnCount, "count", 1, "Number of ISBNs to generate (1-100, default: 1)")

	// Set command group
	isbnCmd.GroupID = "tools"
}

func runISBN(cmd *cobra.Command, args []string) error {
	// Create the ISBN tool
	tool := tools.NewISBNTool()

	// Build parameters
	params := make(map[string]interface{})

	// Add operation
	if isbnOperation != "" {
		params["operation"] = isbnOperation
	}

	// Add input for validation
	if isbnInput != "" {
		params["input"] = isbnInput
	}

	// Add format
	if isbnFormat != "" {
		params["format"] = isbnFormat
	}

	// Add count (always add, even if 0, so validation can handle it)
	params["count"] = float64(isbnCount)

	// Execute the tool
	result, err := tool.Execute(params)
	if err != nil {
		return fmt.Errorf("ISBN tool execution failed: %v", err)
	}

	// Handle the result based on operation
	if isbnOperation == "validate" {
		if resultMap, ok := result.(map[string]interface{}); ok {
			if valid, ok := resultMap["valid"].(bool); ok {
				if valid {
					fmt.Printf("Valid ISBN: %s\n", resultMap["isbn"])
					if format, ok := resultMap["format"].(string); ok {
						fmt.Printf("   Format: %s\n", format)
					}
				} else {
					fmt.Printf("Invalid ISBN: %s\n", resultMap["error"])
					if input, ok := resultMap["input"].(string); ok {
						fmt.Printf("   Input: %s\n", input)
					}
				}
			} else {
				fmt.Printf("ISBN validation result: %v\n", result)
			}
		} else {
			fmt.Printf("ISBN validation result: %v\n", result)
		}
	} else if isbnOperation == "generate" {
		if isbnCount == 1 {
			// Single ISBN
			if isbn, ok := result.(string); ok {
				fmt.Printf("Generated ISBN: %s\n", isbn)
			} else {
				fmt.Printf("Generated ISBN: %v\n", result)
			}
		} else {
			// Multiple ISBNs
			if isbns, ok := result.([]string); ok {
				fmt.Printf("Generated %d ISBNs:\n", len(isbns))
				for i, isbn := range isbns {
					fmt.Printf("  %d. %s\n", i+1, isbn)
				}
			} else {
				fmt.Printf("Generated ISBNs: %v\n", result)
			}
		}
	} else {
		fmt.Printf("ISBN result: %v\n", result)
	}

	return nil
}
