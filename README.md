# Please Go Modules experiment

[![GitHub Workflow Status](https://img.shields.io/github/workflow/status/sagikazarmark/please-go-modules/CI?style=flat-square)](https://github.com/sagikazarmark/please-go-modules/actions?query=workflow%3ACI)
![Please Version](https://img.shields.io/badge/please%20version-%3E=15.16.0-B879FF.svg?style=flat-square)

Experiments with [Please](https://please.build) and Go modules.


## Usage

### Generate `BUILD` file from `go.mod`

Add the following snippet to your `BUILD` file in the root of your repository:

```starlark
http_archive(
    name = "pleasegomod",
    urls = [f"https://github.com/sagikazarmark/please-go-modules/releases/download/v0.0.18/godeps_{CONFIG.HOSTOS}_{CONFIG.HOSTARCH}.tar.gz"],
)
```

Add the following snippet to your `.plzconfig` (or to a `.plzconfig.experimental` file and use `--profile experimental` in Please commands):

```
[please]
version = 15.16.0

[featureflags]
PleaseGoInstall = true

[alias "godeps"]
desc = Generate third-party dependency rules for a Go project
cmd = run ///pleasegomod//:godeps -- -dir third_party/go -clean -builtin
```

Run the following:

```bash
plz godeps
```

The above command will generate build targets in `third_party/go` for your third party dependencies.


### Update BUILD files to use dependencies

You can combine the above with [wollemi](https://github.com/tcncloud/wollemi) that can generate/update
`BUILD` files in your project to use third-party dependencies.

Add the following to the `BUILD` file in the project root:

```starlark
github_repo(
    name = "pleasings2",
    repo = "sagikazarmark/mypleasings",
    revision = "f644350aab5e0090ab69022beeb2feadfcdc223e",
)
```

Then create a `tools/BUILD` file with the following content:

```starlark
subinclude("///pleasings2//go:tools")

# Copied here from https://github.com/sagikazarmark/mypleasings/blob/27b6451ea99d160aec03f242be5261978770b4e1/tools/go/BUILD
# Custom go toolchain doesn't work with subrepos: https://github.com/thought-machine/please/issues/1547
wollemi_wrapper(
    name = "wollemi-wrapper",
    binary = "///pleasings2//tools/go:wollemi",
    labels = ["go"],
    visibility = ["PUBLIC"],
)

sh_cmd(
    name = "plz-tidy",
    cmd = [
        "$(out_exe ///pleasegomod//:godeps) -dir third_party/go -clean -builtin -wollemi",
        "$(out_exe :wollemi-wrapper) gofmt ./...",
    ],
    deps = [
        "///pleasegomod//:godeps",
        ":wollemi-wrapper",
    ],
)
```

(**Note:** the wollemi command might not work perfectly with Go submodules. You need to run wollemi for each module separately)

Finally, add an alias:

```
[alias "tidy"]
desc = Tidy generates build targets for dependencies and makes sure that BUILD files are up-to-date.
cmd = run //tools:plz-tidy
```


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
