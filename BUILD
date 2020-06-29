go_module(
    name = "module",
    install = [
        "github.com/pkg/errors",
        "logur.dev/logur",
        "logur.dev/logur/logtesting",
        "logur.dev/logur/conformance",
        "logur.dev/adapter/logrus",
        "github.com/sirupsen/logrus",
        "golang.org/x/sys",
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
        go_mod_dep(":module", "logur.dev/adapter/logrus"),
        go_mod_dep(":module", "github.com/sirupsen/logrus"),
        go_mod_dep(":module", "golang.org/x/sys"),
        "//lib",
    ],
)
