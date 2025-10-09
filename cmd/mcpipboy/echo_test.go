// Package main provides tests for the echo command
package main

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestEchoCommand(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedOutput string
		expectError    bool
	}{
		{
			name:           "echo simple message",
			args:           []string{"echo", "test message"},
			expectedOutput: "test message",
			expectError:    false,
		},
		{
			name:           "echo with special characters",
			args:           []string{"echo", "Hello, World!"},
			expectedOutput: "Hello, World!",
			expectError:    false,
		},
		{
			name:        "echo with no arguments",
			args:        []string{"echo"},
			expectError: true,
		},
		{
			name:        "echo with multiple arguments",
			args:        []string{"echo", "arg1", "arg2"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Execute the CLI via go run
			args := append([]string{"run", "."}, tt.args...)
			cmd := exec.Command("go", args...)
			output, err := cmd.CombinedOutput()

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but command succeeded")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v\nOutput: %s", err, string(output))
				return
			}

			outputStr := strings.TrimSpace(string(output))
			if !strings.Contains(outputStr, tt.expectedOutput) {
				t.Errorf("Expected output to contain %q, got %q", tt.expectedOutput, outputStr)
			}
		})
	}
}

// TestRunEcho tests the runEcho function directly with buffer (for coverage)
func TestRunEcho(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedOutput string
		expectError    bool
	}{
		{
			name:           "echo simple message",
			args:           []string{"test message"},
			expectedOutput: "test message\n",
			expectError:    false,
		},
		{
			name:           "echo with special characters",
			args:           []string{"Hello, World!"},
			expectedOutput: "Hello, World!\n",
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a buffer to capture output
			var buf bytes.Buffer

			// Call runEcho directly
			err := runEcho(nil, tt.args, &buf)

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

			// Check the output
			output := buf.String()
			if output != tt.expectedOutput {
				t.Errorf("Expected output %q, got %q", tt.expectedOutput, output)
			}
		})
	}
}
