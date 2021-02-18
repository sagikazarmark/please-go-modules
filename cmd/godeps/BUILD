go_binary(
    name = "godeps",
    srcs = [
        "generate.go",
        "generate_builtin.go",
        "main.go",
        "os_config.go",
        "platform.go",
    ],
    visibility = ["PUBLIC"],
    deps = [
        "//pkg/depgraph",
        "//pkg/golist",
        "//pkg/sumfile",
        "//third_party/go:github.com__bazelbuild__buildtools__build",
    ],
)
