module replacements

go 1.21.0

// comment line 1
require (
	example.com/a/b v1.0.0
	example.com/c/d/v2 v2.2.2
	example.com/api v1.1.1
)

// comment line 2
replace example.com/a/b => example.com/a/b v1.1.0

// comment line 3
replace example.com/c/d/v2 => other.com/c/d/v2 v2.4.0

// sub package
replace example.com/api => ./api
