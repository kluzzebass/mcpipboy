// Package tools provides the tool registry system
package tools

import (
	"testing"
)

func TestRegistry(t *testing.T) {
	registry := NewRegistry()

	// Test initial state
	if len(registry.List()) != 0 {
		t.Error("New registry should be empty")
	}
}

func TestRegistryRegister(t *testing.T) {
	registry := NewRegistry()
	tool := NewEchoTool()

	// Test successful registration
	err := registry.Register(tool)
	if err != nil {
		t.Errorf("Register should not error: %v", err)
	}

	// Test that tool is now in registry
	if len(registry.List()) != 1 {
		t.Error("Registry should contain one tool")
	}

	if registry.List()[0] != "echo" {
		t.Error("Registry should contain echo tool")
	}
}

func TestRegistryGet(t *testing.T) {
	registry := NewRegistry()
	tool := NewEchoTool()

	// Register tool
	registry.Register(tool)

	// Test getting existing tool
	retrievedTool, exists := registry.Get("echo")
	if !exists {
		t.Error("Should find echo tool")
	}

	if retrievedTool.Name() != "echo" {
		t.Error("Retrieved tool should be echo")
	}

	// Test getting non-existent tool
	_, exists = registry.Get("nonexistent")
	if exists {
		t.Error("Should not find nonexistent tool")
	}
}

func TestRegistryGetMetadata(t *testing.T) {
	registry := NewRegistry()

	// Register multiple tools
	registry.Register(NewEchoTool())
	registry.Register(NewVersionTool())

	metadata := registry.GetMetadata()
	if len(metadata) != 2 {
		t.Errorf("Expected 2 metadata entries, got %d", len(metadata))
	}
}
