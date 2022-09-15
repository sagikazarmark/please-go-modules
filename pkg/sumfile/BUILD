go_library(
    name = "sumfile",
    srcs = glob(
        ["*.go"],
        exclude = ["*_test.go"],
    ),
    visibility = ["PUBLIC"],
)

go_test(
    name = "sumfile_test",
    srcs = glob(["*_test.go"]),
    external = True,
    deps = [":sumfile"],
)
