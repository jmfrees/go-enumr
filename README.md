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

`go-enumr` automates the boilerplate needed to make these structs behave like proper enums, generating `String()`, `MarshalText()`, `UnmarshalText()`, and helper functions like `Values()` and `Parse()`.

## Installation

```bash
go install github.com/jmfrees/go-enumr/cmd/enumr@latest
```

## Usage

There are two ways to use `go-enumr`:

1.  **Directive Mode (Recommended)**: Define your enum values in comments.
2.  **Manual Mode**: Define your enum values in a `var` block.

### 1. Directive Mode (Recommended)

Use `//enumr:<Name>` comments to define your enum instances. The tool will generate the `var` block for you.

```go
package payment

//go:generate enumr -type=Method

//enumr:CreditCard   Code:CC   Description:"Credit Card"  IsCredit:true
//enumr:PayPal       Code:PP   Description:PayPal
//enumr:BankTransfer Code:BT   Description:"Bank Transfer"
type Method struct {
    Code        string
    Description string
    IsCredit    bool
}
```

Run the generator:

```bash
go generate ./...
```

This creates `method_enum.go` containing:

```go
var (
    CreditCard   = Method{Code: "CC", Description: "Credit Card", IsCredit: true}
    PayPal       = Method{Code: "PP", Description: "PayPal"}
    BankTransfer = Method{Code: "BT", Description: "Bank Transfer"}
)

// ... generated methods ...
```

**Syntax Rules:**
*   `Key:Value` sets a field (e.g., `Code:CC`).
*   `Key:"Value with spaces"` sets a string field with spaces.
*   `Key:true` sets a boolean field (e.g., `IsCredit:true`).
*   `Key:"[]string{\"a\", \"b\"}"` sets complex types (slices, structs) by quoting the Go syntax.
*   Fields not specified default to their zero value.

### 2. Manual Mode

If you need complex initialization (e.g., function calls, external imports) or want to document individual instances, you can define the variables yourself.

**When to use Manual Mode:**
*   You need to use functions like `time.Now()`.
*   You need to import other packages.
*   You want to add godoc comments to specific enum instances.

```go
package payment

//go:generate enumr -type=Method

type Method struct {
    Code string
    Created time.Time
}

var (
    CreditCard = Method{"CC", time.Now()}
)
```

`go-enumr` will detect the existing `var` block and generate the helper methods for it.

## Generated Code

The tool generates the following for your type:

- `func (t Type) String() string`: Returns the enum name (e.g., "credit_card").
- `func (t Type) MarshalText() ([]byte, error)`: Implements `encoding.TextMarshaler`.
- `func (t *Type) UnmarshalText([]byte) error`: Implements `encoding.TextUnmarshaler`.
- `func TypeValues() []Type`: Returns a slice of all enum instances.
- `func ParseType(s string) (Type, error)`: Helper to parse a string into an enum instance.

## CLI Options

- `-type`: (Required) Comma-separated list of type names to generate code for.
- `-format`: (Optional) Casing format for the string representation.
  - `snake_case` (default)
  - `camelCase`
  - `PascalCase`
  - `SNAKE_CASE`
  - `Title Case`
- `-output`: (Optional) Output file name or directory. Defaults to `<type>_enum.go` in the package directory.

## Best Practices

Since Go structs cannot be `const`, these enums are defined as `var`. While technically mutable, the convention is to treat them as immutable constants.

## License

MIT
