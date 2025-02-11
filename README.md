# Go-Enumr

`go-enumr` is meant to be a more flexible implementation of the `stringer` library by Rob Pike. 

```go
package type

//go:generate enumr -type=Type
type Type struct {
	name string
	ext int
}

// Define some constant instances of the enum
var (
	Foo = Type{"foo", 100)
	Bar = Type{"bar", 200)
	Baz = Type{"baz", 300)
)
```

This will create a new file in the same package called `type_string.go` with the following content:

```go
package type

func (t Type) MarshalText() (text []byte, err error) {
	return []byte(t.String()), nil
}

func (t *Type) UnmarshalText(text []byte) error {
	trimmedText := strings.ReplaceAll(strings.ToLower(string(text)), "\"", "")
	switch trimmedText {
	case "Foo":
		*t = Foo
	case "Bar":
		*t = Bar
	case "Baz":
		*t = Baz
	default:
		return fmt.Errorf("unsupported type")
	}
	return nil
}
```
