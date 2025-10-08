package tools

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

// IMOTool implements IMO number validation and generation
type IMOTool struct{}

// NewIMOTool creates a new IMO tool instance
func NewIMOTool() *IMOTool {
	return &IMOTool{}
}

// Name returns the tool name
func (i *IMOTool) Name() string {
	return "imo"
}

// Description returns the tool description
func (i *IMOTool) Description() string {
	return "Generate and validate International Maritime Organization (IMO) numbers. IMO numbers are 7-digit numbers with a check digit calculated using a weighted sum algorithm."
}

// Execute runs the IMO tool
func (i *IMOTool) Execute(params map[string]interface{}) (interface{}, error) {
	operation, _ := params["operation"].(string)
	if operation == "" {
		operation = "validate" // Default to validate
	}

	switch operation {
	case "validate":
		return i.validateIMO(params)
	case "generate":
		return i.generateIMO(params)
	default:
		return nil, fmt.Errorf("invalid operation: %s. Must be 'validate' or 'generate'", operation)
	}
}

// validateIMO validates an IMO number
func (i *IMOTool) validateIMO(params map[string]interface{}) (interface{}, error) {
	input, _ := params["input"].(string)
	if input == "" {
		return nil, fmt.Errorf("input parameter is required for validation")
	}

	// Clean the input (remove spaces, dashes, etc.)
	cleanInput := strings.ReplaceAll(strings.ReplaceAll(input, " ", ""), "-", "")

	// Check if it's exactly 7 digits
	if len(cleanInput) != 7 {
		return map[string]interface{}{
			"valid": false,
			"error": "IMO number must be exactly 7 digits",
			"input": input,
		}, nil
	}

	// Parse as integer and validate
	imoNumber, err := strconv.Atoi(cleanInput)
	if err != nil {
		return map[string]interface{}{
			"valid": false,
			"error": "IMO number must contain only digits",
			"input": input,
		}, nil
	}

	// Extract digits using integer arithmetic
	digits := make([]int, 7)
	temp := imoNumber
	for i := range 7 {
		digits[6-i] = temp % 10
		temp /= 10
	}

	// Calculate expected check digit using weighted sum
	expectedCheckDigit := i.calculateCheckDigit(digits[:6])
	actualCheckDigit := digits[6]

	if expectedCheckDigit != actualCheckDigit {
		return map[string]interface{}{
			"valid": false,
			"error": fmt.Sprintf("invalid check digit. Expected %d, got %d", expectedCheckDigit, actualCheckDigit),
			"input": input,
		}, nil
	}

	return map[string]interface{}{
		"valid": true,
		"imo":   cleanInput,
		"input": input,
	}, nil
}

// generateIMO generates IMO numbers
func (i *IMOTool) generateIMO(params map[string]interface{}) (interface{}, error) {
	count, _ := params["count"].(int)
	if count <= 0 {
		count = 1
	}

	// Validate count
	if count > 100 {
		return nil, fmt.Errorf("count cannot exceed 100")
	}

	results := make([]string, count)

	for idx := range count {
		// Generate 6 random digits
		digits := make([]int, 6)
		for j := range 6 {
			digits[j] = rand.Intn(10)
		}

		// Calculate check digit using weighted sum
		checkDigit := i.calculateCheckDigit(digits)

		// Build the IMO number using integer arithmetic
		imo := 0
		for j := range 6 {
			imo = imo*10 + digits[j]
		}
		imo = imo*10 + checkDigit

		// Format as 7-digit string with leading zeros if needed
		results[idx] = fmt.Sprintf("%07d", imo)
	}

	if count == 1 {
		return results[0], nil
	}

	return results, nil
}

// ValidateParams validates the input parameters
func (i *IMOTool) ValidateParams(params map[string]interface{}) error {
	// Validate operation
	if operation, ok := params["operation"]; ok {
		if operationStr, ok := operation.(string); ok {
			if operationStr != "validate" && operationStr != "generate" {
				return fmt.Errorf("operation must be 'validate' or 'generate'")
			}
		} else {
			return fmt.Errorf("operation must be a string")
		}
	}

	// Validate count for generation
	if count, ok := params["count"]; ok {
		if countInt, ok := count.(int); ok {
			if countInt < 1 {
				return fmt.Errorf("count must be at least 1")
			}
			if countInt > 100 {
				return fmt.Errorf("count cannot exceed 100")
			}
		} else {
			return fmt.Errorf("count must be an integer")
		}
	}

	// Validate input for validation
	if operation, ok := params["operation"]; ok {
		if operationStr, ok := operation.(string); ok && operationStr == "validate" {
			if input, ok := params["input"]; !ok {
				return fmt.Errorf("input parameter is required for validation")
			} else if _, ok := input.(string); !ok {
				return fmt.Errorf("input must be a string")
			}
		}
	}

	return nil
}

// GetInputSchema returns the JSON schema for tool input parameters
func (i *IMOTool) GetInputSchema() map[string]interface{} {
	return CreateJSONSchema([]ParameterDefinition{
		{
			Name:        "operation",
			Type:        "string",
			Description: "Operation to perform: 'validate' or 'generate'",
			Required:    false,
			Enum:        []string{"validate", "generate"},
		},
		{
			Name:        "input",
			Type:        "string",
			Description: "IMO number to validate (required for validation operation)",
			Required:    false,
		},
		{
			Name:        "count",
			Type:        "integer",
			Description: "Number of IMO numbers to generate (default: 1, max: 100)",
			Required:    false,
		},
	})
}

// GetOutputSchema returns the JSON schema for tool output
func (i *IMOTool) GetOutputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"result": map[string]interface{}{
				"description": "Generated IMO number(s) or validation result",
			},
		},
	}
}

// calculateCheckDigit calculates the check digit for the first 6 digits using weighted sum
// Weights: 7, 6, 5, 4, 3, 2, 1 (from left to right)
func (i *IMOTool) calculateCheckDigit(digits []int) int {
	weights := []int{7, 6, 5, 4, 3, 2, 1}
	sum := 0
	for j := range 6 {
		sum += digits[j] * weights[j]
	}
	return sum % 10
}
