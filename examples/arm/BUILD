go_binary(
    name = "example",
    srcs = ["main.go"],
    deps = [
        "//examples/arm/third_party/go:emperror.dev__errors__match",
        "//examples/arm/third_party/go:github.com__containerd__containerd__sys",
        "//examples/arm/third_party/go:github.com__golang__snappy",
        "//examples/arm/third_party/go:github.com__mattn__go-sqlite3",
        "//examples/arm/third_party/go:github.com__opencontainers__runc__libcontainer__system",
        "//examples/arm/third_party/go:golang.org__x__sys__unix",
        "//examples/arm/third_party/go:google.golang.org__grpc",
    ],
)

go_test(
    name = "example_test",
    srcs = [
        "main.go",
        "main_test.go",
    ],
    deps = [
        "//examples/arm/third_party/go:emperror.dev__errors__match",
        "//examples/arm/third_party/go:github.com__containerd__containerd__sys",
        "//examples/arm/third_party/go:github.com__golang__snappy",
        "//examples/arm/third_party/go:github.com__mattn__go-sqlite3",
        "//examples/arm/third_party/go:github.com__opencontainers__runc__libcontainer__system",
        "//examples/arm/third_party/go:github.com__stretchr__testify__assert",
        "//examples/arm/third_party/go:golang.org__x__sys__unix",
        "//examples/arm/third_party/go:google.golang.org__grpc",
    ],
)
