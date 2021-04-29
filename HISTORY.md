# History

This repository started as an experiment for Go modules support in [Please](https://please.build/).
Then Please added experimental Go modules support, so this tool became a rule generator.

This file preserves the history of the experiment.

---

## Add custom Go rules

Some releases might need to come with custom Go rules (eg. temporarily fixing problems, etc).

1. Add `"go_rules.build_defs",` to `srcs` in `dist/build_file`
1. Add `"go_rules.build_defs",` to `srcs` in `build_defs/BUILD`
2. Add a `go_rules.build_defs` file to `build_defs/`


## Change log

This experiment went through a couple versions, hit a few dead ends.
Every major turn should be documented here so in the future we know how and why decisions have been made.


## Latest revision - 2021-02-17

Please introduced a `go_module` rule that basically removes the need to generate `go_library`
rules up front with static build flags.


## `go_library` with custom go module download - 2020-09-03

After a number of attempts to make Please Go module compatible using the high level Go toolchain (`build`, `install`),
I gave up using the high level tools and turned to generating standard `go_library` rules for packages instead.

It turned out to be a great idea, especially since Please improved its Go rules to support library archives saved to non-standard GOPATH locations.

The only additional rule needed is `go_module_download` (and `go_downloaded_source`, but that's just a helper).
It downloads the module source using `go mod download`. The source is then used by a set of generated `go_library` rules.


## `go_get` with `go mod download` - 2020-06-29

Using the Go tooling for downloading turned out to be a dead end,
because in module mode object files can only be built and saved for direct dependencies, but not transitive ones.
(Chances are this is not true, but I couldn't find the right solution)

After looking at [Gazelle](https://github.com/bazelbuild/bazel-gazelle), it turns out that they don't use go modules for package management either.
This is kind of understandable, since go modules break the incremental build idea behind Bazel/Please (every change to `go.mod` would result it massive rebuilds).

Gazelle actually uses a small tool, called [fetch_repo](https://github.com/bazelbuild/bazel-gazelle/tree/5c00b77/cmd/fetch_repo) to download packages.
It supports regular `go get` mode (and more) and works with modules as well. It uses a little trick to do that: it creates a temporary module and downloads the package using `go mod download`.

I was able to integrate this tool into the existing `go_get` rule with a few changes. The working title is `go_module_get`.
If this tool gets integrated into the upstream `go_get` rule, it might worth checking if it can replace the rest of the current download logic.


## First few revisions - 2020 June

The idea is to reuse Go modules and `go.mod` to download and manage dependencies,
therefor I experimented with a couple rules.

The final rule turned out to be less automated than I'd hoped,
but worked well and used caching as much as possible.

Basically a `go_module` rule executed `go mod download` and then built the packages listed in the rule.

This worked quite well with direct dependencies, but unfortunately Go modules doesn't allow building and saving the output of transitive dependencies.
