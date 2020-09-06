subinclude("//build_defs")

go_library(
    name = "external_test_issue",
    srcs = ["pkg.go"],
)

go_test(
    name = "test",
    srcs = ["pkg_test.go"],
    deps = [":external_test_issue"]
)

go_library(
    name = "external_test_issue_test_external",
    srcs = ["pkg.go", "pkg_test.go"],
    import_path = CONFIG.GO_IMPORT_PATH + "/" + package_name(),
)

go_test(
    name = "test_external",
    srcs = ["pkg_external_test.go", "pkg_test.go"],
    external = True,
    deps = [":external_test_issue"],
    labels = ["manual"],
)

go_test(
    name = "test_external2",
    srcs = ["pkg_external_test.go"],
    external = True,
    deps = [":external_test_issue_test_external"],
)
