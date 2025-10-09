package main

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestRunCreditCard(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected string
		hasError bool
	}{
		{
			name:     "validate_valid_visa",
			args:     []string{"--operation", "validate", "--input", "4532015112830366"},
			expected: "", // Will show validation result
			hasError: false,
		},
		{
			name:     "validate_valid_mastercard",
			args:     []string{"--operation", "validate", "--input", "5555555555554444"},
			expected: "", // Will show validation result
			hasError: false,
		},
		{
			name:     "validate_invalid_card",
			args:     []string{"--operation", "validate", "--input", "4532015112830367"},
			expected: "", // Will show error
			hasError: false,
		},
		{
			name:     "generate_single_visa",
			args:     []string{"--operation", "generate", "--card-type", "visa"},
			expected: "", // Will be a credit card number
			hasError: false,
		},
		{
			name:     "generate_multiple_cards",
			args:     []string{"--operation", "generate", "--card-type", "mastercard", "--count", "3"},
			expected: "", // Will be 3 credit card numbers
			hasError: false,
		},
		{
			name:     "generate_amex",
			args:     []string{"--operation", "generate", "--card-type", "amex"},
			expected: "", // Will be an Amex card number
			hasError: false,
		},
		{
			name:     "generate_without_type",
			args:     []string{"--operation", "generate", "--count", "2"},
			expected: "", // Will be 2 random type credit card numbers
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
			args:     []string{"--operation", "generate", "--card-type", "visa", "--count", "101"},
			expected: "",
			hasError: true,
		},
		{
			name:     "invalid_card_type",
			args:     []string{"--operation", "generate", "--card-type", "invalid", "--count", "1"},
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
			// Execute the CLI via go run
			args := append([]string{"run", ".", "creditcard"}, tt.args...)
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

func TestCreditCardCmdFlags(t *testing.T) {
	// Test that flags are properly defined
	if creditCardCmd.Flags().Lookup("operation") == nil {
		t.Error("--operation flag not found")
	}
	if creditCardCmd.Flags().Lookup("input") == nil {
		t.Error("--input flag not found")
	}
	if creditCardCmd.Flags().Lookup("card-type") == nil {
		t.Error("--card-type flag not found")
	}
	if creditCardCmd.Flags().Lookup("count") == nil {
		t.Error("--count flag not found")
	}
}

func TestCreditCardCmdHelp(t *testing.T) {
	// Test that help text is properly set
	if creditCardCmd.Short == "" {
		t.Error("Short description should not be empty")
	}
	if creditCardCmd.Long == "" {
		t.Error("Long description should not be empty")
	}
	if creditCardCmd.Use == "" {
		t.Error("Use string should not be empty")
	}
}

func TestCreditCardCmdGroup(t *testing.T) {
	// Test that command is in the tools group
	if creditCardCmd.GroupID != "tools" {
		t.Errorf("Expected group ID 'tools', got '%s'", creditCardCmd.GroupID)
	}
}

// TestRunCreditCardUnit tests the runCreditCard function directly with buffer (for coverage)
func TestRunCreditCardUnit(t *testing.T) {
	tests := []struct {
		name        string
		operation   string
		input       string
		cardType    string
		count       int
		expectError bool
	}{
		{
			name:        "validate valid visa",
			operation:   "validate",
			input:       "4532015112830366",
			expectError: false,
		},
		{
			name:        "validate invalid card",
			operation:   "validate",
			input:       "4532015112830367",
			expectError: false,
		},
		{
			name:        "generate single card",
			operation:   "generate",
			cardType:    "visa",
			count:       1,
			expectError: false,
		},
		{
			name:        "generate multiple cards",
			operation:   "generate",
			count:       3,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset and set global variables
			creditCardOperation = tt.operation
			creditCardInput = tt.input
			creditCardType = tt.cardType
			creditCardCount = tt.count
			if creditCardCount == 0 {
				creditCardCount = 1
			}

			// Create a buffer to capture output
			var buf bytes.Buffer

			// Call runCreditCard directly
			err := runCreditCard(nil, nil, &buf)

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
