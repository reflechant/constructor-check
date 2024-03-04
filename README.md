# constructor-check

A go/analysis-compatible linter that reports struct type composite literals constructed manually while a constructor is defined in the same package as the type.

A constructor for type `T` (only structs are supported at the moment) is a function with name starting with "New" that returns a value of type `T` or `*T`.

## Usage

### Manually

1. Run `make build` or build the main.go in cmd/constructor-check however you see fit.
2. Run `./constructor_check [-flag] [package]`. You may run the command without parameters to get help message.

### With go vet
1. Run `make build` or build the main.go in cmd/constructor-check however you see fit.
2. Run `go vet --vettool=<path-to-constructor-check> [package]`

### With golangci-lint
Follow the instructions how to include a private plugin [here](https://golangci-lint.run/contributing/new-linters/#configure-a-plugin).

## Todo

- Check derived types (type T2 T)
- Check type aliases (type T2 = T)
- Use different diagnostic message on zero and nil values (?)
- Add flags to switch zero/nil values warnings (?)
- Maybe check constructor returned value instead of its name to extract type (they often rename types without renaming constructors)
- How about reporting constructors with function names inconsistent with type names ?
- Support other constructor signatures
    - (T, error) and (*T, error)
    - (T, bool) and (*T, bool)
- Work towards being included to golangci-lint