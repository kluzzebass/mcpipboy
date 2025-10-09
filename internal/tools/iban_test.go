package tools

import (
	"testing"
)

func TestIBANTool_Name(t *testing.T) {
	tool := NewIBANTool()
	if tool.Name() != "iban" {
		t.Errorf("Expected name 'iban', got '%s'", tool.Name())
	}
}

func TestIBANTool_Description(t *testing.T) {
	tool := NewIBANTool()
	desc := tool.Description()
	if desc == "" {
		t.Error("Description should not be empty")
	}
}

func TestIBANTool_ValidateParams(t *testing.T) {
	tool := NewIBANTool()

	tests := []struct {
		name     string
		params   map[string]interface{}
		expected string
	}{
		{
			name:     "valid validate operation",
			params:   map[string]interface{}{"operation": "validate", "input": "GB82WEST12345698765432"},
			expected: "",
		},
		{
			name:     "valid generate operation",
			params:   map[string]interface{}{"operation": "generate", "country-code": "GB", "count": 5.0},
			expected: "",
		},
		{
			name:     "invalid operation",
			params:   map[string]interface{}{"operation": "invalid"},
			expected: "invalid operation: invalid. Supported operations: validate, generate",
		},
		{
			name:     "missing input for validate",
			params:   map[string]interface{}{"operation": "validate"},
			expected: "input parameter is required for validation",
		},
		{
			name:     "invalid count too low",
			params:   map[string]interface{}{"operation": "generate", "count": 0.0},
			expected: "count must be between 1 and 100",
		},
		{
			name:     "invalid count too high",
			params:   map[string]interface{}{"operation": "generate", "count": 101.0},
			expected: "count must be between 1 and 100",
		},
		{
			name:     "invalid country code",
			params:   map[string]interface{}{"operation": "generate", "country-code": "XX"},
			expected: "invalid country code: XX. Must be a valid ISO 3166-1 alpha-2 country code",
		},
		{
			name:     "valid country code",
			params:   map[string]interface{}{"operation": "generate", "country-code": "GB"},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tool.ValidateParams(tt.params)
			if tt.expected == "" {
				if err != nil {
					t.Errorf("Expected no error, got: %v", err)
				}
			} else {
				if err == nil || err.Error() != tt.expected {
					t.Errorf("Expected error '%s', got: %v", tt.expected, err)
				}
			}
		})
	}
}

func TestIBANTool_Execute_Validate(t *testing.T) {
	tool := NewIBANTool()

	tests := []struct {
		name     string
		params   map[string]interface{}
		expected map[string]interface{}
		hasError bool
	}{
		{
			name:   "valid UK IBAN",
			params: map[string]interface{}{"operation": "validate", "input": "GB82WEST12345698765432"},
			expected: map[string]interface{}{
				"valid":   true,
				"country": "GB",
			},
			hasError: false,
		},
		{
			name:   "valid German IBAN",
			params: map[string]interface{}{"operation": "validate", "input": "DE89370400440532013000"},
			expected: map[string]interface{}{
				"valid":   true,
				"country": "DE",
			},
			hasError: false,
		},
		{
			name:   "valid French IBAN",
			params: map[string]interface{}{"operation": "validate", "input": "FR1420041010050500013M02606"},
			expected: map[string]interface{}{
				"valid":   true,
				"country": "FR",
			},
			hasError: false,
		},
		{
			name:   "invalid IBAN - wrong check digits",
			params: map[string]interface{}{"operation": "validate", "input": "GB82WEST12345698765433"},
			expected: map[string]interface{}{
				"valid": false,
				"error": "invalid check digits",
			},
			hasError: false,
		},
		{
			name:   "invalid IBAN - too short",
			params: map[string]interface{}{"operation": "validate", "input": "GB82WEST"},
			expected: map[string]interface{}{
				"valid": false,
				"error": "IBAN must be between 15 and 34 characters",
			},
			hasError: false,
		},
		{
			name:   "invalid IBAN - doesn't start with letters",
			params: map[string]interface{}{"operation": "validate", "input": "8282WEST12345698765432"},
			expected: map[string]interface{}{
				"valid": false,
				"error": "IBAN must start with 2 letters (country code)",
			},
			hasError: false,
		},
		{
			name:   "invalid IBAN - contains invalid characters",
			params: map[string]interface{}{"operation": "validate", "input": "GB82WEST12345698765432!"},
			expected: map[string]interface{}{
				"valid": false,
				"error": "IBAN must contain only letters and numbers",
			},
			hasError: false,
		},
		{
			name:   "IBAN with spaces",
			params: map[string]interface{}{"operation": "validate", "input": "GB82 WEST 1234 5698 7654 32"},
			expected: map[string]interface{}{
				"valid":   true,
				"country": "GB",
			},
			hasError: false,
		},
		{
			name:   "IBAN with lowercase",
			params: map[string]interface{}{"operation": "validate", "input": "gb82west12345698765432"},
			expected: map[string]interface{}{
				"valid":   true,
				"country": "GB",
			},
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tool.Execute(tt.params)
			if tt.hasError {
				if err == nil {
					t.Error("Expected error, got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			resultMap, ok := result.(map[string]interface{})
			if !ok {
				t.Errorf("Expected map result, got %T", result)
				return
			}

			// Check expected fields
			for key, expectedValue := range tt.expected {
				if actualValue, exists := resultMap[key]; !exists {
					t.Errorf("Expected field '%s' not found in result", key)
				} else if actualValue != expectedValue {
					t.Errorf("Expected field '%s' to be '%v', got '%v'", key, expectedValue, actualValue)
				}
			}
		})
	}
}

func TestIBANTool_Execute_Generate(t *testing.T) {
	tool := NewIBANTool()

	tests := []struct {
		name     string
		params   map[string]interface{}
		expected string
		hasError bool
	}{
		{
			name:     "generate single IBAN",
			params:   map[string]interface{}{"operation": "generate"},
			expected: "string",
			hasError: false,
		},
		{
			name:     "generate single IBAN with country",
			params:   map[string]interface{}{"operation": "generate", "country-code": "GB"},
			expected: "string",
			hasError: false,
		},
		{
			name:     "generate multiple IBANs",
			params:   map[string]interface{}{"operation": "generate", "count": 3.0},
			expected: "[]string",
			hasError: false,
		},
		{
			name:     "generate with invalid country",
			params:   map[string]interface{}{"operation": "generate", "country-code": "XX"},
			expected: "",
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tool.Execute(tt.params)
			if tt.hasError {
				if err == nil {
					t.Error("Expected error, got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Check result type
			switch tt.expected {
			case "string":
				if _, ok := result.(string); !ok {
					t.Errorf("Expected string result, got %T", result)
				}
			case "[]string":
				if _, ok := result.([]string); !ok {
					t.Errorf("Expected []string result, got %T", result)
				}
			}
		})
	}
}

func TestIBANTool_GetInputSchema(t *testing.T) {
	tool := NewIBANTool()
	schema := tool.GetInputSchema()

	if schema["type"] != "object" {
		t.Errorf("Expected schema type 'object', got '%v'", schema["type"])
	}

	properties, ok := schema["properties"].(map[string]interface{})
	if !ok {
		t.Error("Expected properties to be a map")
	}

	// Check required properties exist
	requiredProps := []string{"operation", "input", "country-code", "count"}
	for _, prop := range requiredProps {
		if _, exists := properties[prop]; !exists {
			t.Errorf("Expected property '%s' in schema", prop)
		}
	}
}

func TestIBANTool_GetOutputSchema(t *testing.T) {
	tool := NewIBANTool()
	schema := tool.GetOutputSchema()

	if schema["type"] != "object" {
		t.Errorf("Expected schema type 'object', got '%v'", schema["type"])
	}

	properties, ok := schema["properties"].(map[string]interface{})
	if !ok {
		t.Error("Expected properties to be a map")
	}

	// Check expected output fields exist
	expectedFields := []string{"valid", "iban", "ibans", "country", "error"}
	for _, field := range expectedFields {
		if _, exists := properties[field]; !exists {
			t.Errorf("Expected field '%s' in output schema", field)
		}
	}
}

func TestIBANTool_GetResources(t *testing.T) {
	tool := NewIBANTool()
	resources := tool.GetResources()

	if len(resources) == 0 {
		t.Error("Expected at least one resource")
	}

	expectedURIs := []string{"iban://countries", "iban://mod97", "iban://examples"}
	for _, uri := range expectedURIs {
		found := false
		for _, resource := range resources {
			if resource.URI == uri {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected resource URI '%s' not found", uri)
		}
	}
}

func TestIBANTool_ReadResource(t *testing.T) {
	tool := NewIBANTool()

	tests := []struct {
		uri      string
		hasError bool
	}{
		{"iban://countries", false},
		{"iban://mod97", false},
		{"iban://examples", false},
		{"iban://invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.uri, func(t *testing.T) {
			content, err := tool.ReadResource(tt.uri)
			if tt.hasError {
				if err == nil {
					t.Error("Expected error, got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if content == "" {
				t.Error("Expected non-empty content")
			}
		})
	}
}

func TestIBANTool_mod97Check(t *testing.T) {
	tool := NewIBANTool()

	tests := []struct {
		iban  string
		valid bool
	}{
		{"GB82WEST12345698765432", true},
		{"DE89370400440532013000", true},
		{"FR1420041010050500013M02606", true},
		{"GB82WEST12345698765433", false}, // Wrong check digits
		{"GB82WEST12345698765431", false}, // Wrong check digits
	}

	for _, tt := range tests {
		t.Run(tt.iban, func(t *testing.T) {
			result := tool.mod97Check(tt.iban)
			if result != tt.valid {
				t.Errorf("Expected %v for IBAN %s, got %v", tt.valid, tt.iban, result)
			}
		})
	}
}

func TestIBANTool_lettersToNumbers(t *testing.T) {
	tool := NewIBANTool()

	tests := []struct {
		input    string
		expected string
	}{
		{"A", "10"},
		{"B", "11"},
		{"Z", "35"},
		{"0", "0"},
		{"9", "9"},
		{"A0B1C2", "100111122"},
		{"WEST12345698765432GB82", "3214282912345698765432161182"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := tool.lettersToNumbers(tt.input)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestIBANTool_calculateMod97(t *testing.T) {
	tool := NewIBANTool()

	tests := []struct {
		input    string
		expected int
	}{
		{"3214282912345698765432161182", 1}, // Valid UK IBAN
		{"3214282912345698765432161183", 2}, // Invalid UK IBAN
		{"123456789012345678901234567890", 52},
		{"97", 0},
		{"98", 1},
		{"99", 2},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := tool.calculateMod97(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestIBANTool_calculateCheckDigits(t *testing.T) {
	tool := NewIBANTool()

	tests := []struct {
		input    string
		expected string
	}{
		{"GB00WEST12345698765432", "82"},     // Should produce check digits 82
		{"DE00DE89370400440532013000", "16"}, // Should produce check digits 16
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := tool.calculateCheckDigits(tt.input)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestIBANTool_getCountryIBANLength(t *testing.T) {
	tool := NewIBANTool()

	tests := []struct {
		country  string
		expected int
	}{
		{"GB", 22},
		{"DE", 22},
		{"FR", 27},
		{"IT", 27},
		{"ES", 24},
		{"NL", 18},
		{"BE", 16},
		{"AT", 20},
		{"CH", 21},
		{"SE", 24},
		{"NO", 15},
		{"DK", 18},
		{"FI", 18},
		{"PL", 28},
		{"CZ", 24},
		{"HU", 28},
		{"RO", 24},
		{"BG", 22},
		{"HR", 21},
		{"SI", 19},
		{"SK", 24},
		{"LT", 20},
		{"LV", 21},
		{"EE", 20},
		{"IE", 22},
		{"PT", 25},
		{"GR", 27},
		{"CY", 28},
		{"MT", 31},
		{"LU", 20},
		{"XX", 0}, // Invalid country
	}

	for _, tt := range tests {
		t.Run(tt.country, func(t *testing.T) {
			result := tool.getCountryIBANLength(tt.country)
			if result != tt.expected {
				t.Errorf("Expected %d for country %s, got %d", tt.expected, tt.country, result)
			}
		})
	}
}

func TestIBANTool_isValidCountryCode(t *testing.T) {
	tool := NewIBANTool()

	tests := []struct {
		country  string
		expected bool
	}{
		{"GB", true},
		{"DE", true},
		{"FR", true},
		{"IT", true},
		{"ES", true},
		{"NL", true},
		{"BE", true},
		{"AT", true},
		{"CH", true},
		{"SE", true},
		{"NO", true},
		{"DK", true},
		{"FI", true},
		{"PL", true},
		{"CZ", true},
		{"HU", true},
		{"RO", true},
		{"BG", true},
		{"HR", true},
		{"SI", true},
		{"SK", true},
		{"LT", true},
		{"LV", true},
		{"EE", true},
		{"IE", true},
		{"PT", true},
		{"GR", true},
		{"CY", true},
		{"MT", true},
		{"LU", true},
		{"XX", false}, // Invalid country
		{"", false},   // Empty country
	}

	for _, tt := range tests {
		t.Run(tt.country, func(t *testing.T) {
			result := tool.isValidCountryCode(tt.country)
			if result != tt.expected {
				t.Errorf("Expected %v for country %s, got %v", tt.expected, tt.country, result)
			}
		})
	}
}

func TestIBANTool_generateRandomBBAN(t *testing.T) {
	tool := NewIBANTool()

	tests := []struct {
		length int
	}{
		{10},
		{15},
		{20},
		{25},
		{30},
	}

	for _, tt := range tests {
		t.Run("length_"+string(rune(tt.length)), func(t *testing.T) {
			result := tool.generateRandomBBAN(tt.length)
			if len(result) != tt.length {
				t.Errorf("Expected length %d, got %d", tt.length, len(result))
			}

			// Check that all characters are alphanumeric
			for _, char := range result {
				if !((char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9')) {
					t.Errorf("Expected alphanumeric character, got '%c'", char)
				}
			}
		})
	}
}

func TestIBANTool_generateSingleIBAN(t *testing.T) {
	tool := NewIBANTool()

	tests := []struct {
		country  string
		hasError bool
	}{
		{"GB", false},
		{"DE", false},
		{"FR", false},
		{"IT", false},
		{"ES", false},
		{"NL", false},
		{"BE", false},
		{"AT", false},
		{"CH", false},
		{"SE", false},
		{"NO", false},
		{"DK", false},
		{"FI", false},
		{"PL", false},
		{"CZ", false},
		{"HU", false},
		{"RO", false},
		{"BG", false},
		{"HR", false},
		{"SI", false},
		{"SK", false},
		{"LT", false},
		{"LV", false},
		{"EE", false},
		{"IE", false},
		{"PT", false},
		{"GR", false},
		{"CY", false},
		{"MT", false},
		{"LU", false},
		{"XX", true}, // Invalid country
	}

	for _, tt := range tests {
		t.Run(tt.country, func(t *testing.T) {
			result, err := tool.generateSingleIBAN(tt.country)
			if tt.hasError {
				if err == nil {
					t.Error("Expected error, got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result == "" {
				t.Error("Expected non-empty IBAN")
			}

			// Check that the IBAN starts with the country code
			if len(result) < 2 || result[:2] != tt.country {
				t.Errorf("Expected IBAN to start with country code %s, got %s", tt.country, result[:2])
			}

			// Check that the IBAN is valid
			if !tool.mod97Check(result) {
				t.Errorf("Generated IBAN %s is not valid", result)
			}
		})
	}
}

func TestIBANTool_GenerateAndValidate(t *testing.T) {
	tool := NewIBANTool()

	// Test that generated IBANs are valid
	for i := 0; i < 10; i++ {
		iban, err := tool.generateSingleIBAN("GB")
		if err != nil {
			t.Errorf("Unexpected error generating IBAN: %v", err)
			continue
		}

		// Validate the generated IBAN
		params := map[string]interface{}{
			"operation": "validate",
			"input":     iban,
		}

		result, err := tool.Execute(params)
		if err != nil {
			t.Errorf("Unexpected error validating generated IBAN: %v", err)
			continue
		}

		resultMap, ok := result.(map[string]interface{})
		if !ok {
			t.Errorf("Expected map result, got %T", result)
			continue
		}

		if valid, ok := resultMap["valid"].(bool); !ok || !valid {
			t.Errorf("Generated IBAN %s is not valid", iban)
		}
	}
}
