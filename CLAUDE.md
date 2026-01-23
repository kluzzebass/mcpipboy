# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

mcpipboy is an MCP (Model Context Protocol) server that provides AI agents with tools for data generation and validation. It operates in two modes:
- **MCP Mode** - JSON-RPC 2.0 server for AI agent integration (Cursor, Claude Desktop, etc.)
- **CLI Mode** - Direct command-line access to all tools for testing and automation

## Common Commands

```bash
# Build
just build              # Build for current platform (output: dist/mcpipboy)
just build-all          # Build static binaries for all platforms

# Test
just test               # Run all tests
just test-coverage      # Run tests with coverage report
go test -v ./internal/tools/...   # Test specific package

# Code quality
just lint               # Run golangci-lint (if installed)
just fmt                # Format code
just vet                # Run go vet
just check              # Run fmt + vet + test

# Development
just dev                # Build + test cycle

# Release (automated via GitHub Actions)
git tag v1.0.0 && git push origin v1.0.0
```

## Versioning and Releases

- **Version source**: Git tags (no VERSION file)
- **Version format**: Semantic versioning with `v` prefix (e.g., `v1.0.0`)
- **Release process**: Push a tag -> GitHub Actions builds all platforms -> creates GitHub release -> updates Homebrew tap
- **Local builds**: Version extracted via `git describe --tags`
- **Library version**: Hardcoded in `version.go` (update when releasing library changes)

## Architecture

### Directory Structure

- `cmd/mcpipboy/` - CLI implementation (Cobra commands, one file per tool)
- `internal/server/` - MCP server implementation (JSON-RPC 2.0 handling)
- `internal/tools/` - Tool implementations and registry
- `plans/` - Development plans with sequential numbering

### Key Components

1. **Tool Interface** (`internal/tools/interfaces.go`) - All tools implement this interface:
   - `Name()`, `Description()` - Tool metadata
   - `Execute(params)` - Run the tool
   - `ValidateParams(params)` - Parameter validation
   - `GetInputSchema()`, `GetOutputSchema()` - JSON schemas for MCP

2. **Tool Registry** (`internal/tools/interfaces.go`) - Manages tool registration and discovery

3. **MCP Server** (`internal/server/server.go`) - Handles JSON-RPC 2.0 protocol, tool listing, and execution

4. **CLI Commands** (`cmd/mcpipboy/`) - Each tool has a corresponding CLI command file

### Adding a New Tool

1. Create `internal/tools/<toolname>.go` implementing the `Tool` interface
2. Create `internal/tools/<toolname>_test.go` with comprehensive tests
3. Create `cmd/mcpipboy/<toolname>.go` for CLI access
4. Create `cmd/mcpipboy/<toolname>_test.go` for CLI tests
5. Register the tool in the MCP server initialization

## Cobra Command Conventions

- Include examples for non-trivial commands using the `Examples` property
- Don't put examples in the `Long` property
- Name command functions after the full command path (e.g., `validate schema` -> `validateSchema()`)

## Testing Patterns

- Use table-driven tests for multiple scenarios
- Test both success and failure paths
- Tools have two test files: one in `internal/tools/` (unit) and one in `cmd/mcpipboy/` (CLI integration)
- MCP server integration tests are in `internal/server/integration_test.go`

## Plan Management

Plans are stored in `plans/` with sequential numbering (`001-descriptive-name.md`). See `.cursor/rules/` for detailed plan format requirements.
