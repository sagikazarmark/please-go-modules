go_module_download(
    name = "errors",
    tag = "download",
    module = "emperror.dev/errors",
    version = "v0.7.0",
    sum = "h1:vf0gZ0j9BJCXQYm0M21DPKeXy9PUjvX1YHXhNzIsqOY=",
    visibility = ["//example3/..."],
)

build_rule(
    name = "errors",
    tag = "source",
    cmd = "",
    deps = ["//example3:_errors#download"],
    output_dirs = ["example3/src/emperror.dev/errors"],
)

go_library(
    name = "errors",
    srcs = [go_downloaded_source("errors", ":_errors#download")],
    visibility = ["PUBLIC"],
    deps = [
        "//example3/github.com/pkg:errors",
        "//example3/go.uber.org:multierr",
    ],
)
