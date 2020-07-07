package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/sagikazarmark/please-go-modules/pkg/golist"
	"github.com/sagikazarmark/please-go-modules/pkg/modgraph"
	"github.com/sagikazarmark/please-go-modules/pkg/sumfile"
)

var (
	stdout = flag.Bool("stdout", false, "Dump rules to the standard output")
	dir    = flag.String("dir", "", "Dump rules into a directory")
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

			fmt.Printf(`go_module_get(
	name = "%s",
	get = "%s",
	version = "%s",
	sum = "%s",
	install = [%s],
	deps = [%s],
	visibility=["PUBLIC"],
)`+"\n\n",
				strings.Replace(module.Module.Path, "/", "_", -1),
				module.Module.Path,
				module.Module.Version,
				module.Sum,
				strings.Join(install, ", "),
				strings.Join(deps, ", "),
			)
		}
	} else if *dir != "" {
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

		for _, module := range moduleList {
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

			fmt.Printf(`%s

go_module_get(
	name = "%s",
	get = "%s",
	version = "%s",
	sum = "%s",
	install = [%s],
	deps = [%s],
	visibility=["PUBLIC"],
)`+"\n\n",
				path.Join(*dir, path.Dir(module.Module.Path), "BUILD"),
				path.Base(module.Module.Path),
				module.Module.Path,
				module.Module.Version,
				module.Sum,
				strings.Join(install, ", "),
				strings.Join(deps, ", "),
			)
		}
	} else {
		log.Fatal("Either -stdout or -dir must be passed")
	}
}
