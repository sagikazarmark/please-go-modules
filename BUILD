subinclude("//build_defs")

moddown_binary("moddown", "0.1.0", visibility = ["PUBLIC"])

github_repo(
    name = "pleasings2",
    repo = "sagikazarmark/mypleasings",
    revision = "0b63e09260f2be22deae7d3896348114da77b6a0",
)

tarball(
    name = "package",
    srcs = [
        "README.md",
        "//cmd/gogetgen",
        "//dist:moddown",
        "//dist:build_file",
        "//build_defs:dist",
    ],
    out = f"gogetgen_{CONFIG.OS}_{CONFIG.ARCH}.tar.gz",
    gzip = True,
    labels = ["manual"],
)

subinclude("///pleasings2//github")

github_release(
    name = "publish",
    assets = [
        "@linux_amd64//:package",
        "@darwin_amd64//:package",
    ],
    labels = ["manual"],
)
