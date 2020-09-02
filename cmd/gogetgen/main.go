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
	genpkg = flag.Bool("genpkg", false, "Generate build targets for each package")
)

func main() {
	flag.Parse()

	if *stdout && *dir != "" {
		panic("stdout and dir are mutually exclusive")
	}

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

	depMap := make(map[string]golist.Package)

	for _, dep := range deps {
		depMap[dep.ImportPath] = dep
	}

	modMap := make(map[string]modgraph.Module)

	for _, module := range moduleList {
		modMap[module.Module.Path] = module
	}

	files := make(map[string]string)
	var filePaths []string

	for _, module := range moduleList {
		if *genpkg {
			filePath := module.Module.Path

			// Get (and create) file
			file, ok := files[filePath]
			if !ok {
				filePaths = append(filePaths, filePath)
			}

			moduleVersion := module.Module.Version
			if module.Module.Replace != nil && module.Module.Replace.Version != "" {
				moduleVersion = module.Module.Replace.Version
			}

			var replace string

			if module.Module.Replace != nil {
				replace = fmt.Sprintf("\n    replace = \"%s\",\n", module.Module.Replace.Path)
			}

			name := path.Base(module.Module.Path)
			visibility := fmt.Sprintf("\n    visibility = [\"//%s/...\"],\n", path.Join(*dir, module.Module.Path))

			if *dir == "" {
				name = strings.Replace(module.Module.Path, "/", "_", -1)
				visibility = ""
			}

			file += fmt.Sprintf(`go_module_download(
    name = "%s",
    tag = "download",
    module = "%s",
    version = "%s",
    sum = "%s",%s%s
)`+"\n",
				name,
				module.Module.Path,
				moduleVersion,
				module.Sum,
				visibility,
				replace,
			)

			// Save file
			files[filePath] = file

			for _, pkg := range module.Packages {
				filePath := pkg.ImportPath

				// Get (and create) file
				file, ok := files[filePath]
				if !ok {
					filePaths = append(filePaths, filePath)
				}

				name := path.Base(pkg.ImportPath)
				if *dir == "" {
					name = strings.Replace(pkg.ImportPath, "/", "_", -1)
				}

				// This isn't the root package, so we need to fetch the source
				if module.Module.Path != pkg.ImportPath {
					moduleSource := fmt.Sprintf("//%s:_%s#download", path.Join(*dir, module.Module.Path), path.Base(module.Module.Path))
					if *dir == "" {
						moduleSource = fmt.Sprintf(":_%s#download", strings.Replace(module.Module.Path, "/", "_", -1))
					}

					packageSource := path.Join(*dir, module.Module.Path, "src", pkg.ImportPath)

					// TODO: need to pass specific go files for the package
					file += fmt.Sprintf(`build_rule(
	name = "%s",
	tag = "source",
    cmd = "",
    deps = ["%s"],
	output_dirs = ["%s"],
    output_is_complete = True,
)`+"\n",
						name,
						moduleSource,
						packageSource,
					)
				} else {
					// TODO: need to pass specific go files for the package
					file += fmt.Sprintf(`go_downloaded_source("%[1]s", ":_%[1]s#download")`+"\n", name)
				}

				var deps []string
				for _, depPath := range pkg.Imports {
					if depMap[depPath].Standard {
						continue
					}

					if depPath == "C" {
						continue
					}

					if *dir == "" {
						deps = append(deps, fmt.Sprintf("%q", ":"+strings.Replace(depPath, "/", "_", -1)))

						continue
					}

					deps = append(deps, fmt.Sprintf("%q", fmt.Sprintf("//%s", path.Join(*dir, depPath))))
				}

				file += fmt.Sprintf(`go_library(
    name = "%s",
    srcs = [":_%[1]s#source"],
    visibility = ["PUBLIC"],
    deps = [%s],
    _import_path = "%s",
)`+"\n",
					name,
					strings.Join(deps, ", "),
					pkg.ImportPath,
				)

				files[filePath] = file
			}
		} else if module.ResolvePackages {
			filePath := path.Dir(module.Module.Path)

			file, ok := files[filePath]
			if !ok {
				file = `package(default_visibility = ["PUBLIC"])` + "\n\n"
				filePaths = append(filePaths, filePath)
			}

			moduleVersion := module.Module.Version
			if module.Module.Replace != nil && module.Module.Replace.Version != "" {
				moduleVersion = module.Module.Replace.Version
			}

			var replace string

			if module.Module.Replace != nil {
				replace = fmt.Sprintf("\n    replace = \"%s\",\n", module.Module.Replace.Path)
			}

			var deps []string
			for i := path.Dir(module.Module.Path); strings.Contains(i, "/"); i = path.Dir(i) {
				if _, ok := modMap[i]; ok {
					deps = append(deps, fmt.Sprintf("%q", fmt.Sprintf(":_%s#download", strings.Replace(i, "/", "_", -1))))
				}
			}

			name := path.Base(module.Module.Path)
			if *dir == "" {
				name = strings.Replace(module.Module.Path, "/", "_", -1)
			}

			file += fmt.Sprintf(`go_module_download(
    name = "%s",
    tag = "download",
    module = "%s",
    version = "%s",
	sum = "%s",%s
	deps = [%s],
)`+"\n",
				name,
				module.Module.Path,
				moduleVersion,
				module.Sum,
				replace,
				strings.Join(deps, ", "),
			)

			files[filePath] = file

			for _, pkg := range module.Packages {
				filePath := path.Dir(pkg.ImportPath)

				file, ok := files[filePath]
				if !ok {
					file = `package(default_visibility = ["PUBLIC"])` + "\n\n"
					filePaths = append(filePaths, filePath)
				}

				install := fmt.Sprintf("%q", strings.TrimPrefix(strings.TrimPrefix(pkg.ImportPath, module.Module.Path), "/"))

				var deps []string
				for _, depPath := range pkg.Imports {
					if depMap[depPath].Standard {
						continue
					}

					if depPath == "C" {
						continue
					}

					if modMap[depMap[depPath].Module.Path].ResolvePackages {
						if path.Dir(pkg.ImportPath) == path.Dir(depPath) {
							deps = append(deps, fmt.Sprintf("%q", ":"+path.Base(depPath)))
						} else {
							deps = append(deps, fmt.Sprintf("%q", fmt.Sprintf("//%s:%s", path.Join(*dir, path.Dir(depPath)), path.Base(depPath))))
						}
					} else {
						if path.Dir(pkg.ImportPath) == path.Dir(depMap[depPath].Module.Path) {
							deps = append(deps, fmt.Sprintf("%q", ":"+path.Base(depMap[depPath].Module.Path)))
						} else if *dir == "" {
							deps = append(deps, fmt.Sprintf("%q", ":"+strings.Replace(depMap[depPath].Module.Path, "/", "_", -1)))
						} else {
							deps = append(deps, fmt.Sprintf("%q", fmt.Sprintf("//%s:%s", path.Join(*dir, path.Dir(depMap[depPath].Module.Path)), path.Base(depMap[depPath].Module.Path))))
						}
					}
				}

				name := path.Base(pkg.ImportPath)
				if *dir == "" {
					name = strings.Replace(pkg.ImportPath, "/", "_", -1)
				}

				file += fmt.Sprintf(`go_get(
    name = "%s",
    module = "%s",
	install = [%s],
	src = %s,
    deps = [%s],
)`+"\n",
					name,
					module.Module.Path,
					install,
					fmt.Sprintf("%q", fmt.Sprintf("//%s:_%s#download", path.Join(*dir, path.Dir(module.Module.Path)), path.Base(module.Module.Path))),
					strings.Join(deps, ", "),
				)

				files[filePath] = file
			}
		} else {
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

			var downloadDeps []string
			for i := path.Dir(module.Module.Path); strings.Contains(i, "/"); i = path.Dir(i) {
				if _, ok := modMap[i]; ok {
					downloadDeps = append(downloadDeps, fmt.Sprintf("%q", fmt.Sprintf(":_%s#download", strings.Replace(i, "/", "_", -1))))
				}
			}

			var deps []string
			for _, dep := range module.Deps {
				if modMap[dep].ResolvePackages {
					for _, pkg := range module.Packages {
						for _, imp := range pkg.Imports {
							if depMap[imp].Standard {
								continue
							}

							if imp == "C" {
								continue
							}

							if depMap[imp].Module.Path == dep {
								if path.Dir(module.Module.Path) == path.Dir(imp) {
									deps = append(deps, fmt.Sprintf("%q", ":"+path.Base(imp)))
								} else if *dir == "" {
									deps = append(deps, fmt.Sprintf("%q", ":"+strings.Replace(imp, "/", "_", -1)))
								} else {
									deps = append(deps, fmt.Sprintf("%q", fmt.Sprintf("//%s:%s", path.Join(*dir, path.Dir(imp)), path.Base(imp))))
								}
							}
						}
					}
				} else {
					if path.Dir(module.Module.Path) == path.Dir(dep) {
						deps = append(deps, fmt.Sprintf("%q", ":"+path.Base(dep)))
					} else if *dir == "" {
						deps = append(deps, fmt.Sprintf("%q", ":"+strings.Replace(dep, "/", "_", -1)))
					} else {
						deps = append(deps, fmt.Sprintf("%q", fmt.Sprintf("//%s:%s", path.Join(*dir, path.Dir(dep)), path.Base(dep))))
					}
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

			name := path.Base(module.Module.Path)
			if *dir == "" {
				name = strings.Replace(module.Module.Path, "/", "_", -1)
			}

			file += fmt.Sprintf(`go_get(
    name = "%s",
    module = "%s",
    version = "%s",
    sum = "%s",%s
    install = [%s],
	deps = [%s],
	download_deps = [%s],
)`+"\n",
				name,
				module.Module.Path,
				moduleVersion,
				module.Sum,
				replace,
				strings.Join(install, ", "),
				strings.Join(deps, ", "),
				strings.Join(downloadDeps, ", "),
			)

			files[filePath] = file
		}
	}

	sort.Strings(filePaths)

	for filePath, fileContent := range files {
		buildFile, err := build.ParseBuild("BUILD", []byte(fileContent))
		if err != nil {
			panic(err)
		}

		files[filePath] = string(build.Format(buildFile))
	}

	if *stdout {
		for _, filePath := range filePaths {
			fmt.Printf("# %s\n\n%s\n\n", filePath, files[filePath])
		}
	} else if *dir != "" {
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
