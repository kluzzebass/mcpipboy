## [mcpipboy Core Tools Implementation] ([Add UUID generation, checksum validation, and data generation tools])

Progress legend: [x] Completed, [>] In progress, [ ] Pending

### [ ] 0) Planning document and alignment
Create comprehensive plan for implementing core utility tools that AI agents commonly need.

1. **Tool Requirements Analysis**
   - [ ] Analyze requirements for UUID generation tools (v1, v4, v5, v7)
   - [ ] Define IMO number validation and generation requirements
   - [ ] Specify MMSI number validation and generation requirements
   - [ ] Plan credit card number validation and generation (Luhn algorithm)
   - [ ] Design ISBN validation and generation (ISBN-10, ISBN-13)
   - [ ] Identify additional checksummed identifier types needed

2. **Architecture Decisions**
   - [ ] Decide on UUID library (standard library vs external)
   - [x] Choose lenient time parsing library (selected: github.com/ijt/go-anytime/v2)
   - [ ] Plan checksum algorithm implementations
   - [ ] Design tool parameter validation strategies
   - [ ] Plan error handling for invalid inputs
   - [ ] Decide on output formats for generated data

3. **Implementation Strategy**
   - [ ] Plan tool implementation order (start with UUID, then checksums)
   - [ ] Design consistent tool interface patterns
   - [ ] Plan testing strategy for each tool type
   - [ ] Design CLI integration for new tools
   - [ ] Plan documentation updates

How to test
- Review plan completeness and technical feasibility
- Validate tool requirements against real-world use cases
- Ensure architecture decisions align with existing codebase patterns

Status: Planning phase - document creation and requirements analysis

---

### [ ] 1) Time Tool
Implement a comprehensive time utility that AI agents commonly need.

1. **Time Tool**
   - [x] Research and select lenient time parsing library (selected: github.com/ijt/go-anytime/v2)
   - [x] Add go-anytime library to go.mod dependencies
   - [ ] Create `internal/tools/time.go` with flexible time functionality
   - [ ] Implement time type options: "now", "today", "timestamp", "unix", "relative"
   - [ ] Add format options: "iso", "rfc3339", "unix", "date", "datetime", "time"
   - [ ] Add timezone support: "utc", "local", or specified timezone
   - [ ] Add relative time calculations: from any timestamp to any timestamp
   - [ ] Add offset support: relative time calculations from input timestamp (+1h, -2d, etc.)
   - [ ] Add comprehensive test coverage in `internal/tools/time_test.go`
   - [ ] Add CLI command in `cmd/mcpipboy/time.go`
   - [ ] Add CLI tests in `cmd/mcpipboy/time_test.go`

How to test
- Run `just test` to ensure all time tests pass
- Test CLI commands: `mcpipboy time`, `mcpipboy time --type today`, `mcpipboy time --format unix`
- Test relative calculations: `mcpipboy time --type relative --from "2025-01-01" --to "2025-12-31"`
- Test offset calculations: `mcpipboy time --type timestamp --input "2025-01-01" --offset "+1h"`
- Verify MCP integration: tool appears in `tools/list` and executes via `tools/call`
- Test different combinations: type + format + timezone + offset/relative calculations
- Verify output accuracy and consistency across different options

Status: Pending - Time tool implementation

---

### [ ] 2) UUID Tool
Implement a comprehensive UUID generation and validation tool with version selection.

1. **UUID Tool**
   - [ ] Create `internal/tools/uuid.go` with flexible UUID functionality
   - [ ] Implement version parameter: "v1", "v4", "v5", "v7", "validate"
   - [ ] Add UUID v1 (time-based) generation with MAC address handling
   - [ ] Add UUID v4 (random) generation with count parameter
   - [ ] Add UUID v5 (name-based SHA-1) generation with namespace and name parameters
   - [ ] Add UUID v7 (time-ordered) generation
   - [ ] Add UUID validation functionality for any version
   - [ ] Add comprehensive test coverage in `internal/tools/uuid_test.go`
   - [ ] Add CLI command in `cmd/mcpipboy/uuid.go`
   - [ ] Add CLI tests in `cmd/mcpipboy/uuid_test.go`

How to test
- Run `just test` to ensure all UUID tests pass
- Test CLI commands: `mcpipboy uuid --version v4`, `mcpipboy uuid --version v1 --count 5`
- Test validation: `mcpipboy uuid --version validate --input "550e8400-e29b-41d4-a716-446655440000"`
- Test v5 generation: `mcpipboy uuid --version v5 --namespace "6ba7b810-9dad-11d1-80b4-00c04fd430c8" --name "example"`
- Verify MCP integration: tool appears in `tools/list` and executes via `tools/call`
- Test edge cases: invalid versions, invalid UUIDs, boundary conditions

Status: Pending - UUID tool implementation

---

### [ ] 3) IMO Tool
Implement International Maritime Organization (IMO) number validation and generation.

1. **IMO Tool**
   - [ ] Create `internal/tools/imo.go` with flexible IMO functionality
   - [ ] Implement operation parameter: "validate", "generate"
   - [ ] Add IMO validation with checksum algorithm (7-digit number with check digit)
   - [ ] Add IMO generation with correct checksum calculation
   - [ ] Add comprehensive test coverage in `internal/tools/imo_test.go`
   - [ ] Add CLI command in `cmd/mcpipboy/imo.go`
   - [ ] Add CLI tests in `cmd/mcpipboy/imo_test.go`

How to test
- Run `just test` to ensure all IMO tests pass
- Test validation: `mcpipboy imo --operation validate --input "1234567"`
- Test generation: `mcpipboy imo --operation generate`
- Verify MCP integration: tool appears in `tools/list` and executes via `tools/call`
- Test with known valid/invalid IMO numbers and checksum validation

Status: Pending - IMO tool implementation

---

### [ ] 4) MMSI Tool
Implement Maritime Mobile Service Identity (MMSI) number validation and generation with country code support.

1. **MMSI Tool**
   - [ ] Create `internal/tools/mmsi.go` with flexible MMSI functionality
   - [ ] Implement operation parameter: "validate", "generate"
   - [ ] Add MMSI validation with format checking (9-digit number, country codes)
   - [ ] Add MMSI generation with optional country code parameter
   - [ ] Implement country code validation and lookup
   - [ ] Add comprehensive test coverage in `internal/tools/mmsi_test.go`
   - [ ] Add CLI command in `cmd/mcpipboy/mmsi.go`
   - [ ] Add CLI tests in `cmd/mcpipboy/mmsi_test.go`

How to test
- Run `just test` to ensure all MMSI tests pass
- Test validation: `mcpipboy mmsi --operation validate --input "123456789"`
- Test generation: `mcpipboy mmsi --operation generate --country-code "US"`
- Test generation without country: `mcpipboy mmsi --operation generate`
- Verify MCP integration: tool appears in `tools/list` and executes via `tools/call`
- Test with known valid/invalid MMSI numbers and country codes

Status: Pending - MMSI tool implementation

---

### [ ] 5) Credit Card Tool
Implement credit card number validation and generation using Luhn algorithm with card type support.

1. **Credit Card Tool**
   - [ ] Create `internal/tools/creditcard.go` with flexible credit card functionality
   - [ ] Implement operation parameter: "validate", "generate"
   - [ ] Add credit card validation with Luhn algorithm and card type detection
   - [ ] Add credit card generation with optional card type parameter (Visa, Mastercard, Amex, etc.)
   - [ ] Implement Luhn algorithm for check digit calculation
   - [ ] Add comprehensive test coverage in `internal/tools/creditcard_test.go`
   - [ ] Add CLI command in `cmd/mcpipboy/creditcard.go`
   - [ ] Add CLI tests in `cmd/mcpipboy/creditcard_test.go`

How to test
- Run `just test` to ensure all credit card tests pass
- Test validation: `mcpipboy creditcard --operation validate --input "4111111111111111"`
- Test generation: `mcpipboy creditcard --operation generate --card-type "visa"`
- Test generation without type: `mcpipboy creditcard --operation generate`
- Verify MCP integration: tool appears in `tools/list` and executes via `tools/call`
- Test with known valid/invalid credit card numbers and different card types

Status: Pending - Credit card tool implementation

---

### [ ] 6) ISBN Tool
Implement International Standard Book Number (ISBN) validation and generation with format support.

1. **ISBN Tool**
   - [ ] Create `internal/tools/isbn.go` with flexible ISBN functionality
   - [ ] Implement operation parameter: "validate", "generate"
   - [ ] Add ISBN validation with format parameter: "isbn10", "isbn13", "auto"
   - [ ] Add ISBN generation with optional format parameter (defaults to ISBN-13)
   - [ ] Implement ISBN-10 checksum algorithm
   - [ ] Implement ISBN-13 checksum algorithm
   - [ ] Add comprehensive test coverage in `internal/tools/isbn_test.go`
   - [ ] Add CLI command in `cmd/mcpipboy/isbn.go`
   - [ ] Add CLI tests in `cmd/mcpipboy/isbn_test.go`

How to test
- Run `just test` to ensure all ISBN tests pass
- Test validation: `mcpipboy isbn --operation validate --input "978-0-123456-78-9" --format "isbn13"`
- Test auto-detection: `mcpipboy isbn --operation validate --input "0-123456-78-9"`
- Test generation: `mcpipboy isbn --operation generate --format "isbn10"`
- Test generation without format: `mcpipboy isbn --operation generate` (defaults to ISBN-13)
- Verify MCP integration: tool appears in `tools/list` and executes via `tools/call`
- Test with known valid/invalid ISBN numbers and format conversions

Status: Pending - ISBN tool implementation

---

### [ ] 7) EAN-13 Tool
Implement European Article Number (EAN-13) validation and generation.

1. **EAN-13 Tool**
   - [ ] Create `internal/tools/ean13.go` with flexible EAN-13 functionality
   - [ ] Implement operation parameter: "validate", "generate"
   - [ ] Add EAN-13 validation with checksum algorithm
   - [ ] Add EAN-13 generation with valid checksum
   - [ ] Add comprehensive test coverage in `internal/tools/ean13_test.go`
   - [ ] Add CLI command in `cmd/mcpipboy/ean13.go`
   - [ ] Add CLI tests in `cmd/mcpipboy/ean13_test.go`

How to test
- Run `just test` to ensure all EAN-13 tests pass
- Test validation: `mcpipboy ean13 --operation validate --input "1234567890123"`
- Test generation: `mcpipboy ean13 --operation generate`
- Verify MCP integration: tool appears in `tools/list` and executes via `tools/call`
- Test with known valid/invalid EAN-13 numbers

Status: Pending - EAN-13 tool implementation

---

### [ ] 8) IBAN Tool
Implement International Bank Account Number (IBAN) validation and generation.

1. **IBAN Tool**
   - [ ] Create `internal/tools/iban.go` with flexible IBAN functionality
   - [ ] Implement operation parameter: "validate", "generate"
   - [ ] Add IBAN validation with MOD-97 checksum algorithm
   - [ ] Add IBAN generation with optional country code parameter
   - [ ] Add comprehensive test coverage in `internal/tools/iban_test.go`
   - [ ] Add CLI command in `cmd/mcpipboy/iban.go`
   - [ ] Add CLI tests in `cmd/mcpipboy/iban_test.go`

How to test
- Run `just test` to ensure all IBAN tests pass
- Test validation: `mcpipboy iban --operation validate --input "GB82WEST12345698765432"`
- Test generation: `mcpipboy iban --operation generate --country-code "GB"`
- Test generation without country: `mcpipboy iban --operation generate`
- Verify MCP integration: tool appears in `tools/list` and executes via `tools/call`
- Test with known valid/invalid IBAN numbers and country codes

Status: Pending - IBAN tool implementation

---

### [ ] 9) VIN Tool
Implement Vehicle Identification Number (VIN) validation and generation.

1. **VIN Tool**
   - [ ] Create `internal/tools/vin.go` with flexible VIN functionality
   - [ ] Implement operation parameter: "validate", "generate"
   - [ ] Add VIN validation with checksum algorithm
   - [ ] Add VIN generation with valid checksum
   - [ ] Add comprehensive test coverage in `internal/tools/vin_test.go`
   - [ ] Add CLI command in `cmd/mcpipboy/vin.go`
   - [ ] Add CLI tests in `cmd/mcpipboy/vin_test.go`

How to test
- Run `just test` to ensure all VIN tests pass
- Test validation: `mcpipboy vin --operation validate --input "1HGBH41JXMN109186"`
- Test generation: `mcpipboy vin --operation generate`
- Verify MCP integration: tool appears in `tools/list` and executes via `tools/call`
- Test with known valid/invalid VIN numbers

Status: Pending - VIN tool implementation

---

### [ ] 10) Documentation and Integration
Update documentation and ensure all new tools are properly integrated.

1. **README Updates**
   - [ ] Update README.md with new tool descriptions
   - [ ] Add usage examples for new tools
   - [ ] Update installation and integration instructions

2. **Tool Registry Updates**
   - [ ] Ensure all new tools are properly registered
   - [ ] Update tool discovery and execution
   - [ ] Verify CLI help and completion

3. **Testing and Validation**
   - [ ] Run full test suite to ensure no regressions
   - [ ] Verify MCP integration with all new tools
   - [ ] Test CLI functionality for all new tools
   - [ ] Validate tool schemas and parameter validation

How to test
- Run `just test` to ensure all tests pass
- Test `mcpipboy serve` and verify all tools appear in `tools/list`
- Test CLI commands for all new tools
- Verify documentation is accurate and complete

Status: Pending - Documentation and integration updates

---

## Success Criteria
- [ ] Time tool with flexible options (type, format, timezone, offset) implemented and working
- [ ] All UUID generation tools (v1, v4, v5, v7) implemented and working
- [ ] IMO number validation and generation tools working
- [ ] MMSI number validation and generation tools working
- [ ] Credit card number validation and generation tools working
- [ ] ISBN validation and generation tools working
- [ ] Additional checksum tools implemented
- [ ] All tools have comprehensive test coverage (>80%)
- [ ] All tools work via both CLI and MCP protocol
- [ ] Documentation updated with all new tools
- [ ] No regressions in existing functionality

## Next Phase Preview
After core tools completion, the next phase will add advanced features:
- Date/time utilities and formatting
- Random data generation (names, addresses, etc.)
- Data validation and sanitization tools
- IP calculation tool (IPv4/IPv6 subnet calculations, CIDR operations, network analysis)
- Crypto tool (password hashing, verification, and various hash algorithms)
- Text similarity tool (Levenshtein distance, fuzzy matching, string comparison)
