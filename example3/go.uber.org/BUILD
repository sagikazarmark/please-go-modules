go_module_download(
    name = "multierr",
    tag = "download",
    module = "go.uber.org/multierr",
    version = "v1.5.0",
    sum = "h1:KCa4XfM8CWFCpxXRGok+Q0SS/0XBhMDbHHGABQLvD2A=",
    visibility = ["//example3/go.uber.org/..."],
)

go_library(
    name = "multierr",
    srcs = [go_downloaded_source("multierr", ":_multierr#download")],
    visibility = ["PUBLIC"],
    deps = [":atomic"],
)

go_module_download(
    name = "atomic",
    tag = "download",
    module = "go.uber.org/atomic",
    version = "v1.6.0",
    sum = "h1:Ezj3JGmsOnG1MoRWQkPBsKLe9DwWD9QeXzTRzzldNVk=",
    visibility = ["//example3/go.uber.org/..."],
)

go_library(
    name = "atomic",
    srcs = [go_downloaded_source("atomic", ":_atomic#download")],
    visibility = ["PUBLIC"],
)
