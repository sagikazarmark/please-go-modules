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
        go_mod_dep(":module", "github.com/pkg/errors"),
        go_mod_dep(":module", "logur.dev/logur"),
        go_mod_dep(":module", "logur.dev/logur/logtesting"),
        go_mod_dep(":module", "logur.dev/logur/conformance"),
        "//lib",
    ],
)
