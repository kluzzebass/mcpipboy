# mcpipboy MVP Development Plan

## Purpose
Create a stdin/stdout MCP (Model Context Protocol) server that provides agentic AIs with essential tools for common tasks they struggle with, including UUID generation, checksummed identifier verification/generation (IMO, MMSI, credit card numbers, ISBN, etc.), and other utility functions.

## High Level Phases
1. **Project Foundation** - Set up build system, dependencies, and basic project structure
2. **MCP Server Core** - Implement basic MCP server with echo tool
3. **Build & Release System** - Create justfile for build automation and semantic versioning
4. **Testing & Validation** - Ensure MVP works correctly with MCP clients

Progress legend: [x] Completed, [>] In progress, [ ] Pending

---

### [ ] 0) Planning document and alignment
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

How to test
- [x] Plan document created and reviewed
- [ ] Stakeholder approval of plan structure

Status: Plan created, awaiting approval

---

### [ ] 1) Project Foundation Setup
Establish the basic project structure, dependencies, and build foundation.

1. **Dependency Management**
   - [ ] Add MCP SDK dependency: `go get github.com/modelcontextprotocol/go-sdk`
   - [ ] Add Cobra dependency: `go get github.com/spf13/cobra@latest`
   - [ ] Add testing dependencies: `go get github.com/stretchr/testify`
   - [ ] Run `go mod tidy` to clean dependencies

2. **Project Structure**
   - [ ] Create `cmd/mcpipboy/` directory for main application
   - [ ] Create `internal/` directory for internal packages
   - [ ] Create `internal/server/` for MCP server logic
   - [ ] Create `internal/tools/` for tool implementations

3. **Basic CLI Structure**
   - [ ] Create `cmd/mcpipboy/main.go` with basic cobra root command
   - [ ] Add version command to display semantic version
   - [ ] Add serve command for MCP server mode

How to test
- `go build ./cmd/mcpipboy` should compile successfully
- `./mcpipboy --help` should display help text
- `./mcpipboy version` should display version information

---

### [ ] 2) MCP Server Core Implementation
Implement the basic MCP server with echo tool functionality.

1. **MCP Server Setup**
   - [ ] Create `internal/server/server.go` with MCP server struct
   - [ ] Implement server initialization and configuration
   - [ ] Set up stdin/stdout communication handlers
   - [ ] Add proper error handling and logging

2. **Echo Tool Implementation**
   - [ ] Create `internal/tools/echo.go` with echo tool struct
   - [ ] Implement tool registration with MCP server
   - [ ] Add input validation for echo message
   - [ ] Implement echo functionality that returns input message

3. **Tool Registration System**
   - [ ] Create `internal/tools/registry.go` for tool management
   - [ ] Implement tool registration interface
   - [ ] Add tool discovery and listing functionality
   - [ ] Ensure proper tool metadata (name, description, parameters)

How to test
- `./mcpipboy serve` should start MCP server without errors
- Send MCP request for tool list should return echo tool
- Send MCP request to call echo tool should return echoed message
- Server should handle malformed requests gracefully

---

### [ ] 3) Build & Release System
Create comprehensive justfile for build automation, testing, and releases.

1. **Justfile Creation**
   - [ ] Create `justfile` in project root
   - [ ] Add build targets: `build`, `build-release`
   - [ ] Add test targets: `test`, `test-coverage`
   - [ ] Add development targets: `dev`, `lint`, `fmt`

2. **Semantic Versioning**
   - [ ] Add version management with `bump-version` target
   - [ ] Implement automatic version injection in builds
   - [ ] Add version validation and consistency checks
   - [ ] Create version bumping for patch, minor, major

3. **Release Automation**
   - [ ] Add `release` target for GitHub releases
   - [ ] Implement cross-platform builds (linux, darwin, windows)
   - [ ] Add artifact generation and upload
   - [ ] Include release notes generation

4. **Development Workflow**
   - [ ] Add `install` target for local development
   - [ ] Add `clean` target for build artifacts
   - [ ] Add `deps` target for dependency management
   - [ ] Add `check` target for pre-commit validation

How to test
- `just --list` should show all available targets
- `just build` should create binary successfully
- `just test` should run all tests and pass
- `just bump-version patch` should increment version correctly

---

### [ ] 4) Testing & Validation
Ensure MVP works correctly with MCP clients and meets requirements.

1. **Unit Testing**
   - [ ] Create tests for echo tool functionality
   - [ ] Add tests for MCP server initialization
   - [ ] Test error handling and edge cases
   - [ ] Achieve >80% test coverage

2. **Integration Testing**
   - [ ] Test MCP server with actual MCP client
   - [ ] Verify stdin/stdout communication works
   - [ ] Test tool discovery and execution
   - [ ] Validate JSON-RPC protocol compliance

3. **Documentation**
   - [ ] Update README with usage instructions
   - [ ] Add API documentation for tools
   - [ ] Create example MCP client integration
   - [ ] Document build and release process

4. **Performance & Reliability**
   - [ ] Test server stability under load
   - [ ] Verify memory usage is reasonable
   - [ ] Test graceful shutdown handling
   - [ ] Validate error recovery mechanisms

How to test
- All unit tests pass with `just test`
- Integration test with MCP client succeeds
- Server handles 100+ consecutive requests without issues
- Documentation is clear and complete

---

## Success Criteria
- [ ] MCP server starts successfully and responds to tool discovery
- [ ] Echo tool works correctly and returns input message
- [ ] Build system creates cross-platform binaries
- [ ] Semantic versioning works for releases
- [ ] All tests pass and coverage is adequate
- [ ] Documentation is complete and accurate

## Next Phase Preview
After MVP completion, the next phase will add the core utility tools:
- UUID generation (v1, v4, v5, v7)
- IMO number validation/generation
- MMSI number validation/generation  
- Credit card number validation/generation
- ISBN validation/generation
- Additional checksummed identifier types
