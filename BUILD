subinclude("//build_defs")

go_modules(
    name = "gomod",
)

go_module2(
    name = "gomod2",
)

build_rule(
        name = "teszt",
        srcs = ["go.mod", "go.sum"],
        cmd = ' && '.join([
            "ls -lah",
            "TMPDIR=/tmp $TOOLS mod download -x 2>&1",
            "chmod -R 0755 $TMP_DIR/go",
        ]),
        output_dirs = ["go"],
        tools = [CONFIG.GO_TOOL],
    )

go_binary(
    name = "bin",
    srcs = ["main.go"],
    #deps = [":gomod", ":github.com__pkg__errors"],
    #deps = [":github.com__pkg__errors"],
    deps = [":gomod", "//lib"],
)

go_get(
    name = "errors",
    get = "github.com/pkg/errors",
)
