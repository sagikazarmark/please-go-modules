go_module_download(
    name = "errors",
    tag = "download",
    module = "github.com/pkg/errors",
    version = "v0.9.1",
    sum = "h1:FEBLx1zS214owpjy7qsBeixbURkuhQAwrK5UwLGTwt4=",
    visibility = ["//example3/github.com/pkg/..."],
)

go_library(
    name = "errors",
    srcs = [go_downloaded_source("errors", ":_errors#download")],
    visibility = ["PUBLIC"],
)
