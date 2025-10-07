# mcpipboy

**Flawlessly generated, rigorously validated - dependable data for the digital wasteland.**

Because even wasteland AIs deserve clean timestamps and proper UUIDs.

## Overview

mcpipboy is a Model Context Protocol (MCP) server that provides AI agents with essential tools for data generation and validation. Built with Go and designed for reliability, it offers a comprehensive suite of utilities that AI agents commonly need but often struggle to implement correctly.

## Features

- **Echo Tool**: Simple message echoing for testing and validation
- **Version Tool**: Returns the current version of mcpipboy
- **MCP Protocol Compliance**: Full JSON-RPC 2.0 compliance with proper error handling
- **Command Line Interface**: Direct tool invocation via CLI for testing and automation
- **Selective Tool Management**: Enable/disable specific tools via command line flags
- **Static Binary Builds**: Self-contained executables for easy deployment

## Installation

### From Releases

Download the latest release for your platform from the [Releases page](https://github.com/kluzzebass/mcpipboy/releases).

### From Source

```bash
git clone https://github.com/kluzzebass/mcpipboy.git
cd mcpipboy
go build -o mcpipboy ./cmd/mcpipboy
```

### Using Go Install

```bash
go install github.com/kluzzebass/mcpipboy/cmd/mcpipboy@latest
```

## Usage

### Command Line Interface

mcpipboy provides both MCP server functionality and direct CLI access to tools:

```bash
# Start the MCP server
mcpipboy serve

# Start with specific tools enabled
mcpipboy serve --enable echo,version

# Start with specific tools disabled
mcpipboy serve --disable echo

# Use tools directly via CLI
mcpipboy echo "Hello, wasteland!"
mcpipboy version
```

### MCP Client Integration

#### Cursor IDE

Add mcpipboy to your Cursor MCP configuration:

1. Open Cursor settings
2. Navigate to MCP settings
3. Add the following configuration:

```json
{
  "mcpServers": {
    "mcpipboy": {
      "command": "mcpipboy",
      "args": ["serve"],
      "env": {}
    }
  }
}
```

#### Claude Desktop

Add mcpipboy to your Claude Desktop MCP configuration:

1. Open Claude Desktop settings
2. Navigate to MCP settings
3. Add the following configuration:

```json
{
  "mcpServers": {
    "mcpipboy": {
      "command": "mcpipboy",
      "args": ["serve"],
      "env": {}
    }
  }
}
```

#### Continue.dev

Add mcpipboy to your Continue configuration:

1. Open Continue settings
2. Navigate to MCP settings
3. Add the following configuration:

```json
{
  "mcpServers": {
    "mcpipboy": {
      "command": "mcpipboy",
      "args": ["serve"],
      "env": {}
    }
  }
}
```

#### Custom MCP Client

For custom MCP clients, mcpipboy communicates via stdin/stdout using JSON-RPC 2.0:

```bash
# Start the server
mcpipboy serve

# Send MCP requests via stdin
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{"tools":{}},"clientInfo":{"name":"test-client","version":"1.0.0"}}}' | mcpipboy serve
```

## Available Tools

- **echo**: Echoes back the input message (parameter: `message`)
- **version**: Returns the current version of mcpipboy (no parameters)

## Development

### Prerequisites

- Go 1.21 or later
- Make or Just (for build automation)

### Building

```bash
# Build the binary
just build

# Build for release (static binary)
just build-release

# Run tests
just test

# Run tests with coverage
just test-coverage
```

### Testing

```bash
# Run all tests
just test

# Run tests with coverage
just test-coverage

# Run integration tests specifically
go test -v -run "TestMCPServerIntegration" ./internal/server/...
```

### Version Management

```bash
# Get current version
just get-version

# Bump patch version (0.1.0 -> 0.1.1)
just bump-patch

# Bump minor version (0.1.0 -> 0.2.0)
just bump-minor

# Bump major version (0.1.0 -> 1.0.0)
just bump-major
```

## Architecture

mcpipboy is built with a modular architecture:

- **`cmd/mcpipboy/`**: Command-line interface using Cobra
- **`internal/server/`**: MCP server implementation
- **`internal/tools/`**: Tool implementations and registry
- **`version.go`**: Version management with embedded VERSION file

### Tool Interface

All tools implement the `tools.Tool` interface:

```go
type Tool interface {
    Name() string
    Description() string
    Execute(params map[string]interface{}) (interface{}, error)
    ValidateParams(params map[string]interface{}) error
    GetInputSchema() map[string]interface{}
    GetOutputSchema() map[string]interface{}
}
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Roadmap

- [ ] UUID generation tools (v1, v4, v5, v7)
- [ ] Checksum validation tools (IMO, MMSI, credit card numbers, ISBN)
- [ ] Date/time utilities
- [ ] Random data generation
- [ ] Data validation and sanitization tools

## Support

For issues, feature requests, or questions, please open an issue on the [GitHub repository](https://github.com/kluzzebass/mcpipboy/issues).
