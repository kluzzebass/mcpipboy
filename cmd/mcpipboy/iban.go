package main

import (
	"fmt"

	"github.com/kluzzebass/mcpipboy/internal/tools"
	"github.com/spf13/cobra"
)

var (
	ibanOperation   string
	ibanInput       string
	ibanCountryCode string
	ibanCount       int
)

// ibanCmd represents the iban command
var ibanCmd = &cobra.Command{
	Use:   "iban",
	Short: "Generate and validate International Bank Account Numbers (IBAN)",
	Long: `Generate and validate International Bank Account Numbers (IBAN) using the MOD-97 checksum algorithm.

Supports major European countries and their IBAN formats:
- United Kingdom (GB): 22 characters
- Germany (DE): 22 characters  
- France (FR): 27 characters
- Italy (IT): 27 characters
- Spain (ES): 24 characters
- Netherlands (NL): 18 characters
- Belgium (BE): 16 characters
- Austria (AT): 20 characters
- Switzerland (CH): 21 characters
- Sweden (SE): 24 characters
- Norway (NO): 15 characters
- Denmark (DK): 18 characters
- Finland (FI): 18 characters
- Poland (PL): 28 characters
- Czech Republic (CZ): 24 characters
- Hungary (HU): 28 characters
- Romania (RO): 24 characters
- Bulgaria (BG): 22 characters
- Croatia (HR): 21 characters
- Slovenia (SI): 19 characters
- Slovakia (SK): 24 characters
- Lithuania (LT): 20 characters
- Latvia (LV): 21 characters
- Estonia (EE): 20 characters
- Ireland (IE): 22 characters
- Portugal (PT): 25 characters
- Greece (GR): 27 characters
- Cyprus (CY): 28 characters
- Malta (MT): 31 characters
- Luxembourg (LU): 20 characters

Examples:
  mcpipboy iban --operation validate --input "GB82WEST12345698765432"
  mcpipboy iban --operation generate --country-code "GB" --count 5
  mcpipboy iban --operation generate --country-code "DE"
  mcpipboy iban --operation validate --input "DE89 3704 0044 0532 0130 00"`,
	RunE: runIBAN,
}

func init() {
	ibanCmd.Flags().StringVar(&ibanOperation, "operation", "validate", "Operation to perform: validate or generate")
	ibanCmd.Flags().StringVar(&ibanInput, "input", "", "IBAN number to validate")
	ibanCmd.Flags().StringVar(&ibanCountryCode, "country-code", "", "Country code for generation (ISO 3166-1 alpha-2, e.g., 'GB', 'DE', 'FR')")
	ibanCmd.Flags().IntVar(&ibanCount, "count", 1, "Number of IBANs to generate (1-100)")

	ibanCmd.GroupID = "tools"
	rootCmd.AddCommand(ibanCmd)
}

func runIBAN(cmd *cobra.Command, args []string) error {
	tool := tools.NewIBANTool()

	// Build parameters map
	params := make(map[string]interface{})
	if ibanOperation != "" {
		params["operation"] = ibanOperation
	}
	if ibanInput != "" {
		params["input"] = ibanInput
	}
	if ibanCountryCode != "" {
		params["country-code"] = ibanCountryCode
	}
	params["count"] = float64(ibanCount)

	// Validate parameters
	if err := tool.ValidateParams(params); err != nil {
		return fmt.Errorf("parameter validation failed: %v", err)
	}

	// Execute the tool
	result, err := tool.Execute(params)
	if err != nil {
		return fmt.Errorf("execution failed: %v", err)
	}

	// Handle the result based on operation
	if ibanOperation == "validate" {
		if resultMap, ok := result.(map[string]interface{}); ok {
			if valid, ok := resultMap["valid"].(bool); ok {
				if valid {
					fmt.Printf("Valid IBAN: %s\n", resultMap["iban"])
					if country, ok := resultMap["country"].(string); ok {
						fmt.Printf("   Country: %s\n", country)
					}
				} else {
					fmt.Printf("Invalid IBAN: %s\n", resultMap["error"])
					if input, ok := resultMap["input"].(string); ok {
						fmt.Printf("   Input: %s\n", input)
					}
				}
			}
		}
	} else if ibanOperation == "generate" {
		if ibanCount == 1 {
			// Single IBAN
			if iban, ok := result.(string); ok {
				fmt.Println(iban)
			}
		} else {
			// Multiple IBANs
			if ibans, ok := result.([]string); ok {
				for _, iban := range ibans {
					fmt.Println(iban)
				}
			}
		}
	}

	return nil
}
