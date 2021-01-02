go_binary(
    name = "example",
    srcs = ["main.go"],
    deps = [
        "//example2/third_party/go/emperror.dev/errors/match",
        "//example2/third_party/go/github.com/containerd/containerd/sys",
        "//example2/third_party/go/github.com/golang/snappy",
        "//example2/third_party/go/github.com/mattn/go-sqlite3",
        "//example2/third_party/go/github.com/opencontainers/runc/libcontainer/system",
        "//example2/third_party/go/golang.org/x/sys/unix",
        "//example2/third_party/go/google.golang.org/grpc",
    ],
)

go_test(
    name = "example_test",
    srcs = ["main_test.go"],
    deps = [
        "//example2/third_party/go/github.com/stretchr/testify/assert",
    ],
)
