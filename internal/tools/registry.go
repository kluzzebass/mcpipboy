// Package tools provides the tool registry system
package tools

import (
	"fmt"
	"sync"
)

// Registry manages the collection of available tools
type Registry struct {
	tools map[string]Tool
	mutex sync.RWMutex
}

// NewRegistry creates a new tool registry
func NewRegistry() *Registry {
	return &Registry{
		tools: make(map[string]Tool),
	}
}

// Register adds a tool to the registry
func (r *Registry) Register(tool Tool) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	name := tool.Name()
	if name == "" {
		return fmt.Errorf("tool name cannot be empty")
	}

	r.tools[name] = tool
	return nil
}

// Get retrieves a tool by name
func (r *Registry) Get(name string) (Tool, bool) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	tool, exists := r.tools[name]
	return tool, exists
}

// List returns all registered tool names
func (r *Registry) List() []string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	names := make([]string, 0, len(r.tools))
	for name := range r.tools {
		names = append(names, name)
	}
	return names
}

// GetMetadata returns metadata for all tools
func (r *Registry) GetMetadata() []ToolMetadata {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	metadata := make([]ToolMetadata, 0, len(r.tools))
	for _, tool := range r.tools {
		metadata = append(metadata, ToolMetadata{
			Name:        tool.Name(),
			Description: tool.Description(),
		})
	}
	return metadata
}
