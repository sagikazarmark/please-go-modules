#go_module(
#    name = "module",
#    install = [
#        "github.com/pkg/errors",
#        "logur.dev/logur",
#        "logur.dev/logur/logtesting",
#        "logur.dev/logur/conformance",
#    ],
#)

go_binary(
    name = "bin",
    srcs = ["main.go"],
    deps = [
        "//lib",
        #go_mod_dep(":module", "github.com/pkg/errors"),
        #go_mod_dep(":module", "logur.dev/logur"),
        #go_mod_dep(":module", "logur.dev/logur/logtesting"),
        #go_mod_dep(":module", "logur.dev/logur/conformance"),

        ":github.com_pkg_errors",
        ":github.com_sirupsen_logrus",
        ":logur.dev_adapter_logrus",
        ":logur.dev_logur",
    ],
)

go_getx(
    name = "github.com_sirupsen_logrus",
    get = "github.com/sirupsen/logrus",
    version = "v1.6.0",
    sum = "h1:UBcNElsrwanuuMsnGSlYmtmgbb23qDR5dG+6X6Oo89I=",
    install = [""],
    deps = [":golang.org_x_sys"],
    visibility=["PUBLIC"],
)

go_getx(
    name = "logur.dev_logur",
    get = "logur.dev/logur",
    version = "v0.16.2",
    sum = "h1:q4MxivaiTXiDHrQyeCH5WkwBLUrd6rM2lZlyztYvi4o=",
    install = [""],
    deps = [],
    visibility=["PUBLIC"],
)

go_getx(
    name = "logur.dev_adapter_logrus",
    get = "logur.dev/adapter/logrus",
    version = "v0.5.0",
    sum = "h1:cxsiceNXQLTKBk0keASgKAvrw9zzKa/XPE0Bn8tHXFI=",
    install = [""],
    deps = [":github.com_sirupsen_logrus", ":logur.dev_logur"],
    visibility=["PUBLIC"],
)

go_getx(
    name = "golang.org_x_tools",
    get = "golang.org/x/tools",
    version = "v0.0.0-20200626171337-aa94e735be7f",
    sum = "h1:JcoF/bowzCDI+MXu1yLqQGNO3ibqWsWq+Sk7pOT218w=",
    install = ["go/vcs"],
    deps = [],
    visibility=["PUBLIC"],
)

go_getx(
    name = "github.com_pkg_errors",
    get = "github.com/pkg/errors",
    version = "v0.9.1",
    sum = "h1:FEBLx1zS214owpjy7qsBeixbURkuhQAwrK5UwLGTwt4=",
    install = [""],
    deps = [],
    visibility=["PUBLIC"],
)

go_getx(
    name = "golang.org_x_sys",
    get = "golang.org/x/sys",
    version = "v0.0.0-20190422165155-953cdadca894",
    sum = "h1:Cz4ceDQGXuKRnVBDTS23GTn/pU5OE2C0WrNTOYK1Uuc=",
    install = ["unix"],
    deps = [],
    visibility=["PUBLIC"],
)
