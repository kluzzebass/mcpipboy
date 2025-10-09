package main

import (
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
			// Parse flags manually
			cmd := isbnCmd
			cmd.ParseFlags(tt.args)

			// Execute the command directly
			err := runISBN(cmd, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("runISBN() error = %v, wantErr %v", err, tt.wantErr)
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
