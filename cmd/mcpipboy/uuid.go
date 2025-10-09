// Package main provides the uuid command for mcpipboy
package main

import (
	"fmt"

	"github.com/kluzzebass/mcpipboy/internal/tools"
	"github.com/spf13/cobra"
)

var uuidCmd = &cobra.Command{
	Use:   "uuid [flags]",
	Short: "Generate and validate UUIDs with various versions",
	Long: `UUID generator and validator provides comprehensive UUID functionality including:

- UUID v1 (time-based) generation
- UUID v4 (random) generation  
- UUID v5 (name-based SHA-1) generation
- UUID v7 (time-ordered) generation
- UUID validation for any version
- Batch generation with count parameter

Examples:
  mcpipboy uuid --version v4
  mcpipboy uuid --version v7 --count 10
  mcpipboy uuid --version v5 --namespace "6ba7b810-9dad-11d1-80b4-00c04fd430c8" --name "example"
  mcpipboy uuid --version validate --input "550e8400-e29b-41d4-a716-446655440000"`,
	RunE: runUUID,
}

var (
	uuidVersion   string
	uuidCount     int
	uuidNamespace string
	uuidName      string
	uuidInput     string
)

func init() {
	// Set group ID for uuid command
	uuidCmd.GroupID = "tools"

	uuidCmd.Flags().StringVar(&uuidVersion, "version", "v4", "UUID version: v1, v4, v5, v7, validate")
	uuidCmd.Flags().IntVar(&uuidCount, "count", 1, "Number of UUIDs to generate (1-1000)")
	uuidCmd.Flags().StringVar(&uuidNamespace, "namespace", "", "Namespace UUID for v5 generation")
	uuidCmd.Flags().StringVar(&uuidName, "name", "", "Name for v5 generation")
	uuidCmd.Flags().StringVar(&uuidInput, "input", "", "UUID string to validate")

	// Add command to root
	rootCmd.AddCommand(uuidCmd)
}

func runUUID(cmd *cobra.Command, args []string) error {
	// Build parameters map
	params := make(map[string]interface{})

	// Add parameters if they were set
	if uuidVersion != "" {
		params["version"] = uuidVersion
	}
	params["count"] = float64(uuidCount)
	if uuidNamespace != "" {
		params["namespace"] = uuidNamespace
	}
	if uuidName != "" {
		params["name"] = uuidName
	}
	if uuidInput != "" {
		params["input"] = uuidInput
	}

	// Create and execute the UUID tool
	tool := tools.NewUUIDTool()

	// Validate parameters
	if err := tool.ValidateParams(params); err != nil {
		return fmt.Errorf("parameter validation failed: %v", err)
	}

	// Execute the tool
	result, err := tool.Execute(params)
	if err != nil {
		return fmt.Errorf("UUID tool execution failed: %v", err)
	}

	// Print the result
	fmt.Println(result)
	return nil
}
