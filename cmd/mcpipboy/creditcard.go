package main

import (
	"fmt"

	"github.com/kluzzebass/mcpipboy/internal/tools"
	"github.com/spf13/cobra"
)

var (
	creditCardOperation string
	creditCardInput     string
	creditCardType      string
	creditCardCount     int
)

// creditCardCmd represents the creditcard command
var creditCardCmd = &cobra.Command{
	Use:   "creditcard",
	Short: "Generate and validate credit card numbers using Luhn algorithm",
	Long: `Generate and validate credit card numbers using the Luhn algorithm with card type support.

Supports major card types:
- Visa (starts with 4)
- Mastercard (51-55, 2221-2720)
- American Express (34, 37)
- Discover (6011, 65, 644-649)
- Diners Club (300-305, 36, 38)
- JCB (3528-3589)

Examples:
  mcpipboy creditcard --operation validate --input "4532015112830366"
  mcpipboy creditcard --operation generate --card-type visa --count 5
  mcpipboy creditcard --operation generate --card-type amex
  mcpipboy creditcard --operation validate --input "5555 5555 5555 4444"`,
	RunE: runCreditCard,
}

func init() {
	creditCardCmd.Flags().StringVar(&creditCardOperation, "operation", "validate", "Operation to perform: validate or generate")
	creditCardCmd.Flags().StringVar(&creditCardInput, "input", "", "Credit card number to validate")
	creditCardCmd.Flags().StringVar(&creditCardType, "card-type", "", "Card type for generation: visa, mastercard, amex, discover, diners, jcb")
	creditCardCmd.Flags().IntVar(&creditCardCount, "count", 1, "Number of credit cards to generate (1-100)")

	creditCardCmd.GroupID = "tools"
	rootCmd.AddCommand(creditCardCmd)
}

func runCreditCard(cmd *cobra.Command, args []string) error {
	tool := tools.NewCreditCardTool()

	// Build parameters map
	params := make(map[string]interface{})
	if creditCardOperation != "" {
		params["operation"] = creditCardOperation
	}
	if creditCardInput != "" {
		params["input"] = creditCardInput
	}
	if creditCardType != "" {
		params["card-type"] = creditCardType
	}
	params["count"] = float64(creditCardCount)

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
	if creditCardOperation == "validate" {
		if resultMap, ok := result.(map[string]interface{}); ok {
			if valid, ok := resultMap["valid"].(bool); ok {
				if valid {
					fmt.Printf("Valid credit card: %s\n", resultMap["card"])
					if cardType, ok := resultMap["type"].(string); ok {
						fmt.Printf("   Type: %s\n", cardType)
					}
				} else {
					fmt.Printf("Invalid credit card: %s\n", resultMap["error"])
					if input, ok := resultMap["input"].(string); ok {
						fmt.Printf("   Input: %s\n", input)
					}
				}
			}
		}
	} else if creditCardOperation == "generate" {
		if creditCardCount == 1 {
			// Single card
			if card, ok := result.(string); ok {
				fmt.Println(card)
			}
		} else {
			// Multiple cards
			if cards, ok := result.([]string); ok {
				for _, card := range cards {
					fmt.Println(card)
				}
			}
		}
	}

	return nil
}
