package tools

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

// CreditCardTool implements credit card number validation and generation
type CreditCardTool struct{}

// NewCreditCardTool creates a new credit card tool instance
func NewCreditCardTool() *CreditCardTool {
	return &CreditCardTool{}
}

// Name returns the tool name
func (c *CreditCardTool) Name() string {
	return "creditcard"
}

// Description returns the tool description
func (c *CreditCardTool) Description() string {
	return "Generate and validate credit card numbers using Luhn algorithm with card type support"
}

// Execute processes the credit card tool request
func (c *CreditCardTool) Execute(params map[string]interface{}) (interface{}, error) {
	// Validate parameters first
	if err := c.ValidateParams(params); err != nil {
		return nil, err
	}

	operation, _ := params["operation"].(string)
	if operation == "" {
		operation = "validate" // Default to validate
	}

	switch operation {
	case "validate":
		return c.validateCreditCard(params)
	case "generate":
		return c.generateCreditCard(params)
	default:
		return nil, fmt.Errorf("invalid operation: %s. Supported operations: validate, generate", operation)
	}
}

// ValidateParams validates the input parameters
func (c *CreditCardTool) ValidateParams(params map[string]interface{}) error {
	// Validate operation
	if operation, ok := params["operation"]; ok {
		if opStr, ok := operation.(string); ok {
			if opStr != "validate" && opStr != "generate" {
				return fmt.Errorf("invalid operation: %s. Supported operations: validate, generate", opStr)
			}
		} else {
			return fmt.Errorf("operation must be a string")
		}
	}

	// Validate input for validation operation
	if operation, ok := params["operation"]; ok {
		if opStr, ok := operation.(string); ok && opStr == "validate" {
			if input, ok := params["input"]; !ok || input == "" {
				return fmt.Errorf("input parameter is required for validation")
			}
		}
	}

	// Validate count
	if count, ok := params["count"]; ok {
		if countFloat, ok := count.(float64); ok {
			if countFloat < 1 || countFloat > 100 {
				return fmt.Errorf("count must be between 1 and 100")
			}
		} else {
			return fmt.Errorf("count must be a number")
		}
	}

	// Validate card type
	if cardType, ok := params["card-type"]; ok {
		if typeStr, ok := cardType.(string); ok {
			if typeStr != "" {
				validTypes := []string{"visa", "mastercard", "amex", "discover", "diners", "jcb"}
				if !contains(validTypes, typeStr) {
					return fmt.Errorf("invalid card type: %s. Supported types: %s", typeStr, strings.Join(validTypes, ", "))
				}
			}
		} else {
			return fmt.Errorf("card-type must be a string")
		}
	}

	return nil
}

// GetInputSchema returns the JSON schema for input parameters
func (c *CreditCardTool) GetInputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"operation": map[string]interface{}{
				"type":        "string",
				"description": "Operation to perform: 'validate' or 'generate'",
				"enum":        []string{"validate", "generate"},
			},
			"input": map[string]interface{}{
				"type":        "string",
				"description": "Credit card number to validate (required for validate operation)",
			},
			"card-type": map[string]interface{}{
				"type":        "string",
				"description": "Card type for generation: visa, mastercard, amex, discover, diners, jcb",
				"enum":        []string{"visa", "mastercard", "amex", "discover", "diners", "jcb"},
			},
			"count": map[string]interface{}{
				"type":        "number",
				"description": "Number of credit cards to generate (1-100, default: 1)",
				"minimum":     1,
				"maximum":     100,
			},
		},
		"required": []string{},
	}
}

// GetOutputSchema returns the JSON schema for output
func (c *CreditCardTool) GetOutputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"valid": map[string]interface{}{
				"type":        "boolean",
				"description": "Whether the credit card number is valid",
			},
			"card": map[string]interface{}{
				"type":        "string",
				"description": "Generated credit card number",
			},
			"cards": map[string]interface{}{
				"type":        "array",
				"description": "Generated credit card numbers",
				"items": map[string]interface{}{
					"type": "string",
				},
			},
			"type": map[string]interface{}{
				"type":        "string",
				"description": "Detected or generated card type",
			},
			"error": map[string]interface{}{
				"type":        "string",
				"description": "Error message if validation fails",
			},
		},
	}
}

// GetResources returns the list of resources this tool provides
func (c *CreditCardTool) GetResources() []Resource {
	return []Resource{
		{
			Name:     "Credit Card Types",
			URI:      "creditcard://types",
			MIMEType: "application/json",
		},
		{
			Name:     "Luhn Algorithm",
			URI:      "creditcard://luhn",
			MIMEType: "application/json",
		},
		{
			Name:     "Credit Card Examples",
			URI:      "creditcard://examples",
			MIMEType: "application/json",
		},
	}
}

// ReadResource reads a specific resource by URI
func (c *CreditCardTool) ReadResource(uri string) (string, error) {
	switch uri {
	case "creditcard://types":
		// Return supported credit card types
		types := map[string]interface{}{
			"types": []map[string]interface{}{
				{
					"name":        "visa",
					"description": "Visa cards",
					"prefixes":    []string{"4"},
					"lengths":     []int{13, 16, 19},
				},
				{
					"name":        "mastercard",
					"description": "Mastercard",
					"prefixes":    []string{"51", "52", "53", "54", "55", "2221-2720"},
					"lengths":     []int{16},
				},
				{
					"name":        "amex",
					"description": "American Express",
					"prefixes":    []string{"34", "37"},
					"lengths":     []int{15},
				},
				{
					"name":        "discover",
					"description": "Discover",
					"prefixes":    []string{"6011", "65", "644-649"},
					"lengths":     []int{16},
				},
				{
					"name":        "diners",
					"description": "Diners Club",
					"prefixes":    []string{"300-305", "36", "38"},
					"lengths":     []int{14},
				},
				{
					"name":        "jcb",
					"description": "JCB",
					"prefixes":    []string{"3528-3589"},
					"lengths":     []int{16},
				},
			},
		}
		jsonData, err := json.Marshal(types)
		if err != nil {
			return "", fmt.Errorf("failed to marshal types: %w", err)
		}
		return string(jsonData), nil
	case "creditcard://luhn":
		// Return Luhn algorithm documentation
		algorithm := map[string]interface{}{
			"name":        "Luhn Algorithm",
			"description": "Credit card number validation algorithm",
			"steps": []string{
				"Starting from the rightmost digit, double every second digit",
				"If doubling results in a two-digit number, add the digits together",
				"Sum all the digits",
				"If the total is divisible by 10, the number is valid",
			},
			"example": map[string]interface{}{
				"number":      "4532015112830366",
				"calculation": "4+1+3+2+0+1+5+1+1+2+8+3+0+3+1+2 = 40",
				"valid":       true,
			},
		}
		jsonData, err := json.Marshal(algorithm)
		if err != nil {
			return "", fmt.Errorf("failed to marshal algorithm: %w", err)
		}
		return string(jsonData), nil
	case "creditcard://examples":
		// Return example credit card numbers
		examples := []map[string]interface{}{
			{
				"number":      "4532015112830366",
				"type":        "visa",
				"valid":       true,
				"description": "Example valid Visa card",
			},
			{
				"number":      "5555555555554444",
				"type":        "mastercard",
				"valid":       true,
				"description": "Example valid Mastercard",
			},
			{
				"number":      "378282246310005",
				"type":        "amex",
				"valid":       true,
				"description": "Example valid American Express",
			},
			{
				"number":      "4532015112830367",
				"type":        "visa",
				"valid":       false,
				"description": "Example invalid Visa card (wrong check digit)",
			},
		}
		jsonData, err := json.Marshal(examples)
		if err != nil {
			return "", fmt.Errorf("failed to marshal examples: %w", err)
		}
		return string(jsonData), nil
	default:
		return "", fmt.Errorf("unknown resource URI: %s", uri)
	}
}

// validateCreditCard validates a credit card number
func (c *CreditCardTool) validateCreditCard(params map[string]interface{}) (interface{}, error) {
	input, _ := params["input"].(string)
	if input == "" {
		return nil, fmt.Errorf("input parameter is required for validation")
	}

	// Clean the input (remove spaces and dashes)
	cleanInput := strings.ReplaceAll(strings.ReplaceAll(input, " ", ""), "-", "")

	// Check if it's all digits
	if !isNumeric(cleanInput) {
		return map[string]interface{}{
			"valid": false,
			"error": "credit card number must contain only digits",
			"input": input,
		}, nil
	}

	// Check length
	if len(cleanInput) < 13 || len(cleanInput) > 19 {
		return map[string]interface{}{
			"valid": false,
			"error": "credit card number must be between 13 and 19 digits",
			"input": input,
		}, nil
	}

	// Validate using Luhn algorithm
	if !c.luhnCheck(cleanInput) {
		return map[string]interface{}{
			"valid": false,
			"error": "invalid check digit",
			"input": input,
		}, nil
	}

	// Detect card type
	cardType := c.detectCardType(cleanInput)

	return map[string]interface{}{
		"valid": true,
		"card":  cleanInput,
		"type":  cardType,
		"input": input,
	}, nil
}

// generateCreditCard generates credit card numbers
func (c *CreditCardTool) generateCreditCard(params map[string]interface{}) (interface{}, error) {
	count := 1
	if c, ok := params["count"].(float64); ok {
		count = int(c)
	}

	cardType := ""
	if ct, ok := params["card-type"].(string); ok {
		cardType = ct
	}

	if count == 1 {
		card, err := c.generateSingleCard(cardType)
		if err != nil {
			return nil, err
		}
		return card, nil
	}

	cards := make([]string, count)
	for i := range count {
		card, err := c.generateSingleCard(cardType)
		if err != nil {
			return nil, err
		}
		cards[i] = card
	}

	return cards, nil
}

// generateSingleCard generates a single credit card number
func (c *CreditCardTool) generateSingleCard(cardType string) (string, error) {
	if cardType == "" {
		// Random card type
		types := []string{"visa", "mastercard", "amex", "discover", "diners", "jcb"}
		cardType = types[rand.Intn(len(types))]
	}

	var prefix string
	var length int

	switch cardType {
	case "visa":
		prefix = "4"
		length = 16
	case "mastercard":
		prefixes := []string{"51", "52", "53", "54", "55"}
		prefix = prefixes[rand.Intn(len(prefixes))]
		length = 16
	case "amex":
		prefixes := []string{"34", "37"}
		prefix = prefixes[rand.Intn(len(prefixes))]
		length = 15
	case "discover":
		prefix = "6011"
		length = 16
	case "diners":
		prefixes := []string{"300", "301", "302", "303", "304", "305", "36", "38"}
		prefix = prefixes[rand.Intn(len(prefixes))]
		length = 14
	case "jcb":
		prefix = "35"
		length = 16
	default:
		return "", fmt.Errorf("unsupported card type: %s", cardType)
	}

	// Generate random digits for the remaining positions
	remainingLength := length - len(prefix) - 1 // -1 for check digit
	randomDigits := ""
	for i := 0; i < remainingLength; i++ {
		randomDigits += strconv.Itoa(rand.Intn(10))
	}

	// Combine prefix and random digits
	partialNumber := prefix + randomDigits

	// Calculate check digit using Luhn algorithm
	checkDigit := c.calculateCheckDigit(partialNumber)

	return partialNumber + strconv.Itoa(checkDigit), nil
}

// luhnCheck validates a credit card number using the Luhn algorithm
func (c *CreditCardTool) luhnCheck(number string) bool {
	sum := 0
	alternate := false

	// Process digits from right to left
	for i := len(number) - 1; i >= 0; i-- {
		digit := int(number[i] - '0')

		if alternate {
			digit *= 2
			if digit > 9 {
				digit = digit/10 + digit%10
			}
		}

		sum += digit
		alternate = !alternate
	}

	return sum%10 == 0
}

// calculateCheckDigit calculates the check digit for a partial credit card number
func (c *CreditCardTool) calculateCheckDigit(partialNumber string) int {
	sum := 0
	alternate := true

	// Process digits from right to left
	for i := len(partialNumber) - 1; i >= 0; i-- {
		digit := int(partialNumber[i] - '0')

		if alternate {
			digit *= 2
			if digit > 9 {
				digit = digit/10 + digit%10
			}
		}

		sum += digit
		alternate = !alternate
	}

	return (10 - (sum % 10)) % 10
}

// detectCardType detects the card type based on the number
func (c *CreditCardTool) detectCardType(number string) string {
	// Visa: starts with 4
	if strings.HasPrefix(number, "4") {
		return "visa"
	}

	// American Express: starts with 34 or 37
	if strings.HasPrefix(number, "34") || strings.HasPrefix(number, "37") {
		return "amex"
	}

	// Mastercard: starts with 51-55 or 2221-2720
	if strings.HasPrefix(number, "5") {
		firstTwo := number[:2]
		if firstTwo >= "51" && firstTwo <= "55" {
			return "mastercard"
		}
	}
	if strings.HasPrefix(number, "2221") || strings.HasPrefix(number, "2720") {
		return "mastercard"
	}

	// Discover: starts with 6011, 65, or 644-649
	if strings.HasPrefix(number, "6011") || strings.HasPrefix(number, "65") {
		return "discover"
	}
	if strings.HasPrefix(number, "644") || strings.HasPrefix(number, "645") || strings.HasPrefix(number, "646") ||
		strings.HasPrefix(number, "647") || strings.HasPrefix(number, "648") || strings.HasPrefix(number, "649") {
		return "discover"
	}

	// Diners Club: starts with 300-305, 36, or 38
	if strings.HasPrefix(number, "300") || strings.HasPrefix(number, "301") || strings.HasPrefix(number, "302") ||
		strings.HasPrefix(number, "303") || strings.HasPrefix(number, "304") || strings.HasPrefix(number, "305") {
		return "diners"
	}
	if strings.HasPrefix(number, "36") || strings.HasPrefix(number, "38") {
		return "diners"
	}

	// JCB: starts with 3528-3589
	if strings.HasPrefix(number, "35") {
		firstFour := number[:4]
		if firstFour >= "3528" && firstFour <= "3589" {
			return "jcb"
		}
	}

	return "unknown"
}

// isNumeric checks if a string contains only digits
func isNumeric(s string) bool {
	for _, char := range s {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}
