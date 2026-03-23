package commands

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"cattlecloud.net/go/babycli"
	"github.com/BurntSushi/toml"
	"github.com/shoenig/go-modtool/modfile"
)

type Tool struct {
	configFile            string // optional config file
	writeFile             bool   // overwrite file(s) in place
	convertPathSeparators bool   // convert path separators to UNIX format
	toolComment           string // tool block
	replaceComment        string // replacement block
	submodulesComment     string // replacement block for submodules
	toolchainComment      string // go toolchain
	excludeComment        string // exclude block
	arguments             []string
	outFile               string
}

func NewTool(c *babycli.Component) *Tool {
	return &Tool{
		configFile:            c.GetString("config"),
		writeFile:             c.GetBool("write-in-place"),
		convertPathSeparators: c.GetBool("unix-paths"),
		toolComment:           c.GetString("tool-comment"),
		replaceComment:        c.GetString("replace-comment"),
		submodulesComment:     c.GetString("submodules-comment"),
		toolchainComment:      c.GetString("toolchain-comment"),
		excludeComment:        c.GetString("exclude-comment"),
		arguments:             c.Arguments(),
	}
}

func (t *Tool) applyConfig() error {
	if t.configFile == "" {
		return nil
	}

	type config struct {
		WriteFile            bool   `toml:"WriteFile"`
		ConverPathSeparators bool   `toml:"ConverPathSeparators"` // TODO: fix typo
		ToolComment          string `toml:"ToolComment"`
		ReplaceComment       string `toml:"ReplaceComment"`
		SubmodulesComment    string `toml:"SubmodulesComment"`
		ToolchainComment     string `toml:"ToolchainComment"`
		ExcludeComment       string `toml:"ExcludeComment"`
	}

	var c config
	if _, err := toml.DecodeFile(t.configFile, &c); err != nil {
		return fmt.Errorf("unable to decode config file: %w", err)
	}

	// override the config file values with cli argument values, which have a higher precedence

	if !t.writeFile {
		t.writeFile = c.WriteFile
	}

	if !t.convertPathSeparators {
		t.convertPathSeparators = c.ConverPathSeparators
	}

	if t.toolComment == "" {
		t.toolComment = c.ToolComment
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

	return nil
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

	// set the go.mod file
	t.outFile = args[0]

	content, err := t.openMod(t.outFile)
	if err != nil {
		return err
	}

	content.Tool.Comment = t.toolComment
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
	t.outFile = args[0]

	next, err := t.openMod(args[1])
	if err != nil {
		return err
	}

	content := modfile.Merge(original, next)
	content.Tool.Comment = t.toolComment
	content.Toolchain.Comment = t.toolchainComment
	content.Replace.Comment = t.replaceComment
	content.ReplaceSub.Comment = t.submodulesComment
	content.Exclude.Comment = t.excludeComment
	return t.write(content)
}

func (t *Tool) write(content *modfile.Content) error {
	if t.writeFile {
		f, err := os.OpenFile(t.outFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
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
