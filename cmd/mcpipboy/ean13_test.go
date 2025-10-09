package main

import (
	"testing"
)

func TestRunEAN13(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "validate valid EAN-13",
			args:    []string{"--operation", "validate", "--input", "1234567890128"},
			wantErr: false,
		},
		{
			name:    "validate valid EAN-13 with formatting",
			args:    []string{"--operation", "validate", "--input", "123-456-789-012-8"},
			wantErr: false,
		},
		{
			name:    "validate invalid EAN-13",
			args:    []string{"--operation", "validate", "--input", "1234567890123"},
			wantErr: false,
		},
		{
			name:    "generate single EAN-13",
			args:    []string{"--operation", "generate"},
			wantErr: false,
		},
		{
			name:    "generate multiple EAN-13s",
			args:    []string{"--operation", "generate", "--count", "3"},
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
			name:    "validate without input",
			args:    []string{"--operation", "validate"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse flags manually
			cmd := ean13Cmd
			cmd.ParseFlags(tt.args)

			// Execute the command directly
			err := runEAN13(cmd, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("runEAN13() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEAN13CmdFlags(t *testing.T) {
	// Test that all expected flags exist
	expectedFlags := []string{"operation", "input", "count"}

	for _, flagName := range expectedFlags {
		flag := ean13Cmd.Flag(flagName)
		if flag == nil {
			t.Errorf("Expected flag '%s' not found", flagName)
		}
	}
}

func TestEAN13CmdHelp(t *testing.T) {
	// Test that the command has help text
	if ean13Cmd.Short == "" {
		t.Error("EAN-13 command should have a short description")
	}

	if ean13Cmd.Long == "" {
		t.Error("EAN-13 command should have a long description")
	}
}

func TestEAN13CmdGroup(t *testing.T) {
	// Test that the command is assigned to the tools group
	if ean13Cmd.GroupID != "tools" {
		t.Errorf("Expected GroupID 'tools', got '%s'", ean13Cmd.GroupID)
	}
}
