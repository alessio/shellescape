// The shellescape package provides the function Quote to escape arbitrary
// strings for a safe use as command line arguments in the most common
// POSIX shells.

package shellescape

import (
	"regexp"
	"strings"
)

var pattern *regexp.Regexp

func init() {
	pattern = regexp.MustCompilePOSIX("[a-zA-Z0-9_^@%+=:,./-]")
}

// Quote replaces all occurrences in the string s of the
// single quote character with an escaped version.
func Quote(s string) string {
	if len(s) == 0 {
		return "''"
	}
	if pattern.MatchString(s) {
		return "'" + strings.Replace(s, "'", "'\"'\"'", -1) + "'"
	}

	return s
}
