package main

import (
	"testing"
)

func TestRunUUID(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected string
		hasError bool
	}{
		{
			name:     "generate_v4_default",
			args:     []string{"--version", "v4"},
			expected: "", // Will be a UUID string
			hasError: false,
		},
		{
			name:     "generate_v4_multiple",
			args:     []string{"--version", "v4", "--count", "3"},
			expected: "", // Will be 3 lines of UUIDs
			hasError: false,
		},
		{
			name:     "generate_v1",
			args:     []string{"--version", "v1", "--count", "1"},
			expected: "", // Will be a UUID v1
			hasError: false,
		},
		{
			name:     "generate_v5_with_namespace_and_name",
			args:     []string{"--version", "v5", "--namespace", "6ba7b810-9dad-11d1-80b4-00c04fd430c8", "--name", "example"},
			expected: "", // Will be a UUID v5
			hasError: false,
		},
		{
			name:     "generate_v7",
			args:     []string{"--version", "v7", "--count", "1"},
			expected: "", // Will be a UUID v7
			hasError: false,
		},
		{
			name:     "validate_valid_uuid",
			args:     []string{"--version", "validate", "--input", "550e8400-e29b-41d4-a716-446655440000"},
			expected: "", // Will show validation result
			hasError: false,
		},
		{
			name:     "validate_invalid_uuid",
			args:     []string{"--version", "validate", "--input", "invalid-uuid"},
			expected: "", // Will show error
			hasError: false,
		},
		{
			name:     "invalid_version",
			args:     []string{"--version", "invalid"},
			expected: "",
			hasError: true,
		},
		{
			name:     "count_too_low",
			args:     []string{"--version", "v4", "--count", "0"},
			expected: "",
			hasError: true,
		},
		{
			name:     "count_too_high",
			args:     []string{"--version", "v4", "--count", "1001"},
			expected: "",
			hasError: true,
		},
		{
			name:     "v5_without_namespace",
			args:     []string{"--version", "v5", "--name", "example"},
			expected: "",
			hasError: true,
		},
		{
			name:     "v5_without_name",
			args:     []string{"--version", "v5", "--namespace", "6ba7b810-9dad-11d1-80b4-00c04fd430c8"},
			expected: "",
			hasError: true,
		},
		{
			name:     "validate_without_input",
			args:     []string{"--version", "validate"},
			expected: "",
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flags
			uuidVersion = "v4"
			uuidCount = 1
			uuidNamespace = ""
			uuidName = ""
			uuidInput = ""

			// Parse flags
			uuidCmd.ParseFlags(tt.args)

			// Test the runUUID function
			err := runUUID(uuidCmd, []string{})

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

func TestUUIDCmdFlags(t *testing.T) {
	// Test that flags are properly defined
	if uuidCmd.Flags().Lookup("version") == nil {
		t.Error("--version flag not found")
	}
	if uuidCmd.Flags().Lookup("count") == nil {
		t.Error("--count flag not found")
	}
	if uuidCmd.Flags().Lookup("namespace") == nil {
		t.Error("--namespace flag not found")
	}
	if uuidCmd.Flags().Lookup("name") == nil {
		t.Error("--name flag not found")
	}
	if uuidCmd.Flags().Lookup("input") == nil {
		t.Error("--input flag not found")
	}
}

func TestUUIDCmdHelp(t *testing.T) {
	// Test that help text is properly set
	if uuidCmd.Short == "" {
		t.Error("Short description should not be empty")
	}
	if uuidCmd.Long == "" {
		t.Error("Long description should not be empty")
	}
	if uuidCmd.Use == "" {
		t.Error("Use string should not be empty")
	}
}

func TestUUIDCmdGroup(t *testing.T) {
	// Test that command is in the tools group
	if uuidCmd.GroupID != "tools" {
		t.Errorf("Expected group ID 'tools', got '%s'", uuidCmd.GroupID)
	}
}
