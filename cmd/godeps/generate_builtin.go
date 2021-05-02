package main

import (
	"fmt"
	"sort"
	"strings"

	buildify "github.com/bazelbuild/buildtools/build"
	"github.com/scylladb/go-set/strset"

	"github.com/sagikazarmark/please-go-modules/pkg/depgraph"
)

func generateBuiltinBuildFiles(moduleList []depgraph.Module, ruleDir string, noExpand bool) (*buildify.File, bool, map[string]string) {
	file := newFile("", subinclude)
	var generateOsConfig bool
	knownDeps := make(map[string]string)

	packageToModule := map[string]string{}

	for _, module := range moduleList {
		for _, pkg := range module.Packages {
			packageToModule[pkg.ImportPath] = module.Path
		}
	}

	for _, module := range moduleList {
		name := sanitizeName(module.Path)

		if !noExpand {
			downloadModule := module.Path

			if module.Replace != "" {
				downloadModule = module.Replace
			}

			rule := &buildify.CallExpr{
				X: &buildify.Ident{Name: "go_mod_download"},
				List: []buildify.Expr{
					&buildify.AssignExpr{
						LHS: &buildify.Ident{Name: "name"},
						Op:  "=",
						RHS: &buildify.StringExpr{Value: name},
					},
					&buildify.AssignExpr{
						LHS: &buildify.Ident{Name: "_tag"},
						Op:  "=",
						RHS: &buildify.StringExpr{Value: "download"},
					},
					&buildify.AssignExpr{
						LHS: &buildify.Ident{Name: "module"},
						Op:  "=",
						RHS: &buildify.StringExpr{Value: downloadModule},
					},
					&buildify.AssignExpr{
						LHS: &buildify.Ident{Name: "version"},
						Op:  "=",
						RHS: &buildify.StringExpr{Value: module.Version},
					},
				},
			}

			file.Stmt = append(file.Stmt, rule)

			for _, pkg := range module.Packages {
				name := sanitizeName(pkg.ImportPath)
				downloadRule := ":_" + sanitizeName(module.Path) + "#download"

				depExpr := platformDepExpr("", pkg.Imports.Common, toPlatformSelectSet("", pkg.Imports.PerPlatform))
				if depExpr == nil {
					depExpr = &buildify.ListExpr{}
				}

				if !pkg.CgoFiles.Empty() {
					generateOsConfig = true
				}

				var install string
				if pkg.ImportPath == module.Path {
					install = "."
				} else {
					install = strings.TrimPrefix(pkg.ImportPath, module.Path+"/")
				}

				rule := &buildify.CallExpr{
					X: &buildify.Ident{Name: "go_module"},
					List: []buildify.Expr{
						&buildify.AssignExpr{
							LHS: &buildify.Ident{Name: "name"},
							Op:  "=",
							RHS: &buildify.StringExpr{Value: name},
						},
						&buildify.AssignExpr{
							LHS: &buildify.Ident{Name: "module"},
							Op:  "=",
							RHS: &buildify.StringExpr{Value: module.Path},
						},
						&buildify.AssignExpr{
							LHS: &buildify.Ident{Name: "download"},
							Op:  "=",
							RHS: &buildify.StringExpr{Value: downloadRule},
						},
						&buildify.AssignExpr{
							LHS: &buildify.Ident{Name: "install"},
							Op:  "=",
							RHS: &buildify.ListExpr{
								List: []buildify.Expr{
									&buildify.StringExpr{Value: install},
								},
							},
						},
						&buildify.AssignExpr{
							LHS: &buildify.Ident{Name: "deps"},
							Op:  "=",
							RHS: depExpr,
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
					},
				}

				var stmt buildify.Expr = rule

				if !pkg.AllPlatforms() {
					generateOsConfig = true

					var os, arch []string
					for _, p := range pkg.Platforms {
						os = append(os, p.OS)
						arch = append(arch, p.Arch)
					}

					os = uniqueStrings(os)
					arch = uniqueStrings(arch)

					stmt = &buildify.IfStmt{
						Cond: &buildify.CallExpr{
							X: &buildify.Ident{Name: "is_platform"},
							List: []buildify.Expr{
								&buildify.AssignExpr{
									LHS: &buildify.Ident{Name: "os"},
									Op:  "=",
									RHS: stringListExpr(os),
								},
								&buildify.AssignExpr{
									LHS: &buildify.Ident{Name: "arch"},
									Op:  "=",
									RHS: stringListExpr(arch),
								},
							},
						},
						True: []buildify.Expr{rule},
					}
				}

				file.Stmt = append(file.Stmt, stmt)

				if ruleDir != "" {
					knownDeps[pkg.ImportPath] = fmt.Sprintf("//%s:%s", ruleDir, name)
				} else {
					knownDeps[pkg.ImportPath] = fmt.Sprintf("//%s", name)
				}
			}
		} else {
			downloadModule := false

			if module.Replace != "" {
				downloadModule = true

				rule := &buildify.CallExpr{
					X: &buildify.Ident{Name: "go_mod_download"},
					List: []buildify.Expr{
						&buildify.AssignExpr{
							LHS: &buildify.Ident{Name: "name"},
							Op:  "=",
							RHS: &buildify.StringExpr{Value: name},
						},
						&buildify.AssignExpr{
							LHS: &buildify.Ident{Name: "_tag"},
							Op:  "=",
							RHS: &buildify.StringExpr{Value: "download"},
						},
						&buildify.AssignExpr{
							LHS: &buildify.Ident{Name: "module"},
							Op:  "=",
							RHS: &buildify.StringExpr{Value: module.Replace},
						},
						&buildify.AssignExpr{
							LHS: &buildify.Ident{Name: "version"},
							Op:  "=",
							RHS: &buildify.StringExpr{Value: module.Version},
						},
					},
				}

				file.Stmt = append(file.Stmt, rule)
			}

			rule := &buildify.CallExpr{
				X: &buildify.Ident{Name: "go_module"},
				List: []buildify.Expr{
					&buildify.AssignExpr{
						LHS: &buildify.Ident{Name: "name"},
						Op:  "=",
						RHS: &buildify.StringExpr{Value: name},
					},
					&buildify.AssignExpr{
						LHS: &buildify.Ident{Name: "module"},
						Op:  "=",
						RHS: &buildify.StringExpr{Value: module.Path},
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
				},
			}

			if downloadModule {
				rule.List = append(rule.List, &buildify.AssignExpr{
					LHS: &buildify.Ident{Name: "download"},
					Op:  "=",
					RHS: &buildify.StringExpr{Value: ":_" + name + "#download"},
				})
			} else {
				rule.List = append(rule.List, &buildify.AssignExpr{
					LHS: &buildify.Ident{Name: "version"},
					Op:  "=",
					RHS: &buildify.StringExpr{Value: module.Version},
				})
			}

			commonPkgsSet := strset.New()
			perPlatformPkgsSet := map[depgraph.Platform]*strset.Set{}
			commonDepsSet := strset.New()
			perPlatformDepsSet := map[depgraph.Platform]*strset.Set{}

			for _, pkg := range module.Packages {
				pkgName := strings.TrimPrefix(pkg.ImportPath, module.Path+"/")
				if pkg.ImportPath == module.Path {
					pkgName = "."
				}

				if !pkg.AllPlatforms() {
					generateOsConfig = true

					for _, platform := range pkg.Platforms {
						if perPlatformPkgsSet[platform] == nil {
							perPlatformPkgsSet[platform] = strset.New()
						}

						perPlatformPkgsSet[platform].Add(pkgName)
					}
				} else {
					commonPkgsSet.Add(pkgName)

					for _, importPath := range pkg.Imports.Common {
						if !module.BelongsTo(importPath) {
							commonDepsSet.Add(packageToModule[importPath])
						}
					}
				}

				for platform, imports := range pkg.Imports.PerPlatform {
					if perPlatformDepsSet[platform] == nil {
						perPlatformDepsSet[platform] = strset.New()
					}

					for _, importPath := range imports {
						if !module.BelongsTo(importPath) {
							perPlatformDepsSet[platform].Add(packageToModule[importPath])
						}
					}
				}

				if ruleDir != "" {
					knownDeps[pkg.ImportPath] = fmt.Sprintf("//%s:%s", ruleDir, name)
				} else {
					knownDeps[pkg.ImportPath] = fmt.Sprintf("//%s", name)
				}
			}

			commonPkgs := commonPkgsSet.List()
			sort.Strings(commonPkgs)
			perPlatformPkgs := map[depgraph.Platform][]string{}
			commonDeps := commonDepsSet.List()
			sort.Strings(commonDeps)
			perPlatformDeps := map[depgraph.Platform][]string{}

			for platform, set := range perPlatformPkgsSet {
				perPlatformPkgs[platform] = set.List()
				sort.Strings(perPlatformPkgs[platform])
			}

			for platform, set := range perPlatformDepsSet {
				perPlatformDeps[platform] = set.List()
				sort.Strings(perPlatformDeps[platform])
			}

			installExpr := platformExpr(commonPkgs, toPlatformSelectSet("", perPlatformPkgs), nil)
			if installExpr == nil {
				installExpr = &buildify.ListExpr{}
			}

			rule.List = append(rule.List, &buildify.AssignExpr{
				LHS: &buildify.Ident{Name: "install"},
				Op:  "=",
				RHS: installExpr,
			})

			depExpr := platformDepExpr("", commonDeps, toPlatformSelectSet("", perPlatformDeps))
			if depExpr == nil {
				depExpr = &buildify.ListExpr{}
			}

			rule.List = append(rule.List, &buildify.AssignExpr{
				LHS: &buildify.Ident{Name: "deps"},
				Op:  "=",
				RHS: depExpr,
			})

			file.Stmt = append(file.Stmt, rule)
		}
	}

	return file, generateOsConfig, knownDeps
}

func sanitizeName(name string) string {
	return strings.NewReplacer("/", "__").Replace(name)
}

func uniqueStrings(s []string) []string {
	tmp := make(map[string]bool)

	for _, v := range s {
		tmp[v] = true
	}

	ns := make([]string, 0, len(tmp))

	for k := range tmp {
		ns = append(ns, k)
	}

	return ns
}
