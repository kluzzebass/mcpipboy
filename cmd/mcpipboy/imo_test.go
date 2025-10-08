package main

import (
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
			expected: "✅ Valid IMO: 1234567",
			hasError: false,
		},
		{
			name:     "validate_invalid_imo",
			args:     []string{"--operation", "validate", "--input", "1234568"},
			expected: "❌ Invalid IMO:",
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
			// Reset flags
			imoOperation = "validate"
			imoInput = ""
			imoCount = 1

			// Parse flags
			imoCmd.ParseFlags(tt.args)

			// Test the runIMO function
			err := runIMO(imoCmd, []string{})

			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
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
