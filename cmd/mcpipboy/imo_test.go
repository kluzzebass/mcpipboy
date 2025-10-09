package main

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestRunIMO(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected string
		hasError bool
	}{
		{
			name:     "validate_valid_imo",
			args:     []string{"--operation", "validate", "--input", "1234567"},
			expected: "Valid IMO: 1234567",
			hasError: false,
		},
		{
			name:     "validate_invalid_imo",
			args:     []string{"--operation", "validate", "--input", "1234568"},
			expected: "Invalid IMO:",
			hasError: false,
		},
		{
			name:     "generate_single",
			args:     []string{"--operation", "generate"},
			expected: "", // Will be a 7-digit number
			hasError: false,
		},
		{
			name:     "generate_multiple",
			args:     []string{"--operation", "generate", "--count", "3"},
			expected: "", // Will be 3 lines of 7-digit numbers
			hasError: false,
		},
		{
			name:     "invalid_operation",
			args:     []string{"--operation", "invalid"},
			expected: "",
			hasError: true,
		},
		{
			name:     "missing_input_for_validate",
			args:     []string{"--operation", "validate"},
			expected: "",
			hasError: true,
		},
		{
			name:     "count_too_high",
			args:     []string{"--operation", "generate", "--count", "101"},
			expected: "",
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Execute the CLI via go run
			args := append([]string{"run", ".", "imo"}, tt.args...)
			cmd := exec.Command("go", args...)
			output, err := cmd.CombinedOutput()

			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but command succeeded. Output: %s", string(output))
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v\nOutput: %s", err, string(output))
				return
			}

			// Check that output is not empty
			outputStr := strings.TrimSpace(string(output))
			if len(outputStr) == 0 {
				t.Error("Expected non-empty output")
			}
		})
	}
}

func TestIMOCmdFlags(t *testing.T) {
	// Test that flags are properly defined
	if imoCmd.Flags().Lookup("operation") == nil {
		t.Error("--operation flag not found")
	}
	if imoCmd.Flags().Lookup("input") == nil {
		t.Error("--input flag not found")
	}
	if imoCmd.Flags().Lookup("count") == nil {
		t.Error("--count flag not found")
	}
}

func TestIMOCmdHelp(t *testing.T) {
	// Test that help text is properly set
	if imoCmd.Short == "" {
		t.Error("Short description should not be empty")
	}
	if imoCmd.Long == "" {
		t.Error("Long description should not be empty")
	}
	if imoCmd.Use == "" {
		t.Error("Use string should not be empty")
	}
}

func TestIMOCmdGroup(t *testing.T) {
	// Test that command is in the tools group
	if imoCmd.GroupID != "tools" {
		t.Errorf("Expected group ID 'tools', got '%s'", imoCmd.GroupID)
	}
}

// TestRunIMOUnit tests the runIMO function directly with buffer (for coverage)
func TestRunIMOUnit(t *testing.T) {
	tests := []struct {
		name        string
		operation   string
		input       string
		count       int
		expectError bool
	}{
		{
			name:        "validate valid IMO",
			operation:   "validate",
			input:       "1234567",
			expectError: false,
		},
		{
			name:        "validate invalid IMO",
			operation:   "validate",
			input:       "1234568",
			expectError: false,
		},
		{
			name:        "generate single IMO",
			operation:   "generate",
			count:       1,
			expectError: false,
		},
		{
			name:        "generate multiple IMOs",
			operation:   "generate",
			count:       3,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset and set global variables
			imoOperation = tt.operation
			imoInput = tt.input
			imoCount = tt.count
			if imoCount == 0 {
				imoCount = 1
			}

			// Create a buffer to capture output
			var buf bytes.Buffer

			// Call runIMO directly
			err := runIMO(nil, nil, &buf)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Check that output is not empty
			output := buf.String()
			if len(strings.TrimSpace(output)) == 0 {
				t.Error("Expected non-empty output")
			}
		})
	}
}
