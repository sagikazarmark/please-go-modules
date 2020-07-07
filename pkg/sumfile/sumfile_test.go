package sumfile

import (
	"reflect"
	"testing"
)

const sum = `logur.dev/adapter/logrus v0.5.0 h1:cxsiceNXQLTKBk0keASgKAvrw9zzKa/XPE0Bn8tHXFI=
logur.dev/adapter/logrus v0.5.0/go.mod h1:9VKOXYYAQU3gjKJj1gs4jwr+YtDlGHGRVJ4tVAWeRhQ=
logur.dev/logur v0.16.1/go.mod h1:DyA5B+b6WjjCcnpE1+HGtTLh2lXooxRq+JmAwXMRK08=
logur.dev/logur v0.16.2 h1:q4MxivaiTXiDHrQyeCH5WkwBLUrd6rM2lZlyztYvi4o=
logur.dev/logur v0.16.2/go.mod h1:DyA5B+b6WjjCcnpE1+HGtTLh2lXooxRq+JmAwXMRK08=
`

func TestParse(t *testing.T) {
	file, err := Parse([]byte(sum))
	if err != nil {
		t.Fatal(err)
	}

	modules := []Module{
		{
			Name: "logur.dev/adapter/logrus",
			Versions: []Version{
				{
					Version:  "v0.5.0",
					Sum:      "h1:cxsiceNXQLTKBk0keASgKAvrw9zzKa/XPE0Bn8tHXFI=",
					GoModSum: "h1:9VKOXYYAQU3gjKJj1gs4jwr+YtDlGHGRVJ4tVAWeRhQ=",
				},
			},
		},
		{
			Name: "logur.dev/logur",
			Versions: []Version{
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
}
