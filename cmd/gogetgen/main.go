package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type ModuleRoot struct {
	Module   *Module
	Packages []Package
	Sum      string
}

func main() {
	// TODO: find a better way to determine the root module name
	var rootModule string
	{
		cmd := exec.Command("go", "list", "-m")
		p, err := cmd.Output()
		if err != nil {
			panic(err)
		}

		rootModule = strings.Split(string(p), "\n")[0]
	}

	cmd := exec.Command("go", "list", "-deps", "-json", "./...")
	p, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	decoder := json.NewDecoder(bytes.NewReader(p))

	packages := make(map[string]Package)
	modules := make(map[string]*ModuleRoot)
	var moduleKeys []string
	pkgToModule := make(map[string]string)

	for err != io.EOF {
		var pkg Package

		err = decoder.Decode(&pkg)
		if err != nil {
			// TODO: handle error
			break
		}

		packages[pkg.ImportPath] = pkg

		// We don't want standard packages
		if pkg.Standard {
			continue
		}

		// We don't care about the root module for now
		if pkg.Module.Path == rootModule {
			continue
		}

		// We don't care about submodules for now either
		if strings.HasPrefix(pkg.Module.Path, rootModule+"/") {
			continue
		}

		moduleRoot, ok := modules[pkg.Module.Path]
		if !ok {
			moduleKeys = append(moduleKeys, pkg.Module.Path)
			moduleRoot = new(ModuleRoot)

			moduleRoot.Module = pkg.Module
			modules[pkg.Module.Path] = moduleRoot
		}

		moduleRoot.Packages = append(moduleRoot.Packages, pkg)
		pkgToModule[pkg.ImportPath] = pkg.Module.Path
	}

	sumFile, err := os.Open("go.sum")
	if err != nil {
		panic(err)
	}
	defer sumFile.Close()

	var lines []string

	scanner := bufio.NewScanner(sumFile)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "/go.mod h1") {
			continue
		}

		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	for _, line := range lines {
		parts := strings.Split(line, " ")
		path, version, sum := parts[0], parts[1], parts[2]

		moduleRoot, ok := modules[path]
		if !ok {
			continue
		}

		if moduleRoot.Module.Version != version {
			continue
		}

		moduleRoot.Sum = sum
	}

	for _, modulePath := range moduleKeys {
		module := modules[modulePath]

		var install []string
		depList := make(map[string]bool)

		if len(module.Packages) > 0 {
			for _, pkg := range module.Packages {
				install = append(install, fmt.Sprintf("%q", strings.TrimPrefix(strings.TrimPrefix(pkg.ImportPath, module.Module.Path), "/")))

				for _, imp := range pkg.Imports {
					if packages[imp].Standard {
						continue
					}

					depList[pkgToModule[imp]] = true
				}
			}
		}

		var deps []string
		for dep, _ := range depList {
			deps = append(deps, fmt.Sprintf("%q", ":"+strings.Replace(dep, "/", "_", -1)))
		}

		fmt.Printf(`go_module_get(
	name = "%s",
	get = "%s",
	version = "%s",
	sum = "%s",
	install = [%s],
	deps = [%s],
	visibility=["PUBLIC"],
)`+"\n\n",
			strings.Replace(module.Module.Path, "/", "_", -1),
			module.Module.Path,
			module.Module.Version,
			module.Sum,
			strings.Join(install, ", "),
			strings.Join(deps, ", "),
		)
	}
}
