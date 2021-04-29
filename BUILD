github_repo(
    name = "pleasings2",
    repo = "sagikazarmark/mypleasings",
    revision = "09137dd3e633e2c4dc37f8b840e824a9ceb35d3e",
)

tarball(
    name = "artifact",
    srcs = [
        "README.md",
        "//cmd/godeps",
        "//dist:build_file",
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
