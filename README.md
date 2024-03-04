# constructor-check

A linter to report ignored constructors. It shows you places where someone is doing T{} or &T{} instead of using NewT declared in the same package as T ( new(T) is not yet reported ).

A constructor for type `T` (only structs are supported at the moment) is a function with name "NewT" that returns a value of type `T` or `*T`. Types returned by constructors are not checked right now, only that type T inferred from the function name exists in the same package.

## Why?

Nil maps are unsafe to write to (unlike slices; why, Rob, why?), you may want to have default values for fields that are not equal to zero values, etc. etc. We create constructors for a reason.

Constructors naming convention is just an established pattern and calling them a matter of discipline. But people make mistakes and overlook things all the time. Errare humanum est. And what if a constructor was created later? Who's going to check all the places where the type is used?

"Make zero values useful" is a good proverb but it's an unachievable utopia as it has always been. You should not create constructors unnecessarily but if you have to you better make sure they are used everywhere for this type.

Yes, you may declare an interface and hide the actual (unexported) type behind it and return an interface from the constructor. Which breaks the "return types, accept interfaces" proverb and arguably is an overengineering of a simple problem. Why create an interface while you can just initialize all the type instances properly? And this linter will help to make sure you do.

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