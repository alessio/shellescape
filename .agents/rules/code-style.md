# Rule: Code Style & Standards

## Go Idioms & Formatting
- **Standard Formatting**: Always format code using standard Go tools (`gofmt`, `goimports`, `gci`).
- **Import Ordering**: Maintain three distinct import sections according to `.golangci.yml`:
  1. Standard library packages (e.g. `bufio`, `fmt`, `os`, `strings`)
  2. External dependencies (e.g. `github.com/google/shlex`)
  3. Local module packages (`al.essio.dev/pkg/shellescape`)
- **No Unused Code**: Ensure no unused parameters, variables, or redundant conversions (`unconvert`, `unparam`, `unused`).

## Documentation & Godoc
- Every exported function, type, constant, and variable must have clear, idiomatic godoc comments starting with the symbol name.
- Package comments in `shellescape.go` and `cmd/escargs/escargs.go` must accurately describe functionality and examples.
- Include executable `Example*` tests in `example_test.go` for new or modified public APIs.

## Error Handling
- Never ignore returned errors silently unless explicitly documented (e.g., in utility help printers with `_, _ = fmt.Fprintln(...)`).
- Handle CLI errors gracefully via `log.Fatal` or explicit error exits with informative prefix messages (`escargs: ...`).

## Linter Compliance
- All code must pass `.golangci.yml` linting checks with zero warnings.
- Key linters enforced: `bodyclose`, `copyloopvar`, `dogsled`, `goconst`, `gocritic`, `gosec`, `govet`, `ineffassign`, `misspell`, `prealloc`, `revive`, `staticcheck`, `thelper`, `unconvert`, `unparam`, `unused`, `wsl`.
