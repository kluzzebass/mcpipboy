package main

import (
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
			// Reset flags for each test
			timeType = "now"
			timeFormat = "iso"
			timezone = "utc"
			timeInput = ""
			timeFrom = ""
			timeTo = ""
			timeOffset = ""

			// Parse flags manually
			cmd := timeCmd
			cmd.ParseFlags(tt.args)

			// Execute the command directly
			err := runTime(cmd, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("time command error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTimeCommandHelp(t *testing.T) {
	// Test that help is properly formatted
	cmd := timeCmd
	cmd.SetArgs([]string{"--help"})

	// This should not error
	err := cmd.Execute()
	if err != nil {
		t.Errorf("time command help should not error: %v", err)
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
			// Reset flags for each test
			timeType = "now"
			timeFormat = "iso"
			timezone = "utc"
			timeInput = ""
			timeFrom = ""
			timeTo = ""
			timeOffset = ""

			// Parse flags manually
			cmd := timeCmd
			cmd.ParseFlags(tt.args)

			// Execute the command directly
			err := runTime(cmd, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("time command error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
