// Copyright (c) Seth Hoenig
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"os"

	"cattlecloud.net/go/babycli"
	"github.com/shoenig/go-modtool/commands"
)

func main() {
	args := babycli.Arguments()
	rc := commands.Invoke(args)
	os.Exit(rc)
}
