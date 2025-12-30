# go-enumr

`go-enumr` is a code generation tool for creating **Rich Enums** in Go.

Unlike standard Go enums (which are just integers), `go-enumr` allows you to use **structs** as enums. This enables you to bundle associated data—like labels, database IDs, configuration values, or color codes—directly with your enum identity.

## Why Struct Enums?

Standard Go enums are great for simple flags:

```go
// Standard Go Enum
type Status int
const (
    Pending Status = iota
    Active
)
// Problem: How do I get the string label? The database ID? The color for the UI?
// You usually end up writing massive switch statements or lookup maps.
```

**Rich Enums** solve this by keeping data together:

```go
// go-enumr Rich Enum
type Status struct {
    ID    int
    Label string
    Color string
}

var (
    Pending = Status{0, "Pending", "#FFA500"}
    Active  = Status{1, "Active",  "#00FF00"}
)
```

`go-enumr` automates the boilerplate needed to make these structs behave like proper enums, generating `String()`, `MarshalText()`, and `UnmarshalText()` methods for you.

## Installation

```bash
go install github.com/jmfrees/go-enumr@latest
```

## Usage

1.  Define your struct type.
2.  Define your enum instances as package-level variables.
3.  Add the `//go:generate` directive.

### Example

```go
package payment

//go:generate enumr -type=Method -format=snake_case
type Method struct {
    Code        string
    Description string
    IsCredit    bool
}

var (
    CreditCard   = Method{"CC", "Credit Card", true}
    PayPal       = Method{"PP", "PayPal", false}
    BankTransfer = Method{"BT", "Bank Transfer", false}
)
```

Run the generator:

```bash
go generate ./...
```

This creates `method_string.go` containing:

- `func (m Method) String() string` -> returns "credit_card", "pay_pal", etc.
- `func (m Method) MarshalText() ([]byte, error)` -> allows JSON/XML marshaling.
- `func (m *Method) UnmarshalText([]byte) error` -> allows JSON/XML unmarshaling back to the correct struct instance.

## CLI Options

- `-type`: (Required) Comma-separated list of type names to generate code for.
- `-format`: (Optional) Casing format for the string representation.
  - `snake_case`
  - `camelCase`
  - `PascalCase`
  - `SNAKE_CASE`
  - `Title Case`
- `-output`: (Optional) Output file name. Defaults to `<type>_string.go` (or `<first_type>_string.go` if multiple types are provided).

## Best Practices

Since Go structs cannot be `const`, these enums are defined as `var`. While technically mutable, the convention is to treat them as immutable constants.

## License

MIT
