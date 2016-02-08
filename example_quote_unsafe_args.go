package shellescape

import (
	"fmt"
	"strings"
)

func ExampleUnsageArgs() {
	filename := "myfile; rm -rf /"
	prog := "/bin/ls"
	unescaped_command := strings.Join([]string{prog, filename}, " ")
	escaped_command := strings.Join([]string{prog, Quote(filename)}, " ")

	fmt.Printf("Unescaped command: %s\n", unescaped_command)
	fmt.Printf("Escaped command: %s\n", escaped_command)
}
