// Package tools provides common interfaces and structures for mcpipboy tools
package tools

import (
	"encoding/json"
	"fmt"
)

// Tool defines the common interface that all tools must implement
type Tool interface {
	// Name returns the tool's name
	Name() string

	// Description returns the tool's description
	Description() string

	// Execute runs the tool with the given parameters
	Execute(params map[string]interface{}) (interface{}, error)

	// ValidateParams validates the input parameters
	ValidateParams(params map[string]interface{}) error

	// GetInputSchema returns the JSON schema for tool input parameters
	GetInputSchema() map[string]interface{}

	// GetOutputSchema returns the JSON schema for tool output
	GetOutputSchema() map[string]interface{}

	// GetResources returns the list of resources this tool provides
	GetResources() []Resource

	// ReadResource reads a specific resource by URI
	ReadResource(uri string) (string, error)
}

// Resource represents a resource that a tool can provide
type Resource struct {
	Name     string `json:"name"`
	URI      string `json:"uri"`
	MIMEType string `json:"mimeType"`
}

// ToolMetadata contains metadata about a tool
type ToolMetadata struct {
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	InputSchema  map[string]interface{} `json:"inputSchema"`
	OutputSchema map[string]interface{} `json:"outputSchema"`
}

// ParameterDefinition defines a tool parameter
type ParameterDefinition struct {
	Name        string      `json:"name"`
	Type        string      `json:"type"`
	Description string      `json:"description"`
	Required    bool        `json:"required"`
	Default     interface{} `json:"default,omitempty"`
	Enum        []string    `json:"enum,omitempty"`
}

// ToolRegistry manages tool registration and discovery
type ToolRegistry struct {
	tools map[string]Tool
}

// NewToolRegistry creates a new tool registry
func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{
		tools: make(map[string]Tool),
	}
}

// RegisterTool registers a tool with the registry
func (r *ToolRegistry) RegisterTool(tool Tool) error {
	if tool == nil {
		return fmt.Errorf("tool cannot be nil")
	}

	name := tool.Name()
	if name == "" {
		return fmt.Errorf("tool name cannot be empty")
	}

	if _, exists := r.tools[name]; exists {
		return fmt.Errorf("tool %s is already registered", name)
	}

	r.tools[name] = tool
	return nil
}

// GetTool retrieves a tool by name
func (r *ToolRegistry) GetTool(name string) (Tool, bool) {
	tool, exists := r.tools[name]
	return tool, exists
}

// ListTools returns a list of all registered tool names
func (r *ToolRegistry) ListTools() []string {
	names := make([]string, 0, len(r.tools))
	for name := range r.tools {
		names = append(names, name)
	}
	return names
}

// GetToolMetadata returns metadata for a specific tool
func (r *ToolRegistry) GetToolMetadata(name string) (*ToolMetadata, error) {
	tool, exists := r.GetTool(name)
	if !exists {
		return nil, fmt.Errorf("tool %s not found", name)
	}

	return &ToolMetadata{
		Name:         tool.Name(),
		Description:  tool.Description(),
		InputSchema:  tool.GetInputSchema(),
		OutputSchema: tool.GetOutputSchema(),
	}, nil
}

// GetAllToolMetadata returns metadata for all registered tools
func (r *ToolRegistry) GetAllToolMetadata() ([]ToolMetadata, error) {
	metadata := make([]ToolMetadata, 0, len(r.tools))
	for _, tool := range r.tools {
		metadata = append(metadata, ToolMetadata{
			Name:         tool.Name(),
			Description:  tool.Description(),
			InputSchema:  tool.GetInputSchema(),
			OutputSchema: tool.GetOutputSchema(),
		})
	}
	return metadata, nil
}

// ValidateToolParameters validates parameters for a specific tool
func (r *ToolRegistry) ValidateToolParameters(toolName string, params map[string]interface{}) error {
	tool, exists := r.GetTool(toolName)
	if !exists {
		return fmt.Errorf("tool %s not found", toolName)
	}

	return tool.ValidateParams(params)
}

// ExecuteTool executes a tool with the given parameters
func (r *ToolRegistry) ExecuteTool(toolName string, params map[string]interface{}) (interface{}, error) {
	tool, exists := r.GetTool(toolName)
	if !exists {
		return nil, fmt.Errorf("tool %s not found", toolName)
	}

	// Validate parameters first
	if err := tool.ValidateParams(params); err != nil {
		return nil, fmt.Errorf("parameter validation failed: %w", err)
	}

	// Execute the tool
	return tool.Execute(params)
}

// CreateJSONSchema creates a JSON schema for tool parameters
func CreateJSONSchema(parameters []ParameterDefinition) map[string]interface{} {
	properties := make(map[string]interface{})
	required := make([]string, 0)

	for _, param := range parameters {
		property := map[string]interface{}{
			"type":        param.Type,
			"description": param.Description,
		}

		if param.Default != nil {
			property["default"] = param.Default
		}

		if len(param.Enum) > 0 {
			property["enum"] = param.Enum
		}

		properties[param.Name] = property

		if param.Required {
			required = append(required, param.Name)
		}
	}

	schema := map[string]interface{}{
		"type":       "object",
		"properties": properties,
	}

	if len(required) > 0 {
		schema["required"] = required
	}

	return schema
}

// ValidateParameter validates a single parameter against its definition
func ValidateParameter(value interface{}, param ParameterDefinition) error {
	// Check if required parameter is missing
	if param.Required && value == nil {
		return fmt.Errorf("required parameter %s is missing", param.Name)
	}

	// If value is nil and not required, it's valid
	if value == nil {
		return nil
	}

	// Type validation
	switch param.Type {
	case "string":
		if _, ok := value.(string); !ok {
			return fmt.Errorf("parameter %s must be a string", param.Name)
		}
	case "number":
		switch value.(type) {
		case int, int8, int16, int32, int64, float32, float64:
			// Valid number types
		default:
			return fmt.Errorf("parameter %s must be a number", param.Name)
		}
	case "boolean":
		if _, ok := value.(bool); !ok {
			return fmt.Errorf("parameter %s must be a boolean", param.Name)
		}
	case "array":
		if _, ok := value.([]interface{}); !ok {
			return fmt.Errorf("parameter %s must be an array", param.Name)
		}
	case "object":
		if _, ok := value.(map[string]interface{}); !ok {
			return fmt.Errorf("parameter %s must be an object", param.Name)
		}
	}

	// Enum validation
	if len(param.Enum) > 0 {
		if strValue, ok := value.(string); ok {
			valid := false
			for _, enumValue := range param.Enum {
				if strValue == enumValue {
					valid = true
					break
				}
			}
			if !valid {
				return fmt.Errorf("parameter %s must be one of: %v", param.Name, param.Enum)
			}
		}
	}

	return nil
}

// ValidateParameters validates multiple parameters against their definitions
func ValidateParameters(params map[string]interface{}, paramDefs []ParameterDefinition) error {
	for _, paramDef := range paramDefs {
		value := params[paramDef.Name]
		if err := ValidateParameter(value, paramDef); err != nil {
			return err
		}
	}

	// Check for unknown parameters
	for paramName := range params {
		found := false
		for _, paramDef := range paramDefs {
			if paramDef.Name == paramName {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("unknown parameter: %s", paramName)
		}
	}

	return nil
}

// ToolResult represents the result of a tool execution
type ToolResult struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// ToJSON converts the tool result to JSON
func (tr *ToolResult) ToJSON() ([]byte, error) {
	return json.Marshal(tr)
}

// NewSuccessResult creates a successful tool result
func NewSuccessResult(data interface{}) *ToolResult {
	return &ToolResult{
		Success: true,
		Data:    data,
	}
}

// NewErrorResult creates an error tool result
func NewErrorResult(err error) *ToolResult {
	return &ToolResult{
		Success: false,
		Error:   err.Error(),
	}
}
