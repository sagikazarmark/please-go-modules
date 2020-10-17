go_binary(
    name = "godeps",
    srcs = ["main.go", "platform.go"],
    visibility = ["PUBLIC"],
    deps = [
        "//pkg/golist",
        "//pkg/depgraph",
        "//pkg/sumfile",
        "//third_party/go/github.com/bazelbuild/buildtools/build",
    ],
)
