go_binary(
    name = "example",
    srcs = ["main.go"],
    deps = [
        "//example/lib",
        "//example/third_party/go/github.com/pkg:errors",
        "//example/third_party/go/github.com/sirupsen:logrus",
        "//example/third_party/go/logur.dev:logur",
        "//example/third_party/go/logur.dev/adapter:logrus",
    ],
)

go_module_download(
    name = "test_download",
    module = "github.com/sagikazarmark/functional-options",
    version = "v0.0.0-20200123060825-f70eefebf3bc",
    sum = "h1:KkEl73Fv7raGkzw5vxw7fQkyuW51joXsNojPlAkRRQw=",
)

go_module_compile(
    name = "test_compile",
    module = "github.com/sagikazarmark/functional-options",
    src = ":test_download",
)
