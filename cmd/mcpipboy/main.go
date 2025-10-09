// Package main provides the mcpipboy CLI application
package main

import (
	"os"

	version "github.com/kluzzebass/mcpipboy"
	"github.com/spf13/cobra"

	// Import all command files to ensure their init() functions run
	_ "github.com/kluzzebass/mcpipboy/cmd/mcpipboy"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mcpipboy",
	Short: "MCP server for AI agent utility tools",
	Long: `Flawlessly generated, rigorously validated - dependable data for the digital wasteland.

mcpipboy is an MCP (Model Context Protocol) server that provides agentic AIs 
with essential tools for common tasks they struggle with, including UUID generation, 
checksummed identifier verification/generation (IMO, MMSI, credit card numbers, ISBN, etc.), 
and other utility functions.`,
	Version: version.Version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Add command groups
	rootCmd.AddGroup(&cobra.Group{
		ID:    "server",
		Title: "Server Commands:",
	})

	rootCmd.AddGroup(&cobra.Group{
		ID:    "tools",
		Title: "Tool Commands:",
	})

	// Add global flags here if needed
}

func main() {
	Execute()
}
