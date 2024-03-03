# constructor-check

A go/analysis-compatible linter that reports type composite literals constructed manually while a constructor is defined in the same package as the type.

A constructor for type `T` (only structs are supported at the moment) is a function with name starting with "New" that returns a value of type `T` or `*T`.

## Current state

The linter is in MVP state. It only reports composite literals in the same package the type is defined.

## Usage

To be described later.

## Todo

- Check types described in other packages
- Check derived types (type T2 T)
- Check type aliases (type T2 = T)
- Use different diagnostic message on zero and nil values
- Add flags to switch zero/nil values warnings
- Maybe check constructor returned value instead of its name to extract type (they often rename types without renaming constructors)
- Support other constructor signatures
    - (T, error) and (*T, error)
    - (T, bool) and (*T, bool)
- Move towards being adopted by golangci-lint