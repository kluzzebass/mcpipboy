package main

import (
	"testing"
)

func TestRunIBAN(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected string
		hasError bool
	}{
		{
			name:     "validate_valid_UK_IBAN",
			args:     []string{"--operation", "validate", "--input", "GB82WEST12345698765432"},
			expected: "", // Will show validation result
			hasError: false,
		},
		{
			name:     "validate_valid_German_IBAN",
			args:     []string{"--operation", "validate", "--input", "DE89370400440532013000"},
			expected: "", // Will show validation result
			hasError: false,
		},
		{
			name:     "validate_valid_French_IBAN",
			args:     []string{"--operation", "validate", "--input", "FR1420041010050500013M02606"},
			expected: "", // Will show validation result
			hasError: false,
		},
		{
			name:     "validate_invalid_IBAN",
			args:     []string{"--operation", "validate", "--input", "GB82WEST12345698765433"},
			expected: "", // Will show error
			hasError: false,
		},
		{
			name:     "validate_IBAN_with_spaces",
			args:     []string{"--operation", "validate", "--input", "GB82 WEST 1234 5698 7654 32"},
			expected: "", // Will show validation result
			hasError: false,
		},
		{
			name:     "validate_IBAN_with_lowercase",
			args:     []string{"--operation", "validate", "--input", "gb82west12345698765432"},
			expected: "", // Will show validation result
			hasError: false,
		},
		{
			name:     "validate_missing_input",
			args:     []string{"--operation", "validate"},
			expected: "",
			hasError: true,
		},
		{
			name:     "validate_invalid_operation",
			args:     []string{"--operation", "invalid", "--input", "GB82WEST12345698765432"},
			expected: "",
			hasError: true,
		},
		{
			name:     "generate_single_IBAN",
			args:     []string{"--operation", "generate"},
			expected: "", // Will be an IBAN
			hasError: false,
		},
		{
			name:     "generate_single_IBAN_with_country",
			args:     []string{"--operation", "generate", "--country-code", "GB"},
			expected: "", // Will be a UK IBAN
			hasError: false,
		},
		{
			name:     "generate_multiple_IBANs",
			args:     []string{"--operation", "generate", "--count", "3"},
			expected: "", // Will be 3 IBANs
			hasError: false,
		},
		{
			name:     "generate_with_valid_country_codes",
			args:     []string{"--operation", "generate", "--country-code", "DE"},
			expected: "", // Will be a German IBAN
			hasError: false,
		},
		{
			name:     "generate_with_invalid_country_code",
			args:     []string{"--operation", "generate", "--country-code", "XX"},
			expected: "",
			hasError: true,
		},
		{
			name:     "generate_with_invalid_count_too_low",
			args:     []string{"--operation", "generate", "--count", "0"},
			expected: "",
			hasError: true,
		},
		{
			name:     "generate_with_invalid_count_too_high",
			args:     []string{"--operation", "generate", "--count", "101"},
			expected: "",
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset global variables
			ibanOperation = ""
			ibanInput = ""
			ibanCountryCode = ""
			ibanCount = 1

			// Parse the test args to set global variables
			for i := 0; i < len(tt.args); i += 2 {
				if i+1 < len(tt.args) {
					switch tt.args[i] {
					case "--operation":
						ibanOperation = tt.args[i+1]
					case "--input":
						ibanInput = tt.args[i+1]
					case "--country-code":
						ibanCountryCode = tt.args[i+1]
					case "--count":
						// Parse count as int
						if tt.args[i+1] == "0" {
							ibanCount = 0
						} else if tt.args[i+1] == "3" {
							ibanCount = 3
						} else if tt.args[i+1] == "101" {
							ibanCount = 101
						}
					}
				}
			}

			// Execute the runIBAN function directly
			err := runIBAN(nil, nil)

			if tt.hasError {
				if err == nil {
					t.Error("Expected error, got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
		})
	}
}

func TestIBANCmd_Flags(t *testing.T) {
	// Test that all flags are properly defined
	cmd := ibanCmd

	// Check that required flags exist
	flags := []string{"operation", "input", "country-code", "count"}
	for _, flag := range flags {
		if cmd.Flags().Lookup(flag) == nil {
			t.Errorf("Expected flag '%s' to be defined", flag)
		}
	}

	// Test flag defaults
	if cmd.Flags().Lookup("operation").DefValue != "validate" {
		t.Errorf("Expected default operation to be 'validate', got '%s'", cmd.Flags().Lookup("operation").DefValue)
	}

	if cmd.Flags().Lookup("count").DefValue != "1" {
		t.Errorf("Expected default count to be '1', got '%s'", cmd.Flags().Lookup("count").DefValue)
	}
}

func TestIBANCmd_Help(t *testing.T) {
	cmd := ibanCmd

	// Test that help text is not empty
	if cmd.Short == "" {
		t.Error("Expected Short description to be non-empty")
	}

	if cmd.Long == "" {
		t.Error("Expected Long description to be non-empty")
	}

	// Test that examples are included
	if cmd.Long != "" && len(cmd.Long) < 100 {
		t.Error("Expected Long description to be substantial")
	}
}

func TestIBANCmd_GroupID(t *testing.T) {
	cmd := ibanCmd

	// Test that the command is in the tools group
	if cmd.GroupID != "tools" {
		t.Errorf("Expected GroupID to be 'tools', got '%s'", cmd.GroupID)
	}
}

func TestIBANCmd_Use(t *testing.T) {
	cmd := ibanCmd

	// Test that the command use is correct
	if cmd.Use != "iban" {
		t.Errorf("Expected Use to be 'iban', got '%s'", cmd.Use)
	}
}
