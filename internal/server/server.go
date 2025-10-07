// Package server provides the MCP server implementation
package server

import (
	"context"
	"log"
	"os"

	"github.com/kluzzebass/mcpipboy/internal/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Server represents the MCP server
type Server struct {
	server *mcp.Server
	tools  map[string]tools.Tool
}

// NewServer creates a new MCP server instance
func NewServer() *Server {
	return &Server{
		tools: make(map[string]tools.Tool),
	}
}

// RegisterTool registers a tool with the server
func (s *Server) RegisterTool(tool tools.Tool) {
	s.tools[tool.Name()] = tool
}

// Start starts the MCP server
func (s *Server) Start(ctx context.Context) error {
	log.Println("Starting MCP server...")

	// Create MCP server
	s.server = mcp.NewServer(&mcp.Implementation{Name: "mcpipboy"}, nil)

	// Register tools
	for _, tool := range s.tools {
		mcp.AddTool(s.server, &mcp.Tool{
			Name:        tool.Name(),
			Description: tool.Description(),
		}, func(ctx context.Context, request *mcp.CallToolRequest, input map[string]interface{}) (*mcp.CallToolResult, map[string]interface{}, error) {
			// Execute the tool
			result, err := tool.Execute(input)
			if err != nil {
				return nil, nil, err
			}

			// Convert result to map
			resultMap := map[string]interface{}{
				"result": result,
			}
			return nil, resultMap, nil
		})
	}

	// Start the server with stdio transport
	transport := &mcp.LoggingTransport{Transport: &mcp.StdioTransport{}, Writer: os.Stderr}
	return s.server.Run(ctx, transport)
}

// Stop stops the MCP server
func (s *Server) Stop() error {
	log.Println("Stopping MCP server...")
	// The server will stop when the context is cancelled
	return nil
}
