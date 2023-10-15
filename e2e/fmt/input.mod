// comment

module example.com/my/module

go 1.21.2

toolchain go1.21

replace golang.org/x/mod => golang.org/x/mod v1.0.0-alpha.1

require (
	github.com/BurntSushi/toml v1.3.2
	github.com/hashicorp/go-set/v2 v2.0.0-alpha.2
	internal.com/api v0.0.0
)

replace (
	github.com/hashicorp/go-set/v2 => github.com/shoenig/go-set/v2 v2.0.0-alpha3
	internal.com/api => ./api
)

require golang.org/x/mod v0.12.0

require (
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/shoenig/regexplus v0.3.0 // indirect
)

require(
	example.com/something/else v0.1.0 // indirect
)

require example.com/more/stuff v1.2.0 // indirect

require (
	github.com/shoenig/semantic v1.2.1
	github.com/shoenig/test v0.6.7
)

exclude (
	github.com/shoenig/test v0.5.0
)