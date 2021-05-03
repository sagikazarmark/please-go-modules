go_binary(
    name = "example",
    srcs = ["main.go"],
    deps = [
        "//examples/noexpand/third_party/go:emperror.dev__errors",
        "//examples/noexpand/third_party/go:github.com__containerd__containerd",
        "//examples/noexpand/third_party/go:github.com__golang__snappy",
        "//examples/noexpand/third_party/go:github.com__mattn__go-sqlite3",
        "//examples/noexpand/third_party/go:github.com__opencontainers__runc",
        "//examples/noexpand/third_party/go:golang.org__x__sys",
        "//examples/noexpand/third_party/go:google.golang.org__grpc",
    ],
)

go_test(
    name = "example_test",
    srcs = [
        "main.go",
        "main_test.go",
    ],
    deps = [
        "//examples/noexpand/third_party/go:emperror.dev__errors",
        "//examples/noexpand/third_party/go:github.com__containerd__containerd",
        "//examples/noexpand/third_party/go:github.com__golang__snappy",
        "//examples/noexpand/third_party/go:github.com__mattn__go-sqlite3",
        "//examples/noexpand/third_party/go:github.com__opencontainers__runc",
        "//examples/noexpand/third_party/go:github.com__stretchr__testify",
        "//examples/noexpand/third_party/go:golang.org__x__sys",
        "//examples/noexpand/third_party/go:google.golang.org__grpc",
    ],
)
