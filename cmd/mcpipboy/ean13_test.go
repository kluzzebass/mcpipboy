package main

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestRunEAN13(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "validate valid EAN-13",
			args:    []string{"--operation", "validate", "--input", "1234567890128"},
			wantErr: false,
		},
		{
			name:    "validate valid EAN-13 with formatting",
			args:    []string{"--operation", "validate", "--input", "123-456-789-012-8"},
			wantErr: false,
		},
		{
			name:    "validate invalid EAN-13",
			args:    []string{"--operation", "validate", "--input", "1234567890123"},
			wantErr: false,
		},
		{
			name:    "generate single EAN-13",
			args:    []string{"--operation", "generate"},
			wantErr: false,
		},
		{
			name:    "generate multiple EAN-13s",
			args:    []string{"--operation", "generate", "--count", "3"},
			wantErr: false,
		},
		{
			name:    "generate with maximum count",
			args:    []string{"--operation", "generate", "--count", "100"},
			wantErr: false,
		},
		{
			name:    "invalid operation",
			args:    []string{"--operation", "invalid"},
			wantErr: true,
		},
		{
			name:    "count too high",
			args:    []string{"--operation", "generate", "--count", "101"},
			wantErr: true,
		},
		{
			name:    "validate without input",
			args:    []string{"--operation", "validate"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Execute the CLI via go run
			args := append([]string{"run", ".", "ean13"}, tt.args...)
			cmd := exec.Command("go", args...)
			output, err := cmd.CombinedOutput()

			if tt.wantErr {
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

func TestEAN13CmdFlags(t *testing.T) {
	// Test that all expected flags exist
	expectedFlags := []string{"operation", "input", "count"}

	for _, flagName := range expectedFlags {
		flag := ean13Cmd.Flag(flagName)
		if flag == nil {
			t.Errorf("Expected flag '%s' not found", flagName)
		}
	}
}

func TestEAN13CmdHelp(t *testing.T) {
	// Test that the command has help text
	if ean13Cmd.Short == "" {
		t.Error("EAN-13 command should have a short description")
	}

	if ean13Cmd.Long == "" {
		t.Error("EAN-13 command should have a long description")
	}
}

func TestEAN13CmdGroup(t *testing.T) {
	// Test that the command is assigned to the tools group
	if ean13Cmd.GroupID != "tools" {
		t.Errorf("Expected GroupID 'tools', got '%s'", ean13Cmd.GroupID)
	}
}

// TestRunEAN13Unit tests the runEAN13 function directly with buffer (for coverage)
func TestRunEAN13Unit(t *testing.T) {
	tests := []struct {
		name        string
		operation   string
		input       string
		count       int
		expectError bool
	}{
		{
			name:        "validate valid EAN-13",
			operation:   "validate",
			input:       "1234567890128",
			expectError: false,
		},
		{
			name:        "validate invalid EAN-13",
			operation:   "validate",
			input:       "1234567890123",
			expectError: false,
		},
		{
			name:        "generate single EAN-13",
			operation:   "generate",
			count:       1,
			expectError: false,
		},
		{
			name:        "generate multiple EAN-13s",
			operation:   "generate",
			count:       5,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset and set global variables
			ean13Operation = tt.operation
			ean13Input = tt.input
			ean13Count = tt.count
			if ean13Count == 0 {
				ean13Count = 1
			}

			// Create a buffer to capture output
			var buf bytes.Buffer

			// Call runEAN13 directly
			err := runEAN13(nil, nil, &buf)

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
