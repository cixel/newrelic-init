package main

import (
	"flag"
	"fmt"
	"os"
)

func usage() {
	_, _ = fmt.Fprintln(os.Stderr, `Usage: instrument [OPTION]... SOURCE

TODO: description

SOURCE must be valid as a package path.

Flags:`)
	flag.PrintDefaults()
}
