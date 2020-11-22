// escargs reads lines from the standard input and prints shell-escaped
// versions. Unlinke xargs, blank lines on the standard input are not
// discarded.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/alessio/shellescape"
)

var (
	flagDiscardBlankLines = flag.Bool("D", false, "ignore blank lines on the standard input.")
)

func main() {
	flag.Parse()

	firstScan := true
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()
		if *flagDiscardBlankLines && len(line) == 0 {
			continue
		}

		if firstScan {
			firstScan = false
		} else {
			fmt.Printf(" ")
		}

		fmt.Printf("%s", shellescape.Quote(line))
	}
}
