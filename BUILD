build_rule(
    name = "wollemi",
    binary = True,
    srcs = [remote_file(
        name = "wollemi",
        _tag = "download",
        url = f"https://github.com/tcncloud/wollemi/releases/download/v0.0.3/wollemi-v0.0.3-{CONFIG.HOSTOS}-{CONFIG.HOSTARCH}.tar.gz"
    )],
    cmd = " && ".join([
        "tar xf $SRCS",
    ]),
    outs = ["wollemi"],
    visibility = ["PUBLIC"],
)

genrule(
    name = 'go.mod',
    outs = ['go.mod'],
    cmd = [
      """cat > $OUTS << EOF
module go.go.go.com/go

go 1.12
EOF
      """
    ],
    labels = [
      'link:plz-out/'
    ]
)
