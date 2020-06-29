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

In order to use the `go_getx` rule in a new project copy the `build_defs` directory,
preload the build definitions in it and follwo these commands:

```bash
# This is necessary because Go filter expressions take plz-out into account and can mess with the following commands
mkdir -p plz-out
cd plz-out
go mod init plz-out
cd -

# This command generates basic go_getx rules for you
# It is VERY-VERY FAR from ideal (take a look at BUILD in this repo for manual edits)
go list -m -mod=readonly all | while read mod; do cat go.sum | grep "$mod " | xargs bash -c 'NAME=$(echo $0 | sed "s|\/|_|g") && echo -e "go_getx(\n    name=\"$NAME\",\n    get=\"$0\",\n    version=\"$1\",\n    sum=\"$2\",\n    visibility=[\"PUBLIC\"],\n)"'; done

# Generate a list of dependencies
go list -m -mod=readonly all | while read mod; do cat go.sum | grep "$mod " | cut -d ' ' -f1 | sed 's/.*/":&",/; s|\/|_|g'; done
```

Ideally, the above lines should be replaced by a code generator that generates third-party `go_get`s for you, resolving transparent dependencies.
One such tool could be [wollemi](https://github.com/tcncloud/wollemi), but it's still in early development.


## Notes / Questions

- The initial `go_get` for fetch_rule is **really** slow. Why?


## Change log

This experiment went through a couple versions, hit a few dead ends.
Every major turn should be documented here so in the future we know how and why decisions have been made.


## Latest revision - 2020-06-29

Using the Go tooling for downloading turned out to be a dead end,
because in module mode object files can only be built and saved for direct dependencies, but not transitive ones.
(Chances are this is not true, but I couldn't find the right solution)

After looking at [Gazelle](https://github.com/bazelbuild/bazel-gazelle), it turns out that they don't use go modules for package management either.
This is kind of understandable, since go modules break the incremental build idea behind Bazel/Please (every change to `go.mod` would result it massive rebuilds).

Gazelle actually uses a small tool, called [fetch_repo](https://github.com/bazelbuild/bazel-gazelle/tree/5c00b77/cmd/fetch_repo) to download packages.
It supports regular `go get` mode (and more) and works with modules as well. It uses a little trick to do that: it creates a temporary module and downloads the package using `go mod download`.

I was able to integrate this tool into the existing `go_get` rule with a few changes. The working title is `go_getx`.
If this tool gets integrated into the upstream `go_get` rule, it might worth checking if it can replace the rest of the current download logic.


## First few revisions - 2020 June

The idea is to reuse Go modules and `go.mod` to download and manage dependencies,
therefor I experimented with a couple rules.

The final rule turned out to be less automated than I'd hoped,
but worked well and used caching as much as possible.

Basically a `go_module` rule executed `go mod download` and then built the packages listed in the rule.

This worked quite well with direct dependencies, but unfortunately Go modules doesn't allow building and saving the output of transitive dependencies.