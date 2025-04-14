// Copyright (c) Seth Hoenig
// SPDX-License-Identifier: MPL-2.0

package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/shoenig/go-modtool/modfile"
)

const (
	exitSuccess = 0
	exitFailure = 1
)

type Tool struct {
	configFile            string // optional config file
	writeFile             bool   // overwrite file(s) in place
	convertPathSeparators bool   // convert path separators to UNIX format
	replaceComment        string // replacement block
	submodulesComment     string // replacement block for submodules
	toolchainComment      string // go toolchain
	excludeComment        string // exclude block
	modFile               string // the go.mod file
}

func (t *Tool) flags() []string {
	flag.BoolVar(&t.writeFile, "w", false, "Write go.mod/go.sum file(s) in place (optional)")
	flag.BoolVar(&t.convertPathSeparators, "p", false, "Convert path separators to UNIX format (optional)")
	flag.StringVar(&t.configFile, "config", "", "Config file (optional)")
	flag.StringVar(&t.replaceComment, "replace-comment", "", "Comment for replace stanza (optional)")
	flag.StringVar(&t.submodulesComment, "submodules-comment", "", "Comment for submodules replace stanza (optional)")
	flag.StringVar(&t.toolchainComment, "toolchain-comment", "", "Comment for go toolchain directive (optional)")
	flag.StringVar(&t.excludeComment, "exclude-comment", "", "Comment for exclude directive (optional)")
	flag.Parse()
	return flag.Args()
}

func (t *Tool) applyConfig() int {
	if t.configFile == "" {
		return 0
	}

	type config struct {
		WriteFile            bool
		ConverPathSeparators bool
		ReplaceComment       string
		SubmodulesComment    string
		ToolchainComment     string
		ExcludeComment       string
	}

	var c config
	_, err := toml.DecodeFile(t.configFile, &c)
	if err != nil {
		fmt.Fprintln(os.Stderr, "crash:", err)
		return exitFailure
	}

	// override default values of args if set in the config file,
	// i.e. the args take precedence

	if !t.writeFile {
		t.writeFile = c.WriteFile
	}

	if !t.convertPathSeparators {
		t.convertPathSeparators = c.ConverPathSeparators
	}

	if t.replaceComment == "" {
		t.replaceComment = c.ReplaceComment
	}

	if t.submodulesComment == "" {
		t.submodulesComment = c.SubmodulesComment
	}

	if t.toolchainComment == "" {
		t.toolchainComment = c.ToolchainComment
	}

	if t.excludeComment == "" {
		t.excludeComment = c.ExcludeComment
	}

	return 0
}

func (t *Tool) Run() int {
	args := t.flags() // initialize
	t.applyConfig()   // read config file if set

	var err error
	switch {
	case len(args) == 0:
		err = errors.New("how did you get here?")
	case len(args) == 1:
		err = errors.New("expects one of 'fmt' or 'merge'")
	case args[0] == "fmt":
		err = t.fmt(args[1:])
	case args[0] == "merge":
		err = t.merge(args[1:])
	default:
		err = errors.New("subcmd must be 'fmt' or 'merge'")
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, "crash:", err)
		return exitFailure
	}

	return exitSuccess
}

func (t *Tool) openMod(file string) (*modfile.Content, error) {
	modFile, err := modfile.Open(file)
	if err != nil {
		return nil, err
	}

	if t.convertPathSeparators {
		for _, replace := range modFile.Replace {
			replace.Old.Path = strings.ReplaceAll(replace.Old.Path, "\\", "/")
			replace.New.Path = strings.ReplaceAll(replace.New.Path, "\\", "/")
		}
	}

	content, err := modfile.Process(modFile)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func (t *Tool) fmt(args []string) error {
	switch len(args) {
	case 0:
		return errors.New("must specify go.mod file to format")
	case 1:
	default:
		return errors.New("must specify only one go.mod file")
	}

	t.modFile = args[0]

	content, err := t.openMod(t.modFile)
	if err != nil {
		return err
	}

	content.Toolchain.Comment = t.toolchainComment
	content.Replace.Comment = t.replaceComment
	content.ReplaceSub.Comment = t.submodulesComment
	content.Exclude.Comment = t.excludeComment
	return t.write(content)
}

func (t *Tool) merge(args []string) error {
	switch len(args) {
	case 0, 1:
		return errors.New("must specify old and new go.mod files to merge")
	case 2:
	default:
		return errors.New("must specify just old and new go.mod files to merge")
	}

	original, err := t.openMod(args[0])
	if err != nil {
		return err
	}

	next, err := t.openMod(args[1])
	if err != nil {
		return err
	}

	content := modfile.Merge(original, next)
	content.Toolchain.Comment = t.toolchainComment
	content.Replace.Comment = t.replaceComment
	content.ReplaceSub.Comment = t.submodulesComment
	content.Exclude.Comment = t.excludeComment
	return t.write(content)
}

func (t *Tool) write(content *modfile.Content) error {
	if t.writeFile {
		f, err := os.OpenFile(t.modFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		if err = content.Write(f); err != nil {
			return err
		}
		if err = f.Sync(); err != nil {
			return err
		}
		if err = f.Close(); err != nil {
			return err
		}
		return nil
	}
	return content.Write(os.Stdout)
}
