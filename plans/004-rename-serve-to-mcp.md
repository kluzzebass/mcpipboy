# Rename serve Command to mcp

**Date:** 2025-10-12

## Purpose

Rename the `serve` command to `mcp` throughout the codebase to better reflect its purpose as an MCP server. This includes renaming files, updating command definitions, documentation, and all references.

## Testing Strategy

Manual testing via CLI to verify the renamed command works correctly. Existing unit tests and integration tests should continue to pass without modification (they test the server functionality, not the command name).

**Execution Mode:** autonomous

Progress legend: [x] Completed, [ ] Pending

---

### [x] 0) Planning document and alignment

This refactoring task is straightforward and low-risk:
- Rename source files (serve.go, serve_test.go)
- Update command definitions, function names, and variable names
- Update all documentation references
- No changes to actual server functionality or test logic

**Architecture Decisions:**
- Keep the internal package structure unchanged (internal/server remains as-is)
- The command name changes from `serve` to `mcp`, but all flags remain the same
- All existing tests continue to work as they test the underlying functionality

**Rejected Alternatives:**
- Keeping `serve` as an alias: Rejected to avoid confusion and maintain a single clear command name
- Renaming the internal/server package: Rejected as "server" accurately describes what it does (serves the MCP protocol), while "mcp" better describes the user-facing command

**How to test:**
- Plan is complete and ready to implement

---

## Command Files

### [x] 1) Rename serve.go to mcp.go

Rename `/Users/kluzz/Code/mcpipboy/cmd/mcpipboy/serve.go` to `mcp.go` and update all internal references:
- Change package comment from "serve command" to "mcp command"
- Rename `serveCmd` variable to `mcpCmd`
- Update `Use: "serve"` to `Use: "mcp"`
- Rename `runServe` function to `runMCP`
- Update any internal comments referencing "serve"

**How to test:**
```bash
cd /Users/kluzz/Code/mcpipboy
go build -o bin/mcpipboy ./cmd/mcpipboy
./bin/mcpipboy mcp --help
```

### [x] 2) Rename serve_test.go to mcp_test.go

Rename `/Users/kluzz/Code/mcpipboy/cmd/mcpipboy/serve_test.go` to `mcp_test.go` and update all internal references:
- Rename test functions from `TestServe*` to `TestMCP*`
- Update test names and comments referencing "serve"
- Update command execution strings from "serve" to "mcp"

**How to test:**
```bash
cd /Users/kluzz/Code/mcpipboy
go test -v ./cmd/mcpipboy/... -run TestMCP
```

---

## Documentation Updates

### [x] 3) Update README.md

Update all references to the `serve` command in `/Users/kluzz/Code/mcpipboy/README.md`:
- Change "Start the MCP server" examples from `mcpipboy serve` to `mcpipboy mcp`
- Update all configuration examples with `args: ["serve"]` to `args: ["mcp"]`
- Update any other references to the serve command

**How to test:**
```bash
grep -n "serve" /Users/kluzz/Code/mcpipboy/README.md
# Should only show references to "server" (noun) not "serve" (verb/command)
```

### [x] 4) Update plan documents

**SKIPPED** - Completed plans (001, 002, 003) are historical documents and should not be modified. They accurately reflect that the command was called `serve` when those plans were executed. Only current/future documentation should reference the new `mcp` command name.

**How to test:**
```bash
# No testing needed - historical plans remain unchanged
```

### [x] 5) Update ROADMAP.md

No changes needed - only contains "MCP server integration" which correctly uses "server" as a noun.

**How to test:**
```bash
grep -n "serve" /Users/kluzz/Code/mcpipboy/ROADMAP.md
# Only shows "MCP server integration" - correct usage
```

---

## Build and Infrastructure

### [x] 6) Update justfile

No changes needed - only contains "MCP server for AI agents" which correctly uses "server" as a noun.

**How to test:**
```bash
grep -n "serve" /Users/kluzz/Code/mcpipboy/justfile
# Only shows "MCP server for AI agents" - correct usage
```

---

## Final Validation

### [x] 7) Run all tests and verify command works

Run the full test suite and verify the renamed command works as expected:
- Build the binary
- Run all tests
- Verify `mcpipboy mcp --help` shows proper documentation
- Test with sample flags: `--enable`, `--disable`, `--debug`

**How to test:**
```bash
cd /Users/kluzz/Code/mcpipboy
just build
just test
./bin/mcpipboy mcp --help
./bin/mcpipboy mcp --enable echo,version --debug --log-file /tmp/mcp-debug.log &
# Send a test request and check it works
kill %1
```

---

## Success Criteria

- All files renamed appropriately (serve.go → mcp.go, serve_test.go → mcp_test.go)
- All function and variable names updated (serveCmd → mcpCmd, runServe → runMCP)
- All documentation updated (README.md, plan documents, ROADMAP.md)
- All tests pass
- Command `mcpipboy mcp` works identically to old `mcpipboy serve`
- No references to `serve` command remain (only "server" as a noun is acceptable)

