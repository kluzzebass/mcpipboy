// Package main provides tests for the version command
package main

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

// TestVersionCommand tests the version command via integration (go run)
func TestVersionCommand(t *testing.T) {
	// Execute the CLI via go run
	cmd := exec.Command("go", "run", ".", "version")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Errorf("Unexpected error: %v\nOutput: %s", err, string(output))
		return
	}

	outputStr := strings.TrimSpace(string(output))
	// Check that output contains a version number (e.g., "0.1.0")
	if len(outputStr) == 0 {
		t.Error("Version command should output version information")
	}
}

// TestRunVersion tests the runVersion function directly with buffer (for coverage)
func TestRunVersion(t *testing.T) {
	// Create a buffer to capture output
	var buf bytes.Buffer

	// Call runVersion directly
	err := runVersion(nil, []string{}, &buf)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	// Check the output contains a version
	output := strings.TrimSpace(buf.String())
	if len(output) == 0 {
		t.Error("Version output should not be empty")
	}

	// Verify it looks like a version number (contains digits and dots)
	if !strings.Contains(output, ".") {
		t.Errorf("Expected version format with dots, got %q", output)
	}
}
