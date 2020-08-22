genrule(
    name = "atomic",
    cmd = 'cp "$PKG/atomic/atomic.a" $OUT',
    outs = ["atomic.a"],
    deps = ["//example3/go.uber.org/atomic"],
    visibility = ["PUBLIC"],
)

genrule(
    name = "multierr",
    cmd = 'cp "$PKG/multierr/multierr.a" $OUT',
    outs = ["multierr.a"],
    deps = ["//example3/go.uber.org/multierr"],
    visibility = ["PUBLIC"],
)
