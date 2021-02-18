go_binary(
    name = "example",
    srcs = ["main.go"],
    deps = [
        "//example3/third_party/go:emperror.dev__errors__match",
        "//example3/third_party/go:github.com__containerd__containerd__sys",
        "//example3/third_party/go:github.com__golang__snappy",
        "//example3/third_party/go:github.com__mattn__go-sqlite3",
        "//example3/third_party/go:github.com__opencontainers__runc__libcontainer__system",
        "//example3/third_party/go:golang.org__x__sys__unix",
        "//example3/third_party/go:google.golang.org__grpc",
    ],
)

go_test(
    name = "example_test",
    srcs = ["main_test.go"],
    deps = ["//example3/third_party/go:github.com__stretchr__testify__assert"],
)
