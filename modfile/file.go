package modfile

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/shoenig/semantic"
	modpkg "golang.org/x/mod/modfile"
)

var (
	zero = semantic.New(0, 0, 0)
)

type Dependency struct {
	Name    string
	Version semantic.Tag
}

func (d Dependency) String() string {
	if d.Version.Equal(zero) {
		return d.Name
	}
	return fmt.Sprintf("%s %s", d.Name, d.Version)
}

type ReplaceStanza struct {
	Comment      string
	Replacements []Replacement
}

func (rs *ReplaceStanza) add(r Replacement) {
	rs.Replacements = append(rs.Replacements, r)
}

func (rs *ReplaceStanza) empty() bool {
	return len(rs.Replacements) == 0
}

type Replacement struct {
	Orig Dependency
	Next Dependency
}

func (r Replacement) String() string {
	return fmt.Sprintf("%s => %s", r.Orig, r.Next)
}

type RequireStanza struct {
	Comment      []string
	Dependencies []Dependency
}

func (ds *RequireStanza) empty() bool {
	return len(ds.Dependencies) == 0
}

func (ds *RequireStanza) add(d Dependency) {
	ds.Dependencies = append(ds.Dependencies, d)
}

type Content struct {
	Module     string
	Go         string
	Toolchain  string
	Direct     RequireStanza
	Indirect   RequireStanza
	Replace    ReplaceStanza
	ReplaceSub ReplaceStanza // sub modules, e.g. "=> ./api"

	// Exclude   []Dependency
	// Retract   []semantic.Tag
}

func (c *Content) String() string {
	return "todo"
}

func (c *Content) Write(w io.Writer) error {
	var err error
	const indent = "    "
	write := func(parts ...string) {
		if err == nil {
			for _, part := range parts {
				_, err = io.WriteString(w, part)
				if err != nil {
					return
				}
			}
		}
	}

	write("module", " ", c.Module, "\n", "\n")
	write("go", " ", c.Go, "\n", "\n")

	if c.Toolchain != "" {
		write("toolchain", " ", c.Toolchain, "\n", "\n")
	}

	if !c.Replace.empty() {
		if c.Replace.Comment != "" {
			write("// ", c.Replace.Comment, "\n")
		}
		write("replace (", "\n")
		for _, d := range c.Replace.Replacements {
			write(indent, d.String(), "\n")
		}
		write(")", "\n", "\n")
	}

	if !c.ReplaceSub.empty() {
		if c.ReplaceSub.Comment != "" {
			write("// ", c.ReplaceSub.Comment, "\n")
		}
		write("replace (", "\n")
		for _, d := range c.ReplaceSub.Replacements {
			write(indent, d.String(), "\n")
		}
		write(")", "\n", "\n")
	}

	if !c.Direct.empty() {
		write("require (", "\n")
		for _, d := range c.Direct.Dependencies {
			write(indent, d.String(), "\n")
		}
		write(")", "\n", "\n")
	}

	if !c.Indirect.empty() {
		write("require (", "\n")
		for _, d := range c.Indirect.Dependencies {
			write(indent, d.String(), " // indirect", "\n")
		}
		write(")", "\n", "\n")
	}

	return err
}

func Open(path string) (*Content, error) {
	b, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = b.Close() }()
	return read(b)
}

func read(r io.Reader) (*Content, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	f, err := modpkg.Parse("go.mod", b, nil)
	if err != nil {
		return nil, err
	}
	return process(f)
}

func commentOf(c []modpkg.Comment) string {
	var sb strings.Builder
	for _, comment := range c {
		sb.WriteString(comment.Token)
		sb.WriteString("\n")
	}
	return sb.String()
}

func process(f *modpkg.File) (*Content, error) {
	c := new(Content)

	c.Module = f.Module.Mod.Path
	c.Go = f.Go.Version

	if f.Toolchain != nil {
		c.Toolchain = f.Toolchain.Name
	}

	// iterate every require block, combining them into just 2, one each for
	// direct and indirect dependencies
	for _, requirement := range f.Require {
		version, ok := semantic.Parse(requirement.Mod.Version)
		if !ok {
			return nil, fmt.Errorf("failed to parse module version %q", requirement.Mod.Version)
		}
		dependency := Dependency{
			Name:    requirement.Mod.Path,
			Version: version,
		}
		if requirement.Indirect {
			c.Indirect.add(dependency)
		} else {
			c.Direct.add(dependency)
		}
	}

	// iterate every replace block, combining them into just 2, one for normal
	// replacements and one for sub modules
	for _, replacement := range f.Replace {
		origVersion, _ := semantic.Parse(replacement.Old.Version) // version optional
		orig := Dependency{
			Name:    replacement.Old.Path,
			Version: origVersion,
		}
		nextVersion, ok := semantic.Parse(replacement.New.Version)
		if !ok && !strings.HasPrefix(replacement.New.Path, "./") {
			return nil, fmt.Errorf("failed to parse replacement version %q", replacement.New.Version)
		}
		next := Dependency{
			Name:    replacement.New.Path,
			Version: nextVersion,
		}
		r := Replacement{Orig: orig, Next: next}
		if strings.HasPrefix(next.Name, "./") {
			c.ReplaceSub.add(r)
		} else {
			c.Replace.add(r)
		}
	}

	// todo: retracts
	// todo: excludes

	return c, nil
}
