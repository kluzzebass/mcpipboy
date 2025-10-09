package tools

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

// EAN13Tool implements EAN-13 validation and generation
type EAN13Tool struct{}

// NewEAN13Tool creates a new EAN-13 tool instance
func NewEAN13Tool() *EAN13Tool {
	return &EAN13Tool{}
}

// Name returns the tool name
func (e *EAN13Tool) Name() string {
	return "ean13"
}

// Description returns the tool description
func (e *EAN13Tool) Description() string {
	return "Generate and validate European Article Numbers (EAN-13) with checksum validation"
}

// Execute processes the EAN-13 tool request
func (e *EAN13Tool) Execute(params map[string]interface{}) (interface{}, error) {
	// Validate parameters first
	if err := e.ValidateParams(params); err != nil {
		return nil, err
	}

	operation, _ := params["operation"].(string)
	if operation == "" {
		operation = "validate" // Default to validate
	}

	switch operation {
	case "validate":
		return e.validateEAN13(params)
	case "generate":
		return e.generateEAN13(params)
	default:
		return nil, fmt.Errorf("invalid operation: %s. Supported operations: validate, generate", operation)
	}
}

// ValidateParams validates the input parameters
func (e *EAN13Tool) ValidateParams(params map[string]interface{}) error {
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

	return nil
}

// GetInputSchema returns the JSON schema for input parameters
func (e *EAN13Tool) GetInputSchema() map[string]interface{} {
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
				"description": "EAN-13 number to validate (required for validate operation)",
			},
			"count": map[string]interface{}{
				"type":        "number",
				"description": "Number of EAN-13s to generate (1-100, default: 1)",
				"minimum":     1,
				"maximum":     100,
			},
		},
		"required": []string{},
	}
}

// GetOutputSchema returns the JSON schema for output
func (e *EAN13Tool) GetOutputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"valid": map[string]interface{}{
				"type":        "boolean",
				"description": "Whether the EAN-13 number is valid",
			},
			"ean13": map[string]interface{}{
				"type":        "string",
				"description": "Generated EAN-13 number",
			},
			"ean13s": map[string]interface{}{
				"type":        "array",
				"description": "Generated EAN-13 numbers",
				"items": map[string]interface{}{
					"type": "string",
				},
			},
			"error": map[string]interface{}{
				"type":        "string",
				"description": "Error message if validation fails",
			},
		},
	}
}

// GetResources returns the list of resources this tool provides
func (e *EAN13Tool) GetResources() []Resource {
	return []Resource{
		{
			Name:     "EAN-13 Algorithm",
			URI:      "ean13://algorithm",
			MIMEType: "application/json",
		},
		{
			Name:     "EAN-13 Examples",
			URI:      "ean13://examples",
			MIMEType: "application/json",
		},
	}
}

// ReadResource reads a specific resource by URI
func (e *EAN13Tool) ReadResource(uri string) (string, error) {
	switch uri {
	case "ean13://algorithm":
		// Return EAN-13 algorithm documentation
		algorithm := map[string]interface{}{
			"name":        "EAN-13 Algorithm",
			"description": "European Article Number validation and generation",
			"format":      "13-digit number with check digit",
			"algorithm": map[string]interface{}{
				"description": "Weighted sum algorithm for check digit calculation",
				"weights":     []int{1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3},
				"formula":     "Check digit = (10 - sum) mod 10",
				"steps": []string{
					"Multiply odd positions by 1, even positions by 3",
					"Sum all weighted digits",
					"Calculate check digit: (10 - sum) mod 10",
				},
				"example": map[string]interface{}{
					"ean13":       "1234567890123",
					"calculation": "1×1 + 2×3 + 3×1 + 4×3 + 5×1 + 6×3 + 7×1 + 8×3 + 9×1 + 0×3 + 1×1 + 2×3 = 1 + 6 + 3 + 12 + 5 + 18 + 7 + 24 + 9 + 0 + 1 + 6 = 92",
					"check_digit": "(10 - 92) mod 10 = 8",
					"result":      "1234567890128",
				},
			},
			"validation": map[string]interface{}{
				"description": "To validate an EAN-13 number, calculate the check digit and compare with the last digit",
				"steps": []string{
					"Extract the first 12 digits",
					"Apply the weighted sum formula",
					"Calculate check digit (10 - sum) mod 10",
					"Compare with the 13th digit",
				},
			},
		}
		jsonData, err := json.Marshal(algorithm)
		if err != nil {
			return "", fmt.Errorf("failed to marshal algorithm: %w", err)
		}
		return string(jsonData), nil
	case "ean13://examples":
		// Return example EAN-13 numbers
		examples := []map[string]interface{}{
			{
				"ean13":       "1234567890128",
				"valid":       true,
				"description": "Example valid EAN-13 number",
			},
			{
				"ean13":       "9780123456786",
				"valid":       true,
				"description": "Example valid EAN-13 (ISBN-13 format)",
			},
			{
				"ean13":       "1234567890123",
				"valid":       false,
				"description": "Example invalid EAN-13 (wrong check digit)",
			},
			{
				"ean13":       "123456789012",
				"valid":       false,
				"description": "Invalid EAN-13 (too short)",
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

// validateEAN13 validates an EAN-13 number
func (e *EAN13Tool) validateEAN13(params map[string]interface{}) (interface{}, error) {
	input, _ := params["input"].(string)
	if input == "" {
		return nil, fmt.Errorf("input parameter is required for validation")
	}

	// Clean the input (remove spaces, dashes, and hyphens)
	cleanInput := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(input, " ", ""), "-", ""), "—", "")

	// Validate EAN-13
	isValid, errorMsg := e.validateEAN13Number(cleanInput)

	if !isValid {
		return map[string]interface{}{
			"valid": false,
			"error": errorMsg,
			"input": input,
		}, nil
	}

	return map[string]interface{}{
		"valid": true,
		"ean13": cleanInput,
		"input": input,
	}, nil
}

// generateEAN13 generates EAN-13 numbers
func (e *EAN13Tool) generateEAN13(params map[string]interface{}) (interface{}, error) {
	count := 1
	if c, ok := params["count"].(float64); ok {
		count = int(c)
	}

	if count == 1 {
		ean13, err := e.generateSingleEAN13()
		if err != nil {
			return nil, err
		}
		return ean13, nil
	}

	ean13s := make([]string, count)
	for idx := range count {
		ean13, err := e.generateSingleEAN13()
		if err != nil {
			return nil, err
		}
		ean13s[idx] = ean13
	}

	return ean13s, nil
}

// validateEAN13Number validates an EAN-13 number
func (e *EAN13Tool) validateEAN13Number(ean13 string) (bool, string) {
	if len(ean13) != 13 {
		return false, "EAN-13 must be exactly 13 characters"
	}

	// Check if all characters are digits
	for i := 0; i < 13; i++ {
		if ean13[i] < '0' || ean13[i] > '9' {
			return false, "EAN-13 must contain only digits"
		}
	}

	// Calculate check digit using EAN-13 algorithm
	sum := 0
	for i := 0; i < 12; i++ {
		digit := int(ean13[i] - '0')
		if i%2 == 0 {
			sum += digit * 1
		} else {
			sum += digit * 3
		}
	}

	checkDigit := (10 - (sum % 10)) % 10
	expectedCheckDigit := strconv.Itoa(checkDigit)

	if string(ean13[12]) != expectedCheckDigit {
		return false, fmt.Sprintf("invalid check digit. Expected %s, got %c", expectedCheckDigit, ean13[12])
	}

	return true, ""
}

// generateSingleEAN13 generates a single EAN-13 number
func (e *EAN13Tool) generateSingleEAN13() (string, error) {
	// Generate 12 random digits
	digits := make([]int, 12)
	for i := range 12 {
		digits[i] = rand.Intn(10)
	}

	// Calculate check digit using EAN-13 algorithm
	sum := 0
	for i := 0; i < 12; i++ {
		if i%2 == 0 {
			sum += digits[i] * 1
		} else {
			sum += digits[i] * 3
		}
	}

	checkDigit := (10 - (sum % 10)) % 10

	// Build the EAN-13
	ean13 := ""
	for i := 0; i < 12; i++ {
		ean13 += strconv.Itoa(digits[i])
	}
	ean13 += strconv.Itoa(checkDigit)

	return ean13, nil
}
