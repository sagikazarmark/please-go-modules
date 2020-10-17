package sumfile_test

import (
	"reflect"
	"testing"

	"github.com/sagikazarmark/please-go-modules/pkg/sumfile"
)

func TestParse(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		const sum = `logur.dev/adapter/logrus v0.5.0 h1:cxsiceNXQLTKBk0keASgKAvrw9zzKa/XPE0Bn8tHXFI=
logur.dev/adapter/logrus v0.5.0/go.mod h1:9VKOXYYAQU3gjKJj1gs4jwr+YtDlGHGRVJ4tVAWeRhQ=
logur.dev/logur v0.16.1/go.mod h1:DyA5B+b6WjjCcnpE1+HGtTLh2lXooxRq+JmAwXMRK08=
logur.dev/logur v0.16.2 h1:q4MxivaiTXiDHrQyeCH5WkwBLUrd6rM2lZlyztYvi4o=
logur.dev/logur v0.16.2/go.mod h1:DyA5B+b6WjjCcnpE1+HGtTLh2lXooxRq+JmAwXMRK08=
`

		file := sumfile.Parse([]byte(sum))
		if len(file.Errors) > 0 {
			t.Fatal(file.Errors)
		}

		modules := []sumfile.Module{
			{
				Name: "logur.dev/adapter/logrus",
				Versions: []sumfile.Version{
					{
						Version:  "v0.5.0",
						Sum:      "h1:cxsiceNXQLTKBk0keASgKAvrw9zzKa/XPE0Bn8tHXFI=",
						GoModSum: "h1:9VKOXYYAQU3gjKJj1gs4jwr+YtDlGHGRVJ4tVAWeRhQ=",
					},
				},
			},
			{
				Name: "logur.dev/logur",
				Versions: []sumfile.Version{
					{
						Version:  "v0.16.1",
						Sum:      "",
						GoModSum: "h1:DyA5B+b6WjjCcnpE1+HGtTLh2lXooxRq+JmAwXMRK08=",
					},
					{
						Version:  "v0.16.2",
						Sum:      "h1:q4MxivaiTXiDHrQyeCH5WkwBLUrd6rM2lZlyztYvi4o=",
						GoModSum: "h1:DyA5B+b6WjjCcnpE1+HGtTLh2lXooxRq+JmAwXMRK08=",
					},
				},
			},
		}

		// t.Logf("%+v", file.Modules)
		// t.Logf("%+v", modules)

		if !reflect.DeepEqual(file.Modules, modules) {
			t.Error("modules do not match")
		}
	})

	t.Run("MissinTrailingNewline", func(t *testing.T) {
		const sum = `logur.dev/adapter/logrus v0.5.0 h1:cxsiceNXQLTKBk0keASgKAvrw9zzKa/XPE0Bn8tHXFI=`

		file := sumfile.Parse([]byte(sum))

		errors := []sumfile.Error{
			{
				Pos: 1,
				Err: "final line missing newline",
			},
		}

		if !reflect.DeepEqual(file.Errors, errors) {
			t.Error("errors do not match")
		}
	})

	t.Run("InvalidNumberOfFields", func(t *testing.T) {
		const sum = `logur.dev/adapter/logrus v0.5.0 h1:cxsiceNXQLTKBk0keASgKAvrw9zzKa/XPE0Bn8tHXFI=
logur.dev/adapter/logrus v0.5.0/go.mod
logur.dev/logur
logur.dev/logur v0.16.2 h1:q4MxivaiTXiDHrQyeCH5WkwBLUrd6rM2lZlyztYvi4o= extra
`

		file := sumfile.Parse([]byte(sum))

		errors := []sumfile.Error{
			{
				Pos: 2,
				Err: "invalid number of fields",
			},
			{
				Pos: 3,
				Err: "invalid number of fields",
			},
			{
				Pos: 4,
				Err: "invalid number of fields",
			},
		}

		if !reflect.DeepEqual(file.Errors, errors) {
			t.Error("errors do not match")
		}
	})
}
