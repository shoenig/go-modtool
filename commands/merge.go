package commands

import (
	"cattlecloud.net/go/babycli"
)

const (
	mergeName        = `merge`
	mergeHelp        = `Use merge to combine two go.mod files.`
	mergeDescription = `
The merge command will combine two go.mod files so that the result
follows a strict format, including ordering of each stanza and
alphabetizing packages within each stanza. A common use case for
the merge command is in merging changes between an open source
and a closed source version of the same repository.
`
)

func newMergeCommand() *babycli.Component {
	return &babycli.Component{
		Name:        mergeName,
		Help:        mergeHelp,
		Description: mergeDescription,
		Function: func(c *babycli.Component) babycli.Code {
			t := NewTool(c)

			// load the config file if set
			if err := t.applyConfig(); err != nil {
				return crash(err)
			}

			// run the merge command
			if err := t.merge(c.Arguments()); err != nil {
				return crash(err)
			}

			return babycli.Success
		},
	}
}
