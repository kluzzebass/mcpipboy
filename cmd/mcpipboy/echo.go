// Package main provides the echo command for mcpipboy
package main

import (
	"fmt"

	"github.com/kluzzebass/mcpipboy/internal/tools"
	"github.com/spf13/cobra"
)

// echoCmd represents the echo command
var echoCmd = &cobra.Command{
	Use:   "echo [message]",
	Short: "Echo back a message",
	Long:  `Echo back the provided message. This is a simple test tool.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runEcho,
}

func init() {
	// Set group ID for echo command
	echoCmd.GroupID = "tools"

	rootCmd.AddCommand(echoCmd)
}

func runEcho(cmd *cobra.Command, args []string) error {
	// Create echo tool instance
	echoTool := tools.NewEchoTool()

	// Prepare parameters
	params := map[string]interface{}{
		"message": args[0],
	}

	// Validate parameters
	if err := echoTool.ValidateParams(params); err != nil {
		return fmt.Errorf("invalid parameters: %v", err)
	}

	// Execute the tool
	result, err := echoTool.Execute(params)
	if err != nil {
		return fmt.Errorf("echo execution failed: %v", err)
	}

	// Output the result
	fmt.Println(result)
	return nil
}
