# Rule: Security Invariants & Shell Escaping Safety

## The POSIX Quoting Model
- **Safe Characters**: Only characters matching `[\w@%+=:,./-]` are considered safe without quoting.
- **Unsafe & Control Characters**: Spaces, tabs, newlines, semicolons, backticks, dollar signs, parentheses, braces, pipes, redirection operators, wildcards (`*`, `?`, `[`), and single quotes MUST be quoted.
- **Single Quote Escaping Pattern**: In POSIX shells (sh, bash, dash, zsh, ksh), a single-quoted string `'...'` treats all characters literally, with no escape sequences allowed inside. Therefore, to represent a single quote `'` within a single-quoted string, `shellescape.Quote` outputs:
  ```go
  "'" + strings.ReplaceAll(s, "'", "'\"'\"'") + "'"
  ```
  This effectively closes the single-quote block, opens a double-quoted block containing a single quote, and immediately resumes the single-quoted block.
- **Empty String**: An empty string `""` must quote to `''` so it is treated as an empty argument rather than disappearing in shell word splitting.

## Unicode & Control Rune Sanitization
- `StripUnsafe(s)`: Uses `unicode.IsPrint(r)` to remove non-printable runes (ASCII control characters, ANSI escape sequences, non-printing unicode code points).
- `StripSpaces(s)`: Uses `unicode.IsSpace(r)` to remove all whitespace characters (spaces, tabs, newlines, form feeds).

## Stream Tokenization
- `ScanTokens`: Handles null-terminated byte streams (`\x00`). Must correctly detect EOF and non-terminated final lines without data loss or memory leaks.

## Zero Runtime Dependencies Rule
- The core package must NEVER depend on external third-party packages. All security and string logic must be self-contained within Go's standard library.
