// Package server provides tests for the MCP server implementation
package server

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// MockTool implements the tools.Tool interface for testing
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

func TestNewServer(t *testing.T) {
	server := NewServer()
	if server == nil {
		t.Error("NewServer() returned nil")
		return
	}
	if server.tools == nil {
		t.Error("Server tools map is nil")
	}
}

func TestRegisterTool(t *testing.T) {
	server := NewServer()
	mockTool := &MockTool{
		name:        "test-tool",
		description: "A test tool",
	}

	server.RegisterTool(mockTool)

	if len(server.tools) != 1 {
		t.Errorf("Expected 1 tool, got %d", len(server.tools))
	}

	if server.tools["test-tool"] != mockTool {
		t.Error("Tool not registered correctly")
	}
}

func TestServerStop(t *testing.T) {
	server := NewServer()
	err := server.Stop()
	if err != nil {
		t.Errorf("Stop() returned error: %v", err)
	}
}

func TestServerStartWithTimeout(t *testing.T) {
	server := NewServer()

	// Add a mock tool
	mockTool := &MockTool{
		name:        "echo",
		description: "Echo tool",
		inputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"message": map[string]interface{}{
					"type":        "string",
					"description": "The message to echo back",
				},
			},
			"required": []string{"message"},
		},
		outputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"result": map[string]interface{}{
					"type":        "string",
					"description": "The echoed message",
				},
			},
		},
		executeFunc: func(params map[string]interface{}) (interface{}, error) {
			message, ok := params["message"].(string)
			if !ok {
				return nil, fmt.Errorf("message parameter required")
			}
			return message, nil
		},
	}
	server.RegisterTool(mockTool)

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// This should start the server but timeout quickly
	err := server.Start(ctx)
	if err != nil && err != context.DeadlineExceeded {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestServerErrorHandling(t *testing.T) {
	server := NewServer()

	// Test registering nil tool
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Registering nil tool should not panic: %v", r)
		}
	}()

	// This should handle nil gracefully
	server.RegisterTool(nil)

	// Server should still be functional
	if server.tools == nil {
		t.Error("Server tools map should not be nil")
	}
}

func TestServerMultipleTools(t *testing.T) {
	server := NewServer()

	// Register multiple tools
	tool1 := &MockTool{
		name:        "tool1",
		description: "First tool",
	}
	tool2 := &MockTool{
		name:        "tool2",
		description: "Second tool",
	}

	server.RegisterTool(tool1)
	server.RegisterTool(tool2)

	if len(server.tools) != 2 {
		t.Errorf("Expected 2 tools, got %d", len(server.tools))
	}

	if server.tools["tool1"] != tool1 {
		t.Error("Tool1 not registered correctly")
	}

	if server.tools["tool2"] != tool2 {
		t.Error("Tool2 not registered correctly")
	}
}

func TestServerDuplicateToolRegistration(t *testing.T) {
	server := NewServer()

	tool1 := &MockTool{
		name:        "duplicate",
		description: "First registration",
	}
	tool2 := &MockTool{
		name:        "duplicate",
		description: "Second registration",
	}

	server.RegisterTool(tool1)
	server.RegisterTool(tool2)

	// Second registration should overwrite the first
	if len(server.tools) != 1 {
		t.Errorf("Expected 1 tool after duplicate registration, got %d", len(server.tools))
	}

	if server.tools["duplicate"] != tool2 {
		t.Error("Second tool registration should overwrite the first")
	}
}

func TestServerStopWithoutStart(t *testing.T) {
	server := NewServer()

	// Stop should work even if server was never started
	err := server.Stop()
	if err != nil {
		t.Errorf("Stop() should work without Start(): %v", err)
	}
}
