# go-modtool

An opinionated tool to format and merge `go.mod` (and soon `go.sum`) files.

## Overview

Although the `go mod tidy` command provided by the Go compiler will format `go.mod` files,
the way it does so is not strongly opinionated. Often, you'll end up with multiple `replace`
stanzas where only one was intended. This `go-modtool` tool provides a way to format `go.mod`
files in a way that is consistent no matter the input.

Projects which need to merge `go.mod` files will also experience pain in how git often
fails to resolve conflicts for `go.mod` and `go.sum` files. This `go-modtool` tool can be
used to merge `go.mod` and `go.sum` with better context files so that merge conflicts are 
automatically resolved, every time.

### Strict go.mod Format

The strongly opinionated format of a `go.mod` file is described below.

- The `module` stanza appears on the first line of the file.

- The `go` version stanza appears after the `module` stanza.

- A `replace` stanza appears next (if necessary) for third party modules.

- A `replace` stanza appears next (if necessary) for local submodules.

- A `require` stanza appears next (if necessary) for direct dependencies.

- A `require` stanza appears next (if necessary) for indirect dependencies.

Comments are disallowed, with exceptions being enabled by specifying them through
arguments to the `go-modtool fmt` command.

## Getting Started

`go-modtool` is a command line tool written in Go. It can be built from source, installed
via `go install`, or downloaded from [Releases](https://github.com/shoenig/go-modtool/releases).

#### Install

Install using `go install`

```shell
go install github.com/shoenig/go-modtool@latest
```

## Subcommands

The `go-modtool` command line tool provides two subcommands, for formatting and merging
`go.mod` files.

#### fmt

The `fmt` subcommand is used to format a `go.mod` file.

```shell
go-modtool fmt go.mod
```

By default the output is printed to standard out. Use the `-w` flag to overwrite the
input `go.mod` file in place.

```shell
go-modtool -w fmt go.mod
```

The following flags enable specifying comments in the resulting `go.mod` file.

- `--replace-comment` - Insert a comment before the `replace` stanza for third party modules.

- `--subs-comment` - Insert a comment before the `replace` stanza for submodules.

#### merge

The `merge` subcommand is used to merge two `go.mod` files.

The motivating use case for this is in merging the `go.mod` file of an OSS version
of a repository with a private ENT version of the same reposotiry, where the ENT
version is a superset of OSS. Conflicts happen because git is not smart enough to
resolve differences caused by module changes in nearby lines.

```shell
go-modtool merge /ent/go.mod /oss/go.mod
```

The same CLI arguments from the `fmt` command apply to the `merge` command.

- `-w` - Write the output to the first `go.mod` file in the input.

- `--replace-comment` - Insert a comment before the `replace` stanza for third party modules.

- `--subs-comment` - Insert a comment before the `replace` stanza for submodules.

Note that it is imperative to run `go mod tidy` after doing a merge to ensure only
the necessary modules remain in the resulting `go.mod` and `go.sum` files.

## Config File

To make CI easier the subcommands above can be described in a TOML config file.

```toml
ReplaceComment    = "An example replace comment."
SubmodulesComment = "An example submodules comment."
WriteFile         = true
```

Use the `--config` flag to specify a path to the config file.

```shell-session
go-modtool -config=modtool.toml fmt go.mod
```

## Contributing

The `github.com/shoenig/go-modtool` module / cli tool is always improving with new features
and bug fixes. For contributing such bug fixes and new features please file an issue.

## License

The `github.com/shoenig/go-modtool` module is open source under the [MPL-2.0](LICENSE) license.

