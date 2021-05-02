go_toolchain(
    name = "go_toolchain",
    version = CONFIG.GO_VERSION,
    architectures = [
        "linux_amd64",
        "darwin_amd64",
    ],
)

remote_file(
    name = "wollemi",
    url = f"https://github.com/tcncloud/wollemi/releases/download/v{CONFIG.WOLLEMI_VERSION}/wollemi-v{CONFIG.WOLLEMI_VERSION}-{CONFIG.HOSTOS}-{CONFIG.HOSTARCH}.tar.gz",
    extract = True,
    exported_files = ["wollemi"],
    binary = True,
)

sh_cmd(
    name = "plz-tidy",
    cmd = [
        "export GOROOT=\\\\$($(out_exe :go_toolchain|go) env GOROOT)",
        "$(out_exe //cmd/godeps) -dir third_party/go -clean -builtin -wollemi",
        "$(out_exe :wollemi) gofmt ./...",
    ],
    deps = [
        "//cmd/godeps",
        ":wollemi",
        ":go_toolchain",
    ],
)
