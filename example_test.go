package shellescape_test

import (
	"fmt"
	"strings"

	"github.com/alessio/shellescape"
)

func ExampleQuote() {
	filename := "myfile; rm -rf /"
	prog := "/bin/ls -lh"
	unsafe := strings.Join([]string{prog, filename}, " ")
	safe := strings.Join([]string{prog, shellescape.Quote(filename)}, " ")

	fmt.Println("unsafe:", unsafe)
	fmt.Println("safe:", safe)

	for i, part := range strings.Split(unsafe, " ") {
		fmt.Printf("unsafe[%d] = %s\n", i, part)
	}

	for i, part := range strings.Split(safe, " ") {
		fmt.Printf("safe[%d] = %s\n", i, part)
	}
	// Output:
	// unsafe: /bin/ls -lh myfile; rm -rf /
	// safe: /bin/ls -lh 'myfile; rm -rf /'
	// unsafe[0] = /bin/ls
	// unsafe[1] = -lh
	// unsafe[2] = myfile;
	// unsafe[3] = rm
	// unsafe[4] = -rf
	// unsafe[5] = /
	// safe[0] = /bin/ls
	// safe[1] = -lh
	// safe[2] = 'myfile;
	// safe[3] = rm
	// safe[4] = -rf
	// safe[5] = /'
}

func ExampleQuoteCommand() {
	filename := "filename with space"
	prog := "/usr/bin/ls"
	args := "-lh"

	unsafe := strings.Join([]string{prog, args, filename}, " ")
	safe := strings.Join([]string{prog, shellescape.QuoteCommand([]string{args, filename})}, " ")

	fmt.Println("unsafe:", unsafe)
	fmt.Println("safe:", safe)
	// Output:
	// unsafe: /usr/bin/ls -lh filename with space
	// safe: /usr/bin/ls -lh 'filename with space'
}
