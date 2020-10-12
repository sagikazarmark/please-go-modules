package modgraph

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"

	"github.com/sagikazarmark/please-go-modules/pkg/golist"
	"github.com/sagikazarmark/please-go-modules/pkg/sumfile"
)

// Module represents a single module with all of it's packages and dependencies.
type Module struct {
	Module          *golist.Module
	Packages        []golist.Package
	Sum             string
	Deps            []string
	ResolvePackages bool
}

// CalculateDepGraph calculates the dependency graph of a
func CalculateDepGraph(module string, deps []golist.Package, sumFile sumfile.File) []Module {
	packages := make(map[string]golist.Package)
	modules := make(map[string]*Module)
	var moduleKeys []string
	pkgToModule := make(map[string]string)
	replaces := make(map[string]string)

	for _, pkg := range deps {
		if (pkg.Name == "main" && strings.HasSuffix(pkg.ImportPath, ".test")) || strings.HasSuffix(pkg.Name, "_test") {
			continue
		}

		if pkg.ForTest != "" {
			pkg.ImportPath = pkg.ForTest
		}

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

			if pkg.Module.Path == "cloud.google.com/go" || pkg.Module.Path == "google.golang.org/genproto" {
				moduleRoot.ResolvePackages = true
			}
		}

		moduleRoot.Packages = append(moduleRoot.Packages, pkg)
		pkgToModule[pkg.ImportPath] = pkg.Module.Path

		if pkg.Module.Replace != nil {
			replaces[pkg.Module.Replace.Path] = pkg.Module.Path
		}
	}

	for _, module := range sumFile.Modules {
		for _, version := range module.Versions {
			if version.Sum == "" {
				continue
			}

			modulePath := module.Name

			// Skip replaced module. Will be processed as the replaced module.
			if moduleRoot, ok := modules[modulePath]; ok && moduleRoot.Module.Replace != nil && moduleRoot.Module.Path != moduleRoot.Module.Replace.Path {
				continue
			}

			if replaced, ok := replaces[modulePath]; ok {
				modulePath = replaced
			}

			if moduleRoot, ok := modules[modulePath]; ok {
				moduleVersion := moduleRoot.Module.Version

				if moduleRoot.Module.Replace != nil && moduleRoot.Module.Replace.Version != "" {
					moduleVersion = moduleRoot.Module.Replace.Version
				}

				if moduleVersion != version.Version {
					continue
				}

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

					if imp == "C" {
						continue
					}

					// Ignore self-references
					if packages[imp].Module != nil && packages[imp].Module.Path == modulePath {
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

type Package struct {
	Path     string
	Dir      string
	Deps     []string
	TestDeps []string

	HasTests            bool
	HasIntegrationTests bool

	hasDep     map[string]bool
	hasTestDep map[string]bool
}

// CalculateInternalDepGraph calculates the dependency graph of internal dependencies, including test packages.
func CalculateInternalDepGraph(module string, deps []golist.Package) []*Package {
	var packages []*Package
	packageIndex := make(map[string]int)

	for _, pkg := range deps {
		if !strings.HasPrefix(pkg.ImportPath, module+"/") {
			continue
		}

		path := pkg.ImportPath

		if pkg.ForTest != "" {
			path = pkg.ForTest
		} else if pkg.Name == "main" && strings.HasSuffix(pkg.ImportPath, ".test") {
			path = strings.TrimSuffix(pkg.ImportPath, ".test")
		}

		i, ok := packageIndex[path]
		if !ok {
			i = len(packages)
			packages = append(packages, &Package{
				Path: path,
				Dir:  pkg.Dir,

				hasDep:     make(map[string]bool),
				hasTestDep: make(map[string]bool),
			})
			packageIndex[path] = i
		}

		_package := packages[i]

		for _, imp := range pkg.Imports {
			if !strings.HasPrefix(imp, module+"/") || strings.Contains(imp, fmt.Sprintf("[%s.test]", path)) {
				continue
			}

			if _, ok := _package.hasDep[imp]; ok {
				continue
			}

			if imp == path {
				continue
			}

			_package.Deps = append(_package.Deps, imp)
			_package.hasDep[imp] = true
		}

		sort.Strings(_package.Deps)

		for _, imp := range pkg.TestImports {
			if !strings.HasPrefix(imp, module+"/") || strings.Contains(imp, fmt.Sprintf("[%s.test]", path)) {
				continue
			}

			if _, ok := _package.hasTestDep[imp]; ok {
				continue
			}

			if imp == path {
				continue
			}

			_package.TestDeps = append(_package.TestDeps, imp)
			_package.hasTestDep[imp] = true
		}

		for _, imp := range pkg.XTestImports {
			if !strings.HasPrefix(imp, module+"/") || strings.Contains(imp, fmt.Sprintf("[%s.test]", path)) {
				continue
			}

			if _, ok := _package.hasTestDep[imp]; ok {
				continue
			}

			if imp == path {
				continue
			}

			_package.TestDeps = append(_package.TestDeps, imp)
			_package.hasTestDep[imp] = true
		}

		sort.Strings(_package.TestDeps)

		for _, file := range pkg.TestGoFiles {
			_package.HasTests = true

			c, _ := ioutil.ReadFile(filepath.Join(pkg.Dir, file))

			if strings.Contains(string(c), "func TestIntegration(t *testing.T) {") {
				_package.HasIntegrationTests = true

				break
			}
		}

		for _, file := range pkg.XTestGoFiles {
			_package.HasTests = true

			c, _ := ioutil.ReadFile(filepath.Join(pkg.Dir, file))

			if strings.Contains(string(c), "func TestIntegration(t *testing.T) {") {
				_package.HasIntegrationTests = true

				break
			}
		}
	}

	return packages
}
