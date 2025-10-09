// Package main provides tests for the serve command
package main

import (
	"testing"
)

func TestServeCommand(t *testing.T) {
	// Test that serve command is properly configured
	if serveCmd.Use != "serve" {
		t.Errorf("Expected serve command use to be 'serve', got '%s'", serveCmd.Use)
	}

	if serveCmd.Short == "" {
		t.Error("Serve command should have a short description")
	}
}

func TestServeCommandFlags(t *testing.T) {
	// Test that flags are properly configured
	enableFlag := serveCmd.Flag("enable")
	if enableFlag == nil {
		t.Error("Serve command should have --enable flag")
	}

	disableFlag := serveCmd.Flag("disable")
	if disableFlag == nil {
		t.Error("Serve command should have --disable flag")
	}
}

func TestRunServe(t *testing.T) {
	// Test that runServe function exists and can be called
	// Note: We don't actually call runServe here because it starts a real MCP server
	// that would block the test. The integration tests in internal/server/ cover
	// the actual server functionality.

	// Test that the serve command is properly configured
	if serveCmd.RunE == nil {
		t.Error("serve command should have a RunE function")
	}

	// Test that the serve command has the expected flags
	if serveCmd.Flag("enable") == nil {
		t.Error("serve command should have --enable flag")
	}

	if serveCmd.Flag("disable") == nil {
		t.Error("serve command should have --disable flag")
	}
}
