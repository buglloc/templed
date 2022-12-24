package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/buglloc/templed/internal/commands"
)

func fatal(err error) {
	_, _ = fmt.Fprintf(os.Stderr, "templed: %v\n", err)
	os.Exit(1)
}

func main() {
	runtime.GOMAXPROCS(1)

	if err := commands.Execute(); err != nil {
		fatal(err)
	}
}
