package main

import (
	"flag"
	"fmt"
	"go/build"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	buildify "github.com/bazelbuild/buildtools/build"

	"github.com/sagikazarmark/please-go-modules/pkg/golist"
	"github.com/sagikazarmark/please-go-modules/pkg/modgraph"
	"github.com/sagikazarmark/please-go-modules/pkg/sumfile"
)

var (
	stdout              = flag.Bool("stdout", false, "Dump rules to the standard output")
	dir                 = flag.String("dir", "", "Dump rules into a directory")
	dryRun              = flag.Bool("dry-run", false, "Do not write anything to file")
	clean               = flag.Bool("clean", false, "Clean target before generating new rules")
	genpkg              = flag.Bool("genpkg", false, "Generate build targets for each package")
	subinclude          = flag.String("subinclude", "", "Include a rule in each file. (Useful when you don't want to duplicate the build definitions)")
	base                = flag.String("base", "", "Prepend this path to the directory")
	disableOptimization = flag.Bool("disable-optimization", false, "Disable build file optimization to debug code generator issues")
)

var supportedOSes = []string{"linux", "darwin"}

func selectOSes() ([]string, string) {
	cos := build.Default.GOOS

	var oses []string
	for _, sos := range supportedOSes {
		if cos == sos {
			continue
		}

		oses = append(oses, sos)
	}

	return oses, cos
}

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

	var ruleDir string
	if *dir != "" {
		ruleDir = path.Join(*base, *dir)
	}

	var generateOsConfig bool
	alternateOSes, currentOS := selectOSes()

	for _, module := range moduleList {
		if *genpkg {
			filePath := module.Module.Path

			// Get (and create) file
			file, ok := files[filePath]
			if !ok {
				if *subinclude != "" {
					file = fmt.Sprintf("subinclude(%q)\n\n", *subinclude)
				}

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
			visibility := fmt.Sprintf("\n    visibility = [\"//%s/...\"],\n", path.Join(ruleDir, module.Module.Path))

			if ruleDir == "" {
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

				// Get (or create) a file
				file, ok := files[filePath]
				if !ok {
					if *subinclude != "" {
						file = fmt.Sprintf("subinclude(%q)\n\n", *subinclude)
					}

					filePaths = append(filePaths, filePath)
				}

				name := path.Base(pkg.ImportPath)
				if ruleDir == "" {
					name = strings.Replace(pkg.ImportPath, "/", "_", -1)
				}

				var moduleSource, packageSource string

				// This isn't the root package, so we need to fetch the source
				if module.Module.Path != pkg.ImportPath {
					moduleSource = fmt.Sprintf("//%s:_%s#download", path.Join(ruleDir, module.Module.Path), path.Base(module.Module.Path))
					if ruleDir == "" {
						moduleSource = fmt.Sprintf(":_%s#download", strings.Replace(module.Module.Path, "/", "_", -1))
					}

					packageSource = strings.TrimPrefix(pkg.ImportPath, module.Module.Path+"/")
				} else {
					moduleSource = fmt.Sprintf(":_%s#download", name)
				}

				var gofiles []string
				for _, gf := range pkg.GoFiles {
					gofiles = append(gofiles, fmt.Sprintf("%q", path.Join(packageSource, gf)))
				}

				// Include ignored files (required for cross-compilation)
				for _, gf := range pkg.IgnoredGoFiles {
					if strings.HasSuffix(gf, "_test.go") {
						continue // Skip test files
					}

					gofiles = append(gofiles, fmt.Sprintf("%q", path.Join(packageSource, gf)))
				}

				file += fmt.Sprintf(`fileexport2(
    name = "%s",
    tag = "go_source",
    srcs = [%s],
    deps = ["%s"],
)`+"\n",
					name,
					strings.Join(gofiles, ", "),
					moduleSource,
				)

				isAsm := len(pkg.SFiles) > 0

				if isAsm {

					var sfiles []string
					for _, gf := range pkg.SFiles {
						sfiles = append(sfiles, fmt.Sprintf("%q", path.Join(packageSource, gf)))
					}

					file += fmt.Sprintf(`fileexport2(
    name = "%s",
    tag = "s_source",
    srcs = [%s],
    deps = ["%s"],
)`+"\n",
						name,
						strings.Join(sfiles, ", "),
						moduleSource,
					)
				}

				isCgo := len(pkg.CgoFiles) > 0

				if isCgo {
					generateOsConfig = true

					var cgofiles []string
					for _, gf := range pkg.CgoFiles {
						cgofiles = append(cgofiles, fmt.Sprintf("%q", path.Join(packageSource, gf)))
					}

					file += fmt.Sprintf(`fileexport2(
    name = "%s",
    tag = "cgo_source",
    srcs = [%s],
    deps = ["%s"],
)`+"\n",
						name,
						strings.Join(cgofiles, ", "),
						moduleSource,
					)

					var cfiles []string
					for _, gf := range pkg.CFiles {
						cfiles = append(cfiles, fmt.Sprintf("%q", path.Join(packageSource, gf)))
					}

					file += fmt.Sprintf(`fileexport2(
    name = "%s",
    tag = "c_source",
    srcs = [%s],
    deps = ["%s"],
)`+"\n",
						name,
						strings.Join(cfiles, ", "),
						moduleSource,
					)

					var hfiles []string
					for _, gf := range pkg.HFiles {
						hfiles = append(hfiles, fmt.Sprintf("%q", path.Join(packageSource, gf)))
					}

					file += fmt.Sprintf(`fileexport2(
    name = "%s",
    tag = "h_source",
    srcs = [%s],
    deps = ["%s"],
)`+"\n",
						name,
						strings.Join(hfiles, ", "),
						moduleSource,
					)
				}

				var deps []string
				for _, depPath := range pkg.Imports {
					if depMap[depPath].Standard {
						continue
					}

					if depPath == "C" {
						continue
					}

					if ruleDir == "" {
						deps = append(deps, fmt.Sprintf("%q", ":"+strings.Replace(depPath, "/", "_", -1)))

						continue
					}

					deps = append(deps, fmt.Sprintf("%q", fmt.Sprintf("//%s", path.Join(ruleDir, depPath))))
				}

				if isCgo {
					cgocflags := map[string][]string{
						currentOS: []string{},
					}

					for _, cf := range pkg.CgoCFLAGS {
						// Some libraries add . to the list of include paths
						// Replace it with PKG, ommit the rest
						// TODO: log ommited flags
						if strings.HasPrefix(cf, "-I") {
							// Module or import path reference
							if strings.Contains(cf, fmt.Sprintf("pkg/mod/%s", module.Module.Path)) ||
								strings.Contains(cf, pkg.ImportPath) {
								cgocflags[currentOS] = append(cgocflags[currentOS], `"-I ${PKG}"`)
							}

							continue
						}

						cgocflags[currentOS] = append(cgocflags[currentOS], fmt.Sprintf("%q", cf))
					}

					for _, os := range alternateOSes {
						ctxt := build.Default
						ctxt.GOOS = os
						ctxt.CgoEnabled = true

						osPkg, err := ctxt.ImportDir(pkg.Dir, build.ImportComment)
						if err != nil {
							panic(err)
						}

						cgocflags[os] = make([]string, 0, len(osPkg.CgoCFLAGS))
						for _, f := range osPkg.CgoCFLAGS {
							// Some libraries add . to the list of include paths
							// Replace it with PKG, ommit the rest
							// TODO: log ommited flags
							if strings.HasPrefix(f, "-I") {
								// Module or import path reference
								if strings.Contains(f, fmt.Sprintf("pkg/mod/%s", module.Module.Path)) ||
									strings.Contains(f, pkg.ImportPath) {
									cgocflags[os] = append(cgocflags[os], `"-I ${PKG}"`)
								}

								continue
							}

							cgocflags[os] = append(cgocflags[os], fmt.Sprintf("%q", f))
						}
					}

					var cgocflagsselect string
					for os, fs := range cgocflags {
						cgocflagsselect += fmt.Sprintf("\n"+`        "//%s:%s": [`+"\n", path.Join(ruleDir, "__config"), os)

						for _, f := range fs {
							cgocflagsselect += fmt.Sprintf(`            %s, `+"\n", f)
						}

						cgocflagsselect += "        ],\n    "
					}

					cgoldflags := map[string][]string{
						currentOS: []string{},
					}

					for _, cf := range pkg.CgoLDFLAGS {
						cgoldflags[currentOS] = append(cgoldflags[currentOS], fmt.Sprintf("%q", cf))
					}

					for _, os := range alternateOSes {
						ctxt := build.Default
						ctxt.GOOS = os
						ctxt.CgoEnabled = true

						osPkg, err := ctxt.ImportDir(pkg.Dir, build.ImportComment)
						if err != nil {
							panic(err)
						}

						cgoldflags[os] = make([]string, 0, len(osPkg.CgoLDFLAGS))
						for _, f := range osPkg.CgoLDFLAGS {
							cgoldflags[os] = append(cgoldflags[os], fmt.Sprintf("%q", f))
						}
					}

					var cgoldflagsselect string
					for os, fs := range cgoldflags {
						cgoldflagsselect += fmt.Sprintf("\n"+`        "//%s:%s": [`+"\n", path.Join(ruleDir, "__config"), os)

						for _, f := range fs {
							cgoldflagsselect += fmt.Sprintf(`            %s, `+"\n", f)
						}

						cgoldflagsselect += "        ],\n    "
					}

					file += fmt.Sprintf(`cgo_library(
    name = "%s",
    srcs = [":_%[1]s#cgo_source"],
    go_srcs = [":_%[1]s#go_source"],
    c_srcs = [":_%[1]s#c_source"],
    hdrs = [":_%[1]s#h_source"],
    compiler_flags = select({%s}),
    linker_flags = select({%s}),
    visibility = ["PUBLIC"],
    deps = [%s],
    import_path = "%s",
)`+"\n",
						name,
						cgocflagsselect,
						cgoldflagsselect,
						strings.Join(deps, ", "),
						pkg.ImportPath,
					)
				} else {
					var asm string
					if isAsm {
						asm = fmt.Sprintf("\n    asm_srcs = [\":_%s#s_source\"],\n", name)
					}

					file += fmt.Sprintf(`go_library(
    name = "%s",
    srcs = [":_%[1]s#go_source"],%s
	visibility = ["PUBLIC"],
	cover = False,
    deps = [%s],
    import_path = "%s",
)`+"\n",
						name,
						asm,
						strings.Join(deps, ", "),
						pkg.ImportPath,
					)
				}

				files[filePath] = file
			}
		} else if module.ResolvePackages { // END OF CURRENTLY WORKING LOGIC
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

	if generateOsConfig {
		filePath := "__config"

		// Get (or create) a file
		file, ok := files[filePath]
		if !ok {
			file = `package(default_visibility = ["PUBLIC"])` + "\n\n"

			filePaths = append(filePaths, filePath)
		}

		for _, os := range supportedOSes {
			file += fmt.Sprintf(`config_setting(name = "%[1]s", values = {"os": "%[1]s"})`+"\n", os)
		}
		files[filePath] = file
	}

	sort.Strings(filePaths)

	for filePath, fileContent := range files {
		buildFileContent := []byte(fileContent)
		if !*disableOptimization {
			buildFile, err := buildify.ParseBuild("BUILD.plz", buildFileContent)
			if err != nil {
				panic(err)
			}

			buildFileContent = buildify.Format(buildFile)
		}

		files[filePath] = string(buildFileContent)
	}

	if *stdout {
		for _, filePath := range filePaths {
			fmt.Printf("# %s\n\n%s\n\n", filePath, files[filePath])
		}
	} else if *dir != "" {
		if *dryRun {
			for _, filePath := range filePaths {
				file := files[filePath]

				fmt.Printf("%s:\n\n%s", path.Join(*dir, filePath, "BUILD.plz"), file)
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

				file, err := os.OpenFile(path.Join(dirPath, "BUILD.plz"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
				if err != nil {
					panic(err)
				}

				buildFileContent := []byte(files[filePath])

				if !*disableOptimization {
					buildFile, err := buildify.ParseBuild("BUILD.plz", []byte(files[filePath]))
					if err != nil {
						panic(err)
					}

					buildFileContent = buildify.Format(buildFile)
				}

				_, err = file.Write(buildFileContent)
				if err != nil {
					panic(err)
				}
			}
		}
	} else {
		log.Fatal("Either -stdout or -dir must be passed")
	}
}
