package main

import (
	"bytes"
	"os/exec"
	"strings"
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
			// Execute the CLI via go run
			args := append([]string{"run", ".", "random"}, tt.args...)
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

// TestRunRandomUnit tests the runRandom function directly with buffer (for coverage)
func TestRunRandomUnit(t *testing.T) {
	tests := []struct {
		name        string
		randType    string
		count       int
		min         float64
		max         float64
		precision   int
		expectError bool
	}{
		{
			name:        "generate integer",
			randType:    "integer",
			count:       1,
			min:         1,
			max:         100,
			expectError: false,
		},
		{
			name:        "generate float",
			randType:    "float",
			count:       1,
			min:         0,
			max:         1,
			precision:   3,
			expectError: false,
		},
		{
			name:        "generate boolean",
			randType:    "boolean",
			count:       1,
			expectError: false,
		},
		{
			name:        "generate multiple integers",
			randType:    "integer",
			count:       5,
			min:         1,
			max:         100,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset and set global variables
			randomType = tt.randType
			randomCount = tt.count
			randomMin = tt.min
			randomMax = tt.max
			randomPrecision = tt.precision
			if randomCount == 0 {
				randomCount = 1
			}

			// Create a buffer to capture output
			var buf bytes.Buffer

			// Call runRandom directly
			err := runRandom(nil, nil, &buf)

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
