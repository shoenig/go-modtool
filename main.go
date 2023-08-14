package main

import (
	"os"

	"github.com/shoenig/go-modtool/cli"
)

func main() {
	tool := new(cli.Tool)
	rc := tool.Run()
	os.Exit(rc)
}
