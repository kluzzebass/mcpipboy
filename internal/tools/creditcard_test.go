package tools

import (
	"strings"
	"testing"
)

func TestCreditCardToolExecute(t *testing.T) {
	tool := NewCreditCardTool()

	tests := []struct {
		name     string
		params   map[string]interface{}
		wantErr  bool
		validate func(t *testing.T, result interface{})
	}{
		{
			name: "validate_valid_visa",
			params: map[string]interface{}{
				"operation": "validate",
				"input":     "4532015112830366",
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
					if cardType, ok := resultMap["type"].(string); !ok || cardType != "visa" {
						t.Errorf("Expected type='visa', got: %v", resultMap["type"])
					}
				} else {
					t.Errorf("Expected map result, got %T", result)
				}
			},
		},
		{
			name: "validate_valid_mastercard",
			params: map[string]interface{}{
				"operation": "validate",
				"input":     "5555555555554444",
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
					if cardType, ok := resultMap["type"].(string); !ok || cardType != "mastercard" {
						t.Errorf("Expected type='mastercard', got: %v", resultMap["type"])
					}
				} else {
					t.Errorf("Expected map result, got %T", result)
				}
			},
		},
		{
			name: "validate_valid_amex",
			params: map[string]interface{}{
				"operation": "validate",
				"input":     "378282246310005",
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
					if cardType, ok := resultMap["type"].(string); !ok || cardType != "amex" {
						t.Errorf("Expected type='amex', got: %v", resultMap["type"])
					}
				} else {
					t.Errorf("Expected map result, got %T", result)
				}
			},
		},
		{
			name: "validate_invalid_card",
			params: map[string]interface{}{
				"operation": "validate",
				"input":     "4532015112830367",
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
					if errorMsg, ok := resultMap["error"].(string); !ok || errorMsg != "credit card number must be between 13 and 19 digits" {
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
				"input":     "4532-0151-1283-0366",
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
			name: "generate_single_visa",
			params: map[string]interface{}{
				"operation": "generate",
				"card-type": "visa",
			},
			wantErr: false,
			validate: func(t *testing.T, result interface{}) {
				if result == nil {
					t.Error("Expected non-nil result")
				}
				if card, ok := result.(string); ok {
					if len(card) != 16 {
						t.Errorf("Expected 16-digit Visa card, got length %d", len(card))
					}
					if !strings.HasPrefix(card, "4") {
						t.Errorf("Expected Visa card to start with 4, got: %s", card)
					}
					// Validate the generated card
					validateParams := map[string]interface{}{
						"operation": "validate",
						"input":     card,
					}
					validateResult, err := tool.Execute(validateParams)
					if err != nil {
						t.Errorf("Generated card validation failed: %v", err)
					}
					validateMap := validateResult.(map[string]interface{})
					if !validateMap["valid"].(bool) {
						t.Errorf("Generated card is not valid: %v", validateMap["error"])
					}
				} else {
					t.Errorf("Expected string result, got %T", result)
				}
			},
		},
		{
			name: "generate_multiple_cards",
			params: map[string]interface{}{
				"operation": "generate",
				"count":     float64(3),
			},
			wantErr: false,
			validate: func(t *testing.T, result interface{}) {
				if result == nil {
					t.Error("Expected non-nil result")
				}
				if cards, ok := result.([]string); ok {
					if len(cards) != 3 {
						t.Errorf("Expected 3 cards, got %d", len(cards))
					}
					// Validate each generated card
					for i, card := range cards {
						if len(card) < 13 || len(card) > 19 {
							t.Errorf("Card %d: expected 13-19 digits, got length %d", i, len(card))
						}
						validateParams := map[string]interface{}{
							"operation": "validate",
							"input":     card,
						}
						validateResult, err := tool.Execute(validateParams)
						if err != nil {
							t.Errorf("Card %d validation failed: %v", i, err)
						}
						validateMap := validateResult.(map[string]interface{})
						if !validateMap["valid"].(bool) {
							t.Errorf("Card %d is not valid: %v", i, validateMap["error"])
						}
					}
				} else {
					t.Errorf("Expected []string result, got %T", result)
				}
			},
		},
		{
			name: "generate_amex_card",
			params: map[string]interface{}{
				"operation": "generate",
				"card-type": "amex",
			},
			wantErr: false,
			validate: func(t *testing.T, result interface{}) {
				if result == nil {
					t.Error("Expected non-nil result")
				}
				if card, ok := result.(string); ok {
					if len(card) != 15 {
						t.Errorf("Expected 15-digit Amex card, got length %d", len(card))
					}
					if !strings.HasPrefix(card, "34") && !strings.HasPrefix(card, "37") {
						t.Errorf("Expected Amex card to start with 34 or 37, got: %s", card)
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
			name: "invalid_card_type",
			params: map[string]interface{}{
				"operation": "generate",
				"card-type": "invalid",
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

func TestCreditCardToolValidateParams(t *testing.T) {
	tool := NewCreditCardTool()

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
				"input":     "4532015112830366",
			},
			wantErr: false,
		},
		{
			name: "valid_generate_operation",
			params: map[string]interface{}{
				"operation": "generate",
				"card-type": "visa",
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
			name: "valid_card_type",
			params: map[string]interface{}{
				"card-type": "mastercard",
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
			name: "invalid_card_type",
			params: map[string]interface{}{
				"card-type": "invalid",
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

func TestCreditCardToolSchemas(t *testing.T) {
	tool := NewCreditCardTool()

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

func TestCreditCardToolEdgeCases(t *testing.T) {
	tool := NewCreditCardTool()

	tests := []struct {
		name     string
		params   map[string]interface{}
		wantErr  bool
		validate func(t *testing.T, result interface{})
	}{
		{
			name: "single_card_with_count_1",
			params: map[string]interface{}{
				"operation": "generate",
				"card-type": "visa",
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
				if cards, ok := result.([]string); ok {
					if len(cards) != 100 {
						t.Errorf("Expected array of length 100, got: %d", len(cards))
					}
				} else {
					t.Errorf("Expected []string result, got %T", result)
				}
			},
		},
		{
			name: "validate_cleaned_input",
			params: map[string]interface{}{
				"operation": "validate",
				"input":     "4532 0151 1283 0366",
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
			name: "validate_dashed_input",
			params: map[string]interface{}{
				"operation": "validate",
				"input":     "4532-0151-1283-0366",
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

func TestCreditCardToolResources(t *testing.T) {
	tool := NewCreditCardTool()

	// Test GetResources
	resources := tool.GetResources()
	if len(resources) != 3 {
		t.Errorf("Expected 3 resources, got %d", len(resources))
	}

	// Test resource names and URIs
	expectedResources := map[string]string{
		"Credit Card Types":    "creditcard://types",
		"Luhn Algorithm":       "creditcard://luhn",
		"Credit Card Examples": "creditcard://examples",
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

	// Test ReadResource for types
	typesContent, err := tool.ReadResource("creditcard://types")
	if err != nil {
		t.Errorf("ReadResource(creditcard://types) failed: %v", err)
	}
	if typesContent == "" {
		t.Error("ReadResource(creditcard://types) returned empty content")
	}

	// Test ReadResource for luhn algorithm
	luhnContent, err := tool.ReadResource("creditcard://luhn")
	if err != nil {
		t.Errorf("ReadResource(creditcard://luhn) failed: %v", err)
	}
	if luhnContent == "" {
		t.Error("ReadResource(creditcard://luhn) returned empty content")
	}

	// Test ReadResource for examples
	examplesContent, err := tool.ReadResource("creditcard://examples")
	if err != nil {
		t.Errorf("ReadResource(creditcard://examples) failed: %v", err)
	}
	if examplesContent == "" {
		t.Error("ReadResource(creditcard://examples) returned empty content")
	}

	// Test ReadResource with unknown URI
	_, err = tool.ReadResource("creditcard://unknown")
	if err == nil {
		t.Error("ReadResource with unknown URI should return error")
	}
}
