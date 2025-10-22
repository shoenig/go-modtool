// Copyright (c) Seth Hoenig
// SPDX-License-Identifier: MPL-2.0

package modfile

import (
	"github.com/hashicorp/go-set/v3"
)

// Merge will combine the changes of o onto c, resolving conflicts
// a way that makes sense.
//
// e.g.
// oss - from the open source repository
// ent - from the enterprise repository
// and we are merging in changes from oss onto ent, which may cause
// conflicts due to ent being a superset of oss.
func Merge(ent, oss *Content) *Content {
	directOSS := set.HashSetFrom(oss.Direct.Dependencies)
	directENT := set.HashSetFrom(ent.Direct.Dependencies)

	indirectOSS := set.HashSetFrom(oss.Indirect.Dependencies)
	indirectENT := set.HashSetFrom(ent.Indirect.Dependencies)

	replaceOSS := set.HashSetFrom(oss.Replace.Replacements)
	replaceENT := set.HashSetFrom(ent.Replace.Replacements)

	subsOSS := set.HashSetFrom(oss.ReplaceSub.Replacements)
	subsENT := set.HashSetFrom(ent.ReplaceSub.Replacements)

	direct := dependencies(directOSS, directENT).Slice()
	indirect := dependencies(indirectOSS, indirectENT).Slice()
	replace := replacements(replaceOSS, replaceENT).Slice()
	subs := replacements(subsOSS, subsENT).Slice()

	c := &Content{
		Module:     ent.Module,
		Go:         oss.Go,
		Toolchain:  oss.Toolchain,
		Direct:     BasicStanza{Dependencies: direct},
		Indirect:   BasicStanza{Dependencies: indirect},
		Replace:    ReplaceStanza{Replacements: replace},
		ReplaceSub: ReplaceStanza{Replacements: subs},
	}

	c.sort()

	return c
}

func dependencies(oss, ent *set.HashSet[Dependency, string]) set.Collection[Dependency] {
	// blunt merge with oss overriding ent
	return oss.Union(ent)
}

func replacements(oss, ent *set.HashSet[Replacement, string]) set.Collection[Replacement] {
	// blunt merge with oss overriding ent
	return oss.Union(ent)
}
