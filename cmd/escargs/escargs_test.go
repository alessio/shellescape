package main

import (
	"bufio"
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestEscape(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    string
		wantErr error
	}{
		{
			name:  "quotes each item",
			input: "hello world\nsafe\n",
			want:  "'hello world' safe",
		},
		{
			name:  "no input",
			input: "",
			want:  "",
		},
		{
			name:    "item longer than the scanner buffer",
			input:   strings.Repeat("a", bufio.MaxScanTokenSize) + "\nsafe\n",
			want:    "",
			wantErr: bufio.ErrTooLong,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer

			err := escape(&buf, strings.NewReader(tt.input))
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("escape() error = %v, want %v", err, tt.wantErr)
			}

			if got := buf.String(); got != tt.want {
				t.Errorf("escape() wrote %q, want %q", got, tt.want)
			}
		})
	}
}
