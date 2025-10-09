# mcpipboy MVP Development Plan

**Date:** 2025-10-07

## Purpose
Create a stdin/stdout MCP (Model Context Protocol) server that provides agentic AIs with essential tools for common tasks they struggle with, including UUID generation, checksummed identifier verification/generation (IMO, MMSI, credit card numbers, ISBN, etc.), and other utility functions.

## High Level Phases
1. **Project Foundation** - Set up build system, dependencies, and basic project structure
2. **MCP Server Core** - Implement basic MCP server with echo tool
3. **Build & Release System** - Create justfile for build automation and semantic versioning
4. **Testing & Validation** - Ensure MVP works correctly with MCP clients

Progress legend: [x] Completed, [>] In progress, [ ] Pending

---

### [x] 0) Planning document and alignment
Establish the development plan and ensure all requirements are captured correctly.

1. **Requirements Analysis**
   - [x] Define core purpose: MCP server for AI agent utility tools
   - [x] Identify key libraries: github.com/modelcontextprotocol/go-sdk, cobra
   - [x] Establish MVP scope: echo tool as proof of concept
   - [x] Plan build system: justfile for automation

2. **Architecture Decisions**
   - [x] Use stdin/stdout for MCP communication
   - [x] Cobra for CLI argument parsing
   - [x] Go modules for dependency management
   - [x] Semantic versioning for releases
   - [x] Individual tool commands for CLI usage (e.g., `mcpipboy echo`)
   - [x] Tool enable/disable system to prevent AI agent confusion
   - [x] Mutual exclusivity between --enable and --disable flags
   - [x] Static builds for easy distribution and deployment

How to test
- [x] Plan document created and reviewed
- [x] Stakeholder approval of plan structure

Status: Plan completed and approved (2025-10-07, commit c81c391)

---

### [x] 1) Project Foundation Setup
Establish the basic project structure, dependencies, and build foundation.

1. **Project Structure**
   - [x] Create `cmd/mcpipboy/` directory for main application
   - [x] Create `internal/` directory for internal packages
   - [x] Create `internal/server/` for MCP server logic
   - [x] Create `internal/tools/` for tool implementations
   - [x] Create `internal/tools/interfaces.go` for common tool interfaces
   - [x] Create `internal/tools/interfaces_test.go` for interface tests
   - [x] Create `internal/tools/echo.go` for echo tool implementation
   - [x] Create `internal/tools/echo_test.go` for echo tool tests
   - [x] Create `internal/tools/version.go` for version tool implementation
   - [x] Create `internal/tools/version_test.go` for version tool tests
   - [x] Create `internal/tools/registry.go` for tool registration system
   - [x] Create `internal/tools/registry_test.go` for registry tests

2. **Dependency Management**
   - [x] Add MCP SDK dependency: `go get github.com/modelcontextprotocol/go-sdk`
   - [x] Add Cobra dependency: `go get github.com/spf13/cobra@latest`
   - [x] Add testing dependencies: `go get github.com/stretchr/testify`
   - [x] Run `go mod tidy` to clean dependencies

3. **Basic CLI Structure**
   - [x] Create `cmd/mcpipboy/main.go` with basic cobra root command
   - [x] Create `cmd/mcpipboy/serve.go` for MCP server mode with --enable/--disable flags
   - [x] Create `cmd/mcpipboy/serve_test.go` for serve command tests
   - [x] Create `cmd/mcpipboy/echo.go` for echo tool command
   - [x] Create `cmd/mcpipboy/echo_test.go` for echo command tests
   - [x] Create `cmd/mcpipboy/version.go` for version tool command
   - [x] Create `cmd/mcpipboy/version_test.go` for version command tests
   - [x] Implement tool enable/disable logic with mutual exclusivity validation
   - [x] Set up command structure to mirror internal tool organization

4. **Build System Setup**
   - [x] Create `justfile` in project root
   - [x] Add build targets: `build`, `build-release` (with static linking)
   - [x] Add test targets: `test`, `test-coverage`
   - [x] Add development targets: `dev`, `lint`, `fmt`
   - [x] Add workflow targets: `install`, `clean`, `deps`, `check`
   - [x] Add version bumping targets: `bump-patch`, `bump-minor`, `bump-major`, `bump-prerelease`
   - [x] Implement hidden default target using `_default` recipe
   - [x] Add descriptive comments for all recipes (no redundant descriptions)
   - [x] Configure release builds to use `dist/` directory (added to .gitignore)
   - [x] Add static linking with `CGO_ENABLED=0` and `-ldflags="-s -w"`

How to test
- `just` should show available recipes (hidden default target)
- `just build` should compile successfully to `bin/`
- `just build-release` should create static binaries in `dist/` for all platforms
- `just test` should run all tests and pass with >80% coverage
- `just test-coverage` should generate HTML coverage report
- `just bump-patch` should increment patch version (0.1.0 -> 0.1.1)
- `just bump-minor` should increment minor version (0.1.0 -> 0.2.0)
- `just bump-major` should increment major version (0.1.0 -> 1.0.0)
- `just bump-prerelease` should add prerelease version (0.1.0 -> 0.1.0-alpha1)
- `just check` should run fmt, vet, and test
- `just clean` should remove bin/, dist/, and coverage files
- All recipes should have descriptive comments (no redundant descriptions)
- Release binaries should be statically linked and work without dependencies

---

### [x] 2) MCP Server Core Implementation
Implement the basic MCP server with echo tool functionality.

5. **MCP Server Setup**
   - [x] Create `internal/server/server.go` with MCP server struct
   - [x] Create `internal/server/server_test.go` for server tests
   - [x] Implement server initialization and configuration
   - [x] Set up stdin/stdout communication handlers
   - [x] Add proper error handling and logging

6. **Tool Interface Definition**
   - [x] Define common tool interface in `internal/tools/interfaces.go`
   - [x] Create tool metadata structures (name, description, parameters)
   - [x] Define tool execution interface for both CLI and MCP usage
   - [x] Add input validation interface for tool parameters
   - [x] Implement tool registration interface

7. **Echo Tool Implementation**
   - [x] Implement echo tool struct conforming to tool interface
   - [x] Add input validation for echo message parameter
   - [x] Implement echo functionality that returns input message
   - [x] Add tool registration with MCP server
   - [x] Ensure CLI and MCP compatibility

8. **Version Tool Implementation**
   - [x] Create `VERSION` file with semantic version
   - [x] Implement version tool struct conforming to tool interface
   - [x] Implement version embedding using go:embed
   - [x] Add version parsing and validation
   - [x] Implement version functionality that returns current version
   - [x] Add tool registration with MCP server

9. **Tool Registration System**
   - [x] Implement tool registry for managing available tools
   - [x] Add tool discovery and listing functionality
   - [x] Implement enable/disable logic for tool filtering
   - [x] Add tool metadata management (name, description, parameters)
   - [x] Ensure proper tool registration and deregistration

How to test
- `./mcpipboy serve` should start MCP server without errors
- Send MCP request for tool list should return echo and version tools
- Send MCP request to call echo tool should return echoed message
- Send MCP request to call version tool should return current version
- Server should handle malformed requests gracefully

---

### [ ] 3) Build & Release System
Create comprehensive build automation, testing, and releases.

10. **Semantic Versioning**
    - [x] Add version management with `bump-version` target
    - [x] Implement automatic version injection in builds
    - [x] Add version validation and consistency checks
    - [x] Create version bumping for patch, minor, major

11. **Release Automation**
    - [x] Add `release` target for GitHub releases
    - [x] Implement cross-platform static builds (linux, darwin, windows)
    - [x] Add artifact generation and upload
    - [x] Include release notes generation

12. **Development Workflow**
    - [x] Add `install` target for local development
    - [x] Add `clean` target for build artifacts
    - [x] Add `deps` target for dependency management
    - [x] Add `check` target for pre-commit validation

How to test
- `just bump-version patch` should increment version correctly
- `just build-release` should create statically linked binaries for all platforms
- `ldd ./bin/mcpipboy` should show "not a dynamic executable" (on Linux)
- `otool -L ./bin/mcpipboy` should show minimal dependencies (on macOS)
- `just release` should create GitHub release with artifacts

---

### [ ] 4) Testing & Validation
Ensure MVP works correctly with MCP clients and meets requirements.

13. **Unit Testing**
    - [x] Create tests for echo tool functionality
    - [x] Add tests for MCP server initialization
    - [x] Test error handling and edge cases
    - [x] Achieve >80% test coverage

How to test
- Run `go test -v ./...` to execute all tests
- Run `go test -coverprofile=coverage.out ./internal/tools/... && go tool cover -func=coverage.out` to check coverage
- All tests should pass with >80% coverage in tools package

Status: **COMPLETE** - Enhanced echo tool tests with comprehensive edge cases (empty strings, unicode, special characters, error conditions). Added server error handling tests (nil tools, duplicate registration, multiple tools). Achieved 88.3% test coverage in tools package, exceeding the >80% requirement. All 25+ tests passing across all packages.

14. **Integration Testing**
    - [x] Test MCP server with actual MCP client
    - [x] Verify stdin/stdout communication works
    - [x] Test tool discovery and execution
    - [x] Validate JSON-RPC protocol compliance

How to test
- Run `go test -v -run "TestMCPServerIntegration" ./internal/server/...` for integration tests
- Run `go test -v -run "TestMCPServerProtocolCompliance" ./internal/server/...` for protocol compliance tests
- Tests verify MCP protocol communication, tool discovery, and tool execution

Status: **COMPLETE** - Created comprehensive integration tests that verify MCP server functionality with actual MCP protocol communication. Tests cover tool discovery (tools/list), tool execution (tools/call for echo and version tools), JSON-RPC protocol compliance, error handling for invalid requests, and unsupported methods. All integration tests passing with full MCP protocol validation.

15. **Documentation**
    - [x] Update README with usage instructions
    - [x] Add API documentation for tools
    - [x] Create example MCP client integration
    - [x] Document build and release process

How to test
- Review README.md for comprehensive usage instructions and integration guides
- Verify all justfile commands are properly documented
- Check that installation and usage instructions are clear and accurate

Status: **COMPLETE** - Created comprehensive README with installation, usage, and integration instructions for Cursor, Claude Desktop, Continue.dev, and custom MCP clients. Simplified documentation approach to focus on essential information in the README rather than extensive separate documentation files. All documentation properly references justfile commands where appropriate.

---

## Success Criteria
- [x] MCP server starts successfully and responds to tool discovery
- [x] Echo tool works correctly and returns input message
- [x] Build system creates cross-platform binaries
- [x] Semantic versioning works for releases
- [x] All tests pass and coverage is adequate
- [x] Documentation is complete and accurate

## Next Phase Preview
After MVP completion, the next phase will add the core utility tools:
- UUID generation (v1, v4, v5, v7)
- IMO number validation/generation
- MMSI number validation/generation  
- Credit card number validation/generation
- ISBN validation/generation
- Additional checksummed identifier types
