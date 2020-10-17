package depgraph

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/sagikazarmark/please-go-modules/pkg/golist"
	"github.com/sagikazarmark/please-go-modules/pkg/sumfile"
)

func TestCalculateDepGraph(t *testing.T) {
	platforms := []Platform{
		{"linux", "amd64"},
		{"darwin", "amd64"},
	}

	var packageLists []GoPackageList

	for _, platform := range platforms {
		file, err := os.Open(fmt.Sprintf("testdata/packages/%s_%s.json", platform.OS, platform.Arch))
		if err != nil {
			t.Fatal(err)
		}

		decoder := json.NewDecoder(file)

		var packages []golist.Package

		for {
			var pkg golist.Package

			err := decoder.Decode(&pkg)
			if err == io.EOF {
				break
			} else if err != nil {
				t.Fatal(err)
			}

			packages = append(packages, pkg)
		}

		packageLists = append(packageLists, GoPackageList{platform, packages})
	}

	sumFileContent, err := ioutil.ReadFile("testdata/go.sum")
	if err != nil {
		t.Fatal(err)
	}

	sumFile := sumfile.Parse(sumFileContent)

	modules := CalculateDepGraph("github.com/sagikazarmark/please-go-modules/example4", packageLists, sumfile.CreateIndex(sumFile))

	t.Logf("%#v", modules)
}
