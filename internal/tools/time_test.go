package tools

import (
	"testing"
	"time"
)

func TestTimeToolExecute(t *testing.T) {
	tool := NewTimeTool()

	tests := []struct {
		name     string
		params   map[string]interface{}
		wantErr  bool
		validate func(t *testing.T, result interface{})
	}{
		{
			name:    "default now",
			params:  map[string]interface{}{},
			wantErr: false,
			validate: func(t *testing.T, result interface{}) {
				if result == nil {
					t.Error("Expected non-nil result")
				}
				if _, ok := result.(string); !ok {
					t.Errorf("Expected string result, got %T", result)
				}
			},
		},
		{
			name: "now with iso format",
			params: map[string]interface{}{
				"type":   "now",
				"format": "iso",
			},
			wantErr: false,
			validate: func(t *testing.T, result interface{}) {
				if result == nil {
					t.Error("Expected non-nil result")
				}
				if str, ok := result.(string); ok {
					if _, err := time.Parse(time.RFC3339, str); err != nil {
						t.Errorf("Expected valid RFC3339 format, got: %s", str)
					}
				} else {
					t.Errorf("Expected string result, got %T", result)
				}
			},
		},
		{
			name: "today with date format",
			params: map[string]interface{}{
				"type":   "today",
				"format": "date",
			},
			wantErr: false,
			validate: func(t *testing.T, result interface{}) {
				if result == nil {
					t.Error("Expected non-nil result")
				}
				if str, ok := result.(string); ok {
					if _, err := time.Parse("2006-01-02", str); err != nil {
						t.Errorf("Expected valid date format, got: %s", str)
					}
				} else {
					t.Errorf("Expected string result, got %T", result)
				}
			},
		},
		{
			name: "timestamp with valid input",
			params: map[string]interface{}{
				"input":  "2025-01-01T12:00:00Z",
				"format": "unix",
			},
			wantErr: false,
			validate: func(t *testing.T, result interface{}) {
				if result == nil {
					t.Error("Expected non-nil result")
				}
				if str, ok := result.(string); ok {
					if str != "1735732800" {
						t.Errorf("Expected unix timestamp 1735732800, got: %s", str)
					}
				} else {
					t.Errorf("Expected string result, got %T", result)
				}
			},
		},
		{
			name: "unix timestamp",
			params: map[string]interface{}{
				"input":  "2025-01-01T12:00:00Z",
				"format": "unix",
			},
			wantErr: false,
			validate: func(t *testing.T, result interface{}) {
				if result == nil {
					t.Error("Expected non-nil result")
				}
				if str, ok := result.(string); ok {
					if str != "1735732800" {
						t.Errorf("Expected unix timestamp '1735732800', got: %s", str)
					}
				} else {
					t.Errorf("Expected string result, got %T", result)
				}
			},
		},
		{
			name: "offset calculation",
			params: map[string]interface{}{
				"input":  "2025-01-01T00:00:00Z",
				"offset": "+1h",
				"format": "datetime",
			},
			wantErr: false,
			validate: func(t *testing.T, result interface{}) {
				if result == nil {
					t.Error("Expected non-nil result")
				}
				if str, ok := result.(string); ok {
					if str != "2025-01-01 01:00:00" {
						t.Errorf("Expected datetime '2025-01-01 01:00:00', got: %s", str)
					}
				} else {
					t.Errorf("Expected string result, got %T", result)
				}
			},
		},
		{
			name: "timezone conversion",
			params: map[string]interface{}{
				"input":    "2025-01-01T12:00:00Z",
				"timezone": "America/New_York",
				"format":   "datetime",
			},
			wantErr: false,
			validate: func(t *testing.T, result interface{}) {
				if result == nil {
					t.Error("Expected non-nil result")
				}
				if str, ok := result.(string); ok {
					// Should be 7:00 AM in New York (UTC-5 in January)
					if str != "2025-01-01 07:00:00" {
						t.Errorf("Expected datetime '2025-01-01 07:00:00', got: %s", str)
					}
				} else {
					t.Errorf("Expected string result, got %T", result)
				}
			},
		},
		{
			name: "invalid offset",
			params: map[string]interface{}{
				"offset": "invalid",
			},
			wantErr: true,
		},
		{
			name: "invalid timezone",
			params: map[string]interface{}{
				"timezone": "Invalid/Timezone",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tool.Execute(tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.validate != nil {
				tt.validate(t, result)
			}
		})
	}
}

func TestTimeToolValidateParams(t *testing.T) {
	tool := NewTimeTool()

	tests := []struct {
		name    string
		params  map[string]interface{}
		wantErr bool
	}{
		{
			name:    "valid empty params",
			params:  map[string]interface{}{},
			wantErr: false,
		},
		{
			name: "valid format",
			params: map[string]interface{}{
				"format": "iso",
			},
			wantErr: false,
		},
		{
			name: "invalid format",
			params: map[string]interface{}{
				"format": "invalid",
			},
			wantErr: true,
		},
		{
			name: "valid timezone",
			params: map[string]interface{}{
				"timezone": "UTC",
			},
			wantErr: false,
		},
		{
			name: "invalid timezone",
			params: map[string]interface{}{
				"timezone": "Invalid/Timezone",
			},
			wantErr: true,
		},
		{
			name: "invalid parameter types",
			params: map[string]interface{}{
				"format": 456,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tool.ValidateParams(tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateParams() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTimeToolSchemas(t *testing.T) {
	tool := NewTimeTool()

	// Test input schema
	inputSchema := tool.GetInputSchema()
	if inputSchema == nil {
		t.Error("Expected non-nil input schema")
	}

	// Test output schema
	outputSchema := tool.GetOutputSchema()
	if outputSchema == nil {
		t.Error("Expected non-nil output schema")
	}

	// Verify schema structure
	if schemaType, ok := outputSchema["type"].(string); !ok || schemaType != "object" {
		t.Error("Expected output schema type to be 'object'")
	}
}

func TestTimeToolEdgeCases(t *testing.T) {
	tool := NewTimeTool()

	tests := []struct {
		name     string
		params   map[string]interface{}
		wantErr  bool
		validate func(t *testing.T, result interface{})
	}{
		{
			name: "lenient timestamp parsing",
			params: map[string]interface{}{
				"input":  "January 1, 2025 at 12:00 PM",
				"format": "date",
			},
			wantErr: false,
			validate: func(t *testing.T, result interface{}) {
				if result == nil {
					t.Error("Expected non-nil result")
				}
				if str, ok := result.(string); ok {
					if str != "2025-01-01" {
						t.Errorf("Expected date '2025-01-01', got: %s", str)
					}
				} else {
					t.Errorf("Expected string result, got %T", result)
				}
			},
		},
		{
			name: "negative offset",
			params: map[string]interface{}{
				"input":  "2025-01-01T12:00:00Z",
				"offset": "-2h",
				"format": "datetime",
			},
			wantErr: false,
			validate: func(t *testing.T, result interface{}) {
				if result == nil {
					t.Error("Expected non-nil result")
				}
				if str, ok := result.(string); ok {
					if str != "2025-01-01 10:00:00" {
						t.Errorf("Expected datetime '2025-01-01 10:00:00', got: %s", str)
					}
				} else {
					t.Errorf("Expected string result, got %T", result)
				}
			},
		},
		{
			name: "complex offset",
			params: map[string]interface{}{
				"input":  "2025-01-01T00:00:00Z",
				"offset": "26h30m",
				"format": "datetime",
			},
			wantErr: false,
			validate: func(t *testing.T, result interface{}) {
				if result == nil {
					t.Error("Expected non-nil result")
				}
				if str, ok := result.(string); ok {
					if str != "2025-01-02 02:30:00" {
						t.Errorf("Expected datetime '2025-01-02 02:30:00', got: %s", str)
					}
				} else {
					t.Errorf("Expected string result, got %T", result)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tool.Execute(tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.validate != nil {
				tt.validate(t, result)
			}
		})
	}
}
