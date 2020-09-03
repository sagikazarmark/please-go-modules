go_module_download(
    name = "atomic",
    tag = "download",
    module = "go.uber.org/atomic",
    version = "v1.6.0",
    sum = "h1:Ezj3JGmsOnG1MoRWQkPBsKLe9DwWD9QeXzTRzzldNVk=",
    visibility = ["//example3/go.uber.org/atomic/..."],
)

go_library(
    name = "atomic",
    srcs = [go_downloaded_source("atomic", ":_atomic#download")],
    visibility = ["PUBLIC"],
    _import_path = "go.uber.org/atomic",
)
