module github.com/x/example

go 1.20

toolchain go1.19.3

require (
	github.com/foo/bar v1.0.0
	github.com/private/blah v0.1.0
	github.com/example/test1/sdk v0.0.0
)

replace github.com/foo/bar => github.com/goo/bar v1.1.0

require (
	github.com/apple/apple v1.2.3 // indirect
)

replace github.com/example/test1/sdk => ./sdk

