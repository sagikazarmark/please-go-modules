# Please Go Modules experiment

Experiments with [Please](https://please.build) and Go modules.


## Usage

```
# Build the binary with all the dependencies
./pleasew build //:bin

# Run the binary
./pleasew run //:bin

# For more details
./pleasew build //:bin --rebuild -v debug --show_all_output
```


## Current approach

The `go_module` rule downloads all dependencies using `go mod download`,
generates a rule for each package (passed to the same rule) which builds those packages.

The generated rules can be referenced by `go_library` and `go_binary` rules.

The name format is `_{module}#install_{package_name}` where slashes in the package name is replaced with underscores.

See details in the source files.


## Notes / Questions

- The package build rule output might collide with other build rules?
- How to handle transitive dependencies?
- How about the generated rule name scheme?
- Can the package list be generated? (related to transitive dependencies)
