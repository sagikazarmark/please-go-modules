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
        "//third_party/go:github.com__scylladb__go-set__strset",
    ],
)

go_test(
    name = "test",
    srcs = glob(["*.go"]),
    data = [
        "//pkg/depgraph/testdata",
        "//pkg/depgraph/testdata/packages",
    ],
    deps = [
        "//pkg/golist",
        "//pkg/sumfile",
        "//third_party/go:github.com__scylladb__go-set__strset",
    ],
)
