package sumfile

// Index indexes module sums by module name and version.
type Index map[string]map[string]string

// Sum returns the module sum for a specific version of a module.
func (i Index) Sum(module string, version string) string {
	sums, ok := i[module]
	if !ok {
		return ""
	}

	return sums[version]
}

// CreateIndex creates a module sum index from a sum file.
func CreateIndex(file File) Index {
	index := make(Index, len(file.Modules))

	for _, module := range file.Modules {
		sums := make(map[string]string, len(module.Versions))

		for _, version := range module.Versions {
			sums[version.Version] = version.Sum
		}

		index[module.Name] = sums
	}

	return index
}
