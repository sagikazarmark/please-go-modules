genrule(
    name = "errors",
    cmd = 'cp "$PKG/errors/errors.a" $OUT',
    outs = ["errors.a"],
    deps = ["//example3/emperror.dev/errors"],
    visibility = ["PUBLIC"],
)
