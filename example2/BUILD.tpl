subinclude("//build_defs")

go_binary(
    name = "example",
    srcs = ["main.go"],
    deps = [
        ":emperror.dev_errors_match",
        ":github.com_golang_snappy",
        ":github.com_mattn_go-sqlite3",
        ":golang.org_x_sys_unix",
        ":google.golang.org_grpc",
    ],
)
