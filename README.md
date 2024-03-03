# constructor-check

A go/analysis-compatible linter that reports type composite literals constructed manually while a constructor is defined in the same package as the type.

A constructor for type `T` (only structs are supported at the moment) is a function with name starting with "New" that returns a value of type `T` or `*T`.

## Current state

The linter is in MVP state. It only reports non-zero, non-nil composite literate in the same package the 

## Usage

To be described later.

## Todo

- Check for composite literals inside containers
- Check types described in other packages
- Check derived types (type T2 T)
- Check type aliases (type T2 = T)
- Check for composite literals in struct fields
- Warn on zero and nil values being potentially unsafe if a constructor is defined for their type
- Add flags to switch zero/nil values warnings
- Maybe check constructor returned value instead of its name to extract type (they often rename types without renaming constructors)
- Support other constructor signatures
    - (T, error) and (*T, error)
    - (T, bool) and (*T, bool)
- Move towards being adopted by golangci-lint