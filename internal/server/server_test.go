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
