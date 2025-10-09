// Package server provides the MCP server implementation
package server

import (
	"context"
	"io"
	"log"

	"github.com/kluzzebass/mcpipboy/internal/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Server represents the MCP server
type Server struct {
	server    *mcp.Server
	tools     map[string]tools.Tool
	debugMode bool
	logWriter io.Writer
}

// NewServer creates a new MCP server instance
func NewServer() *Server {
	return &Server{
		tools:     make(map[string]tools.Tool),
		debugMode: false,
		logWriter: nil,
	}
}

// SetDebugMode enables or disables debug logging
func (s *Server) SetDebugMode(enabled bool) {
	s.debugMode = enabled
}

// SetLogWriter sets the writer for debug logs
func (s *Server) SetLogWriter(w io.Writer) {
	s.logWriter = w
}

// RegisterTool registers a tool with the server
func (s *Server) RegisterTool(tool tools.Tool) {
	if tool == nil {
		return // Ignore nil tools
	}
	s.tools[tool.Name()] = tool
}

// Start starts the MCP server
func (s *Server) Start(ctx context.Context) error {
	if s.debugMode && s.logWriter != nil {
		log.SetOutput(s.logWriter)
		log.Println("Starting MCP server...")
	}

	// Create MCP server
	s.server = mcp.NewServer(&mcp.Implementation{Name: "mcpipboy"}, nil)

	// Register tools
	for _, tool := range s.tools {
		mcp.AddTool(s.server, &mcp.Tool{
			Name:        tool.Name(),
			Description: tool.Description(),
			InputSchema: tool.GetInputSchema(),
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
	// Only use LoggingTransport if debug mode is enabled
	var transport mcp.Transport
	if s.debugMode && s.logWriter != nil {
		transport = &mcp.LoggingTransport{Transport: &mcp.StdioTransport{}, Writer: s.logWriter}
	} else {
		transport = &mcp.StdioTransport{}
	}
	return s.server.Run(ctx, transport)
}

// Stop stops the MCP server
func (s *Server) Stop() error {
	if s.debugMode && s.logWriter != nil {
		log.Println("Stopping MCP server...")
	}
	// The server will stop when the context is cancelled
	return nil
}
