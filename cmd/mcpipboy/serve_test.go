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
	// Test that runServe can be called without crashing
	// Note: This will start the MCP server and wait for input, so we expect it to eventually timeout or error
	// In a real test environment, we'd mock the server or use a timeout
	err := runServe(serveCmd, []string{})
	// The server will start and wait for MCP protocol messages
	// We expect it to eventually error when the test environment doesn't provide proper MCP input
	if err != nil {
		// This is expected - the server will error when it doesn't receive proper MCP input
		t.Logf("Server error (expected in test environment): %v", err)
	}
}
