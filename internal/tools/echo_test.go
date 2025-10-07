// Package tools provides the echo tool implementation
package tools

import (
	"testing"
)

func TestEchoTool(t *testing.T) {
	tool := NewEchoTool()

	// Test name and description
	if tool.Name() != "echo" {
		t.Errorf("Expected name 'echo', got '%s'", tool.Name())
	}

	if tool.Description() == "" {
		t.Error("Description should not be empty")
	}
}

func TestEchoToolValidateParams(t *testing.T) {
	tool := NewEchoTool()

	// Test valid parameters
	validParams := map[string]interface{}{
		"message": "test message",
	}

	if err := tool.ValidateParams(validParams); err != nil {
		t.Errorf("Valid parameters should not error: %v", err)
	}

	// Test missing message parameter
	invalidParams := map[string]interface{}{}

	if err := tool.ValidateParams(invalidParams); err == nil {
		t.Error("Missing message parameter should error")
	}
}

func TestEchoToolExecute(t *testing.T) {
	tool := NewEchoTool()

	// Test valid execution
	params := map[string]interface{}{
		"message": "test message",
	}

	result, err := tool.Execute(params)
	if err != nil {
		t.Errorf("Execute should not error: %v", err)
	}

	if result != "test message" {
		t.Errorf("Expected 'test message', got '%v'", result)
	}
}

func TestEchoToolExecuteEdgeCases(t *testing.T) {
	tool := NewEchoTool()

	tests := []struct {
		name        string
		params      map[string]interface{}
		expectError bool
		expectedMsg string
	}{
		{
			name:        "empty string message",
			params:      map[string]interface{}{"message": ""},
			expectError: false,
			expectedMsg: "",
		},
		{
			name:        "long message",
			params:      map[string]interface{}{"message": "This is a very long message that contains multiple words and should be echoed back exactly as provided without any modifications or truncation"},
			expectError: false,
			expectedMsg: "This is a very long message that contains multiple words and should be echoed back exactly as provided without any modifications or truncation",
		},
		{
			name:        "message with special characters",
			params:      map[string]interface{}{"message": "Hello, world! @#$%^&*()_+-=[]{}|;':\",./<>?"},
			expectError: false,
			expectedMsg: "Hello, world! @#$%^&*()_+-=[]{}|;':\",./<>?",
		},
		{
			name:        "message with unicode",
			params:      map[string]interface{}{"message": "Hello 世界! 🌍"},
			expectError: false,
			expectedMsg: "Hello 世界! 🌍",
		},
		{
			name:        "message with newlines",
			params:      map[string]interface{}{"message": "Line 1\nLine 2\nLine 3"},
			expectError: false,
			expectedMsg: "Line 1\nLine 2\nLine 3",
		},
		{
			name:        "non-string message type",
			params:      map[string]interface{}{"message": 123},
			expectError: true,
		},
		{
			name:        "nil message",
			params:      map[string]interface{}{"message": nil},
			expectError: true,
		},
		{
			name:        "missing message parameter",
			params:      map[string]interface{}{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tool.Execute(tt.params)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for %s, but got none", tt.name)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for %s: %v", tt.name, err)
				}
				if result != tt.expectedMsg {
					t.Errorf("Expected '%s', got '%v'", tt.expectedMsg, result)
				}
			}
		})
	}
}

func TestEchoToolValidateParamsEdgeCases(t *testing.T) {
	tool := NewEchoTool()

	tests := []struct {
		name        string
		params      map[string]interface{}
		expectError bool
	}{
		{
			name:        "valid string message",
			params:      map[string]interface{}{"message": "test"},
			expectError: false,
		},
		{
			name:        "empty string message",
			params:      map[string]interface{}{"message": ""},
			expectError: false,
		},
		{
			name:        "integer message",
			params:      map[string]interface{}{"message": 42},
			expectError: true,
		},
		{
			name:        "boolean message",
			params:      map[string]interface{}{"message": true},
			expectError: true,
		},
		{
			name:        "nil message",
			params:      map[string]interface{}{"message": nil},
			expectError: true,
		},
		{
			name:        "missing message",
			params:      map[string]interface{}{},
			expectError: true,
		},
		{
			name:        "extra parameters (should be ignored)",
			params:      map[string]interface{}{"message": "test", "extra": "ignored"},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tool.ValidateParams(tt.params)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for %s, but got none", tt.name)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for %s: %v", tt.name, err)
				}
			}
		})
	}
}

func TestEchoToolSchemas(t *testing.T) {
	tool := NewEchoTool()

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

	// Verify input schema structure
	if inputSchema["type"] != "object" {
		t.Errorf("Expected input schema type 'object', got '%v'", inputSchema["type"])
	}

	// Verify output schema structure
	if outputSchema["type"] != "object" {
		t.Errorf("Expected output schema type 'object', got '%v'", outputSchema["type"])
	}
}
