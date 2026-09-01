---
name: shellescape-dev
description: >-
  Use this skill when developing, testing, building, benchmarking, or running linters for the shellescape Go library and escargs CLI tool.
---

# Development Workflow for `shellescape`

This skill provides step-by-step procedures for the local development cycle of `shellescape` and `escargs`.

## 1. Environment Setup

Ensure Go 1.18 or newer is accessible in the environment:

```bash
export PATH=$PATH:/usr/local/go/bin:/opt/homebrew/bin:$HOME/go/bin
go version
```

## 2. Running Tests & Race Detection

Execute all unit and example tests with race detection:

```bash
go test -v -race ./...
```

To run coverage analysis:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

## 3. Running Benchmarks

Run benchmarks to verify there are no performance regressions or unexpected memory allocations:

```bash
go test -v -bench=. -benchmem ./...
```

Compare memory allocations:
- `Quote` on simple strings should be `0 allocs/op`.
- `Quote` on complex strings should allocate minimally.

## 4. Running Code Quality & Linters

Run static analysis using `golangci-lint`:

```bash
golangci-lint run ./...
```

If `golangci-lint` is not installed globally, you can run standard Go vetting:

```bash
go vet ./...
```

## 5. Building the CLI Utility (`escargs`)

Build using the Makefile or Go toolchain:

```bash
# Build via Makefile
make build

# Or directly build binary
go build -v -ldflags="-X 'main.version=dev'" -o escargs ./cmd/escargs
```

Test the binary interactively or via stdin:

```bash
echo -e "foo\nbar with spaces\n; rm -rf /" | ./escargs
# Expected output: foo 'bar with spaces' '; rm -rf /'
```

## 6. Cleaning Build Artifacts

```bash
make clean
```
