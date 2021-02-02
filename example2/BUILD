go_binary(
    name = "example",
    srcs = ["main.go"],
    deps = [
        "//example2/third_party/go:emperror.dev_errors_match",
        "//example2/third_party/go:github.com_containerd_containerd_sys",
        "//example2/third_party/go:github.com_golang_snappy",
        "//example2/third_party/go:github.com_mattn_go-sqlite3",
        "//example2/third_party/go:github.com_opencontainers_runc_libcontainer_system",
        "//example2/third_party/go:golang.org_x_sys_unix",
        "//example2/third_party/go:google.golang.org_grpc",
    ],
)

go_test(
    name = "example_test",
    srcs = ["main_test.go"],
    deps = [
        "//example2/third_party/go/github.com/stretchr/testify/assert",
    ],
)
