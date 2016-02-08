# shellescape
Escape arbitrary strings for safe use as command line arguments.

## Contents of the package

This package provides the `shellescape.Quote()` function that returns a
shell-escaped copy of a string.

This work was inspired by the Python original package [shellescape] 
(https://pypi.python.org/pypi/shellescape).

## Example

```go
package main

import (
        "github.com/alessio/shellescape"
        "fmt"
)

func main() {
        filename := "somefile; rm -rf ~"
        fmt.Printf("ls -l %s\n", shellescape.Quote(filename))
}
```

Run it:

```shell
$ go run bin/program.go
ls -l 'somefile; rm -rf ~'
```
