# Please Go Modules rule generator

[![GitHub Workflow Status](https://img.shields.io/github/workflow/status/sagikazarmark/please-go-modules/CI?style=flat-square)](https://github.com/sagikazarmark/please-go-modules/actions?query=workflow%3ACI)
![Please Version](https://img.shields.io/badge/please%20version-%3E=16.0.0-B879FF.svg?style=flat-square)


## Usage

### Generate `BUILD` file from `go.mod`

Add the following snippet to your `BUILD` file in the root of your repository:

```starlark
http_archive(
    name = "pleasegomod",
    urls = [f"https://github.com/sagikazarmark/please-go-modules/releases/download/v0.0.19/godeps_{CONFIG.HOSTOS}_{CONFIG.HOSTARCH}.tar.gz"],
)
```

Add the following snippet to your `.plzconfig`:

```
[please]
version = 16.0.0

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
