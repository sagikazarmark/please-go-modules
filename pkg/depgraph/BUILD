go_library(
    name = "depgraph",
    srcs = glob(
        ["*.go"],
        exclude = ["*_test.go"],
    ),
    visibility = ["PUBLIC"],
    deps = [
        "//pkg/golist",
        "//pkg/sumfile",
        "//third_party/go/github.com/scylladb/go-set/strset",
    ],
)
