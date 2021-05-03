package depgraph

import (
	"sort"
	"strings"

	"github.com/sagikazarmark/please-go-modules/pkg/golist"
	"github.com/sagikazarmark/please-go-modules/pkg/sumfile"
	"github.com/scylladb/go-set/strset"
)

// GoPackageList is a list of Go packages for a specific platform.
type GoPackageList struct {
	Platform Platform
	Packages []golist.Package
}

// Platform identifies a specific target platform (ie. linux amd64).
type Platform struct {
	OS   string
	Arch string
}

func (p Platform) String() string {
	return p.OS + "_" + p.Arch
}

// Module is a Go module.
type Module struct {
	Path    string
	Replace string
	Version string
	Sum     string

	Packages []Package2

	pkgIndex map[string]bool
}

// SourcePath returns the module path that contains the source for this module.
// This is useful in case the module is replaced.
func (m Module) SourcePath() string {
	if m.Replace != "" {
		return m.Replace
	}

	return m.Path
}

// BelongsTo checks if an import path belongs to this module.
func (m Module) BelongsTo(importPath string) bool {
	_, ok := m.pkgIndex[importPath]

	return ok
}

// Package is a Go package with information for building the package on all supported platforms.
type Package struct {
	ImportPath string

	PlatformPackages map[Platform]golist.Package
}

// IsASM determines whether the package contains any assembly code.
func (p Package) IsASM() bool {
	for _, pkg := range p.PlatformPackages {
		if len(pkg.SFiles) > 0 {
			return true
		}
	}

	return false
}

// IsCGO determines whether the package contains any cgo code.
func (p Package) IsCGO() bool {
	for _, pkg := range p.PlatformPackages {
		if len(pkg.CgoFiles) > 0 {
			return true
		}
	}

	return false
}

// Package2 is a Go package with information for building the package on all supported platforms.
type Package2 struct {
	ImportPath string // import path of package in dir

	// Source files
	GoFiles  PlatformStringList // .go source files (excluding CgoFiles, TestGoFiles, XTestGoFiles)
	CgoFiles PlatformStringList // .go source files that import "C"
	CFiles   PlatformStringList // .c source files
	CXXFiles PlatformStringList // .cc, .cxx and .cpp source files
	HFiles   PlatformStringList // .h, .hh, .hpp and .hxx source files
	SFiles   PlatformStringList // .s source files

	// Cgo directives
	CgoCFLAGS   PlatformStringList // cgo: flags for C compiler
	CgoCPPFLAGS PlatformStringList // cgo: flags for C preprocessor
	CgoCXXFLAGS PlatformStringList // cgo: flags for C++ compiler
	CgoLDFLAGS  PlatformStringList // cgo: flags for linker

	// Dependency information
	Imports PlatformStringList // import paths used by this package

	// Platform information
	Platforms    []Platform
	allPlatforms bool

	Module Module
}

// IsASM determines whether the package contains any assembly code.
func (p Package2) IsASM() bool {
	if len(p.SFiles.Common) > 0 {
		return true
	}

	for _, files := range p.SFiles.PerPlatform {
		if len(files) > 0 {
			return true
		}
	}

	return false
}

// IsCGO determines whether the package contains any cgo code.
func (p Package2) IsCGO() bool {
	if len(p.CgoFiles.Common) > 0 {
		return true
	}

	for _, files := range p.CgoFiles.PerPlatform {
		if len(files) > 0 {
			return true
		}
	}

	return false
}

// AllPlatforms determines whether the package should be compiled on all platforms.
func (p Package2) AllPlatforms() bool {
	return p.allPlatforms
}

// PlatformStringList is a list of strings (ie. files, compiler flags, etc) for all supported platforms,
// divided into the intersection of all lists and the differences (for each platform) with said intersection.
type PlatformStringList struct {
	Common []string

	PerPlatform map[Platform][]string
}

// Empty checks whether the list has any items.
func (l PlatformStringList) Empty() bool {
	if len(l.Common) > 0 {
		return false
	}

	for _, set := range l.PerPlatform {
		if len(set) > 0 {
			return false
		}
	}

	return true
}

// CalculateDepGraph calculates the dependency graph of an application.
func CalculateDepGraph(rootModule string, packageLists []GoPackageList, sums sumfile.Index) []Module {
	allPackagesIdx := make(map[Platform]map[string]golist.Package)
	platformsIdx := make([]Platform, 0, len(packageLists))
	var packagesToProcess []string
	pkgToModule := make(map[string]string)

	modules := make(map[string]Module)
	var moduleKeys []string

	for _, packageList := range packageLists {
		allPackagesIdx[packageList.Platform] = make(map[string]golist.Package)
		platformsIdx = append(platformsIdx, packageList.Platform)

		for _, pkg := range packageList.Packages {
			if (pkg.Name == "main" && strings.HasSuffix(pkg.ImportPath, ".test")) || strings.HasSuffix(pkg.Name, "_test") {
				continue
			}

			if pkg.ForTest != "" {
				pkg.ImportPath = pkg.ForTest
			}

			// Record all packages at least once
			allPackagesIdx[packageList.Platform][pkg.ImportPath] = pkg

			// Filter unwanted packages
			if !packageFilter(rootModule, pkg) {
				continue
			}

			// Ensure the module is recorded
			module, ok := modules[pkg.Module.Path]
			if !ok {
				module = Module{
					Path:    pkg.Module.Path,
					Version: pkg.Module.Version,

					pkgIndex: map[string]bool{},
				}

				if pkg.Module.Replace != nil {
					module.Replace = pkg.Module.Replace.Path
					module.Version = pkg.Module.Replace.Version
				}

				module.Sum = sums.Sum(module.SourcePath(), module.Version)

				modules[module.Path] = module
				moduleKeys = append(moduleKeys, module.Path)
			}

			if _, ok := pkgToModule[pkg.ImportPath]; !ok {
				packagesToProcess = append(packagesToProcess, pkg.ImportPath)
				pkgToModule[pkg.ImportPath] = module.Path
			}
		}
	}

	sort.Strings(packagesToProcess)

	for _, packageToProcess := range packagesToProcess {
		platformVariants := make(map[Platform]golist.Package)

		allPlatforms := true
		var pkgPlatforms []Platform

		for _, platform := range platformsIdx {
			p, ok := allPackagesIdx[platform][packageToProcess]
			if !ok {
				platformVariants[platform] = golist.Package{}

				allPlatforms = false

				continue
			}

			platformVariants[platform] = p
			pkgPlatforms = append(pkgPlatforms, platform)
		}

		pkg := Package2{
			ImportPath: packageToProcess,

			GoFiles: calculatePlatformStringList(platformVariants, func(_ Platform, p golist.Package) []string { return p.GoFiles }),

			CgoFiles: calculatePlatformStringList(platformVariants, func(_ Platform, p golist.Package) []string { return p.CgoFiles }),
			CFiles:   calculatePlatformStringList(platformVariants, func(_ Platform, p golist.Package) []string { return p.CFiles }),

			CXXFiles: calculatePlatformStringList(platformVariants, func(_ Platform, p golist.Package) []string { return p.CXXFiles }),
			HFiles:   calculatePlatformStringList(platformVariants, func(_ Platform, p golist.Package) []string { return p.HFiles }),
			SFiles:   calculatePlatformStringList(platformVariants, func(_ Platform, p golist.Package) []string { return p.SFiles }),

			CgoCFLAGS:   calculatePlatformStringList(platformVariants, func(_ Platform, p golist.Package) []string { return p.CgoCFLAGS }),
			CgoCPPFLAGS: calculatePlatformStringList(platformVariants, func(_ Platform, p golist.Package) []string { return p.CgoCPPFLAGS }),
			CgoCXXFLAGS: calculatePlatformStringList(platformVariants, func(_ Platform, p golist.Package) []string { return p.CgoCXXFLAGS }),
			CgoLDFLAGS:  calculatePlatformStringList(platformVariants, func(_ Platform, p golist.Package) []string { return p.CgoLDFLAGS }),

			Imports: calculatePlatformStringList(platformVariants, func(platform Platform, pkg golist.Package) []string {
				imports := []string{}

				for _, i := range pkg.Imports {
					importedPkg, ok := allPackagesIdx[platform][i]
					if !ok {
						continue
					}

					if packageFilter(rootModule, importedPkg) {
						imports = append(imports, i)
					}
				}

				return imports
			}),

			Platforms:    pkgPlatforms,
			allPlatforms: allPlatforms,
		}

		module := modules[pkgToModule[packageToProcess]]

		module.Packages = append(module.Packages, pkg)
		module.pkgIndex[pkg.ImportPath] = true

		modules[pkgToModule[packageToProcess]] = module
	}

	sort.Strings(moduleKeys)

	moduleList := make([]Module, 0, len(moduleKeys))

	for _, moduleKey := range moduleKeys {
		moduleList = append(moduleList, modules[moduleKey])
	}

	return moduleList
}

func packageFilter(rootModule string, pkg golist.Package) bool {
	// We don't want standard packages
	if pkg.Standard {
		return false
	}

	// We don't care about the root module for now
	if pkg.Module.Path == rootModule {
		return false
	}

	// We don't care about submodules either
	if strings.HasPrefix(pkg.Module.Path, rootModule+"/") {
		return false
	}

	return true
}

func calculatePlatformStringList(pkgs map[Platform]golist.Package, ex func(Platform, golist.Package) []string) PlatformStringList {
	var sets []*strset.Set
	var setIndex []Platform
	diffSets := make(map[Platform][]string)

	for platform, p := range pkgs {
		sets = append(sets, strset.New(ex(platform, p)...))
		setIndex = append(setIndex, platform)
	}

	iset := strset.Intersection(sets...)

	for i, set := range sets {
		l := strset.Difference(set, iset).List()

		sort.Strings(l)

		diffSets[setIndex[i]] = l
	}

	c := iset.List()

	sort.Strings(c)

	return PlatformStringList{
		Common:      c,
		PerPlatform: diffSets,
	}
}

func calculateCgoPlatformStringList(pkgs map[Platform]golist.Package, ex func(Platform, golist.Package) []string) PlatformStringList {
	var sets []*strset.Set
	var setIndex []Platform
	diffSets := make(map[Platform][]string)

	for platform, p := range pkgs {
		sets = append(sets, strset.New(ex(platform, p)...))
		setIndex = append(setIndex, platform)
	}

	iset := strset.Union(sets...)

	for i, set := range sets {
		l := strset.Difference(set, iset).List()

		sort.Strings(l)

		diffSets[setIndex[i]] = l
	}

	c := iset.List()

	sort.Strings(c)

	return PlatformStringList{
		Common:      c,
		PerPlatform: diffSets,
	}
}
