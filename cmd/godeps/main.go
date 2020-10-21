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

	buildify "github.com/bazelbuild/buildtools/build"

	"github.com/sagikazarmark/please-go-modules/pkg/depgraph"
	"github.com/sagikazarmark/please-go-modules/pkg/golist"
	"github.com/sagikazarmark/please-go-modules/pkg/sumfile"
)

var (
	stdout     = flag.Bool("stdout", false, "Dump rules to the standard output")
	dir        = flag.String("dir", "", "Dump rules into a directory")
	dryRun     = flag.Bool("dry-run", false, "Do not write anything to file")
	clean      = flag.Bool("clean", false, "Clean target before generating new rules")
	subinclude = flag.String("subinclude", "", "Include a rule in each file. (Useful when you don't want to duplicate the build definitions)")
	base       = flag.String("base", "", "Prepend this path to the directory")
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

	deps := make([]depgraph.GoPackageList, len(SupportedPlatforms))

	for _, platform := range SupportedPlatforms {
		options := golist.ListOptions{
			Packages: []string{fmt.Sprintf("%s/...", rootModule)},
			Deps:     true,
			Test:     true,
			OS:       platform.OS,
			Arch:     platform.Arch,
		}

		platformDeps, err := golist.List(options)
		if err != nil {
			panic(err)
		}

		deps = append(deps, depgraph.GoPackageList{
			Platform: depgraph.Platform{
				OS:   platform.OS,
				Arch: platform.Arch,
			},
			Packages: platformDeps,
		})
	}

	sumFile, err := sumfile.Load()
	if err != nil {
		panic(err)
	}

	moduleList := depgraph.CalculateDepGraph(rootModule, deps, sumfile.CreateIndex(*sumFile))

	buildFiles := make(map[string]*buildify.File)
	var filePaths []string

	var ruleDir string
	if *dir != "" {
		ruleDir = path.Join(*base, *dir)
	}

	var generateOsConfig bool

	for _, module := range moduleList {
		filePath := module.Path

		// Get (and create) file
		file, ok := buildFiles[filePath]
		if !ok {
			file = newFile(filePath, subinclude)
			filePaths = append(filePaths, filePath)
			buildFiles[filePath] = file
		}

		name := path.Base(module.Path)
		visibility := fmt.Sprintf("//%s/...", path.Join(ruleDir, module.Path))

		if ruleDir == "" {
			name = strings.Replace(module.Path, "/", "_", -1)
			visibility = ""
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
					RHS: &buildify.StringExpr{Value: module.Path},
				},
				&buildify.AssignExpr{
					LHS: &buildify.Ident{Name: "version"},
					Op:  "=",
					RHS: &buildify.StringExpr{Value: module.Version},
				},
				&buildify.AssignExpr{
					LHS: &buildify.Ident{Name: "sum"},
					Op:  "=",
					RHS: &buildify.StringExpr{Value: module.Sum},
				},
			},
		}

		if module.Replace != "" {
			rule.List = append(rule.List, &buildify.AssignExpr{
				LHS: &buildify.Ident{Name: "replace"},
				Op:  "=",
				RHS: &buildify.StringExpr{Value: module.Replace},
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
			if module.Path != pkg.ImportPath {
				moduleSource = fmt.Sprintf("//%s:_%s#download", path.Join(ruleDir, module.Path), path.Base(module.Path))
				if ruleDir == "" {
					moduleSource = fmt.Sprintf(":_%s#download", strings.Replace(module.Path, "/", "_", -1))
				}

				packageSource = strings.TrimPrefix(pkg.ImportPath, module.Path+"/")
			} else {
				moduleSource = fmt.Sprintf(":_%s#download", name)
			}

			file.Stmt = append(file.Stmt, platformSourceFileRule(
				pkg.GoFiles.Common,
				toPlatformSelectSet(ruleDir, pkg.GoFiles.PerPlatform),
				name,
				"go_source",
				moduleSource,
				packageSource,
			))

			isAsm := pkg.IsASM()

			if isAsm {
				file.Stmt = append(file.Stmt, platformSourceFileRule(
					pkg.SFiles.Common,
					toPlatformSelectSet(ruleDir, pkg.SFiles.PerPlatform),
					name,
					"s_source",
					moduleSource,
					packageSource,
				))
			}

			isCgo := pkg.IsCGO()

			if isCgo {
				generateOsConfig = true

				file.Stmt = append(file.Stmt, platformSourceFileRule(
					pkg.CgoFiles.Common,
					toPlatformSelectSet(ruleDir, pkg.CgoFiles.PerPlatform),
					name,
					"cgo_source",
					moduleSource,
					packageSource,
				))

				file.Stmt = append(file.Stmt, platformSourceFileRule(
					pkg.CFiles.Common,
					toPlatformSelectSet(ruleDir, pkg.CFiles.PerPlatform),
					name,
					"c_source",
					moduleSource,
					packageSource,
				))

				file.Stmt = append(file.Stmt, platformSourceFileRule(
					pkg.HFiles.Common,
					toPlatformSelectSet(ruleDir, pkg.HFiles.PerPlatform),
					name,
					"h_source",
					moduleSource,
					packageSource,
				))
			}

			depExpr := platformDepExpr(ruleDir, pkg.Imports.Common, toPlatformSelectSet(ruleDir, pkg.Imports.PerPlatform))
			if depExpr == nil {
				depExpr = &buildify.ListExpr{}
			}

			if isCgo {
				cgoldflagsExpr := platformList(pkg.CgoLDFLAGS.Common, toPlatformSelectSet(ruleDir, pkg.CgoLDFLAGS.PerPlatform))
				if cgoldflagsExpr == nil {
					cgoldflagsExpr = &buildify.ListExpr{}
				}

				cgocFlagsExpr := platformCgocFlagsExpr(
					pkg.CgoCFLAGS.Common,
					toPlatformSelectSet(ruleDir, pkg.CgoCFLAGS.PerPlatform),
					pkg,
					module,
				)
				if cgocFlagsExpr == nil {
					cgocFlagsExpr = &buildify.ListExpr{}
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
							RHS: depExpr,
						},
						&buildify.AssignExpr{
							LHS: &buildify.Ident{Name: "compiler_flags"},
							Op:  "=",
							RHS: cgocFlagsExpr,
						},
						&buildify.AssignExpr{
							LHS: &buildify.Ident{Name: "linker_flags"},
							Op:  "=",
							RHS: cgoldflagsExpr,
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
							RHS: depExpr,
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

			if ruleDir != "" {
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
			}

			filePaths = append(filePaths, filePath)
			buildFiles[filePath] = file
		}

		for _, platform := range SupportedPlatforms {
			ruleName := platform.String()
			if ruleDir == "" {
				ruleName = "__config_" + ruleName
			}

			rule := &buildify.CallExpr{
				X: &buildify.Ident{Name: "config_setting"},
				List: []buildify.Expr{
					&buildify.AssignExpr{
						LHS: &buildify.Ident{Name: "name"},
						Op:  "=",
						RHS: &buildify.StringExpr{Value: ruleName},
					},
					&buildify.AssignExpr{
						LHS: &buildify.Ident{Name: "values"},
						Op:  "=",
						RHS: &buildify.DictExpr{
							List: []*buildify.KeyValueExpr{
								{
									Key:   &buildify.StringExpr{Value: "os"},
									Value: &buildify.StringExpr{Value: platform.OS},
								},
								{
									Key:   &buildify.StringExpr{Value: "cpu"},
									Value: &buildify.StringExpr{Value: platform.Arch},
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

func stringMapListSelect(s map[string][]string) buildify.Expr {
	keys := make([]string, 0, len(s))

	var hasValues bool

	for key, value := range s {
		keys = append(keys, key)

		if len(value) > 0 {
			hasValues = true
		}
	}

	if !hasValues {
		return nil
	}

	sort.Strings(keys)

	dict := &buildify.DictExpr{}

	for _, key := range keys {
		dict.List = append(dict.List, &buildify.KeyValueExpr{
			Key:   &buildify.StringExpr{Value: key},
			Value: stringListExpr(s[key]),
		})
	}

	dict.List = append(dict.List, &buildify.KeyValueExpr{
		Key:   &buildify.StringExpr{Value: "default"},
		Value: &buildify.ListExpr{},
	})

	return &buildify.CallExpr{
		X:    &buildify.Ident{Name: "select"},
		List: []buildify.Expr{dict},
	}
}

func platformSourceFileRule(
	commonFiles []string,
	platformFiles map[string][]string,
	name string,
	tag string,
	moduleSource string,
	packageSource string,
) *buildify.CallExpr {
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

	var commonFileList []string
	for _, file := range commonFiles {
		commonFileList = append(commonFileList, path.Join(packageSource, file))
	}

	platformFileList := make(map[string][]string, len(platformFiles))

	for platform, list := range platformFiles {
		newList := make([]string, 0, len(list))

		for _, file := range list {
			newList = append(newList, path.Join(packageSource, file))
		}

		platformFileList[platform] = newList
	}

	srcs := platformList(commonFileList, platformFileList)
	if srcs != nil {
		rule.List = append(rule.List, &buildify.AssignExpr{
			LHS: &buildify.Ident{Name: "srcs"},
			Op:  "=",
			RHS: srcs,
		})
	}

	return rule
}

func platformDepExpr(ruleDir string, common []string, platform map[string][]string) buildify.Expr {
	newCommon := make([]string, 0, len(common))

	for _, v := range common {
		newCommon = append(newCommon, formatDepPath(ruleDir, v))
	}

	newPlatform := make(map[string][]string, len(platform))

	for platform, list := range platform {
		newList := make([]string, 0, len(list))

		for _, v := range list {
			newList = append(newList, formatDepPath(ruleDir, v))
		}

		newPlatform[platform] = newList
	}

	return platformList(newCommon, newPlatform)
}

func formatDepPath(ruleDir string, depPath string) string {
	if ruleDir == "" {
		return ":" + strings.Replace(depPath, "/", "_", -1)
	}

	return fmt.Sprintf("//%s", path.Join(ruleDir, depPath))
}

func platformCgocFlagsExpr(common []string, platform map[string][]string, pkg depgraph.Package2, mod depgraph.Module) buildify.Expr {
	newCommon := filterCgoCFlags(common, pkg, mod)

	newPlatform := make(map[string][]string, len(platform))

	for platform, list := range platform {
		newPlatform[platform] = filterCgoCFlags(list, pkg, mod)
	}

	return platformList(newCommon, newPlatform)
}

func filterCgoCFlags(flags []string, pkg depgraph.Package2, mod depgraph.Module) []string {
	result := []string{}

	for _, f := range flags {
		// Some libraries add . to the list of include paths
		// Replace it with PKG, ommit the rest
		// TODO: log ommited flags
		if strings.HasPrefix(f, "-I") {
			// Module or import path reference
			if strings.Contains(f, fmt.Sprintf("pkg/mod/%s", mod.Path)) || strings.Contains(f, pkg.ImportPath) {
				result = append(result, "-I ${PKG}")
			}

			continue
		}

		result = append(result, f)
	}

	return result
}

func platformList(common []string, perPlatform map[string][]string) buildify.Expr {
	commonList := &buildify.ListExpr{}

	for _, s := range common {
		commonList.List = append(commonList.List, &buildify.StringExpr{Value: s})
	}

	platformSelect := stringMapListSelect(perPlatform)

	if len(commonList.List) > 0 && platformSelect != nil {
		return &buildify.BinaryExpr{
			X:  commonList,
			Op: "+",
			Y:  platformSelect,
		}
	} else if len(commonList.List) > 0 {
		return commonList
	} else if platformSelect != nil {
		return platformSelect
	}

	return nil
}

func toPlatformSelectSet(ruleDir string, sets map[depgraph.Platform][]string) map[string][]string {
	platformSets := make(map[string][]string, len(sets))

	for platform, set := range sets {
		if ruleDir == "" {
			platformSets[fmt.Sprintf(":__config_%s", platform.String())] = set

			continue
		}

		platformSets[fmt.Sprintf("//%s:%s", path.Join(ruleDir, "__config"), platform.String())] = set
	}

	return platformSets
}
