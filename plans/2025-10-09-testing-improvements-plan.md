## Testing Improvements Plan (Testable CLI Commands with Proper Output Handling)

Progress legend: [x] Completed, [>] In progress, [ ] Pending

### [x] 0) Planning document and alignment

Plan created to refactor CLI command testing to use proper output streams instead of hardcoded `fmt.Println` calls.

**Current Problem:**
- CLI commands use `fmt.Println()` which always writes to stdout
- Tests either spam output (when calling functions directly) or can't measure coverage (when using `go run`)
- No way to capture and verify command output in tests
- Coverage dropped from 79.5% to 18.2% in cmd/mcpipboy after refactoring to use `go run`

**Proposed Solution:**
- Refactor all `run*()` functions to accept `io.Writer` parameter
- Change all `fmt.Println()` calls to `fmt.Fprintln(writer, ...)`
- In real CLI usage, pass `os.Stdout`
- In tests, pass `&bytes.Buffer{}` to capture output
- Verify exact output in tests
- Restore proper coverage measurement while keeping clean test output

**Benefits:**
1. Tests can capture and verify exact output
2. No output spam during test runs
3. Proper coverage measurement (functions execute in-process)
4. More flexible output handling (stdout, files, buffers, etc.)
5. Idiomatic Go testing pattern
6. Better error message testing

Status: Complete - Plan created

---

### [ ] 1) Refactor Command Signatures

Update all CLI command `run*()` functions to accept `io.Writer` parameter.

1. **Update function signatures**
   - [ ] Add `out io.Writer` parameter to `runEcho()`
   - [ ] Add `out io.Writer` parameter to `runVersion()`
   - [ ] Add `out io.Writer` parameter to `runTime()`
   - [ ] Add `out io.Writer` parameter to `runRandom()`
   - [ ] Add `out io.Writer` parameter to `runUUID()`
   - [ ] Add `out io.Writer` parameter to `runCreditCard()`
   - [ ] Add `out io.Writer` parameter to `runISBN()`
   - [ ] Add `out io.Writer` parameter to `runEAN13()`
   - [ ] Add `out io.Writer` parameter to `runIBAN()`
   - [ ] Add `out io.Writer` parameter to `runIMO()`
   - [ ] Add `out io.Writer` parameter to `runMMSI()`

2. **Update command definitions**
   - [ ] Update all `RunE` fields in command definitions to pass `os.Stdout`
   - [ ] Ensure all commands use `RunE` instead of `Run` for proper error handling

How to test
- Verify all commands compile without errors
- Check that command execution still works via CLI

---

### [ ] 2) Refactor Output Calls

Replace all `fmt.Println()` and `fmt.Printf()` calls with stream-based equivalents.

1. **Replace output calls in echo.go**
   - [ ] Change `fmt.Println()` to `fmt.Fprintln(out, ...)`
   - [ ] Change `fmt.Printf()` to `fmt.Fprintf(out, ...)`

2. **Replace output calls in version.go**
   - [ ] Change `fmt.Println()` to `fmt.Fprintln(out, ...)`

3. **Replace output calls in time.go**
   - [ ] Change `fmt.Println()` to `fmt.Fprintln(out, ...)`
   - [ ] Change `fmt.Printf()` to `fmt.Fprintf(out, ...)`

4. **Replace output calls in random.go**
   - [ ] Change `fmt.Println()` to `fmt.Fprintln(out, ...)`
   - [ ] Change `fmt.Printf()` to `fmt.Fprintf(out, ...)`

5. **Replace output calls in uuid.go**
   - [ ] Change `fmt.Println()` to `fmt.Fprintln(out, ...)`
   - [ ] Change `fmt.Printf()` to `fmt.Fprintf(out, ...)`

6. **Replace output calls in creditcard.go**
   - [ ] Change `fmt.Println()` to `fmt.Fprintln(out, ...)`
   - [ ] Change `fmt.Printf()` to `fmt.Fprintf(out, ...)`

7. **Replace output calls in isbn.go**
   - [ ] Change `fmt.Println()` to `fmt.Fprintln(out, ...)`
   - [ ] Change `fmt.Printf()` to `fmt.Fprintf(out, ...)`

8. **Replace output calls in ean13.go**
   - [ ] Change `fmt.Println()` to `fmt.Fprintln(out, ...)`
   - [ ] Change `fmt.Printf()` to `fmt.Fprintf(out, ...)`

9. **Replace output calls in iban.go**
   - [ ] Change `fmt.Println()` to `fmt.Fprintln(out, ...)`
   - [ ] Change `fmt.Printf()` to `fmt.Fprintf(out, ...)`

10. **Replace output calls in imo.go**
    - [ ] Change `fmt.Println()` to `fmt.Fprintln(out, ...)`
    - [ ] Change `fmt.Printf()` to `fmt.Fprintf(out, ...)`

11. **Replace output calls in mmsi.go**
    - [ ] Change `fmt.Println()` to `fmt.Fprintln(out, ...)`
    - [ ] Change `fmt.Printf()` to `fmt.Fprintf(out, ...)`

How to test
- Run `just build` to verify all commands compile
- Test manual CLI execution to ensure output still works
- Verify no output is lost or changed

---

### [ ] 3) Update Command Definitions

Update all Cobra command definitions to pass `os.Stdout` to the `run*()` functions.

1. **Update command RunE fields**
   - [ ] Update `echoCmd.RunE` to pass `os.Stdout`
   - [ ] Update `versionCmd.RunE` to pass `os.Stdout`
   - [ ] Update `timeCmd.RunE` to pass `os.Stdout`
   - [ ] Update `randomCmd.RunE` to pass `os.Stdout`
   - [ ] Update `uuidCmd.RunE` to pass `os.Stdout`
   - [ ] Update `creditCardCmd.RunE` to pass `os.Stdout`
   - [ ] Update `isbnCmd.RunE` to pass `os.Stdout`
   - [ ] Update `ean13Cmd.RunE` to pass `os.Stdout`
   - [ ] Update `ibanCmd.RunE` to pass `os.Stdout`
   - [ ] Update `imoCmd.RunE` to pass `os.Stdout`
   - [ ] Update `mmsiCmd.RunE` to pass `os.Stdout`

2. **Verify command execution**
   - [ ] Test each command manually via CLI
   - [ ] Ensure output appears correctly on stdout
   - [ ] Verify error handling still works

How to test
- Run `mcpipboy echo "test"` and verify output
- Run various commands to ensure they all produce output
- Test error cases to ensure error messages appear

---

### [ ] 4) Refactor CLI Tests

Update all CLI tests to use `bytes.Buffer` for output capture while keeping integration tests.

**Test Strategy:**
- Keep existing `go run` tests as integration tests (verify end-to-end CLI behavior)
- Add new unit tests that call `run*()` functions directly with buffers (for coverage)
- Integration tests verify the full CLI experience
- Unit tests verify exact output and provide coverage measurement

1. **Refactor echo_test.go**
   - [ ] Keep existing `TestEchoCommand` integration test (uses `go run`)
   - [ ] Add new `TestRunEcho` unit test with `bytes.Buffer`
   - [ ] Call `runEcho()` directly with buffer
   - [ ] Verify exact output content
   - [ ] Test error cases

2. **Refactor version_test.go**
   - [ ] Keep integration test if present
   - [ ] Add unit test with buffer
   - [ ] Call `runVersion()` with buffer
   - [ ] Verify version string in output

3. **Refactor iban_test.go**
   - [ ] Keep existing `TestRunIBAN` as integration test (uses `go run`)
   - [ ] Add new `TestRunIBANUnit` with buffer for coverage
   - [ ] Call `runIBAN()` directly with buffer
   - [ ] Verify exact output for validation
   - [ ] Verify exact output for generation
   - [ ] Test error messages

4. **Refactor uuid_test.go**
   - [ ] Keep integration test
   - [ ] Add unit test with buffer
   - [ ] Call `runUUID()` with buffer
   - [ ] Verify UUID format in output
   - [ ] Test all UUID versions (v1, v4, v5, v7)

5. **Refactor creditcard_test.go**
   - [ ] Keep integration test
   - [ ] Add unit test with buffer
   - [ ] Call `runCreditCard()` with buffer
   - [ ] Verify card number output
   - [ ] Test validation output format

6. **Refactor isbn_test.go**
   - [ ] Keep integration test
   - [ ] Add unit test with buffer
   - [ ] Call `runISBN()` with buffer
   - [ ] Verify ISBN output format
   - [ ] Test both ISBN-10 and ISBN-13

7. **Refactor ean13_test.go**
   - [ ] Keep integration test
   - [ ] Add unit test with buffer
   - [ ] Call `runEAN13()` with buffer
   - [ ] Verify EAN-13 output format

8. **Refactor imo_test.go**
   - [ ] Keep integration test
   - [ ] Add unit test with buffer
   - [ ] Call `runIMO()` with buffer
   - [ ] Verify IMO number output

9. **Refactor mmsi_test.go**
   - [ ] Keep integration test
   - [ ] Add unit test with buffer
   - [ ] Call `runMMSI()` with buffer
   - [ ] Verify MMSI output format

10. **Refactor time_test.go**
    - [ ] Keep integration test
    - [ ] Add unit test with buffer
    - [ ] Call `runTime()` with buffer
    - [ ] Verify time format output
    - [ ] Test all time operations

11. **Refactor random_test.go**
    - [ ] Keep integration test
    - [ ] Add unit test with buffer
    - [ ] Call `runRandom()` with buffer
    - [ ] Verify random value output format

How to test
- Run `just test` to ensure all tests pass
- Verify no output spam during test runs
- Check exact output content is validated in tests

---

### [ ] 5) Verify Coverage Improvement

Run coverage tests and verify we've restored proper coverage measurement.

1. **Run coverage tests**
   - [ ] Run `just test-coverage` 
   - [ ] Verify `cmd/mcpipboy` coverage is back above 75%
   - [ ] Verify `internal/tools` coverage remains above 80%
   - [ ] Check total coverage is above 80%

2. **Verify test quality**
   - [ ] All tests pass without errors
   - [ ] No output spam during test runs
   - [ ] Tests verify exact output content
   - [ ] Error cases are properly tested

3. **Generate coverage report**
   - [ ] Run `just coverage-html`
   - [ ] Review coverage.html for any gaps
   - [ ] Add tests for any uncovered critical paths

How to test
- Run `just test-coverage` and check percentages
- Run `just test` and verify clean output
- Review coverage.html for completeness
- Verify all success criteria are met

Status: Pending - Testing improvements

---

## Success Criteria

- [ ] All `run*()` functions accept `io.Writer` parameter
- [ ] All commands pass `os.Stdout` when executing normally
- [ ] All tests use `bytes.Buffer` to capture output
- [ ] All tests verify exact output content
- [ ] No output spam during test runs
- [ ] cmd/mcpipboy coverage restored above 75%
- [ ] internal/tools coverage remains above 80%
- [ ] Total coverage above 80%
- [ ] All tests pass with `just test`
- [ ] Test coverage works with `just test-coverage`
- [ ] No regressions in CLI functionality

## Decisions

**Why io.Writer instead of returning strings?**
- More flexible - can write to stdout, files, buffers, etc.
- Idiomatic Go pattern for output handling
- Allows streaming output for large datasets
- Better separation of concerns (logic vs output)

**Why refactor now?**
- Current tests spam output making them hard to read
- Coverage measurement is broken with external process execution
- Proper output handling makes code more maintainable
- Sets good pattern for future commands

**What about the serve command?**
- Serve command already uses stdin/stdout properly
- Server tests will remain as integration tests
- Server coverage issues are due to pipe management, not fixable with this approach

## Notes

- This refactoring improves code quality and testability
- It follows Go best practices for testable CLI applications
- The io.Writer pattern is used throughout the Go standard library
- Tests will be more comprehensive and maintainable
- Coverage measurement will be accurate and meaningful

