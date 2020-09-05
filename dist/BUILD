subinclude("//build_defs")

moddown_binary("moddown", "0.1.0", visibility = ["PUBLIC"], host = False)

genrule(
    name = "build_file",
    srcs = ["build_file"],
    outs = ["BUILD"],
    cmd = 'cp "$SRCS" "$OUTS"',
    visibility = ["PUBLIC"],
)
