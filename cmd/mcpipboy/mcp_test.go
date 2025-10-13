// Package main provides tests for the mcp command
package main

import (
	"testing"
)

func TestMCPCommand(t *testing.T) {
	// Test that mcp command is properly configured
	if mcpCmd.Use != "mcp" {
		t.Errorf("Expected mcp command use to be 'mcp', got '%s'", mcpCmd.Use)
	}

	if mcpCmd.Short == "" {
		t.Error("MCP command should have a short description")
	}
}

func TestMCPCommandFlags(t *testing.T) {
	// Test that flags are properly configured
	enableFlag := mcpCmd.Flag("enable")
	if enableFlag == nil {
		t.Error("MCP command should have --enable flag")
	}

	disableFlag := mcpCmd.Flag("disable")
	if disableFlag == nil {
		t.Error("MCP command should have --disable flag")
	}
}

func TestRunMCP(t *testing.T) {
	// Test that runMCP function exists and can be called
	// Note: We don't actually call runMCP here because it starts a real MCP server
	// that would block the test. The integration tests in internal/server/ cover
	// the actual server functionality.

	// Test that the mcp command is properly configured
	if mcpCmd.RunE == nil {
		t.Error("mcp command should have a RunE function")
	}

	// Test that the mcp command has the expected flags
	if mcpCmd.Flag("enable") == nil {
		t.Error("mcp command should have --enable flag")
	}

	if mcpCmd.Flag("disable") == nil {
		t.Error("mcp command should have --disable flag")
	}
}

