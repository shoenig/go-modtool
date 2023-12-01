module example.com/my/module

// An old comment

go 1.21.2

toolchain go1.21

// The old e2e replace comment
replace (
    golang.org/x/mod => golang.org/x/mod v1.0.0-alpha.1
    github.com/hashicorp/go-set/v2 => github.com/shoenig/go-set/v2 v2.0.0-alpha3
)

// The old e2e submodules comment
replace (
    internal.com/api => ./api
)

require (
    github.com/BurntSushi/toml v1.4.2
    github.com/hashicorp/go-set/v2 v2.0.0-alpha.2
    internal.com/api v0.0.0
    golang.org/x/mod v0.12.0
    github.com/shoenig/semantic v1.2.1
    github.com/shoenig/test v0.6.7
)

require (
    github.com/google/go-cmp v0.5.9 // indirect
    github.com/shoenig/regexplus v0.3.0 // indirect
    example.com/something/else v0.1.0 // indirect
    example.com/more/stuff v1.2.0 // indirect
)

// The e2e exclude comment
exclude (
    github.com/shoenig/test v0.5.0
)

