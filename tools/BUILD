go_toolchain(
    name = "go_toolchain",
    architectures = [
        "linux_amd64",
        "darwin_amd64",
        "darwin_arm64",
    ],
    version = CONFIG.GO_VERSION,
)

remote_file(
    name = "wollemi",
    binary = True,
    exported_files = ["wollemi"],
    extract = True,
    url = f"https://github.com/tcncloud/wollemi/releases/download/v{CONFIG.WOLLEMI_VERSION}/wollemi-v{CONFIG.WOLLEMI_VERSION}-{CONFIG.HOSTOS}-{CONFIG.HOSTARCH}.tar.gz",
)

sh_cmd(
    name = "plz-tidy",
    cmd = [
        "export GOROOT=\\\\$($(out_exe :go_toolchain|go) env GOROOT)",
        "$(out_exe //cmd/godeps) -dir third_party/go -clean -builtin -wollemi",
        "$(out_exe :wollemi) gofmt ./...",
    ],
    deps = [
        ":go_toolchain",
        ":wollemi",
        "//cmd/godeps",
    ],
)
