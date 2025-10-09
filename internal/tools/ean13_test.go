package tools

import (
	"strings"
	"testing"
)

func TestEAN13ToolExecute(t *testing.T) {
	tool := NewEAN13Tool()

	tests := []struct {
		name     string
		params   map[string]interface{}
		wantErr  bool
		validate func(t *testing.T, result interface{})
	}{
		{
			name: "validate_valid_ean13",
			params: map[string]interface{}{
				"operation": "validate",
				"input":     "1234567890128",
			},
			wantErr: false,
			validate: func(t *testing.T, result interface{}) {
				if result == nil {
					t.Error("Expected non-nil result")
				}
				if resultMap, ok := result.(map[string]interface{}); ok {
					if valid, ok := resultMap["valid"].(bool); !ok || !valid {
						t.Errorf("Expected valid=true, got: %v", resultMap["valid"])
					}
				} else {
					t.Errorf("Expected map result, got %T", result)
				}
			},
		},
		{
			name: "validate_valid_ean13_with_spaces",
			params: map[string]interface{}{
				"operation": "validate",
				"input":     "123 456 789 012 8",
			},
			wantErr: false,
			validate: func(t *testing.T, result interface{}) {
				if result == nil {
					t.Error("Expected non-nil result")
				}
				if resultMap, ok := result.(map[string]interface{}); ok {
					if valid, ok := resultMap["valid"].(bool); !ok || !valid {
						t.Errorf("Expected valid=true, got: %v", resultMap["valid"])
					}
				} else {
					t.Errorf("Expected map result, got %T", result)
				}
			},
		},
		{
			name: "validate_valid_ean13_with_dashes",
			params: map[string]interface{}{
				"operation": "validate",
				"input":     "123-456-789-012-8",
			},
			wantErr: false,
			validate: func(t *testing.T, result interface{}) {
				if result == nil {
					t.Error("Expected non-nil result")
				}
				if resultMap, ok := result.(map[string]interface{}); ok {
					if valid, ok := resultMap["valid"].(bool); !ok || !valid {
						t.Errorf("Expected valid=true, got: %v", resultMap["valid"])
					}
				} else {
					t.Errorf("Expected map result, got %T", result)
				}
			},
		},
		{
			name: "validate_invalid_ean13",
			params: map[string]interface{}{
				"operation": "validate",
				"input":     "1234567890123",
			},
			wantErr: false,
			validate: func(t *testing.T, result interface{}) {
				if result == nil {
					t.Error("Expected non-nil result")
				}
				if resultMap, ok := result.(map[string]interface{}); ok {
					if valid, ok := resultMap["valid"].(bool); !ok || valid {
						t.Errorf("Expected valid=false, got: %v", resultMap["valid"])
					}
					if _, ok := resultMap["error"]; !ok {
						t.Error("Expected error field in result")
					}
				} else {
					t.Errorf("Expected map result, got %T", result)
				}
			},
		},
		{
			name: "validate_invalid_length",
			params: map[string]interface{}{
				"operation": "validate",
				"input":     "123456789012",
			},
			wantErr: false,
			validate: func(t *testing.T, result interface{}) {
				if result == nil {
					t.Error("Expected non-nil result")
				}
				if resultMap, ok := result.(map[string]interface{}); ok {
					if valid, ok := resultMap["valid"].(bool); !ok || valid {
						t.Errorf("Expected valid=false, got: %v", resultMap["valid"])
					}
					if errorMsg, ok := resultMap["error"].(string); !ok || !strings.Contains(errorMsg, "exactly 13 characters") {
						t.Errorf("Expected length error, got: %v", resultMap["error"])
					}
				} else {
					t.Errorf("Expected map result, got %T", result)
				}
			},
		},
		{
			name: "validate_non_numeric",
			params: map[string]interface{}{
				"operation": "validate",
				"input":     "123456789012a",
			},
			wantErr: false,
			validate: func(t *testing.T, result interface{}) {
				if result == nil {
					t.Error("Expected non-nil result")
				}
				if resultMap, ok := result.(map[string]interface{}); ok {
					if valid, ok := resultMap["valid"].(bool); !ok || valid {
						t.Errorf("Expected valid=false, got: %v", resultMap["valid"])
					}
					if errorMsg, ok := resultMap["error"].(string); !ok || !strings.Contains(errorMsg, "only digits") {
						t.Errorf("Expected digits error, got: %v", resultMap["error"])
					}
				} else {
					t.Errorf("Expected map result, got %T", result)
				}
			},
		},
		{
			name: "generate_single_ean13",
			params: map[string]interface{}{
				"operation": "generate",
			},
			wantErr: false,
			validate: func(t *testing.T, result interface{}) {
				if result == nil {
					t.Error("Expected non-nil result")
				}
				if ean13, ok := result.(string); ok {
					if len(ean13) != 13 {
						t.Errorf("Expected 13-character EAN-13, got length %d", len(ean13))
					}
					// Validate the generated EAN-13
					validateParams := map[string]interface{}{
						"operation": "validate",
						"input":     ean13,
					}
					validateResult, err := tool.Execute(validateParams)
					if err != nil {
						t.Errorf("Generated EAN-13 validation failed: %v", err)
					}
					validateMap := validateResult.(map[string]interface{})
					if !validateMap["valid"].(bool) {
						t.Errorf("Generated EAN-13 is not valid: %v", validateMap["error"])
					}
				} else {
					t.Errorf("Expected string result, got %T", result)
				}
			},
		},
		{
			name: "generate_multiple_ean13s",
			params: map[string]interface{}{
				"operation": "generate",
				"count":     float64(3),
			},
			wantErr: false,
			validate: func(t *testing.T, result interface{}) {
				if result == nil {
					t.Error("Expected non-nil result")
				}
				if ean13s, ok := result.([]string); ok {
					if len(ean13s) != 3 {
						t.Errorf("Expected 3 EAN-13s, got %d", len(ean13s))
					}
					// Validate each generated EAN-13
					for i, ean13 := range ean13s {
						if len(ean13) != 13 {
							t.Errorf("EAN-13 %d: expected 13 characters, got length %d", i, len(ean13))
						}
						validateParams := map[string]interface{}{
							"operation": "validate",
							"input":     ean13,
						}
						validateResult, err := tool.Execute(validateParams)
						if err != nil {
							t.Errorf("EAN-13 %d validation failed: %v", i, err)
						}
						validateMap := validateResult.(map[string]interface{})
						if !validateMap["valid"].(bool) {
							t.Errorf("EAN-13 %d is not valid: %v", i, validateMap["error"])
						}
					}
				} else {
					t.Errorf("Expected []string result, got %T", result)
				}
			},
		},
		{
			name: "invalid_operation",
			params: map[string]interface{}{
				"operation": "invalid",
			},
			wantErr: true,
		},
		{
			name: "validate_without_input",
			params: map[string]interface{}{
				"operation": "validate",
			},
			wantErr: true,
		},
		{
			name: "count_too_high",
			params: map[string]interface{}{
				"operation": "generate",
				"count":     float64(101),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tool.Execute(tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && tt.validate != nil {
				tt.validate(t, result)
			}
		})
	}
}

func TestEAN13ToolValidateParams(t *testing.T) {
	tool := NewEAN13Tool()

	tests := []struct {
		name    string
		params  map[string]interface{}
		wantErr bool
	}{
		{
			name:    "valid_empty_params",
			params:  map[string]interface{}{},
			wantErr: false,
		},
		{
			name: "valid_validate_operation",
			params: map[string]interface{}{
				"operation": "validate",
				"input":     "1234567890128",
			},
			wantErr: false,
		},
		{
			name: "valid_generate_operation",
			params: map[string]interface{}{
				"operation": "generate",
				"count":     float64(5),
			},
			wantErr: false,
		},
		{
			name: "valid_count",
			params: map[string]interface{}{
				"count": float64(10),
			},
			wantErr: false,
		},
		{
			name: "invalid_operation",
			params: map[string]interface{}{
				"operation": "invalid",
			},
			wantErr: true,
		},
		{
			name: "validate_without_input",
			params: map[string]interface{}{
				"operation": "validate",
			},
			wantErr: true,
		},
		{
			name: "count_too_low",
			params: map[string]interface{}{
				"count": float64(0),
			},
			wantErr: true,
		},
		{
			name: "count_too_high",
			params: map[string]interface{}{
				"count": float64(101),
			},
			wantErr: true,
		},
		{
			name: "invalid_parameter_types",
			params: map[string]interface{}{
				"operation": 123,
				"count":     "invalid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tool.ValidateParams(tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateParams() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEAN13ToolSchemas(t *testing.T) {
	tool := NewEAN13Tool()

	// Test input schema
	inputSchema := tool.GetInputSchema()
	if inputSchema == nil {
		t.Error("Expected non-nil input schema")
	}

	// Test output schema
	outputSchema := tool.GetOutputSchema()
	if outputSchema == nil {
		t.Error("Expected non-nil output schema")
	}

	// Verify schema structure
	if schemaType, ok := inputSchema["type"].(string); !ok || schemaType != "object" {
		t.Error("Input schema should have type 'object'")
	}

	if schemaType, ok := outputSchema["type"].(string); !ok || schemaType != "object" {
		t.Error("Output schema should have type 'object'")
	}
}

func TestEAN13ToolEdgeCases(t *testing.T) {
	tool := NewEAN13Tool()

	tests := []struct {
		name     string
		params   map[string]interface{}
		wantErr  bool
		validate func(t *testing.T, result interface{})
	}{
		{
			name: "single_ean13_with_count_1",
			params: map[string]interface{}{
				"operation": "generate",
				"count":     float64(1),
			},
			wantErr: false,
			validate: func(t *testing.T, result interface{}) {
				if result == nil {
					t.Error("Expected non-nil result")
				}
				if _, ok := result.(string); !ok {
					t.Errorf("Expected single string result, got %T", result)
				}
			},
		},
		{
			name: "maximum_count",
			params: map[string]interface{}{
				"operation": "generate",
				"count":     float64(100),
			},
			wantErr: false,
			validate: func(t *testing.T, result interface{}) {
				if result == nil {
					t.Error("Expected non-nil result")
				}
				if ean13s, ok := result.([]string); ok {
					if len(ean13s) != 100 {
						t.Errorf("Expected array of length 100, got: %d", len(ean13s))
					}
				} else {
					t.Errorf("Expected []string result, got %T", result)
				}
			},
		},
		{
			name: "validate_dashed_input",
			params: map[string]interface{}{
				"operation": "validate",
				"input":     "123-456-789-012-8",
			},
			wantErr: false,
			validate: func(t *testing.T, result interface{}) {
				if result == nil {
					t.Error("Expected non-nil result")
				}
				if resultMap, ok := result.(map[string]interface{}); ok {
					if valid, ok := resultMap["valid"].(bool); !ok || !valid {
						t.Errorf("Expected valid=true, got: %v", resultMap["valid"])
					}
				} else {
					t.Errorf("Expected map result, got %T", result)
				}
			},
		},
		{
			name: "validate_spaced_input",
			params: map[string]interface{}{
				"operation": "validate",
				"input":     "123 456 789 012 8",
			},
			wantErr: false,
			validate: func(t *testing.T, result interface{}) {
				if result == nil {
					t.Error("Expected non-nil result")
				}
				if resultMap, ok := result.(map[string]interface{}); ok {
					if valid, ok := resultMap["valid"].(bool); !ok || !valid {
						t.Errorf("Expected valid=true, got: %v", resultMap["valid"])
					}
				} else {
					t.Errorf("Expected map result, got %T", result)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tool.Execute(tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && tt.validate != nil {
				tt.validate(t, result)
			}
		})
	}
}

func TestEAN13ToolResources(t *testing.T) {
	tool := NewEAN13Tool()

	// Test GetResources
	resources := tool.GetResources()
	if len(resources) != 2 {
		t.Errorf("Expected 2 resources, got %d", len(resources))
	}

	// Test resource names and URIs
	expectedResources := map[string]string{
		"EAN-13 Algorithm": "ean13://algorithm",
		"EAN-13 Examples":  "ean13://examples",
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
	algorithmContent, err := tool.ReadResource("ean13://algorithm")
	if err != nil {
		t.Errorf("ReadResource(ean13://algorithm) failed: %v", err)
	}
	if algorithmContent == "" {
		t.Error("ReadResource(ean13://algorithm) returned empty content")
	}

	// Test ReadResource for examples
	examplesContent, err := tool.ReadResource("ean13://examples")
	if err != nil {
		t.Errorf("ReadResource(ean13://examples) failed: %v", err)
	}
	if examplesContent == "" {
		t.Error("ReadResource(ean13://examples) returned empty content")
	}

	// Test ReadResource with unknown URI
	_, err = tool.ReadResource("ean13://unknown")
	if err == nil {
		t.Error("ReadResource with unknown URI should return error")
	}
}
