genrule(
    name = "source",
    cmd = "",
    deps = ["//example3/emperror.dev/errors:_errors#download"],
    output_dirs = ["example3/emperror.dev/errors/src/emperror.dev/errors/match"],
)

go_library(
    name = "match",
    srcs = [":source"],
    visibility = ["PUBLIC"],
    deps = ["//example3/emperror.dev/errors"],
    _import_path = "emperror.dev/errors/match",
)
