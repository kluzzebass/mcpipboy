package main

import (
	"bytes"
	"os/exec"
	"strings"
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
			// Execute the CLI via go run
			args := append([]string{"run", ".", "iban"}, tt.args...)
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

			// Optionally check that output is not empty
			outputStr := strings.TrimSpace(string(output))
			if len(outputStr) == 0 {
				t.Error("Expected non-empty output")
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

// TestRunIBANUnit tests the runIBAN function directly with buffer (for coverage)
func TestRunIBANUnit(t *testing.T) {
	tests := []struct {
		name           string
		operation      string
		input          string
		countryCode    string
		count          int
		expectedOutput string
		expectError    bool
	}{
		{
			name:           "validate valid UK IBAN",
			operation:      "validate",
			input:          "GB82WEST12345698765432",
			expectedOutput: "Valid IBAN: GB82WEST12345698765432",
			expectError:    false,
		},
		{
			name:           "validate invalid IBAN",
			operation:      "validate",
			input:          "GB82WEST12345698765433",
			expectedOutput: "Invalid IBAN",
			expectError:    false,
		},
		{
			name:        "generate single IBAN",
			operation:   "generate",
			countryCode: "GB",
			count:       1,
			expectError: false,
		},
		{
			name:        "validate without input",
			operation:   "validate",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset global variables
			ibanOperation = tt.operation
			ibanInput = tt.input
			ibanCountryCode = tt.countryCode
			ibanCount = tt.count
			if ibanCount == 0 {
				ibanCount = 1
			}

			// Create a buffer to capture output
			var buf bytes.Buffer

			// Call runIBAN directly
			err := runIBAN(nil, nil, &buf)

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

			// Check the output if expected output is specified
			if tt.expectedOutput != "" {
				output := buf.String()
				if !strings.Contains(output, tt.expectedOutput) {
					t.Errorf("Expected output to contain %q, got %q", tt.expectedOutput, output)
				}
			}
		})
	}
}
