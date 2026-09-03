// escargs reads lines from the standard input and prints shell-escaped
// versions. Unlike xargs, blank lines on the standard input are not
// discarded.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"al.essio.dev/pkg/shellescape"
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

	var input io.Reader = os.Stdin

	if argFile != "" {
		f, err := os.Open(argFile)
		if err != nil {
			log.Fatal(err)
		}

		input = f
	}

	if err := escape(os.Stdout, input); err != nil {
		log.Fatal(err)
	}
}

// escape writes the shell-escaped items read from r to w, separated by
// spaces.
func escape(w io.Writer, r io.Reader) error {
	firstScan := true
	scanner := bufio.NewScanner(r)

	if nullSeparator {
		scanner.Split(shellescape.ScanTokens)
	}

	for scanner.Scan() {
		line := scanner.Text()
		if discardBlankLines && len(line) == 0 {
			continue
		}

		if firstScan {
			firstScan = false
		} else {
			_, _ = fmt.Fprint(w, " ")
		}

		_, _ = fmt.Fprint(w, shellescape.Quote(line))
	}

	// A scan error stops the loop early and truncates the output.
	return scanner.Err()
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
	fmt.Fprintln(os.Stderr, "Copyright (C) 2020-2024 Alessio Treglia <alessio@debian.org>")
}
