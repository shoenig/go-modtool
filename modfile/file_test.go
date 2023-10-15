// Copyright (c) Seth Hoenig
// SPDX-License-Identifier: MPL-2.0

package modfile

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/shoenig/test/must"
)

func cat(t *testing.T, path string) string {
	b, err := os.ReadFile(path)
	must.NoError(t, err)
	return string(b)
}

func TestContent_Write(t *testing.T) {
	cases := []struct {
		name string
	}{
		{name: "simple"},
		{name: "parens"},
		{name: "comments"}, // no comments
		{name: "replace"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			in := filepath.Join("tests", tc.name+".mod")
			exp := filepath.Join("tests", tc.name+".expect")
			c, err := Open(in)
			must.NoError(t, err)

			var b bytes.Buffer
			err = c.Write(&b)
			must.NoError(t, err)

			must.Eq(t, cat(t, exp), b.String())
		})
	}
}
