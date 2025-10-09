package main

import (
	"testing"
)

func TestRunMMSI(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected string
		hasError bool
	}{
		{
			name:     "validate_valid_mmsi",
			args:     []string{"--operation", "validate", "--input", "366123456"},
			expected: "", // Will show validation result
			hasError: false,
		},
		{
			name:     "validate_invalid_mmsi",
			args:     []string{"--operation", "validate", "--input", "12345678"},
			expected: "", // Will show error
			hasError: false,
		},
		{
			name:     "generate_single",
			args:     []string{"--operation", "generate", "--country-code", "US"},
			expected: "", // Will be a MMSI number
			hasError: false,
		},
		{
			name:     "generate_multiple",
			args:     []string{"--operation", "generate", "--country-code", "GB", "--count", "3"},
			expected: "", // Will be 3 MMSI numbers
			hasError: false,
		},
		{
			name:     "generate_without_country",
			args:     []string{"--operation", "generate", "--count", "2"},
			expected: "", // Will be 2 MMSI numbers with default US country
			hasError: false,
		},
		{
			name:     "invalid_operation",
			args:     []string{"--operation", "invalid"},
			expected: "",
			hasError: true,
		},
		{
			name:     "count_too_high",
			args:     []string{"--operation", "generate", "--country-code", "US", "--count", "101"},
			expected: "",
			hasError: true,
		},
		{
			name:     "invalid_country_code",
			args:     []string{"--operation", "generate", "--country-code", "XX", "--count", "1"},
			expected: "",
			hasError: true,
		},
		{
			name:     "validate_without_input",
			args:     []string{"--operation", "validate"},
			expected: "",
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flags
			mmsiOperation = "validate"
			mmsiInput = ""
			mmsiCountryCode = "US"
			mmsiCount = 1

			// Parse flags
			mmsiCmd.ParseFlags(tt.args)

			// Test the runMMSI function
			err := runMMSI(mmsiCmd, []string{})
			
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

func TestMMSICmdFlags(t *testing.T) {
	// Test that flags are properly defined
	if mmsiCmd.Flags().Lookup("operation") == nil {
		t.Error("--operation flag not found")
	}
	if mmsiCmd.Flags().Lookup("input") == nil {
		t.Error("--input flag not found")
	}
	if mmsiCmd.Flags().Lookup("country-code") == nil {
		t.Error("--country-code flag not found")
	}
	if mmsiCmd.Flags().Lookup("count") == nil {
		t.Error("--count flag not found")
	}
}

func TestMMSICmdHelp(t *testing.T) {
	// Test that help text is properly set
	if mmsiCmd.Short == "" {
		t.Error("Short description should not be empty")
	}
	if mmsiCmd.Long == "" {
		t.Error("Long description should not be empty")
	}
	if mmsiCmd.Use == "" {
		t.Error("Use string should not be empty")
	}
}

func TestMMSICmdGroup(t *testing.T) {
	// Test that command is in the tools group
	if mmsiCmd.GroupID != "tools" {
		t.Errorf("Expected group ID 'tools', got '%s'", mmsiCmd.GroupID)
	}
}
