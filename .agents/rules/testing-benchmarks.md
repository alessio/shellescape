# Rule: Testing & Benchmarking Standards

## Testing Guidelines
- **Table-Driven Tests**: Use structured table-driven tests with descriptive subtest names (`t.Run(tt.name, func(t *testing.T) { ... })`).
- **Parallel Tests**: Enable `t.Parallel()` for subtests where tests are independent and thread-safe.
- **Race Detection**: Always run tests with `-race` flag enabled:
  ```bash
  go test -race -v ./...
  ```
- **Example Tests**: Maintain runnable `Example*` functions in `example_test.go` with `// Output:` comments. These are tested automatically by `go test` and rendered directly into pkg.go.dev documentation.
- **Coverage**: Strive for 100% test coverage on `shellescape.go`.

## Benchmarking Guidelines
- All core functions (`Quote`, `QuoteCommand`, `StripUnsafe`, `StripSpaces`, `ScanTokens`) must have matching benchmarks in `shellescape_test.go`.
- Include benchmarks for both typical inputs and large strings/slices (e.g., 1000 items/repeats) to detect performance regressions.
- Execute benchmarks with memory allocation metrics:
  ```bash
  go test -bench=. -benchmem ./...
  ```
- Any optimization PR must demonstrate non-increasing memory allocations (`B/op` and `allocs/op`).
