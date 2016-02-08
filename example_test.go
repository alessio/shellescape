package shellescape_test

import (
	"fmt"
	"strings"

	"github.com/alessio/shellescape"
)

func ExampleQuote() {
	filename := "myfile; rm -rf /"
	prog := "/bin/ls"
	unescapedCommand := strings.Join([]string{prog, filename}, " ")
	escapedCommand := strings.Join([]string{prog, shellescape.Quote(filename)}, " ")

	fmt.Println(unescapedCommand)
	fmt.Println(escapedCommand)
	// Output:
	// /bin/ls myfile; rm -rf /
	// /bin/ls 'myfile; rm -rf /'
}
