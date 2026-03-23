package commands

import (
	"fmt"
	"os"

	"cattlecloud.net/go/babycli"
)

var (
	// Version is set at release time to the latest git tag.
	Version = "development"
)

const (
	name        = "go-modtool"
	description = `An opinionated tool for formatting and merging Go's go.mod and go.sum files.`
)

func Invoke(args []string) babycli.Code {
	return babycli.New(&babycli.Configuration{
		Arguments: args,
		Version:   Version,
		Globals: babycli.Flags{
			{
				Type: babycli.StringFlag,
				Long: "config",
				Help: "Specify TOML configuration file",
			},
			{
				Type:  babycli.BooleanFlag,
				Long:  "write-in-place",
				Short: "w",
				Help:  "Write go.mod/go.sum file(s) in place",
			},
			{
				Type:  babycli.BooleanFlag,
				Long:  "unix-paths",
				Short: "p",
				Help:  "Convert path separators to UNIX format",
			},
			{
				Type: babycli.StringFlag,
				Long: "tool-comment",
				Help: "Comment for tool stanza",
			},
			{
				Type: babycli.StringFlag,
				Long: "replace-comment",
				Help: "Comment for replace stanza",
			},
			{
				Type: babycli.StringFlag,
				Long: "submodules-comment",
				Help: "Comment for submodules replace stanza",
			},
			{
				Type: babycli.StringFlag,
				Long: "toolchain-comment",
				Help: "Comment for go toolchain directive",
			},
			{
				Type: babycli.StringFlag,
				Long: "exclude-comment",
				Help: "Comment for exclude directive",
			},
		},
		Top: &babycli.Component{
			Name:        name,
			Description: description,
			Components: babycli.Components{
				newFmtCommand(),
				newMergeCommand(),
			},
		},
	}).Run()
}

func crash(err error) babycli.Code {
	fmt.Fprintln(os.Stderr, "error:", err)
	fmt.Fprintln(os.Stderr, "error: use -h for help")
	return babycli.Failure
}
