genrule(
    name = "errors",
    cmd = 'cp "$PKG/errors/errors.a" $OUT',
    outs = ["errors.a"],
    deps = ["//example3/github.com/pkg/errors"],
    visibility = ["PUBLIC"],
)
