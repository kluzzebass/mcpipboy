// Package tools provides the version tool implementation
package tools

import (
	"testing"
)

func TestVersionTool(t *testing.T) {
	tool := NewVersionTool()

	// Test name and description
	if tool.Name() != "version" {
		t.Errorf("Expected name 'version', got '%s'", tool.Name())
	}

	if tool.Description() == "" {
		t.Error("Description should not be empty")
	}
}

func TestVersionToolValidateParams(t *testing.T) {
	tool := NewVersionTool()

	// Test that version tool accepts any parameters (or none)
	params := map[string]interface{}{}

	if err := tool.ValidateParams(params); err != nil {
		t.Errorf("Version tool should accept empty parameters: %v", err)
	}
}

func TestVersionToolExecute(t *testing.T) {
	tool := NewVersionTool()

	// Test execution
	params := map[string]interface{}{}

	result, err := tool.Execute(params)
	if err != nil {
		t.Errorf("Execute should not error: %v", err)
	}

	if result == nil {
		t.Error("Version result should not be nil")
	}
}
