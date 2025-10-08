package main

import (
	"testing"
)

func TestRunRandom(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected string
		hasError bool
	}{
		{
			name:     "generate_integer_default",
			args:     []string{"--type", "integer"},
			expected: "", // Will be a random integer
			hasError: false,
		},
		{
			name:     "generate_integer_with_range",
			args:     []string{"--type", "integer", "--min", "10", "--max", "20"},
			expected: "", // Will be a random integer between 10-20
			hasError: false,
		},
		{
			name:     "generate_multiple_integers",
			args:     []string{"--type", "integer", "--count", "5"},
			expected: "", // Will be 5 random integers
			hasError: false,
		},
		{
			name:     "generate_float_default",
			args:     []string{"--type", "float"},
			expected: "", // Will be a random float
			hasError: false,
		},
		{
			name:     "generate_float_with_precision",
			args:     []string{"--type", "float", "--precision", "3"},
			expected: "", // Will be a random float with 3 decimal places
			hasError: false,
		},
		{
			name:     "generate_boolean",
			args:     []string{"--type", "boolean"},
			expected: "", // Will be true or false
			hasError: false,
		},
		{
			name:     "generate_multiple_booleans",
			args:     []string{"--type", "boolean", "--count", "3"},
			expected: "", // Will be 3 random booleans
			hasError: false,
		},
		{
			name:     "invalid_type",
			args:     []string{"--type", "invalid"},
			expected: "",
			hasError: true,
		},
		{
			name:     "count_too_low",
			args:     []string{"--type", "integer", "--count", "0"},
			expected: "",
			hasError: true,
		},
		{
			name:     "count_too_high",
			args:     []string{"--type", "integer", "--count", "1001"},
			expected: "",
			hasError: true,
		},
		{
			name:     "invalid_range_min_greater_than_max",
			args:     []string{"--type", "integer", "--min", "20", "--max", "10"},
			expected: "",
			hasError: true,
		},
		{
			name:     "precision_too_high",
			args:     []string{"--type", "float", "--precision", "11"},
			expected: "",
			hasError: true,
		},
		{
			name:     "precision_negative",
			args:     []string{"--type", "float", "--precision", "-1"},
			expected: "",
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flags
			randomType = "integer"
			randomCount = 1
			randomMin = 0.0
			randomMax = 100.0
			randomPrecision = 2.0

			// Parse flags
			randomCmd.ParseFlags(tt.args)

			// Test the runRandom function
			err := runRandom(randomCmd, []string{})

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

func TestRandomCmdFlags(t *testing.T) {
	// Test that flags are properly defined
	if randomCmd.Flags().Lookup("type") == nil {
		t.Error("--type flag not found")
	}
	if randomCmd.Flags().Lookup("count") == nil {
		t.Error("--count flag not found")
	}
	if randomCmd.Flags().Lookup("min") == nil {
		t.Error("--min flag not found")
	}
	if randomCmd.Flags().Lookup("max") == nil {
		t.Error("--max flag not found")
	}
	if randomCmd.Flags().Lookup("precision") == nil {
		t.Error("--precision flag not found")
	}
}

func TestRandomCmdHelp(t *testing.T) {
	// Test that help text is properly set
	if randomCmd.Short == "" {
		t.Error("Short description should not be empty")
	}
	if randomCmd.Long == "" {
		t.Error("Long description should not be empty")
	}
	if randomCmd.Use == "" {
		t.Error("Use string should not be empty")
	}
}

func TestRandomCmdGroup(t *testing.T) {
	// Test that command is in the tools group
	if randomCmd.GroupID != "tools" {
		t.Errorf("Expected group ID 'tools', got '%s'", randomCmd.GroupID)
	}
}
