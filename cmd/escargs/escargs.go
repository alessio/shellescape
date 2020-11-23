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
	discardBlankLines bool
	nullSeparator     bool
	argFile           string
)

func main() {
	flag.BoolVar(&discardBlankLines, "D", false, "ignore blank lines on the input stream.")
	flag.BoolVar(&nullSeparator, "0", false, "input items are terminated by a null character instead of by new line.")
	flag.StringVar(&argFile, "a", "", "read arguments from file, not standard input.")
	flag.Parse()

	firstScan := true
	scanner := bufio.NewScanner(os.Stdin)

	if argFile != "" {
		f, err := os.Open(argFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "escargs: %v\n", err)
			os.Exit(1)
		}

		scanner = bufio.NewScanner(f)
	}

	if nullSeparator {
		scanner.Split(splitNullTerminatedItems)
	}

	for scanner.Scan() {
		line := scanner.Text()
		if discardBlankLines && len(line) == 0 {
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
