go_library(
    name = "sumfile",
    srcs = glob(["*.go"], exclude = ["*_test.go"]),
    visibility = ["PUBLIC"],
    deps = [],
)

go_test(
    name = "sumfile_test",
    srcs = glob(["*_test.go"]),
    deps = [":sumfile"],
)
