package main

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestRunISBN(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "validate valid ISBN-10",
			args:    []string{"--operation", "validate", "--input", "0-123456-78-9", "--format", "isbn10"},
			wantErr: false,
		},
		{
			name:    "validate valid ISBN-13",
			args:    []string{"--operation", "validate", "--input", "978-0-123456-78-9", "--format", "isbn13"},
			wantErr: false,
		},
		{
			name:    "validate with auto format",
			args:    []string{"--operation", "validate", "--input", "9780123456789"},
			wantErr: false,
		},
		{
			name:    "validate invalid ISBN",
			args:    []string{"--operation", "validate", "--input", "0-123456-78-8", "--format", "isbn10"},
			wantErr: false,
		},
		{
			name:    "generate single ISBN-13",
			args:    []string{"--operation", "generate"},
			wantErr: false,
		},
		{
			name:    "generate single ISBN-10",
			args:    []string{"--operation", "generate", "--format", "isbn10"},
			wantErr: false,
		},
		{
			name:    "generate multiple ISBNs",
			args:    []string{"--operation", "generate", "--format", "isbn13", "--count", "3"},
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
			name:    "invalid format",
			args:    []string{"--operation", "generate", "--format", "invalid"},
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
			args := append([]string{"run", ".", "isbn"}, tt.args...)
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

func TestISBNCmdFlags(t *testing.T) {
	// Test that all expected flags exist
	expectedFlags := []string{"operation", "input", "format", "count"}

	for _, flagName := range expectedFlags {
		flag := isbnCmd.Flag(flagName)
		if flag == nil {
			t.Errorf("Expected flag '%s' not found", flagName)
		}
	}
}

func TestISBNCmdHelp(t *testing.T) {
	// Test that the command has help text
	if isbnCmd.Short == "" {
		t.Error("ISBN command should have a short description")
	}

	if isbnCmd.Long == "" {
		t.Error("ISBN command should have a long description")
	}
}

func TestISBNCmdGroup(t *testing.T) {
	// Test that the command is assigned to the tools group
	if isbnCmd.GroupID != "tools" {
		t.Errorf("Expected GroupID 'tools', got '%s'", isbnCmd.GroupID)
	}
}

// TestRunISBNUnit tests the runISBN function directly with buffer (for coverage)
func TestRunISBNUnit(t *testing.T) {
	tests := []struct {
		name        string
		operation   string
		input       string
		format      string
		count       int
		expectError bool
	}{
		{
			name:        "validate valid ISBN-10",
			operation:   "validate",
			input:       "0123456789",
			expectError: false,
		},
		{
			name:        "validate invalid ISBN",
			operation:   "validate",
			input:       "1234567890",
			expectError: false,
		},
		{
			name:        "generate single ISBN-13",
			operation:   "generate",
			format:      "isbn13",
			count:       1,
			expectError: false,
		},
		{
			name:        "generate multiple ISBNs",
			operation:   "generate",
			count:       3,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset and set global variables
			isbnOperation = tt.operation
			isbnInput = tt.input
			isbnFormat = tt.format
			isbnCount = tt.count
			if isbnCount == 0 {
				isbnCount = 1
			}

			// Create a buffer to capture output
			var buf bytes.Buffer

			// Call runISBN directly
			err := runISBN(nil, nil, &buf)

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
