package tools

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ijt/go-anytime"
)

// TimeTool implements comprehensive time functionality
type TimeTool struct{}

// NewTimeTool creates a new time tool instance
func NewTimeTool() *TimeTool {
	return &TimeTool{}
}

// Name returns the tool's name
func (t *TimeTool) Name() string {
	return "time"
}

// Description returns the tool's description
func (t *TimeTool) Description() string {
	return "Comprehensive time utility with parsing, formatting, and calculations. Note: Dates before year 1000 are not supported due to parsing library limitations."
}

// Execute runs the time tool
func (t *TimeTool) Execute(params map[string]interface{}) (interface{}, error) {
	// Parse parameters
	format, _ := params["format"].(string)
	if format == "" {
		format = "iso"
	}

	timezone, _ := params["timezone"].(string)
	if timezone == "" {
		timezone = "local"
	}

	input, _ := params["input"].(string)
	offset, _ := params["offset"].(string)

	// Get base time - use go-anytime for all input parsing
	var baseTime time.Time
	var err error

	if input != "" {
		// Parse the input timestamp using go-anytime
		// Use UTC as reference to avoid timezone confusion
		baseTime, err = anytime.Parse(input, time.Now().UTC())
		if err != nil {
			// Provide helpful hints for common parsing issues
			errorMsg := fmt.Sprintf("failed to parse timestamp: %v", err)

			// Add hints for common issues
			if strings.Contains(err.Error(), "expected natural date") {
				errorMsg += "\n\nHint: If you were trying to parse a year prior to 1000, note that dates before year 1000 are not supported. Otherwise, try using a more standard date format like 'YYYY-MM-DD' or 'January 1, 2025'."
			} else if strings.Contains(err.Error(), "left unparsed") {
				errorMsg += "\n\nHint: Try using a more standard date format like 'YYYY-MM-DD' or 'January 1, 2025'."
			}

			return nil, fmt.Errorf("%s", errorMsg)
		}
	} else {
		// No input provided, use current time
		baseTime = time.Now()
	}

	// Apply timezone
	if timezone != "local" {
		if timezone == "utc" {
			baseTime = baseTime.UTC()
		} else {
			loc, err := time.LoadLocation(timezone)
			if err != nil {
				return nil, fmt.Errorf("invalid timezone: %s", timezone)
			}
			baseTime = baseTime.In(loc)
		}
	}

	// Apply offset if provided
	if offset != "" {
		duration, err := time.ParseDuration(offset)
		if err != nil {
			return nil, fmt.Errorf("invalid offset format: %s", offset)
		}
		baseTime = baseTime.Add(duration)
	}

	// Format the result
	result, err := t.formatTime(baseTime, format)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// formatTime formats a time according to the specified format
func (t *TimeTool) formatTime(tm time.Time, format string) (string, error) {
	switch format {
	case "iso":
		return tm.Format(time.RFC3339), nil
	case "rfc3339":
		return tm.Format(time.RFC3339), nil
	case "unix":
		return strconv.FormatInt(tm.Unix(), 10), nil
	case "date":
		return tm.Format("2006-01-02"), nil
	case "datetime":
		return tm.Format("2006-01-02 15:04:05"), nil
	case "time":
		return tm.Format("15:04:05"), nil
	case "weekday":
		return tm.Format("Monday, January 2, 2006"), nil
	default:
		return "", fmt.Errorf("invalid format: %s", format)
	}
}

// ValidateParams validates the input parameters
func (t *TimeTool) ValidateParams(params map[string]interface{}) error {
	// Type parameter is no longer needed - go-anytime handles everything

	// Validate format
	if format, ok := params["format"]; ok {
		if formatStr, ok := format.(string); ok {
			validFormats := []string{"iso", "rfc3339", "unix", "date", "datetime", "time", "weekday"}
			if !contains(validFormats, formatStr) {
				return fmt.Errorf("invalid format: %s, must be one of: %s", formatStr, strings.Join(validFormats, ", "))
			}
		} else {
			return fmt.Errorf("format parameter must be a string")
		}
	}

	// Validate timezone
	if timezone, ok := params["timezone"]; ok {
		if timezoneStr, ok := timezone.(string); ok {
			if timezoneStr != "utc" && timezoneStr != "local" {
				// Try to load the timezone to validate it
				if _, err := time.LoadLocation(timezoneStr); err != nil {
					return fmt.Errorf("invalid timezone: %s", timezoneStr)
				}
			}
		} else {
			return fmt.Errorf("timezone parameter must be a string")
		}
	}

	// Validate offset
	if offset, ok := params["offset"]; ok {
		if offsetStr, ok := offset.(string); ok {
			if _, err := time.ParseDuration(offsetStr); err != nil {
				return fmt.Errorf("invalid offset format: %s", offsetStr)
			}
		} else {
			return fmt.Errorf("offset parameter must be a string")
		}
	}

	return nil
}

// GetInputSchema returns the JSON schema for tool input parameters
func (t *TimeTool) GetInputSchema() map[string]interface{} {
	return CreateJSONSchema([]ParameterDefinition{
		{
			Name:        "format",
			Type:        "string",
			Description: "Output format: iso, rfc3339, unix, date, datetime, time, weekday",
			Required:    false,
		},
		{
			Name:        "timezone",
			Type:        "string",
			Description: "Timezone: utc, local, or IANA timezone name",
			Required:    false,
		},
		{
			Name:        "input",
			Type:        "string",
			Description: "Input timestamp (any format supported by go-anytime). Note: Dates before year 1000 are not supported.",
			Required:    false,
		},
		{
			Name:        "offset",
			Type:        "string",
			Description: "Time offset (e.g., +1h, -2d, +30m)",
			Required:    false,
		},
	})
}

// GetOutputSchema returns the JSON schema for tool output
func (t *TimeTool) GetOutputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"result": map[string]interface{}{
				"type":        "string",
				"description": "Formatted time string or duration information",
			},
		},
	}
}

// GetResources returns the list of resources this tool provides
func (t *TimeTool) GetResources() []Resource {
	return []Resource{
		{
			Name:     "Time Formats",
			URI:      "time://formats",
			MIMEType: "application/json",
		},
		{
			Name:     "Time Examples",
			URI:      "time://examples",
			MIMEType: "application/json",
		},
	}
}

// ReadResource reads a specific resource by URI
func (t *TimeTool) ReadResource(uri string) (string, error) {
	switch uri {
	case "time://formats":
		// Return supported output formats
		formats := []map[string]interface{}{
			{
				"format":      "iso",
				"name":        "ISO 8601",
				"description": "ISO 8601 format (2006-01-02T15:04:05Z07:00)",
				"example":     "2025-01-09T14:30:00Z",
			},
			{
				"format":      "rfc3339",
				"name":        "RFC 3339",
				"description": "RFC 3339 format (2006-01-02T15:04:05Z07:00)",
				"example":     "2025-01-09T14:30:00Z",
			},
			{
				"format":      "unix",
				"name":        "Unix Timestamp",
				"description": "Unix timestamp in seconds",
				"example":     "1736430600",
			},
			{
				"format":      "date",
				"name":        "Date Only",
				"description": "Date in YYYY-MM-DD format",
				"example":     "2025-01-09",
			},
			{
				"format":      "datetime",
				"name":        "Date and Time",
				"description": "Date and time in YYYY-MM-DD HH:MM:SS format",
				"example":     "2025-01-09 14:30:00",
			},
			{
				"format":      "time",
				"name":        "Time Only",
				"description": "Time in HH:MM:SS format",
				"example":     "14:30:00",
			},
			{
				"format":      "weekday",
				"name":        "Day of Week",
				"description": "Day of the week",
				"example":     "Thursday",
			},
		}
		jsonData, err := json.Marshal(formats)
		if err != nil {
			return "", fmt.Errorf("failed to marshal formats: %w", err)
		}
		return string(jsonData), nil
	case "time://examples":
		// Return example time strings that can be parsed
		examples := []map[string]interface{}{
			{
				"input":       "now",
				"description": "Current time",
			},
			{
				"input":       "today",
				"description": "Today at midnight",
			},
			{
				"input":       "2025-01-09",
				"description": "Date in YYYY-MM-DD format",
			},
			{
				"input":       "2025-01-09T14:30:00Z",
				"description": "ISO 8601 format with UTC timezone",
			},
			{
				"input":       "January 9, 2025",
				"description": "Natural language date",
			},
			{
				"input":       "yesterday",
				"description": "Yesterday",
			},
			{
				"input":       "tomorrow",
				"description": "Tomorrow",
			},
			{
				"input":       "next Monday",
				"description": "Next occurrence of Monday",
			},
			{
				"input":       "last Friday",
				"description": "Last occurrence of Friday",
			},
			{
				"input":       "3 days ago",
				"description": "Relative time expression",
			},
		}
		jsonData, err := json.Marshal(examples)
		if err != nil {
			return "", fmt.Errorf("failed to marshal examples: %w", err)
		}
		return string(jsonData), nil
	default:
		return "", fmt.Errorf("unknown resource URI: %s", uri)
	}
}

// contains checks if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
