# Project Memory & Agent Guidelines: `shellescape`

Welcome to `shellescape` (`al.essio.dev/pkg/shellescape`). This repository provides a high-reliability, zero-dependency Go library and a CLI tool (`escargs`) to safely escape arbitrary strings for POSIX shell command-line execution and protect systems against shell injection vulnerabilities.

---

## 1. Project Overview & Architecture

### Core Components
1. **Core Library (`shellescape.go`)**:
   - `Quote(s string) string`: Escapes arbitrary strings using POSIX single-quote wrapping (`'...'`) and single-quote substitution (`'\''` via `'\"'\"'`).
   - `QuoteCommand(args []string) string`: Quotes a slice of string arguments and joins them with spaces.
   - `StripUnsafe(s string) string`: Removes non-printable characters/control runes using `unicode.IsPrint`.
   - `StripSpaces(s string) string`: Strips whitespace runes using `unicode.IsSpace`.
   - `ScanTokens(data []byte, atEOF bool) (int, []byte, error)`: Custom `bufio.Scanner` split function for null-delimited (`\x00`) data streams.

2. **CLI Utility (`cmd/escargs/escargs.go`)**:
   - Command-line tool `escargs` that reads tokens from standard input or a specified file (`-a`) and prints shell-escaped tokens separated by spaces.
   - Flags:
     - `-0`: Null-character (`\x00`) input item separator.
     - `-D`: Ignore/discard blank lines.
     - `-a <file>`: Read input from file instead of stdin.
     - `-h`: Display help message.
     - `-V`: Print version and copyright information.

3. **Module & Dependencies**:
   - Module Path: `al.essio.dev/pkg/shellescape`
   - Minimum Go Version: `go 1.18`
   - **Zero Runtime Dependencies**: The core package and CLI rely purely on Go standard library packages (`bytes`, `flag`, `fmt`, `log`, `os`, `regexp`, `strings`, `unicode`, `bufio`).
   - Test Dependency: `github.com/google/shlex` is used strictly in `example_test.go` for validation examples.

---

## 2. Invariants & Security Guarantees

- **POSIX Escaping Safety**: When escaping a string, every character must be enclosed in single quotes unless it matches safe characters `[\w@%+=:,./-]`. Any embedded single quote `'` is escaped as `'"'"'` (closing the single quote, opening double quotes with single quote inside, and reopening single quote).
- **Empty String Escaping**: `Quote("")` must always return `''`.
- **Zero Allocations & High Performance**: Performance is critical. Keep allocations minimal and benchmark all string manipulation changes.
- **Zero External Dependencies**: Never introduce external third-party dependencies to `shellescape.go` or `cmd/escargs/`.

---

## 3. Directory Layout

```text
.
├── .agents/                 # AI Agent instructions, roles, rules, and skills
│   ├── roles/               # Persona and task-specific role definitions
│   │   ├── cli-maintainer.md
│   │   ├── core-developer.md
│   │   ├── release-engineer.md
│   │   └── security-auditor.md
│   ├── rules/               # Contextual workspace rules
│   │   ├── code-style.md
│   │   ├── release-management.md
│   │   ├── security-invariants.md
│   │   └── testing-benchmarks.md
│   └── skills/              # Workflow action runbooks
│       ├── release-workflow/
│       │   └── SKILL.md
│       ├── security-audit/
│       │   └── SKILL.md
│       └── shellescape-dev/
│           └── SKILL.md
├── .github/workflows/       # CI/CD workflows (build, lint, codeql, release)
├── cmd/
│   └── escargs/             # CLI application source
│       └── escargs.go
├── example_test.go          # Godoc testable examples
├── go.mod                   # Go module definition
├── go.sum                   # Go checksums
├── Makefile                 # Make build/install automation
├── shellescape.go           # Core package logic
└── shellescape_test.go      # Unit tests and benchmarks
```

---

## 4. Agent Roles & Responsibilities

When taking on tasks in this repository, identify which role best fits the scope:

1. **[Core Developer](file:///.agents/roles/core-developer.md)**:
   - Maintains `shellescape.go`, string performance, API stability, and zero-dependency guarantees.
2. **[CLI Maintainer](file:///.agents/roles/cli-maintainer.md)**:
   - Maintains `cmd/escargs`, flag parsing, stream processing, UNIX pipe behavior, and error reporting.
3. **[Security Auditor](file:///.agents/roles/security-auditor.md)**:
   - Analyzes POSIX shell escaping mechanics, fuzzing tests, injection vectors, and Unicode/control rune handling.
4. **[Release Engineer](file:///.agents/roles/release-engineer.md)**:
   - Manages semantic versioning, `.goreleaser.yml`, GitHub Actions pipelines, cross-compilation builds, and tag management.

---

## 5. Development & Maintenance Commands

### Building
```bash
# Build package and escargs binary
make build
# Or directly via Go
go build ./...
go build ./cmd/escargs
```

### Testing
```bash
# Run unit and example tests with race detector
go test -v -race ./...

# Run test coverage
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

### Benchmarks
```bash
# Run all benchmarks with memory allocation metrics
go test -v -bench=. -benchmem ./...
```

### Linting
```bash
# Run golangci-lint
golangci-lint run ./...
```

### Clean
```bash
make clean
```
