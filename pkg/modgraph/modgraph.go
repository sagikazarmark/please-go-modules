package modgraph

import (
	"sort"
	"strings"

	"github.com/sagikazarmark/please-go-modules/pkg/golist"
	"github.com/sagikazarmark/please-go-modules/pkg/sumfile"
)

// Module represents a single module with all of it's packages and dependencies.
type Module struct {
	Module   *golist.Module
	Packages []golist.Package
	Sum      string
	Deps     []string
}

// CalculateDepGraph calculates the dependency graph of a
func CalculateDepGraph(module string, deps []golist.Package, sumFile sumfile.File) []Module {
	packages := make(map[string]golist.Package)
	modules := make(map[string]*Module)
	var moduleKeys []string
	pkgToModule := make(map[string]string)

	for _, pkg := range deps {
		packages[pkg.ImportPath] = pkg

		// We don't want standard packages
		if pkg.Standard {
			continue
		}

		// We don't care about the root module for now
		if pkg.Module.Path == module {
			continue
		}

		// We don't care about submodules for now either
		if strings.HasPrefix(pkg.Module.Path, module+"/") {
			continue
		}

		moduleRoot, ok := modules[pkg.Module.Path]
		if !ok {
			moduleKeys = append(moduleKeys, pkg.Module.Path)
			moduleRoot = new(Module)

			moduleRoot.Module = pkg.Module
			modules[pkg.Module.Path] = moduleRoot
		}

		moduleRoot.Packages = append(moduleRoot.Packages, pkg)
		pkgToModule[pkg.ImportPath] = pkg.Module.Path
	}

	for _, module := range sumFile.Modules {
		for _, version := range module.Versions {
			if version.Sum == "" {
				continue
			}

			if moduleRoot, ok := modules[module.Name]; ok && moduleRoot.Module.Version == version.Version {
				moduleRoot.Sum = version.Sum
			}
		}
	}

	for _, modulePath := range moduleKeys {
		module := modules[modulePath]

		depList := make(map[string]bool)

		if len(module.Packages) > 0 {
			for _, pkg := range module.Packages {
				for _, imp := range pkg.Imports {
					if packages[imp].Standard {
						continue
					}

					depList[pkgToModule[imp]] = true
				}
			}
		}

		var deps []string
		for dep := range depList {
			deps = append(deps, dep)
		}

		module.Deps = deps
	}

	sort.Strings(moduleKeys)

	moduleList := make([]Module, 0, len(moduleKeys))

	for _, moduleKey := range moduleKeys {
		moduleList = append(moduleList, *modules[moduleKey])
	}

	return moduleList
}
