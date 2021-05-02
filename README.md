# Please Go Modules rule generator

[![GitHub Workflow Status](https://img.shields.io/github/workflow/status/sagikazarmark/please-go-modules/CI?style=flat-square)](https://github.com/sagikazarmark/please-go-modules/actions?query=workflow%3ACI)
![Please Version](https://img.shields.io/badge/please%20version-%3E=16.0.1-B879FF.svg?style=flat-square)


## Usage

### Generate `BUILD` file from `go.mod`

Add the following snippet to your `tools/BUILD` file in the root of your repository:

```starlark
remote_file(
    name = "godeps",
    url = f"https://github.com/sagikazarmark/please-go-modules/releases/latest/download/godeps_{CONFIG.HOSTOS}_{CONFIG.HOSTARCH}.tar.gz",
    extract = True,
    exported_files = ["godeps"],
    binary = True,
)
```

Add the following snippet to your `.plzconfig`:

```
[please]
version = 16.0.1

[alias "godeps"]
desc = Generate third-party dependency rules for a Go project
cmd = run //tools:godeps -- -dir third_party/go -clean -builtin
```

Run the following:

```bash
plz godeps
```

The above command will generate build targets in `third_party/go` for your third party dependencies.


### Update BUILD files to use dependencies

You can combine the above with [wollemi](https://github.com/tcncloud/wollemi) that can generate/update
`BUILD` files in your project to use third-party dependencies.

Add the following content to your `tools/BUILD` file:

```starlark
go_toolchain(
    name = "go_toolchain",
    version = "1.16.3",
)

WOLLEMI_VERSION = "v0.7.0"
remote_file(
    name = "wollemi",
    url = f"https://github.com/tcncloud/wollemi/releases/download/{WOLLEMI_VERSION}/wollemi-{WOLLEMI_VERSION}-{CONFIG.HOSTOS}-{CONFIG.HOSTARCH}.tar.gz",
    extract = True,
    exported_files = ["wollemi"],
    binary = True,
)

sh_cmd(
    name = "plz-tidy",
    cmd = [
        "export GOROOT=\\\\$($(out_exe :go_toolchain|go) env GOROOT)",
        "$(out_exe :godeps) -dir third_party/go -clean -builtin -wollemi",
        "$(out_exe :wollemi) gofmt ./...",
    ],
    deps = [
        ":godeps",
        ":wollemi",
        ":go_toolchain",
    ],
)
```

**Note:** You can remove any references to `go_toolchain` if you want to use Go installed on your system.

Finally, add an alias:

```
[alias "tidy"]
desc = Tidy generates build targets for dependencies and makes sure that BUILD files are up-to-date.
cmd = run //tools:plz-tidy
```

and run:

```bash
plz tidy
```

**Note:** the wollemi command might not work perfectly with Go submodules. You need to run wollemi for each module separately.
