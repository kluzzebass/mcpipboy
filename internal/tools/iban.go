package tools

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

// IBANCountry represents a country with its IBAN information
type IBANCountry struct {
	Code   string
	Name   string
	Length int
}

// IBANTool implements IBAN validation and generation
type IBANTool struct {
	countries []IBANCountry
}

// NewIBANTool creates a new IBAN tool instance
func NewIBANTool() *IBANTool {
	tool := &IBANTool{}
	tool.populateCountries()
	return tool
}

// Name returns the tool name
func (i *IBANTool) Name() string {
	return "iban"
}

// Description returns the tool description
func (i *IBANTool) Description() string {
	return "Generate and validate International Bank Account Numbers (IBAN) with MOD-97 checksum algorithm"
}

// Execute processes the IBAN tool request
func (i *IBANTool) Execute(params map[string]interface{}) (interface{}, error) {
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
		return i.validateIBAN(params)
	case "generate":
		return i.generateIBAN(params)
	default:
		return nil, fmt.Errorf("invalid operation: %s. Supported operations: validate, generate", operation)
	}
}

// ValidateParams validates the input parameters
func (i *IBANTool) ValidateParams(params map[string]interface{}) error {
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

	// Validate country code
	if countryCode, ok := params["country-code"]; ok {
		if ccStr, ok := countryCode.(string); ok {
			if ccStr != "" {
				if !i.isValidCountryCode(ccStr) {
					return fmt.Errorf("invalid country code: %s. Must be a valid ISO 3166-1 alpha-2 country code", ccStr)
				}
			}
		} else {
			return fmt.Errorf("country-code must be a string")
		}
	}

	return nil
}

// GetInputSchema returns the JSON schema for input parameters
func (i *IBANTool) GetInputSchema() map[string]interface{} {
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
				"description": "IBAN number to validate (required for validate operation)",
			},
			"country-code": map[string]interface{}{
				"type":        "string",
				"description": "Country code for generation (ISO 3166-1 alpha-2, e.g., 'GB', 'DE', 'FR')",
			},
			"count": map[string]interface{}{
				"type":        "number",
				"description": "Number of IBANs to generate (1-100, default: 1)",
				"minimum":     1,
				"maximum":     100,
			},
		},
		"required": []string{},
	}
}

// GetOutputSchema returns the JSON schema for output
func (i *IBANTool) GetOutputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"valid": map[string]interface{}{
				"type":        "boolean",
				"description": "Whether the IBAN is valid",
			},
			"iban": map[string]interface{}{
				"type":        "string",
				"description": "Generated IBAN number",
			},
			"ibans": map[string]interface{}{
				"type":        "array",
				"description": "Generated IBAN numbers",
				"items": map[string]interface{}{
					"type": "string",
				},
			},
			"country": map[string]interface{}{
				"type":        "string",
				"description": "Country code of the IBAN",
			},
			"error": map[string]interface{}{
				"type":        "string",
				"description": "Error message if validation fails",
			},
		},
	}
}

// GetResources returns the list of resources this tool provides
func (i *IBANTool) GetResources() []Resource {
	return []Resource{
		{
			Name:     "IBAN Country Codes",
			URI:      "iban://countries",
			MIMEType: "application/json",
		},
		{
			Name:     "MOD-97 Algorithm",
			URI:      "iban://mod97",
			MIMEType: "application/json",
		},
		{
			Name:     "IBAN Examples",
			URI:      "iban://examples",
			MIMEType: "application/json",
		},
	}
}

// ReadResource reads a specific resource by URI
func (i *IBANTool) ReadResource(uri string) (string, error) {
	switch uri {
	case "iban://countries":
		// Return supported country codes
		countries := map[string]interface{}{
			"countries": i.getCountriesData(),
		}
		jsonData, err := json.Marshal(countries)
		if err != nil {
			return "", fmt.Errorf("failed to marshal countries: %w", err)
		}
		return string(jsonData), nil
	case "iban://mod97":
		// Return MOD-97 algorithm documentation
		algorithm := map[string]interface{}{
			"name":        "MOD-97 Algorithm",
			"description": "IBAN checksum validation algorithm",
			"steps": []string{
				"Move the first 4 characters to the end of the string",
				"Replace each letter with its numeric value (A=10, B=11, ..., Z=35)",
				"Calculate the remainder when dividing by 97",
				"If the remainder is 1, the IBAN is valid",
			},
			"example": map[string]interface{}{
				"iban":       "GB82WEST12345698765432",
				"rearranged": "WEST12345698765432GB82",
				"numeric":    "3214282912345698765432161182",
				"remainder":  1,
				"valid":      true,
			},
		}
		jsonData, err := json.Marshal(algorithm)
		if err != nil {
			return "", fmt.Errorf("failed to marshal algorithm: %w", err)
		}
		return string(jsonData), nil
	case "iban://examples":
		// Return example IBAN numbers
		examples := []map[string]interface{}{
			{
				"iban":        "GB82WEST12345698765432",
				"country":     "GB",
				"valid":       true,
				"description": "Example valid UK IBAN",
			},
			{
				"iban":        "DE89370400440532013000",
				"country":     "DE",
				"valid":       true,
				"description": "Example valid German IBAN",
			},
			{
				"iban":        "FR1420041010050500013M02606",
				"country":     "FR",
				"valid":       true,
				"description": "Example valid French IBAN",
			},
			{
				"iban":        "GB82WEST12345698765433",
				"country":     "GB",
				"valid":       false,
				"description": "Example invalid UK IBAN (wrong check digits)",
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

// validateIBAN validates an IBAN number
func (i *IBANTool) validateIBAN(params map[string]interface{}) (interface{}, error) {
	input, _ := params["input"].(string)
	if input == "" {
		return nil, fmt.Errorf("input parameter is required for validation")
	}

	// Clean the input (remove spaces and convert to uppercase)
	cleanInput := strings.ToUpper(strings.ReplaceAll(input, " ", ""))

	// Check basic format (at least 15 characters, starts with 2 letters)
	if len(cleanInput) < 15 || len(cleanInput) > 34 {
		return map[string]interface{}{
			"valid": false,
			"error": "IBAN must be between 15 and 34 characters",
			"input": input,
		}, nil
	}

	// Check if it starts with 2 letters
	if len(cleanInput) < 2 || !isAlpha(cleanInput[:2]) {
		return map[string]interface{}{
			"valid": false,
			"error": "IBAN must start with 2 letters (country code)",
			"input": input,
		}, nil
	}

	// Check if the rest are alphanumeric
	if len(cleanInput) > 2 && !isAlphanumeric(cleanInput[2:]) {
		return map[string]interface{}{
			"valid": false,
			"error": "IBAN must contain only letters and numbers",
			"input": input,
		}, nil
	}

	// Extract country code
	countryCode := cleanInput[:2]

	// Validate using MOD-97 algorithm
	if !i.mod97Check(cleanInput) {
		return map[string]interface{}{
			"valid": false,
			"error": "invalid check digits",
			"input": input,
		}, nil
	}

	return map[string]interface{}{
		"valid":   true,
		"iban":    cleanInput,
		"country": countryCode,
		"input":   input,
	}, nil
}

// generateIBAN generates IBAN numbers
func (i *IBANTool) generateIBAN(params map[string]interface{}) (interface{}, error) {
	count := 1
	if c, ok := params["count"].(float64); ok {
		count = int(c)
	}

	countryCode := ""
	if cc, ok := params["country-code"].(string); ok {
		countryCode = cc
	}

	if count == 1 {
		iban, err := i.generateSingleIBAN(countryCode)
		if err != nil {
			return nil, err
		}
		return iban, nil
	}

	ibans := make([]string, count)
	for j := range count {
		iban, err := i.generateSingleIBAN(countryCode)
		if err != nil {
			return nil, err
		}
		ibans[j] = iban
	}

	return ibans, nil
}

// generateSingleIBAN generates a single IBAN number
func (i *IBANTool) generateSingleIBAN(countryCode string) (string, error) {
	if countryCode == "" {
		// Random country code
		countryCode = i.getRandomCountryCode()
	}

	// Get the expected length for this country
	expectedLength := i.getCountryIBANLength(countryCode)
	if expectedLength == 0 {
		return "", fmt.Errorf("unsupported country code: %s", countryCode)
	}

	// Generate random BBAN (Basic Bank Account Number)
	// Length is total - 4 (2 for country code + 2 for check digits)
	bbanLength := expectedLength - 4
	bban := i.generateRandomBBAN(bbanLength)

	// Create the IBAN without check digits
	ibanWithoutChecks := countryCode + "00" + bban

	// Calculate check digits
	checkDigits := i.calculateCheckDigits(ibanWithoutChecks)

	// Return the complete IBAN
	return countryCode + checkDigits + bban, nil
}

// mod97Check validates an IBAN using the MOD-97 algorithm
func (i *IBANTool) mod97Check(iban string) bool {
	// Move first 4 characters to the end
	rearranged := iban[4:] + iban[:4]

	// Convert to numeric string
	numeric := i.lettersToNumbers(rearranged)

	// Calculate remainder using big integer arithmetic
	remainder := i.calculateMod97(numeric)

	return remainder == 1
}

// calculateMod97 calculates the remainder when dividing by 97
func (i *IBANTool) calculateMod97(numeric string) int {
	remainder := 0
	for _, char := range numeric {
		digit := int(char - '0')
		remainder = (remainder*10 + digit) % 97
	}
	return remainder
}

// lettersToNumbers converts letters to their numeric values (A=10, B=11, ..., Z=35)
func (i *IBANTool) lettersToNumbers(s string) string {
	var result strings.Builder
	for _, char := range s {
		if char >= 'A' && char <= 'Z' {
			// Convert letter to number (A=10, B=11, ..., Z=35)
			value := int(char - 'A' + 10)
			result.WriteString(strconv.Itoa(value))
		} else {
			// Keep digit as is
			result.WriteRune(char)
		}
	}
	return result.String()
}

// calculateCheckDigits calculates the check digits for an IBAN
func (i *IBANTool) calculateCheckDigits(ibanWithoutChecks string) string {
	// Move first 4 characters to the end
	rearranged := ibanWithoutChecks[4:] + ibanWithoutChecks[:4]

	// Convert to numeric string
	numeric := i.lettersToNumbers(rearranged)

	// Calculate remainder
	remainder := i.calculateMod97(numeric)

	// Check digits are 98 - remainder
	checkDigits := 98 - remainder

	// Format as 2-digit string
	return fmt.Sprintf("%02d", checkDigits)
}

// generateRandomBBAN generates a random Basic Bank Account Number
func (i *IBANTool) generateRandomBBAN(length int) string {
	var result strings.Builder
	for j := 0; j < length; j++ {
		// Generate random alphanumeric character
		if rand.Intn(2) == 0 {
			// Generate random digit
			result.WriteString(strconv.Itoa(rand.Intn(10)))
		} else {
			// Generate random letter
			result.WriteRune(rune('A' + rand.Intn(26)))
		}
	}
	return result.String()
}

// populateCountries initializes the countries data
func (i *IBANTool) populateCountries() {
	i.countries = []IBANCountry{
		{"GB", "United Kingdom", 22},
		{"DE", "Germany", 22},
		{"FR", "France", 27},
		{"IT", "Italy", 27},
		{"ES", "Spain", 24},
		{"NL", "Netherlands", 18},
		{"BE", "Belgium", 16},
		{"AT", "Austria", 20},
		{"CH", "Switzerland", 21},
		{"SE", "Sweden", 24},
		{"NO", "Norway", 15},
		{"DK", "Denmark", 18},
		{"FI", "Finland", 18},
		{"PL", "Poland", 28},
		{"CZ", "Czech Republic", 24},
		{"HU", "Hungary", 28},
		{"RO", "Romania", 24},
		{"BG", "Bulgaria", 22},
		{"HR", "Croatia", 21},
		{"SI", "Slovenia", 19},
		{"SK", "Slovakia", 24},
		{"LT", "Lithuania", 20},
		{"LV", "Latvia", 21},
		{"EE", "Estonia", 20},
		{"IE", "Ireland", 22},
		{"PT", "Portugal", 25},
		{"GR", "Greece", 27},
		{"CY", "Cyprus", 28},
		{"MT", "Malta", 31},
		{"LU", "Luxembourg", 20},
	}
}

// getCountryIBANLength returns the expected IBAN length for a country
func (i *IBANTool) getCountryIBANLength(countryCode string) int {
	for _, country := range i.countries {
		if country.Code == countryCode {
			return country.Length
		}
	}
	return 0
}

// isValidCountryCode checks if a country code is valid
func (i *IBANTool) isValidCountryCode(countryCode string) bool {
	return i.getCountryIBANLength(countryCode) > 0
}

// getCountriesData returns the countries data for JSON serialization
func (i *IBANTool) getCountriesData() []map[string]interface{} {
	var countriesData []map[string]interface{}
	for _, country := range i.countries {
		countriesData = append(countriesData, map[string]interface{}{
			"code":   country.Code,
			"name":   country.Name,
			"length": country.Length,
		})
	}
	return countriesData
}

// getRandomCountryCode returns a random country code from the supported countries
func (i *IBANTool) getRandomCountryCode() string {
	if len(i.countries) == 0 {
		return "GB" // Fallback to UK if no countries loaded
	}
	return i.countries[rand.Intn(len(i.countries))].Code
}

// isAlpha checks if a string contains only letters
func isAlpha(s string) bool {
	for _, char := range s {
		if char < 'A' || char > 'Z' {
			return false
		}
	}
	return true
}

// isAlphanumeric checks if a string contains only letters and numbers
func isAlphanumeric(s string) bool {
	for _, char := range s {
		if !((char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9')) {
			return false
		}
	}
	return true
}
