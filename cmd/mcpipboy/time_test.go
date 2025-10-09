package main

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestTimeCommand(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "default time command",
			args:    []string{},
			wantErr: false,
		},
		{
			name:    "time with default (now)",
			args:    []string{},
			wantErr: false,
		},
		{
			name:    "time with format date",
			args:    []string{"--format", "date"},
			wantErr: false,
		},
		{
			name:    "time with timestamp input",
			args:    []string{"--input", "2025-01-01T12:00:00Z", "--format", "unix"},
			wantErr: false,
		},
		{
			name:    "time with unix timestamp",
			args:    []string{"--input", "2025-01-01T12:00:00Z", "--format", "datetime"},
			wantErr: false,
		},
		{
			name:    "time with offset",
			args:    []string{"--input", "2025-01-01T00:00:00Z", "--offset", "+1h", "--format", "datetime"},
			wantErr: false,
		},
		{
			name:    "time with timezone",
			args:    []string{"--input", "2025-01-01T12:00:00Z", "--timezone", "America/New_York", "--format", "datetime"},
			wantErr: false,
		},
		{
			name:    "invalid offset",
			args:    []string{"--offset", "invalid"},
			wantErr: true,
		},
		{
			name:    "invalid timezone",
			args:    []string{"--timezone", "Invalid/Timezone"},
			wantErr: true,
		},
		{
			name:    "invalid format",
			args:    []string{"--format", "invalid"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Execute the CLI via go run
			args := append([]string{"run", ".", "time"}, tt.args...)
			cmd := exec.Command("go", args...)
			output, err := cmd.CombinedOutput()

			if tt.wantErr {
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

func TestTimeCommandHelp(t *testing.T) {
	// Test that help text is properly set
	if timeCmd.Short == "" {
		t.Error("Time command should have a short description")
	}
	
	if timeCmd.Long == "" {
		t.Error("Time command should have a long description")
	}
	
	if timeCmd.Use == "" {
		t.Error("Time command should have usage text")
	}
}

func TestTimeCommandValidation(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "valid combination",
			args:    []string{"time", "--type", "timestamp", "--input", "2025-01-01", "--format", "date", "--timezone", "UTC"},
			wantErr: false,
		},
		{
			name:    "lenient timestamp parsing",
			args:    []string{"time", "--type", "timestamp", "--input", "January 1, 2025 at 12:00 PM", "--format", "date"},
			wantErr: false,
		},
		{
			name:    "complex offset",
			args:    []string{"time", "--type", "timestamp", "--input", "2025-01-01T00:00:00Z", "--offset", "+26h30m", "--format", "datetime"},
			wantErr: false,
		},
		{
			name:    "negative offset",
			args:    []string{"time", "--type", "timestamp", "--input", "2025-01-01T12:00:00Z", "--offset", "-2h", "--format", "datetime"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Execute the CLI via go run
			args := append([]string{"run", "."}, tt.args...)
			cmd := exec.Command("go", args...)
			output, err := cmd.CombinedOutput()

			if tt.wantErr {
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

// TestRunTimeUnit tests the runTime function directly with buffer (for coverage)
func TestRunTimeUnit(t *testing.T) {
	tests := []struct {
		name        string
		timeType    string
		format      string
		timezone    string
		input       string
		offset      string
		expectError bool
	}{
		{
			name:        "get current time",
			timeType:    "now",
			format:      "iso",
			timezone:    "utc",
			expectError: false,
		},
		{
			name:        "parse timestamp",
			timeType:    "timestamp",
			input:       "2025-01-01T12:00:00Z",
			format:      "date",
			expectError: false,
		},
		{
			name:        "time with offset",
			timeType:    "timestamp",
			input:       "2025-01-01T12:00:00Z",
			offset:      "+1h",
			format:      "datetime",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset and set global variables
			timeType = tt.timeType
			timeFormat = tt.format
			timezone = tt.timezone
			timeInput = tt.input
			timeOffset = tt.offset
			timeFrom = ""
			timeTo = ""

			// Create a buffer to capture output
			var buf bytes.Buffer

			// Call runTime directly
			err := runTime(nil, nil, &buf)

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
