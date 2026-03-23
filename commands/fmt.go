package commands

import (
	"cattlecloud.net/go/babycli"
)

const (
	fmtName        = `fmt`
	fmtHelp        = `Use fmt to format a go.mod file.`
	fmtDescription = `
The fmt command will rewrite a go.mod file so that it follows a
strict format, including ordering of each stanza and alphabetizing
packages within each stanza.`
)

func newFmtCommand() *babycli.Component {
	return &babycli.Component{
		Name:        fmtName,
		Help:        fmtHelp,
		Description: fmtDescription,
		Function: func(c *babycli.Component) babycli.Code {
			t := NewTool(c)

			// load the config file if set
			if err := t.applyConfig(); err != nil {
				return crash(err)
			}

			// run the fmt command
			if err := t.fmt(c.Arguments()); err != nil {
				return crash(err)
			}

			return babycli.Success
		},
	}
}
