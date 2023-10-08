package modfile

import (
	"testing"

	"github.com/shoenig/semantic"
	"github.com/shoenig/test/must"
)

func TestMerge(t *testing.T) {
	cases := []struct {
		name string
		ent  *Content
		oss  *Content
		exp  *Content
	}{
		{
			name: "empty",
			ent:  new(Content),
			oss:  new(Content),
			exp: &Content{
				Direct:     RequireStanza{Dependencies: make([]Dependency, 0)},
				Indirect:   RequireStanza{Dependencies: make([]Dependency, 0)},
				Replace:    ReplaceStanza{Replacements: make([]Replacement, 0)},
				ReplaceSub: ReplaceStanza{Replacements: make([]Replacement, 0)},
			},
		},
		{
			name: "upgrade",
			ent: &Content{
				Module: "module",
				Go:     "1.20",
				Direct: RequireStanza{Dependencies: []Dependency{
					{Name: "example.com/one/blue", Version: semantic.New(1, 0, 0)},
				}},
				Indirect:   RequireStanza{Dependencies: make([]Dependency, 0)},
				Replace:    ReplaceStanza{Replacements: make([]Replacement, 0)},
				ReplaceSub: ReplaceStanza{Replacements: make([]Replacement, 0)},
			},
			oss: &Content{
				Module: "module",
				Go:     "1.20",
				Direct: RequireStanza{Dependencies: []Dependency{
					{Name: "example.com/one/blue", Version: semantic.New(1, 0, 1)},
				}},
				Indirect:   RequireStanza{Dependencies: make([]Dependency, 0)},
				Replace:    ReplaceStanza{Replacements: make([]Replacement, 0)},
				ReplaceSub: ReplaceStanza{Replacements: make([]Replacement, 0)},
			},
			exp: &Content{
				Module: "module",
				Go:     "1.20",
				Direct: RequireStanza{Dependencies: []Dependency{
					{Name: "example.com/one/blue", Version: semantic.New(1, 0, 1)},
				}},
				Indirect:   RequireStanza{Dependencies: make([]Dependency, 0)},
				Replace:    ReplaceStanza{Replacements: make([]Replacement, 0)},
				ReplaceSub: ReplaceStanza{Replacements: make([]Replacement, 0)},
			},
		},
		{
			name: "mix",
			ent: &Content{
				Module:    "example.com/project/ent",
				Go:        "1.20",
				Toolchain: ToolchainStanza{Version: "1.20"},
				Direct: RequireStanza{Dependencies: []Dependency{
					{Name: "example.com/one/blue", Version: semantic.New(1, 0, 0)},
					{Name: "example.com/one/green", Version: semantic.New(0, 3, 2)},
					{Name: "example.com/private/orange", Version: semantic.New(1, 2, 3)}, // ent only
					{Name: "example.com/two/purple", Version: semantic.New(1, 4, 5)},
				}},
				Indirect: RequireStanza{Dependencies: []Dependency{
					{Name: "example.com/x/one", Version: semantic.New(0, 0, 1)}, // ent only
					{Name: "example.com/x/two", Version: semantic.New(0, 2, 0)},
				}},
				Replace: ReplaceStanza{Replacements: []Replacement{
					{
						Orig: Dependency{Name: "example.com/one/blue"},
						Next: Dependency{Name: "internal.com/blue", Version: semantic.New(0, 0, 1)},
					},
				}},
				ReplaceSub: ReplaceStanza{Replacements: []Replacement{
					{
						Orig: Dependency{Name: "example.com/foo/bar/api"},
						Next: Dependency{Name: "./api"},
					},
				}},
			},
			oss: &Content{
				Module:    "example.com/project/oss",
				Go:        "1.21",                           // upgrade
				Toolchain: ToolchainStanza{Version: "1.21"}, // upgrade
				Direct: RequireStanza{Dependencies: []Dependency{
					{Name: "example.com/one/pink", Version: semantic.New(1, 0, 1)}, // introduction
					{Name: "example.com/one/blue", Version: semantic.New(1, 0, 0)},
					{Name: "example.com/one/green", Version: semantic.New(0, 3, 3)}, // upgrade
					{Name: "example.com/two/purple", Version: semantic.New(1, 4, 5)},
				}},
				Indirect: RequireStanza{Dependencies: []Dependency{
					{Name: "example.com/x/two", Version: semantic.New(0, 2, 1)},   // upgrade
					{Name: "example.com/x/three", Version: semantic.New(0, 1, 1)}, // introduction
				}},
				Replace: ReplaceStanza{Replacements: []Replacement{
					//
				}},
				ReplaceSub: ReplaceStanza{Replacements: []Replacement{
					{
						Orig: Dependency{Name: "example.com/foo/bar/api"},
						Next: Dependency{Name: "./api"},
					},
				}},
			},
			exp: &Content{
				Module:    "example.com/project/ent",
				Go:        "1.21",
				Toolchain: ToolchainStanza{Version: "1.21"},
				Direct: RequireStanza{Dependencies: []Dependency{
					{Name: "example.com/one/blue", Version: semantic.New(1, 0, 0)},
					{Name: "example.com/one/green", Version: semantic.New(0, 3, 3)},      // upgrade
					{Name: "example.com/one/pink", Version: semantic.New(1, 0, 1)},       // introduction
					{Name: "example.com/private/orange", Version: semantic.New(1, 2, 3)}, // ent only
					{Name: "example.com/two/purple", Version: semantic.New(1, 4, 5)},
				}},
				Indirect: RequireStanza{Dependencies: []Dependency{
					{Name: "example.com/x/one", Version: semantic.New(0, 0, 1)},   // ent only
					{Name: "example.com/x/three", Version: semantic.New(0, 1, 1)}, // introduction
					{Name: "example.com/x/two", Version: semantic.New(0, 2, 1)},   // upgrade
				}},
				Replace: ReplaceStanza{Replacements: []Replacement{
					{
						Orig: Dependency{Name: "example.com/one/blue"},
						Next: Dependency{Name: "internal.com/blue", Version: semantic.New(0, 0, 1)},
					},
				}},
				ReplaceSub: ReplaceStanza{Replacements: []Replacement{
					{
						Orig: Dependency{Name: "example.com/foo/bar/api"},
						Next: Dependency{Name: "./api"},
					},
				}},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := Merge(tc.ent, tc.oss)
			must.Eq(t, tc.exp, result)
		})
	}
}
