subinclude("//build_defs")

moddown_binary(
    "moddown",
    CONFIG.MODDOWN_VERSION,
    visibility = ["PUBLIC"],
)

github_repo(
    name = "pleasings2",
    repo = "sagikazarmark/mypleasings",
    revision = "3dea080ba8070af5a9dac6e382703f862e97a4f0",
)

tarball(
    name = "artifact",
    srcs = [
        "README.md",
        "//build_defs:dist",
        "//cmd/godeps",
        "//dist:build_file",
        "//dist:moddown",
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
