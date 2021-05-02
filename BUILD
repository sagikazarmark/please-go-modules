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

text_file(
    name = "release_notes_template",
    content = """
Add the following to you `tools/BUILD` file:

```
remote_file(
    name = "godeps",
    url = f"https://github.com/sagikazarmark/please-go-modules/releases/download/REPLACE_VERSION/godeps_{CONFIG.HOSTOS}_{CONFIG.HOSTARCH}.tar.gz",
    hashes = [
REPLACE_HASHES    ],
    extract = True,
    exported_files = ["godeps"],
    binary = True,
)
```
""",
)

genrule(
    name = "release_notes",
    srcs = [":release_notes_template"],
    outs = ["release_notes"],
    cmd = [
        "export HASHES=$(cat checksums.txt | cut -f1 -d' ' | sed 's/\\(.*\\)/        \"\\1\",/g' | tr '\\n' '|')",
        "sed \"s/REPLACE_HASHES/$HASHES/g; s/REPLACE_VERSION/$GIT_TAG/g\" \"$SRCS\" | tr '|' '\\n' > \"$OUTS\"",
    ],
    pass_env = ["GIT_TAG"],
    deps = [":artifacts"],
)

subinclude("///pleasings2//github")

github_release(
    name = "publish",
    assets = [":artifacts"],
    labels = ["manual"],
    notes = ":release_notes",
)
