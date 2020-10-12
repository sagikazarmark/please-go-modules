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
	stdout     = flag.Bool("stdout", false, "Dump rules to the standard output")
	dir        = flag.String("dir", "", "Dump rules into a directory")
	dryRun     = flag.Bool("dry-run", false, "Do not write anything to file")
	clean      = flag.Bool("clean", false, "Clean target before generating new rules")
	genpkg     = flag.Bool("genpkg", false, "Generate build targets for each package")
	subinclude = flag.String("subinclude", "", "Include a rule in each file. (Useful when you don't want to duplicate the build definitions)")
	base       = flag.String("base", "", "Prepend this path to the directory")
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

	buildFiles := make(map[string]*buildify.File)
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
			file, ok := buildFiles[filePath]
			if !ok {
				file = newFile(filePath, subinclude)
				filePaths = append(filePaths, filePath)
				buildFiles[filePath] = file
			}

			name := path.Base(module.Module.Path)
			visibility := fmt.Sprintf("//%s/...", path.Join(ruleDir, module.Module.Path))

			if ruleDir == "" {
				name = strings.Replace(module.Module.Path, "/", "_", -1)
				visibility = ""
			}

			version := module.Module.Version
			if module.Module.Replace != nil && module.Module.Replace.Version != "" {
				version = module.Module.Replace.Version
			}

			var replace string

			if module.Module.Replace != nil {
				replace = module.Module.Replace.Path
			}

			rule := &buildify.CallExpr{
				X: &buildify.Ident{Name: "go_module_download"},
				List: []buildify.Expr{
					&buildify.AssignExpr{
						LHS: &buildify.Ident{Name: "name"},
						Op:  "=",
						RHS: &buildify.StringExpr{Value: name},
					},
					&buildify.AssignExpr{
						LHS: &buildify.Ident{Name: "tag"},
						Op:  "=",
						RHS: &buildify.StringExpr{Value: "download"},
					},
					&buildify.AssignExpr{
						LHS: &buildify.Ident{Name: "module"},
						Op:  "=",
						RHS: &buildify.StringExpr{Value: module.Module.Path},
					},
					&buildify.AssignExpr{
						LHS: &buildify.Ident{Name: "version"},
						Op:  "=",
						RHS: &buildify.StringExpr{Value: version},
					},
					&buildify.AssignExpr{
						LHS: &buildify.Ident{Name: "sum"},
						Op:  "=",
						RHS: &buildify.StringExpr{Value: module.Sum},
					},
				},
			}

			if replace != "" {
				rule.List = append(rule.List, &buildify.AssignExpr{
					LHS: &buildify.Ident{Name: "replace"},
					Op:  "=",
					RHS: &buildify.StringExpr{Value: replace},
				})
			}

			if visibility != "" {
				rule.List = append(rule.List, &buildify.AssignExpr{
					LHS: &buildify.Ident{Name: "visibility"},
					Op:  "=",
					RHS: &buildify.ListExpr{
						List: []buildify.Expr{
							&buildify.StringExpr{Value: visibility},
						},
					},
				})
			}

			file.Stmt = append(file.Stmt, rule)

			for _, pkg := range module.Packages {
				filePath := pkg.ImportPath

				// Get (and create) file
				file, ok := buildFiles[filePath]
				if !ok {
					file = newFile(filePath, subinclude)
					filePaths = append(filePaths, filePath)
					buildFiles[filePath] = file
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
				for _, f := range pkg.GoFiles {
					gofiles = append(gofiles, f)
				}

				// Include ignored files (required for cross-compilation)
				for _, f := range pkg.IgnoredGoFiles {
					if strings.HasSuffix(f, "_test.go") {
						continue // Skip test files
					}

					gofiles = append(gofiles, f)
				}

				file.Stmt = append(file.Stmt, sourceFileRule(gofiles, name, "go_source", moduleSource, packageSource))

				isAsm := len(pkg.SFiles) > 0

				if isAsm {
					file.Stmt = append(file.Stmt, sourceFileRule(pkg.SFiles, name, "s_source", moduleSource, packageSource))
				}

				isCgo := len(pkg.CgoFiles) > 0

				if isCgo {
					generateOsConfig = true

					file.Stmt = append(file.Stmt, sourceFileRule(pkg.CgoFiles, name, "cgo_source", moduleSource, packageSource))
					file.Stmt = append(file.Stmt, sourceFileRule(pkg.CFiles, name, "c_source", moduleSource, packageSource))
					file.Stmt = append(file.Stmt, sourceFileRule(pkg.HFiles, name, "h_source", moduleSource, packageSource))
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
						deps = append(deps, ":"+strings.Replace(depPath, "/", "_", -1))

						continue
					}

					deps = append(deps, fmt.Sprintf("//%s", path.Join(ruleDir, depPath)))
				}

				if isCgo {
					cgocflags := map[string][]string{
						fmt.Sprintf("//%s:%s", path.Join(ruleDir, "__config"), currentOS): filterCgoCFlags(pkg.CgoCFLAGS, pkg, module),
					}

					cgoldflags := map[string][]string{
						fmt.Sprintf("//%s:%s", path.Join(ruleDir, "__config"), currentOS): pkg.CgoLDFLAGS,
					}

					for _, os := range alternateOSes {
						ctxt := build.Default
						ctxt.GOOS = os
						ctxt.CgoEnabled = true

						osPkg, err := ctxt.ImportDir(pkg.Dir, build.ImportComment)
						if err != nil {
							panic(err)
						}

						cgocflags[fmt.Sprintf("//%s:%s", path.Join(ruleDir, "__config"), os)] = filterCgoCFlags(osPkg.CgoCFLAGS, pkg, module)
						cgoldflags[fmt.Sprintf("//%s:%s", path.Join(ruleDir, "__config"), os)] = osPkg.CgoLDFLAGS
					}

					rule := &buildify.CallExpr{
						X: &buildify.Ident{Name: "cgo_library"},
						List: []buildify.Expr{
							&buildify.AssignExpr{
								LHS: &buildify.Ident{Name: "name"},
								Op:  "=",
								RHS: &buildify.StringExpr{Value: name},
							},
							&buildify.AssignExpr{
								LHS: &buildify.Ident{Name: "srcs"},
								Op:  "=",
								RHS: &buildify.ListExpr{
									List: []buildify.Expr{
										&buildify.StringExpr{Value: fmt.Sprintf(":_%s#cgo_source", name)},
									},
								},
							},
							&buildify.AssignExpr{
								LHS: &buildify.Ident{Name: "go_srcs"},
								Op:  "=",
								RHS: &buildify.ListExpr{
									List: []buildify.Expr{
										&buildify.StringExpr{Value: fmt.Sprintf(":_%s#go_source", name)},
									},
								},
							},
							&buildify.AssignExpr{
								LHS: &buildify.Ident{Name: "c_srcs"},
								Op:  "=",
								RHS: &buildify.ListExpr{
									List: []buildify.Expr{
										&buildify.StringExpr{Value: fmt.Sprintf(":_%s#c_source", name)},
									},
								},
							},
							&buildify.AssignExpr{
								LHS: &buildify.Ident{Name: "hdrs"},
								Op:  "=",
								RHS: &buildify.ListExpr{
									List: []buildify.Expr{
										&buildify.StringExpr{Value: fmt.Sprintf(":_%s#h_source", name)},
									},
								},
							},
							&buildify.AssignExpr{
								LHS: &buildify.Ident{Name: "visibility"},
								Op:  "=",
								RHS: &buildify.ListExpr{
									List: []buildify.Expr{
										&buildify.StringExpr{Value: "PUBLIC"},
									},
								},
							},
							&buildify.AssignExpr{
								LHS: &buildify.Ident{Name: "import_path"},
								Op:  "=",
								RHS: &buildify.StringExpr{Value: pkg.ImportPath},
							},
							&buildify.AssignExpr{
								LHS: &buildify.Ident{Name: "deps"},
								Op:  "=",
								RHS: stringListExpr(deps),
							},
							&buildify.AssignExpr{
								LHS: &buildify.Ident{Name: "compiler_flags"},
								Op:  "=",
								RHS: stringMapListSelect(cgocflags),
							},
							&buildify.AssignExpr{
								LHS: &buildify.Ident{Name: "linker_flags"},
								Op:  "=",
								RHS: stringMapListSelect(cgoldflags),
							},
						},
					}

					file.Stmt = append(file.Stmt, rule)
				} else {
					rule := &buildify.CallExpr{
						X: &buildify.Ident{Name: "go_library"},
						List: []buildify.Expr{
							&buildify.AssignExpr{
								LHS: &buildify.Ident{Name: "name"},
								Op:  "=",
								RHS: &buildify.StringExpr{Value: name},
							},
							&buildify.AssignExpr{
								LHS: &buildify.Ident{Name: "srcs"},
								Op:  "=",
								RHS: &buildify.ListExpr{
									List: []buildify.Expr{
										&buildify.StringExpr{Value: fmt.Sprintf(":_%s#go_source", name)},
									},
								},
							},
							&buildify.AssignExpr{
								LHS: &buildify.Ident{Name: "visibility"},
								Op:  "=",
								RHS: &buildify.ListExpr{
									List: []buildify.Expr{
										&buildify.StringExpr{Value: "PUBLIC"},
									},
								},
							},
							&buildify.AssignExpr{
								LHS: &buildify.Ident{Name: "cover"},
								Op:  "=",
								RHS: &buildify.Ident{Name: "False"},
							},
							&buildify.AssignExpr{
								LHS: &buildify.Ident{Name: "import_path"},
								Op:  "=",
								RHS: &buildify.StringExpr{Value: pkg.ImportPath},
							},
							&buildify.AssignExpr{
								LHS: &buildify.Ident{Name: "deps"},
								Op:  "=",
								RHS: stringListExpr(deps),
							},
						},
					}

					if isAsm {
						rule.List = append(rule.List, &buildify.AssignExpr{
							LHS: &buildify.Ident{Name: "asm_srcs"},
							Op:  "=",
							RHS: &buildify.ListExpr{
								List: []buildify.Expr{
									&buildify.StringExpr{Value: fmt.Sprintf(":_%s#s_source", name)},
								},
							},
						})
					}

					file.Stmt = append(file.Stmt, rule)
				}
			}
		} else {
			panic("non-genpkg mode is unsopported ath the moment")
		}
	}

	if generateOsConfig {
		filePath := "__config"

		// Get (and create) file
		file, ok := buildFiles[filePath]
		if !ok {
			file = &buildify.File{
				Path: filePath,
				Type: buildify.TypeBuild,
			}

			file.Stmt = append(file.Stmt, &buildify.CallExpr{
				X: &buildify.Ident{Name: "package"},
				List: []buildify.Expr{
					&buildify.AssignExpr{
						LHS: &buildify.Ident{Name: "default_visibility"},
						Op:  "=",
						RHS: &buildify.ListExpr{
							List: []buildify.Expr{
								&buildify.StringExpr{Value: "PUBLIC"},
							},
						},
					},
				},
			})

			filePaths = append(filePaths, filePath)
			buildFiles[filePath] = file
		}

		for _, os := range supportedOSes {
			rule := &buildify.CallExpr{
				X: &buildify.Ident{Name: "config_setting"},
				List: []buildify.Expr{
					&buildify.AssignExpr{
						LHS: &buildify.Ident{Name: "name"},
						Op:  "=",
						RHS: &buildify.StringExpr{Value: os},
					},
					&buildify.AssignExpr{
						LHS: &buildify.Ident{Name: "values"},
						Op:  "=",
						RHS: &buildify.DictExpr{
							List: []*buildify.KeyValueExpr{
								{
									Key:   &buildify.StringExpr{Value: "os"},
									Value: &buildify.StringExpr{Value: os},
								},
							},
						},
					},
				},
			}

			file.Stmt = append(file.Stmt, rule)
		}
	}

	sort.Strings(filePaths)

	genFiles := make(map[string][]byte, len(buildFiles))

	for filePath, buildFile := range buildFiles {
		genFiles[filePath] = buildify.Format(buildFile)
	}

	if *stdout {
		for _, filePath := range filePaths {
			fmt.Printf("# %s\n\n%s\n\n", filePath, genFiles[filePath])
		}
	} else if *dir != "" {
		if *dryRun {
			for _, filePath := range filePaths {
				file := genFiles[filePath]

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

				buildFileContent := []byte(genFiles[filePath])

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

func newFile(filePath string, subinclude *string) *buildify.File {
	file := &buildify.File{
		Path: filePath,
		Type: buildify.TypeBuild,
	}

	if *subinclude != "" {
		file.Stmt = append(file.Stmt, &buildify.CallExpr{
			X: &buildify.Ident{Name: "subinclude"},
			List: []buildify.Expr{
				&buildify.StringExpr{Value: *subinclude},
			},
		})
	}

	return file
}

func stringListExpr(s []string) *buildify.ListExpr {
	list := &buildify.ListExpr{}

	for _, e := range s {
		list.List = append(list.List, &buildify.StringExpr{Value: e})
	}

	return list
}

func sourceFileRule(files []string, name string, tag string, moduleSource string, packageSource string) *buildify.CallExpr {
	rule := &buildify.CallExpr{
		X: &buildify.Ident{Name: "fileexport"},
		List: []buildify.Expr{
			&buildify.AssignExpr{
				LHS: &buildify.Ident{Name: "name"},
				Op:  "=",
				RHS: &buildify.StringExpr{Value: name},
			},
			&buildify.AssignExpr{
				LHS: &buildify.Ident{Name: "tag"},
				Op:  "=",
				RHS: &buildify.StringExpr{Value: tag},
			},
			&buildify.AssignExpr{
				LHS: &buildify.Ident{Name: "deps"},
				Op:  "=",
				RHS: &buildify.ListExpr{
					List: []buildify.Expr{
						&buildify.StringExpr{Value: moduleSource},
					},
				},
			},
		},
	}

	srcsList := &buildify.ListExpr{}

	for _, file := range files {
		srcsList.List = append(srcsList.List, &buildify.StringExpr{Value: path.Join(packageSource, file)})
	}

	rule.List = append(rule.List, &buildify.AssignExpr{
		LHS: &buildify.Ident{Name: "srcs"},
		Op:  "=",
		RHS: srcsList,
	})

	return rule
}

func filterCgoCFlags(flags []string, pkg golist.Package, mod modgraph.Module) []string {
	result := []string{}

	for _, f := range flags {
		// Some libraries add . to the list of include paths
		// Replace it with PKG, ommit the rest
		// TODO: log ommited flags
		if strings.HasPrefix(f, "-I") {
			// Module or import path reference
			if strings.Contains(f, fmt.Sprintf("pkg/mod/%s", mod.Module.Path)) || strings.Contains(f, pkg.ImportPath) {
				result = append(result, "-I ${PKG}")
			}

			continue
		}

		result = append(result, f)
	}

	return result
}

func stringMapListSelect(s map[string][]string) buildify.Expr {
	keys := make([]string, 0, len(s))

	for key := range s {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	dict := &buildify.DictExpr{}

	for _, key := range keys {
		dict.List = append(dict.List, &buildify.KeyValueExpr{
			Key:   &buildify.StringExpr{Value: key},
			Value: stringListExpr(s[key]),
		})
	}

	return &buildify.CallExpr{
		X:    &buildify.Ident{Name: "select"},
		List: []buildify.Expr{dict},
	}
}
