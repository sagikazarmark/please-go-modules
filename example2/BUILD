subinclude("//build_defs")

go_binary(
    name = "example",
    srcs = ["main.go"],
    deps = [
        ":emperror.dev_errors_match",
        ":github.com_golang_snappy",
        ":github.com_mattn_go-sqlite3",
        ":golang.org_x_sys_unix",
        ":google.golang.org_grpc",
    ],
)
# __config

config_setting(
    name = "__config_linux_amd64",
    values = {
        "os": "linux",
        "cpu": "amd64",
    },
)

config_setting(
    name = "__config_darwin_amd64",
    values = {
        "os": "darwin",
        "cpu": "amd64",
    },
)


# emperror.dev/errors

go_module_download(
    name = "emperror.dev_errors",
    module = "emperror.dev/errors",
    sum = "h1:vf0gZ0j9BJCXQYm0M21DPKeXy9PUjvX1YHXhNzIsqOY=",
    tag = "download",
    version = "v0.7.0",
)

fileexport(
    name = "emperror.dev_errors",
    srcs = [
        "error_details.go",
        "errors.go",
        "errors_multi.go",
        "stack.go",
        "wrap.go",
        "wrap_go1_13.go",
    ],
    tag = "go_source",
    deps = [":_emperror.dev_errors#download"],
)

go_library(
    name = "emperror.dev_errors",
    srcs = [":_emperror.dev_errors#go_source"],
    cover = False,
    import_path = "emperror.dev/errors",
    visibility = ["PUBLIC"],
    deps = [
        ":github.com_pkg_errors",
        ":go.uber.org_multierr",
    ],
)


# emperror.dev/errors/match

fileexport(
    name = "emperror.dev_errors_match",
    srcs = ["match/match.go"],
    tag = "go_source",
    deps = [":_emperror.dev_errors#download"],
)

go_library(
    name = "emperror.dev_errors_match",
    srcs = [":_emperror.dev_errors_match#go_source"],
    cover = False,
    import_path = "emperror.dev/errors/match",
    visibility = ["PUBLIC"],
    deps = [":emperror.dev_errors"],
)


# github.com/golang/protobuf

go_module_download(
    name = "github.com_golang_protobuf",
    module = "github.com/golang/protobuf",
    sum = "h1:gyjaxf+svBWX08ZjK86iN9geUJF0H6gp2IRKX6Nf6/I=",
    tag = "download",
    version = "v1.3.3",
)


# github.com/golang/protobuf/proto

fileexport(
    name = "github.com_golang_protobuf_proto",
    srcs = [
        "proto/clone.go",
        "proto/decode.go",
        "proto/deprecated.go",
        "proto/discard.go",
        "proto/encode.go",
        "proto/equal.go",
        "proto/extensions.go",
        "proto/lib.go",
        "proto/message_set.go",
        "proto/pointer_unsafe.go",
        "proto/properties.go",
        "proto/table_marshal.go",
        "proto/table_merge.go",
        "proto/table_unmarshal.go",
        "proto/text.go",
        "proto/text_parser.go",
    ],
    tag = "go_source",
    deps = [":_github.com_golang_protobuf#download"],
)

go_library(
    name = "github.com_golang_protobuf_proto",
    srcs = [":_github.com_golang_protobuf_proto#go_source"],
    cover = False,
    import_path = "github.com/golang/protobuf/proto",
    visibility = ["PUBLIC"],
    deps = [],
)


# github.com/golang/protobuf/ptypes

fileexport(
    name = "github.com_golang_protobuf_ptypes",
    srcs = [
        "ptypes/any.go",
        "ptypes/doc.go",
        "ptypes/duration.go",
        "ptypes/timestamp.go",
    ],
    tag = "go_source",
    deps = [":_github.com_golang_protobuf#download"],
)

go_library(
    name = "github.com_golang_protobuf_ptypes",
    srcs = [":_github.com_golang_protobuf_ptypes#go_source"],
    cover = False,
    import_path = "github.com/golang/protobuf/ptypes",
    visibility = ["PUBLIC"],
    deps = [
        ":github.com_golang_protobuf_proto",
        ":github.com_golang_protobuf_ptypes_any",
        ":github.com_golang_protobuf_ptypes_duration",
        ":github.com_golang_protobuf_ptypes_timestamp",
    ],
)


# github.com/golang/protobuf/ptypes/any

fileexport(
    name = "github.com_golang_protobuf_ptypes_any",
    srcs = ["ptypes/any/any.pb.go"],
    tag = "go_source",
    deps = [":_github.com_golang_protobuf#download"],
)

go_library(
    name = "github.com_golang_protobuf_ptypes_any",
    srcs = [":_github.com_golang_protobuf_ptypes_any#go_source"],
    cover = False,
    import_path = "github.com/golang/protobuf/ptypes/any",
    visibility = ["PUBLIC"],
    deps = [":github.com_golang_protobuf_proto"],
)


# github.com/golang/protobuf/ptypes/duration

fileexport(
    name = "github.com_golang_protobuf_ptypes_duration",
    srcs = ["ptypes/duration/duration.pb.go"],
    tag = "go_source",
    deps = [":_github.com_golang_protobuf#download"],
)

go_library(
    name = "github.com_golang_protobuf_ptypes_duration",
    srcs = [":_github.com_golang_protobuf_ptypes_duration#go_source"],
    cover = False,
    import_path = "github.com/golang/protobuf/ptypes/duration",
    visibility = ["PUBLIC"],
    deps = [":github.com_golang_protobuf_proto"],
)


# github.com/golang/protobuf/ptypes/timestamp

fileexport(
    name = "github.com_golang_protobuf_ptypes_timestamp",
    srcs = ["ptypes/timestamp/timestamp.pb.go"],
    tag = "go_source",
    deps = [":_github.com_golang_protobuf#download"],
)

go_library(
    name = "github.com_golang_protobuf_ptypes_timestamp",
    srcs = [":_github.com_golang_protobuf_ptypes_timestamp#go_source"],
    cover = False,
    import_path = "github.com/golang/protobuf/ptypes/timestamp",
    visibility = ["PUBLIC"],
    deps = [":github.com_golang_protobuf_proto"],
)


# github.com/golang/snappy

go_module_download(
    name = "github.com_golang_snappy",
    module = "github.com/golang/snappy",
    sum = "h1:Qgr9rKW7uDUkrbSmQeiDsGa8SjGyCOGtuasMWwvp2P4=",
    tag = "download",
    version = "v0.0.1",
)

fileexport(
    name = "github.com_golang_snappy",
    srcs = [
        "decode.go",
        "decode_amd64.go",
        "encode.go",
        "encode_amd64.go",
        "snappy.go",
    ],
    tag = "go_source",
    deps = [":_github.com_golang_snappy#download"],
)

fileexport(
    name = "github.com_golang_snappy",
    srcs = [
        "decode_amd64.s",
        "encode_amd64.s",
    ],
    tag = "s_source",
    deps = [":_github.com_golang_snappy#download"],
)

go_library(
    name = "github.com_golang_snappy",
    srcs = [":_github.com_golang_snappy#go_source"],
    asm_srcs = [":_github.com_golang_snappy#s_source"],
    cover = False,
    import_path = "github.com/golang/snappy",
    visibility = ["PUBLIC"],
    deps = [],
)


# github.com/mattn/go-sqlite3

go_module_download(
    name = "github.com_mattn_go-sqlite3",
    module = "github.com/mattn/go-sqlite3",
    sum = "h1:A2EQLwjYf/hfYaM20FVjs1UewCTTFR7RmjEHkLjldIA=",
    tag = "download",
    version = "v1.14.2",
)

fileexport(
    name = "github.com_mattn_go-sqlite3",
    srcs = [
        "convert.go",
        "doc.go",
        "sqlite3_func_crypt.go",
        "sqlite3_go18.go",
        "sqlite3_opt_preupdate.go",
        "sqlite3_opt_preupdate_omit.go",
    ],
    tag = "go_source",
    deps = [":_github.com_mattn_go-sqlite3#download"],
)

fileexport(
    name = "github.com_mattn_go-sqlite3",
    srcs = [
        "backup.go",
        "callback.go",
        "error.go",
        "sqlite3.go",
        "sqlite3_context.go",
        "sqlite3_load_extension.go",
        "sqlite3_opt_userauth_omit.go",
        "sqlite3_other.go",
        "sqlite3_type.go",
    ],
    tag = "cgo_source",
    deps = [":_github.com_mattn_go-sqlite3#download"],
)

fileexport(
    name = "github.com_mattn_go-sqlite3",
    srcs = [
        "sqlite3-binding.c",
        "sqlite3_opt_unlock_notify.c",
    ],
    tag = "c_source",
    deps = [":_github.com_mattn_go-sqlite3#download"],
)

fileexport(
    name = "github.com_mattn_go-sqlite3",
    srcs = [
        "sqlite3-binding.h",
        "sqlite3ext.h",
    ],
    tag = "h_source",
    deps = [":_github.com_mattn_go-sqlite3#download"],
)

cgo_library(
    name = "github.com_mattn_go-sqlite3",
    srcs = [":_github.com_mattn_go-sqlite3#cgo_source"],
    hdrs = [":_github.com_mattn_go-sqlite3#h_source"],
    c_srcs = [":_github.com_mattn_go-sqlite3#c_source"],
    compiler_flags = [
        "-DHAVE_USLEEP=1",
        "-DSQLITE_DEFAULT_WAL_SYNCHRONOUS=1",
        "-DSQLITE_ENABLE_UPDATE_DELETE_LIMIT",
        "-DSQLITE_ENABLE_FTS3_PARENTHESIS",
        "-DSQLITE_OMIT_DEPRECATED",
        "-DSQLITE_ENABLE_FTS3",
        "-I ${PKG}",
        "-std=gnu99",
        "-DSQLITE_ENABLE_RTREE",
        "-DSQLITE_DISABLE_INTRINSIC",
        "-DSQLITE_ENABLE_FTS4_UNICODE61",
        "-DSQLITE_TRACE_SIZE_LIMIT=15",
        "-DSQLITE_THREADSAFE=1",
        "-Wno-deprecated-declarations",
    ] + select({
        ":__config_darwin_amd64": [],
        ":__config_linux_amd64": [
            "-DHAVE_PWRITE64=1",
            "-DHAVE_PREAD64=1",
        ],
        "default": [],
    }),
    go_srcs = [":_github.com_mattn_go-sqlite3#go_source"],
    import_path = "github.com/mattn/go-sqlite3",
    linker_flags = select({
        ":__config_darwin_amd64": [],
        ":__config_linux_amd64": ["-ldl"],
        "default": [],
    }),
    visibility = ["PUBLIC"],
    deps = [],
)


# github.com/pkg/errors

go_module_download(
    name = "github.com_pkg_errors",
    module = "github.com/pkg/errors",
    sum = "h1:J8lpUdobwIeCI7OiSxHqEwJUKvJwicL5+3v1oe2Yb4k=",
    tag = "download",
    version = "v0.9.0",
)

fileexport(
    name = "github.com_pkg_errors",
    srcs = [
        "errors.go",
        "go113.go",
        "stack.go",
    ],
    tag = "go_source",
    deps = [":_github.com_pkg_errors#download"],
)

go_library(
    name = "github.com_pkg_errors",
    srcs = [":_github.com_pkg_errors#go_source"],
    cover = False,
    import_path = "github.com/pkg/errors",
    visibility = ["PUBLIC"],
    deps = [],
)


# go.uber.org/atomic

go_module_download(
    name = "go.uber.org_atomic",
    module = "go.uber.org/atomic",
    sum = "h1:OI5t8sDa1Or+q8AeE+yKeB/SDYioSHAgcVljj9JIETY=",
    tag = "download",
    version = "v1.5.0",
)

fileexport(
    name = "go.uber.org_atomic",
    srcs = [
        "atomic.go",
        "error.go",
        "string.go",
    ],
    tag = "go_source",
    deps = [":_go.uber.org_atomic#download"],
)

go_library(
    name = "go.uber.org_atomic",
    srcs = [":_go.uber.org_atomic#go_source"],
    cover = False,
    import_path = "go.uber.org/atomic",
    visibility = ["PUBLIC"],
    deps = [],
)


# go.uber.org/multierr

go_module_download(
    name = "go.uber.org_multierr",
    module = "go.uber.org/multierr",
    sum = "h1:f3WCSC2KzAcBXGATIxAB1E2XuCpNU255wNKZ505qi3E=",
    tag = "download",
    version = "v1.4.0",
)

fileexport(
    name = "go.uber.org_multierr",
    srcs = [
        "error.go",
        "go113.go",
    ],
    tag = "go_source",
    deps = [":_go.uber.org_multierr#download"],
)

go_library(
    name = "go.uber.org_multierr",
    srcs = [":_go.uber.org_multierr#go_source"],
    cover = False,
    import_path = "go.uber.org/multierr",
    visibility = ["PUBLIC"],
    deps = [":go.uber.org_atomic"],
)


# golang.org/x/net

go_module_download(
    name = "golang.org_x_net",
    module = "golang.org/x/net",
    sum = "h1:3G+cUijn7XD+S4eJFddp53Pv7+slrESplyjG25HgL+k=",
    tag = "download",
    version = "v0.0.0-20200324143707-d3edc9973b7e",
)


# golang.org/x/net/http/httpguts

fileexport(
    name = "golang.org_x_net_http_httpguts",
    srcs = [
        "http/httpguts/guts.go",
        "http/httpguts/httplex.go",
    ],
    tag = "go_source",
    deps = [":_golang.org_x_net#download"],
)

go_library(
    name = "golang.org_x_net_http_httpguts",
    srcs = [":_golang.org_x_net_http_httpguts#go_source"],
    cover = False,
    import_path = "golang.org/x/net/http/httpguts",
    visibility = ["PUBLIC"],
    deps = [":golang.org_x_net_idna"],
)


# golang.org/x/net/http2

fileexport(
    name = "golang.org_x_net_http2",
    srcs = [
        "http2/ciphers.go",
        "http2/client_conn_pool.go",
        "http2/databuffer.go",
        "http2/errors.go",
        "http2/flow.go",
        "http2/frame.go",
        "http2/go111.go",
        "http2/gotrack.go",
        "http2/headermap.go",
        "http2/http2.go",
        "http2/pipe.go",
        "http2/server.go",
        "http2/transport.go",
        "http2/write.go",
        "http2/writesched.go",
        "http2/writesched_priority.go",
        "http2/writesched_random.go",
    ],
    tag = "go_source",
    deps = [":_golang.org_x_net#download"],
)

go_library(
    name = "golang.org_x_net_http2",
    srcs = [":_golang.org_x_net_http2#go_source"],
    cover = False,
    import_path = "golang.org/x/net/http2",
    visibility = ["PUBLIC"],
    deps = [
        ":golang.org_x_net_http2_hpack",
        ":golang.org_x_net_http_httpguts",
        ":golang.org_x_net_idna",
    ],
)


# golang.org/x/net/http2/hpack

fileexport(
    name = "golang.org_x_net_http2_hpack",
    srcs = [
        "http2/hpack/encode.go",
        "http2/hpack/hpack.go",
        "http2/hpack/huffman.go",
        "http2/hpack/tables.go",
    ],
    tag = "go_source",
    deps = [":_golang.org_x_net#download"],
)

go_library(
    name = "golang.org_x_net_http2_hpack",
    srcs = [":_golang.org_x_net_http2_hpack#go_source"],
    cover = False,
    import_path = "golang.org/x/net/http2/hpack",
    visibility = ["PUBLIC"],
    deps = [],
)


# golang.org/x/net/idna

fileexport(
    name = "golang.org_x_net_idna",
    srcs = [
        "idna/idna10.0.0.go",
        "idna/punycode.go",
        "idna/tables12.00.go",
        "idna/trie.go",
        "idna/trieval.go",
    ],
    tag = "go_source",
    deps = [":_golang.org_x_net#download"],
)

go_library(
    name = "golang.org_x_net_idna",
    srcs = [":_golang.org_x_net_idna#go_source"],
    cover = False,
    import_path = "golang.org/x/net/idna",
    visibility = ["PUBLIC"],
    deps = [
        ":golang.org_x_text_secure_bidirule",
        ":golang.org_x_text_unicode_bidi",
        ":golang.org_x_text_unicode_norm",
    ],
)


# golang.org/x/net/internal/timeseries

fileexport(
    name = "golang.org_x_net_internal_timeseries",
    srcs = ["internal/timeseries/timeseries.go"],
    tag = "go_source",
    deps = [":_golang.org_x_net#download"],
)

go_library(
    name = "golang.org_x_net_internal_timeseries",
    srcs = [":_golang.org_x_net_internal_timeseries#go_source"],
    cover = False,
    import_path = "golang.org/x/net/internal/timeseries",
    visibility = ["PUBLIC"],
    deps = [],
)


# golang.org/x/net/trace

fileexport(
    name = "golang.org_x_net_trace",
    srcs = [
        "trace/events.go",
        "trace/histogram.go",
        "trace/trace.go",
    ],
    tag = "go_source",
    deps = [":_golang.org_x_net#download"],
)

go_library(
    name = "golang.org_x_net_trace",
    srcs = [":_golang.org_x_net_trace#go_source"],
    cover = False,
    import_path = "golang.org/x/net/trace",
    visibility = ["PUBLIC"],
    deps = [":golang.org_x_net_internal_timeseries"],
)


# golang.org/x/sys

go_module_download(
    name = "golang.org_x_sys",
    module = "golang.org/x/sys",
    sum = "h1:Fqb3ao1hUmOR3GkUOg/Y+BadLwykBIzs5q8Ez2SbHyc=",
    tag = "download",
    version = "v0.0.0-20200905004654-be1d3432aa8f",
)


# golang.org/x/sys/internal/unsafeheader

fileexport(
    name = "golang.org_x_sys_internal_unsafeheader",
    srcs = ["internal/unsafeheader/unsafeheader.go"],
    tag = "go_source",
    deps = [":_golang.org_x_sys#download"],
)

go_library(
    name = "golang.org_x_sys_internal_unsafeheader",
    srcs = [":_golang.org_x_sys_internal_unsafeheader#go_source"],
    cover = False,
    import_path = "golang.org/x/sys/internal/unsafeheader",
    visibility = ["PUBLIC"],
    deps = [],
)


# golang.org/x/sys/unix

fileexport(
    name = "golang.org_x_sys_unix",
    srcs = [
        "unix/aliases.go",
        "unix/ioctl.go",
        "unix/syscall.go",
        "unix/fdset.go",
        "unix/endian_little.go",
        "unix/timestruct.go",
        "unix/syscall_unix_gc.go",
        "unix/sockcmsg_unix.go",
        "unix/dirent.go",
        "unix/sockcmsg_unix_other.go",
        "unix/race0.go",
        "unix/syscall_unix.go",
        "unix/pagesize_unix.go",
        "unix/str.go",
        "unix/env_unix.go",
        "unix/constants.go",
    ] + select({
        ":__config_darwin_amd64": [
            "unix/readdirent_getdirentries.go",
            "unix/syscall_darwin.go",
            "unix/ztypes_darwin_amd64.go",
            "unix/syscall_darwin_libSystem.go",
            "unix/syscall_darwin.1_13.go",
            "unix/fcntl_darwin.go",
            "unix/dev_darwin.go",
            "unix/zerrors_darwin_amd64.go",
            "unix/zsyscall_darwin_amd64.1_13.go",
            "unix/zsyscall_darwin_amd64.go",
            "unix/syscall_darwin_amd64.go",
            "unix/syscall_bsd.go",
        ],
        ":__config_linux_amd64": [
            "unix/syscall_linux_amd64.go",
            "unix/zerrors_linux_amd64.go",
            "unix/ztypes_linux.go",
            "unix/zsyscall_linux_amd64.go",
            "unix/readdirent_getdents.go",
            "unix/affinity_linux.go",
            "unix/sockcmsg_linux.go",
            "unix/ztypes_linux_amd64.go",
            "unix/fcntl.go",
            "unix/zerrors_linux.go",
            "unix/syscall_linux_amd64_gc.go",
            "unix/syscall_linux_gc.go",
            "unix/syscall_linux.go",
            "unix/dev_linux.go",
            "unix/zsyscall_linux.go",
            "unix/bluetooth_linux.go",
            "unix/zsysnum_linux_amd64.go",
            "unix/zptrace_x86_linux.go",
        ],
        "default": [],
    }),
    tag = "go_source",
    deps = [":_golang.org_x_sys#download"],
)

fileexport(
    name = "golang.org_x_sys_unix",
    srcs = select({
        ":__config_darwin_amd64": [
            "unix/asm_darwin_amd64.s",
            "unix/zsyscall_darwin_amd64.1_13.s",
            "unix/zsyscall_darwin_amd64.s",
        ],
        ":__config_linux_amd64": ["unix/asm_linux_amd64.s"],
        "default": [],
    }),
    tag = "s_source",
    deps = [":_golang.org_x_sys#download"],
)

go_library(
    name = "golang.org_x_sys_unix",
    srcs = [":_golang.org_x_sys_unix#go_source"],
    asm_srcs = [":_golang.org_x_sys_unix#s_source"],
    cover = False,
    import_path = "golang.org/x/sys/unix",
    visibility = ["PUBLIC"],
    deps = [":golang.org_x_sys_internal_unsafeheader"],
)


# golang.org/x/text

go_module_download(
    name = "golang.org_x_text",
    module = "golang.org/x/text",
    sum = "h1:g61tztE5qeGQ89tm6NTjjM9VPIm088od1l6aSorWRWg=",
    tag = "download",
    version = "v0.3.0",
)


# golang.org/x/text/secure/bidirule

fileexport(
    name = "golang.org_x_text_secure_bidirule",
    srcs = [
        "secure/bidirule/bidirule.go",
        "secure/bidirule/bidirule10.0.0.go",
    ],
    tag = "go_source",
    deps = [":_golang.org_x_text#download"],
)

go_library(
    name = "golang.org_x_text_secure_bidirule",
    srcs = [":_golang.org_x_text_secure_bidirule#go_source"],
    cover = False,
    import_path = "golang.org/x/text/secure/bidirule",
    visibility = ["PUBLIC"],
    deps = [
        ":golang.org_x_text_transform",
        ":golang.org_x_text_unicode_bidi",
    ],
)


# golang.org/x/text/transform

fileexport(
    name = "golang.org_x_text_transform",
    srcs = ["transform/transform.go"],
    tag = "go_source",
    deps = [":_golang.org_x_text#download"],
)

go_library(
    name = "golang.org_x_text_transform",
    srcs = [":_golang.org_x_text_transform#go_source"],
    cover = False,
    import_path = "golang.org/x/text/transform",
    visibility = ["PUBLIC"],
    deps = [],
)


# golang.org/x/text/unicode/bidi

fileexport(
    name = "golang.org_x_text_unicode_bidi",
    srcs = [
        "unicode/bidi/bidi.go",
        "unicode/bidi/bracket.go",
        "unicode/bidi/core.go",
        "unicode/bidi/prop.go",
        "unicode/bidi/tables10.0.0.go",
        "unicode/bidi/trieval.go",
    ],
    tag = "go_source",
    deps = [":_golang.org_x_text#download"],
)

go_library(
    name = "golang.org_x_text_unicode_bidi",
    srcs = [":_golang.org_x_text_unicode_bidi#go_source"],
    cover = False,
    import_path = "golang.org/x/text/unicode/bidi",
    visibility = ["PUBLIC"],
    deps = [],
)


# golang.org/x/text/unicode/norm

fileexport(
    name = "golang.org_x_text_unicode_norm",
    srcs = [
        "unicode/norm/composition.go",
        "unicode/norm/forminfo.go",
        "unicode/norm/input.go",
        "unicode/norm/iter.go",
        "unicode/norm/normalize.go",
        "unicode/norm/readwriter.go",
        "unicode/norm/tables10.0.0.go",
        "unicode/norm/transform.go",
        "unicode/norm/trie.go",
    ],
    tag = "go_source",
    deps = [":_golang.org_x_text#download"],
)

go_library(
    name = "golang.org_x_text_unicode_norm",
    srcs = [":_golang.org_x_text_unicode_norm#go_source"],
    cover = False,
    import_path = "golang.org/x/text/unicode/norm",
    visibility = ["PUBLIC"],
    deps = [":golang.org_x_text_transform"],
)


# google.golang.org/genproto

go_module_download(
    name = "google.golang.org_genproto",
    module = "google.golang.org/genproto",
    sum = "h1:gSJIx1SDwno+2ElGhA4+qG2zF97qiUzTM+rQ0klBOcE=",
    tag = "download",
    version = "v0.0.0-20190819201941-24fa4b261c55",
)


# google.golang.org/genproto/googleapis/rpc/status

fileexport(
    name = "google.golang.org_genproto_googleapis_rpc_status",
    srcs = ["googleapis/rpc/status/status.pb.go"],
    tag = "go_source",
    deps = [":_google.golang.org_genproto#download"],
)

go_library(
    name = "google.golang.org_genproto_googleapis_rpc_status",
    srcs = [":_google.golang.org_genproto_googleapis_rpc_status#go_source"],
    cover = False,
    import_path = "google.golang.org/genproto/googleapis/rpc/status",
    visibility = ["PUBLIC"],
    deps = [
        ":github.com_golang_protobuf_proto",
        ":github.com_golang_protobuf_ptypes_any",
    ],
)


# google.golang.org/grpc

go_module_download(
    name = "google.golang.org_grpc",
    module = "google.golang.org/grpc",
    sum = "h1:zWTV+LMdc3kaiJMSTOFz2UgSBgx8RNQoTGiZu3fR9S0=",
    tag = "download",
    version = "v1.32.0",
)

fileexport(
    name = "google.golang.org_grpc",
    srcs = [
        "backoff.go",
        "balancer_conn_wrappers.go",
        "call.go",
        "clientconn.go",
        "codec.go",
        "dialoptions.go",
        "doc.go",
        "interceptor.go",
        "picker_wrapper.go",
        "pickfirst.go",
        "preloader.go",
        "proxy.go",
        "resolver_conn_wrapper.go",
        "rpc_util.go",
        "server.go",
        "service_config.go",
        "stream.go",
        "trace.go",
        "version.go",
    ],
    tag = "go_source",
    deps = [":_google.golang.org_grpc#download"],
)

go_library(
    name = "google.golang.org_grpc",
    srcs = [":_google.golang.org_grpc#go_source"],
    cover = False,
    import_path = "google.golang.org/grpc",
    visibility = ["PUBLIC"],
    deps = [
        ":golang.org_x_net_trace",
        ":google.golang.org_grpc_backoff",
        ":google.golang.org_grpc_balancer",
        ":google.golang.org_grpc_balancer_base",
        ":google.golang.org_grpc_balancer_roundrobin",
        ":google.golang.org_grpc_codes",
        ":google.golang.org_grpc_connectivity",
        ":google.golang.org_grpc_credentials",
        ":google.golang.org_grpc_encoding",
        ":google.golang.org_grpc_encoding_proto",
        ":google.golang.org_grpc_grpclog",
        ":google.golang.org_grpc_internal",
        ":google.golang.org_grpc_internal_backoff",
        ":google.golang.org_grpc_internal_balancerload",
        ":google.golang.org_grpc_internal_binarylog",
        ":google.golang.org_grpc_internal_buffer",
        ":google.golang.org_grpc_internal_channelz",
        ":google.golang.org_grpc_internal_envconfig",
        ":google.golang.org_grpc_internal_grpcrand",
        ":google.golang.org_grpc_internal_grpcsync",
        ":google.golang.org_grpc_internal_grpcutil",
        ":google.golang.org_grpc_internal_resolver_dns",
        ":google.golang.org_grpc_internal_resolver_passthrough",
        ":google.golang.org_grpc_internal_serviceconfig",
        ":google.golang.org_grpc_internal_transport",
        ":google.golang.org_grpc_keepalive",
        ":google.golang.org_grpc_metadata",
        ":google.golang.org_grpc_peer",
        ":google.golang.org_grpc_resolver",
        ":google.golang.org_grpc_serviceconfig",
        ":google.golang.org_grpc_stats",
        ":google.golang.org_grpc_status",
        ":google.golang.org_grpc_tap",
    ],
)


# google.golang.org/grpc/attributes

fileexport(
    name = "google.golang.org_grpc_attributes",
    srcs = ["attributes/attributes.go"],
    tag = "go_source",
    deps = [":_google.golang.org_grpc#download"],
)

go_library(
    name = "google.golang.org_grpc_attributes",
    srcs = [":_google.golang.org_grpc_attributes#go_source"],
    cover = False,
    import_path = "google.golang.org/grpc/attributes",
    visibility = ["PUBLIC"],
    deps = [],
)


# google.golang.org/grpc/backoff

fileexport(
    name = "google.golang.org_grpc_backoff",
    srcs = ["backoff/backoff.go"],
    tag = "go_source",
    deps = [":_google.golang.org_grpc#download"],
)

go_library(
    name = "google.golang.org_grpc_backoff",
    srcs = [":_google.golang.org_grpc_backoff#go_source"],
    cover = False,
    import_path = "google.golang.org/grpc/backoff",
    visibility = ["PUBLIC"],
    deps = [],
)


# google.golang.org/grpc/balancer

fileexport(
    name = "google.golang.org_grpc_balancer",
    srcs = ["balancer/balancer.go"],
    tag = "go_source",
    deps = [":_google.golang.org_grpc#download"],
)

go_library(
    name = "google.golang.org_grpc_balancer",
    srcs = [":_google.golang.org_grpc_balancer#go_source"],
    cover = False,
    import_path = "google.golang.org/grpc/balancer",
    visibility = ["PUBLIC"],
    deps = [
        ":google.golang.org_grpc_connectivity",
        ":google.golang.org_grpc_credentials",
        ":google.golang.org_grpc_internal",
        ":google.golang.org_grpc_metadata",
        ":google.golang.org_grpc_resolver",
        ":google.golang.org_grpc_serviceconfig",
    ],
)


# google.golang.org/grpc/balancer/base

fileexport(
    name = "google.golang.org_grpc_balancer_base",
    srcs = [
        "balancer/base/balancer.go",
        "balancer/base/base.go",
    ],
    tag = "go_source",
    deps = [":_google.golang.org_grpc#download"],
)

go_library(
    name = "google.golang.org_grpc_balancer_base",
    srcs = [":_google.golang.org_grpc_balancer_base#go_source"],
    cover = False,
    import_path = "google.golang.org/grpc/balancer/base",
    visibility = ["PUBLIC"],
    deps = [
        ":google.golang.org_grpc_balancer",
        ":google.golang.org_grpc_connectivity",
        ":google.golang.org_grpc_grpclog",
        ":google.golang.org_grpc_resolver",
    ],
)


# google.golang.org/grpc/balancer/grpclb/state

fileexport(
    name = "google.golang.org_grpc_balancer_grpclb_state",
    srcs = ["balancer/grpclb/state/state.go"],
    tag = "go_source",
    deps = [":_google.golang.org_grpc#download"],
)

go_library(
    name = "google.golang.org_grpc_balancer_grpclb_state",
    srcs = [":_google.golang.org_grpc_balancer_grpclb_state#go_source"],
    cover = False,
    import_path = "google.golang.org/grpc/balancer/grpclb/state",
    visibility = ["PUBLIC"],
    deps = [":google.golang.org_grpc_resolver"],
)


# google.golang.org/grpc/balancer/roundrobin

fileexport(
    name = "google.golang.org_grpc_balancer_roundrobin",
    srcs = ["balancer/roundrobin/roundrobin.go"],
    tag = "go_source",
    deps = [":_google.golang.org_grpc#download"],
)

go_library(
    name = "google.golang.org_grpc_balancer_roundrobin",
    srcs = [":_google.golang.org_grpc_balancer_roundrobin#go_source"],
    cover = False,
    import_path = "google.golang.org/grpc/balancer/roundrobin",
    visibility = ["PUBLIC"],
    deps = [
        ":google.golang.org_grpc_balancer",
        ":google.golang.org_grpc_balancer_base",
        ":google.golang.org_grpc_grpclog",
        ":google.golang.org_grpc_internal_grpcrand",
    ],
)


# google.golang.org/grpc/binarylog/grpc_binarylog_v1

fileexport(
    name = "google.golang.org_grpc_binarylog_grpc_binarylog_v1",
    srcs = ["binarylog/grpc_binarylog_v1/binarylog.pb.go"],
    tag = "go_source",
    deps = [":_google.golang.org_grpc#download"],
)

go_library(
    name = "google.golang.org_grpc_binarylog_grpc_binarylog_v1",
    srcs = [":_google.golang.org_grpc_binarylog_grpc_binarylog_v1#go_source"],
    cover = False,
    import_path = "google.golang.org/grpc/binarylog/grpc_binarylog_v1",
    visibility = ["PUBLIC"],
    deps = [
        ":github.com_golang_protobuf_proto",
        ":github.com_golang_protobuf_ptypes_duration",
        ":github.com_golang_protobuf_ptypes_timestamp",
    ],
)


# google.golang.org/grpc/codes

fileexport(
    name = "google.golang.org_grpc_codes",
    srcs = [
        "codes/code_string.go",
        "codes/codes.go",
    ],
    tag = "go_source",
    deps = [":_google.golang.org_grpc#download"],
)

go_library(
    name = "google.golang.org_grpc_codes",
    srcs = [":_google.golang.org_grpc_codes#go_source"],
    cover = False,
    import_path = "google.golang.org/grpc/codes",
    visibility = ["PUBLIC"],
    deps = [],
)


# google.golang.org/grpc/connectivity

fileexport(
    name = "google.golang.org_grpc_connectivity",
    srcs = ["connectivity/connectivity.go"],
    tag = "go_source",
    deps = [":_google.golang.org_grpc#download"],
)

go_library(
    name = "google.golang.org_grpc_connectivity",
    srcs = [":_google.golang.org_grpc_connectivity#go_source"],
    cover = False,
    import_path = "google.golang.org/grpc/connectivity",
    visibility = ["PUBLIC"],
    deps = [":google.golang.org_grpc_grpclog"],
)


# google.golang.org/grpc/credentials

fileexport(
    name = "google.golang.org_grpc_credentials",
    srcs = [
        "credentials/credentials.go",
        "credentials/go12.go",
        "credentials/tls.go",
    ],
    tag = "go_source",
    deps = [":_google.golang.org_grpc#download"],
)

go_library(
    name = "google.golang.org_grpc_credentials",
    srcs = [":_google.golang.org_grpc_credentials#go_source"],
    cover = False,
    import_path = "google.golang.org/grpc/credentials",
    visibility = ["PUBLIC"],
    deps = [
        ":github.com_golang_protobuf_proto",
        ":google.golang.org_grpc_attributes",
        ":google.golang.org_grpc_internal",
        ":google.golang.org_grpc_internal_credentials",
    ],
)


# google.golang.org/grpc/encoding

fileexport(
    name = "google.golang.org_grpc_encoding",
    srcs = ["encoding/encoding.go"],
    tag = "go_source",
    deps = [":_google.golang.org_grpc#download"],
)

go_library(
    name = "google.golang.org_grpc_encoding",
    srcs = [":_google.golang.org_grpc_encoding#go_source"],
    cover = False,
    import_path = "google.golang.org/grpc/encoding",
    visibility = ["PUBLIC"],
    deps = [],
)


# google.golang.org/grpc/encoding/proto

fileexport(
    name = "google.golang.org_grpc_encoding_proto",
    srcs = ["encoding/proto/proto.go"],
    tag = "go_source",
    deps = [":_google.golang.org_grpc#download"],
)

go_library(
    name = "google.golang.org_grpc_encoding_proto",
    srcs = [":_google.golang.org_grpc_encoding_proto#go_source"],
    cover = False,
    import_path = "google.golang.org/grpc/encoding/proto",
    visibility = ["PUBLIC"],
    deps = [
        ":github.com_golang_protobuf_proto",
        ":google.golang.org_grpc_encoding",
    ],
)


# google.golang.org/grpc/grpclog

fileexport(
    name = "google.golang.org_grpc_grpclog",
    srcs = [
        "grpclog/component.go",
        "grpclog/grpclog.go",
        "grpclog/logger.go",
        "grpclog/loggerv2.go",
    ],
    tag = "go_source",
    deps = [":_google.golang.org_grpc#download"],
)

go_library(
    name = "google.golang.org_grpc_grpclog",
    srcs = [":_google.golang.org_grpc_grpclog#go_source"],
    cover = False,
    import_path = "google.golang.org/grpc/grpclog",
    visibility = ["PUBLIC"],
    deps = [":google.golang.org_grpc_internal_grpclog"],
)


# google.golang.org/grpc/internal

fileexport(
    name = "google.golang.org_grpc_internal",
    srcs = ["internal/internal.go"],
    tag = "go_source",
    deps = [":_google.golang.org_grpc#download"],
)

go_library(
    name = "google.golang.org_grpc_internal",
    srcs = [":_google.golang.org_grpc_internal#go_source"],
    cover = False,
    import_path = "google.golang.org/grpc/internal",
    visibility = ["PUBLIC"],
    deps = [
        ":google.golang.org_grpc_connectivity",
        ":google.golang.org_grpc_serviceconfig",
    ],
)


# google.golang.org/grpc/internal/backoff

fileexport(
    name = "google.golang.org_grpc_internal_backoff",
    srcs = ["internal/backoff/backoff.go"],
    tag = "go_source",
    deps = [":_google.golang.org_grpc#download"],
)

go_library(
    name = "google.golang.org_grpc_internal_backoff",
    srcs = [":_google.golang.org_grpc_internal_backoff#go_source"],
    cover = False,
    import_path = "google.golang.org/grpc/internal/backoff",
    visibility = ["PUBLIC"],
    deps = [
        ":google.golang.org_grpc_backoff",
        ":google.golang.org_grpc_internal_grpcrand",
    ],
)


# google.golang.org/grpc/internal/balancerload

fileexport(
    name = "google.golang.org_grpc_internal_balancerload",
    srcs = ["internal/balancerload/load.go"],
    tag = "go_source",
    deps = [":_google.golang.org_grpc#download"],
)

go_library(
    name = "google.golang.org_grpc_internal_balancerload",
    srcs = [":_google.golang.org_grpc_internal_balancerload#go_source"],
    cover = False,
    import_path = "google.golang.org/grpc/internal/balancerload",
    visibility = ["PUBLIC"],
    deps = [":google.golang.org_grpc_metadata"],
)


# google.golang.org/grpc/internal/binarylog

fileexport(
    name = "google.golang.org_grpc_internal_binarylog",
    srcs = [
        "internal/binarylog/binarylog.go",
        "internal/binarylog/binarylog_testutil.go",
        "internal/binarylog/env_config.go",
        "internal/binarylog/method_logger.go",
        "internal/binarylog/sink.go",
    ],
    tag = "go_source",
    deps = [":_google.golang.org_grpc#download"],
)

go_library(
    name = "google.golang.org_grpc_internal_binarylog",
    srcs = [":_google.golang.org_grpc_internal_binarylog#go_source"],
    cover = False,
    import_path = "google.golang.org/grpc/internal/binarylog",
    visibility = ["PUBLIC"],
    deps = [
        ":github.com_golang_protobuf_proto",
        ":github.com_golang_protobuf_ptypes",
        ":google.golang.org_grpc_binarylog_grpc_binarylog_v1",
        ":google.golang.org_grpc_grpclog",
        ":google.golang.org_grpc_internal_grpcutil",
        ":google.golang.org_grpc_metadata",
        ":google.golang.org_grpc_status",
    ],
)


# google.golang.org/grpc/internal/buffer

fileexport(
    name = "google.golang.org_grpc_internal_buffer",
    srcs = ["internal/buffer/unbounded.go"],
    tag = "go_source",
    deps = [":_google.golang.org_grpc#download"],
)

go_library(
    name = "google.golang.org_grpc_internal_buffer",
    srcs = [":_google.golang.org_grpc_internal_buffer#go_source"],
    cover = False,
    import_path = "google.golang.org/grpc/internal/buffer",
    visibility = ["PUBLIC"],
    deps = [],
)


# google.golang.org/grpc/internal/channelz

fileexport(
    name = "google.golang.org_grpc_internal_channelz",
    srcs = [
        "internal/channelz/funcs.go",
        "internal/channelz/logging.go",
        "internal/channelz/types.go",
    ] + select({
        ":__config_darwin_amd64": [
            "internal/channelz/types_nonlinux.go",
            "internal/channelz/util_nonlinux.go",
        ],
        ":__config_linux_amd64": [
            "internal/channelz/types_linux.go",
            "internal/channelz/util_linux.go",
        ],
        "default": [],
    }),
    tag = "go_source",
    deps = [":_google.golang.org_grpc#download"],
)

go_library(
    name = "google.golang.org_grpc_internal_channelz",
    srcs = [":_google.golang.org_grpc_internal_channelz#go_source"],
    cover = False,
    import_path = "google.golang.org/grpc/internal/channelz",
    visibility = ["PUBLIC"],
    deps = [
        ":google.golang.org_grpc_connectivity",
        ":google.golang.org_grpc_credentials",
        ":google.golang.org_grpc_grpclog",
    ] + select({
        ":__config_darwin_amd64": [],
        ":__config_linux_amd64": [":golang.org_x_sys_unix"],
        "default": [],
    }),
)


# google.golang.org/grpc/internal/credentials

fileexport(
    name = "google.golang.org_grpc_internal_credentials",
    srcs = [
        "internal/credentials/spiffe.go",
        "internal/credentials/syscallconn.go",
        "internal/credentials/util.go",
    ],
    tag = "go_source",
    deps = [":_google.golang.org_grpc#download"],
)

go_library(
    name = "google.golang.org_grpc_internal_credentials",
    srcs = [":_google.golang.org_grpc_internal_credentials#go_source"],
    cover = False,
    import_path = "google.golang.org/grpc/internal/credentials",
    visibility = ["PUBLIC"],
    deps = [":google.golang.org_grpc_grpclog"],
)


# google.golang.org/grpc/internal/envconfig

fileexport(
    name = "google.golang.org_grpc_internal_envconfig",
    srcs = ["internal/envconfig/envconfig.go"],
    tag = "go_source",
    deps = [":_google.golang.org_grpc#download"],
)

go_library(
    name = "google.golang.org_grpc_internal_envconfig",
    srcs = [":_google.golang.org_grpc_internal_envconfig#go_source"],
    cover = False,
    import_path = "google.golang.org/grpc/internal/envconfig",
    visibility = ["PUBLIC"],
    deps = [],
)


# google.golang.org/grpc/internal/grpclog

fileexport(
    name = "google.golang.org_grpc_internal_grpclog",
    srcs = [
        "internal/grpclog/grpclog.go",
        "internal/grpclog/prefixLogger.go",
    ],
    tag = "go_source",
    deps = [":_google.golang.org_grpc#download"],
)

go_library(
    name = "google.golang.org_grpc_internal_grpclog",
    srcs = [":_google.golang.org_grpc_internal_grpclog#go_source"],
    cover = False,
    import_path = "google.golang.org/grpc/internal/grpclog",
    visibility = ["PUBLIC"],
    deps = [],
)


# google.golang.org/grpc/internal/grpcrand

fileexport(
    name = "google.golang.org_grpc_internal_grpcrand",
    srcs = ["internal/grpcrand/grpcrand.go"],
    tag = "go_source",
    deps = [":_google.golang.org_grpc#download"],
)

go_library(
    name = "google.golang.org_grpc_internal_grpcrand",
    srcs = [":_google.golang.org_grpc_internal_grpcrand#go_source"],
    cover = False,
    import_path = "google.golang.org/grpc/internal/grpcrand",
    visibility = ["PUBLIC"],
    deps = [],
)


# google.golang.org/grpc/internal/grpcsync

fileexport(
    name = "google.golang.org_grpc_internal_grpcsync",
    srcs = ["internal/grpcsync/event.go"],
    tag = "go_source",
    deps = [":_google.golang.org_grpc#download"],
)

go_library(
    name = "google.golang.org_grpc_internal_grpcsync",
    srcs = [":_google.golang.org_grpc_internal_grpcsync#go_source"],
    cover = False,
    import_path = "google.golang.org/grpc/internal/grpcsync",
    visibility = ["PUBLIC"],
    deps = [],
)


# google.golang.org/grpc/internal/grpcutil

fileexport(
    name = "google.golang.org_grpc_internal_grpcutil",
    srcs = [
        "internal/grpcutil/encode_duration.go",
        "internal/grpcutil/metadata.go",
        "internal/grpcutil/method.go",
        "internal/grpcutil/target.go",
    ],
    tag = "go_source",
    deps = [":_google.golang.org_grpc#download"],
)

go_library(
    name = "google.golang.org_grpc_internal_grpcutil",
    srcs = [":_google.golang.org_grpc_internal_grpcutil#go_source"],
    cover = False,
    import_path = "google.golang.org/grpc/internal/grpcutil",
    visibility = ["PUBLIC"],
    deps = [
        ":google.golang.org_grpc_metadata",
        ":google.golang.org_grpc_resolver",
    ],
)


# google.golang.org/grpc/internal/resolver/dns

fileexport(
    name = "google.golang.org_grpc_internal_resolver_dns",
    srcs = [
        "internal/resolver/dns/dns_resolver.go",
        "internal/resolver/dns/go113.go",
    ],
    tag = "go_source",
    deps = [":_google.golang.org_grpc#download"],
)

go_library(
    name = "google.golang.org_grpc_internal_resolver_dns",
    srcs = [":_google.golang.org_grpc_internal_resolver_dns#go_source"],
    cover = False,
    import_path = "google.golang.org/grpc/internal/resolver/dns",
    visibility = ["PUBLIC"],
    deps = [
        ":google.golang.org_grpc_balancer_grpclb_state",
        ":google.golang.org_grpc_grpclog",
        ":google.golang.org_grpc_internal_envconfig",
        ":google.golang.org_grpc_internal_grpcrand",
        ":google.golang.org_grpc_resolver",
        ":google.golang.org_grpc_serviceconfig",
    ],
)


# google.golang.org/grpc/internal/resolver/passthrough

fileexport(
    name = "google.golang.org_grpc_internal_resolver_passthrough",
    srcs = ["internal/resolver/passthrough/passthrough.go"],
    tag = "go_source",
    deps = [":_google.golang.org_grpc#download"],
)

go_library(
    name = "google.golang.org_grpc_internal_resolver_passthrough",
    srcs = [":_google.golang.org_grpc_internal_resolver_passthrough#go_source"],
    cover = False,
    import_path = "google.golang.org/grpc/internal/resolver/passthrough",
    visibility = ["PUBLIC"],
    deps = [":google.golang.org_grpc_resolver"],
)


# google.golang.org/grpc/internal/serviceconfig

fileexport(
    name = "google.golang.org_grpc_internal_serviceconfig",
    srcs = ["internal/serviceconfig/serviceconfig.go"],
    tag = "go_source",
    deps = [":_google.golang.org_grpc#download"],
)

go_library(
    name = "google.golang.org_grpc_internal_serviceconfig",
    srcs = [":_google.golang.org_grpc_internal_serviceconfig#go_source"],
    cover = False,
    import_path = "google.golang.org/grpc/internal/serviceconfig",
    visibility = ["PUBLIC"],
    deps = [
        ":google.golang.org_grpc_balancer",
        ":google.golang.org_grpc_grpclog",
        ":google.golang.org_grpc_serviceconfig",
    ],
)


# google.golang.org/grpc/internal/status

fileexport(
    name = "google.golang.org_grpc_internal_status",
    srcs = ["internal/status/status.go"],
    tag = "go_source",
    deps = [":_google.golang.org_grpc#download"],
)

go_library(
    name = "google.golang.org_grpc_internal_status",
    srcs = [":_google.golang.org_grpc_internal_status#go_source"],
    cover = False,
    import_path = "google.golang.org/grpc/internal/status",
    visibility = ["PUBLIC"],
    deps = [
        ":github.com_golang_protobuf_proto",
        ":github.com_golang_protobuf_ptypes",
        ":google.golang.org_genproto_googleapis_rpc_status",
        ":google.golang.org_grpc_codes",
    ],
)


# google.golang.org/grpc/internal/syscall

fileexport(
    name = "google.golang.org_grpc_internal_syscall",
    srcs = select({
        ":__config_darwin_amd64": ["internal/syscall/syscall_nonlinux.go"],
        ":__config_linux_amd64": ["internal/syscall/syscall_linux.go"],
        "default": [],
    }),
    tag = "go_source",
    deps = [":_google.golang.org_grpc#download"],
)

go_library(
    name = "google.golang.org_grpc_internal_syscall",
    srcs = [":_google.golang.org_grpc_internal_syscall#go_source"],
    cover = False,
    import_path = "google.golang.org/grpc/internal/syscall",
    visibility = ["PUBLIC"],
    deps = [":google.golang.org_grpc_grpclog"] + select({
        ":__config_darwin_amd64": [],
        ":__config_linux_amd64": [":golang.org_x_sys_unix"],
        "default": [],
    }),
)


# google.golang.org/grpc/internal/transport

fileexport(
    name = "google.golang.org_grpc_internal_transport",
    srcs = [
        "internal/transport/bdp_estimator.go",
        "internal/transport/controlbuf.go",
        "internal/transport/defaults.go",
        "internal/transport/flowcontrol.go",
        "internal/transport/handler_server.go",
        "internal/transport/http2_client.go",
        "internal/transport/http2_server.go",
        "internal/transport/http_util.go",
        "internal/transport/transport.go",
    ],
    tag = "go_source",
    deps = [":_google.golang.org_grpc#download"],
)

go_library(
    name = "google.golang.org_grpc_internal_transport",
    srcs = [":_google.golang.org_grpc_internal_transport#go_source"],
    cover = False,
    import_path = "google.golang.org/grpc/internal/transport",
    visibility = ["PUBLIC"],
    deps = [
        ":github.com_golang_protobuf_proto",
        ":golang.org_x_net_http2",
        ":golang.org_x_net_http2_hpack",
        ":google.golang.org_genproto_googleapis_rpc_status",
        ":google.golang.org_grpc_codes",
        ":google.golang.org_grpc_credentials",
        ":google.golang.org_grpc_grpclog",
        ":google.golang.org_grpc_internal",
        ":google.golang.org_grpc_internal_channelz",
        ":google.golang.org_grpc_internal_grpcrand",
        ":google.golang.org_grpc_internal_grpcutil",
        ":google.golang.org_grpc_internal_syscall",
        ":google.golang.org_grpc_keepalive",
        ":google.golang.org_grpc_metadata",
        ":google.golang.org_grpc_peer",
        ":google.golang.org_grpc_resolver",
        ":google.golang.org_grpc_stats",
        ":google.golang.org_grpc_status",
        ":google.golang.org_grpc_tap",
    ],
)


# google.golang.org/grpc/keepalive

fileexport(
    name = "google.golang.org_grpc_keepalive",
    srcs = ["keepalive/keepalive.go"],
    tag = "go_source",
    deps = [":_google.golang.org_grpc#download"],
)

go_library(
    name = "google.golang.org_grpc_keepalive",
    srcs = [":_google.golang.org_grpc_keepalive#go_source"],
    cover = False,
    import_path = "google.golang.org/grpc/keepalive",
    visibility = ["PUBLIC"],
    deps = [],
)


# google.golang.org/grpc/metadata

fileexport(
    name = "google.golang.org_grpc_metadata",
    srcs = ["metadata/metadata.go"],
    tag = "go_source",
    deps = [":_google.golang.org_grpc#download"],
)

go_library(
    name = "google.golang.org_grpc_metadata",
    srcs = [":_google.golang.org_grpc_metadata#go_source"],
    cover = False,
    import_path = "google.golang.org/grpc/metadata",
    visibility = ["PUBLIC"],
    deps = [],
)


# google.golang.org/grpc/peer

fileexport(
    name = "google.golang.org_grpc_peer",
    srcs = ["peer/peer.go"],
    tag = "go_source",
    deps = [":_google.golang.org_grpc#download"],
)

go_library(
    name = "google.golang.org_grpc_peer",
    srcs = [":_google.golang.org_grpc_peer#go_source"],
    cover = False,
    import_path = "google.golang.org/grpc/peer",
    visibility = ["PUBLIC"],
    deps = [":google.golang.org_grpc_credentials"],
)


# google.golang.org/grpc/resolver

fileexport(
    name = "google.golang.org_grpc_resolver",
    srcs = ["resolver/resolver.go"],
    tag = "go_source",
    deps = [":_google.golang.org_grpc#download"],
)

go_library(
    name = "google.golang.org_grpc_resolver",
    srcs = [":_google.golang.org_grpc_resolver#go_source"],
    cover = False,
    import_path = "google.golang.org/grpc/resolver",
    visibility = ["PUBLIC"],
    deps = [
        ":google.golang.org_grpc_attributes",
        ":google.golang.org_grpc_credentials",
        ":google.golang.org_grpc_serviceconfig",
    ],
)


# google.golang.org/grpc/serviceconfig

fileexport(
    name = "google.golang.org_grpc_serviceconfig",
    srcs = ["serviceconfig/serviceconfig.go"],
    tag = "go_source",
    deps = [":_google.golang.org_grpc#download"],
)

go_library(
    name = "google.golang.org_grpc_serviceconfig",
    srcs = [":_google.golang.org_grpc_serviceconfig#go_source"],
    cover = False,
    import_path = "google.golang.org/grpc/serviceconfig",
    visibility = ["PUBLIC"],
    deps = [],
)


# google.golang.org/grpc/stats

fileexport(
    name = "google.golang.org_grpc_stats",
    srcs = [
        "stats/handlers.go",
        "stats/stats.go",
    ],
    tag = "go_source",
    deps = [":_google.golang.org_grpc#download"],
)

go_library(
    name = "google.golang.org_grpc_stats",
    srcs = [":_google.golang.org_grpc_stats#go_source"],
    cover = False,
    import_path = "google.golang.org/grpc/stats",
    visibility = ["PUBLIC"],
    deps = [":google.golang.org_grpc_metadata"],
)


# google.golang.org/grpc/status

fileexport(
    name = "google.golang.org_grpc_status",
    srcs = ["status/status.go"],
    tag = "go_source",
    deps = [":_google.golang.org_grpc#download"],
)

go_library(
    name = "google.golang.org_grpc_status",
    srcs = [":_google.golang.org_grpc_status#go_source"],
    cover = False,
    import_path = "google.golang.org/grpc/status",
    visibility = ["PUBLIC"],
    deps = [
        ":google.golang.org_genproto_googleapis_rpc_status",
        ":google.golang.org_grpc_codes",
        ":google.golang.org_grpc_internal_status",
    ],
)


# google.golang.org/grpc/tap

fileexport(
    name = "google.golang.org_grpc_tap",
    srcs = ["tap/tap.go"],
    tag = "go_source",
    deps = [":_google.golang.org_grpc#download"],
)

go_library(
    name = "google.golang.org_grpc_tap",
    srcs = [":_google.golang.org_grpc_tap#go_source"],
    cover = False,
    import_path = "google.golang.org/grpc/tap",
    visibility = ["PUBLIC"],
    deps = [],
)


