package golist

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
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
	Packages       []string
	Deps           bool
	Test           bool
	IgnoreNonFatal bool
	OS             string
	Arch           string
}

// GetOS returns the OS defined in the options,
// falling back to runtime.GOOS.
func (o ListOptions) GetOS() string {
	if o.OS == "" {
		return runtime.GOOS
	}

	return o.OS
}

// GetArch returns the arch defined in the options,
// falling back to runtime.GOARCH.
func (o ListOptions) GetArch() string {
	if o.Arch == "" {
		return runtime.GOARCH
	}

	return o.Arch
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

	if options.IgnoreNonFatal {
		args = append(args, "-e")
	}

	args = append(args, options.Packages...)

	cmd := exec.Command("go", args...)
	cmd.Env = os.Environ()

	cmd.Env = append(cmd.Env, "GOOS="+options.GetOS(), "GOARCH="+options.GetArch(), "CGO_ENABLED=1")

	p, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return ParsePackages(p)
}
