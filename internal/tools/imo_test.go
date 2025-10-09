package tools

import (
	"fmt"
	"testing"
)

func TestIMOToolName(t *testing.T) {
	tool := NewIMOTool()
	if tool.Name() != "imo" {
		t.Errorf("Expected name 'imo', got '%s'", tool.Name())
	}
}

func TestIMOToolDescription(t *testing.T) {
	tool := NewIMOTool()
	desc := tool.Description()
	if desc == "" {
		t.Error("Description should not be empty")
	}
	if len(desc) < 50 {
		t.Error("Description should be more detailed")
	}
}

func TestIMOToolValidateParams(t *testing.T) {
	tool := NewIMOTool()

	tests := []struct {
		name     string
		params   map[string]interface{}
		expected error
	}{
		{
			name:     "valid_validate_operation",
			params:   map[string]interface{}{"operation": "validate", "input": "1234567"},
			expected: nil,
		},
		{
			name:     "valid_generate_operation",
			params:   map[string]interface{}{"operation": "generate", "count": 5},
			expected: nil,
		},
		{
			name:     "invalid_operation",
			params:   map[string]interface{}{"operation": "invalid"},
			expected: fmt.Errorf("operation must be 'validate' or 'generate'"),
		},
		{
			name:     "count_too_low",
			params:   map[string]interface{}{"operation": "generate", "count": 0},
			expected: fmt.Errorf("count must be at least 1"),
		},
		{
			name:     "count_too_high",
			params:   map[string]interface{}{"operation": "generate", "count": 101},
			expected: fmt.Errorf("count cannot exceed 100"),
		},
		{
			name:     "validate_missing_input",
			params:   map[string]interface{}{"operation": "validate"},
			expected: fmt.Errorf("input parameter is required for validation"),
		},
		{
			name:     "validate_input_not_string",
			params:   map[string]interface{}{"operation": "validate", "input": 123},
			expected: fmt.Errorf("input must be a string"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tool.ValidateParams(tt.params)
			if (err == nil) != (tt.expected == nil) {
				t.Errorf("Expected error %v, got %v", tt.expected, err)
			}
			if err != nil && tt.expected != nil && err.Error() != tt.expected.Error() {
				t.Errorf("Expected error message '%s', got '%s'", tt.expected.Error(), err.Error())
			}
		})
	}
}

func TestIMOToolValidateIMO(t *testing.T) {
	tool := NewIMOTool()

	tests := []struct {
		name     string
		input    string
		expected map[string]interface{}
	}{
		{
			name:  "valid_imo",
			input: "1234567",
			expected: map[string]interface{}{
				"valid": true,
				"imo":   "1234567",
				"input": "1234567",
			},
		},
		{
			name:  "valid_imo_with_spaces",
			input: "123 456 7",
			expected: map[string]interface{}{
				"valid": true,
				"imo":   "1234567",
				"input": "123 456 7",
			},
		},
		{
			name:  "valid_imo_with_dashes",
			input: "123-456-7",
			expected: map[string]interface{}{
				"valid": true,
				"imo":   "1234567",
				"input": "123-456-7",
			},
		},
		{
			name:  "too_short",
			input: "123456",
			expected: map[string]interface{}{
				"valid": false,
				"error": "IMO number must be exactly 7 digits",
				"input": "123456",
			},
		},
		{
			name:  "too_long",
			input: "12345678",
			expected: map[string]interface{}{
				"valid": false,
				"error": "IMO number must be exactly 7 digits",
				"input": "12345678",
			},
		},
		{
			name:  "non_numeric",
			input: "123456a",
			expected: map[string]interface{}{
				"valid": false,
				"error": "IMO number must contain only digits",
				"input": "123456a",
			},
		},
		{
			name:  "invalid_check_digit",
			input: "1234568", // Should be 7 for valid check digit
			expected: map[string]interface{}{
				"valid": false,
				"error": "invalid check digit. Expected 7, got 8",
				"input": "1234568",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := map[string]interface{}{
				"operation": "validate",
				"input":     tt.input,
			}

			result, err := tool.Execute(params)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			resultMap, ok := result.(map[string]interface{})
			if !ok {
				t.Fatalf("Expected map result, got %T", result)
			}

			// Check each expected field
			for key, expectedValue := range tt.expected {
				if actualValue, exists := resultMap[key]; !exists {
					t.Errorf("Missing field '%s' in result", key)
				} else if actualValue != expectedValue {
					t.Errorf("Field '%s': expected %v, got %v", key, expectedValue, actualValue)
				}
			}
		})
	}
}

func TestIMOToolGenerateIMO(t *testing.T) {
	tool := NewIMOTool()

	tests := []struct {
		name     string
		params   map[string]interface{}
		validate func(interface{}) error
	}{
		{
			name:   "generate_single",
			params: map[string]interface{}{"operation": "generate"},
			validate: func(result interface{}) error {
				imo, ok := result.(string)
				if !ok {
					return fmt.Errorf("Expected string result, got %T", result)
				}
				if len(imo) != 7 {
					return fmt.Errorf("Expected 7-digit IMO, got length %d", len(imo))
				}
				// Validate the generated IMO
				validateParams := map[string]interface{}{
					"operation": "validate",
					"input":     imo,
				}
				validateResult, err := tool.Execute(validateParams)
				if err != nil {
					return fmt.Errorf("Generated IMO validation failed: %v", err)
				}
				validateMap := validateResult.(map[string]interface{})
				if !validateMap["valid"].(bool) {
					return fmt.Errorf("Generated IMO is not valid: %v", validateMap["error"])
				}
				return nil
			},
		},
		{
			name:   "generate_multiple",
			params: map[string]interface{}{"operation": "generate", "count": 5},
			validate: func(result interface{}) error {
				imos, ok := result.([]string)
				if !ok {
					return fmt.Errorf("Expected []string result, got %T", result)
				}
				if len(imos) != 5 {
					return fmt.Errorf("Expected 5 IMOs, got %d", len(imos))
				}
				// Validate each generated IMO
				for i, imo := range imos {
					if len(imo) != 7 {
						return fmt.Errorf("IMO %d: expected 7 digits, got length %d", i, len(imo))
					}
					validateParams := map[string]interface{}{
						"operation": "validate",
						"input":     imo,
					}
					validateResult, err := tool.Execute(validateParams)
					if err != nil {
						return fmt.Errorf("IMO %d validation failed: %v", i, err)
					}
					validateMap := validateResult.(map[string]interface{})
					if !validateMap["valid"].(bool) {
						return fmt.Errorf("IMO %d is not valid: %v", i, validateMap["error"])
					}
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tool.Execute(tt.params)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if err := tt.validate(result); err != nil {
				t.Errorf("Validation failed: %v", err)
			}
		})
	}
}

func TestIMOToolEdgeCases(t *testing.T) {
	tool := NewIMOTool()

	tests := []struct {
		name     string
		params   map[string]interface{}
		expected error
	}{
		{
			name:     "count_exceeds_limit",
			params:   map[string]interface{}{"operation": "generate", "count": 101},
			expected: fmt.Errorf("count cannot exceed 100"),
		},
		{
			name:     "empty_input",
			params:   map[string]interface{}{"operation": "validate", "input": ""},
			expected: fmt.Errorf("input parameter is required for validation"),
		},
		{
			name:     "missing_operation_defaults_to_validate",
			params:   map[string]interface{}{"input": "1234567"},
			expected: nil, // Should work with default operation
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tool.Execute(tt.params)
			if (err == nil) != (tt.expected == nil) {
				t.Errorf("Expected error %v, got %v", tt.expected, err)
			}
			if err != nil && tt.expected != nil && err.Error() != tt.expected.Error() {
				t.Errorf("Expected error message '%s', got '%s'", tt.expected.Error(), err.Error())
			}
		})
	}
}

func TestIMOToolSchema(t *testing.T) {
	tool := NewIMOTool()

	// Test input schema
	inputSchema := tool.GetInputSchema()
	if inputSchema == nil {
		t.Error("Input schema should not be nil")
	}

	// Test output schema
	outputSchema := tool.GetOutputSchema()
	if outputSchema == nil {
		t.Error("Output schema should not be nil")
	}

	// Validate schema structure
	if schemaType, ok := inputSchema["type"]; !ok || schemaType != "object" {
		t.Error("Input schema should have type 'object'")
	}

	if schemaType, ok := outputSchema["type"]; !ok || schemaType != "object" {
		t.Error("Output schema should have type 'object'")
	}
}

func TestIMOToolResources(t *testing.T) {
	tool := NewIMOTool()

	// Test GetResources
	resources := tool.GetResources()
	if len(resources) != 2 {
		t.Errorf("Expected 2 resources, got %d", len(resources))
	}

	// Test resource names and URIs
	expectedResources := map[string]string{
		"IMO Algorithm": "imo://algorithm",
		"IMO Examples":  "imo://examples",
	}

	for _, resource := range resources {
		expectedURI, exists := expectedResources[resource.Name]
		if !exists {
			t.Errorf("Unexpected resource name: %s", resource.Name)
		}
		if resource.URI != expectedURI {
			t.Errorf("Expected URI %s, got %s", expectedURI, resource.URI)
		}
		if resource.MIMEType != "application/json" {
			t.Errorf("Expected MIME type 'application/json', got %s", resource.MIMEType)
		}
	}

	// Test ReadResource for algorithm
	algorithmContent, err := tool.ReadResource("imo://algorithm")
	if err != nil {
		t.Errorf("ReadResource(imo://algorithm) failed: %v", err)
	}
	if algorithmContent == "" {
		t.Error("ReadResource(imo://algorithm) returned empty content")
	}

	// Test ReadResource for examples
	examplesContent, err := tool.ReadResource("imo://examples")
	if err != nil {
		t.Errorf("ReadResource(imo://examples) failed: %v", err)
	}
	if examplesContent == "" {
		t.Error("ReadResource(imo://examples) returned empty content")
	}

	// Test ReadResource with unknown URI
	_, err = tool.ReadResource("imo://unknown")
	if err == nil {
		t.Error("ReadResource with unknown URI should return error")
	}
}
