// Package sumfile implements a parser for go.sum files.
//
// The go.sum syntax is described in
// https://golang.org/cmd/go/#hdr-Module_authentication_using_go_sum
package sumfile

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"sort"
	"strings"
)

// File is the parsed, interpreted form of a go.sum file.
type File struct {
	Modules []Module
}

// Module is a single module in a go.sum file.
// It contains hashes for all versions in the file.
type Module struct {
	Name string

	Versions []Version
}

// Version contains the hash for both the package itself and for it's go.mod file.
type Version struct {
	Version  string
	Sum      string
	GoModSum string
}

// Parse parses the data into a File struct.
func Parse(data []byte) (*File, error) {
	lines := strings.Split(string(data), "\n")
	if lines[len(lines)-1] != "" {
		return nil, errors.New("final line missing newline")
	}

	lines = lines[:len(lines)-1]

	sort.Strings(lines)

	errs := make([]string, len(lines))

	var file File

	var currentModule Module
	var currentVersion Version

	for i, line := range lines {
		fields := strings.Fields(line)
		if len(fields) != 3 {
			errs[i] = "invalid number of fields"

			continue
		}

		name, version, sum := fields[0], fields[1], fields[2]

		var isGoMod bool

		if strings.HasSuffix(version, "/go.mod") {
			isGoMod = true
			version = strings.TrimSuffix(version, "/go.mod")
		}

		if name != currentModule.Name || version != currentVersion.Version {
			currentModule.Versions = append(currentModule.Versions, currentVersion)

			currentVersion = Version{
				Version: version,
			}
		}

		if name != currentModule.Name {
			if i != 0 {
				file.Modules = append(file.Modules, currentModule)
			}

			currentModule = Module{
				Name: name,
			}
		}

		if isGoMod {
			currentVersion.GoModSum = sum
		} else {
			currentVersion.Sum = sum
		}
	}

	currentModule.Versions = append(currentModule.Versions, currentVersion)
	file.Modules = append(file.Modules, currentModule)

	errStr := "invalid sum file:"
	var isErr bool

	for i, err := range errs {
		if err != "" {
			isErr = true
			errStr = fmt.Sprintf("%s\n%d: %s", errStr, i+1, err)
		}
	}

	if isErr {
		return nil, errors.New(errStr)
	}

	return &file, nil
}

// Load loads and parses a sum file from the current module.
func Load() (*File, error) {
	cmd := exec.Command("go", "env", "GOMOD")
	p, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	file := strings.TrimSuffix(strings.Trim(string(p), "\n"), ".mod") + ".sum"

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return Parse(data)
}
