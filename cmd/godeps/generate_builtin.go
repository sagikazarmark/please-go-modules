package main

import (
	"fmt"
	"strings"

	buildify "github.com/bazelbuild/buildtools/build"

	"github.com/sagikazarmark/please-go-modules/pkg/depgraph"
)

func generateBuiltinBuildFiles(moduleList []depgraph.Module, ruleDir string) (*buildify.File, bool, map[string]string) {
	file := newFile("", subinclude)
	var generateOsConfig bool
	knownDeps := make(map[string]string)

	for _, module := range moduleList {
		name := sanitizeName(module.Path)

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
