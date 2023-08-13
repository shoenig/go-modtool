package cli

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/shoenig/test/must"
)

// setup copies the input source to a temp file we can modify
// in the test
func setup(t *testing.T) string {
	// copy the input and return that filepath
	dir := t.TempDir()
	f := filepath.Join(dir, "go.mod")
	file, err := os.OpenFile(f, os.O_CREATE|os.O_WRONLY, 0644)
	must.NoError(t, err)

	orig, err := os.Open("tests/a.mod")
	must.NoError(t, err)

	_, err = io.Copy(file, orig)
	must.NoError(t, err)
	must.Close(t, file)

	return f
}

func cat(t *testing.T, filename string) string {
	b, err := os.ReadFile(filename)
	must.NoError(t, err)
	return string(b)
}

func compare(t *testing.T, exp, actual string) {
	expContent := cat(t, exp)
	actualContent := cat(t, actual)
	must.Eq(t, expContent, actualContent)
}

func TestTool_fmt(t *testing.T) {
	modFile := setup(t)

	tool := &Tool{
		writeFile:   true,
		repsComment: "This is a comment about replacements.",
		subsComment: "This is a comment about submodules.",
		modFile:     modFile,
	}
	args := []string{modFile}

	err := tool.fmt(args)
	must.NoError(t, err)

	compare(t, "tests/a.expect", modFile)
}

func TestTool_merge(t *testing.T) {
	// eh todo
}
