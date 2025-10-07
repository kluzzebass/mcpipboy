// Package main provides the serve command for mcpipboy
package main

import (
	"fmt"
	"strings"

	"github.com/kluzzebass/mcpipboy/internal/tools"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the MCP server",
	Long: `Start the MCP (Model Context Protocol) server that provides tools to AI agents.
The server communicates via stdin/stdout and can be configured to enable or disable specific tools.`,
	RunE: runServe,
}

var (
	enableTools  []string
	disableTools []string
)

func init() {
	// Set group ID for serve command
	serveCmd.GroupID = "server"

	rootCmd.AddCommand(serveCmd)

	// Add flags for tool enable/disable
	serveCmd.Flags().StringSliceVar(&enableTools, "enable", []string{}, "Comma-separated list of tools to enable (mutually exclusive with --disable)")
	serveCmd.Flags().StringSliceVar(&disableTools, "disable", []string{}, "Comma-separated list of tools to disable (mutually exclusive with --enable)")

	// Mark flags as mutually exclusive
	serveCmd.MarkFlagsMutuallyExclusive("enable", "disable")

	// Add dynamic completion for tool names
	serveCmd.RegisterFlagCompletionFunc("enable", toolCompletionFunc)
	serveCmd.RegisterFlagCompletionFunc("disable", toolCompletionFunc)
}

// getAvailableTools returns a registry with all available tools registered
func getAvailableTools() *tools.Registry {
	registry := tools.NewRegistry()

	// Register all available tools
	registry.Register(tools.NewEchoTool())
	registry.Register(tools.NewVersionTool())
	// TODO: Add more tools as they are implemented

	return registry
}

// toolCompletionFunc provides shell completion for tool names
func toolCompletionFunc(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	registry := getAvailableTools()
	return registry.List(), cobra.ShellCompDirectiveNoFileComp
}

func runServe(cmd *cobra.Command, args []string) error {
	// Validate that enable and disable are not both used
	if len(enableTools) > 0 && len(disableTools) > 0 {
		return fmt.Errorf("--enable and --disable flags are mutually exclusive")
	}

	// Get available tools for validation
	registry := getAvailableTools()
	availableTools := registry.List()

	// Validate enable tools
	if len(enableTools) > 0 {
		for _, tool := range enableTools {
			if !contains(availableTools, tool) {
				return fmt.Errorf("invalid tool: %s. Available tools: %s", tool, strings.Join(availableTools, ", "))
			}
		}
	}

	// Validate disable tools
	if len(disableTools) > 0 {
		for _, tool := range disableTools {
			if !contains(availableTools, tool) {
				return fmt.Errorf("invalid tool: %s. Available tools: %s", tool, strings.Join(availableTools, ", "))
			}
		}
	}

	// TODO: Implement MCP server startup
	fmt.Printf("Starting MCP server...\n")
	fmt.Printf("Version: %s\n", rootCmd.Version)
	fmt.Printf("Available tools: %s\n", strings.Join(availableTools, ", "))

	if len(enableTools) > 0 {
		fmt.Printf("Enabled tools: %s\n", strings.Join(enableTools, ", "))
	} else if len(disableTools) > 0 {
		fmt.Printf("Disabled tools: %s\n", strings.Join(disableTools, ", "))
	} else {
		fmt.Println("All tools enabled by default")
	}

	// TODO: Start actual MCP server
	return fmt.Errorf("MCP server not yet implemented")
}

// contains checks if a string slice contains a specific string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
