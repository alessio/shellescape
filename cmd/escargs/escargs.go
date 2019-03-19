// escargs reads lines from the standard input and prints shell-escaped
// versions. Unlinke xargs, blank lines on the standard input are not
// discarded.
package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/alessio/shellescape"
)

func main() {
	firstScan := true
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if firstScan {
			firstScan = false
		} else {
			fmt.Printf(" ")
		}
		fmt.Printf("%s", shellescape.Quote(scanner.Text()))
	}
}
