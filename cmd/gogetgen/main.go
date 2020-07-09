package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/bazelbuild/buildtools/build"

	"github.com/sagikazarmark/please-go-modules/pkg/golist"
	"github.com/sagikazarmark/please-go-modules/pkg/modgraph"
	"github.com/sagikazarmark/please-go-modules/pkg/sumfile"
)

var (
	stdout = flag.Bool("stdout", false, "Dump rules to the standard output")
	dir    = flag.String("dir", "", "Dump rules into a directory")
	dryRun = flag.Bool("dry-run", false, "Do not write anything to file")
	clean  = flag.Bool("clean", false, "Clean target before generating new rules")
)

func main() {
	flag.Parse()

	rootModule, err := golist.CurrentModule()
	if err != nil {
		panic(err)
	}

	deps, err := golist.Deps(rootModule)
	if err != nil {
		panic(err)
	}

	sumFile, err := sumfile.Load()
	if err != nil {
		panic(err)
	}

	moduleList := modgraph.CalculateDepGraph(rootModule, deps, *sumFile)

	if *stdout {
		for _, module := range moduleList {
			var install []string

			if len(module.Packages) > 0 {
				for _, pkg := range module.Packages {
					install = append(install, fmt.Sprintf("%q", strings.TrimPrefix(strings.TrimPrefix(pkg.ImportPath, module.Module.Path), "/")))
				}
			}

			var deps []string
			for _, dep := range module.Deps {
				deps = append(deps, fmt.Sprintf("%q", ":"+strings.Replace(dep, "/", "_", -1)))
			}

			moduleVersion := module.Module.Version
			if module.Module.Replace != nil && module.Module.Replace.Version != "" {
				moduleVersion = module.Module.Replace.Version
			}

			var replace string

			if module.Module.Replace != nil {
				replace = fmt.Sprintf("\n    replace = \"%s\",\n", module.Module.Replace.Path)
			}

			fmt.Printf(`go_get(
    name = "%s",
    get = "%s",
    version = "%s",
    sum = "%s",%s
    install = [%s],
    deps = [%s],
    visibility=["PUBLIC"],
)`+"\n\n",
				strings.Replace(module.Module.Path, "/", "_", -1),
				module.Module.Path,
				moduleVersion,
				module.Sum,
				replace,
				strings.Join(install, ", "),
				strings.Join(deps, ", "),
			)
		}
	} else if *dir != "" {
		files := make(map[string]string)
		var filePaths []string

		for _, module := range moduleList {
			filePath := path.Dir(module.Module.Path)

			file, ok := files[filePath]
			if !ok {
				file = `package(default_visibility = ["PUBLIC"])` + "\n\n"
				filePaths = append(filePaths, filePath)
			}

			var install []string

			if len(module.Packages) > 0 {
				for _, pkg := range module.Packages {
					install = append(install, fmt.Sprintf("%q", strings.TrimPrefix(strings.TrimPrefix(pkg.ImportPath, module.Module.Path), "/")))
				}
			}

			var deps []string
			for _, dep := range module.Deps {
				if path.Dir(module.Module.Path) == path.Dir(dep) {
					deps = append(deps, fmt.Sprintf("%q", ":"+path.Base(dep)))
				} else {
					deps = append(deps, fmt.Sprintf("%q", fmt.Sprintf("//%s:%s", path.Join(*dir, path.Dir(dep)), path.Base(dep))))
				}
			}

			moduleVersion := module.Module.Version
			if module.Module.Replace != nil && module.Module.Replace.Version != "" {
				moduleVersion = module.Module.Replace.Version
			}

			var replace string

			if module.Module.Replace != nil {
				replace = fmt.Sprintf("\n    replace = \"%s\",\n", module.Module.Replace.Path)
			}

			file += fmt.Sprintf(`go_get(
    name = "%s",
    get = "%s",
    version = "%s",
    sum = "%s",%s
    install = [%s],
    deps = [%s],
)`+"\n",
				path.Base(module.Module.Path),
				module.Module.Path,
				moduleVersion,
				module.Sum,
				replace,
				strings.Join(install, ", "),
				strings.Join(deps, ", "),
			)

			files[filePath] = file
		}

		sort.Strings(filePaths)

		if *dryRun {
			for _, filePath := range filePaths {
				file := files[filePath]

				fmt.Printf("%s:\n\n%s", path.Join(*dir, filePath, "BUILD"), file)
			}
		} else {
			if filepath.IsAbs(*dir) {
				log.Fatal("Absolute path not allowed")
			}

			// TODO: disable every path outside of the module root

			if *clean {
				err := os.RemoveAll(*dir)
				if err != nil {
					panic(err)
				}
			}

			err := os.MkdirAll(*dir, 0755)
			if err != nil {
				panic(err)
			}

			for _, filePath := range filePaths {
				dirPath := path.Join(*dir, filePath)

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
	} else {
		log.Fatal("Either -stdout or -dir must be passed")
	}
}
