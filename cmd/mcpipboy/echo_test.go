// Package main provides tests for the echo command
package main

import (
	"testing"
)

func TestEchoCommand(t *testing.T) {
	// Test that echo command is properly configured
	if echoCmd.Use != "echo [message]" {
		t.Errorf("Expected echo command use to be 'echo [message]', got '%s'", echoCmd.Use)
	}

	if echoCmd.Short == "" {
		t.Error("Echo command should have a short description")
	}
}

func TestRunEcho(t *testing.T) {
	// Test valid echo execution
	err := runEcho(echoCmd, []string{"test message"})
	if err != nil {
		t.Errorf("runEcho should not error with valid input: %v", err)
	}
}

func TestRunEchoInvalidArgs(t *testing.T) {
	// Test that echo requires exactly one argument
	// Note: Cobra handles argument validation, so we test the command validation
	if echoCmd.Args == nil {
		t.Error("Echo command should have argument validation")
	}

	// Test that the command validates arguments correctly
	err := echoCmd.Args(echoCmd, []string{})
	if err == nil {
		t.Error("Echo command should error with no arguments")
	}

	err = echoCmd.Args(echoCmd, []string{"arg1", "arg2"})
	if err == nil {
		t.Error("Echo command should error with multiple arguments")
	}
}
