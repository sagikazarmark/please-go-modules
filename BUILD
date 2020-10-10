subinclude("//build_defs")

moddown_binary("moddown", "0.3.0", visibility = ["PUBLIC"])

github_repo(
    name = "pleasings2",
    repo = "sagikazarmark/mypleasings",
    revision = "f8a12721c6f929db3e227e07c152d428ac47ab1b",
)

tarball(
    name = "artifact",
    srcs = [
        "README.md",
        "//cmd/godeps",
        "//dist:moddown",
        "//dist:build_file",
        "//build_defs:dist",
    ],
    out = f"godeps_{CONFIG.OS}_{CONFIG.ARCH}.tar.gz",
    gzip = True,
    labels = ["manual"],
)

subinclude("///pleasings2//misc")

build_artifacts(
    name = "artifacts",
    artifacts = [
        "@linux_amd64//:artifact",
        "@darwin_amd64//:artifact",
    ],
    labels = ["manual"],
)

subinclude("///pleasings2//github")

github_release(
    name = "publish",
    assets = [":artifacts"],
    labels = ["manual"],
)
