package shellescape

import (
	"fmt"
	"testing"
)

func assertEqual(t *testing.T, s, expected string) {
	if s != expected {
		t.Error(fmt.Sprintf("%q (expected: %q)", s, expected))
	}
}

func TestEmptyString(t *testing.T) {
	s := Quote("")
	expected := "''"
	assertEqual(t, s, expected)
}

func TestDoubleQuotedString(t *testing.T) {
	s := Quote(`"double quoted"`)
	expected := `'"double quoted"'`
	assertEqual(t, s, expected)
}

func TestSingleQuotedString(t *testing.T) {
	s := Quote(`'single quoted'`)
	expected := `''"'"'single quoted'"'"''`
	assertEqual(t, s, expected)
}

func TestUnquotedString(t *testing.T) {
	s := Quote(`no quotes`)
	expected := `'no quotes'`
	assertEqual(t, s, expected)
}
