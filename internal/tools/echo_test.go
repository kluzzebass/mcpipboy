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
