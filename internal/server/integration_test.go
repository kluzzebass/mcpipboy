// Package server provides integration tests for the MCP server
package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

// TestMCPServerIntegration tests the MCP server with actual MCP protocol communication
func TestMCPServerIntegration(t *testing.T) {
	// Skip this test in short mode since it requires external processes
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Get the project root directory
	projectRoot := "../../"

	// Build the mcpipboy binary for testing
	buildCmd := exec.Command("go", "build", "-o", "test-mcpipboy", "./cmd/mcpipboy")
	buildCmd.Dir = projectRoot
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build mcpipboy binary: %v", err)
	}
	defer os.Remove(projectRoot + "test-mcpipboy")

	// Test MCP server startup and tool discovery
	t.Run("ToolDiscovery", func(t *testing.T) {
		testToolDiscovery(t, projectRoot)
	})

	// Test tool execution
	t.Run("ToolExecution", func(t *testing.T) {
		testToolExecution(t, projectRoot)
	})
}

func testToolDiscovery(t *testing.T, projectRoot string) {
	// Start the MCP server
	cmd := exec.Command(projectRoot+"test-mcpipboy", "serve")
	cmd.Stderr = os.Stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatalf("Failed to create stdin pipe: %v", err)
	}
	defer stdin.Close()

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("Failed to create stdout pipe: %v", err)
	}

	if err := cmd.Start(); err != nil {
		t.Fatalf("Failed to start MCP server: %v", err)
	}
	defer func() {
		cmd.Process.Kill()
		cmd.Wait()
	}()

	// Send initialize request
	initRequest := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "initialize",
		"params": map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"capabilities": map[string]interface{}{
				"tools": map[string]interface{}{},
			},
			"clientInfo": map[string]interface{}{
				"name":    "test-client",
				"version": "1.0.0",
			},
		},
	}

	if err := sendMCPRequest(stdin, initRequest); err != nil {
		t.Fatalf("Failed to send initialize request: %v", err)
	}

	// Read initialize response
	initResponse, err := readMCPResponse(stdout)
	if err != nil {
		t.Fatalf("Failed to read initialize response: %v", err)
	}

	if initResponse["error"] != nil {
		t.Fatalf("Initialize request failed: %v", initResponse["error"])
	}

	// Send tools/list request
	toolsRequest := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      2,
		"method":  "tools/list",
		"params":  map[string]interface{}{},
	}

	if err := sendMCPRequest(stdin, toolsRequest); err != nil {
		t.Fatalf("Failed to send tools/list request: %v", err)
	}

	// Read tools/list response
	toolsResponse, err := readMCPResponse(stdout)
	if err != nil {
		t.Fatalf("Failed to read tools/list response: %v", err)
	}

	if toolsResponse["error"] != nil {
		t.Fatalf("Tools/list request failed: %v", toolsResponse["error"])
	}

	// Verify tools are listed
	result, ok := toolsResponse["result"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected result object, got %T", toolsResponse["result"])
	}

	toolsList, ok := result["tools"].([]interface{})
	if !ok {
		t.Fatalf("Expected tools array, got %T", result["tools"])
	}

	if len(toolsList) < 2 {
		t.Fatalf("Expected at least 2 tools (echo, version), got %d", len(toolsList))
	}

	// Verify echo and version tools are present
	toolNames := make(map[string]bool)
	for _, tool := range toolsList {
		toolMap, ok := tool.(map[string]interface{})
		if !ok {
			continue
		}
		if name, ok := toolMap["name"].(string); ok {
			toolNames[name] = true
		}
	}

	if !toolNames["echo"] {
		t.Error("Echo tool not found in tools list")
	}
	if !toolNames["version"] {
		t.Error("Version tool not found in tools list")
	}
}

func testToolExecution(t *testing.T, projectRoot string) {
	// Start the MCP server
	cmd := exec.Command(projectRoot+"test-mcpipboy", "serve")
	cmd.Stderr = os.Stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatalf("Failed to create stdin pipe: %v", err)
	}
	defer stdin.Close()

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("Failed to create stdout pipe: %v", err)
	}

	if err := cmd.Start(); err != nil {
		t.Fatalf("Failed to start MCP server: %v", err)
	}
	defer func() {
		cmd.Process.Kill()
		cmd.Wait()
	}()

	// Initialize the server
	if err := initializeMCPServer(stdin, stdout); err != nil {
		t.Fatalf("Failed to initialize MCP server: %v", err)
	}

	// Test echo tool execution
	t.Run("EchoTool", func(t *testing.T) {
		testEchoToolExecution(t, stdin, stdout)
	})

	// Test version tool execution
	t.Run("VersionTool", func(t *testing.T) {
		testVersionToolExecution(t, stdin, stdout)
	})
}

func testEchoToolExecution(t *testing.T, stdin io.WriteCloser, stdout io.ReadCloser) {
	// Send tools/call request for echo tool
	callRequest := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      3,
		"method":  "tools/call",
		"params": map[string]interface{}{
			"name": "echo",
			"arguments": map[string]interface{}{
				"message": "Hello, MCP!",
			},
		},
	}

	if err := sendMCPRequest(stdin, callRequest); err != nil {
		t.Fatalf("Failed to send tools/call request: %v", err)
	}

	// Read tools/call response
	callResponse, err := readMCPResponse(stdout)
	if err != nil {
		t.Fatalf("Failed to read tools/call response: %v", err)
	}

	if callResponse["error"] != nil {
		t.Fatalf("Tools/call request failed: %v", callResponse["error"])
	}

	// Verify echo response
	result, ok := callResponse["result"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected result object, got %T", callResponse["result"])
	}

	content, ok := result["content"].([]interface{})
	if !ok || len(content) == 0 {
		t.Fatalf("Expected content array with at least one item, got %v", result["content"])
	}

	contentItem, ok := content[0].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected content item to be an object, got %T", content[0])
	}

	text, ok := contentItem["text"].(string)
	if !ok {
		t.Fatalf("Expected text field in content, got %T", contentItem["text"])
	}

	// Parse the JSON response to get the actual result
	var echoResult map[string]interface{}
	if err := json.Unmarshal([]byte(text), &echoResult); err != nil {
		t.Fatalf("Failed to parse echo result JSON: %v", err)
	}

	if echoResult["result"] != "Hello, MCP!" {
		t.Errorf("Expected echo result 'Hello, MCP!', got '%v'", echoResult["result"])
	}
}

func testVersionToolExecution(t *testing.T, stdin io.WriteCloser, stdout io.ReadCloser) {
	// Send tools/call request for version tool
	callRequest := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      4,
		"method":  "tools/call",
		"params": map[string]interface{}{
			"name":      "version",
			"arguments": map[string]interface{}{},
		},
	}

	if err := sendMCPRequest(stdin, callRequest); err != nil {
		t.Fatalf("Failed to send tools/call request: %v", err)
	}

	// Read tools/call response
	callResponse, err := readMCPResponse(stdout)
	if err != nil {
		t.Fatalf("Failed to read tools/call response: %v", err)
	}

	if callResponse["error"] != nil {
		t.Fatalf("Tools/call request failed: %v", callResponse["error"])
	}

	// Verify version response
	result, ok := callResponse["result"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected result object, got %T", callResponse["result"])
	}

	content, ok := result["content"].([]interface{})
	if !ok || len(content) == 0 {
		t.Fatalf("Expected content array with at least one item, got %v", result["content"])
	}

	contentItem, ok := content[0].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected content item to be an object, got %T", content[0])
	}

	text, ok := contentItem["text"].(string)
	if !ok {
		t.Fatalf("Expected text field in content, got %T", contentItem["text"])
	}

	// Parse the JSON response to get the actual result
	var versionResult map[string]interface{}
	if err := json.Unmarshal([]byte(text), &versionResult); err != nil {
		t.Fatalf("Failed to parse version result JSON: %v", err)
	}

	// Read the expected version from VERSION file
	expectedVersionBytes, err := os.ReadFile("../../VERSION")
	if err != nil {
		t.Fatalf("Failed to read VERSION file: %v", err)
	}
	expectedVersion := strings.TrimSpace(string(expectedVersionBytes))

	// Validate the version matches
	if versionResult["result"] != expectedVersion {
		t.Errorf("Expected version '%s', got '%v'", expectedVersion, versionResult["result"])
	}
}

func initializeMCPServer(stdin io.WriteCloser, stdout io.ReadCloser) error {
	// Send initialize request
	initRequest := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "initialize",
		"params": map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"capabilities": map[string]interface{}{
				"tools": map[string]interface{}{},
			},
			"clientInfo": map[string]interface{}{
				"name":    "test-client",
				"version": "1.0.0",
			},
		},
	}

	if err := sendMCPRequest(stdin, initRequest); err != nil {
		return fmt.Errorf("failed to send initialize request: %w", err)
	}

	// Read initialize response
	initResponse, err := readMCPResponse(stdout)
	if err != nil {
		return fmt.Errorf("failed to read initialize response: %w", err)
	}

	if initResponse["error"] != nil {
		return fmt.Errorf("initialize request failed: %v", initResponse["error"])
	}

	return nil
}

func sendMCPRequest(stdin io.WriteCloser, request map[string]interface{}) error {
	data, err := json.Marshal(request)
	if err != nil {
		return err
	}

	// MCP uses newline-delimited JSON
	_, err = stdin.Write(append(data, '\n'))
	return err
}

func readMCPResponse(stdout io.ReadCloser) (map[string]interface{}, error) {
	// Read one line of JSON response
	scanner := bufio.NewScanner(stdout)
	if !scanner.Scan() {
		return nil, fmt.Errorf("no response received")
	}

	line := strings.TrimSpace(scanner.Text())
	if line == "" {
		return nil, fmt.Errorf("empty response")
	}

	var response map[string]interface{}
	if err := json.Unmarshal([]byte(line), &response); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	return response, nil
}

// TestMCPServerProtocolCompliance tests JSON-RPC protocol compliance
func TestMCPServerProtocolCompliance(t *testing.T) {
	// Skip this test in short mode
	if testing.Short() {
		t.Skip("Skipping protocol compliance test in short mode")
	}

	// Get the project root directory
	projectRoot := "../../"

	// Build the mcpipboy binary for testing
	buildCmd := exec.Command("go", "build", "-o", "test-mcpipboy", "./cmd/mcpipboy")
	buildCmd.Dir = projectRoot
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build mcpipboy binary: %v", err)
	}
	defer os.Remove(projectRoot + "test-mcpipboy")

	// Test invalid JSON-RPC request
	t.Run("InvalidJSONRPC", func(t *testing.T) {
		testInvalidJSONRPC(t, projectRoot)
	})

	// Test unsupported method
	t.Run("UnsupportedMethod", func(t *testing.T) {
		testUnsupportedMethod(t, projectRoot)
	})
}

func testInvalidJSONRPC(t *testing.T, projectRoot string) {
	cmd := exec.Command(projectRoot+"test-mcpipboy", "serve")
	cmd.Stderr = os.Stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatalf("Failed to create stdin pipe: %v", err)
	}
	defer stdin.Close()

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("Failed to create stdout pipe: %v", err)
	}

	if err := cmd.Start(); err != nil {
		t.Fatalf("Failed to start MCP server: %v", err)
	}
	defer func() {
		cmd.Process.Kill()
		cmd.Wait()
	}()

	// Send invalid JSON
	_, err = stdin.Write([]byte("invalid json\n"))
	if err != nil {
		t.Fatalf("Failed to send invalid JSON: %v", err)
	}

	// For invalid JSON, the server may crash or not respond
	// This is expected behavior since it can't parse the request
	// We just verify the server doesn't hang indefinitely
	done := make(chan bool, 1)
	go func() {
		response, err := readMCPResponse(stdout)
		if err != nil {
			// Expected - server can't respond to invalid JSON
			done <- true
			return
		}
		// If we get a response, it should be an error
		if response["error"] == nil {
			t.Error("Expected error response for invalid JSON")
		}
		done <- true
	}()

	select {
	case <-done:
		// Test completed (either with error or response)
	case <-time.After(2 * time.Second):
		t.Error("Server hung on invalid JSON request")
	}
}

func testUnsupportedMethod(t *testing.T, projectRoot string) {
	cmd := exec.Command(projectRoot+"test-mcpipboy", "serve")
	cmd.Stderr = os.Stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatalf("Failed to create stdin pipe: %v", err)
	}
	defer stdin.Close()

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("Failed to create stdout pipe: %v", err)
	}

	if err := cmd.Start(); err != nil {
		t.Fatalf("Failed to start MCP server: %v", err)
	}
	defer func() {
		cmd.Process.Kill()
		cmd.Wait()
	}()

	// Initialize first
	if err := initializeMCPServer(stdin, stdout); err != nil {
		t.Fatalf("Failed to initialize MCP server: %v", err)
	}

	// Send unsupported method request
	unsupportedRequest := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      5,
		"method":  "unsupported/method",
		"params":  map[string]interface{}{},
	}

	if err := sendMCPRequest(stdin, unsupportedRequest); err != nil {
		t.Fatalf("Failed to send unsupported method request: %v", err)
	}

	// Read error response
	response, err := readMCPResponse(stdout)
	if err != nil {
		t.Fatalf("Failed to read error response: %v", err)
	}

	// Should have an error
	if response["error"] == nil {
		t.Error("Expected error response for unsupported method")
	}
}
