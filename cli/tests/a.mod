
module github.com/shoenig/example


go 1.21

// a comment
require example.com/a/b v1.0.0

// another comment
require(
	 example.com/c/d v1.1.0
	example.com/e/f v1.2.0
)

require example.com/g/h v1.3.0 // indirect

replace example.com/c/d => example.com/x/y v0.4.0

require github.com/shoenig/example/api v1.1.1

replace (
	github.com/shoenig/example/api => ./api
)
