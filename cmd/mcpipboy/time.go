package main

import (
	"fmt"

	"github.com/kluzzebass/mcpipboy/internal/tools"
	"github.com/spf13/cobra"
)

var timeCmd = &cobra.Command{
	Use:   "time [flags]",
	Short: "Comprehensive time utility with parsing, formatting, and calculations",
	Long: `Time tool provides comprehensive time functionality including:

- Current time in various formats
- Time parsing with lenient input handling
- Timezone conversions
- Relative time calculations
- Time offset operations

Examples:
  mcpipboy time                                    # Current time in ISO format
  mcpipboy time --type today --format date         # Today's date
  mcpipboy time --type timestamp --input "2025-01-01" --format unix
  mcpipboy time --type relative --from "2025-01-01" --to "2025-12-31"
  mcpipboy time --type now --offset "+1h" --format datetime
  mcpipboy time --type timestamp --input "2025-01-01T12:00:00Z" --timezone "America/New_York"`,
	GroupID: "tools",
	Args:    cobra.ExactArgs(0),
	RunE:    runTime,
}

var (
	timeType   string
	timeFormat string
	timezone   string
	timeInput  string
	timeFrom   string
	timeTo     string
	timeOffset string
)

func init() {
	timeCmd.Flags().StringVar(&timeType, "type", "now", "Time type: now, today, timestamp, unix, relative")
	timeCmd.Flags().StringVar(&timeFormat, "format", "iso", "Output format: iso, rfc3339, unix, date, datetime, time")
	timeCmd.Flags().StringVar(&timezone, "timezone", "utc", "Timezone: utc, local, or IANA timezone name")
	timeCmd.Flags().StringVar(&timeInput, "input", "", "Input timestamp (required for timestamp/unix types)")
	timeCmd.Flags().StringVar(&timeFrom, "from", "", "Start timestamp for relative calculations")
	timeCmd.Flags().StringVar(&timeTo, "to", "", "End timestamp for relative calculations")
	timeCmd.Flags().StringVar(&timeOffset, "offset", "", "Time offset (e.g., +1h, -2d, +30m)")

	// Validation is handled in runTime function
}

func runTime(cmd *cobra.Command, args []string) error {
	// Build parameters map
	params := make(map[string]interface{})
	
	if timeType != "" {
		params["type"] = timeType
	}
	if timeFormat != "" {
		params["format"] = timeFormat
	}
	if timezone != "" {
		params["timezone"] = timezone
	}
	if timeInput != "" {
		params["input"] = timeInput
	}
	if timeFrom != "" {
		params["from"] = timeFrom
	}
	if timeTo != "" {
		params["to"] = timeTo
	}
	if timeOffset != "" {
		params["offset"] = timeOffset
	}

	// Create and execute the tool
	tool := tools.NewTimeTool()
	
	// Validate parameters
	if err := tool.ValidateParams(params); err != nil {
		return fmt.Errorf("parameter validation failed: %v", err)
	}

	// Execute the tool
	result, err := tool.Execute(params)
	if err != nil {
		return fmt.Errorf("time tool execution failed: %v", err)
	}

	// Output the result
	if resultMap, ok := result.(map[string]interface{}); ok {
		// Handle relative time results
		if duration, ok := resultMap["duration"].(string); ok {
			fmt.Printf("Duration: %s\n", duration)
			if seconds, ok := resultMap["seconds"].(float64); ok {
				fmt.Printf("Seconds: %.0f\n", seconds)
			}
			if from, ok := resultMap["from"].(string); ok {
				fmt.Printf("From: %s\n", from)
			}
			if to, ok := resultMap["to"].(string); ok {
				fmt.Printf("To: %s\n", to)
			}
		}
	} else if resultStr, ok := result.(string); ok {
		// Handle simple string results
		fmt.Println(resultStr)
	} else {
		fmt.Printf("%v\n", result)
	}

	return nil
}
