// Package main provides the version command for mcpipboy
package main

import (
	"fmt"
	"io"
	"os"

	"github.com/kluzzebass/mcpipboy/internal/tools"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display the current version",
	Long:  `Display the current version of mcpipboy.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runVersion(cmd, args, os.Stdout)
	},
}

func init() {
	// Set group ID for version command
	versionCmd.GroupID = "tools"

	rootCmd.AddCommand(versionCmd)
}

func runVersion(cmd *cobra.Command, args []string, out io.Writer) error {
	// Create version tool instance
	versionTool := tools.NewVersionTool()

	// Prepare parameters (version tool doesn't need any)
	params := map[string]interface{}{}

	// Validate parameters
	if err := versionTool.ValidateParams(params); err != nil {
		return fmt.Errorf("invalid parameters: %v", err)
	}

	// Execute the tool
	result, err := versionTool.Execute(params)
	if err != nil {
		return fmt.Errorf("version execution failed: %v", err)
	}

	// Output the result
	fmt.Fprintln(out, result)
	return nil
}
