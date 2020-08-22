go_module_download(
    name = "errors",
    tag = "download",
    module = "emperror.dev/errors",
    version = "v0.7.0",
    sum = "h1:vf0gZ0j9BJCXQYm0M21DPKeXy9PUjvX1YHXhNzIsqOY=",
    visibility = ["//example3/..."],
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

genrule(
    name = "match",
    cmd = 'cp "$PKG/match/match.a" $OUT',
    outs = ["match.a"],
    deps = ["//example3/emperror.dev/errors/match"],
    visibility = ["PUBLIC"],
)
