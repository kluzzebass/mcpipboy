package tools

import (
	"encoding/json"
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

// GetResources returns the list of resources this tool provides
func (i *IMOTool) GetResources() []Resource {
	return []Resource{
		{
			Name:     "IMO Algorithm",
			URI:      "imo://algorithm",
			MIMEType: "application/json",
		},
		{
			Name:     "IMO Examples",
			URI:      "imo://examples",
			MIMEType: "application/json",
		},
	}
}

// ReadResource reads a specific resource by URI
func (i *IMOTool) ReadResource(uri string) (string, error) {
	switch uri {
	case "imo://algorithm":
		// Return IMO algorithm documentation
		algorithm := map[string]interface{}{
			"name":        "IMO Number Algorithm",
			"description": "International Maritime Organization number validation and generation",
			"format":      "7-digit number with check digit",
			"algorithm": map[string]interface{}{
				"description": "Weighted sum algorithm for check digit calculation",
				"weights":     []int{7, 6, 5, 4, 3, 2, 1},
				"formula":     "Check digit = (7×d1 + 6×d2 + 5×d3 + 4×d4 + 3×d5 + 2×d6) mod 10",
				"example": map[string]interface{}{
					"digits":      []int{1, 2, 3, 4, 5, 6},
					"calculation": "7×1 + 6×2 + 5×3 + 4×4 + 3×5 + 2×6 = 7 + 12 + 15 + 16 + 15 + 12 = 77",
					"check_digit": "77 mod 10 = 7",
					"result":      "1234567",
				},
			},
			"validation": map[string]interface{}{
				"description": "To validate an IMO number, calculate the check digit and compare with the last digit",
				"steps": []string{
					"Extract the first 6 digits",
					"Apply the weighted sum formula",
					"Calculate check digit (sum mod 10)",
					"Compare with the 7th digit",
				},
			},
		}
		jsonData, err := json.Marshal(algorithm)
		if err != nil {
			return "", fmt.Errorf("failed to marshal algorithm: %w", err)
		}
		return string(jsonData), nil
	case "imo://examples":
		// Return example IMO numbers
		examples := []map[string]interface{}{
			{
				"imo":         "1234567",
				"valid":       true,
				"description": "Example valid IMO number",
			},
			{
				"imo":         "9074729",
				"valid":       true,
				"description": "Another valid IMO number",
			},
			{
				"imo":         "1234568",
				"valid":       false,
				"description": "Example invalid IMO number (wrong check digit)",
			},
			{
				"imo":         "123456",
				"valid":       false,
				"description": "Invalid IMO number (too short)",
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
