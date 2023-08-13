package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/shoenig/modmerge/modfile"
)

const (
	exitSuccess = 0
	exitFailure = 1
)

type Tool struct {
	repsComment string // replacement block
	subsComment string // replacement block for submodules
}

func (t *Tool) flags() []string {
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
		fmt.Println("failure:", err)
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

	original, err := modfile.Open(args[0])
	if err != nil {
		return err
	}

	original.Replace.Comment = t.repsComment
	original.ReplaceSub.Comment = t.subsComment

	original.Write(os.Stdout)

	return nil
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

	fmt.Println("original", original)
	fmt.Println("next", next)

	return nil
}

func (t *Tool) write(content *modfile.Content) error {
	return nil
}
