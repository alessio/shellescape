// escargs reads lines from the standard input and prints shell-escaped
// versions. Unlinke xargs, blank lines on the standard input are not
// discarded.
package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"

	"github.com/alessio/shellescape"
)

var (
	flagDiscardBlankLines = flag.Bool("D", false, "ignore blank lines on the standard input.")
	flagNull              = flag.Bool("0", false, "input items are terminated by a null character instead of by new line.")
)

func main() {
	flag.Parse()

	firstScan := true
	scanner := bufio.NewScanner(os.Stdin)

	if *flagNull {
		scanner.Split(splitNullTerminatedItems)
	}

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

func splitNullTerminatedItems(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// Return nothing if at end of file and no data passed.
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	// Find the index of the input of a null character.
	if i := bytes.IndexByte(data, '\x00'); i >= 0 {
		return i + 1, data[0:i], nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}

	// Request more data.
	return 0, nil, nil
}
