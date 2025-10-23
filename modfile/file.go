// Copyright (c) Seth Hoenig
// SPDX-License-Identifier: MPL-2.0

package modfile

import (
	_ "embed"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
	"text/template"

	"github.com/shoenig/semantic"
	modpkg "golang.org/x/mod/modfile"
)

var (
	zero = semantic.New(0, 0, 0)

	//go:embed go.mod.tmpl
	goModTemplateBody string

	goModTemplate = template.Must(template.New("go.mod.tmpl").Parse(goModTemplateBody))
)

type Dependency struct {
	Replacement bool
	Tool        bool
	Name        string
	Version     semantic.Tag
}

func (d Dependency) Hash() string {
	return d.Name
}

func (d Dependency) String() string {
	if d.Version.Equal(zero) && (d.Replacement || d.Tool) {
		return d.Name
	}
	return fmt.Sprintf("%s %s", d.Name, d.Version)
}

func (d Dependency) Cmp(o Dependency) int {
	if d.Name < o.Name {
		return -1
	} else if d.Name > o.Name {
		return 1
	}
	if d.Version.Less(o.Version) {
		return -1
	} else if o.Version.Less(d.Version) {
		return 1
	}
	return 0
}

type ReplaceStanza struct {
	Comment      string
	Replacements []Replacement
}

func (rs *ReplaceStanza) add(r Replacement) {
	rs.Replacements = append(rs.Replacements, r)
}

func (rs *ReplaceStanza) sort() {
	slices.SortFunc(rs.Replacements, func(a, b Replacement) int {
		return a.Cmp(b)
	})
}

type Replacement struct {
	Orig Dependency
	Next Dependency
}

func (r Replacement) Hash() string {
	return r.Orig.Name
}

func (r Replacement) Cmp(o Replacement) int {
	return r.Orig.Cmp(o.Orig)
}

func (r Replacement) String() string {
	return fmt.Sprintf("%s => %s", r.Orig, r.Next)
}

type BasicStanza struct {
	Comment      string
	Dependencies []Dependency
}

func (ds *BasicStanza) add(d Dependency) {
	ds.Dependencies = append(ds.Dependencies, d)
}

func (ds *BasicStanza) sort() {
	slices.SortFunc(ds.Dependencies, func(a, b Dependency) int {
		return a.Cmp(b)
	})
}

type ToolchainStanza struct {
	Comment string
	Version string // arbitrary
}

type Content struct {
	Module     string
	Go         string
	Toolchain  ToolchainStanza
	Direct     BasicStanza
	Indirect   BasicStanza
	Replace    ReplaceStanza
	ReplaceSub ReplaceStanza // sub modules, e.g. "=> ./api"
	Exclude    BasicStanza
	Tool       BasicStanza

	// Retract   []semantic.Tag
}

func (c *Content) sort() {
	c.Direct.sort()
	c.Indirect.sort()
	c.Replace.sort()
	c.ReplaceSub.sort()
	c.Exclude.sort()
	c.Tool.sort()
}

func (c *Content) String() string {
	return "todo"
}

func (c *Content) Write(w io.Writer) error {
	return goModTemplate.Execute(w, c)
}

func Open(path string) (*modpkg.File, error) {
	b, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = b.Close() }()
	return read(b)
}

func read(r io.Reader) (*modpkg.File, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	f, err := modpkg.Parse("go.mod", b, nil)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func Process(f *modpkg.File) (*Content, error) {
	c := new(Content)

	c.Module = f.Module.Mod.Path
	c.Go = f.Go.Version

	if f.Toolchain != nil {
		c.Toolchain = ToolchainStanza{
			Version: f.Toolchain.Name,
		}
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

	isLocal := func(name string) bool {
		return strings.HasPrefix(name, "./") || strings.HasPrefix(name, "../")
	}

	// iterate every replace block, combining them into just 2, one for normal
	// replacements and one for sub modules
	for _, replacement := range f.Replace {
		origVersion, _ := semantic.Parse(replacement.Old.Version) // version optional
		orig := Dependency{
			Replacement: true,
			Name:        replacement.Old.Path,
			Version:     origVersion,
		}
		nextVersion, ok := semantic.Parse(replacement.New.Version)
		if !ok && !isLocal(replacement.New.Path) {
			return nil, fmt.Errorf("failed to parse replacement version %q", replacement.New.Version)
		}
		next := Dependency{
			Replacement: true,
			Name:        replacement.New.Path,
			Version:     nextVersion,
		}
		r := Replacement{Orig: orig, Next: next}
		if isLocal(next.Name) {
			c.ReplaceSub.add(r)
		} else {
			c.Replace.add(r)
		}
	}

	// iterate every exclude block, combining them into one
	for _, exclude := range f.Exclude {
		version, ok := semantic.Parse(exclude.Mod.Version)
		if !ok {
			return nil, fmt.Errorf("failed to parse exclude version %q", exclude.Mod.Version)
		}
		dependency := Dependency{
			Name:    exclude.Mod.Path,
			Version: version,
		}
		c.Exclude.add(dependency)
	}

	// iterate every exclude block, combining them into one
	for _, tool := range f.Tool {
		dependency := Dependency{
			Tool: true,
			Name: tool.Path,
		}
		c.Tool.add(dependency)
	}

	c.sort()

	// todo: retracts

	return c, nil
}
