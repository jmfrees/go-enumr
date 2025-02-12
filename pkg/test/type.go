package test

//go:generate ./../../enumr -type=Type -format=snake_case
type Type struct {
	v1 int
	v2 string
}

var (
	Foo        = Type{100, "a"}
	Bar        = Type{200, "b"}
	Baz        = Type{300, "c"}
	LongerName = Type{300, "d"}
)

var Test = 1
