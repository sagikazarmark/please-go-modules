go_module_download(
    name = "multierr",
    tag = "download",
    module = "go.uber.org/multierr",
    version = "v1.5.0",
    sum = "h1:KCa4XfM8CWFCpxXRGok+Q0SS/0XBhMDbHHGABQLvD2A=",
    visibility = ["//example3/go.uber.org/multierr/..."],
)

go_library(
    name = "multierr",
    srcs = [go_downloaded_source("multierr", ":_multierr#download")],
    visibility = ["PUBLIC"],
    deps = ["//example3/go.uber.org/atomic"],
    _import_path = "go.uber.org/multierr",
)
