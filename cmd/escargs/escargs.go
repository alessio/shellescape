// escargs reads lines from the standard input and prints shell-escaped
// versions. Unlinke xargs, blank lines on the standard input are not
// discarded.
package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/alessio/shellescape"
)

var (
	discardBlankLines bool
	nullSeparator     bool
	argFile           string
	helpMode          bool
	versionMode       bool
)

var version = "UNRELEASED"

func init() {
	flag.BoolVar(&discardBlankLines, "D", false, "ignore blank lines on the input stream.")
	flag.BoolVar(&nullSeparator, "0", false, "input items are terminated by a null character instead of by new line.")
	flag.StringVar(&argFile, "a", "", "read arguments from file, not standard input.")
	flag.BoolVar(&helpMode, "h", false, "display this help and exit.")
	flag.BoolVar(&versionMode, "V", false, "output version information and exit.")
	flag.Usage = usage
	flag.ErrHelp = nil
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("escargs: ")
	log.SetOutput(os.Stderr)
	flag.Parse()

	if helpMode {
		usage()
		return
	}

	if versionMode {
		outputVersion()
		return
	}

	firstScan := true
	scanner := bufio.NewScanner(os.Stdin)

	if argFile != "" {
		f, err := os.Open(argFile)
		if err != nil {
			log.Fatal(err)
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

func usage() {
	usageString := `Usage: escargs [-0ad]
Escape arbitrary strings for safe use as command line arguments.
		
Options:`
	_, _ = fmt.Fprintln(os.Stderr, usageString)

	flag.PrintDefaults()
}

func outputVersion() {
	fmt.Fprintf(os.Stderr, "escargs version %s\n", version)
	fmt.Fprintln(os.Stderr, "Copyright (C) 2020-2023 Alessio Treglia <alessio@debian.org>")
}
