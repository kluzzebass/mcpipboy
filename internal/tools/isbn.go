package tools

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

// ISBNTool implements ISBN validation and generation
type ISBNTool struct{}

// NewISBNTool creates a new ISBN tool instance
func NewISBNTool() *ISBNTool {
	return &ISBNTool{}
}

// Name returns the tool name
func (i *ISBNTool) Name() string {
	return "isbn"
}

// Description returns the tool description
func (i *ISBNTool) Description() string {
	return "Generate and validate International Standard Book Numbers (ISBN-10 and ISBN-13) with format support"
}

// Execute processes the ISBN tool request
func (i *ISBNTool) Execute(params map[string]interface{}) (interface{}, error) {
	// Validate parameters first
	if err := i.ValidateParams(params); err != nil {
		return nil, err
	}

	operation, _ := params["operation"].(string)
	if operation == "" {
		operation = "validate" // Default to validate
	}

	switch operation {
	case "validate":
		return i.validateISBN(params)
	case "generate":
		return i.generateISBN(params)
	default:
		return nil, fmt.Errorf("invalid operation: %s. Supported operations: validate, generate", operation)
	}
}

// ValidateParams validates the input parameters
func (i *ISBNTool) ValidateParams(params map[string]interface{}) error {
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

	// Validate format
	if format, ok := params["format"]; ok {
		if formatStr, ok := format.(string); ok {
			if formatStr != "" {
				validFormats := []string{"isbn10", "isbn13", "auto"}
				if !contains(validFormats, formatStr) {
					return fmt.Errorf("invalid format: %s. Supported formats: %s", formatStr, strings.Join(validFormats, ", "))
				}
			}
		} else {
			return fmt.Errorf("format must be a string")
		}
	}

	return nil
}

// GetInputSchema returns the JSON schema for input parameters
func (i *ISBNTool) GetInputSchema() map[string]interface{} {
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
				"description": "ISBN number to validate (required for validate operation)",
			},
			"format": map[string]interface{}{
				"type":        "string",
				"description": "ISBN format: isbn10, isbn13, or auto (default: auto for validation, isbn13 for generation)",
				"enum":        []string{"isbn10", "isbn13", "auto"},
			},
			"count": map[string]interface{}{
				"type":        "number",
				"description": "Number of ISBNs to generate (1-100, default: 1)",
				"minimum":     1,
				"maximum":     100,
			},
		},
		"required": []string{},
	}
}

// GetOutputSchema returns the JSON schema for output
func (i *ISBNTool) GetOutputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"valid": map[string]interface{}{
				"type":        "boolean",
				"description": "Whether the ISBN number is valid",
			},
			"isbn": map[string]interface{}{
				"type":        "string",
				"description": "Generated ISBN number",
			},
			"isbns": map[string]interface{}{
				"type":        "array",
				"description": "Generated ISBN numbers",
				"items": map[string]interface{}{
					"type": "string",
				},
			},
			"format": map[string]interface{}{
				"type":        "string",
				"description": "Detected or generated ISBN format (ISBN-10 or ISBN-13)",
			},
			"error": map[string]interface{}{
				"type":        "string",
				"description": "Error message if validation fails",
			},
		},
	}
}

// GetResources returns the list of resources this tool provides
func (i *ISBNTool) GetResources() []Resource {
	return []Resource{
		{
			Name:     "ISBN Formats",
			URI:      "isbn://formats",
			MIMEType: "application/json",
		},
		{
			Name:     "ISBN Algorithms",
			URI:      "isbn://algorithms",
			MIMEType: "application/json",
		},
		{
			Name:     "ISBN Examples",
			URI:      "isbn://examples",
			MIMEType: "application/json",
		},
	}
}

// ReadResource reads a specific resource by URI
func (i *ISBNTool) ReadResource(uri string) (string, error) {
	switch uri {
	case "isbn://formats":
		// Return ISBN format information
		formats := map[string]interface{}{
			"formats": []map[string]interface{}{
				{
					"name":        "ISBN-10",
					"description": "10-digit ISBN with weighted sum check digit",
					"length":      10,
					"algorithm":   "Weighted sum (10×d1 + 9×d2 + ... + 2×d9) mod 11",
					"check_digit": "0-9 or X (10)",
				},
				{
					"name":        "ISBN-13",
					"description": "13-digit ISBN with EAN-13 check digit",
					"length":      13,
					"algorithm":   "EAN-13 (odd×1 + even×3) mod 10",
					"check_digit": "0-9",
				},
			},
		}
		jsonData, err := json.Marshal(formats)
		if err != nil {
			return "", fmt.Errorf("failed to marshal formats: %w", err)
		}
		return string(jsonData), nil
	case "isbn://algorithms":
		// Return ISBN algorithm documentation
		algorithms := map[string]interface{}{
			"isbn10": map[string]interface{}{
				"name":        "ISBN-10 Algorithm",
				"description": "Weighted sum algorithm for ISBN-10 validation",
				"steps": []string{
					"Multiply each digit by its position weight (10, 9, 8, ..., 2)",
					"Sum all weighted digits",
					"Calculate check digit: (11 - sum) mod 11",
					"Use X for check digit 10",
				},
				"example": map[string]interface{}{
					"isbn":        "0-123456-78-9",
					"calculation": "10×0 + 9×1 + 8×2 + 7×3 + 6×4 + 5×5 + 4×6 + 3×7 + 2×8 = 165",
					"check_digit": "(11 - 165) mod 11 = 1",
					"result":      "0-123456-78-1",
				},
			},
			"isbn13": map[string]interface{}{
				"name":        "ISBN-13 Algorithm",
				"description": "EAN-13 algorithm for ISBN-13 validation",
				"steps": []string{
					"Multiply odd positions by 1, even positions by 3",
					"Sum all weighted digits",
					"Calculate check digit: (10 - sum) mod 10",
				},
				"example": map[string]interface{}{
					"isbn":        "978-0-123456-78-9",
					"calculation": "9×1 + 7×3 + 8×1 + 0×3 + 1×1 + 2×3 + 3×1 + 4×3 + 5×1 + 6×3 + 7×1 + 8×3 = 120",
					"check_digit": "(10 - 120) mod 10 = 0",
					"result":      "978-0-123456-78-0",
				},
			},
		}
		jsonData, err := json.Marshal(algorithms)
		if err != nil {
			return "", fmt.Errorf("failed to marshal algorithms: %w", err)
		}
		return string(jsonData), nil
	case "isbn://examples":
		// Return example ISBN numbers
		examples := []map[string]interface{}{
			{
				"isbn":        "0-123456-78-9",
				"format":      "ISBN-10",
				"valid":       true,
				"description": "Example valid ISBN-10",
			},
			{
				"isbn":        "978-0-123456-78-9",
				"format":      "ISBN-13",
				"valid":       true,
				"description": "Example valid ISBN-13",
			},
			{
				"isbn":        "0-123456-78-X",
				"format":      "ISBN-10",
				"valid":       true,
				"description": "Example valid ISBN-10 with X check digit",
			},
			{
				"isbn":        "0-123456-78-8",
				"format":      "ISBN-10",
				"valid":       false,
				"description": "Example invalid ISBN-10 (wrong check digit)",
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

// validateISBN validates an ISBN number
func (i *ISBNTool) validateISBN(params map[string]interface{}) (interface{}, error) {
	input, _ := params["input"].(string)
	if input == "" {
		return nil, fmt.Errorf("input parameter is required for validation")
	}

	format, _ := params["format"].(string)
	if format == "" {
		format = "auto" // Default to auto-detection
	}

	// Clean the input (remove spaces, dashes, and hyphens)
	cleanInput := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(input, " ", ""), "-", ""), "—", "")

	// Auto-detect format if needed
	if format == "auto" {
		if len(cleanInput) == 10 {
			format = "isbn10"
		} else if len(cleanInput) == 13 {
			format = "isbn13"
		} else {
			return map[string]interface{}{
				"valid": false,
				"error": "unable to auto-detect ISBN format (must be 10 or 13 digits)",
				"input": input,
			}, nil
		}
	}

	// Validate based on format
	var isValid bool
	var errorMsg string

	switch format {
	case "isbn10":
		isValid, errorMsg = i.validateISBN10(cleanInput)
	case "isbn13":
		isValid, errorMsg = i.validateISBN13(cleanInput)
	default:
		return nil, fmt.Errorf("invalid format: %s", format)
	}

	if !isValid {
		return map[string]interface{}{
			"valid": false,
			"error": errorMsg,
			"input": input,
		}, nil
	}

	return map[string]interface{}{
		"valid":  true,
		"isbn":   cleanInput,
		"format": strings.ToUpper(format),
		"input":  input,
	}, nil
}

// generateISBN generates ISBN numbers
func (i *ISBNTool) generateISBN(params map[string]interface{}) (interface{}, error) {
	count := 1
	if c, ok := params["count"].(float64); ok {
		count = int(c)
	}

	format := "isbn13" // Default to ISBN-13
	if f, ok := params["format"].(string); ok && f != "" {
		format = f
	}

	if count == 1 {
		isbn, err := i.generateSingleISBN(format)
		if err != nil {
			return nil, err
		}
		return isbn, nil
	}

	isbns := make([]string, count)
	for idx := range count {
		isbn, err := i.generateSingleISBN(format)
		if err != nil {
			return nil, err
		}
		isbns[idx] = isbn
	}

	return isbns, nil
}

// generateSingleISBN generates a single ISBN number
func (i *ISBNTool) generateSingleISBN(format string) (string, error) {
	switch format {
	case "isbn10":
		return i.generateISBN10()
	case "isbn13":
		return i.generateISBN13()
	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}
}

// validateISBN10 validates an ISBN-10 number
func (i *ISBNTool) validateISBN10(isbn string) (bool, string) {
	if len(isbn) != 10 {
		return false, "ISBN-10 must be exactly 10 characters"
	}

	// Check if all characters except the last are digits
	for i := 0; i < 9; i++ {
		if isbn[i] < '0' || isbn[i] > '9' {
			return false, "ISBN-10 must contain only digits (except last character)"
		}
	}

	// Check last character (can be digit or X)
	lastChar := isbn[9]
	if lastChar != 'X' && (lastChar < '0' || lastChar > '9') {
		return false, "ISBN-10 last character must be digit or X"
	}

	// Calculate check digit
	sum := 0
	for i := 0; i < 9; i++ {
		sum += int(isbn[i]-'0') * (10 - i)
	}

	checkDigit := (11 - (sum % 11)) % 11
	var expectedCheckDigit string
	if checkDigit == 10 {
		expectedCheckDigit = "X"
	} else {
		expectedCheckDigit = strconv.Itoa(checkDigit)
	}

	if string(lastChar) != expectedCheckDigit {
		return false, fmt.Sprintf("invalid check digit. Expected %s, got %c", expectedCheckDigit, lastChar)
	}

	return true, ""
}

// validateISBN13 validates an ISBN-13 number
func (i *ISBNTool) validateISBN13(isbn string) (bool, string) {
	if len(isbn) != 13 {
		return false, "ISBN-13 must be exactly 13 characters"
	}

	// Check if all characters are digits
	for i := 0; i < 13; i++ {
		if isbn[i] < '0' || isbn[i] > '9' {
			return false, "ISBN-13 must contain only digits"
		}
	}

	// Calculate check digit using EAN-13 algorithm
	sum := 0
	for i := 0; i < 12; i++ {
		digit := int(isbn[i] - '0')
		if i%2 == 0 {
			sum += digit * 1
		} else {
			sum += digit * 3
		}
	}

	checkDigit := (10 - (sum % 10)) % 10
	expectedCheckDigit := strconv.Itoa(checkDigit)

	if string(isbn[12]) != expectedCheckDigit {
		return false, fmt.Sprintf("invalid check digit. Expected %s, got %c", expectedCheckDigit, isbn[12])
	}

	return true, ""
}

// generateISBN10 generates a random ISBN-10 number
func (i *ISBNTool) generateISBN10() (string, error) {
	// Generate 9 random digits
	digits := make([]int, 9)
	for i := range 9 {
		digits[i] = rand.Intn(10)
	}

	// Calculate check digit
	sum := 0
	for i := 0; i < 9; i++ {
		sum += digits[i] * (10 - i)
	}

	checkDigit := (11 - (sum % 11)) % 11

	// Build the ISBN
	isbn := ""
	for i := 0; i < 9; i++ {
		isbn += strconv.Itoa(digits[i])
	}

	if checkDigit == 10 {
		isbn += "X"
	} else {
		isbn += strconv.Itoa(checkDigit)
	}

	return isbn, nil
}

// generateISBN13 generates a random ISBN-13 number
func (i *ISBNTool) generateISBN13() (string, error) {
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

	// Build the ISBN
	isbn := ""
	for i := 0; i < 12; i++ {
		isbn += strconv.Itoa(digits[i])
	}
	isbn += strconv.Itoa(checkDigit)

	return isbn, nil
}
