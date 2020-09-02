go_module_download(
    name = "multierr",
    module = "go.uber.org/multierr",
    sum = "h1:f3WCSC2KzAcBXGATIxAB1E2XuCpNU255wNKZ505qi3E=",
    tag = "download",
    version = "v1.4.0",
    visibility = ["//third_party/go/go.uber.org/multierr/..."],
)

go_downloaded_source("multierr", ":_multierr#download")

go_library(
    name = "multierr",
    srcs = [":_multierr#source"],
    _import_path = "go.uber.org/multierr",
    visibility = ["PUBLIC"],
    deps = ["//third_party/go/go.uber.org/atomic"],
)
