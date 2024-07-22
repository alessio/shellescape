package shellescape_test

import (
	"fmt"
	"strings"

	"al.essio.dev/pkg/shellescape"
	"github.com/google/shlex"
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

func ExampleQuoteCommand_simple() {
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

func ExampleQuoteCommand() {
	filename := "myfile; rm -rf /"
	unsafe := fmt.Sprintf("ls -l %s", filename)
	command := fmt.Sprintf("ls -l %s", shellescape.Quote(filename))
	splitCommand, _ := shlex.Split(command)

	fmt.Println("unsafe:", unsafe)
	fmt.Println("command:", command)
	fmt.Println("splitCommand:", splitCommand)

	remoteCommandUnsafe := fmt.Sprintf("ssh host.domain %s", command)
	remoteCommand := fmt.Sprintf("ssh host.domain %s", shellescape.Quote(command))
	splitRemoteCommand, _ := shlex.Split(remoteCommand)

	fmt.Println("remoteCommandUnsafe:", remoteCommandUnsafe)
	fmt.Println("remoteCommand:", remoteCommand)
	fmt.Println("splitRemoteCommand:", splitRemoteCommand)

	lastSplit, _ := shlex.Split(splitRemoteCommand[2])
	fmt.Println("lastSplit[0]:", lastSplit[0])
	fmt.Println("lastSplit[1]:", lastSplit[1])
	fmt.Println("lastSplit[2]:", lastSplit[2])

	// unsafe: ls -l myfile; rm -rf /
	// command: ls -l 'myfile; rm -rf /'
	// splitCommand: [ls -l myfile; rm -rf /]
	// remoteCommandUnsafe: ssh host.domain ls -l 'myfile; rm -rf /'
	// remoteCommand: ssh host.domain 'ls -l '"'"'myfile; rm -rf /'"'"''
	// splitRemoteCommand: [ssh host.domain ls -l 'myfile; rm -rf /']
	// lastSplit[0]: ls
	// lastSplit[1]: -l
	// lastSplit[2]: myfile; rm -rf /
}

func ExampleStripUnsafe() {
	safeString := `"printable!" #$%^characters '' 12321312"`
	unsafeString := "these runes shall be removed: \u0000\u0081\u001f"

	fmt.Println("safe:", shellescape.StripUnsafe(safeString))
	fmt.Println("unsafe:", shellescape.StripUnsafe(unsafeString))
	// Output:
	// safe: "printable!" #$%^characters '' 12321312"
	// unsafe: these runes shall be removed:
}
