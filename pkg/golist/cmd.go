package golist

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// CurrentModule returns the name of the current module (regardless of the current directory).
func CurrentModule() (string, error) {
	cmd := exec.Command("go", "list", "-m")
	p, err := cmd.Output()
	if err != nil {
		return "", err
	}

	module := strings.Split(string(p), "\n")[0]

	if module == "" {
		return "", errors.New("failed to determine module name")
	}

	return module, nil
}

// Deps returns the list of dependencies.
func Deps(module string) ([]Package, error) {
	if module == "" {
		module = "."
	}

	cmd := exec.Command("go", "list", "-deps", "-json", fmt.Sprintf("%s/...", module))
	p, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return ParsePackages(p)
}
