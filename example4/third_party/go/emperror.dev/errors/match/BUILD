build_rule(
    name = "match",
    cmd = "",
    output_dirs = ["third_party/go/emperror.dev/errors/src/emperror.dev/errors/match"],
    output_is_complete = True,
    tag = "source",
    deps = ["//third_party/go/emperror.dev/errors:_errors#download"],
)

go_library(
    name = "match",
    srcs = [":_match#source"],
    import_path = "emperror.dev/errors/match",
    visibility = ["PUBLIC"],
    deps = ["//third_party/go/emperror.dev/errors"],
)
