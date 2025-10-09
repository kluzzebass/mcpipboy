# mcpipboy

**Flawlessly generated, rigorously validated - dependable data for the digital wasteland.**

Because even wasteland AIs deserve clean timestamps and proper UUIDs.

## Overview

mcpipboy is a Model Context Protocol (MCP) server that provides AI agents with essential tools for data generation and validation. Built with Go and designed for reliability, it offers a comprehensive suite of utilities that AI agents commonly need but often struggle to implement correctly.

## Features

### Core Tools
- **Echo Tool**: Simple message echoing for testing and validation
- **Version Tool**: Returns the current version of mcpipboy
- **Time Tool**: Flexible time operations (current time, parsing, formatting, timezone conversion)
- **Random Tool**: Generate random data (integers, strings, UUIDs, passwords)
- **UUID Tool**: Generate and validate UUIDs (v1, v4, v5, v7)

### Validation & Generation Tools
- **Credit Card Tool**: Generate and validate credit card numbers with Luhn algorithm
- **ISBN Tool**: Generate and validate ISBN-10 and ISBN-13 numbers
- **EAN-13 Tool**: Generate and validate EAN-13 barcodes
- **IBAN Tool**: Generate and validate International Bank Account Numbers with MOD-97 checksum
- **IMO Tool**: Generate and validate International Maritime Organization numbers
- **MMSI Tool**: Generate and validate Maritime Mobile Service Identity numbers

### Infrastructure
- **MCP Protocol Compliance**: Full JSON-RPC 2.0 compliance with proper error handling
- **Command Line Interface**: Direct tool invocation via CLI for testing and automation
- **Selective Tool Management**: Enable/disable specific tools via command line flags
- **Static Binary Builds**: Self-contained executables for easy deployment
- **Comprehensive Testing**: Full test coverage with unit and integration tests

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

# Time operations
mcpipboy time --type current
mcpipboy time --type parse --input "2024-01-15T10:30:00Z"
mcpipboy time --type format --input "2024-01-15T10:30:00Z" --format "2006-01-02 15:04:05"

# Random data generation
mcpipboy random --type integer --min 1 --max 100
mcpipboy random --type string --length 10
mcpipboy random --type password --length 16

# UUID operations
mcpipboy uuid --operation generate --version v4
mcpipboy uuid --operation validate --input "550e8400-e29b-41d4-a716-446655440000"

# Credit card operations
mcpipboy creditcard --operation validate --input "4532015112830366"
mcpipboy creditcard --operation generate --count 3

# ISBN operations
mcpipboy isbn --operation validate --input "978-0-306-40615-7"
mcpipboy isbn --operation generate --type isbn13 --count 2

# EAN-13 operations
mcpipboy ean13 --operation validate --input "1234567890123"
mcpipboy ean13 --operation generate --count 5

# IBAN operations
mcpipboy iban --operation validate --input "GB82WEST12345698765432"
mcpipboy iban --operation generate --country GB --count 3

# IMO operations
mcpipboy imo --operation validate --input "9176181"
mcpipboy imo --operation generate --count 5

# MMSI operations
mcpipboy mmsi --operation validate --input "123456789"
mcpipboy mmsi --operation generate --country US --count 3
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

### Core Tools
- **echo**: Echoes back the input message (parameter: `message`)
- **version**: Returns the current version of mcpipboy (no parameters)

### Time & Date Tools
- **time**: Flexible time operations
  - `current`: Get current time in various formats
  - `parse`: Parse time strings with automatic format detection
  - `format`: Format time strings to specific layouts
  - `convert`: Convert between timezones
  - `offset`: Apply time offsets (add/subtract duration)

### Random Data Tools
- **random**: Generate random data
  - `integer`: Random integers within a range
  - `string`: Random strings with specified length and character set
  - `uuid`: Generate random UUIDs
  - `password`: Generate secure passwords with customizable rules

### UUID Tools
- **uuid**: UUID generation and validation
  - `generate`: Generate UUIDs (v1, v4, v5, v7)
  - `validate`: Validate UUID format and version
  - `parse`: Parse UUID strings and extract components

### Validation & Generation Tools
- **creditcard**: Credit card number operations
  - `validate`: Validate credit card numbers using Luhn algorithm
  - `generate`: Generate valid credit card numbers
  - `type`: Detect card type (Visa, MasterCard, etc.)

- **isbn**: ISBN operations
  - `validate`: Validate ISBN-10 and ISBN-13 numbers
  - `generate`: Generate valid ISBN numbers
  - `convert`: Convert between ISBN-10 and ISBN-13

- **ean13**: EAN-13 barcode operations
  - `validate`: Validate EAN-13 barcodes with checksum
  - `generate`: Generate valid EAN-13 barcodes
  - `decode`: Decode EAN-13 country and manufacturer info

- **iban**: International Bank Account Number operations
  - `validate`: Validate IBANs with MOD-97 checksum
  - `generate`: Generate valid IBANs for specified countries
  - `decode`: Decode IBAN country and bank information

- **imo**: International Maritime Organization number operations
  - `validate`: Validate IMO numbers with checksum
  - `generate`: Generate valid IMO numbers
  - `decode`: Decode IMO number components

- **mmsi**: Maritime Mobile Service Identity operations
  - `validate`: Validate MMSI numbers
  - `generate`: Generate MMSI numbers for specified countries
  - `decode`: Decode MMSI country and vessel type information

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

### Completed
- [x] Core tools (echo, version)
- [x] Time tool with flexible operations (current, parse, format, convert, offset)
- [x] Random data generation (integers, strings, UUIDs, passwords)
- [x] UUID generation and validation (v1, v4, v5, v7)
- [x] Credit card validation and generation with Luhn algorithm
- [x] ISBN validation and generation (ISBN-10, ISBN-13)
- [x] EAN-13 barcode validation and generation
- [x] IBAN validation and generation with MOD-97 checksum
- [x] IMO number validation and generation
- [x] MMSI number validation and generation
- [x] Comprehensive CLI interface for all tools
- [x] MCP server integration
- [x] Full test coverage

### Future Enhancements
- [ ] Enhanced date/time utilities and formatting
- [ ] Random data generation (names, addresses, realistic test data)
- [ ] Data validation and sanitization tools
- [ ] IP calculation tool (IPv4/IPv6 subnet calculations, CIDR operations, network analysis)
- [ ] Crypto tool (password hashing, verification, hash algorithms)
- [ ] Text similarity tool (Levenshtein distance, fuzzy matching, string comparison)
- [ ] Unit converter tool (length, weight, temperature, currency, conversions)
- [ ] Additional barcode formats (UPC, Code 128, etc.)
- [ ] More international standards (SWIFT codes, etc.)
- [ ] Performance optimizations for large-scale generation
- [ ] Enhanced MCP protocol features
- [ ] Plugin architecture for custom tools

## Support

For issues, feature requests, or questions, please open an issue on the [GitHub repository](https://github.com/kluzzebass/mcpipboy/issues).
