package tools

import (
	"strings"
	"testing"
)

func TestISBNToolExecute(t *testing.T) {
	tool := NewISBNTool()

	tests := []struct {
		name     string
		params   map[string]interface{}
		wantErr  bool
		validate func(t *testing.T, result interface{})
	}{
		{
			name: "validate_valid_isbn10",
			params: map[string]interface{}{
				"operation": "validate",
				"input":     "0-123456-78-9",
				"format":    "isbn10",
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
					if format, ok := resultMap["format"].(string); !ok || format != "ISBN10" {
						t.Errorf("Expected format='ISBN10', got: %v", resultMap["format"])
					}
				} else {
					t.Errorf("Expected map result, got %T", result)
				}
			},
		},
		{
			name: "validate_valid_isbn13",
			params: map[string]interface{}{
				"operation": "validate",
				"input":     "978-0-123456-78-6",
				"format":    "isbn13",
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
					if format, ok := resultMap["format"].(string); !ok || format != "ISBN13" {
						t.Errorf("Expected format='ISBN13', got: %v", resultMap["format"])
					}
				} else {
					t.Errorf("Expected map result, got %T", result)
				}
			},
		},
		{
			name: "validate_isbn10_with_x",
			params: map[string]interface{}{
				"operation": "validate",
				"input":     "100000001X",
				"format":    "isbn10",
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
			name: "validate_auto_detection_isbn10",
			params: map[string]interface{}{
				"operation": "validate",
				"input":     "0123456789",
				"format":    "auto",
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
			name: "validate_auto_detection_isbn13",
			params: map[string]interface{}{
				"operation": "validate",
				"input":     "9780123456786",
				"format":    "auto",
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
			name: "validate_invalid_isbn10",
			params: map[string]interface{}{
				"operation": "validate",
				"input":     "0-123456-78-8",
				"format":    "isbn10",
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
				"input":     "123456789",
				"format":    "isbn10",
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
					if errorMsg, ok := resultMap["error"].(string); !ok || !strings.Contains(errorMsg, "exactly 10 characters") {
						t.Errorf("Expected length error, got: %v", resultMap["error"])
					}
				} else {
					t.Errorf("Expected map result, got %T", result)
				}
			},
		},
		{
			name: "validate_cleaned_input",
			params: map[string]interface{}{
				"operation": "validate",
				"input":     "0 123 456 78 9",
				"format":    "isbn10",
			},
			wantErr: false,
			validate: func(t *testing.T, result interface{}) {
				if result == nil {
					t.Error("Expected non-nil result")
				}
				if resultMap, ok := result.(map[string]interface{}); ok {
					if valid, ok := resultMap["valid"].(bool); !ok || !valid {
						t.Errorf("Expected valid=true (after cleaning), got: %v", resultMap["valid"])
					}
				} else {
					t.Errorf("Expected map result, got %T", result)
				}
			},
		},
		{
			name: "generate_single_isbn10",
			params: map[string]interface{}{
				"operation": "generate",
				"format":    "isbn10",
			},
			wantErr: false,
			validate: func(t *testing.T, result interface{}) {
				if result == nil {
					t.Error("Expected non-nil result")
				}
				if isbn, ok := result.(string); ok {
					if len(isbn) != 10 {
						t.Errorf("Expected 10-character ISBN-10, got length %d", len(isbn))
					}
					// Validate the generated ISBN
					validateParams := map[string]interface{}{
						"operation": "validate",
						"input":     isbn,
						"format":    "isbn10",
					}
					validateResult, err := tool.Execute(validateParams)
					if err != nil {
						t.Errorf("Generated ISBN validation failed: %v", err)
					}
					validateMap := validateResult.(map[string]interface{})
					if !validateMap["valid"].(bool) {
						t.Errorf("Generated ISBN is not valid: %v", validateMap["error"])
					}
				} else {
					t.Errorf("Expected string result, got %T", result)
				}
			},
		},
		{
			name: "generate_single_isbn13",
			params: map[string]interface{}{
				"operation": "generate",
				"format":    "isbn13",
			},
			wantErr: false,
			validate: func(t *testing.T, result interface{}) {
				if result == nil {
					t.Error("Expected non-nil result")
				}
				if isbn, ok := result.(string); ok {
					if len(isbn) != 13 {
						t.Errorf("Expected 13-character ISBN-13, got length %d", len(isbn))
					}
					// Validate the generated ISBN
					validateParams := map[string]interface{}{
						"operation": "validate",
						"input":     isbn,
						"format":    "isbn13",
					}
					validateResult, err := tool.Execute(validateParams)
					if err != nil {
						t.Errorf("Generated ISBN validation failed: %v", err)
					}
					validateMap := validateResult.(map[string]interface{})
					if !validateMap["valid"].(bool) {
						t.Errorf("Generated ISBN is not valid: %v", validateMap["error"])
					}
				} else {
					t.Errorf("Expected string result, got %T", result)
				}
			},
		},
		{
			name: "generate_multiple_isbns",
			params: map[string]interface{}{
				"operation": "generate",
				"format":    "isbn13",
				"count":     float64(3),
			},
			wantErr: false,
			validate: func(t *testing.T, result interface{}) {
				if result == nil {
					t.Error("Expected non-nil result")
				}
				if isbns, ok := result.([]string); ok {
					if len(isbns) != 3 {
						t.Errorf("Expected 3 ISBNs, got %d", len(isbns))
					}
					// Validate each generated ISBN
					for i, isbn := range isbns {
						if len(isbn) != 13 {
							t.Errorf("ISBN %d: expected 13 characters, got length %d", i, len(isbn))
						}
						validateParams := map[string]interface{}{
							"operation": "validate",
							"input":     isbn,
							"format":    "isbn13",
						}
						validateResult, err := tool.Execute(validateParams)
						if err != nil {
							t.Errorf("ISBN %d validation failed: %v", i, err)
						}
						validateMap := validateResult.(map[string]interface{})
						if !validateMap["valid"].(bool) {
							t.Errorf("ISBN %d is not valid: %v", i, validateMap["error"])
						}
					}
				} else {
					t.Errorf("Expected []string result, got %T", result)
				}
			},
		},
		{
			name: "generate_default_format",
			params: map[string]interface{}{
				"operation": "generate",
			},
			wantErr: false,
			validate: func(t *testing.T, result interface{}) {
				if result == nil {
					t.Error("Expected non-nil result")
				}
				if isbn, ok := result.(string); ok {
					if len(isbn) != 13 {
						t.Errorf("Expected 13-character ISBN-13 (default), got length %d", len(isbn))
					}
				} else {
					t.Errorf("Expected string result, got %T", result)
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
		{
			name: "invalid_format",
			params: map[string]interface{}{
				"operation": "generate",
				"format":    "invalid",
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

func TestISBNToolValidateParams(t *testing.T) {
	tool := NewISBNTool()

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
				"input":     "0-123456-78-9",
				"format":    "isbn10",
			},
			wantErr: false,
		},
		{
			name: "valid_generate_operation",
			params: map[string]interface{}{
				"operation": "generate",
				"format":    "isbn13",
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
			name: "valid_format",
			params: map[string]interface{}{
				"format": "auto",
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
			name: "invalid_format",
			params: map[string]interface{}{
				"format": "invalid",
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

func TestISBNToolSchemas(t *testing.T) {
	tool := NewISBNTool()

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

func TestISBNToolEdgeCases(t *testing.T) {
	tool := NewISBNTool()

	tests := []struct {
		name     string
		params   map[string]interface{}
		wantErr  bool
		validate func(t *testing.T, result interface{})
	}{
		{
			name: "single_isbn_with_count_1",
			params: map[string]interface{}{
				"operation": "generate",
				"format":    "isbn10",
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
				"format":    "isbn13",
				"count":     float64(100),
			},
			wantErr: false,
			validate: func(t *testing.T, result interface{}) {
				if result == nil {
					t.Error("Expected non-nil result")
				}
				if isbns, ok := result.([]string); ok {
					if len(isbns) != 100 {
						t.Errorf("Expected array of length 100, got: %d", len(isbns))
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
				"input":     "0-123-456-78-9",
				"format":    "isbn10",
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
				"input":     "0 123 456 78 9",
				"format":    "isbn10",
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

func TestISBNToolResources(t *testing.T) {
	tool := NewISBNTool()

	// Test GetResources
	resources := tool.GetResources()
	if len(resources) != 3 {
		t.Errorf("Expected 3 resources, got %d", len(resources))
	}

	// Test resource names and URIs
	expectedResources := map[string]string{
		"ISBN Formats":    "isbn://formats",
		"ISBN Algorithms": "isbn://algorithms",
		"ISBN Examples":   "isbn://examples",
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

	// Test ReadResource for formats
	formatsContent, err := tool.ReadResource("isbn://formats")
	if err != nil {
		t.Errorf("ReadResource(isbn://formats) failed: %v", err)
	}
	if formatsContent == "" {
		t.Error("ReadResource(isbn://formats) returned empty content")
	}

	// Test ReadResource for algorithms
	algorithmsContent, err := tool.ReadResource("isbn://algorithms")
	if err != nil {
		t.Errorf("ReadResource(isbn://algorithms) failed: %v", err)
	}
	if algorithmsContent == "" {
		t.Error("ReadResource(isbn://algorithms) returned empty content")
	}

	// Test ReadResource for examples
	examplesContent, err := tool.ReadResource("isbn://examples")
	if err != nil {
		t.Errorf("ReadResource(isbn://examples) failed: %v", err)
	}
	if examplesContent == "" {
		t.Error("ReadResource(isbn://examples) returned empty content")
	}

	// Test ReadResource with unknown URI
	_, err = tool.ReadResource("isbn://unknown")
	if err == nil {
		t.Error("ReadResource with unknown URI should return error")
	}
}
