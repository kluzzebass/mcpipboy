package tools

import (
	"testing"
)

func TestRandomTool(t *testing.T) {
	tool := NewRandomTool()

	if tool.Name() != "random" {
		t.Errorf("Expected name 'random', got: %s", tool.Name())
	}

	if tool.Description() == "" {
		t.Error("Expected non-empty description")
	}
}

func TestRandomToolExecute(t *testing.T) {
	tool := NewRandomTool()

	tests := []struct {
		name     string
		params   map[string]interface{}
		wantErr  bool
		validate func(t *testing.T, result interface{})
	}{
		{
			name:    "default integer generation",
			params:  map[string]interface{}{},
			wantErr: false,
			validate: func(t *testing.T, result interface{}) {
				if result == nil {
					t.Error("Expected non-nil result")
				}
				if _, ok := result.(int64); !ok {
					t.Errorf("Expected int64 result, got %T", result)
				}
			},
		},
		{
			name: "integer with range",
			params: map[string]interface{}{
				"type": "integer",
				"min":  float64(10),
				"max":  float64(20),
			},
			wantErr: false,
			validate: func(t *testing.T, result interface{}) {
				if result == nil {
					t.Error("Expected non-nil result")
				}
				if val, ok := result.(int64); ok {
					if val < 10 || val > 20 {
						t.Errorf("Expected value between 10-20, got: %d", val)
					}
				} else {
					t.Errorf("Expected int64 result, got %T", result)
				}
			},
		},
		{
			name: "multiple integers",
			params: map[string]interface{}{
				"type":  "integer",
				"count": float64(5),
				"min":   float64(1),
				"max":   float64(10),
			},
			wantErr: false,
			validate: func(t *testing.T, result interface{}) {
				if result == nil {
					t.Error("Expected non-nil result")
				}
				if arr, ok := result.([]int64); ok {
					if len(arr) != 5 {
						t.Errorf("Expected array of length 5, got: %d", len(arr))
					}
					for _, val := range arr {
						if val < 1 || val > 10 {
							t.Errorf("Expected value between 1-10, got: %d", val)
						}
					}
				} else {
					t.Errorf("Expected []int64 result, got %T", result)
				}
			},
		},
		{
			name: "float generation",
			params: map[string]interface{}{
				"type":      "float",
				"min":       float64(0),
				"max":       float64(1),
				"precision": float64(2),
			},
			wantErr: false,
			validate: func(t *testing.T, result interface{}) {
				if result == nil {
					t.Error("Expected non-nil result")
				}
				if val, ok := result.(float64); ok {
					if val < 0 || val > 1 {
						t.Errorf("Expected value between 0-1, got: %f", val)
					}
				} else {
					t.Errorf("Expected float64 result, got %T", result)
				}
			},
		},
		{
			name: "multiple floats",
			params: map[string]interface{}{
				"type":      "float",
				"count":     float64(3),
				"min":       float64(0),
				"max":       float64(10),
				"precision": float64(1),
			},
			wantErr: false,
			validate: func(t *testing.T, result interface{}) {
				if result == nil {
					t.Error("Expected non-nil result")
				}
				if arr, ok := result.([]float64); ok {
					if len(arr) != 3 {
						t.Errorf("Expected array of length 3, got: %d", len(arr))
					}
					for _, val := range arr {
						if val < 0 || val > 10 {
							t.Errorf("Expected value between 0-10, got: %f", val)
						}
					}
				} else {
					t.Errorf("Expected []float64 result, got %T", result)
				}
			},
		},
		{
			name: "boolean generation",
			params: map[string]interface{}{
				"type": "boolean",
			},
			wantErr: false,
			validate: func(t *testing.T, result interface{}) {
				if result == nil {
					t.Error("Expected non-nil result")
				}
				if _, ok := result.(bool); !ok {
					t.Errorf("Expected bool result, got %T", result)
				}
			},
		},
		{
			name: "multiple booleans",
			params: map[string]interface{}{
				"type":  "boolean",
				"count": float64(4),
			},
			wantErr: false,
			validate: func(t *testing.T, result interface{}) {
				if result == nil {
					t.Error("Expected non-nil result")
				}
				if arr, ok := result.([]bool); ok {
					if len(arr) != 4 {
						t.Errorf("Expected array of length 4, got: %d", len(arr))
					}
				} else {
					t.Errorf("Expected []bool result, got %T", result)
				}
			},
		},
		{
			name: "invalid type",
			params: map[string]interface{}{
				"type": "invalid",
			},
			wantErr: true,
		},
		{
			name: "invalid count too low",
			params: map[string]interface{}{
				"count": float64(0),
			},
			wantErr: true,
		},
		{
			name: "invalid count too high",
			params: map[string]interface{}{
				"count": float64(1001),
			},
			wantErr: true,
		},
		{
			name: "invalid range min >= max",
			params: map[string]interface{}{
				"type": "integer",
				"min":  float64(10),
				"max":  float64(5),
			},
			wantErr: true,
		},
		{
			name: "invalid precision too high",
			params: map[string]interface{}{
				"type":      "float",
				"precision": float64(11),
			},
			wantErr: true,
		},
		{
			name: "invalid precision negative",
			params: map[string]interface{}{
				"type":      "float",
				"precision": float64(-1),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tool.Execute(tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && tt.validate != nil {
				tt.validate(t, result)
			}
		})
	}
}

func TestRandomToolValidateParams(t *testing.T) {
	tool := NewRandomTool()

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
			name: "valid integer type",
			params: map[string]interface{}{
				"type": "integer",
			},
			wantErr: false,
		},
		{
			name: "valid float type",
			params: map[string]interface{}{
				"type": "float",
			},
			wantErr: false,
		},
		{
			name: "valid boolean type",
			params: map[string]interface{}{
				"type": "boolean",
			},
			wantErr: false,
		},
		{
			name: "valid count",
			params: map[string]interface{}{
				"count": float64(10),
			},
			wantErr: false,
		},
		{
			name: "valid range",
			params: map[string]interface{}{
				"type": "integer",
				"min":  float64(1),
				"max":  float64(100),
			},
			wantErr: false,
		},
		{
			name: "valid precision",
			params: map[string]interface{}{
				"type":      "float",
				"precision": float64(3),
			},
			wantErr: false,
		},
		{
			name: "invalid type",
			params: map[string]interface{}{
				"type": "invalid",
			},
			wantErr: true,
		},
		{
			name: "invalid count too low",
			params: map[string]interface{}{
				"count": float64(0),
			},
			wantErr: true,
		},
		{
			name: "invalid count too high",
			params: map[string]interface{}{
				"count": float64(1001),
			},
			wantErr: true,
		},
		{
			name: "invalid range min >= max",
			params: map[string]interface{}{
				"type": "integer",
				"min":  float64(10),
				"max":  float64(5),
			},
			wantErr: true,
		},
		{
			name: "invalid precision too high",
			params: map[string]interface{}{
				"type":      "float",
				"precision": float64(11),
			},
			wantErr: true,
		},
		{
			name: "invalid parameter types",
			params: map[string]interface{}{
				"type":  123,
				"count": "invalid",
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

func TestRandomToolSchemas(t *testing.T) {
	tool := NewRandomTool()

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
	if schemaType, ok := inputSchema["type"].(string); !ok || schemaType != "object" {
		t.Error("Input schema should have type 'object'")
	}

	if schemaType, ok := outputSchema["type"].(string); !ok || schemaType != "object" {
		t.Error("Output schema should have type 'object'")
	}
}

func TestRandomToolEdgeCases(t *testing.T) {
	tool := NewRandomTool()

	tests := []struct {
		name     string
		params   map[string]interface{}
		wantErr  bool
		validate func(t *testing.T, result interface{})
	}{
		{
			name: "single integer with count 1",
			params: map[string]interface{}{
				"type":  "integer",
				"count": float64(1),
				"min":   float64(5),
				"max":   float64(5),
			},
			wantErr: false,
			validate: func(t *testing.T, result interface{}) {
				if result == nil {
					t.Error("Expected non-nil result")
				}
				if val, ok := result.(int64); !ok {
					t.Errorf("Expected single int64 result, got %T", result)
				} else if val != 5 {
					t.Errorf("Expected value 5, got: %d", val)
				}
			},
		},
		{
			name: "single float with count 1",
			params: map[string]interface{}{
				"type":      "float",
				"count":     float64(1),
				"min":       float64(0.5),
				"max":       float64(0.5),
				"precision": float64(0),
			},
			wantErr: false,
			validate: func(t *testing.T, result interface{}) {
				if result == nil {
					t.Error("Expected non-nil result")
				}
				if val, ok := result.(float64); !ok {
					t.Errorf("Expected single float64 result, got %T", result)
				} else if val != 1.0 {
					t.Errorf("Expected value 1.0 (0.5 rounded to nearest integer), got: %f", val)
				}
			},
		},
		{
			name: "single boolean with count 1",
			params: map[string]interface{}{
				"type":  "boolean",
				"count": float64(1),
			},
			wantErr: false,
			validate: func(t *testing.T, result interface{}) {
				if result == nil {
					t.Error("Expected non-nil result")
				}
				if _, ok := result.(bool); !ok {
					t.Errorf("Expected single bool result, got %T", result)
				}
			},
		},
		{
			name: "maximum count",
			params: map[string]interface{}{
				"type":  "integer",
				"count": float64(1000),
				"min":   float64(1),
				"max":   float64(2),
			},
			wantErr: false,
			validate: func(t *testing.T, result interface{}) {
				if result == nil {
					t.Error("Expected non-nil result")
				}
				if arr, ok := result.([]int64); ok {
					if len(arr) != 1000 {
						t.Errorf("Expected array of length 1000, got: %d", len(arr))
					}
				} else {
					t.Errorf("Expected []int64 result, got %T", result)
				}
			},
		},
		{
			name: "zero precision",
			params: map[string]interface{}{
				"type":      "float",
				"precision": float64(0),
				"min":       float64(1.7),
				"max":       float64(1.7),
			},
			wantErr: false,
			validate: func(t *testing.T, result interface{}) {
				if result == nil {
					t.Error("Expected non-nil result")
				}
				if val, ok := result.(float64); ok {
					// Should be rounded to integer
					if val != 2.0 {
						t.Errorf("Expected value 2.0 (rounded), got: %f", val)
					}
				} else {
					t.Errorf("Expected float64 result, got %T", result)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tool.Execute(tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && tt.validate != nil {
				tt.validate(t, result)
			}
		})
	}
}
