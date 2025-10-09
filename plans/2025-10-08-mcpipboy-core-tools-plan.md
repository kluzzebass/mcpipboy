## [mcpipboy Core Tools Implementation] ([Add UUID generation, checksum validation, and data generation tools])

Progress legend: [x] Completed, [>] In progress, [ ] Pending

### [x] 0) Planning document and alignment
Create comprehensive plan for implementing core utility tools that AI agents commonly need.

1. **Tool Requirements Analysis**
   - [x] Analyze requirements for UUID generation tools (v1, v4, v5, v7)
     - UUID v1: Time-based with MAC address, handle privacy concerns
     - UUID v4: Cryptographically secure random (most common)
     - UUID v5: Name-based SHA-1 with namespace support (DNS, URL, OID, X.500)
     - UUID v7: Time-ordered with random component (better for DB indexing)
     - Validation: Format check, version bits, variant bits verification
   - [x] Define IMO number validation and generation requirements
     - Format: 7-digit number (6 digits + 1 check digit)
     - Check digit: Weighted sum algorithm
     - Validation: Length, digits only, check digit verification
     - Generation: 6 random digits + calculated check digit
   - [x] Specify MMSI number validation and generation requirements
     - Format: 9-digit number (3-digit country code + 6-digit station ID)
     - Country codes: 201-775 (maritime nations)
     - Validation: Length, digits only, country code range check
     - Generation: Optional country code parameter, random station ID
   - [x] Plan credit card number validation and generation (Luhn algorithm)
     - Card types: Visa (4), Mastercard (51-55, 2221-2720), AmEx (34,37), Discover (6011,65,644-649)
     - Lengths: 13-19 digits (varies by type)
     - Validation: Luhn algorithm, prefix matching, length verification
     - Generation: Optional card type, appropriate prefix, Luhn check digit
   - [x] Design ISBN validation and generation (ISBN-10, ISBN-13)
     - ISBN-10: 10 digits with weighted sum check digit (0-9 or X)
     - ISBN-13: 13 digits with EAN-13 check digit algorithm
     - Validation: Format detection, check digit verification
     - Generation: Optional format parameter (defaults to ISBN-13)
   - [x] Identify additional checksummed identifier types needed
     - EAN-13: 13-digit product barcode with EAN-13 check digit
     - IBAN: Country code + check digits + account number with MOD-97 algorithm
     - VIN: 17-character vehicle ID with weighted sum check digit (no I,O,Q)

2. **Architecture Decisions**
   - [x] Decide on UUID library (standard library vs external)
     - Decision: Use Go standard library only (crypto/rand, crypto/sha1, time, net)
     - No external UUID library needed - can implement v1, v4, v5, v7 with stdlib
   - [x] Choose lenient time parsing library (selected: github.com/ijt/go-anytime/v2)
   - [x] Plan checksum algorithm implementations
     - Luhn Algorithm: Credit cards - double every 2nd digit, sum, check digit = (10 - sum%10) % 10
     - Weighted Sum: IMO (7×1+6×2+...), ISBN-10 (10×d1+9×d2+...), VIN (position-based multipliers)
     - EAN-13: Odd×1 + Even×3, check digit = (10 - sum%10) % 10
     - MOD-97: IBAN - move first 4 chars to end, replace letters with numbers, remainder = 1
   - [x] Design tool parameter validation strategies
     - Use existing ValidateParams method in Tool interface
     - Parameter types: operation, version, type, format, card-type, country-code, count
     - Validation at CLI and MCP server levels before tool execution
     - Descriptive error messages for invalid parameters
     - Leverage existing ValidateParameters function from interfaces.go
   - [x] Plan error handling for invalid inputs
     - Parameter validation errors: Invalid types, missing params, out-of-range values
     - Format validation errors: Wrong length, invalid characters, malformed input
     - Algorithm errors: Checksum failures, generation failures
     - System errors: Random generation failures, network access issues
     - Use NewErrorResult for consistent error responses with context
     - Graceful degradation where possible (e.g., random MAC if real MAC unavailable)
   - [x] Decide on output formats for generated data
     - Single item: Return generated/validated item directly
     - Multiple items: Return array when count > 1
     - Validation results: Boolean + details (valid/invalid, card type, country, etc.)
     - Consistent structure: Use NewSuccessResult wrapper for all successful responses
     - Examples: UUID string/array, credit card validation object, MMSI string/array

3. **Implementation Strategy**
   - [x] Plan tool implementation order (start with UUID, then checksums)
     - Time Tool → Random Number Generator → UUID → IMO → MMSI → Credit Card → ISBN → EAN-13 → IBAN → VIN → Documentation
     - Start with foundational utilities (Time, Random), then most common (UUID), then simple checksums, then complex ones
     - Each tool builds on previous patterns and utilities
   - [x] Design consistent tool interface patterns
     - Tool struct pattern: [Name]Tool with New[Name]Tool() constructor
     - Required methods: Name(), Description(), Execute(), ValidateParams(), GetInputSchema(), GetOutputSchema()
     - CLI pattern: [name]Cmd with Use, Short, Long, GroupID, Args, RunE
     - Consistent parameter validation using existing ValidateParameters function
     - Standard error handling with NewErrorResult and NewSuccessResult
   - [x] Plan testing strategy for each tool type
     - Unit tests: internal/tools/[name]_test.go for each tool
     - CLI tests: cmd/mcpipboy/[name]_test.go for each command
     - Integration tests: Update internal/server/integration_test.go
     - Test categories: Valid inputs, invalid inputs, edge cases, error conditions, schema validation
     - Follow existing test patterns from echo and version tools
   - [x] Design CLI integration for new tools
     - Tool registration: Add to getAvailableTools() function in cmd/mcpipboy/serve.go
     - Command registration: Add to cmd/mcpipboy/main.go in "Tool Commands" group
     - Flag handling: Use Cobra's built-in validation and completion
     - Help integration: Automatic help generation and command grouping
     - Shell completion: Automatic tab completion with Cobra
   - [x] Plan documentation updates
     - README.md: Update "Available Tools" section with new tool descriptions
     - Tool descriptions: Brief descriptions with key parameters and usage examples
     - CLI examples: Common use cases for each tool
     - Integration instructions: Update MCP client examples if needed
     - Keep concise: Avoid detailed API documentation per previous feedback

How to test
- Review plan completeness and technical feasibility
- Validate tool requirements against real-world use cases
- Ensure architecture decisions align with existing codebase patterns

Status: Planning phase - document creation and requirements analysis

---

## Phase 1: Core Tools Implementation

Implement a comprehensive time utility that AI agents commonly need.

4. **Time Tool**
   - [x] Research and select lenient time parsing library (selected: github.com/ijt/go-anytime v1)
   - [x] Add go-anytime library to go.mod dependencies
   - [x] Create `internal/tools/time.go` with flexible time functionality
   - [x] Implement simplified API: input parsing with go-anytime, output formatting with our logic
   - [x] Add format options: "iso", "rfc3339", "unix", "date", "datetime", "time"
   - [x] Add timezone support: "utc", "local", or specified timezone
   - [x] Add offset support: relative time calculations from input timestamp (+1h, -2d, etc.)
   - [x] Add comprehensive test coverage in `internal/tools/time_test.go`
   - [x] Add CLI command in `cmd/mcpipboy/time.go`
   - [x] Add CLI tests in `cmd/mcpipboy/time_test.go`

How to test
- Run `just test` to ensure all time tests pass
- Test CLI commands: `mcpipboy time`, `mcpipboy time --input "2025-01-01" --format unix`
- Test offset calculations: `mcpipboy time --input "2025-01-01" --offset "+1h"`
- Test timezone conversion: `mcpipboy time --input "2025-01-01T12:00:00Z" --timezone "America/New_York"`
- Verify MCP integration: tool appears in `tools/list` and executes via `tools/call`
- Test different combinations: input + format + timezone + offset
- Verify output accuracy and consistency across different options

Status: Complete - Time tool implementation with go-anytime v1 integration

---

Implement a comprehensive random number generator with various types and distributions.

5. **Random Number Generator Tool** [x]
   - [x] Create `internal/tools/random.go` with flexible random number functionality
   - [x] Implement type parameter: "integer", "float", "boolean"
   - [x] Add integer generation with min/max range and count parameter
   - [x] Add float generation with min/max range, precision, and count parameter
   - [x] Add boolean generation with count parameter
   - [x] Add comprehensive test coverage in `internal/tools/random_test.go`
   - [x] Add CLI command in `cmd/mcpipboy/random.go`
   - [x] Register random tool with MCP server

How to test
- Run `just test` to ensure all random number tests pass
- Test integer generation: `mcpipboy random --type integer --min 1 --max 100 --count 10`
- Test float generation: `mcpipboy random --type float --min 0.0 --max 1.0 --precision 2 --count 5`
- Test boolean generation: `mcpipboy random --type boolean --count 10`
- Test distributions: `mcpipboy random --type integer --distribution "normal" --mean 50 --stddev 10 --count 5`
- Verify MCP integration: tool appears in `tools/list` and executes via `tools/call`
- Test edge cases: invalid ranges, boundary conditions, different distributions

Status: Pending - Random number generator tool implementation

---

Implement a comprehensive UUID generation and validation tool with version selection.

6. **UUID Tool** [x]
   - [x] Create `internal/tools/uuid.go` with flexible UUID functionality
   - [x] Implement version parameter: "v1", "v4", "v5", "v7", "validate"
   - [x] Add UUID v1 (time-based) generation with MAC address handling and count parameter
   - [x] Add UUID v4 (random) generation with count parameter
   - [x] Add UUID v5 (name-based SHA-1) generation with namespace and name parameters and count parameter
   - [x] Add UUID v7 (time-ordered) generation with count parameter
   - [x] Add UUID validation functionality for any version
   - [x] Add enhanced validation with timestamp extraction (v1, v7) and MAC address (v1)
   - [x] Add comprehensive test coverage in `internal/tools/uuid_test.go`
   - [x] Add CLI command in `cmd/mcpipboy/uuid.go`
   - [x] Register UUID tool with MCP server

How to test
- Run `just test` to ensure all UUID tests pass
- Test CLI commands: `mcpipboy uuid --version v4`, `mcpipboy uuid --version v7 --count 10`
- Test validation: `mcpipboy uuid --version validate --input "550e8400-e29b-41d4-a716-446655440000"`
- Test v5 generation: `mcpipboy uuid --version v5 --namespace "6ba7b810-9dad-11d1-80b4-00c04fd430c8" --name "example"`
- Test enhanced validation: `mcpipboy uuid --version validate --input "6ba7b810-9dad-11d1-80b4-00c04fd430c8"`
- Verify MCP integration: tool appears in `tools/list` and executes via `tools/call`
- Test edge cases: invalid versions, invalid UUIDs, boundary conditions

Status: Complete - UUID tool with enhanced validation and metadata extraction

---

Implement International Maritime Organization (IMO) number validation and generation.

7. **IMO Tool** [x]
   - [x] Create `internal/tools/imo.go` with flexible IMO functionality
   - [x] Implement operation parameter: "validate", "generate"
   - [x] Add IMO validation with checksum algorithm (7-digit number with check digit)
   - [x] Add IMO generation with correct checksum calculation and count parameter
   - [x] Add comprehensive test coverage in `internal/tools/imo_test.go`
   - [x] Add CLI command in `cmd/mcpipboy/imo.go`
   - [x] Add CLI tests in `cmd/mcpipboy/imo_test.go`
   - [x] Register IMO tool with MCP server

How to test
- Run `just test` to ensure all IMO tests pass
- Test CLI commands: `mcpipboy imo --operation validate --input "1234567"`
- Test generation: `mcpipboy imo --operation generate --count 5`
- Test validation: `mcpipboy imo --operation validate --input "1234568"` (should fail)
- Verify MCP integration: tool appears in `tools/list` and executes via `tools/call`
- Test edge cases: invalid operations, missing input, count limits

Status: Complete - IMO tool with validation and generation functionality

---

Implement Maritime Mobile Service Identity (MMSI) number validation and generation with country code support.

8. **MMSI Tool**
   - [x] Create `internal/tools/mmsi.go` with flexible MMSI functionality
   - [x] Implement operation parameter: "validate", "generate"
   - [x] Add MMSI validation with format checking (9-digit number, country codes)
   - [x] Add MMSI generation with optional country code parameter and count parameter
   - [x] Implement country code validation and lookup
   - [x] Add comprehensive test coverage in `internal/tools/mmsi_test.go`
   - [x] Add CLI command in `cmd/mcpipboy/mmsi.go`
   - [x] **Enhanced with type-specific generation** - Added 16 supported MMSI types (ship, sar-aircraft, us-federal, etc.) with dedicated generation functions and CLI `--type` parameter
   - [x] Add CLI tests in `cmd/mcpipboy/mmsi_test.go`

How to test
- Run `just test` to ensure all MMSI tests pass
- Test validation: `mcpipboy mmsi --operation validate --input "123456789"`
- Test generation: `mcpipboy mmsi --operation generate --country-code "US" --count 3`
- Test generation without country: `mcpipboy mmsi --operation generate --count 5`
- Verify MCP integration: tool appears in `tools/list` and executes via `tools/call`
- Test with known valid/invalid MMSI numbers and country codes

Status: Complete - MMSI tool with comprehensive type-specific generation and validation

---

Implement credit card number validation and generation using Luhn algorithm with card type support.

9. **Credit Card Tool**
   - [x] Create `internal/tools/creditcard.go` with flexible credit card functionality
   - [x] Implement operation parameter: "validate", "generate"
   - [x] Add credit card validation with Luhn algorithm and card type detection
   - [x] Add credit card generation with optional card type parameter (Visa, Mastercard, Amex, etc.) and count parameter
   - [x] Implement Luhn algorithm for check digit calculation
   - [x] Add comprehensive test coverage in `internal/tools/creditcard_test.go`
   - [x] Add CLI command in `cmd/mcpipboy/creditcard.go`
   - [x] Add CLI tests in `cmd/mcpipboy/creditcard_test.go`
   - [x] Register credit card tool with MCP server

How to test
- Run `just test` to ensure all credit card tests pass
- Test validation: `mcpipboy creditcard --operation validate --input "4111111111111111"`
- Test generation: `mcpipboy creditcard --operation generate --card-type "visa" --count 3`
- Test generation without type: `mcpipboy creditcard --operation generate --count 5`
- Verify MCP integration: tool appears in `tools/list` and executes via `tools/call`
- Test with known valid/invalid credit card numbers and different card types

Status: Complete - Credit card tool with Luhn algorithm validation and generation

---

Implement International Standard Book Number (ISBN) validation and generation with format support.

10. **ISBN Tool**
   - [x] Create `internal/tools/isbn.go` with flexible ISBN functionality
   - [x] Implement operation parameter: "validate", "generate"
   - [x] Add ISBN validation with format parameter: "isbn10", "isbn13", "auto"
   - [x] Add ISBN generation with optional format parameter (defaults to ISBN-13) and count parameter
   - [x] Implement ISBN-10 checksum algorithm
   - [x] Implement ISBN-13 checksum algorithm
   - [x] Add comprehensive test coverage in `internal/tools/isbn_test.go`
   - [x] Add CLI command in `cmd/mcpipboy/isbn.go`
   - [x] Add CLI tests in `cmd/mcpipboy/isbn_test.go`

How to test
- Run `just test` to ensure all ISBN tests pass
- Test validation: `mcpipboy isbn --operation validate --input "978-0-123456-78-9" --format "isbn13"`
- Test auto-detection: `mcpipboy isbn --operation validate --input "0-123456-78-9"`
- Test generation: `mcpipboy isbn --operation generate --format "isbn10" --count 3`
- Test generation without format: `mcpipboy isbn --operation generate --count 5` (defaults to ISBN-13)
- Verify MCP integration: tool appears in `tools/list` and executes via `tools/call`
- Test with known valid/invalid ISBN numbers and format conversions

Status: Complete - ISBN tool implementation with comprehensive validation, generation, and resource system

---

Implement European Article Number (EAN-13) validation and generation.

11. **EAN-13 Tool**
   - [ ] Create `internal/tools/ean13.go` with flexible EAN-13 functionality
   - [ ] Implement operation parameter: "validate", "generate"
   - [ ] Add EAN-13 validation with checksum algorithm
   - [ ] Add EAN-13 generation with valid checksum and count parameter
   - [ ] Add comprehensive test coverage in `internal/tools/ean13_test.go`
   - [ ] Add CLI command in `cmd/mcpipboy/ean13.go`
   - [ ] Add CLI tests in `cmd/mcpipboy/ean13_test.go`

How to test
- Run `just test` to ensure all EAN-13 tests pass
- Test validation: `mcpipboy ean13 --operation validate --input "1234567890123"`
- Test generation: `mcpipboy ean13 --operation generate --count 5`
- Verify MCP integration: tool appears in `tools/list` and executes via `tools/call`
- Test with known valid/invalid EAN-13 numbers

Status: Pending - EAN-13 tool implementation

---

Implement International Bank Account Number (IBAN) validation and generation.

12. **IBAN Tool**
   - [ ] Create `internal/tools/iban.go` with flexible IBAN functionality
   - [ ] Implement operation parameter: "validate", "generate"
   - [ ] Add IBAN validation with MOD-97 checksum algorithm
   - [ ] Add IBAN generation with optional country code parameter and count parameter
   - [ ] Add comprehensive test coverage in `internal/tools/iban_test.go`
   - [ ] Add CLI command in `cmd/mcpipboy/iban.go`
   - [ ] Add CLI tests in `cmd/mcpipboy/iban_test.go`

How to test
- Run `just test` to ensure all IBAN tests pass
- Test validation: `mcpipboy iban --operation validate --input "GB82WEST12345698765432"`
- Test generation: `mcpipboy iban --operation generate --country-code "GB" --count 3`
- Test generation without country: `mcpipboy iban --operation generate --count 5`
- Verify MCP integration: tool appears in `tools/list` and executes via `tools/call`
- Test with known valid/invalid IBAN numbers and country codes

Status: Pending - IBAN tool implementation

---

Implement Vehicle Identification Number (VIN) validation and generation.

13. **VIN Tool**
   - [ ] Create `internal/tools/vin.go` with flexible VIN functionality
   - [ ] Implement operation parameter: "validate", "generate"
   - [ ] Add VIN validation with checksum algorithm
   - [ ] Add VIN generation with valid checksum and count parameter
   - [ ] Add comprehensive test coverage in `internal/tools/vin_test.go`
   - [ ] Add CLI command in `cmd/mcpipboy/vin.go`
   - [ ] Add CLI tests in `cmd/mcpipboy/vin_test.go`

How to test
- Run `just test` to ensure all VIN tests pass
- Test validation: `mcpipboy vin --operation validate --input "1HGBH41JXMN109186"`
- Test generation: `mcpipboy vin --operation generate --count 5`
- Verify MCP integration: tool appears in `tools/list` and executes via `tools/call`
- Test with known valid/invalid VIN numbers

Status: Pending - VIN tool implementation

---

Update documentation and ensure all new tools are properly integrated.

14. **README Updates**
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
- Unit converter tool (length, weight, temperature, currency, etc.)
