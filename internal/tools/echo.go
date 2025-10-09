// Package tools provides the echo tool implementation
package tools

import (
	"fmt"
)

// EchoTool implements the echo functionality
type EchoTool struct{}

// NewEchoTool creates a new echo tool instance
func NewEchoTool() *EchoTool {
	return &EchoTool{}
}

// Name returns the tool's name
func (e *EchoTool) Name() string {
	return "echo"
}

// Description returns the tool's description
func (e *EchoTool) Description() string {
	return "Echoes back the input message"
}

// Execute runs the echo tool
func (e *EchoTool) Execute(params map[string]interface{}) (interface{}, error) {
	message, ok := params["message"].(string)
	if !ok {
		return nil, fmt.Errorf("message parameter is required and must be a string")
	}
	return message, nil
}

// ValidateParams validates the input parameters
func (e *EchoTool) ValidateParams(params map[string]interface{}) error {
	if message, ok := params["message"]; !ok {
		return fmt.Errorf("message parameter is required")
	} else if _, ok := message.(string); !ok {
		return fmt.Errorf("message parameter must be a string")
	}
	return nil
}

// GetInputSchema returns the JSON schema for tool input parameters
func (e *EchoTool) GetInputSchema() map[string]interface{} {
	return CreateJSONSchema([]ParameterDefinition{
		{
			Name:        "message",
			Type:        "string",
			Description: "The message to echo back",
			Required:    true,
		},
	})
}

// GetOutputSchema returns the JSON schema for tool output
func (e *EchoTool) GetOutputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"result": map[string]interface{}{
				"type":        "string",
				"description": "The echoed message",
			},
		},
	}
}

// GetResources returns the list of resources this tool provides
func (e *EchoTool) GetResources() []Resource {
	return []Resource{}
}

// ReadResource reads a specific resource by URI
func (e *EchoTool) ReadResource(uri string) (string, error) {
	return "", fmt.Errorf("no resources available for echo tool")
}
