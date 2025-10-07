// Package tools provides common interfaces and structures for mcpipboy tools
package tools

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
}

// ToolMetadata contains metadata about a tool
type ToolMetadata struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}
