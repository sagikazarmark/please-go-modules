// Package sumfile implements a parser for go.sum files.
//
// The go.sum syntax is described in
// https://golang.org/cmd/go/#hdr-Module_authentication_using_go_sum
package sumfile

import (
	"io/ioutil"
	"os/exec"
	"sort"
	"strings"
)

// File is the parsed, interpreted form of a go.sum file.
type File struct {
	Modules []Module

	Errors []Error
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

// Error represents an error ocurred when parsing a specific entry in the sum file.
type Error struct {
	Pos int    // position of error (if present, line)
	Err string // the error itself
}

// Parse parses the data into a File struct.
func Parse(data []byte) File {
	lines := strings.Split(string(data), "\n")
	if lines[len(lines)-1] != "" {
		return File{
			Errors: []Error{
				{
					Pos: len(lines),
					Err: "final line missing newline",
				},
			},
		}
	}

	lines = lines[:len(lines)-1]

	sort.Strings(lines)

	var errs []Error

	var file File

	var currentModule Module
	var currentVersion Version

	for i, line := range lines {
		fields := strings.Fields(line)
		if len(fields) != 3 {
			errs = append(errs, Error{
				Pos: i + 1,
				Err: "invalid number of fields",
			})

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
	file.Errors = errs

	return file
}

// Load loads and parses a sum file from the current module.
func Load() (*File, error) {
	cmd := exec.Command("go", "env", "GOMOD")
	p, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	filePath := strings.TrimSuffix(strings.Trim(string(p), "\n"), ".mod") + ".sum"

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	file := Parse(data)

	return &file, nil
}
