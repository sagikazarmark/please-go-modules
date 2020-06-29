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

        ":github.com_davecgh_go-spew",
        ":github.com_konsorten_go-windows-terminal-sequences",
        ":github.com_pkg_errors",
        ":github.com_pmezard_go-difflib",
        ":github.com_sirupsen_logrus",
        ":github.com_stretchr_testify",
        ":golang.org_x_sys",
        ":golang.org_x_tools",
        ":logur.dev_adapter_logrus",
        ":logur.dev_logur",
    ],
)

go_getx(
    name="github.com_davecgh_go-spew",
    get="github.com/davecgh/go-spew/...", # MANUAL
    version="v1.1.1",
    sum="h1:vj9j/u1bqnvCEfJOwUhtlOARqs3+rkHYY13jYWTU97c=",
    visibility=["PUBLIC"],
)
go_getx(
    name="github.com_konsorten_go-windows-terminal-sequences",
    get="github.com/konsorten/go-windows-terminal-sequences",
    version="v1.0.3",
    sum="h1:CE8S1cTafDpPvMhIxNJKvHsGVBgn1xWYf1NbHQhywc8=",
    visibility=["PUBLIC"],
)
go_getx(
    name="github.com_pkg_errors",
    get="github.com/pkg/errors",
    version="v0.9.1",
    sum="h1:FEBLx1zS214owpjy7qsBeixbURkuhQAwrK5UwLGTwt4=",
    visibility=["PUBLIC"],
)
go_getx(
    name="github.com_pmezard_go-difflib",
    get="github.com/pmezard/go-difflib/...", # MANUAL
    version="v1.0.0",
    sum="h1:4DBwDE0NGyQoBHbLQYPwSUPoCMWR5BEzIk/f1lZbAQM=",
    visibility=["PUBLIC"],
)
go_getx(
    name="github.com_sirupsen_logrus",
    get="github.com/sirupsen/logrus",
    version="v1.6.0",
    sum="h1:UBcNElsrwanuuMsnGSlYmtmgbb23qDR5dG+6X6Oo89I=",
    visibility=["PUBLIC"],
    deps=[":golang.org_x_sys", ":github.com_pmezard_go-difflib"], # MANUAL
)
go_getx(
    name="github.com_stretchr_testify",
    get="github.com/stretchr/testify",
    version="v1.2.2",
    sum="h1:bSDNvY7ZPG5RlJ8otE/7V6gMiyenm9RtJ7IUVIAoJ1w=",
    visibility=["PUBLIC"],
    deps=[":github.com_pmezard_go-difflib", ":github.com_davecgh_go-spew", ":github.com_stretchr_objx"], # MANUAL
)
go_getx(
    name="golang.org_x_sys",
    get="golang.org/x/sys/...", # MANUAL
    version="v0.0.0-20190422165155-953cdadca894",
    sum="h1:Cz4ceDQGXuKRnVBDTS23GTn/pU5OE2C0WrNTOYK1Uuc=",
    visibility=["PUBLIC"],
)
go_getx(
    name="golang.org_x_tools",
    get="golang.org/x/tools", # MANUAL
    install=["go/vcs"], # MANUAL
    version="v0.0.0-20200626171337-aa94e735be7f",
    sum="h1:JcoF/bowzCDI+MXu1yLqQGNO3ibqWsWq+Sk7pOT218w=",
    visibility=["PUBLIC"],
)
go_getx(
    name="logur.dev_adapter_logrus",
    get="logur.dev/adapter/logrus",
    version="v0.5.0",
    sum="h1:cxsiceNXQLTKBk0keASgKAvrw9zzKa/XPE0Bn8tHXFI=",
    visibility=["PUBLIC"],
    deps=[":github.com_sirupsen_logrus", ":logur.dev_logur"], # MANUAL
)
go_getx(
    name="logur.dev_logur",
    get="logur.dev/logur",
    version="v0.16.2",
    sum="h1:q4MxivaiTXiDHrQyeCH5WkwBLUrd6rM2lZlyztYvi4o=",
    visibility=["PUBLIC"],
)


# MANUAL
go_getx(
    name="github.com_stretchr_objx",
    get="github.com/stretchr/objx/...",
    version="v0.1.0",
    sum="h1:4G4v2dO3VZwixGIRoQ5Lfboy6nUhCyYzaqnIAPPhYs4=",
)
