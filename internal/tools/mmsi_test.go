package tools

import (
	"strings"
	"testing"
)

func TestMMSIToolName(t *testing.T) {
	tool := NewMMSITool()
	if tool.Name() != "mmsi" {
		t.Errorf("Expected name 'mmsi', got '%s'", tool.Name())
	}
}

func TestMMSIToolDescription(t *testing.T) {
	tool := NewMMSITool()
	desc := tool.Description()
	if desc == "" {
		t.Error("Description should not be empty")
	}
	if !strings.Contains(desc, "MMSI") {
		t.Error("Description should contain 'MMSI'")
	}
}

func TestMMSIToolValidateParams(t *testing.T) {
	tool := NewMMSITool()

	tests := []struct {
		name    string
		params  map[string]interface{}
		wantErr bool
	}{
		{
			name:    "valid_validate_operation",
			params:  map[string]interface{}{"operation": "validate", "input": "123456789"},
			wantErr: false,
		},
		{
			name:    "valid_generate_operation",
			params:  map[string]interface{}{"operation": "generate", "country-code": "US", "count": 5},
			wantErr: false,
		},
		{
			name:    "invalid_operation",
			params:  map[string]interface{}{"operation": "invalid"},
			wantErr: true,
		},
		{
			name:    "count_too_low",
			params:  map[string]interface{}{"operation": "generate", "count": 0},
			wantErr: true,
		},
		{
			name:    "count_too_high",
			params:  map[string]interface{}{"operation": "generate", "count": 101},
			wantErr: true,
		},
		{
			name:    "validate_missing_input",
			params:  map[string]interface{}{"operation": "validate"},
			wantErr: true,
		},
		{
			name:    "validate_input_not_string",
			params:  map[string]interface{}{"operation": "validate", "input": 123456789},
			wantErr: true,
		},
		{
			name:    "invalid_country_code",
			params:  map[string]interface{}{"operation": "generate", "country-code": "XX"},
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

func TestMMSIToolValidateMMSI(t *testing.T) {
	tool := NewMMSITool()

	tests := []struct {
		name     string
		input    string
		expected bool
		hasError bool
	}{
		{
			name:     "valid_mmsi",
			input:    "366123456",
			expected: true,
			hasError: false,
		},
		{
			name:     "valid_mmsi_with_spaces",
			input:    "366 123 456",
			expected: true,
			hasError: false,
		},
		{
			name:     "valid_mmsi_with_dashes",
			input:    "366-123-456",
			expected: true,
			hasError: false,
		},
		{
			name:     "too_short",
			input:    "12345678",
			expected: false,
			hasError: false,
		},
		{
			name:     "too_long",
			input:    "1234567890",
			expected: false,
			hasError: false,
		},
		{
			name:     "non_numeric",
			input:    "12345678a",
			expected: false,
			hasError: false,
		},
		{
			name:     "too_small",
			input:    "12345678",
			expected: false,
			hasError: false,
		},
		{
			name:     "too_large",
			input:    "1234567890",
			expected: false,
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := map[string]interface{}{
				"operation": "validate",
				"input":     tt.input,
			}

			result, err := tool.validateMMSI(params)

			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			resultMap, ok := result.(map[string]interface{})
			if !ok {
				t.Errorf("Expected map result, got %T", result)
				return
			}

			valid, ok := resultMap["valid"].(bool)
			if !ok {
				t.Errorf("Expected valid field to be bool, got %T", resultMap["valid"])
				return
			}

			if valid != tt.expected {
				t.Errorf("Expected valid=%v, got %v", tt.expected, valid)
			}
		})
	}
}

func TestMMSIToolGenerateMMSI(t *testing.T) {
	tool := NewMMSITool()

	tests := []struct {
		name        string
		params      map[string]interface{}
		expectError bool
	}{
		{
			name: "generate_single",
			params: map[string]interface{}{
				"operation":    "generate",
				"country-code": "US",
				"count":        1,
			},
			expectError: false,
		},
		{
			name: "generate_multiple",
			params: map[string]interface{}{
				"operation":    "generate",
				"country-code": "GB",
				"count":        3,
			},
			expectError: false,
		},
		{
			name: "generate_without_country",
			params: map[string]interface{}{
				"operation": "generate",
				"count":     2,
			},
			expectError: false,
		},
		{
			name: "invalid_country_code",
			params: map[string]interface{}{
				"operation":    "generate",
				"country-code": "XX",
				"count":        1,
			},
			expectError: true,
		},
		{
			name: "count_exceeds_limit",
			params: map[string]interface{}{
				"operation":    "generate",
				"country-code": "US",
				"count":        101,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tool.generateMMSI(tt.params)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result == nil {
				t.Errorf("Expected result but got nil")
				return
			}

			// Check if result is string (single) or slice (multiple)
			switch v := result.(type) {
			case string:
				if len(v) != 9 {
					t.Errorf("Expected 9-digit MMSI, got %d digits: %s", len(v), v)
				}
			case []string:
				if len(v) == 0 {
					t.Errorf("Expected non-empty slice")
				}
				for i, mmsi := range v {
					if len(mmsi) != 9 {
						t.Errorf("Expected 9-digit MMSI at index %d, got %d digits: %s", i, len(mmsi), mmsi)
					}
				}
			default:
				t.Errorf("Expected string or []string, got %T", result)
			}
		})
	}
}

func TestMMSIToolExecute(t *testing.T) {
	tool := NewMMSITool()

	tests := []struct {
		name        string
		params      map[string]interface{}
		expectError bool
	}{
		{
			name: "validate_operation",
			params: map[string]interface{}{
				"operation": "validate",
				"input":     "366123456",
			},
			expectError: false,
		},
		{
			name: "generate_operation",
			params: map[string]interface{}{
				"operation":    "generate",
				"country-code": "US",
				"count":        1,
			},
			expectError: false,
		},
		{
			name: "invalid_operation",
			params: map[string]interface{}{
				"operation": "invalid",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tool.Execute(tt.params)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result == nil {
				t.Errorf("Expected result but got nil")
			}
		})
	}
}

func TestMMSIToolEdgeCases(t *testing.T) {
	tool := NewMMSITool()

	tests := []struct {
		name        string
		params      map[string]interface{}
		expectError bool
	}{
		{
			name: "count_exceeds_limit",
			params: map[string]interface{}{
				"operation":    "generate",
				"country-code": "US",
				"count":        101,
			},
			expectError: true,
		},
		{
			name: "empty_input",
			params: map[string]interface{}{
				"operation": "validate",
				"input":     "",
			},
			expectError: true,
		},
		{
			name: "missing_operation_defaults_to_validate",
			params: map[string]interface{}{
				"input": "366123456",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tool.Execute(tt.params)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestMMSIToolSchema(t *testing.T) {
	tool := NewMMSITool()

	// Test input schema
	inputSchema := tool.GetInputSchema()
	if inputSchema == nil {
		t.Error("Input schema should not be nil")
	}

	// Test output schema
	outputSchema := tool.GetOutputSchema()
	if outputSchema == nil {
		t.Error("Output schema should not be nil")
	}
}

func TestMMSIToolTypeDetection(t *testing.T) {
	tool := NewMMSITool()

	tests := []struct {
		name     string
		mmsi     string
		expected string
	}{
		// Ship stations (non-US)
		{"regular_ship", "232123456", "Ship Station"},
		{"inmarsat_bcm", "232123000", "Ship Station (Inmarsat B/C/M)"},
		{"inmarsat_c", "232123450", "Ship Station (Inmarsat C)"},

		// Group stations
		{"group_ship", "036612345", "Group Ship Station"},
		{"coast_station", "003661234", "Coast Station"},
		{"group_coast", "003660000", "Group Coast Station"},

		// US Federal MMSIs
		{"us_coast_guard_group_ship", "036699999", "US Coast Guard Group Ship Station"},
		{"us_coast_guard_group_coast", "003669999", "US Coast Guard Group Coast Station"},
		{"us_federal", "366912345", "US Federal MMSI"},
		{"us_ship_international", "366123000", "US Ship Station (International/Inmarsat)"},
		{"us_ship_other", "366123450", "US Ship Station (Other)"},
		{"us_ship_regular", "366123456", "US Ship Station"},

		// Special devices
		{"sar_aircraft", "111123456", "SAR Aircraft"},
		{"handheld_vhf", "812345678", "Handheld VHF Transceiver"},
		{"sar_transponder", "970123456", "SAR Transponder (AIS-SART)"},
		{"man_overboard", "972123456", "Man Overboard Device"},
		{"epirb_ais", "974123456", "EPIRB-AIS"},
		{"craft_associated", "981234567", "Craft Associated with Parent Ship"},
		{"navigational_aid", "991234567", "Navigational Aid (AtoN)"},
		{"free_form", "912345678", "Free-form Device"},

		// Invalid
		{"too_short", "12345678", "Invalid"},
		{"too_long", "1234567890", "Invalid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tool.determineMMSIType(tt.mmsi)
			if result != tt.expected {
				t.Errorf("Expected type '%s', got '%s' for MMSI '%s'", tt.expected, result, tt.mmsi)
			}
		})
	}
}

func TestMMSIToolValidateWithType(t *testing.T) {
	tool := NewMMSITool()

	// Test regular ship station (non-US)
	result, err := tool.Execute(map[string]interface{}{
		"operation": "validate",
		"input":     "232123456",
	})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	validationResult, ok := result.(map[string]interface{})
	if !ok {
		t.Fatal("Expected map result")
	}
	if !validationResult["valid"].(bool) {
		t.Error("Expected valid MMSI")
	}
	if validationResult["type"] != "Ship Station" {
		t.Errorf("Expected type 'Ship Station', got '%s'", validationResult["type"])
	}

	// Test Inmarsat B/C/M ship station (non-US)
	result, err = tool.Execute(map[string]interface{}{
		"operation": "validate",
		"input":     "232123000",
	})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	validationResult, ok = result.(map[string]interface{})
	if !ok {
		t.Fatal("Expected map result")
	}
	if !validationResult["valid"].(bool) {
		t.Error("Expected valid MMSI")
	}
	if validationResult["type"] != "Ship Station (Inmarsat B/C/M)" {
		t.Errorf("Expected type 'Ship Station (Inmarsat B/C/M)', got '%s'", validationResult["type"])
	}

	// Test SAR Transponder
	result, err = tool.Execute(map[string]interface{}{
		"operation": "validate",
		"input":     "970123456",
	})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	validationResult, ok = result.(map[string]interface{})
	if !ok {
		t.Fatal("Expected map result")
	}
	if !validationResult["valid"].(bool) {
		t.Error("Expected valid MMSI")
	}
	if validationResult["type"] != "SAR Transponder (AIS-SART)" {
		t.Errorf("Expected type 'SAR Transponder (AIS-SART)', got '%s'", validationResult["type"])
	}
}
