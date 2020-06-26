subinclude("//build_defs")

go_module(
    name = "module",
    install = [
        "github.com/pkg/errors",
        "logur.dev/logur",
        "logur.dev/logur/logtesting",
        "logur.dev/logur/conformance",
    ],
)

go_binary(
    name = "bin",
    srcs = ["main.go"],
    deps = [
        ":_module#install_github.com_pkg_errors",
        ":_module#install_logur.dev_logur",
        ":_module#install_logur.dev_logur_logtesting",
        ":_module#install_logur.dev_logur_conformance",
        "//lib",
    ],
)
