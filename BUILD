subinclude("//build_defs")

moddown_binary("moddown", "0.3.0", visibility = ["PUBLIC"])

github_repo(
    name = "pleasings2",
    repo = "sagikazarmark/mypleasings",
    revision = "4c40fa674130e6d92bcdb4ef9bd17954fdbf3fab",
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

sha256sum(
    name = "checksums.txt",
    srcs = [
        "@linux_amd64//:artifact",
        "@darwin_amd64//:artifact",
    ],
    out = "checksums.txt",
    labels = ["manual"],
)

filegroup(
    name = "artifacts",
    srcs = [
        "@linux_amd64//:artifact",
        "@darwin_amd64//:artifact",
        ":checksums.txt",
    ],
    labels = ["manual"],
)

subinclude("///pleasings2//github")

github_release(
    name = "publish",
    assets = [
        "@linux_amd64//:artifact",
        "@darwin_amd64//:artifact",
        ":checksums.txt",
    ],
    labels = ["manual"],
)
