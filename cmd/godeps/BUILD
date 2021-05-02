go_binary(
    name = "godeps",
    srcs = glob(["*.go"], exclude = ["*_test.go"]),
    visibility = ["PUBLIC"],
    deps = [
        "//pkg/depgraph",
        "//pkg/golist",
        "//pkg/sumfile",
        "//third_party/go:github.com__bazelbuild__buildtools__build",
        "//third_party/go:github.com__scylladb__go-set__strset",
    ],
)
