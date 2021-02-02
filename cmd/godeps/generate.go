package main

import (
	"fmt"
	"path"
	"strings"

	buildify "github.com/bazelbuild/buildtools/build"

	"github.com/sagikazarmark/please-go-modules/pkg/depgraph"
)

func generateBuildFiles(moduleList []depgraph.Module, ruleDir string) (map[string]*buildify.File, []string, bool) {
	buildFiles := make(map[string]*buildify.File)
	var filePaths []string
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

			isAsm := !pkg.SFiles.Empty()

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

			isCgo := !pkg.CgoFiles.Empty()

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

				if !pkg.CFiles.Empty() {
					file.Stmt = append(file.Stmt, platformSourceFileRule(
						pkg.CFiles.Common,
						toPlatformSelectSet(ruleDir, pkg.CFiles.PerPlatform),
						name,
						"c_source",
						moduleSource,
						packageSource,
					))
				}

				if !pkg.HFiles.Empty() {
					file.Stmt = append(file.Stmt, platformSourceFileRule(
						pkg.HFiles.Common,
						toPlatformSelectSet(ruleDir, pkg.HFiles.PerPlatform),
						name,
						"h_source",
						moduleSource,
						packageSource,
					))
				}
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

				cgoSourceName := fmt.Sprintf(":_%s#cgo_source", name)
				var srcs buildify.Expr = &buildify.ListExpr{
					List: []buildify.Expr{
						&buildify.StringExpr{Value: cgoSourceName},
					},
				}

				// There are no common cgo files
				// In case only some of the platforms needs cgo,
				// we use a select to fall back to go_library in cgo_library
				// TODO: we only need a select if there are no common cgo files AND not all platforms have cgo files
				if len(pkg.CgoFiles.Common) == 0 {
					perPlatform := make(map[depgraph.Platform][]string, len(pkg.CgoFiles.PerPlatform))

					for platform, set := range pkg.CgoFiles.PerPlatform {
						if len(set) > 0 {
							perPlatform[platform] = []string{cgoSourceName}
						}
					}

					srcs = stringMapListSelect(toPlatformSelectSet(ruleDir, perPlatform))
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
							RHS: srcs,
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

				if !pkg.CFiles.Empty() {
					rule.List = append(rule.List, &buildify.AssignExpr{
						LHS: &buildify.Ident{Name: "c_srcs"},
						Op:  "=",
						RHS: &buildify.ListExpr{
							List: []buildify.Expr{
								&buildify.StringExpr{Value: fmt.Sprintf(":_%s#c_source", name)},
							},
						},
					})
				}

				if !pkg.HFiles.Empty() {
					rule.List = append(rule.List, &buildify.AssignExpr{
						LHS: &buildify.Ident{Name: "hdrs"},
						Op:  "=",
						RHS: &buildify.ListExpr{
							List: []buildify.Expr{
								&buildify.StringExpr{Value: fmt.Sprintf(":_%s#h_source", name)},
							},
						},
					})
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
					asmSourceName := fmt.Sprintf(":_%s#s_source", name)
					var srcs buildify.Expr = &buildify.ListExpr{
						List: []buildify.Expr{
							&buildify.StringExpr{Value: asmSourceName},
						},
					}

					// There are no common asm files
					// In case only some of the platforms needs asm
					// TODO: we only need a select if there are no common asm files AND not all platforms have asm files
					if len(pkg.SFiles.Common) == 0 {
						perPlatform := make(map[depgraph.Platform][]string, len(pkg.SFiles.PerPlatform))

						for platform, set := range pkg.SFiles.PerPlatform {
							if len(set) > 0 {
								perPlatform[platform] = []string{asmSourceName}
							}
						}

						srcs = stringMapListSelect(toPlatformSelectSet(ruleDir, perPlatform))
					}

					rule.List = append(rule.List, &buildify.AssignExpr{
						LHS: &buildify.Ident{Name: "asm_srcs"},
						Op:  "=",
						RHS: srcs,
					})
				}

				// There are no common go files
				// A rule might only be a dependency of a specific platform
				// TODO: we only need a select if there are no common go files AND not all platforms require the rule
				if len(pkg.GoFiles.Common) == 0 {
					perPlatform := make([]buildify.Expr, 0, len(pkg.CgoFiles.PerPlatform))

					for platform, set := range pkg.GoFiles.PerPlatform {
						if len(set) > 0 {
							perPlatform = append(perPlatform, &buildify.DictExpr{
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
							})
						}
					}

					ifRule := &buildify.IfStmt{
						Cond: &buildify.CallExpr{
							X: &buildify.Ident{Name: "select_config"},
							List: []buildify.Expr{
								&buildify.ListExpr{
									List: perPlatform,
								},
							},
						},
						True: []buildify.Expr{rule},
					}

					file.Stmt = append(file.Stmt, ifRule)
				} else {
					file.Stmt = append(file.Stmt, rule)
				}
			}
		}
	}

	return buildFiles, filePaths, generateOsConfig
}
