go_module_download(
    name = "atomic",
    module = "go.uber.org/atomic",
    sum = "h1:OI5t8sDa1Or+q8AeE+yKeB/SDYioSHAgcVljj9JIETY=",
    tag = "download",
    version = "v1.5.0",
    visibility = ["//third_party/go/go.uber.org/atomic/..."],
)

go_downloaded_source("atomic", ":_atomic#download")

go_library(
    name = "atomic",
    srcs = [":_atomic#source"],
    import_path = "go.uber.org/atomic",
    visibility = ["PUBLIC"],
    deps = [],
)
