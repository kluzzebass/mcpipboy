package tools

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// RandomTool implements comprehensive random number generation
type RandomTool struct{}

// NewRandomTool creates a new random tool instance
func NewRandomTool() *RandomTool {
	return &RandomTool{}
}

// Name returns the tool's name
func (r *RandomTool) Name() string {
	return "random"
}

// Description returns the tool's description
func (r *RandomTool) Description() string {
	return "Generate random numbers with various types and distributions"
}

// Execute runs the random tool
func (r *RandomTool) Execute(params map[string]interface{}) (interface{}, error) {
	// Parse parameters
	typeParam, _ := params["type"].(string)
	if typeParam == "" {
		typeParam = "integer"
	}

	count, countProvided := params["count"].(float64)

	// Validate count if provided
	if countProvided && (count < 1 || count > 1000) {
		return nil, fmt.Errorf("count must be between 1 and 1000")
	}

	// Set default if not provided
	if !countProvided {
		count = 1
	}

	// Generate random numbers based on type
	switch typeParam {
	case "integer":
		return r.generateIntegers(params, int(count))
	case "float":
		return r.generateFloats(params, int(count))
	case "boolean":
		return r.generateBooleans(int(count))
	default:
		return nil, fmt.Errorf("invalid type: %s, must be one of: integer, float, boolean", typeParam)
	}
}

// generateIntegers generates random integers
func (r *RandomTool) generateIntegers(params map[string]interface{}, count int) (interface{}, error) {
	min, _ := params["min"].(float64)
	max, _ := params["max"].(float64)

	// Default range if not specified
	if min == 0 && max == 0 {
		min = 0
		max = 100
	}

	// Validate range
	if min > max {
		return nil, fmt.Errorf("min must be less than or equal to max")
	}

	// Convert to int64 for generation
	minInt := int64(min)
	maxInt := int64(max)

	// Generate random integers
	rand.Seed(time.Now().UnixNano())
	var results []int64
	for i := 0; i < count; i++ {
		var value int64
		if minInt == maxInt {
			// If min equals max, return that exact value
			value = minInt
		} else {
			value = minInt + rand.Int63n(maxInt-minInt+1)
		}
		results = append(results, value)
	}

	// Return single value if count is 1, otherwise return array
	if count == 1 {
		return results[0], nil
	}
	return results, nil
}

// generateFloats generates random floats
func (r *RandomTool) generateFloats(params map[string]interface{}, count int) (interface{}, error) {
	min, _ := params["min"].(float64)
	max, _ := params["max"].(float64)
	precision, precisionProvided := params["precision"].(float64)

	// Default range if not specified
	if min == 0 && max == 0 {
		min = 0.0
		max = 1.0
	}

	// Default precision if not provided
	if !precisionProvided {
		precision = 2
	}

	// Validate range
	if min > max {
		return nil, fmt.Errorf("min must be less than or equal to max")
	}

	// Validate precision
	if precision < 0 || precision > 10 {
		return nil, fmt.Errorf("precision must be between 0 and 10")
	}

	// Generate random floats
	rand.Seed(time.Now().UnixNano())
	var results []float64
	for i := 0; i < count; i++ {
		var value float64
		if min == max {
			// If min equals max, use that value but still apply precision rounding
			value = min
		} else {
			value = min + rand.Float64()*(max-min)
		}
		// Round to specified precision
		if precision == 0.0 {
			// Round to nearest integer
			value = float64(int64(value + 0.5))
		} else {
			multiplier := 1.0
			for j := 0; j < int(precision); j++ {
				multiplier *= 10
			}
			value = float64(int64(value*multiplier+0.5)) / multiplier
		}
		results = append(results, value)
	}

	// Return single value if count is 1, otherwise return array
	if count == 1 {
		return results[0], nil
	}
	return results, nil
}

// generateBooleans generates random booleans
func (r *RandomTool) generateBooleans(count int) (interface{}, error) {
	// Generate random booleans
	rand.Seed(time.Now().UnixNano())
	var results []bool
	for i := 0; i < count; i++ {
		value := rand.Intn(2) == 1
		results = append(results, value)
	}

	// Return single value if count is 1, otherwise return array
	if count == 1 {
		return results[0], nil
	}
	return results, nil
}

// ValidateParams validates the input parameters
func (r *RandomTool) ValidateParams(params map[string]interface{}) error {
	// Validate type
	if typeParam, ok := params["type"]; ok {
		if typeStr, ok := typeParam.(string); ok {
			validTypes := []string{"integer", "float", "boolean"}
			if !contains(validTypes, typeStr) {
				return fmt.Errorf("invalid type: %s, must be one of: %s", typeStr, strings.Join(validTypes, ", "))
			}
		} else {
			return fmt.Errorf("type parameter must be a string")
		}
	}

	// Validate count
	if count, ok := params["count"]; ok {
		if countFloat, ok := count.(float64); ok {
			if countFloat < 1 || countFloat > 1000 {
				return fmt.Errorf("count must be between 1 and 1000")
			}
		} else {
			return fmt.Errorf("count parameter must be a number")
		}
	}

	// Validate min/max for integer and float types
	typeParam, _ := params["type"].(string)
	if typeParam == "" {
		typeParam = "integer"
	}

	if typeParam == "integer" || typeParam == "float" {
		min, minOk := params["min"]
		max, maxOk := params["max"]

		if minOk && maxOk {
			minFloat, minValid := min.(float64)
			maxFloat, maxValid := max.(float64)

			if !minValid || !maxValid {
				return fmt.Errorf("min and max parameters must be numbers")
			}

			if minFloat > maxFloat {
				return fmt.Errorf("min must be less than or equal to max")
			}
		}
	}

	// Validate precision for float type
	if typeParam == "float" {
		if precision, ok := params["precision"]; ok {
			if precisionFloat, ok := precision.(float64); ok {
				if precisionFloat < 0 || precisionFloat > 10 {
					return fmt.Errorf("precision must be between 0 and 10")
				}
			} else {
				return fmt.Errorf("precision parameter must be a number")
			}
		}
	}

	return nil
}

// GetInputSchema returns the JSON schema for tool input parameters
func (r *RandomTool) GetInputSchema() map[string]interface{} {
	return CreateJSONSchema([]ParameterDefinition{
		{
			Name:        "type",
			Type:        "string",
			Description: "Type of random number: integer, float, boolean",
			Required:    false,
		},
		{
			Name:        "count",
			Type:        "number",
			Description: "Number of random values to generate (1-1000)",
			Required:    false,
		},
		{
			Name:        "min",
			Type:        "number",
			Description: "Minimum value (for integer/float types)",
			Required:    false,
		},
		{
			Name:        "max",
			Type:        "number",
			Description: "Maximum value (for integer/float types)",
			Required:    false,
		},
		{
			Name:        "precision",
			Type:        "number",
			Description: "Decimal places for float values (0-10)",
			Required:    false,
		},
	})
}

// GetOutputSchema returns the JSON schema for tool output
func (r *RandomTool) GetOutputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"result": map[string]interface{}{
				"type":        "array",
				"description": "Array of random values (or single value if count=1)",
				"items": map[string]interface{}{
					"oneOf": []map[string]interface{}{
						{"type": "integer"},
						{"type": "number"},
						{"type": "boolean"},
					},
				},
			},
		},
	}
}
