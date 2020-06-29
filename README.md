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


## Change log

This experiment went through a couple versions, hit a few dead ends.
Every major turn should be documented here so in the future we know how and why decisions have been made.

## Latest revision - 2020-06-29



## First few revisions - 2020 June

The idea is to reuse Go modules and `go.mod` to download and manage dependencies,
therefor I experimented with a couple rules.

The final rule turned out to be less automated than I'd hoped,
but worked well and used caching as much as possible.

Basically a `go_module` rule executed `go mod download` and then built the packages listed in the rule.

This worked quite well with direct dependencies, but unfortunately Go modules doesn't allow building and saving the output of transitive dependencies.