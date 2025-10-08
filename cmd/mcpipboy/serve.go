// Package main provides the serve command for mcpipboy
package main

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/kluzzebass/mcpipboy/internal/server"
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
func getAvailableTools() *tools.ToolRegistry {
	registry := tools.NewToolRegistry()

	// Register all available tools
	registry.RegisterTool(tools.NewEchoTool())
	registry.RegisterTool(tools.NewVersionTool())
	registry.RegisterTool(tools.NewTimeTool())
	// TODO: Add more tools as they are implemented

	return registry
}

// toolCompletionFunc provides shell completion for tool names
func toolCompletionFunc(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	registry := getAvailableTools()
	return registry.ListTools(), cobra.ShellCompDirectiveNoFileComp
}

func runServe(cmd *cobra.Command, args []string) error {
	// Validate that enable and disable are not both used
	if len(enableTools) > 0 && len(disableTools) > 0 {
		return fmt.Errorf("--enable and --disable flags are mutually exclusive")
	}

	// Get available tools for validation
	registry := getAvailableTools()
	availableTools := registry.ListTools()

	// Validate enable tools
	if len(enableTools) > 0 {
		for _, tool := range enableTools {
			if !slices.Contains(availableTools, tool) {
				return fmt.Errorf("invalid tool: %s. Available tools: %s", tool, strings.Join(availableTools, ", "))
			}
		}
	}

	// Validate disable tools
	if len(disableTools) > 0 {
		for _, tool := range disableTools {
			if !slices.Contains(availableTools, tool) {
				return fmt.Errorf("invalid tool: %s. Available tools: %s", tool, strings.Join(availableTools, ", "))
			}
		}
	}

	// Use the same registry we created for validation
	// (registry is already created above with all tools registered)

	// Determine which tools to enable
	var enabledTools []string
	if len(enableTools) > 0 {
		// Only enable specified tools
		enabledTools = enableTools
	} else if len(disableTools) > 0 {
		// Enable all tools except disabled ones
		for _, tool := range availableTools {
			if !slices.Contains(disableTools, tool) {
				enabledTools = append(enabledTools, tool)
			}
		}
	} else {
		// Enable all tools by default
		enabledTools = availableTools
	}

	// Create MCP server
	server := server.NewServer()

	// Register only the enabled tools with the MCP server
	for _, toolName := range enabledTools {
		if tool, exists := registry.GetTool(toolName); exists {
			server.RegisterTool(tool)
		}
	}

	// Start the MCP server
	ctx := context.Background()
	return server.Start(ctx)
}
