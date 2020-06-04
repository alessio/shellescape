package shellescape_test

import (
	"fmt"
	"strings"

	"github.com/alessio/shellescape"
)

func ExampleQuote() {
	filename := "myfile; rm -rf /"
	prog := "/bin/ls -lh"
	unescapedCommand := strings.Join([]string{prog, filename}, " ")
	escapedCommand := strings.Join([]string{prog, shellescape.Quote(filename)}, " ")

	fmt.Println(unescapedCommand)
	fmt.Println(escapedCommand)
	fmt.Println(strings.Split(unescapedCommand," "))
	fmt.Println(strings.Split(escapedCommand," "))
	// Output:
	// /bin/ls -lh myfile; rm -rf /
	// /bin/ls -lh 'myfile; rm -rf /'
	// [/bin/ls -lh myfile; rm -rf /]
	// [/bin/ls -lh 'myfile; rm -rf /']
}
