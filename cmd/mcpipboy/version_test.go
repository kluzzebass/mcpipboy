// Package main provides tests for the version command
package main

import (
	"testing"
)

func TestVersionCommand(t *testing.T) {
	// Test that version command is properly configured
	if versionCmd.Use != "version" {
		t.Errorf("Expected version command use to be 'version', got '%s'", versionCmd.Use)
	}

	if versionCmd.Short == "" {
		t.Error("Version command should have a short description")
	}
}

func TestRunVersion(t *testing.T) {
	// Test valid version execution
	err := runVersion(versionCmd, []string{})
	if err != nil {
		t.Errorf("runVersion should not error: %v", err)
	}
}
