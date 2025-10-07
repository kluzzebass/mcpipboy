// Package tools provides the version tool implementation
package tools

import version "github.com/kluzzebass/mcpipboy"

// VersionTool implements the version functionality
type VersionTool struct{}

// NewVersionTool creates a new version tool instance
func NewVersionTool() *VersionTool {
	return &VersionTool{}
}

// Name returns the tool's name
func (v *VersionTool) Name() string {
	return "version"
}

// Description returns the tool's description
func (v *VersionTool) Description() string {
	return "Returns the current version of mcpipboy"
}

// Execute runs the version tool
func (v *VersionTool) Execute(params map[string]interface{}) (interface{}, error) {
	return version.Version, nil
}

// ValidateParams validates the input parameters
func (v *VersionTool) ValidateParams(params map[string]interface{}) error {
	// Version tool doesn't require any parameters
	return nil
}

// GetInputSchema returns the JSON schema for tool input parameters
func (v *VersionTool) GetInputSchema() map[string]interface{} {
	return CreateJSONSchema([]ParameterDefinition{})
}

// GetOutputSchema returns the JSON schema for tool output
func (v *VersionTool) GetOutputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"result": map[string]interface{}{
				"type":        "string",
				"description": "The current version of mcpipboy",
			},
		},
	}
}
