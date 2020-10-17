subinclude("//build_defs")

moddown_binary(
    "moddown",
    CONFIG.MODDOWN_VERSION,
    host = False,
    visibility = ["PUBLIC"],
)

genrule(
    name = "build_file",
    srcs = ["build_file"],
    outs = ["BUILD"],
    cmd = 'cp "$SRCS" "$OUTS"',
    visibility = ["PUBLIC"],
)
