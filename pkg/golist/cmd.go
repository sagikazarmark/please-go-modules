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

	options := ListOptions{
		Packages: []string{fmt.Sprintf("%s/...", module)},
		Deps:     true,
		Test:     true,
	}

	return List(options)
}

// DepsWithoutTests returns the list of dependencies without tests.
func DepsWithoutTests(module string) ([]Package, error) {
	if module == "" {
		module = "."
	}

	options := ListOptions{
		Packages: []string{fmt.Sprintf("%s/...", module)},
		Deps:     true,
		Test:     false,
	}

	return List(options)
}

// ListOptions customizes a go list execution.
type ListOptions struct {
	Packages []string
	Deps     bool
	Test     bool
}

// List lists named packages.
func List(options ListOptions) ([]Package, error) {
	args := []string{"list", "-json"}

	if options.Deps {
		args = append(args, "-deps")
	}

	if options.Test {
		args = append(args, "-test")
	}

	args = append(args, options.Packages...)

	cmd := exec.Command("go", args...)
	p, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return ParsePackages(p)
}
