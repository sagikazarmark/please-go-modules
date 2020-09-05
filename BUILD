subinclude("//build_defs")

moddown_binary("moddown", "0.1.0", visibility = ["PUBLIC"])

github_repo(
    name = "pleasings2",
    repo = "sagikazarmark/mypleasings",
    revision = "e4cc66bc0cd5b2bc86fc9bd058319a5c864c4261",
)

tarball(
    name = "package",
    srcs = [
        "README.md",
        "//cmd/gogetgen",
        "//dist:moddown",
        "//dist:build_file",
        "//build_defs:go_module",
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
