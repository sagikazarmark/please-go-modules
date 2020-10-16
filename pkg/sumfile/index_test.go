package sumfile_test

import (
	"testing"

	"github.com/sagikazarmark/please-go-modules/pkg/sumfile"
)

func TestIndex(t *testing.T) {
	file := sumfile.File{
		Modules: []sumfile.Module{
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
		},
	}

	index := sumfile.CreateIndex(file)

	t.Run("OK", func(t *testing.T) {
		sum := index.Sum("logur.dev/logur", "v0.16.2")

		if got, want := sum, "h1:q4MxivaiTXiDHrQyeCH5WkwBLUrd6rM2lZlyztYvi4o="; got != want {
			t.Errorf("unexpected sum\nactual:   %q\nexpected: %q", got, want)
		}
	})

	t.Run("NoSum", func(t *testing.T) {
		sum := index.Sum("logur.dev/logur", "v0.16.1")

		if got, want := sum, ""; got != want {
			t.Errorf("unexpected sum\nactual:   %q\nexpected: %q", got, want)
		}
	})

	t.Run("UnknownVersion", func(t *testing.T) {
		sum := index.Sum("logur.dev/logur", "v99.99.99")

		if got, want := sum, ""; got != want {
			t.Errorf("unexpected sum\nactual:   %q\nexpected: %q", got, want)
		}
	})

	t.Run("UnknownModule", func(t *testing.T) {
		sum := index.Sum("example.com/nothing", "v99.99.99")

		if got, want := sum, ""; got != want {
			t.Errorf("unexpected sum\nactual:   %q\nexpected: %q", got, want)
		}
	})
}
