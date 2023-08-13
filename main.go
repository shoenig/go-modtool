package main

import (
	"os"

	"github.com/shoenig/modmerge/cli"
)

func main() {
	// subs.Register(cli.)
	// read $1 and $2 (old, new)
	// combine
	// write out

	tool := new(cli.Tool)
	rc := tool.Run()
	os.Exit(rc)
}
