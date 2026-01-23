// Package main provides the mcpipboy CLI application
package main

import (
	"os"

	"github.com/kluzzebass/mcpipboy"
	"github.com/spf13/cobra"
)

// Version information - set via ldflags at build time
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
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
	Version: getVersionInfo(),
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// getVersionInfo returns detailed version information
func getVersionInfo() string {
	v := version
	// If version was not set by build process, use library version
	if v == "dev" {
		v = mcpipboy.Version()
	}

	// Add commit and date if available
	if commit != "none" || date != "unknown" {
		v += "\n"
		if commit != "none" {
			v += "  commit: " + commit + "\n"
		}
		if date != "unknown" {
			v += "  built:  " + date
		}
	}
	return v
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
