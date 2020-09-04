moddown_binary(
    "moddown",
    "0.1.0",
    visibility = ["PUBLIC"],
)

go_binary(
    name = "example4",
    srcs = ["main.go"],
    deps = [
        "//example/third_party/go/emperror.dev/errors/match",
        "//example/third_party/go/github.com/golang/snappy",
        "//example/third_party/go/github.com/mattn/go-sqlite3",
    ],
)
