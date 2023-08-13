package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/shoenig/modtool/modfile"
)

const (
	exitSuccess = 0
	exitFailure = 1
)

type Tool struct {
	writeFile   bool   // overwrite file(s) in place
	repsComment string // replacement block
	subsComment string // replacement block for submodules
	modFile     string // the go.mod file
}

func (t *Tool) flags() []string {
	flag.BoolVar(&t.writeFile, "w", false, "Write go.mod/go.sum file(s) in place")
	flag.StringVar(&t.repsComment, "replace-comment", "", "Comment for replace stanza")
	flag.StringVar(&t.subsComment, "subs-comment", "", "Comment for submodules replace stanza")
	flag.Parse()
	return flag.Args()
}

func (t *Tool) Run() int {
	args := t.flags() // initialize

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

func (t *Tool) fmt(args []string) error {
	switch len(args) {
	case 0:
		return errors.New("must specify go.mod file to format")
	case 1:
	default:
		return errors.New("must specify only one go.mod file")
	}

	t.modFile = args[0]
	content, err := modfile.Open(t.modFile)
	if err != nil {
		return err
	}

	content.Replace.Comment = t.repsComment
	content.ReplaceSub.Comment = t.subsComment
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

	original, err := modfile.Open(args[0])
	if err != nil {
		return err
	}

	next, err := modfile.Open(args[1])
	if err != nil {
		return err
	}

	content := modfile.Merge(original, next)
	content.Replace.Comment = t.repsComment
	content.ReplaceSub.Comment = t.subsComment
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
