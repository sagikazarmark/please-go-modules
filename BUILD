subinclude("//build_defs")

moddown_binary("moddown", "0.3.0", visibility = ["PUBLIC"])

github_repo(
    name = "pleasings2",
    repo = "sagikazarmark/mypleasings",
    revision = "67e70f85bf5f8a8c9e5e78f13b3729e56836d209",
)

tarball(
    name = "artifact",
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
