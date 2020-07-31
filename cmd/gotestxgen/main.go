package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/bazelbuild/buildtools/build"

	"github.com/sagikazarmark/please-go-modules/pkg/golist"
	"github.com/sagikazarmark/please-go-modules/pkg/modgraph"
)

var (
	dryRun     = flag.Bool("dry-run", false, "Do not write anything to file")
	ignorePath = flag.String("ignore-path", "", "Ignore paths starting with these strings")
)

func main() {
	flag.Parse()

	rootModule, err := golist.CurrentModule()
	if err != nil {
		panic(err)
	}

	deps, err := golist.DepsWithoutTests(rootModule)
	if err != nil {
		panic(err)
	}

	pkgList := modgraph.CalculateInternalDepGraph(rootModule, deps)

	ignorePaths := strings.Split(*ignorePath, "|")

	files := make(map[string]string)
	var filePaths []string

PkgLoop:
	for _, pkg := range pkgList {
		for _, ignorePath := range ignorePaths {
			p := fmt.Sprintf("%s/%s", rootModule, ignorePath)

			if pkg.Path == p || strings.HasPrefix(pkg.Path, p+"/") {
				continue PkgLoop
			}
		}

		file, ok := files[pkg.Dir]
		if !ok {
			file = `subinclude("///pleasings2//go:compat")

`
			filePaths = append(filePaths, pkg.Dir)
		}

		var deps []string
		for _, dep := range pkg.Deps {
			deps = append(deps, fmt.Sprintf("%q", "//"+strings.TrimPrefix(dep, rootModule+"/")))
		}

		var testDeps []string
		for _, dep := range pkg.TestDeps {
			testDeps = append(testDeps, fmt.Sprintf("%q", "//"+strings.TrimPrefix(dep, rootModule+"/")))
		}

		file += fmt.Sprintf(`filegroup(
    name = "%s",
    srcs = glob(["*.go"], exclude = ["*_test.go"]),
    deps = [%s],
    visibility = ["PUBLIC"],
)
`,
			path.Base(pkg.Path),
			strings.Join(deps, ", "),
		)

		if pkg.HasTests {
			file += fmt.Sprintf(`
go_test(
    name = "test",
    srcs = [":%s"],
    deps = [%s],
)
`,
				path.Base(pkg.Path),
				strings.Join(testDeps, ", "),
			)
		}

		if pkg.HasIntegrationTests {
			file += fmt.Sprintf(`
go_test(
    name = "integration_test",
    srcs = [":%s"],
    deps = [%s],
    flags = "-test.run ^TestIntegration$",
    labels = ["integration"],
)
`,
				path.Base(pkg.Path),
				strings.Join(testDeps, ", "),
			)
		}

		files[pkg.Dir] = file
	}

	sort.Strings(filePaths)

	if *dryRun {
		for _, filePath := range filePaths {
			file := files[filePath]

			fmt.Printf("%s:\n\n%s", path.Join(filePath, "BUILD"), file)
		}
	} else {
		for _, filePath := range filePaths {
			dirPath := path.Join(filePath)

			err := os.MkdirAll(dirPath, 0755)
			if err != nil {
				panic(err)
			}

			file, err := os.OpenFile(path.Join(dirPath, "BUILD"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
			if err != nil {
				panic(err)
			}

			buildFile, err := build.ParseBuild("BUILD", []byte(files[filePath]))
			if err != nil {
				panic(err)
			}

			_, err = file.Write(build.Format(buildFile))
			if err != nil {
				panic(err)
			}
		}
	}
}
