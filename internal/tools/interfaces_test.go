// Package tools provides tests for common interfaces and structures for mcpipboy tools
package tools

import (
	"fmt"
	"testing"
)

// MockTool implements the Tool interface for testing
type MockTool struct {
	name         string
	description  string
	inputSchema  map[string]interface{}
	outputSchema map[string]interface{}
	executeFunc  func(params map[string]interface{}) (interface{}, error)
	validateFunc func(params map[string]interface{}) error
}

func (m *MockTool) Name() string {
	return m.name
}

func (m *MockTool) Description() string {
	return m.description
}

func (m *MockTool) Execute(params map[string]interface{}) (interface{}, error) {
	if m.executeFunc != nil {
		return m.executeFunc(params)
	}
	return "mock result", nil
}

func (m *MockTool) ValidateParams(params map[string]interface{}) error {
	if m.validateFunc != nil {
		return m.validateFunc(params)
	}
	return nil
}

func (m *MockTool) GetInputSchema() map[string]interface{} {
	return m.inputSchema
}

func (m *MockTool) GetOutputSchema() map[string]interface{} {
	return m.outputSchema
}

func TestToolRegistry(t *testing.T) {
	registry := NewToolRegistry()

	// Test empty registry
	if len(registry.ListTools()) != 0 {
		t.Error("New registry should be empty")
	}

	// Test registering a tool
	mockTool := &MockTool{
		name:        "test-tool",
		description: "A test tool",
		inputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"message": map[string]interface{}{
					"type": "string",
				},
			},
		},
		outputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"result": map[string]interface{}{
					"type": "string",
				},
			},
		},
	}

	err := registry.RegisterTool(mockTool)
	if err != nil {
		t.Errorf("Failed to register tool: %v", err)
	}

	// Test listing tools
	tools := registry.ListTools()
	if len(tools) != 1 || tools[0] != "test-tool" {
		t.Errorf("Expected [test-tool], got %v", tools)
	}

	// Test getting tool
	tool, exists := registry.GetTool("test-tool")
	if !exists {
		t.Error("Tool should exist")
	}
	if tool.Name() != "test-tool" {
		t.Errorf("Expected test-tool, got %s", tool.Name())
	}

	// Test getting non-existent tool
	_, exists = registry.GetTool("non-existent")
	if exists {
		t.Error("Non-existent tool should not exist")
	}

	// Test duplicate registration
	err = registry.RegisterTool(mockTool)
	if err == nil {
		t.Error("Should not allow duplicate tool registration")
	}

	// Test registering nil tool
	err = registry.RegisterTool(nil)
	if err == nil {
		t.Error("Should not allow nil tool registration")
	}

	// Test registering tool with empty name
	emptyNameTool := &MockTool{name: ""}
	err = registry.RegisterTool(emptyNameTool)
	if err == nil {
		t.Error("Should not allow tool with empty name")
	}
}

func TestToolMetadata(t *testing.T) {
	registry := NewToolRegistry()

	mockTool := &MockTool{
		name:        "metadata-test",
		description: "Test metadata",
		inputSchema: map[string]interface{}{
			"type": "object",
		},
		outputSchema: map[string]interface{}{
			"type": "object",
		},
	}

	registry.RegisterTool(mockTool)

	// Test getting metadata
	metadata, err := registry.GetToolMetadata("metadata-test")
	if err != nil {
		t.Errorf("Failed to get metadata: %v", err)
	}

	if metadata.Name != "metadata-test" {
		t.Errorf("Expected metadata-test, got %s", metadata.Name)
	}

	if metadata.Description != "Test metadata" {
		t.Errorf("Expected 'Test metadata', got %s", metadata.Description)
	}

	// Test getting metadata for non-existent tool
	_, err = registry.GetToolMetadata("non-existent")
	if err == nil {
		t.Error("Should return error for non-existent tool")
	}

	// Test getting all metadata
	allMetadata, err := registry.GetAllToolMetadata()
	if err != nil {
		t.Errorf("Failed to get all metadata: %v", err)
	}

	if len(allMetadata) != 1 {
		t.Errorf("Expected 1 metadata entry, got %d", len(allMetadata))
	}
}

func TestToolExecution(t *testing.T) {
	registry := NewToolRegistry()

	executed := false
	mockTool := &MockTool{
		name:        "execution-test",
		description: "Test execution",
		executeFunc: func(params map[string]interface{}) (interface{}, error) {
			executed = true
			message, ok := params["message"].(string)
			if !ok {
				return nil, fmt.Errorf("message parameter required")
			}
			return "echo: " + message, nil
		},
		validateFunc: func(params map[string]interface{}) error {
			if _, ok := params["message"]; !ok {
				return fmt.Errorf("message parameter required")
			}
			return nil
		},
	}

	registry.RegisterTool(mockTool)

	// Test successful execution
	params := map[string]interface{}{
		"message": "hello",
	}

	result, err := registry.ExecuteTool("execution-test", params)
	if err != nil {
		t.Errorf("Execution failed: %v", err)
	}

	if !executed {
		t.Error("Tool should have been executed")
	}

	if result != "echo: hello" {
		t.Errorf("Expected 'echo: hello', got %v", result)
	}

	// Test execution with invalid parameters
	invalidParams := map[string]interface{}{
		"invalid": "param",
	}

	_, err = registry.ExecuteTool("execution-test", invalidParams)
	if err == nil {
		t.Error("Should fail with invalid parameters")
	}

	// Test execution of non-existent tool
	_, err = registry.ExecuteTool("non-existent", params)
	if err == nil {
		t.Error("Should fail for non-existent tool")
	}
}

func TestCreateJSONSchema(t *testing.T) {
	params := []ParameterDefinition{
		{
			Name:        "message",
			Type:        "string",
			Description: "The message to echo",
			Required:    true,
		},
		{
			Name:        "count",
			Type:        "number",
			Description: "Number of times to repeat",
			Required:    false,
			Default:     1,
		},
		{
			Name:        "format",
			Type:        "string",
			Description: "Output format",
			Required:    false,
			Enum:        []string{"text", "json"},
		},
	}

	schema := CreateJSONSchema(params)

	// Check schema structure
	if schema["type"] != "object" {
		t.Error("Schema type should be object")
	}

	properties, ok := schema["properties"].(map[string]interface{})
	if !ok {
		t.Error("Schema should have properties")
	}

	// Check message property
	messageProp, ok := properties["message"].(map[string]interface{})
	if !ok {
		t.Error("Message property should exist")
	}
	if messageProp["type"] != "string" {
		t.Error("Message property type should be string")
	}

	// Check required fields
	required, ok := schema["required"].([]string)
	if !ok {
		t.Error("Schema should have required fields")
	}
	if len(required) != 1 || required[0] != "message" {
		t.Errorf("Expected [message], got %v", required)
	}
}

func TestValidateParameter(t *testing.T) {
	// Test string parameter
	param := ParameterDefinition{
		Name:     "message",
		Type:     "string",
		Required: true,
	}

	err := ValidateParameter("hello", param)
	if err != nil {
		t.Errorf("Valid string should pass: %v", err)
	}

	err = ValidateParameter(123, param)
	if err == nil {
		t.Error("Invalid type should fail")
	}

	// Test required parameter
	err = ValidateParameter(nil, param)
	if err == nil {
		t.Error("Missing required parameter should fail")
	}

	// Test optional parameter
	param.Required = false
	err = ValidateParameter(nil, param)
	if err != nil {
		t.Errorf("Optional parameter should allow nil: %v", err)
	}

	// Test enum validation
	param.Enum = []string{"option1", "option2"}
	err = ValidateParameter("option1", param)
	if err != nil {
		t.Errorf("Valid enum value should pass: %v", err)
	}

	err = ValidateParameter("invalid", param)
	if err == nil {
		t.Error("Invalid enum value should fail")
	}
}

func TestValidateParameters(t *testing.T) {
	paramDefs := []ParameterDefinition{
		{
			Name:     "message",
			Type:     "string",
			Required: true,
		},
		{
			Name:     "count",
			Type:     "number",
			Required: false,
		},
	}

	// Test valid parameters
	params := map[string]interface{}{
		"message": "hello",
		"count":   5,
	}

	err := ValidateParameters(params, paramDefs)
	if err != nil {
		t.Errorf("Valid parameters should pass: %v", err)
	}

	// Test missing required parameter
	invalidParams := map[string]interface{}{
		"count": 5,
	}

	err = ValidateParameters(invalidParams, paramDefs)
	if err == nil {
		t.Error("Missing required parameter should fail")
	}

	// Test unknown parameter
	unknownParams := map[string]interface{}{
		"message": "hello",
		"unknown": "value",
	}

	err = ValidateParameters(unknownParams, paramDefs)
	if err == nil {
		t.Error("Unknown parameter should fail")
	}
}

func TestToolResult(t *testing.T) {
	// Test success result
	successResult := NewSuccessResult("test data")
	if !successResult.Success {
		t.Error("Success result should be successful")
	}
	if successResult.Data != "test data" {
		t.Errorf("Expected 'test data', got %v", successResult.Data)
	}

	// Test error result
	testErr := fmt.Errorf("test error")
	errorResult := NewErrorResult(testErr)
	if errorResult.Success {
		t.Error("Error result should not be successful")
	}
	if errorResult.Error != "test error" {
		t.Errorf("Expected 'test error', got %s", errorResult.Error)
	}

	// Test JSON conversion
	jsonData, err := successResult.ToJSON()
	if err != nil {
		t.Errorf("Failed to convert to JSON: %v", err)
	}

	expectedJSON := `{"success":true,"data":"test data"}`
	if string(jsonData) != expectedJSON {
		t.Errorf("Expected %s, got %s", expectedJSON, string(jsonData))
	}
}
