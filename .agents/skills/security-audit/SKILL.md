---
name: security-audit
description: >-
  Use this skill when auditing shell escaping logic, testing against command injection attack vectors, or verifying POSIX shell safety guarantees.
---

# Security Audit Procedures for `shellescape`

This skill guides security audits, invariant checks, and edge-case testing for shell injection vulnerabilities in `shellescape`.

## 1. Security Invariants Checklist

Verify that every one of the following dangerous characters forces single-quote wrapping:

- [ ] Shell metacharacters: `;`, `&`, `|`, `&&`, `||`, `(` , `)`, `{`, `}`
- [ ] Redirection symbols: `<`, `>`, `>>`, `>&`, `<<`
- [ ] Expansions: `$`, `${...}`, `$(...)`, `` `...` ``
- [ ] Wildcards: `*`, `?`, `[`, `]`
- [ ] Whitespace: Spaces, tabs `\t`, newlines `\n`, carriage returns `\r`, form feeds `\f`
- [ ] Quotes: `'` (escaped as `'"'"'`), `"`
- [ ] Escape sequences: `\`

## 2. Adversarial Test Cases

When auditing or adding tests, include these test cases in `shellescape_test.go`:

| Payload Category | Example Input | Expected Output | Rationale |
| :--- | :--- | :--- | :--- |
| Command Injection | `test; rm -rf /` | `'test; rm -rf /'` | Semicolon must be inside single quotes |
| Subshell Execution | `$(cat /etc/passwd)` | `'$(cat /etc/passwd)'` | Subshell syntax must not expand |
| Backticks | `` `whoami` `` | `''`whoami`'` | Backticks must not execute |
| Single Quote Breakout | `'; id; '` | `''"'"'; id; '"'"''` | Quote must not allow closing string context |
| Environment Variable | `$PATH` | `'$PATH'` | Must not expand environment variable |
| Newline Injection | `foo\nbar` | `'foo\nbar'` | Multiline arguments must be safely enclosed |
| Empty String | `""` | `''` | Prevents argument vanishing in shell argv |

## 3. Verifying Escaping with Shell Lexers

Use `github.com/google/shlex` or subshell execution to confirm that the parsed command produces exactly one argument token with the unescaped content.

Example verification snippet in Go:

```go
package main

import (
    "fmt"
    "github.com/google/shlex"
    "al.essio.dev/pkg/shellescape"
)

func verifyPayload(input string) bool {
    escaped := shellescape.Quote(input)
    cmd := "echo " + escaped
    tokens, err := shlex.Split(cmd)
    if err != nil || len(tokens) != 2 {
        return false
    }
    return tokens[1] == input
}
```

## 4. Unicode & Control Character Audit

Verify that `StripUnsafe` and `StripSpaces` function properly:
- `StripUnsafe` removes non-printable runes (`\u0000`, `\u001f`, `\u0081`) while preserving all printable ASCII and Unicode.
- `StripSpaces` removes all whitespace according to `unicode.IsSpace`.
